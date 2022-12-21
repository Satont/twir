import { Group, Avatar, Text, Button, Flex, Alert } from '@mantine/core';
import { IconBrandSpotify, IconLogin, IconLogout } from '@tabler/icons';

import { IntegrationCard } from './card';

import { useSpotifyIntegration } from '@/services/api/integrations';
import {useTranslation} from "next-i18next";

export const SpotifyIntegration: React.FC = () => {
  const manager = useSpotifyIntegration();
  const { t } = useTranslation("integrations")

  // const { data: integration } = manager.getIntegration();
  const { data: profile } = manager.getProfile();

  async function login() {
    const link = await manager.getAuthLink();
    if (link) {
      window.location.replace(link);
    }
  }

  return (
    <IntegrationCard
      title="Spotify"
      icon={IconBrandSpotify}
      iconColor="green"
      header={
        <Flex direction="row" gap="sm">
          {profile && (
            <Button
              compact
              leftIcon={<IconLogout />}
              variant="outline"
              color="red"
              onClick={manager.logout}
            >
              {t("logout")}
            </Button>
          )}
          <Button compact leftIcon={<IconLogin />} variant="outline" color="green" onClick={login}>
            {t("login")}
          </Button>
        </Flex>
      }
    >
      {!profile && <Alert>{t("notLoggedIn")}</Alert>}
      {profile && (
        <Group position="apart" mt={10}>
          <Text weight={500} size={30}>
            Satont WorldWide
          </Text>
          {profile.images && (
            <Avatar src={profile.images[0].url} h={150} w={150} style={{ borderRadius: 900 }} />
          )}
        </Group>
      )}
    </IntegrationCard>
  );
};
