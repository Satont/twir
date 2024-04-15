/* eslint-disable */
import type { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  Time: { input: any; output: any; }
  Upload: { input: File; output: File; }
};

export type AdminNotification = Notification & {
  __typename?: 'AdminNotification';
  createdAt: Scalars['Time']['output'];
  id: Scalars['ID']['output'];
  text: Scalars['String']['output'];
  twitchProfile?: Maybe<TwirUserTwitchInfo>;
  userId?: Maybe<Scalars['ID']['output']>;
};

export type AdminNotificationsParams = {
  page?: InputMaybe<Scalars['Int']['input']>;
  perPage?: InputMaybe<Scalars['Int']['input']>;
  search?: InputMaybe<Scalars['String']['input']>;
  type?: InputMaybe<NotificationType>;
};

export type AdminNotificationsResponse = {
  __typename?: 'AdminNotificationsResponse';
  notifications: Array<AdminNotification>;
  total: Scalars['Int']['output'];
};

export type AuthenticatedUser = TwirUser & {
  __typename?: 'AuthenticatedUser';
  apiKey: Scalars['String']['output'];
  botId?: Maybe<Scalars['ID']['output']>;
  hideOnLandingPage: Scalars['Boolean']['output'];
  id: Scalars['ID']['output'];
  isBanned: Scalars['Boolean']['output'];
  isBotAdmin: Scalars['Boolean']['output'];
  isBotModerator?: Maybe<Scalars['Boolean']['output']>;
  isEnabled?: Maybe<Scalars['Boolean']['output']>;
  twitchProfile: TwirUserTwitchInfo;
};

export type Badge = {
  __typename?: 'Badge';
  createdAt: Scalars['String']['output'];
  enabled: Scalars['Boolean']['output'];
  ffzSlot: Scalars['Int']['output'];
  fileUrl: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  name: Scalars['String']['output'];
  /** IDS of users which has this badge */
  users?: Maybe<Array<Scalars['String']['output']>>;
};

export enum ChannelRolePermissionEnum {
  CanAccessDashboard = 'CAN_ACCESS_DASHBOARD',
  ManageCommands = 'MANAGE_COMMANDS',
  ManageGreetings = 'MANAGE_GREETINGS',
  ManageIntegrations = 'MANAGE_INTEGRATIONS',
  ManageKeywords = 'MANAGE_KEYWORDS',
  ManageModeration = 'MANAGE_MODERATION',
  ManageOverlays = 'MANAGE_OVERLAYS',
  ManageSongRequests = 'MANAGE_SONG_REQUESTS',
  ManageTimers = 'MANAGE_TIMERS',
  ManageVariables = 'MANAGE_VARIABLES',
  UpdateChannelCategory = 'UPDATE_CHANNEL_CATEGORY',
  UpdateChannelTitle = 'UPDATE_CHANNEL_TITLE',
  ViewCommands = 'VIEW_COMMANDS',
  ViewGreetings = 'VIEW_GREETINGS',
  ViewIntegrations = 'VIEW_INTEGRATIONS',
  ViewKeywords = 'VIEW_KEYWORDS',
  ViewModeration = 'VIEW_MODERATION',
  ViewOverlays = 'VIEW_OVERLAYS',
  ViewSongRequests = 'VIEW_SONG_REQUESTS',
  ViewTimers = 'VIEW_TIMERS',
  ViewVariables = 'VIEW_VARIABLES'
}

export type Command = {
  __typename?: 'Command';
  aliases?: Maybe<Array<Scalars['String']['output']>>;
  allowedUsersIds?: Maybe<Array<Scalars['String']['output']>>;
  cooldown?: Maybe<Scalars['Int']['output']>;
  cooldownRolesIds?: Maybe<Array<Scalars['String']['output']>>;
  cooldownType: Scalars['String']['output'];
  default: Scalars['Boolean']['output'];
  defaultName?: Maybe<Scalars['String']['output']>;
  deniedUsersIds?: Maybe<Array<Scalars['String']['output']>>;
  description?: Maybe<Scalars['String']['output']>;
  enabled: Scalars['Boolean']['output'];
  enabledCategories?: Maybe<Array<Scalars['String']['output']>>;
  id: Scalars['ID']['output'];
  isReply: Scalars['Boolean']['output'];
  keepResponsesOrder: Scalars['Boolean']['output'];
  module: Scalars['String']['output'];
  name: Scalars['String']['output'];
  onlineOnly: Scalars['Boolean']['output'];
  requiredMessages: Scalars['Int']['output'];
  requiredUsedChannelPoints: Scalars['Int']['output'];
  requiredWatchTime: Scalars['Int']['output'];
  responses?: Maybe<Array<CommandResponse>>;
  rolesIds?: Maybe<Array<Scalars['String']['output']>>;
  visible: Scalars['Boolean']['output'];
};

