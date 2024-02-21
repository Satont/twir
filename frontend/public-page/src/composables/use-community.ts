import { useQuery } from '@tanstack/vue-query';
import {
	GetUsersRequest_Order,
	GetUsersRequest_SortBy,
} from '@twir/api/messages/community/community';
import { type ComputedRef, type Ref, unref } from 'vue';

import { unprotectedClient } from '@/api/twirp.js';
import { useStreamerProfile } from '@/composables/use-streamer-profile';

const sortBy = {
	'watched': GetUsersRequest_SortBy.Watched,
	'messages': GetUsersRequest_SortBy.Messages,
	'emotes': GetUsersRequest_SortBy.Emotes,
	'usedChannelPoints': GetUsersRequest_SortBy.UsedChannelPoints,
};

export type SortKey = keyof typeof sortBy

export type GetCommunityUsersOpts = {
	limit: number;
	page: number;
	desc: boolean;
	sortBy: SortKey;
}

export const useCommunityUsers = (options: Ref<GetCommunityUsersOpts> | ComputedRef<GetCommunityUsersOpts>) => {
	const { data: profile } = useStreamerProfile();

	return useQuery({
		queryKey: ['communityUsers', options],
		queryFn: async () => {
			const rawOpts = unref(options);

			const order = rawOpts.desc ? GetUsersRequest_Order.Desc : GetUsersRequest_Order.Asc;
			const call = await unprotectedClient.communityGetUsers({
				limit: rawOpts.limit,
				page: rawOpts.page,
				order,
				sortBy: sortBy[rawOpts.sortBy],
				channelId: profile.value!.id,
			}, { timeout: 5000 });
			return call.response;
		},
		enabled: () => !!profile.value,
	});
};
