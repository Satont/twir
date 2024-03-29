<script setup lang="ts">
import { IconRotate, IconTrophyFilled } from '@tabler/icons-vue';
import { NCard, NButton } from 'naive-ui';
import { ref, onMounted, onUnmounted } from 'vue';
import { useI18n } from 'vue-i18n';

import { EditableGiveaway } from './types';

import { useClearGiveawayParticipants, useParticipants } from '@/api';


const props = defineProps<{
	giveaway: EditableGiveaway | null
}>();

const emit = defineEmits<{
	openWinners: []
}>();

const { t } = useI18n();

const searchValue = ref('');

const { data: participants, refetch: refreshParticipants } = useParticipants(props.giveaway?.id ?? '', searchValue.value);

let intervalId = 0;

function refreshUsers() {
	intervalId = setInterval(async () => {
		if (!props.giveaway?.isFinished) {
			await refreshParticipants();
		}
	}, 5000);
}

const resetParticipants = useClearGiveawayParticipants(props.giveaway?.id ?? '');

async function resetUsers() {
	await resetParticipants.mutateAsync({
		giveawayId: props.giveaway?.id ?? '',
	});
}

onMounted(refreshUsers);
onUnmounted(() => {
	clearInterval(intervalId);
});

async function onWinnersClick() {
	emit('openWinners');
}

</script>

<template>
	<n-card :title="t('giveaways.users.title')" content-style="padding: 0;" header-style="padding: 10px;" style="min-width: 300px; height: 100%" segmented>
		<div style="padding: 10px">
			<div style="display: flex; flex-direction: row; align-items: flex-start; justify-content: flex-start;">
				<!-- <n-input v-mode:value="searchValue" type="text" placeholder="eg. TwirApp" /> -->
				<n-button type="tertiary" style="display: flex; align-items: center; justify-content: center; height: 34px; margin: 0 5px;" @click="onWinnersClick">
					<IconTrophyFilled />
				</n-button>
				<n-button :disabled="giveaway?.isFinished" type="tertiary" style="display: flex; align-items: center; justify-content: center; height: 34px; margin: 0 5px;" @click="resetUsers">
					<IconRotate />
				</n-button>
			</div>
		</div>
		<div style="padding: 10px; overflow-y: scoll">
			<ul style="list-style-type: none">
				<!-- TODO: this is not winners but participants -->
				<li v-for="participant in participants?.winners" :key="participant.displayName">
					{{ participant.displayName }}
				</li>
			</ul>
		</div>
	</n-card>
</template>