export type CommandResponse = {
  __typename?: 'CommandResponse';
  commandId: Scalars['ID']['output'];
  id: Scalars['ID']['output'];
  order: Scalars['Int']['output'];
  text: Scalars['String']['output'];
};

export type CreateCommandInput = {
  aliases?: InputMaybe<Array<Scalars['String']['input']>>;
  description?: InputMaybe<Scalars['String']['input']>;
  name: Scalars['String']['input'];
  responses?: InputMaybe<Array<CreateCommandResponseInput>>;
};

export type CreateCommandResponseInput = {
  order: Scalars['Int']['input'];
  text: Scalars['String']['input'];
};

export type DashboardStats = {
  __typename?: 'DashboardStats';
  categoryId: Scalars['ID']['output'];
  categoryName: Scalars['String']['output'];
  chatMessages: Scalars['Int']['output'];
  followers: Scalars['Int']['output'];
  requestedSongs: Scalars['Int']['output'];
  startedAt?: Maybe<Scalars['Time']['output']>;
  subs: Scalars['Int']['output'];
  title: Scalars['String']['output'];
  usedEmotes: Scalars['Int']['output'];
  viewers?: Maybe<Scalars['Int']['output']>;
};

export type Greeting = {
  __typename?: 'Greeting';
  enabled: Scalars['Boolean']['output'];
  id: Scalars['ID']['output'];
  isReply: Scalars['Boolean']['output'];
  text: Scalars['String']['output'];
  twitchProfile: TwirUserTwitchInfo;
  userId: Scalars['String']['output'];
};

export type GreetingsCreateInput = {
  enabled: Scalars['Boolean']['input'];
  isReply: Scalars['Boolean']['input'];
  text: Scalars['String']['input'];
  userId: Scalars['String']['input'];
};

export type GreetingsUpdateInput = {
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  isReply?: InputMaybe<Scalars['Boolean']['input']>;
  text?: InputMaybe<Scalars['String']['input']>;
  userId?: InputMaybe<Scalars['String']['input']>;
};

export type IntegrationData = IntegrationDataDiscord | IntegrationDataDonationAlerts | IntegrationDataLastfm | IntegrationDataSevenTv | IntegrationDataSpotify | IntegrationDataStreamLabs | IntegrationDataValorant | IntegrationDataVk;

export type IntegrationDataDiscord = {
  __typename?: 'IntegrationDataDiscord';
  guilds: Array<IntegrationDataDiscordGuild>;
};

export type IntegrationDataDiscordGuild = {
  __typename?: 'IntegrationDataDiscordGuild';
  additionalUsersIdsForLiveCheck: Array<Scalars['String']['output']>;
  additionalUsersTwitchProfilesForLiveCheck: Array<TwirUserTwitchInfo>;
  channels: Array<IntegrationDataDiscordGuildChannel>;
  icon: Scalars['String']['output'];
  id: Scalars['String']['output'];
  liveNotificationAdditional_twitch_users_ids: Array<Scalars['String']['output']>;
  liveNotificationChannels_ids: Array<Scalars['String']['output']>;
  liveNotificationEnabled: Scalars['Boolean']['output'];
  liveNotificationMessage: Scalars['String']['output'];
  liveNotificationShowPreview: Scalars['Boolean']['output'];
  liveNotificationShowProfileImage: Scalars['Boolean']['output'];
  liveNotificationShowViewers: Scalars['Boolean']['output'];
  liveNotificationShow_category: Scalars['Boolean']['output'];
  liveNotificationShow_title: Scalars['Boolean']['output'];
  name: Scalars['String']['output'];
  offlineNotificationMessage: Scalars['String']['output'];
  roles: Array<IntegrationDataDiscordGuildRole>;
  shouldDeleteMessageOnOffline: Scalars['Boolean']['output'];
};

export type IntegrationDataDiscordGuildChannel = {
  __typename?: 'IntegrationDataDiscordGuildChannel';
  id: Scalars['String']['output'];
  name: Scalars['String']['output'];
};

export type IntegrationDataDiscordGuildRole = {
  __typename?: 'IntegrationDataDiscordGuildRole';
  id: Scalars['String']['output'];
  name: Scalars['String']['output'];
};

