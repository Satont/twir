import { useQuery } from '@tanstack/vue-query'
import { useSubscription } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { computed } from 'vue'

import { protectedApiClient } from './twirp'

import type { DashboardEventsSubscription, DashboardStats } from '@/gql/graphql'

import { graphql } from '@/gql'

export function useDashboardStats() {
	return useQuery({
		queryKey: ['dashboardStats'],
		queryFn: async () => {
			const call = await protectedApiClient.getDashboardStats({})

			return call.response
		},
	})
}

export type DashboardWidgetEvent = DashboardEventsSubscription['dashboardWidgetsEvents']['events'][number]

export const useDashboardEvents = createGlobalState(() => {
	const { data, isPaused, fetching } = useSubscription({
		query: graphql(`
			subscription DashboardEvents {
				dashboardWidgetsEvents {
					events {
						userId
						data {
							donationAmount
							donationCurrency
							donationMessage
							donationUserName
							raidedViewersCount
							raidedFromUserName
							raidedFromDisplayName
							followUserName
							followUserDisplayName
							redemptionTitle
							redemptionInput
							redemptionUserName
							redemptionUserDisplayName
							redemptionCost
							subLevel
							subUserName
							subUserDisplayName
							reSubLevel
							reSubUserName
							reSubUserDisplayName
							reSubMonths
							reSubStreak
							subGiftLevel
							subGiftUserName
							subGiftUserDisplayName
							subGiftTargetUserName
							subGiftTargetUserDisplayName
							firstUserMessageUserName
							firstUserMessageUserDisplayName
							firstUserMessageMessage
							banReason
							banEndsInMinutes
							bannedUserName
							bannedUserLogin
							moderatorName
							moderatorDisplayName
							message
							userLogin
							userName
						}
						createdAt
						type
					}
				}
			}
		`),
	})

	const events = computed(() => {
		return data.value?.dashboardWidgetsEvents.events
	})

	return {
		events,
		isPaused,
		fetching,
	}
})

export function useRealtimeDashboardStats() {
	const { data, isPaused, fetching } = useSubscription({
		query: graphql(`
			subscription dashboardStats {
				dashboardStats {
					categoryId
					categoryName
					viewers
					startedAt
					title
					chatMessages
					followers
					usedEmotes
					requestedSongs
					subs
				}
			}
		`),
	})

	const stats = computed<DashboardStats>(() => {
		return data.value?.dashboardStats ?? {}
	})

	return { stats, isPaused, fetching }
}
