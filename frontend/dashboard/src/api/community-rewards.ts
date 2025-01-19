import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { unref } from 'vue'

import type { RewardsRedemptionsHistoryQuery, TwitchRedemptionsOpts } from '@/gql/graphql'
import type { MaybeRef } from 'vue'

import { graphql } from '@/gql/gql.js'

export type Redemption = RewardsRedemptionsHistoryQuery['rewardsRedemptionsHistory']['redemptions'][0]

export const useCommunityRewardsApi = createGlobalState(() => {
	const useHistory = (opts: MaybeRef<TwitchRedemptionsOpts>) => useQuery({
		query: graphql(`
			query RewardsRedemptionsHistory($opts: TwitchRedemptionsOpts!) {
				rewardsRedemptionsHistory(opts: $opts) {
					redemptions {
						id
						channelId
						user {
							id
							displayName
							login
							profileImageUrl
						}
						reward {
							id
							cost
							imageUrls
							title
							usedTimes
						}
						redeemedAt
						prompt
					}
					total
				}
			}
		`),
		get variables() {
			return {
				opts: unref(opts),
			}
		},
	})

	return {
		useHistory,
	}
})
