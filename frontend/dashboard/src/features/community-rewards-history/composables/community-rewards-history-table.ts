import {
	type ColumnDef,
	getCoreRowModel,
	getFacetedRowModel,
	getFacetedUniqueValues,
	getFilteredRowModel,
	getPaginationRowModel,
	getSortedRowModel,
	useVueTable,
} from '@tanstack/vue-table'
import { createGlobalState } from '@vueuse/core'
import { computed, h } from 'vue'

import CommunityRewardsTableRewardCell from '../ui/cells/community-rewards-history-table-reward-cell.vue'

import type { Redemption } from '@/api/community-rewards'
import type { TwitchRedemptionsOpts } from '@/gql/graphql'

import { useProfile } from '@/api/auth.js'
import { useCommunityRewardsApi } from '@/api/community-rewards'
import { usePagination } from '@/composables/use-pagination'
import UsersTableCellUser from '@/features/admin-panel/manage-users/ui/users-table-cell-user.vue'
import { valueUpdater } from '@/helpers/value-updater'

export const useCommunityRewardsTable = createGlobalState(() => {
	const communityRewardsApi = useCommunityRewardsApi()
	const { data: profile } = useProfile()
	// пример откуда взять реварды. Но ты же наверняка будешь получение выносить в сам компонент(композабл?) селекта
	// const { data: rewards } = useTwitchRewardsNew()

	const { pagination, setPagination } = usePagination()
	const params = computed<TwitchRedemptionsOpts>(() => {
		return {
			byChannelId: profile.value?.selectedDashboardId,
			userSearch: undefined,
			page: pagination.value.pageIndex,
			perPage: pagination.value.pageSize,
			rewardsIds: [], // можно взять айдиники для селекта с rewards
		}
	})
	const historyResult = communityRewardsApi.useHistory(params)

	const history = computed<Redemption[]>(() => {
		return historyResult.data.value?.rewardsRedemptionsHistory.redemptions ?? []
	})
	const total = computed(() => {
		return historyResult.data.value?.rewardsRedemptionsHistory.total ?? 0
	})
	const pageCount = computed(() => {
		return Math.ceil(total.value / pagination.value.pageSize)
	})

	const tableColumns = computed<ColumnDef<Redemption>[]>(() => [
		{
			accessorKey: 'reward',
			size: 20,
			header: () => 'Reward',
			cell: ({ row }) => {
				return h(CommunityRewardsTableRewardCell, {
					name: row.original.reward.title,
					imageUrl: row.original.reward.imageUrls?.at(-1),
				})
			},
		},
		{
			accessorKey: 'user',
			size: 20,
			header: () => 'User',
			cell: ({ row }) => {
				return h('a', {
					class: 'flex flex-col',
					href: `https://twitch.tv/${row.original.user.login}`,
					target: '_blank',
				}, h(UsersTableCellUser, {
					avatar: row.original.user.profileImageUrl,
					name: row.original.user.login,
					displayName: row.original.user.displayName,
				}))
			},
		},
		{
			accessorKey: 'cost',
			size: 5,
			header: () => 'Cost',
			cell: ({ row }) => row.original.reward.cost,
		},
		{
			accessorKey: 'input',
			size: 65,
			header: () => 'User input',
			cell: ({ row }) => row.original.prompt,
		},
	])

	const table = useVueTable({
		get pageCount() {
			return pageCount.value
		},
		get data() {
			return history.value
		},
		get columns() {
			return tableColumns.value
		},
		manualPagination: true,
		enableRowSelection: true,
		onPaginationChange: (updaterOrValue) => valueUpdater(updaterOrValue, pagination),
		// onSortingChange: updaterOrValue => valueUpdater(updaterOrValue, sorting),
		// onColumnFiltersChange: updaterOrValue => valueUpdater(updaterOrValue, columnFilters),
		// onColumnVisibilityChange: updaterOrValue => valueUpdater(updaterOrValue, columnVisibility),
		// onRowSelectionChange: updaterOrValue => valueUpdater(updaterOrValue, rowSelection),
		getCoreRowModel: getCoreRowModel(),
		getFilteredRowModel: getFilteredRowModel(),
		getPaginationRowModel: getPaginationRowModel(),
		getSortedRowModel: getSortedRowModel(),
		getFacetedRowModel: getFacetedRowModel(),
		getFacetedUniqueValues: getFacetedUniqueValues(),
	})

	return {
		table,
		total,
		pageCount,
		pagination,
		setPagination,
		isLoading: historyResult.fetching,
	}
})
