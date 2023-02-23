import { Column, Entity, OneToMany, PrimaryGeneratedColumn } from 'typeorm';

// eslint-disable-next-line import/no-cycle
import { ChannelRolePermission } from './ChannelRolePermission';

enum RolePermissionEnum {
  ADMINISTRATOR = 'ADMINISTRATOR',

  UPDATE_CHANNEL_TITLE = 'UPDATE_CHANNEL_TITLE',
  UPDATE_CHANNEL_CATEGORY = 'UPDATE_CHANNEL_CATEGORY',

  VIEW_COMMANDS = 'VIEW_COMMANDS',
  MANAGE_COMMANDS = 'MANAGE_COMMANDS',

  VIEW_KEYWORDS = 'VIEW_KEYWORDS',
  MANAGE_KEYWORDS = 'MANAGE_KEYWORDS',

  VIEW_TIMERS = 'VIEW_TIMERS',
  MANAGE_TIMERS = 'MANAGE_TIMERS',

  VIEW_INTEGRATIONS = 'VIEW_INTEGRATIONS',
  MANAGE_INTEGRATIONS = 'MANAGE_INTEGRATIONS',

  VIEW_SONG_REQUESTS = 'VIEW_SONG_REQUESTS',
  MANAGE_SONG_REQUESTS = 'MANAGE_SONG_REQUESTS',

  VIEW_MODERATION = 'VIEW_MODERATION',
  MANAGE_MODERATION = 'MANAGE_MODERATION',

  VIEW_VARIABLES = 'VIEW_VARIABLES',
  MANAGE_VARIABLES = 'MANAGE_VARIABLES',

  VIEW_GREETINGS = 'VIEW_GREETINGS',
  MANAGE_GREETINGS = 'MANAGE_GREETINGS',
}

@Entity('roles_permissions')
export class RolePermission {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column('enum', { enum: RolePermissionEnum, default: RolePermissionEnum.ADMINISTRATOR, unique: true })
  permission: RolePermissionEnum;

  @OneToMany(() => ChannelRolePermission, _ => _.permission)
  channelRolePermissions: ChannelRolePermission[];
}