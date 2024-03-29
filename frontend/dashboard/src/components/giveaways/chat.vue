<script lang="ts" setup>
import { IconSun, IconMoon } from '@tabler/icons-vue';
import { useLocalStorage } from '@vueuse/core';
import { NButton, NCard } from 'naive-ui';
import { computed } from 'vue';

import { useProfile, useTwitchGetUsers } from '@/api/index.js';
import { type Theme } from '@/hooks/index.js';

const { data: profile } = useProfile();

const chatTheme = useLocalStorage<Theme>('twirTwitchChatTheme', 'dark');
const toggleTheme = () => {
	chatTheme.value = chatTheme.value === 'light' ? 'dark' : 'light';
};

const selectedTwitchId = computed(() => profile.value?.selectedDashboardId ?? '');
const selectedDashboardTwitchUser = useTwitchGetUsers({ ids: selectedTwitchId });

const chatUrl = computed(() => {
	if (!selectedDashboardTwitchUser.data.value?.users.length) return;

	const user = selectedDashboardTwitchUser.data.value.users.at(0)!;

	let url = `https://www.twitch.tv/embed/${user.login}/chat?parent=${window.location.host}`;
	console.log(url);

	if (chatTheme.value === 'dark') {
		url += '&darkpopout';
	}

	return url;
});
</script>

<template>
	<n-card
		content-style="padding: 0;"
		header-style="padding: 10px;"
		style="min-width: 300px; height: 100%"
	>
		<template #header-extra>
			<n-button size="small" text @click="toggleTheme">
				<IconSun v-if="chatTheme === 'dark'" color="orange" />
				<IconMoon v-else />
			</n-button>
		</template>

		<iframe
			v-if="chatUrl"
			:src="chatUrl"
			width="100%"
			height="100%"
			frameborder="0"
			scrolling="no"
			allowfullscreen="true"
		>
		</iframe>
	</n-card>
</template>
