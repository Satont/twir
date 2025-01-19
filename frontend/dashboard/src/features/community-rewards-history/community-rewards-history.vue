<script setup lang="ts">
import { useCommunityRewardsTable } from './composables/community-rewards-history-table'
import CommunityRewardsPage from './ui/community-rewards-history-page.vue'

import Pagination from '@/components/pagination.vue'

const rewardsTable = useCommunityRewardsTable()
</script>

<template>
	<CommunityRewardsPage>
		<template #pagination>
			<Pagination
				:total="rewardsTable.total.value"
				:table="rewardsTable.table"
				:pagination="rewardsTable.pagination.value"
				@update:page="(page) => rewardsTable.pagination.value.pageIndex = page"
				@update:page-size="(pageSize) => rewardsTable.pagination.value.pageSize = pageSize"
			/>
		</template>
	</CommunityRewardsPage>
</template>
