import { ChannelEvent, EventType } from '@tsuwari/typeorm/entities/ChannelEvent';
import { ChannelDonationEvent } from '@tsuwari/typeorm/entities/channelEvents/Donation';
import { ChannelIntegration } from '@tsuwari/typeorm/entities/ChannelIntegration';
import Centrifuge from 'centrifuge';
import ws from 'ws';
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import { XMLHttpRequest } from 'xmlhttprequest';

import { donatePayStore, typeorm } from '../index.js';
import { sendMessage } from '../libs/sender.js';
import { DonationAlerts } from './donationAlerts';

global.XMLHttpRequest = XMLHttpRequest;

type Event = {
  data: {
    notification: {
      type: 'donation',
      vars: {
        name: string,
        comment: string,
        sum: number,
        currency: 'string'
      }
    }
  }
}

export class DonatePay {
  #centrifuge: Centrifuge;
  #subscription: Centrifuge.Subscription;

  constructor(private readonly twitchUserId: string, private readonly apiKey: string) {}

  async connect() {
    const userData = await this.#getUserData();

    this.#centrifuge = new Centrifuge('wss://centrifugo.donatepay.ru:43002/connection/websocket', {
      subscribeEndpoint: 'https://donatepay.ru/api/v2/socket/token',
      subscribeParams:   {
        access_token: this.apiKey,
      },
      disableWithCredentials: true,
      websocket: ws,
    });

    this.#centrifuge.setToken(userData.token);

    this.#subscription = this.#centrifuge.subscribe(`$public:${userData.id}`, async (message: Event) => {
      if (message.data.notification.type !== 'donation') return;

      const { vars } = message.data.notification;

      const event = await typeorm.getRepository(ChannelEvent).save({
        channelId: this.twitchUserId,
        type: EventType.DONATION,
      });

      await typeorm.getRepository(ChannelDonationEvent).save({
        event,
        amount: vars.sum,
        currency: vars.currency,
        toUserId: this.twitchUserId,
        message: vars.comment,
        username: vars.name,
      });

      const msg = vars.comment || '';
      await sendMessage({
        channelId: this.twitchUserId,
        message: `${vars.name ?? 'Anonymous'}: ${vars.sum}${vars.currency} ${msg}`,
        color: 'orange',
      });
    });

    this.#centrifuge.on('connect', () => {
      console.info(`DonatePay: connected to channel ${this.twitchUserId}`);
    });

    this.#centrifuge.connect();
  }

  async disconnect() {
    await this.#subscription.unsubscribe();
    this.#centrifuge.disconnect();
  }

  async #getUserId() {
    const getUserParams = new URLSearchParams({
      access_token: this.apiKey,
    });
    const req = await fetch(`https://donatepay.ru/api/v1/user?${getUserParams}`);

    if (!req.ok) {
      throw new Error('incorrect response');
    }

    const data = await req.json();

    if (!data.data?.id) {
      throw new Error('incorrect response');
    }

    return data.data.id;
  }

  async #getUserData() {
    const userId = await this.#getUserId().catch(() => null);

    if (!userId) {
      console.error(`DonatePay: something wen't wrong when getting token of ${this.twitchUserId}`);
    }

    const req = await fetch('https://donatepay.ru/api/v2/socket/token', {
      method: 'post',
      body: JSON.stringify({
        access_token: this.apiKey,
      }),
      headers: {
        'Content-Type': 'application/json',
      },
    });
    const data = await req.json();

    return {
      token: data.token,
      id: userId,
    };
  }
}

export async function addDonatePayIntegration(integration: ChannelIntegration) {
  if (
    !integration.integration ||
    !integration.apiKey
  ) {
    return;
  }

  if (donatePayStore.get(integration.channelId)) {
    return;
  }

  const instance = new DonatePay(integration.channelId, integration.apiKey);
  await instance.connect();

  return instance;
}