export type IntegrationDataDonationAlerts = {
  __typename?: 'IntegrationDataDonationAlerts';
  avatar: Scalars['String']['output'];
  username: Scalars['String']['output'];
};

export type IntegrationDataLastfm = {
  __typename?: 'IntegrationDataLastfm';
  avatar: Scalars['String']['output'];
  username: Scalars['String']['output'];
};

export type IntegrationDataSevenTv = {
  __typename?: 'IntegrationDataSevenTv';
  botSevenTvProfile: IntegrationDataSevenTvProfile;
  deleteEmotesOnlyAddedByApp: Scalars['Boolean']['output'];
  emoteSet_id: Scalars['String']['output'];
  isEditor: Scalars['Boolean']['output'];
  userSevenTvProfile: IntegrationDataSevenTvProfile;
};

export type IntegrationDataSevenTvProfile = {
  __typename?: 'IntegrationDataSevenTvProfile';
  avatar: Scalars['String']['output'];
  id: Scalars['String']['output'];
  username: Scalars['String']['output'];
};

export type IntegrationDataSpotify = {
  __typename?: 'IntegrationDataSpotify';
  avatar: Scalars['String']['output'];
  username: Scalars['String']['output'];
};

export type IntegrationDataStreamLabs = {
  __typename?: 'IntegrationDataStreamLabs';
  avatar: Scalars['String']['output'];
  username: Scalars['String']['output'];
};

export type IntegrationDataValorant = {
  __typename?: 'IntegrationDataValorant';
  avatar: Scalars['String']['output'];
  username: Scalars['String']['output'];
};

export type IntegrationDataVk = {
  __typename?: 'IntegrationDataVk';
  avatar: Scalars['String']['output'];
  username: Scalars['String']['output'];
};

export enum IntegrationService {
  Discord = 'DISCORD',
  Donatello = 'DONATELLO',
  Donatepay = 'DONATEPAY',
  DonateStream = 'DONATE_STREAM',
  Donationalerts = 'DONATIONALERTS',
  Faceit = 'FACEIT',
  Lastfm = 'LASTFM',
  Nightbot = 'NIGHTBOT',
  Spotify = 'SPOTIFY',
  Streamlabs = 'STREAMLABS',
  Valorant = 'VALORANT',
  Vk = 'VK'
}

export type Keyword = {
  __typename?: 'Keyword';
  cooldown: Scalars['Int']['output'];
  enabled: Scalars['Boolean']['output'];
  id: Scalars['ID']['output'];
  isRegularExpression: Scalars['Boolean']['output'];
  isReply: Scalars['Boolean']['output'];
  response?: Maybe<Scalars['String']['output']>;
  text: Scalars['String']['output'];
  usages: Scalars['Int']['output'];
};

export type KeywordCreateInput = {
  cooldown?: InputMaybe<Scalars['Int']['input']>;
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  isRegularExpression?: InputMaybe<Scalars['Boolean']['input']>;
  isReply?: InputMaybe<Scalars['Boolean']['input']>;
  response?: InputMaybe<Scalars['String']['input']>;
  text: Scalars['String']['input'];
  usageCount?: InputMaybe<Scalars['Int']['input']>;
};

export type KeywordUpdateInput = {
  cooldown?: InputMaybe<Scalars['Int']['input']>;
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  isRegularExpression?: InputMaybe<Scalars['Boolean']['input']>;
  isReply?: InputMaybe<Scalars['Boolean']['input']>;
  response?: InputMaybe<Scalars['String']['input']>;
  text?: InputMaybe<Scalars['String']['input']>;
  usageCount?: InputMaybe<Scalars['Int']['input']>;
};

export type Mutation = {
  __typename?: 'Mutation';
  badgesAddUser: Scalars['Boolean']['output'];
  badgesCreate: Badge;
  badgesDelete: Scalars['Boolean']['output'];
  badgesRemoveUser: Scalars['Boolean']['output'];
  badgesUpdate: Badge;
  createCommand: Command;
  greetingsCreate: Greeting;
  greetingsRemove: Scalars['Boolean']['output'];
  greetingsUpdate: Greeting;
  keywordCreate: Keyword;
  keywordRemove: Scalars['Boolean']['output'];
  keywordUpdate: Keyword;
  notificationsCreate: AdminNotification;
  notificationsDelete: Scalars['Boolean']['output'];
  notificationsUpdate: AdminNotification;
  removeCommand: Scalars['Boolean']['output'];
  switchUserAdmin: Scalars['Boolean']['output'];
  switchUserBan: Scalars['Boolean']['output'];
  timersCreate: Timer;
  timersRemove: Scalars['Boolean']['output'];
  timersUpdate: Timer;
  updateCommand: Command;
};


