import { useQueryClient } from '@tanstack/vue-query'
import { createRequest, useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { computed } from 'vue'

import { useMutation } from '@/composables/use-mutation.js'
import { graphql } from '@/gql'
import { ChannelRolePermissionEnum } from '@/gql/graphql.js'
import { urqlClient, useUrqlClient } from '@/plugins/urql.js'

export const profileQuery = createRequest(graphql(`
	query AuthenticatedUser {
		authenticatedUser {
			id
			isBotAdmin
			isBanned
			isEnabled
			isBotModerator
			hideOnLandingPage
			botId
			apiKey
			twitchProfile {
				description
				displayName
				login
				profileImageUrl
			}
			selectedDashboardId
			selectedDashboardTwitchUser {
				login
				displayName
				profileImageUrl
				offlineImageUrl
			}
			availableDashboards {
				id
				flags
				twitchProfile {
					login
					displayName
					profileImageUrl
				}
			}
		}
	}
`), {})

export const userInvalidateQueryKey = 'UserInvalidateQueryKey'

export const useProfile = createGlobalState(() => {
	const { data: response, executeQuery, fetching } = useQuery({
		query: profileQuery.query,
		variables: {},
		context: {
			key: profileQuery.key,
			additionalTypenames: [userInvalidateQueryKey],
		},
	})

	const computedUser = computed(() => {
		const user = response.value?.authenticatedUser
		if (!user) return null

		return {
			id: user.id,
			avatar: user.twitchProfile.profileImageUrl,
			login: user.twitchProfile.login,
			displayName: user.twitchProfile.displayName,
			apiKey: user.apiKey,
			isBotAdmin: user.isBotAdmin,
			isEnabled: user.isEnabled,
			isBanned: user.isBanned,
			isBotModerator: user.isBotModerator,
			botId: user.botId,
			selectedDashboardId: user.selectedDashboardId,
			selectedDashboardTwitchUser: user.selectedDashboardTwitchUser,
			hideOnLandingPage: user.hideOnLandingPage,
			availableDashboards: user.availableDashboards,
		}
	})

	return { data: computedUser, executeQuery, isLoading: fetching }
})

export function useLogout() {
	const { executeMutation } = useMutation(graphql(`
		mutation userLogout {
			logout
		}
	`))

	async function execute() {
		const result = await executeMutation({})
		if (result.error) throw new Error(result.error.toString())
		window.location.replace('/')
	}

	return execute
}

export const useUserSettings = createGlobalState(() => {
	const userPublicSettingsInvalidateKey = 'UserPublicSettingsInvalidateKey'

	const usePublicQuery = () => useQuery({
		query: graphql(`
			query userPublicSettings {
				userPublicSettings {
					description
					socialLinks {
						href
						title
					}
				}
			}
		`),
		variables: {},
		context: {
			additionalTypenames: [userPublicSettingsInvalidateKey],
		},
	})

	const usePublicMutation = () => useMutation(graphql(`
		mutation userPublicSettingsUpdate($opts: UserUpdatePublicSettingsInput!) {
			authenticatedUserUpdatePublicPage(opts: $opts)
		}
	`), [userPublicSettingsInvalidateKey])

	const useApiKeyGenerateMutation = () => useMutation(graphql(`
		mutation userRegenerateApiKey {
			authenticatedUserRegenerateApiKey
		}
	`), [userInvalidateQueryKey])

	const useUserUpdateMutation = () => useMutation(graphql(`
		mutation userUpdateSettings($opts: UserUpdateSettingsInput!) {
			authenticatedUserUpdateSettings(opts: $opts)
		}
	`), [userInvalidateQueryKey, userPublicSettingsInvalidateKey])

	return {
		usePublicQuery,
		usePublicMutation,
		useApiKeyGenerateMutation,
		useUserUpdateMutation,
	}
})

export const useDashboard = createGlobalState(() => {
	const urqlClient = useUrqlClient()

	const mutationSetDashboard = useMutation(graphql(`
		mutation SetDashboard($dashboardId: String!) {
			authenticatedUserSelectDashboard(dashboardId: $dashboardId)
		}
	`))

	const queryClient = useQueryClient()

	async function setDashboard(dashboardId: string) {
		await mutationSetDashboard.executeMutation({ dashboardId })
		urqlClient.reInitClient()
		await queryClient.invalidateQueries()
		await queryClient.resetQueries()
	}

	return {
		setDashboard,
	}
})

type Flag = { perm: ChannelRolePermissionEnum, description: string } | 'delimiter'

export const PERMISSIONS_FLAGS: Flag[] = [
	{ perm: ChannelRolePermissionEnum.CanAccessDashboard, description: 'All permissions' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.UpdateChannelTitle, description: 'Can update channel title' },
	{ perm: ChannelRolePermissionEnum.UpdateChannelCategory, description: 'Can update channel category' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewCommands, description: 'Can view commands' },
	{ perm: ChannelRolePermissionEnum.ManageCommands, description: 'Can manage commands' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewKeywords, description: 'Can view keywords' },
	{ perm: ChannelRolePermissionEnum.ManageKeywords, description: 'Can manage keywords' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewTimers, description: 'Can view timers' },
	{ perm: ChannelRolePermissionEnum.ManageTimers, description: 'Can manage timers' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewIntegrations, description: 'Can view integrations' },
	{ perm: ChannelRolePermissionEnum.ManageIntegrations, description: 'Can manage integrations' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewSongRequests, description: 'Can view song requests' },
	{ perm: ChannelRolePermissionEnum.ManageSongRequests, description: 'Can manage song requests' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewModeration, description: 'Can view moderation settings' },
	{ perm: ChannelRolePermissionEnum.ManageModeration, description: 'Can manage moderation settings' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewVariables, description: 'Can view variables' },
	{ perm: ChannelRolePermissionEnum.ManageVariables, description: 'Can manage variables' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewGreetings, description: 'Can view greetings' },
	{ perm: ChannelRolePermissionEnum.ManageGreetings, description: 'Can manage greetings' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewOverlays, description: 'Can view overlays' },
	{ perm: ChannelRolePermissionEnum.ManageOverlays, description: 'Can manage overlays' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewRoles, description: 'Can view roles' },
	{ perm: ChannelRolePermissionEnum.ManageRoles, description: 'Can manage roles' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewEvents, description: 'Can view events' },
	{ perm: ChannelRolePermissionEnum.ManageEvents, description: 'Can manage events' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewAlerts, description: 'Can view alerts' },
	{ perm: ChannelRolePermissionEnum.ManageAlerts, description: 'Can manage alerts' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewGames, description: 'Can view games' },
	{ perm: ChannelRolePermissionEnum.ManageGames, description: 'Can manage games' },
]

export function useUserAccessFlagChecker(flag: ChannelRolePermissionEnum) {
	const { data: profile } = useProfile()

	return computed(() => {
		if (!profile.value?.availableDashboards || !profile.value?.selectedDashboardId) return false

		if (profile.value.id === profile.value.selectedDashboardId) {
			return true
		}

		if (profile.value.isBotAdmin) return true

		const dashboard = profile.value?.availableDashboards.find(dashboard => {
			return dashboard.id === profile.value?.selectedDashboardId
		})
		if (!dashboard) return false

		if (dashboard.flags.includes(ChannelRolePermissionEnum.CanAccessDashboard)) return true
		return dashboard.flags.includes(flag)
	})
}

export async function userAccessFlagChecker(flag: ChannelRolePermissionEnum) {
	const { data: profile } = await urqlClient.value.executeQuery(profileQuery)

	if (profile?.authenticatedUser.isBotAdmin) return true
	if (!profile || !profile?.authenticatedUser.selectedDashboardId) return false
	if (profile.authenticatedUser.selectedDashboardId === profile.authenticatedUser.id) return true

	const dashboard = profile.authenticatedUser.availableDashboards.find(dashboard => {
		return dashboard.id === profile.authenticatedUser.selectedDashboardId
	})
	if (!dashboard) return false

	if (dashboard.flags.includes(ChannelRolePermissionEnum.CanAccessDashboard)) return true
	return dashboard.flags.includes(flag)
}
