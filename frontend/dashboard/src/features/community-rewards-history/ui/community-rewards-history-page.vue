<script setup lang="ts">
import CommunityRewardsTable from './community-rewards-history-table.vue'
</script>

<template>
	<div class="flex flex-col w-full gap-4">
		<slot name="pagination" />
		<CommunityRewardsTable />
		<slot name="pagination" />
	</div>
</template>
