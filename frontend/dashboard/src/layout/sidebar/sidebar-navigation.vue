<script lang="ts" setup>
import { IconCalendarCog } from '@tabler/icons-vue'
import { useLocalStorage } from '@vueuse/core'
import {
	AudioLines,
	Bell,
	Blend,
	Box,
	ChevronRight,
	Dices,
	Import,
	LayoutDashboard,
	MessageCircleHeart,
	MessageCircleWarning,
	Package,
	PackageCheck,
	PackagePlus,
	Shield,
	Smile,
	SparklesIcon,
	Timer,
	UserCog,
	Users,
	Variable,
	WholeWord,
} from 'lucide-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

import { useUserAccessFlagChecker } from '@/api'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible'
import {
	SidebarGroup,
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
	SidebarMenuSub,
	SidebarMenuSubItem,
	useSidebar,
} from '@/components/ui/sidebar'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

const { t } = useI18n()
const currentRoute = useRoute()
const sidebar = useSidebar()

const canViewIntegrations = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewIntegrations)
const canViewEvents = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewEvents)
const canViewOverlays = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewOverlays)
const canViewSongRequests = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewSongRequests)
const canViewCommands = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewCommands)
const canViewTimers = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewTimers)
const canViewKeywords = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewKeywords)
const canViewVariables = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewVariables)
const canViewGreetings = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewGreetings)
// const canViewRoles = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewRoles)
const canViewAlerts = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewAlerts)
const canViewGames = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewGames)
const canViewModeration = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewModeration)

const twirSidebarOpenedStates = useLocalStorage<Record<string, boolean>>('twir-sidebar-opened-states', {
	commands: false,
	community: false,
})

const links = computed(() => {
	return [
		{
			name: t('sidebar.dashboard'),
			icon: LayoutDashboard,
			disabled: false,
			path: '/dashboard',
		},
		{
			name: t('sidebar.integrations'),
			icon: Box,
			disabled: !canViewIntegrations.value,
			path: '/dashboard/integrations',
		},
		{
			name: t('sidebar.alerts'),
			icon: Bell,
			disabled: !canViewAlerts.value,
			path: '/dashboard/alerts',
		},
		{
			name: t('sidebar.chatAlerts'),
			icon: MessageCircleWarning,
			disabled: !canViewEvents.value,
			path: '/dashboard/events/chat-alerts',
		},
		{
			name: t('sidebar.events'),
			icon: IconCalendarCog,
			disabled: !canViewEvents.value,
			path: '/dashboard/events/custom',
		},
		{
			name: t('sidebar.overlays'),
			icon: Blend,
			disabled: !canViewOverlays.value,
			path: '/dashboard/overlays',
		},
		{
			name: t('sidebar.songRequests'),
			icon: AudioLines,
			disabled: !canViewSongRequests.value,
			path: '/dashboard/song-requests',
		},
		{
			name: t('sidebar.games'),
			icon: Dices,
			disabled: !canViewGames.value,
			path: '/dashboard/games',
		},
		{
			name: t('sidebar.commands.label'),
			icon: Package,
			disabled: !canViewCommands.value,
			path: '/dashboard/commands',
			openStateKey: 'commands',
			child: [
				{
					name: t('sidebar.commands.custom'),
					icon: PackagePlus,
					path: '/dashboard/commands/custom',
				},
				{
					name: t('sidebar.commands.builtin'),
					icon: PackageCheck,
					path: '/dashboard/commands/builtin',
				},
			],
		},
		{
			name: t('sidebar.moderation'),
			icon: Shield,
			disabled: !canViewModeration.value,
			path: '/dashboard/moderation',
		},
		{
			name: t('sidebar.community'),
			icon: Users,
			path: '/dashboard/community',
			openStateKey: 'community',
			child: [
				{
					name: t('community.users.title'),
					icon: Users,
					path: '/dashboard/community?tab=users',
				},
				{
					name: t('sidebar.roles'),
					icon: UserCog,
					path: '/dashboard/community?tab=permissions',
				},
				{
					name: t('community.emotesStatistic.title'),
					icon: Smile,
					path: '/dashboard/community?tab=emotes-stats',
				},
				{
					name: 'Rewards history',
					icon: SparklesIcon,
					path: '/dashboard/community?tab=rewards-history',
				},
			],
		},
		{
			name: t('sidebar.timers'),
			icon: Timer,
			disabled: !canViewTimers.value,
			path: '/dashboard/timers',
		},
		{
			name: t('sidebar.keywords'),
			icon: WholeWord,
			disabled: !canViewKeywords.value,
			path: '/dashboard/keywords',
		},
		{
			name: t('sidebar.variables'),
			icon: Variable,
			disabled: !canViewVariables.value,
			path: '/dashboard/variables',
		},
		{
			name: t('sidebar.greetings'),
			icon: MessageCircleHeart,
			disabled: !canViewGreetings.value,
			path: '/dashboard/greetings',
		},
		{
			name: t('sidebar.import'),
			icon: Import,
			path: '/dashboard/import',
		},
	]
})

function goToRoute() {
	if (sidebar.isMobile.value) {
		sidebar.setOpenMobile(false)
	}
}
</script>

<template>
	<SidebarGroup>
		<SidebarMenu>
			<SidebarMenuItem
				v-for="item in links"
				:key="item.name"
			>
				<SidebarMenuButton
					v-if="!item.child"
					as-child
					:tooltip="item.name"
					:variant="currentRoute.path === item.path ? 'active' : 'default'"
					@click="goToRoute"
				>
					<RouterLink :to="item.path!">
						<component :is="item.icon" />
						<span>{{ item.name }}</span>
					</RouterLink>
				</SidebarMenuButton>
				<Collapsible
					v-else-if="item.openStateKey"
					v-model:open="twirSidebarOpenedStates[item.openStateKey]"
					as-child
					class="group/collapsible"
				>
					<SidebarMenuItem>
						<CollapsibleTrigger as-child>
							<SidebarMenuButton :tooltip="item.name" :variant="currentRoute.path.startsWith(item.path) ? 'active' : 'default'">
								<component :is="item.icon" />
								<span>{{ item.name }}</span>
								<ChevronRight class="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90" />
							</SidebarMenuButton>
						</CollapsibleTrigger>
						<CollapsibleContent>
							<SidebarMenuSub>
								<SidebarMenuSubItem
									v-for="child in item.child"
									:key="child.name"
								>
									<SidebarMenuButton as-child @click="goToRoute">
										<RouterLink :to="child.path!">
											<component :is="child.icon" />
											<span>{{ child.name }}</span>
										</RouterLink>
									</SidebarMenuButton>
								</SidebarMenuSubItem>
							</SidebarMenuSub>
						</CollapsibleContent>
					</SidebarMenuItem>
				</Collapsible>
			</SidebarMenuItem>
		</SidebarMenu>
	</SidebarGroup>
</template>
