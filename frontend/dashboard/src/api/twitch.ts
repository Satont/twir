import { useMutation, useQuery } from '@tanstack/vue-query'
import { isRef } from 'vue'

import type { GetResponse as RewardsResponse } from '@twir/api/messages/rewards/rewards'
import type { TwitchGetUsersResponse, TwitchSearchChannelsRequest, TwitchSearchChannelsResponse } from '@twir/api/messages/twitch/twitch'
import type { ComputedRef, MaybeRef, Ref } from 'vue'

import { protectedApiClient, unprotectedApiClient } from '@/api/twirp.js'

type TwitchIn = Ref<string[]> | Ref<string> | ComputedRef<string> | ComputedRef<string[]> | string[]
export function useTwitchGetUsers(opts: {
	ids?: TwitchIn
	names?: TwitchIn
}) {
	return useQuery({
		queryKey: ['twitch', 'search', 'users', opts.ids, opts.names],
		queryFn: async (): Promise<TwitchGetUsersResponse> => {
			let ids = isRef(opts?.ids)
				? Array.isArray(opts.ids.value) ? opts.ids.value : [opts.ids.value]
				: opts?.ids ?? ['']
			let names = isRef(opts?.names)
				? Array.isArray(opts.names.value) ? opts.names.value : [opts.names.value]
				: opts?.names ?? ['']

			names = names.filter(n => n !== '')
			ids = ids.filter(n => n !== '')

			if (ids.length === 0 && names.length === 0) {
				return {
					users: [],
				}
			}

			const call = await unprotectedApiClient.twitchGetUsers({
				ids,
				names,
			})

			return call.response
		},
	})
}

export function useTwitchSearchChannels(params: Ref<TwitchSearchChannelsRequest>) {
	return useQuery({
		queryKey: ['twitch', 'search', 'channels', params],
		queryFn: async (): Promise<TwitchSearchChannelsResponse> => {
			const rawParams = isRef(params) ? params.value : params

			if (!rawParams.query) {
				return { channels: [] }
			}

			const call = await unprotectedApiClient.twitchSearchChannels(rawParams)
			return call.response
		},
	})
}

export function useTwitchRewards() {
	return useQuery({
		queryKey: ['twitchRewards'],
		queryFn: async (): Promise<RewardsResponse> => {
			const call = await protectedApiClient.rewardsGet({})
			return call.response
		},
	})
}

export function useTwitchSearchCategories(query: string | Ref<string>) {
	return useQuery({
		queryKey: ['twitchSearchCategories', query || ''],
		queryFn: async () => {
			const input = isRef(query) ? query.value : query
			if (!input) return { categories: [] }

			const call = await protectedApiClient.twitchSearchCategories({ query: input })
			return call.response
		},
	})
}

export function useTwitchGetCategories(ids: MaybeRef<string[]> | ComputedRef<string[]>) {
	return useQuery({
		queryKey: ['twitchGetCategories', ids || ''],
		queryFn: async () => {
			const input = isRef(ids) ? ids.value : ids
			if (!input) return { categories: [] }

			const call = await protectedApiClient.twitchGetCategories({ ids: input })
			return call.response
		},
	})
}

export function twitchSetChannelInformationMutation() {
	return useMutation({
		mutationKey: ['twitchSetChannelInformation'],
		mutationFn: async (req: { categoryId: string, title: string }) => {
			await protectedApiClient.twitchSetChannelInformation(req)
		},
	})
}
