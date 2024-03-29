<script setup lang="ts">
import { computed } from 'vue';

import { useGiveawaysManager } from '@/api/index.js';
import List from '@/components/giveaways/list.vue';

const giveawaysManager = useGiveawaysManager();
const { data: giveawaysResponse } = giveawaysManager.getAll({});
const giveaways = computed(() => {
	return giveawaysResponse.value?.giveaways ?? [];
});

</script>

<template>
	<list :giveaways="giveaways" :showHeader="true" :showCreateButton="true" />
</template>
