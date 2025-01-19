<script setup lang='ts'>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import type { PageLayoutTab } from '@/layout/page-layout.vue'

import { useUserAccessFlagChecker } from '@/api'
import CommunityEmotesStatistic
	from '@/features/community-emotes-statistic/community-emotes-statistic.vue'
import CommunityRewardsHistory from '@/features/community-rewards-history/community-rewards-history.vue'
import CommunityRoles from '@/features/community-roles/community-roles.vue'
import CommunityUsers from '@/features/community-users/community-users.vue'
import { ChannelRolePermissionEnum } from '@/gql/graphql'
import PageLayout from '@/layout/page-layout.vue'

const { t } = useI18n()

const canViewRoles = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewRoles)

const tabs = computed<PageLayoutTab[]>(() => ([
	{
		title: t('community.users.title'),
		component: CommunityUsers,
		name: 'users',
	},
	{
		title: t('sidebar.roles'),
		component: CommunityRoles,
		name: 'permissions',
		disabled: !canViewRoles.value,
	},
	{
		title: t('community.emotesStatistic.title'),
		component: CommunityEmotesStatistic,
		name: 'emotes-stats',
	},
	{
		title: 'Rewards history',
		component: CommunityRewardsHistory,
		name: 'rewards-history',
	},
]))
</script>

<template>
	<PageLayout :tabs="tabs" active-tab="users">
		<template #title>
			{{ t('community.title') }}
		</template>
	</PageLayout>
</template>
