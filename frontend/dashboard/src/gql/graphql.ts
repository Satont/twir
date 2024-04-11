/* eslint-disable */
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
  DateTime: { input: any; output: any; }
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

export type Mutation = {
  __typename?: 'Mutation';
  createCommand: Command;
  createNotification: Notification;
  removeCommand: Scalars['Boolean']['output'];
  updateCommand: Command;
};


export type MutationCreateCommandArgs = {
  opts: CreateCommandInput;
};


export type MutationCreateNotificationArgs = {
  text: Scalars['String']['input'];
  userId?: InputMaybe<Scalars['String']['input']>;
};


export type MutationRemoveCommandArgs = {
  id: Scalars['String']['input'];
};


export type MutationUpdateCommandArgs = {
  id: Scalars['String']['input'];
  opts: UpdateCommandOpts;
};

export type Notification = {
  __typename?: 'Notification';
  id: Scalars['ID']['output'];
  text: Scalars['String']['output'];
  userId: Scalars['ID']['output'];
};

export type Query = {
  __typename?: 'Query';
  authedUser: AuthenticatedUser;
  commands: Array<Command>;
  notifications: Array<Notification>;
  /** finding users on twitch with filter does they exists in database */
  twirUsers: TwirUsersResponse;
};


export type QueryNotificationsArgs = {
  userId: Scalars['String']['input'];
};


export type QueryTwirUsersArgs = {
  opts: TwirUsersSearchParams;
};

export type Subscription = {
  __typename?: 'Subscription';
  /** `newNotification` will return a stream of `Notification` objects. */
  newNotification: Notification;
};

export type TwirAdminUser = TwirUser & {
  __typename?: 'TwirAdminUser';
  apiKey: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  isBanned: Scalars['Boolean']['output'];
  isBotAdmin: Scalars['Boolean']['output'];
  isBotModerator: Scalars['Boolean']['output'];
  twitchProfile: TwirUserTwitchInfo;
};

export type TwirUser = {
  id: Scalars['ID']['output'];
  twitchProfile: TwirUserTwitchInfo;
};

export type TwirUserTwitchInfo = {
  __typename?: 'TwirUserTwitchInfo';
  description: Scalars['String']['output'];
  displayName: Scalars['String']['output'];
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
