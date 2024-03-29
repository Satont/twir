import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query';
import { ChooseWinnersRequest, ClearParticipantsRequest } from '@twir/api/messages/giveaways/giveaways';

import { protectedApiClient } from '@/api/twirp';

export const useParticipants = (id: string, query: string) => useQuery({
	queryKey: ['participants', id],
	queryFn: async () => {
		const req = await protectedApiClient.giveawaysGetParticipants({
			giveawayId: id,
			query: query,
		});
		return req.response;
	},
	enabled: !!id,
});

export const useClearGiveawayParticipants = (id: string) => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationKey: ['chooseWinners'],
		mutationFn: async (opts: ClearParticipantsRequest) => {
			await protectedApiClient.giveawaysClearParticipants({
				...opts,
			});
		},
		onSuccess: async () => {
			await queryClient.invalidateQueries(
				['participants', id],
			);
			await queryClient.invalidateQueries(
				winnersKey(id),
			);
		},
	});
};

const winnersKey = (id: string) => ['giveawayWinners', id];

export const useChooseGiveawayWinners = (id: string) => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationKey: winnersKey(id),
		mutationFn: async (opts: ChooseWinnersRequest) => {
			const req = await protectedApiClient.giveawaysChooseWinners({
				...opts,
			});

			return req.response;
		},
		onSuccess: async () => {
			await queryClient.invalidateQueries(
				winnersKey(id),
			);
		},
	});
};

export const useGiveawaysWinners = (giveawayId: string) => useQuery({
	queryKey: winnersKey(giveawayId),
	queryFn: async () => {
		const req = await protectedApiClient.giveawaysGetWinners({
			giveawayId: giveawayId,
		});
		return req.response;
	},
});