export type MutationBadgesAddUserArgs = {
  id: Scalars['ID']['input'];
  userId: Scalars['String']['input'];
};


export type MutationBadgesCreateArgs = {
  opts: TwirBadgeCreateOpts;
};


export type MutationBadgesDeleteArgs = {
  id: Scalars['ID']['input'];
};


export type MutationBadgesRemoveUserArgs = {
  id: Scalars['ID']['input'];
  userId: Scalars['String']['input'];
};


export type MutationBadgesUpdateArgs = {
  id: Scalars['ID']['input'];
  opts: TwirBadgeUpdateOpts;
};


export type MutationCreateCommandArgs = {
  opts: CreateCommandInput;
};


export type MutationGreetingsCreateArgs = {
  opts: GreetingsCreateInput;
};


export type MutationGreetingsRemoveArgs = {
  id: Scalars['String']['input'];
};


export type MutationGreetingsUpdateArgs = {
  id: Scalars['String']['input'];
  opts: GreetingsUpdateInput;
};


export type MutationKeywordCreateArgs = {
  opts: KeywordCreateInput;
};


export type MutationKeywordRemoveArgs = {
  id: Scalars['String']['input'];
};


export type MutationKeywordUpdateArgs = {
  id: Scalars['String']['input'];
  opts: KeywordUpdateInput;
};


export type MutationNotificationsCreateArgs = {
  text: Scalars['String']['input'];
  userId?: InputMaybe<Scalars['String']['input']>;
};


export type MutationNotificationsDeleteArgs = {
  id: Scalars['ID']['input'];
};


export type MutationNotificationsUpdateArgs = {
  id: Scalars['ID']['input'];
  opts: NotificationUpdateOpts;
};


export type MutationRemoveCommandArgs = {
  id: Scalars['String']['input'];
};


export type MutationSwitchUserAdminArgs = {
  userId: Scalars['ID']['input'];
};


export type MutationSwitchUserBanArgs = {
  userId: Scalars['ID']['input'];
};


export type MutationTimersCreateArgs = {
  opts: TimerCreateInput;
};


export type MutationTimersRemoveArgs = {
  id: Scalars['String']['input'];
};


export type MutationTimersUpdateArgs = {
  id: Scalars['String']['input'];
  opts: TimerUpdateInput;
};


export type MutationUpdateCommandArgs = {
  id: Scalars['String']['input'];
  opts: UpdateCommandOpts;
};

export type Notification = {
  createdAt: Scalars['Time']['output'];
  id: Scalars['ID']['output'];
  text: Scalars['String']['output'];
  userId?: Maybe<Scalars['ID']['output']>;
};

export enum NotificationType {
  Global = 'GLOBAL',
  User = 'USER'
}

export type NotificationUpdateOpts = {
  text?: InputMaybe<Scalars['String']['input']>;
};

export type Query = {
  __typename?: 'Query';
  authenticatedUser: AuthenticatedUser;
  commands: Array<Command>;
  greetings: Array<Greeting>;
  integrationsGetData?: Maybe<IntegrationData>;
  integrationsGetServiceAuthLink?: Maybe<Scalars['String']['output']>;
  keywords: Array<Keyword>;
  notificationsByAdmin: AdminNotificationsResponse;
  notificationsByUser: Array<UserNotification>;
  timers: Array<Timer>;
  /** Twir badges */
  twirBadges: Array<Badge>;
  /** finding users on twitch with filter does they exists in database */
  twirUsers: TwirUsersResponse;
};


export type QueryIntegrationsGetDataArgs = {
  service: IntegrationService;
};


export type QueryIntegrationsGetServiceAuthLinkArgs = {
  service: IntegrationService;
};


export type QueryNotificationsByAdminArgs = {
  opts: AdminNotificationsParams;
};


export type QueryTwirUsersArgs = {
  opts: TwirUsersSearchParams;
};

export type Subscription = {
  __typename?: 'Subscription';
  dashboardStats: DashboardStats;
  /** `newNotification` will return a stream of `Notification` objects. */
  newNotification: UserNotification;
};

