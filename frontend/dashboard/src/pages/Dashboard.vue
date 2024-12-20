<script setup lang="ts">
import { GridItem, GridLayout } from 'grid-layout-plus'
import { SquarePen } from 'lucide-vue-next'
import { NDropdown } from 'naive-ui'
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue'

import AuditLogs from '@/components/dashboard/audit-logs.vue'
import BotStatus from '@/components/dashboard/bot-status.vue'
import Chat from '@/components/dashboard/chat.vue'
import Events from '@/components/dashboard/events.vue'
import { useWidgets } from '@/components/dashboard/widgets.js'
import { Button } from '@/components/ui/button'
import { useIsMobile } from '@/composables/use-is-mobile'
import Stream from '@/features/dashboard/widgets/stream.vue'

const { isMobile } = useIsMobile()
const widgets = useWidgets()
const visibleWidgets = computed(() => widgets.value.filter((v) => v.visible))
const dropdownOptions = computed(() => {
	return widgets.value
		.filter((v) => !v.visible)
		.map((v) => ({ label: v.i, key: v.i }))
})

function addWidget(key: string) {
	const item = widgets.value.find(v => v.i === key)
	if (!item) return

	const widgetsLength = visibleWidgets.value.length

	item.visible = true
	item.x = (widgetsLength * 2) % 12
	item.y = widgetsLength + 12
}

const showEmptyItem = ref(false)

function onMouseUp() {
	showEmptyItem.value = false
}

onMounted(async () => {
	await nextTick()

	document.querySelectorAll('.vgl-item__resizer').forEach((el) => {
		el.addEventListener('mousedown', () => {
			showEmptyItem.value = true
		})

		window.addEventListener('mouseup', onMouseUp)
	})
})

onBeforeUnmount(() => {
	window.removeEventListener('mouseup', onMouseUp)
})
</script>

<template>
	<BotStatus />
	<div class="w-full h-full pl-1">
		<GridLayout
			v-model:layout="widgets"
			:row-height="30"
			:use-css-transforms="false"
		>
			<GridItem
				v-for="item in visibleWidgets"
				:key="item.i"
				:x="item.x"
				:y="item.y"
				:w="item.w"
				:h="item.h"
				:i="item.i"
				:min-w="item.minW"
				:min-h="item.minH"
				drag-allow-from=".widgets-draggable-handle"
			>
				<div v-if="showEmptyItem" class="w-full h-full absolute z-50"></div>
				<Chat v-if="item.i === 'chat'" :item="item" class="h-full" />
				<Stream v-if="item.i === 'stream'" :item="item" class="h-full" />
				<Events v-if="item.i === 'events'" :item="item" class="h-full" />
				<AuditLogs v-if="item.i === 'audit-logs'" :item="item" class="h-full" />
			</GridItem>
		</GridLayout>

		<div
			v-if="dropdownOptions.length"
			class="fixed right-[2rem] bottom-[2rem] z-50"
			:class="[{ '!right-[6rem]': isMobile }]"
		>
			<NDropdown
				size="huge"
				trigger="click"
				:options="dropdownOptions"
				@select="addWidget"
			>
				<Button variant="secondary" class="h-14 w-14" size="icon">
					<SquarePen class="size-8" />
				</Button>
			</NDropdown>
		</div>
	</div>
</template>

<style scoped>
.vgl-layout {
	@apply w-full;
}

:deep(.vgl-item__resizer) {
	z-index: 51;
}
</style>
