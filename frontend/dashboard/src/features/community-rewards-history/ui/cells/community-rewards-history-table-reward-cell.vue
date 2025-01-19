<script setup lang="ts">
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'

defineProps<{
	imageUrl?: string
	name: string
}>()
</script>

<template>
	<div class="flex items-center gap-4 max-sm:justify-start">
		<Avatar class="size-9">
			<AvatarImage v-if="imageUrl" :src="imageUrl" :alt="name" loading="lazy" />
			<AvatarFallback>{{ name }}</AvatarFallback>
		</Avatar>
		{{ name }}
	</div>
</template>
