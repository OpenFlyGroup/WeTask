import { Module } from '@nestjs/common';
import { ClientsModule, Transport } from '@nestjs/microservices';
import { BoardsController } from './boards.controller';
import { ColumnsController } from './columns.controller';

@Module({
  imports: [
    ClientsModule.register([
      {
        name: 'BOARDS_SERVICE',
        transport: Transport.RMQ,
        options: {
          urls: [process.env.RABBITMQ_URL || 'amqp://rabbitmq:5672'],
          queue: 'boards_queue',
          queueOptions: {
            durable: true,
          },
        },
      },
    ]),
  ],
  controllers: [BoardsController, ColumnsController],
})
export class BoardsModule {}
