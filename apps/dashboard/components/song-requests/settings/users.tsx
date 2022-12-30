import { ActionIcon, Card, Divider, Flex, NumberInput, Popover, ScrollArea, Text, TextInput } from '@mantine/core';
import { IconCheck, IconPlus, IconX } from '@tabler/icons';
import React, { useState } from 'react';

import { useYouTubeSettingsFormContext } from '@/components/song-requests/settings/form';
import { YouTubeSettingsListButtonButton } from '@/components/song-requests/settings/listButton';

export const YouTubeUsersSettings: React.FC = () => {
  const form = useYouTubeSettingsFormContext();

  const [newDenyUser, setNewDenyUser] = useState('');

  function insertDenyUser() {
    form.insertListItem('denyList.users', { userName: newDenyUser, id: '' });
    setNewDenyUser('');
  }

  const [filterUsers, setFilterUsers] = useState('');

  return (
    <Card>
      <Card.Section p={'xs'}><Text>Users</Text></Card.Section>
      <Divider/>
      <Card.Section p={'md'}>
        <NumberInput label="Maximum songs by user in queue" {...form.getInputProps('user.maxRequests')} />
        <NumberInput
          label="Minimal watch time of user for request song (minutes)" {...form.getInputProps('user.minWatchTime')} />
        <NumberInput label="Minimal messages by user for request song" {...form.getInputProps('user.minMessages')} />
        <NumberInput
          label="Minimal follow time for request song (minutes)" {...form.getInputProps('user.minFollowTime')} />

        <Divider style={{ marginTop: 10 }}/>

        <Flex direction="row" justify="space-between" style={{ marginTop: 10 }}>
          <Text size="sm">Denied users for request</Text>
          <Popover width={200} position="bottom" withArrow shadow="md">
            <Popover.Target>
              <ActionIcon
                color="green"
                size={'sm'}
              ><IconPlus/></ActionIcon>
            </Popover.Target>
            <Popover.Dropdown>
              <Flex direction={'row'} gap={'sm'}>
                <TextInput placeholder="enter username" onChange={(v) => setNewDenyUser(v.currentTarget.value)}/>
                <ActionIcon onClick={() => insertDenyUser()}><IconCheck/></ActionIcon>
              </Flex>
            </Popover.Dropdown>
          </Popover>
        </Flex>

        {form.values.denyList.users?.length
          ? <TextInput
            style={{ marginTop: 10 }}
            placeholder="filter..."
            onChange={(v) => setFilterUsers(v.target.value)}
          />
          : ''
        }


        <ScrollArea type={'always'} style={{ marginTop: 10 }}>
          <Flex direction={'column'} style={{ maxHeight: 300 }}>
            {form.values.denyList.users.length
              ? form.values.denyList.users.map((u, i) => <YouTubeSettingsListButtonButton
                key={i}
                text={u.userName}
                onClick={() => form.removeListItem('denyList.users', i)}
                icon={IconX}
              />)
              : ''
            }
          </Flex>
        </ScrollArea>

      </Card.Section>
    </Card>
  );
};