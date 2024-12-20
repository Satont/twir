<script setup lang="ts">
import { useIntervalFn, useLocalStorage } from '@vueuse/core'
import { intervalToDuration } from 'date-fns'
import {
	HeartIcon,
	MessageCircleIcon,
	MusicIcon,
	SettingsIcon,
	SmileIcon,
	SquarePenIcon,
	StarIcon,
} from 'lucide-vue-next'
import { computed, onBeforeUnmount, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import Widget from './ui/widget.vue'

import type { FunctionalComponent } from 'vue'

import { useRealtimeDashboardStats } from '@/api'
import { useProfile } from '@/api/index.js'
import {
	DropdownMenu,
	DropdownMenuCheckboxItem,
	DropdownMenuContent,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { padTo2Digits } from '@/helpers/convertMillisToTime'
import { useIsPopup } from '@/popup-layout/use-is-popup'

const { t } = useI18n()
const { stats } = useRealtimeDashboardStats()
const { data: profile } = useProfile()
const { isPopup } = useIsPopup()

const currentTime = ref(new Date())
const { pause: pauseUptimeInterval } = useIntervalFn(() => {
	currentTime.value = new Date()
}, 1000)

const uptime = computed(() => {
	if (!stats.value?.startedAt) return '00:00:00'

	const duration = intervalToDuration({
		start: new Date(stats.value.startedAt),
		end: currentTime.value,
	})

	const mappedDuration = [duration.hours ?? 0, duration.minutes ?? 0, duration.seconds ?? 0]
	if (duration.days !== undefined && duration.days !== 0) mappedDuration.unshift(duration.days)

	return mappedDuration
		.map(v => padTo2Digits(v!))
		.filter(v => typeof v !== 'undefined')
		.join(':')
})

onBeforeUnmount(() => {
	pauseUptimeInterval()
})

const settings = useLocalStorage<{
	showPreview: boolean
	visibility: Record<string, boolean>
}>('twirWidgetsStream', {
	showPreview: true,
	visibility: {
		messages: true,
		followers: true,
		subs: true,
		usedEmotes: true,
		requestedSongs: true,
	},
})

const statsItems = computed<{
	key: string
	name: string
	count: number
	icon: FunctionalComponent
}[]>(() => [
	{
		key: 'messages',
		name: t(`dashboard.statsWidgets.messages`),
		count: stats.value.chatMessages ?? 0,
		icon: MessageCircleIcon,
	},
	{
		key: 'followers',
		name: t(`dashboard.statsWidgets.followers`),
		count: stats.value.followers ?? 0,
		icon: HeartIcon,
	},
	{
		key: 'subs',
		name: t(`dashboard.statsWidgets.subs`),
		count: stats.value.subs ?? 0,
		icon: StarIcon,
	},
	{
		key: 'usedEmotes',
		name: t(`dashboard.statsWidgets.usedEmotes`),
		count: stats.value.usedEmotes ?? 0,
		icon: SmileIcon,
	},
	{
		key: 'requestedSongs',
		name: t(`dashboard.statsWidgets.requestedSongs`),
		count: stats.value.requestedSongs ?? 0,
		icon: MusicIcon,
	},
])
const numberFormatter = new Intl.NumberFormat('fr-FR', {
	useGrouping: true,
})

const streamUrl = computed(() => {
	if (!profile.value) return

	const user = profile.value.selectedDashboardTwitchUser

	return `https://player.twitch.tv/?channel=${user.login}&parent=${window.location.host}&autoplay=false`
})

const popupHref = computed(() => {
	if (!profile.value) return

	return `${window.location.origin}/dashboard/popup/widgets/stream?apiKey=${profile.value.apiKey}`
})
</script>

<template>
	<Widget :popupHref>
		<template #title>
			{{ $attrs?.item?.i as string }}
		</template>

		<template #settings>
			<DropdownMenu>
				<DropdownMenuTrigger as-child>
					<button
						class="p-1 rounded hover:bg-white/10 outline-none focus-visible:ring-2 ring-white/15 active:bg-white/15"
					>
						<SettingsIcon
							class="size-4 text-white/60 stroke-[1.5]"
							absolute-stroke-width
						/>
					</button>
				</DropdownMenuTrigger>
				<DropdownMenuContent class="w-56">
					<DropdownMenuLabel>Appearance</DropdownMenuLabel>
					<DropdownMenuCheckboxItem
						v-model:checked="settings.showPreview"
					>
						Show stream preview
					</DropdownMenuCheckboxItem>
					<DropdownMenuSeparator />
					<DropdownMenuLabel>Showed items</DropdownMenuLabel>
					<DropdownMenuCheckboxItem
						v-for="(_, key) of settings.visibility"
						:key="key"
						v-model:checked="settings.visibility[key]"
					>
						{{ t(`dashboard.statsWidgets.${key}`) }}
					</DropdownMenuCheckboxItem>
				</DropdownMenuContent>
			</DropdownMenu>
		</template>

		<template #content>
			<div
				class="flex flex-col sm:flex-row divide-y bg-[#252525] sm:divide-x sm:divide-y-0 divide-white/10"
			>
				<div class="flex flex-1 min-w-0">
					<div class="flex flex-col gap-0.5 px-3 py-2 flex-1 min-w-0">
						<span class="text-base font-medium text-white truncate">
							{{ stats?.title ?? 'N/A' }}
						</span>
						<span class="text-sm text-white/60 truncate">
							{{ stats?.categoryName ?? 'N/A' }}
						</span>
					</div>
					<div class="p-2">
						<button
							class="p-1 rounded hover:bg-white/10 outline-none focus-visible:ring-2 ring-white/15 active:bg-white/15"
						>
							<SquarePenIcon class="size-4 text-white/60 stroke-[1.5]" absolute-stroke-width />
						</button>
					</div>
				</div>
				<div class="flex divide-x divide-white/10">
					<div class="pl-3 pr-4 py-2 flex flex-col gap-0.5 flex-1">
						<span class="text-sm text-white/60 truncate"> 	{{ t(`dashboard.statsWidgets.uptime`) }} </span>
						<span class="text-base font-medium text-white truncate">
							{{ uptime }}
						</span>
					</div>
					<div class="pl-3 pr-4 py-2 flex flex-col gap-0.5 flex-1">
						<span class="text-sm text-white/60 truncate"> {{ t(`dashboard.statsWidgets.viewers`) }} </span>
						<span class="text-base font-medium text-white truncate"> {{ stats.viewers ?? 0 }} </span>
					</div>
				</div>
			</div>

			<template v-if="settings.showPreview && !isPopup">
				<iframe
					v-if="stats.startedAt"
					:src="streamUrl"
					width="100%"
					height="100%"
					frameborder="0"
					scrolling="no"
					allowfullscreen="true"
					class="aspect-video w-full"
				>
				</iframe>
				<img
					v-else-if="profile?.selectedDashboardTwitchUser.offlineImageUrl"
					class="aspect-video w-full"
					:src="profile.selectedDashboardTwitchUser.offlineImageUrl"
				/>
				<div v-else class="aspect-video max-h-16 w-full bg-card text-4xl flex justify-center items-center">
					Offline
				</div>
			</template>

			<ul class="grid sm:grid-cols-2 lg:grid-cols-3 bg-[#343434] gap-y-px">
				<template
					v-for="stat of statsItems"
					:key="stat.name"
				>
					<li
						v-if="settings.visibility[stat.key]"
						class="p-2 flex items-center gap-1 bg-[#252525] sm:last:col-span-2"
					>
						<component :is="stat.icon" class="size-8 stroke-2 text-white/50 m-2" absolute-stroke-width />
						<div class="flex flex-col gap-0.5">
							<span class="text-sm text-white/60">{{ stat.name }}</span>
							<span class="text-white font-medium text-2xl tracking-wide">
								{{ numberFormatter.format(stat.count) }}
							</span>
						</div>
					</li>
				</template>
			</ul>
		</template>
	</Widget>
</template>
