<script setup lang="ts">
import { onMounted } from 'vue';

import { EditableGiveaway } from './types';

import { useGiveawaysWinners } from '@/api/giveaways';

const props = defineProps<{
	giveaway: EditableGiveaway | null
}>();
const winners = useGiveawaysWinners(props.giveaway?.id ?? '');

onMounted(() => {
	winners.refetch();
});

</script>

<template>
	<div
		:mask-closable="false" :segmented="true" preset="card" title="Winners" class="modal" :style="{
			width: '500px',
			height: '400px',
			overflowY: 'auto',
		}"
	>
		<ul style="list-style-type: none">
			<!-- TODO: this is not winners but participants -->
			<li v-for="winner in winners.data.value?.winners" :key="winner.displayName">
				{{ winner.displayName }}
			</li>
		</ul>
	</div>
</template>

