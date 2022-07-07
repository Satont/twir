import { PrismaService } from '@tsuwari/prisma';
import { DotaGame, RedisService, TwitchApiService } from '@tsuwari/shared';

import { app } from '../../index.js';
import { DefaultCommand } from '../types.js';

const prisma = app.get(PrismaService);
const staticApi = app.get(TwitchApiService);
const redis = app.get(RedisService);

const messages = Object.freeze({
  GAME_NOT_FOUND: 'Game not found.',
});

const getGames = async (accounts: string[]) => {
  const rps = await Promise.all(accounts.map(a => redis.get(`dotaRps:${a}`)));
  if (!rps.filter(r => r !== null).length) {
    return messages.GAME_NOT_FOUND;
  }
  const cachedRps = rps.filter(r => r !== null).map(r => JSON.parse(r!));

  const cachedGames = await Promise.all(accounts.map(a => redis.get(`dotaMatches:${a}`)));
  if (!cachedGames.filter(r => r !== null).length) {
    return messages.GAME_NOT_FOUND;
  }

  const parsedGames = cachedGames.filter(g => g !== null).map(g => JSON.parse(g!) as DotaGame);
  const dbGames = await prisma.dotaMatch.findMany({
    where: {
      lobbyId: {
        in: cachedRps.map(r => r.lobbyId),
      },
      players: {
        hasSome: accounts.map(a => Number(a)),
      },
    },
    orderBy: {
      startedAt: 'desc',
    },
    include: {
      gameMode: true,
    },
    take: 2,
  });

  return dbGames.map(g => {
    const cachedGame = parsedGames.find(game => game.match_id === g.match_id)!;

    return {
      ...g,
      players: cachedGame.players,
    };
  });
};

export const dota: DefaultCommand[] = [{
  name: 'np',
  permission: 'VIEWER',
  visible: false,
  handler: async (state, params) => {
    if (!state.channelId) return;

    const accounts = await prisma.dotaAccount.findMany({
      where: {
        channelId: state.channelId,
      },
    });

    if (!accounts.length) return 'You have not added account.';

    const games = await getGames(accounts.map(a => a.id));

    if (typeof games === 'string') return games;
    if (!games.length) return messages.GAME_NOT_FOUND;

    console.dir(games, { depth: null });
    return games
      .map(g => {
        const avgMmr = g.gameMode.id === 22 ? ` (${g.avarage_mmr}mmr)` : '';
        return `${g.gameMode.name}${avgMmr}`;
      })
      .join(' | ');
  },
}];