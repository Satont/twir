import { Inject, Injectable, Logger } from '@nestjs/common';
import { Interval } from '@nestjs/schedule';
import { config } from '@tsuwari/config';
import { PrismaService } from '@tsuwari/prisma';
import { ClientProxy } from '@tsuwari/shared';
import _ from 'lodash';

@Injectable()
export class DotaService {
  #logger = new Logger(DotaService.name);

  constructor(private readonly prisma: PrismaService, @Inject('NATS') private nats: ClientProxy) { }

  @Interval('dota', config.isDev ? 10000 : 1 * 60 * 1000)
  async cacheDota() {
    const accounts = await this.prisma.dotaAccount.findMany({
      where: {
        channel: {
          isEnabled: true,
        },
      },
    });
    const chunks = _.chunk(accounts.map(a => a.id), 50);
    this.#logger.log(`Getting information about ${accounts.length} accs.`);

    for (const chunk of chunks) {
      await this.nats.emit('dota.cacheAccountsMatches', chunk).toPromise();
    }
  }
}
