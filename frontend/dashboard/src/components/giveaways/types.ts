import { Giveaway } from '@twir/api/messages/giveaways/giveaways';

export type EditableGiveaway = Omit<
	Giveaway,
	'channelId' |
	'id' |
	'createdAt' |
	'finishedAt'
> & {
	id?: string
};