export type Timer = {
  __typename?: 'Timer';
  enabled: Scalars['Boolean']['output'];
  id: Scalars['ID']['output'];
  messageInterval: Scalars['Int']['output'];
  name: Scalars['String']['output'];
  responses: Array<TimerResponse>;
  timeInterval: Scalars['Int']['output'];
};

export type TimerCreateInput = {
  enabled: Scalars['Boolean']['input'];
  messageInterval: Scalars['Int']['input'];
  name: Scalars['String']['input'];
  responses: Array<TimerResponseCreateInput>;
  timeInterval: Scalars['Int']['input'];
};

export type TimerResponse = {
  __typename?: 'TimerResponse';
  id: Scalars['ID']['output'];
  isAnnounce: Scalars['Boolean']['output'];
  text: Scalars['String']['output'];
};

export type TimerResponseCreateInput = {
  isAnnounce: Scalars['Boolean']['input'];
  text: Scalars['String']['input'];
};

export type TimerResponseUpdateInput = {
  isAnnounce: Scalars['Boolean']['input'];
  text: Scalars['String']['input'];
};

export type TimerUpdateInput = {
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  messageInterval?: InputMaybe<Scalars['Int']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  responses?: InputMaybe<Array<TimerResponseUpdateInput>>;
  timeInterval?: InputMaybe<Scalars['Int']['input']>;
};

export type TwirAdminUser = TwirUser & {
  __typename?: 'TwirAdminUser';
  apiKey: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  isBanned: Scalars['Boolean']['output'];
  isBotAdmin: Scalars['Boolean']['output'];
  isBotEnabled: Scalars['Boolean']['output'];
  isBotModerator: Scalars['Boolean']['output'];
  twitchProfile: TwirUserTwitchInfo;
};

export type TwirBadgeCreateOpts = {
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  ffzSlot: Scalars['Int']['input'];
  file: Scalars['Upload']['input'];
  name: Scalars['String']['input'];
};

export type TwirBadgeUpdateOpts = {
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  ffzSlot?: InputMaybe<Scalars['Int']['input']>;
  file?: InputMaybe<Scalars['Upload']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
};

export type TwirUser = {
  id: Scalars['ID']['output'];
  twitchProfile: TwirUserTwitchInfo;
};

export type TwirUserTwitchInfo = {
  __typename?: 'TwirUserTwitchInfo';
  description: Scalars['String']['output'];
  displayName: Scalars['String']['output'];
  id: Scalars['String']['output'];
  login: Scalars['String']['output'];
  profileImageUrl: Scalars['String']['output'];
};

export type TwirUsersResponse = {
  __typename?: 'TwirUsersResponse';
  total: Scalars['Int']['output'];
  users: Array<TwirAdminUser>;
};

export type TwirUsersSearchParams = {
  badges?: InputMaybe<Array<Scalars['String']['input']>>;
  isBanned?: InputMaybe<Scalars['Boolean']['input']>;
  isBotAdmin?: InputMaybe<Scalars['Boolean']['input']>;
  isBotEnabled?: InputMaybe<Scalars['Boolean']['input']>;
  page?: InputMaybe<Scalars['Int']['input']>;
  perPage?: InputMaybe<Scalars['Int']['input']>;
  search?: InputMaybe<Scalars['String']['input']>;
};

export type UpdateCommandOpts = {
  aliases?: InputMaybe<Array<Scalars['String']['input']>>;
  allowedUsersIds?: InputMaybe<Array<Scalars['String']['input']>>;
  cooldown?: InputMaybe<Scalars['Int']['input']>;
  cooldownRolesIds?: InputMaybe<Array<Scalars['String']['input']>>;
  cooldownType?: InputMaybe<Scalars['String']['input']>;
  deniedUsersIds?: InputMaybe<Array<Scalars['String']['input']>>;
  description?: InputMaybe<Scalars['String']['input']>;
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  enabledCategories?: InputMaybe<Array<Scalars['String']['input']>>;
  isReply?: InputMaybe<Scalars['Boolean']['input']>;
  keepResponsesOrder?: InputMaybe<Scalars['Boolean']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  onlineOnly?: InputMaybe<Scalars['Boolean']['input']>;
  requiredMessages?: InputMaybe<Scalars['Int']['input']>;
  requiredUsedChannelPoints?: InputMaybe<Scalars['Int']['input']>;
  requiredWatchTime?: InputMaybe<Scalars['Int']['input']>;
  responses?: InputMaybe<Array<CreateCommandResponseInput>>;
  rolesIds?: InputMaybe<Array<Scalars['String']['input']>>;
  visible?: InputMaybe<Scalars['Boolean']['input']>;
};

