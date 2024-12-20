<script setup lang="ts">
import { ExternalLinkIcon, GripVerticalIcon, XIcon } from 'lucide-vue-next'
import { useAttrs } from 'vue'

import { useProfile } from '@/api'
import { type WidgetItem, useWidgets } from '@/components/dashboard/widgets'
import { useIsPopup } from '@/popup-layout/use-is-popup'

const props = defineProps<{
	popupHref?: string
}>()
const attrs = useAttrs() as { item: WidgetItem, [x: string]: unknown } | undefined

const widgets = useWidgets()
const { isPopup } = useIsPopup()

function hideWidget() {
	if (!attrs) return

	const item = widgets.value.find(item => item.i === attrs.item.i)
	if (!item) return
	item.visible = false
}

const { data: profile } = useProfile()

function openPopup() {
	if (!profile.value || !props.popupHref) return

	const height = 800
	const width = 500
	const top = Math.max(0, (screen.height - height) / 2)
	const left = Math.max(0, (screen.width - width) / 2)

	window.open(
		props.popupHref,
		'_blank',
		`height=${height},width=${width},top=${top},left=${left},status=0,location=0,menubar=0,toolbar=0`,
	)
}
</script>

<template>
	<div class="flex flex-col w-full overflow-hidden" :class="{ 'rounded-md': !isPopup }">
		<div
			v-if="!isPopup"
			class="flex justify-between items-center bg-[#343434] px-2 h-8 gap-1 w-full"
		>
			<div class="flex items-center gap-1 cursor-move widgets-draggable-handle">
				<GripVerticalIcon class="h-5 w-5 text-white/50 shrink-0 stroke-2" absolute-stroke-width />
				<span class="text-xs tracking-wide text-white/60 uppercase font-medium">
					<slot name="title" />
				</span>
			</div>
			<div class="flex">
				<slot name="settings" />
				<button
					v-if="popupHref"
					type="button"
					class="p-1 rounded hover:bg-white/10 outline-none focus-visible:ring-2 ring-white/15 active:bg-white/15"
					@click="openPopup"
				>
					<ExternalLinkIcon class="size-4 stroke-[1.5] text-white/60" absolute-stroke-width />
				</button>
				<button
					type="button"
					class="p-1 rounded hover:bg-white/10 outline-none focus-visible:ring-2 ring-white/15 active:bg-white/15"
					@click="hideWidget"
				>
					<XIcon class="size-4 stroke-[1.5] text-white/60" absolute-stroke-width />
				</button>
			</div>
		</div>

		<slot name="content" />
	</div>
</template>

<style scoped>

</style>
