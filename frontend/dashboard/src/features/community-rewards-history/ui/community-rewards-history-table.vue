<script setup lang="ts">
import { useCommunityRewardsTable } from '../composables/community-rewards-history-table'

import Table from '@/components/table.vue'

const rewardsTable = useCommunityRewardsTable()
</script>

<template>
	<Table
		:table="rewardsTable.table"
		:is-loading="rewardsTable.isLoading.value"
	/>
</template>
