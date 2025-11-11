import { NestFactory } from '@nestjs/core';
import { MicroserviceOptions, Transport } from '@nestjs/microservices';
import { AuthServiceModule } from './auth-service.module';

async function bootstrap() {
  const app = await NestFactory.createMicroservice<MicroserviceOptions>(
    AuthServiceModule,
    {
      transport: Transport.RMQ,
      options: {
        urls: [process.env.RABBITMQ_URL || 'amqp://localhost:5672'],
        queue: 'auth_queue',
        queueOptions: {
          durable: true,
        },
        noAck: false,
      },
    },
  );

  await app.listen();
  await new Promise((resolve) => setTimeout(resolve, 500));
  console.log('Auth Service is listening on RabbitMQ');
}
bootstrap();
