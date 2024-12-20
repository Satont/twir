import { createRouter, createWebHistory } from 'vue-router'

import { urqlClient } from './urql.js'

import type { RouteRecordRaw } from 'vue-router'

import { profileQuery, userAccessFlagChecker } from '@/api/auth.js'
import { ChannelRolePermissionEnum } from '@/gql/graphql.js'

export function newRouter() {
	const routes: ReadonlyArray<RouteRecordRaw> = [
		{
			path: '/dashboard/integrations/:integrationName',
			component: () => import('../pages/IntegrationsCallback.vue'),
		},
		{
			path: '/dashboard/popup',
			component: () => import('../popup-layout/popup-layout.vue'),
			children: [
				{
					path: '/dashboard/popup/widgets/eventslist',
					component: () => import('../components/dashboard/events.vue'),
					props: { popup: true },
				},
				{
					path: '/dashboard/popup/widgets/audit-log',
					component: () => import('../components/dashboard/audit-logs.vue'),
					props: { popup: true },
				},
				{
					path: '/dashboard/popup/widgets/stream',
					component: () => import('@/features/dashboard/widgets/stream.vue'),
					props: { popup: true },
				},
			],
		},
		{
			path: '/dashboard',
			component: () => import('../layout/layout.vue'),
			children: [
				{
					path: '/dashboard',
					component: () => import('../pages/Dashboard.vue'),
					meta: { noPadding: true },
				},
				{
					name: 'Integrations',
					path: '/dashboard/integrations',
					component: () => import('../pages/Integrations.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewIntegrations },
				},
				{
					path: '/dashboard/commands/:system',
					component: () => import('../features/commands/commands.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewCommands, noPadding: true },
				},
				{
					path: '/dashboard/commands/:system/:id',
					component: () => import('../features/commands/commands-edit.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ManageCommands, noPadding: true },
				},
				{
					path: '/dashboard/timers',
					component: () => import('../features/timers/timers.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewTimers, noPadding: true },
				},
				{
					path: '/dashboard/timers/:id',
					component: () => import('../features/timers/timers-edit.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ManageTimers, noPadding: true },
				},
				{
					path: '/dashboard/keywords',
					component: () => import('../pages/Keywords.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewKeywords, noPadding: true },
				},
				{
					path: '/dashboard/variables',
					component: () => import('../features/variables/variables.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewVariables, noPadding: true },
				},
				{
					path: '/dashboard/variables/:id',
					component: () => import('../features/variables/variables-edit.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ManageVariables, noPadding: true },
				},
				{
					path: '/dashboard/greetings',
					component: () => import('../pages/greetings.vue'),
					meta: {
						noPadding: true,
						neededPermission: ChannelRolePermissionEnum.ViewGreetings,
					},
				},
				{
					path: '/dashboard/community',
					component: () => import('../pages/community.vue'),
					meta: { noPadding: true },
				},
				{
					path: '/dashboard/community/roles',
					component: () => import('../features/community-roles/community-roles.vue'),
					meta: {
						neededPermission: ChannelRolePermissionEnum.ViewRoles,
						noPadding: true,
					},
				},
				{
					path: '/dashboard/song-requests',
					component: () => import('../pages/SongRequests.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewSongRequests },
				},
				{
					path: '/dashboard/overlays',
					component: () => import('../pages/Overlays.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewOverlays },
				},
				{
					name: 'ChatOverlay',
					path: '/dashboard/overlays/chat',
					component: () => import('../pages/overlays/chat/Chat.vue'),
					meta: {
						noPadding: true,
						neededPermission: ChannelRolePermissionEnum.ManageOverlays,
					},
				},
				{
					name: 'KappagenOverlay',
					path: '/dashboard/overlays/kappagen',
					component: () => import('../pages/overlays/kappagen/Kappagen.vue'),
					meta: {
						noPadding: true,
						neededPermission: ChannelRolePermissionEnum.ManageOverlays,
					},
				},
				{
					name: 'BrbOverlay',
					path: '/dashboard/overlays/brb',
					component: () => import('../pages/overlays/brb/Brb.vue'),
					meta: {
						noPadding: true,
						neededPermission: ChannelRolePermissionEnum.ManageOverlays,
					},
				},
				{
					name: 'DudesOverlay',
					path: '/dashboard/overlays/dudes',
					component: () => import('../pages/overlays/dudes/dudes-settings.vue'),
					meta: {
						fullScreen: true,
						neededPermission: ChannelRolePermissionEnum.ManageOverlays,
					},
				},
				{
					name: 'FaceitStatuOverlay',
					path: '/dashboard/overlays/faceit-stats',
					component: () => import('../features/overlays/faceit-stats/builder.vue'),
					meta: {
						noPadding: true,
					},
				},
				{
					path: '/dashboard/events/chat-alerts',
					component: () => import('../pages/chat-alerts.vue'),
					meta: {
						noPadding: true,
						neededPermission: ChannelRolePermissionEnum.ViewEvents,
					},
				},
				{
					path: '/dashboard/events/custom',
					component: () => import('../pages/Events.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewEvents },
				},
				{
					path: '/dashboard/alerts',
					component: () => import('../pages/alerts.vue'),
					meta: {
						noPadding: true,
						neededPermission: ChannelRolePermissionEnum.ViewAlerts,
					},
				},
				{
					path: '/dashboard/games',
					component: () => import('../pages/Games.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewGames },
				},
				{
					path: '/dashboard/files',
					component: () => import('../pages/Files.vue'),
				},
				{
					name: 'RegistryOverlayEdit',
					path: '/dashboard/registry/overlays/:id',
					component: () => import('../components/registry/overlays/edit.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewOverlays },
				},
				{
					name: 'Moderation',
					path: '/dashboard/moderation',
					component: () => import('../pages/Moderation.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ManageModeration },
				},
				{
					name: 'Settings',
					path: '/dashboard/settings',
					component: () => import('../pages/user-settings/user-settings.vue'),
					meta: { noPadding: true },
				},
				{
					name: 'AdminPanel',
					path: '/dashboard/admin',
					component: () => import('../pages/admin-panel.vue'),
					meta: { noPadding: true, adminOnly: true },
				},
				{
					name: 'Import',
					path: '/dashboard/import',
					component: () => import('../pages/Import.vue'),
				},
				{
					name: 'Notifications',
					path: '/dashboard/notifications',
					component: () => import('../pages/notifications.vue'),
					meta: { noPadding: true },
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
	]

	const router = createRouter({
		history: createWebHistory(),
		routes,
	})

	router.beforeEach(async (to, _, next) => {
		if (to.path.startsWith('/dashboard/popup')) return next()

		try {
			const profileRequest = await urqlClient.value.executeQuery(profileQuery)
			if (!profileRequest.data) {
				return window.location.replace('/')
			}

			if (to.meta.adminOnly && !profileRequest.data.authenticatedUser.isBotAdmin) {
				return next({ name: 'NotFound' })
			}

			if (!to.meta.neededPermission) return next()

			const hasAccess = await userAccessFlagChecker(to.meta.neededPermission)
			if (hasAccess) {
				return next()
			}

			return next({ name: 'Forbidden' })
		} catch (error) {
			console.log(error)
			window.location.replace('/')
		}
	})

	return router
}
