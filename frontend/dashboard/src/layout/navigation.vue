<script setup lang="ts">
import {
	IconActivity,
	IconBell,
	IconBox,
	IconCalendarEvent,
	IconClockHour7,
	IconCommand,
	IconDashboard,
	IconDeviceDesktop,
	IconHammer,
	IconDeviceGamepad2,
	IconHeadphones,
	IconKey,
	IconMessageCircle2,
	IconPencilPlus,
	IconShieldHalfFilled, IconSpeakerphone,
	IconSword,
	IconUsers,
	IconGift,
} from '@tabler/icons-vue';
import { MenuDividerOption, MenuOption, NBadge, NMenu, NDivider } from 'naive-ui';
import { computed, h, onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { RouterLink, useRouter } from 'vue-router';

import DashboardsMenu from './dashboardsMenu.vue';
import { renderIcon } from '../helpers/index.js';

import { useUserAccessFlagChecker } from '@/api';

const { t } = useI18n();

const activeKey = ref<string | null>('/');

const canViewIntegrations = useUserAccessFlagChecker('VIEW_INTEGRATIONS');
const canViewEvents = useUserAccessFlagChecker('VIEW_EVENTS');
const canViewOverlays = useUserAccessFlagChecker('VIEW_OVERLAYS');
const canViewSongRequests = useUserAccessFlagChecker('VIEW_SONG_REQUESTS');
const canViewCommands = useUserAccessFlagChecker('VIEW_COMMANDS');
const canViewTimers = useUserAccessFlagChecker('VIEW_TIMERS');
const canViewKeywords = useUserAccessFlagChecker('VIEW_KEYWORDS');
const canViewVariables = useUserAccessFlagChecker('VIEW_VARIABLES');
const canViewGreetings = useUserAccessFlagChecker('VIEW_GREETINGS');
const canViewRoles = useUserAccessFlagChecker('VIEW_ROLES');
const canViewAlerts = useUserAccessFlagChecker('VIEW_ALERTS');
const canViewGames = useUserAccessFlagChecker('VIEW_GAMES');
const canViewGiveaways = useUserAccessFlagChecker('VIEW_GIVEAWAYS');

const menuOptions = computed<(MenuOption | MenuDividerOption)[]>(() => {
	return [
		{
			label: t('sidebar.dashboard'),
			icon: renderIcon(IconDashboard),
			path: '/dashboard',
			isNew: false,
		},
		{
			label: t('sidebar.integrations'),
			icon: renderIcon(IconBox),
			path: '/dashboard/integrations',
			disabled: !canViewIntegrations.value,
		},
		{
			label: t('sidebar.alerts'),
			icon: renderIcon(IconBell),
			path: '/dashboard/alerts',
			disabled: !canViewAlerts.value,
		},
		{
			label: t('sidebar.chatAlerts'),
			icon: renderIcon(IconMessageCircle2),
			path: '/dashboard/events/chat-alerts',
			disabled: !canViewEvents.value,
		},
		{
			label: t('sidebar.events'),
			icon: renderIcon(IconCalendarEvent),
			disabled: !canViewEvents.value,
			path: '/dashboard/events/custom',
		},
		{
			label: t('sidebar.overlays'),
			icon: renderIcon(IconDeviceDesktop),
			path: '/dashboard/overlays',
			disabled: !canViewOverlays.value,
		},
		{
			label: t('sidebar.songRequests'),
			icon: renderIcon(IconHeadphones),
			path: '/dashboard/song-requests',
			disabled: !canViewSongRequests.value,
		},
		{
			label: t('sidebar.games'),
			icon: renderIcon(IconDeviceGamepad2),
			path: '/dashboard/games',
			disabled: !canViewGames.value,
		},
		{
			label: t('sidebar.commands.label'),
			icon: renderIcon(IconCommand),
			disabled: !canViewCommands.value,
			children: [
				{
					label: t('sidebar.commands.custom'),
					icon: renderIcon(IconPencilPlus),
					path: '/dashboard/commands/custom',
				},
				{
					label: t('sidebar.commands.builtin'),
					icon: renderIcon(IconHammer),
					path: '/dashboard/commands/builtin',
				},
			],
		},
		{
			label: t('sidebar.moderation'),
			icon: renderIcon(IconSword),
			path: '/dashboard/moderation',
		},
		{
			label: t('sidebar.users'),
			icon: renderIcon(IconUsers),
			path: '/dashboard/community/users',
		},
		{
			label: t('sidebar.roles'),
			icon: renderIcon(IconShieldHalfFilled),
			path: '/dashboard/community/roles',
			disabled: !canViewRoles.value,
		},
		{
			label: t('sidebar.timers'),
			icon: renderIcon(IconClockHour7),
			path: '/dashboard/timers',
			disabled: !canViewTimers.value,
		},
		{
			label: t('sidebar.keywords'),
			icon: renderIcon(IconKey),
			path: '/dashboard/keywords',
			disabled: !canViewKeywords.value,
		},
		{
			label: t('sidebar.variables'),
			icon: renderIcon(IconActivity),
			path: '/dashboard/variables',
			disabled: !canViewVariables.value,
		},
		{
			label: t('sidebar.greetings'),
			icon: renderIcon(IconSpeakerphone),
			path: '/dashboard/greetings',
			disabled: !canViewGreetings.value,
		},
		{
			label: t('sidebar.giveaways'),
			icon: renderIcon(IconGift),
			path: '/dashboard/giveaways',
			disabled: !canViewGiveaways.value,
			isNew: true,
		},
	].map((item) => ({
		...item,
		key: item.path ?? item.label,
		extra: item.disabled ? 'No perms' : undefined,
		label: !item.path || item.disabled ? item.label ?? undefined : () => h(
			RouterLink,
			{
				to: {
					path: item.path,
				},
			},
			{
				default: () => item.isNew
					? h(NBadge, {
						type: 'info',
						value: 'new',
						processing: true,
						offset: [17, 5],
					}, { default: () => item.label })
					: item.label,
			},
		),
		children: item.children?.map((child) => ({
			...child,
			key: child.path,
			label: item.disabled ? child.label : () => h(
				RouterLink,
				{
					to: {
						path: child.path,
					},
				},
				{ default: () => child.label },
			),
		})),
	}));
});

const router = useRouter();

onMounted(async () => {
	await router.isReady();
	activeKey.value = router.currentRoute.value.path;
});
</script>

<template>
	<div>
		<dashboards-menu />

		<n-divider style="margin-top: 0; margin-bottom: 5px;" />

		<n-menu
			v-model:value="activeKey"
			:collapsed-width="64"
			:collapsed-icon-size="22"
			:options="menuOptions"
		/>
	</div>
</template>

<style scoped>
:deep(.n-menu-item-content-header) {
	align-self: stretch;
	display: flex;
	align-items: center;
}
</style>
