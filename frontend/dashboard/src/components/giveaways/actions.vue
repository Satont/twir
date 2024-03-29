<script setup lang="ts">
import { IconPencil, IconTrash } from '@tabler/icons-vue';
import { Giveaway } from '@twir/api/messages/giveaways/giveaways';
import { useI18n } from 'vue-i18n';

import { useGiveawaysManager } from '@/api';
import { Button } from '@/components/ui/button';


const emits = defineEmits<{ edit: [] }>();
const props = defineProps<{ row: Giveaway }>();

const manager = useGiveawaysManager();
const deleter = manager.deleteOne;

const { t } = useI18n();

function edit() {
	emits('edit');
}

async function deleteClick() {
	await deleter.mutateAsync({ id: props.row.id });
}

</script>

<template>
	<div class="flex items-center gap-1">
		<Button variant="ghost" size="icon" @click="edit">
			<IconPencil class="h-4 w-4" />
		</Button>
		<Button variant="ghost" size="icon" @click="deleteClick">
			<IconTrash class="h-4 w-4" />
		</Button>
	</div>
</template>