export type UserNotification = Notification & {
  __typename?: 'UserNotification';
  createdAt: Scalars['Time']['output'];
  id: Scalars['ID']['output'];
  text: Scalars['String']['output'];
  userId?: Maybe<Scalars['ID']['output']>;
};

export type CreateBadgeMutationVariables = Exact<{
  opts: TwirBadgeCreateOpts;
}>;


export type CreateBadgeMutation = { __typename?: 'Mutation', badgesCreate: { __typename?: 'Badge', id: string } };

export type DeleteBadgeMutationVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type DeleteBadgeMutation = { __typename?: 'Mutation', badgesDelete: boolean };

export type UpdateBadgeMutationVariables = Exact<{
  id: Scalars['ID']['input'];
  opts: TwirBadgeUpdateOpts;
}>;


export type UpdateBadgeMutation = { __typename?: 'Mutation', badgesUpdate: { __typename?: 'Badge', id: string } };

export type AddUserBadgeMutationVariables = Exact<{
  id: Scalars['ID']['input'];
  userId: Scalars['String']['input'];
}>;


export type AddUserBadgeMutation = { __typename?: 'Mutation', badgesAddUser: boolean };

export type RemoveUserBadgeMutationVariables = Exact<{
  id: Scalars['ID']['input'];
  userId: Scalars['String']['input'];
}>;


export type RemoveUserBadgeMutation = { __typename?: 'Mutation', badgesRemoveUser: boolean };

export type BadgesGetAllQueryVariables = Exact<{ [key: string]: never; }>;


export type BadgesGetAllQuery = { __typename?: 'Query', twirBadges: Array<{ __typename?: 'Badge', id: string, name: string, createdAt: string, fileUrl: string, enabled: boolean, ffzSlot: number, users?: Array<string> | null }> };

export type GetAllNotificationsQueryVariables = Exact<{ [key: string]: never; }>;


export type GetAllNotificationsQuery = { __typename?: 'Query', notificationsByUser: Array<{ __typename?: 'UserNotification', id: string, text: string, createdAt: any }> };

export type NotificationsSubscriptionSubscriptionVariables = Exact<{ [key: string]: never; }>;


export type NotificationsSubscriptionSubscription = { __typename?: 'Subscription', newNotification: { __typename?: 'UserNotification', id: string, text: string, createdAt: any } };

export type NotificationsByAdminQueryVariables = Exact<{
  opts: AdminNotificationsParams;
}>;


export type NotificationsByAdminQuery = { __typename?: 'Query', notificationsByAdmin: { __typename?: 'AdminNotificationsResponse', total: number, notifications: Array<{ __typename?: 'AdminNotification', id: string, text: string, userId?: string | null, createdAt: any, twitchProfile?: { __typename?: 'TwirUserTwitchInfo', displayName: string, profileImageUrl: string } | null }> } };

export type CreateNotificationMutationVariables = Exact<{
  text: Scalars['String']['input'];
  userId?: InputMaybe<Scalars['String']['input']>;
}>;


export type CreateNotificationMutation = { __typename?: 'Mutation', notificationsCreate: { __typename?: 'AdminNotification', id: string } };

export type DeleteNotificationMutationVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type DeleteNotificationMutation = { __typename?: 'Mutation', notificationsDelete: boolean };

export type UpdateNotificationsMutationVariables = Exact<{
  id: Scalars['ID']['input'];
  opts: NotificationUpdateOpts;
}>;


export type UpdateNotificationsMutation = { __typename?: 'Mutation', notificationsUpdate: { __typename?: 'AdminNotification', id: string } };

export type UsersGetAllQueryVariables = Exact<{
  opts: TwirUsersSearchParams;
}>;


export type UsersGetAllQuery = { __typename?: 'Query', twirUsers: { __typename?: 'TwirUsersResponse', total: number, users: Array<{ __typename?: 'TwirAdminUser', id: string, isBanned: boolean, isBotAdmin: boolean, isBotEnabled: boolean, twitchProfile: { __typename?: 'TwirUserTwitchInfo', login: string, displayName: string, profileImageUrl: string } }> } };

export type UserSwitchBanMutationVariables = Exact<{
  userId: Scalars['ID']['input'];
}>;


export type UserSwitchBanMutation = { __typename?: 'Mutation', switchUserBan: boolean };

