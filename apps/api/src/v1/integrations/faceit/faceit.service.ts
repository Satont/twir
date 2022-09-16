import { HttpException, Injectable } from '@nestjs/common';
import { ChannelIntegration } from '@tsuwari/typeorm/entities/ChannelIntegration';
import { Integration, IntegrationService } from '@tsuwari/typeorm/entities/Integration';

import { typeorm } from '../../../index.js';
import { FaceitUpdateDto } from './dto/update.js';

@Injectable()
export class FaceitService {
  async getIntegration(channelId: string) {
    const integration = await typeorm.getRepository(ChannelIntegration).findOne({
      where: { channelId, integration: { service: IntegrationService.FACEIT } },
    });

    return integration;
  }

  async updateIntegration(channelId: string, body: FaceitUpdateDto) {
    const integrationService = await typeorm.getRepository(Integration).findOneBy({
      service: IntegrationService.FACEIT,
    });

    if (!integrationService)
      throw new HttpException(
        `Faceit not enabled on our backed. Please, make patience or contact us`,
        404,
      );

    body.data.game = body.data.game ?? 'csgo';

    const repository = typeorm.getRepository(ChannelIntegration);
    let integration = await this.getIntegration(channelId);
    if (!integration) {
      const newIntegration = repository.create({
        channelId,
        enabled: body.enabled,
        data: { ...body.data },
        integrationId: integrationService.id,
      });
      await repository.save(newIntegration);
      integration = newIntegration;
    } else {
      await repository.update(
        { id: integration.id },
        {
          enabled: body.enabled,
          data: body.data as any,
        },
      );
    }

    return repository.findOneBy({ id: integration.id });
  }
}
