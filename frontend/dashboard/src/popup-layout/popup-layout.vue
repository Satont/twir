<script setup lang="ts">
import {
	NConfigProvider,
	NDialogProvider,
	NMessageProvider,
	NNotificationProvider,
	darkTheme,
	lightTheme,
} from 'naive-ui'
import { computed, onMounted } from 'vue'

import { Toaster as Sonner } from '@/components/ui/sonner'
import { Toaster } from '@/components/ui/toast'
import { useTheme } from '@/composables/use-theme'
import { useIsPopup } from '@/popup-layout/use-is-popup'

const { theme } = useTheme()
const themeStyles = computed(() => theme.value === 'dark' ? darkTheme : lightTheme)

const { setIsPopup } = useIsPopup()

onMounted(() => {
	setIsPopup(true)
})
</script>

<template>
	<NConfigProvider
		:theme="themeStyles"
		class="h-full"
		:breakpoints="{ 'xs': 0, 's': 640, 'm': 1024, 'l': 1280, 'xl': 1536, 'xxl': 1920, '2xl': 2560 }"
	>
		<NNotificationProvider>
			<NMessageProvider :duration="2500" :closable="true">
				<NDialogProvider>
					<RouterView />
					<Toaster />
					<Sonner />
				</NDialogProvider>
			</NMessageProvider>
		</NNotificationProvider>
	</NConfigProvider>
</template>