export type UserSwitchAdminMutationVariables = Exact<{
  userId: Scalars['ID']['input'];
}>;


export type UserSwitchAdminMutation = { __typename?: 'Mutation', switchUserAdmin: boolean };

export type DashboardStatsSubscriptionVariables = Exact<{ [key: string]: never; }>;


export type DashboardStatsSubscription = { __typename?: 'Subscription', dashboardStats: { __typename?: 'DashboardStats', categoryId: string, categoryName: string, viewers?: number | null, startedAt?: any | null, title: string, chatMessages: number, followers: number, usedEmotes: number, requestedSongs: number, subs: number } };

export type GetAllTimersQueryVariables = Exact<{ [key: string]: never; }>;


export type GetAllTimersQuery = { __typename?: 'Query', timers: Array<{ __typename?: 'Timer', id: string, name: string, enabled: boolean, timeInterval: number, messageInterval: number, responses: Array<{ __typename?: 'TimerResponse', text: string, isAnnounce: boolean }> }> };

export type CreateTimerMutationVariables = Exact<{
  opts: TimerCreateInput;
}>;


export type CreateTimerMutation = { __typename?: 'Mutation', timersCreate: { __typename?: 'Timer', id: string } };

export type UpdateTimerMutationVariables = Exact<{
  id: Scalars['String']['input'];
  opts: TimerUpdateInput;
}>;


export type UpdateTimerMutation = { __typename?: 'Mutation', timersUpdate: { __typename?: 'Timer', id: string } };

export type RemoveTimerMutationVariables = Exact<{
  id: Scalars['String']['input'];
}>;


export type RemoveTimerMutation = { __typename?: 'Mutation', timersRemove: boolean };


