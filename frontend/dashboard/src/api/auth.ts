import {
	QueryClient,
	QueryOptions,
	useMutation,
	useQuery,
	useQueryClient,
} from '@tanstack/vue-query';
import type { Dashboard, Profile } from '@twir/api/messages/auth/auth';
import { computed } from 'vue';

import { protectedApiClient } from './twirp.js';

export const profileQueryOptions = {
	queryKey: ['authProfile'],
	queryFn: async () => {
		const call = await protectedApiClient.authUserProfile({});
		return call.response;
	},
	retry: false,
} satisfies QueryOptions;

export const useProfile = () => useQuery<Profile | null>(profileQueryOptions);


export const useLogout = () => useMutation({
	mutationKey: ['authLogout'],
	mutationFn: async () => {
		await protectedApiClient.authLogout({});
	},
});


export const dashboardsQueryOptions = {
	queryKey: ['authDashboards'],
	queryFn: async () => {
		const call = await protectedApiClient.authGetDashboards({});
		return call.response;
	},
} satisfies QueryOptions;

export const useDashboards = () => useQuery(dashboardsQueryOptions);

export const useSetDashboard = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationKey: ['authSetDashboard'],
		mutationFn: async (dashboardId: string) => {
			await protectedApiClient.authSetDashboard({ dashboardId });
		},
		onSuccess: async () => {
			await queryClient.invalidateQueries();
		},
	});
};


export const PERMISSIONS_FLAGS = {
	CAN_ACCESS_DASHBOARD: 'All permissions',
	empty1: '',

	UPDATE_CHANNEL_TITLE: 'Can update channel title',
	UPDATE_CHANNEL_CATEGORY: 'Can update channel category',

	VIEW_COMMANDS: 'Can view commands',
	MANAGE_COMMANDS: 'Can manage commands',

	VIEW_KEYWORDS: 'Can view keywords',
	MANAGE_KEYWORDS: 'Can manage keywords',

	VIEW_TIMERS: 'Can view timers',
	MANAGE_TIMERS: 'Can manage timers',

	VIEW_INTEGRATIONS: 'Can view integrations',
	MANAGE_INTEGRATIONS: 'Can manage integrations',

	VIEW_SONG_REQUESTS: 'Can view song requests',
	MANAGE_SONG_REQUESTS: 'Can manage song requests',

	VIEW_MODERATION: 'Can view moderation settings',
	MANAGE_MODERATION: 'Can manage moderation settings',

	VIEW_VARIABLES: 'Can view variables',
	MANAGE_VARIABLES: 'Can manage variables',

	VIEW_GREETINGS: 'Can view greetings',
	MANAGE_GREETINGS: 'Can manage greetings',

	VIEW_OVERLAYS: 'Can view overlays',
	MANAGE_OVERLAYS: 'Can manage overlays',

	VIEW_ROLES: 'Can view roles',
	MANAGE_ROLES: 'Can manage roles',

	VIEW_EVENTS: 'Can view events',
	MANAGE_EVENTS: 'Can manage events',

	VIEW_ALERTS: 'Can view alerts',
	MANAGE_ALERTS: 'Can manage alerts',

	VIEW_GAMES: 'Can view games',
	MANAGE_GAMES: 'Can manage games',

	VIEW_GIVEAWAYS: 'Can view giveaways',
	MANAGE_GIVEAWAYS: 'Can manage giveaways',
};

export type PermissionsType = keyof typeof PERMISSIONS_FLAGS

export const useUserAccessFlagChecker = (flag: PermissionsType) => {
	const profile = useProfile();
	const dashboards = useDashboards();

	return computed(() => {
		if (!dashboards.data.value?.dashboards || !profile.data.value?.selectedDashboardId) return false;

		if (profile.data.value.id == profile.data.value.selectedDashboardId) {
			return true;
		}

		const dashboard = dashboards.data.value.dashboards.find(d => d.id === profile.data.value!.selectedDashboardId);
		if (!dashboard) return false;

		if (dashboard.flags.includes('CAN_ACCESS_DASHBOARD')) return true;
		return dashboard.flags.includes(flag);
	});
};

export const userAccessFlagChecker = async (queryClient: QueryClient, flag: PermissionsType) => {
	const profile = await queryClient.getQueryData(profileQueryOptions.queryKey) as Profile | null;
	const { dashboards } = await queryClient.getQueryData(dashboardsQueryOptions.queryKey) as {
		dashboards: Dashboard[]
	};

	if (!dashboards || !profile || !profile.selectedDashboardId) return false;
	if (profile.selectedDashboardId == profile.id) return true;

	const dashboard = dashboards.find(d => d.id === profile.selectedDashboardId);
	if (!dashboard) return false;

	if (dashboard.flags.includes('CAN_ACCESS_DASHBOARD')) return true;
	return dashboard.flags.includes(flag);
};
