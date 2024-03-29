import { QueryClient } from '@tanstack/vue-query';
import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';

import {
	userAccessFlagChecker,
	profileQueryOptions,
	dashboardsQueryOptions,
} from '@/api';

export const newRouter = (queryClient: QueryClient) => {
	const routes: ReadonlyArray<RouteRecordRaw> = [
		{
			path: '/dashboard/integrations/:integrationName',
			component: () => import('../pages/IntegrationsCallback.vue'),
		},
		{
			path: '/dashboard',
			component: () => import('../layout/layout.vue'),
			children: [
				{
					path: '/dashboard',
					component: () => import('../pages/Dashboard.vue'),
					meta: {
						noPadding: true,
					},
				},
				{
					name: 'Integrations',
					path: '/dashboard/integrations',
					component: () => import('../pages/Integrations.vue'),
					meta: { neededPermission: 'VIEW_INTEGRATIONS' },
				},
				{
					path: '/dashboard/commands/:system',
					component: () => import('../pages/Commands.vue'),
					meta: { neededPermission: 'VIEW_COMMANDS' },
				},
				{
					path: '/dashboard/timers',
					component: () => import('../pages/Timers.vue'),
					meta: { neededPermission: 'VIEW_TIMERS' },
				},
				{
					path: '/dashboard/keywords',
					component: () => import('../pages/Keywords.vue'),
					meta: { neededPermission: 'VIEW_KEYWORDS' },
				},
				{
					path: '/dashboard/variables',
					component: () => import('../pages/Variables.vue'),
					meta: { neededPermission: 'VIEW_VARIABLES' },
				},
				{
					path: '/dashboard/greetings',
					component: () => import('../pages/Greetings.vue'),
					meta: { neededPermission: 'VIEW_GREETINGS' },
				},
				{
					path: '/dashboard/community/users',
					component: () => import('../pages/CommunityUsers.vue'),
				},
				{
					path: '/dashboard/community/roles',
					component: () => import('../pages/CommunityRoles.vue'),
					meta: { neededPermission: 'VIEW_ROLES' },
				},
				{
					path: '/dashboard/song-requests',
					component: () => import('../pages/SongRequests.vue'),
					meta: { neededPermission: 'VIEW_SONG_REQUESTS' },
				},
				{
					path: '/dashboard/overlays',
					component: () => import('../pages/Overlays.vue'),
					meta: { neededPermission: 'VIEW_OVERLAYS' },
				},
				{
					name: 'ChatOverlay',
					path: '/dashboard/overlays/chat',
					component: () => import('../pages/overlays/chat/Chat.vue'),
					meta: {
						neededPermission: 'MANAGE_OVERLAYS',
						noPadding: true,
					},
				},
				{
					name: 'KappagenOverlay',
					path: '/dashboard/overlays/kappagen',
					component: () => import('../pages/overlays/kappagen/Kappagen.vue'),
					meta: {
						neededPermission: 'MANAGE_OVERLAYS',
						noPadding: true,
					},
				},
				{
					name: 'BrbOverlay',
					path: '/dashboard/overlays/brb',
					component: () => import('../pages/overlays/brb/Brb.vue'),
					meta: {
						neededPermission: 'MANAGE_OVERLAYS',
						noPadding: true,
					},
				},
				{
					name: 'DudesOverlay',
					path: '/dashboard/overlays/dudes',
					component: () => import('../pages/overlays/dudes/dudes-settings.vue'),
					meta: {
						neededPermission: 'MANAGE_OVERLAYS',
						fullScreen: true,
					},
				},
				{
					path: '/dashboard/events/chat-alerts',
					component: () => import('../pages/ChatAlerts.vue'),
					meta: { neededPermission: 'VIEW_EVENTS' },
				},
				{
					path: '/dashboard/events/custom',
					component: () => import('../pages/Events.vue'),
					meta: { neededPermission: 'VIEW_EVENTS' },
				},
				{
					path: '/dashboard/alerts',
					component: () => import('../pages/Alerts.vue'),
					meta: { neededPermission: 'VIEW_ALERTS' },
				},
				{
					path: '/dashboard/games',
					component: () => import('../pages/Games.vue'),
					meta: { neededPermission: 'VIEW_GAMES' },
				},
				{
					path: '/dashboard/files',
					component: () => import('../pages/Files.vue'),
				},
				{
					name: 'RegistryOverlayEdit',
					path: '/dashboard/registry/overlays/:id',
					component: () => import('../components/registry/overlays/edit.vue'),
					meta: { neededPermission: 'MANAGE_OVERLAYS' },
				},
				{
					name: 'Moderation',
					path: '/dashboard/moderation',
					component: () => import('../pages/Moderation.vue'),
					meta: { neededPermission: 'MANAGE_MODERATION' },
				},
				{
					name: 'Settings',
					path: '/dashboard/settings',
					component: () => import('../pages/UserSettings.vue'),
					meta: { noPadding: true },
				},
				{
					name: 'Giveaways',
					path: '/dashboard/giveaways',
					component: () => import('../pages/Giveaways.vue'),
					meta: { neededPermission: 'MANAGE_GIVEAWAYS' },
				},
				{
					name: 'Forbidden',
					path: '/dashboard/forbidden',
					component: () => import('../pages/NoAccess.vue'),
					meta: { fullScreen: true },
				},
				{
					path: '/:pathMatch(.*)*',
					name: 'NotFound',
					component: () => import('../pages/NotFound.vue'),
					meta: { fullScreen: true },
				},
			],
		},
	];

	const router = createRouter({
		history: createWebHistory(),
		// eslint-disable-next-line @typescript-eslint/ban-ts-comment
		// @ts-ignore
		routes,
	});

	router.beforeEach(async (to, _, next) => {
		try {
			const profile = await queryClient.ensureQueryData(profileQueryOptions);
			await queryClient.ensureQueryData(dashboardsQueryOptions);

			if (!profile) {
				return window.location.replace('/');
			}

			if (!to.meta.neededPermission) return next();

			const hasAccess = await userAccessFlagChecker(queryClient, to.meta.neededPermission);
			if (hasAccess) {
				return next();
			}

			return next({ name: 'Forbidden' });
		} catch (error) {
			console.log(error);
			window.location.replace('/');
		}
	});

	return router;
};