export const CreateBadgeDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"CreateBadge"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"opts"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"TwirBadgeCreateOpts"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"badgesCreate"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"opts"},"value":{"kind":"Variable","name":{"kind":"Name","value":"opts"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}}]}}]} as unknown as DocumentNode<CreateBadgeMutation, CreateBadgeMutationVariables>;
export const DeleteBadgeDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"DeleteBadge"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"badgesDelete"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}]}]}}]} as unknown as DocumentNode<DeleteBadgeMutation, DeleteBadgeMutationVariables>;
export const UpdateBadgeDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UpdateBadge"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"opts"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"TwirBadgeUpdateOpts"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"badgesUpdate"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}},{"kind":"Argument","name":{"kind":"Name","value":"opts"},"value":{"kind":"Variable","name":{"kind":"Name","value":"opts"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}}]}}]} as unknown as DocumentNode<UpdateBadgeMutation, UpdateBadgeMutationVariables>;
export const AddUserBadgeDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"AddUserBadge"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"badgesAddUser"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}},{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}}]}]}}]} as unknown as DocumentNode<AddUserBadgeMutation, AddUserBadgeMutationVariables>;
export const RemoveUserBadgeDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"RemoveUserBadge"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"badgesRemoveUser"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}},{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}}]}]}}]} as unknown as DocumentNode<RemoveUserBadgeMutation, RemoveUserBadgeMutationVariables>;
export const BadgesGetAllDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"BadgesGetAll"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"twirBadges"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"fileUrl"}},{"kind":"Field","name":{"kind":"Name","value":"enabled"}},{"kind":"Field","name":{"kind":"Name","value":"ffzSlot"}},{"kind":"Field","name":{"kind":"Name","value":"users"}}]}}]}}]} as unknown as DocumentNode<BadgesGetAllQuery, BadgesGetAllQueryVariables>;
export const GetAllNotificationsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetAllNotifications"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"notificationsByUser"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"text"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]}}]} as unknown as DocumentNode<GetAllNotificationsQuery, GetAllNotificationsQueryVariables>;
export const NotificationsSubscriptionDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"subscription","name":{"kind":"Name","value":"NotificationsSubscription"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"newNotification"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"text"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]}}]} as unknown as DocumentNode<NotificationsSubscriptionSubscription, NotificationsSubscriptionSubscriptionVariables>;
export const NotificationsByAdminDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"NotificationsByAdmin"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"opts"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"AdminNotificationsParams"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"notificationsByAdmin"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"opts"},"value":{"kind":"Variable","name":{"kind":"Name","value":"opts"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"notifications"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"text"}},{"kind":"Field","name":{"kind":"Name","value":"userId"}},{"kind":"Field","name":{"kind":"Name","value":"twitchProfile"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"displayName"}},{"kind":"Field","name":{"kind":"Name","value":"profileImageUrl"}}]}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]}}]}}]} as unknown as DocumentNode<NotificationsByAdminQuery, NotificationsByAdminQueryVariables>;
export const CreateNotificationDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"CreateNotification"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"text"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"notificationsCreate"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"text"},"value":{"kind":"Variable","name":{"kind":"Name","value":"text"}}},{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}}]}}]} as unknown as DocumentNode<CreateNotificationMutation, CreateNotificationMutationVariables>;
export const DeleteNotificationDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"DeleteNotification"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"notificationsDelete"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}]}]}}]} as unknown as DocumentNode<DeleteNotificationMutation, DeleteNotificationMutationVariables>;
export const UpdateNotificationsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UpdateNotifications"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"opts"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"NotificationUpdateOpts"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"notificationsUpdate"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}},{"kind":"Argument","name":{"kind":"Name","value":"opts"},"value":{"kind":"Variable","name":{"kind":"Name","value":"opts"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}}]}}]} as unknown as DocumentNode<UpdateNotificationsMutation, UpdateNotificationsMutationVariables>;
export const UsersGetAllDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"UsersGetAll"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"opts"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"TwirUsersSearchParams"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"twirUsers"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"opts"},"value":{"kind":"Variable","name":{"kind":"Name","value":"opts"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"users"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"isBanned"}},{"kind":"Field","name":{"kind":"Name","value":"isBotAdmin"}},{"kind":"Field","name":{"kind":"Name","value":"isBotEnabled"}},{"kind":"Field","name":{"kind":"Name","value":"twitchProfile"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"login"}},{"kind":"Field","name":{"kind":"Name","value":"displayName"}},{"kind":"Field","name":{"kind":"Name","value":"profileImageUrl"}}]}}]}}]}}]}}]} as unknown as DocumentNode<UsersGetAllQuery, UsersGetAllQueryVariables>;
export const UserSwitchBanDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UserSwitchBan"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"switchUserBan"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}}]}]}}]} as unknown as DocumentNode<UserSwitchBanMutation, UserSwitchBanMutationVariables>;
export const UserSwitchAdminDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UserSwitchAdmin"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"switchUserAdmin"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}}]}]}}]} as unknown as DocumentNode<UserSwitchAdminMutation, UserSwitchAdminMutationVariables>;
export const DashboardStatsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"subscription","name":{"kind":"Name","value":"dashboardStats"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"dashboardStats"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"categoryId"}},{"kind":"Field","name":{"kind":"Name","value":"categoryName"}},{"kind":"Field","name":{"kind":"Name","value":"viewers"}},{"kind":"Field","name":{"kind":"Name","value":"startedAt"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"chatMessages"}},{"kind":"Field","name":{"kind":"Name","value":"followers"}},{"kind":"Field","name":{"kind":"Name","value":"usedEmotes"}},{"kind":"Field","name":{"kind":"Name","value":"requestedSongs"}},{"kind":"Field","name":{"kind":"Name","value":"subs"}}]}}]}}]} as unknown as DocumentNode<DashboardStatsSubscription, DashboardStatsSubscriptionVariables>;
export const GetAllTimersDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetAllTimers"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"timers"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"enabled"}},{"kind":"Field","name":{"kind":"Name","value":"timeInterval"}},{"kind":"Field","name":{"kind":"Name","value":"messageInterval"}},{"kind":"Field","name":{"kind":"Name","value":"responses"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"text"}},{"kind":"Field","name":{"kind":"Name","value":"isAnnounce"}}]}}]}}]}}]} as unknown as DocumentNode<GetAllTimersQuery, GetAllTimersQueryVariables>;
export const CreateTimerDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"CreateTimer"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"opts"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"TimerCreateInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"timersCreate"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"opts"},"value":{"kind":"Variable","name":{"kind":"Name","value":"opts"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}}]}}]} as unknown as DocumentNode<CreateTimerMutation, CreateTimerMutationVariables>;
export const UpdateTimerDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UpdateTimer"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"opts"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"TimerUpdateInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"timersUpdate"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}},{"kind":"Argument","name":{"kind":"Name","value":"opts"},"value":{"kind":"Variable","name":{"kind":"Name","value":"opts"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}}]}}]} as unknown as DocumentNode<UpdateTimerMutation, UpdateTimerMutationVariables>;
export const RemoveTimerDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"RemoveTimer"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"timersRemove"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}]}]}}]} as unknown as DocumentNode<RemoveTimerMutation, RemoveTimerMutationVariables>;