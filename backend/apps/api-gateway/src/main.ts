import { NestFactory } from '@nestjs/core';
import { ValidationPipe } from '@nestjs/common';
import { SwaggerModule, DocumentBuilder } from '@nestjs/swagger';
import { NestExpressApplication } from '@nestjs/platform-express';
import { ApiGatewayModule } from './api-gateway.module';

async function bootstrap() {
  const app = await NestFactory.create<NestExpressApplication>(
    ApiGatewayModule,
    { cors: true },
  );

  app.setGlobalPrefix('api');
  app.useGlobalPipes(
    new ValidationPipe({
      whitelist: true,
      transform: true,
    }),
  );

  // app.enableCors({
  //   origin: true,
  //   credentials: true,
  // });

  // Swagger/OpenAPI configuration
  const config = new DocumentBuilder()
    .setTitle('Task Tracker API')
    .setDescription(
      'Микросервисный backend для таск-трекера. API документация с полным описанием всех эндпоинтов.',
    )
    .setVersion('1.0')
    .addBearerAuth(
      {
        type: 'http',
        scheme: 'bearer',
        bearerFormat: 'JWT',
        name: 'JWT',
        description: 'Enter JWT token',
        in: 'header',
      },
      'JWT-auth',
    )
    .addTag('auth', 'Аутентификация и авторизация')
    .addTag('users', 'Управление пользователями')
    .addTag('teams', 'Управление командами')
    .addTag('boards', 'Управление досками')
    .addTag('columns', 'Управление колонками')
    .addTag('tasks', 'Управление задачами')
    .build();

  const document = SwaggerModule.createDocument(app, config);
  // Serve Swagger UI assets locally from node_modules (works inside Docker without Internet)
  // eslint-disable-next-line @typescript-eslint/no-var-requires
  const swaggerDistPath: string = require('swagger-ui-dist').absolutePath();
  app.useStaticAssets(swaggerDistPath, { prefix: '/swagger-ui/' });
  SwaggerModule.setup('api/docs', app, document, {
    customCssUrl: '/swagger-ui/swagger-ui.css',
    customfavIcon: '/swagger-ui/favicon-32x32.png',
    customJs: [
      '/swagger-ui/swagger-ui-bundle.js',
      '/swagger-ui/swagger-ui-standalone-preset.js',
    ],
  });

  const port = process.env.PORT || 3000;
  await app.listen(port);
  console.log(`API Gateway is running on http://localhost:${port}/api`);
  console.log(`WebSocket is available on ws://localhost:${port}`);
  console.log(`Swagger documentation: http://localhost:${port}/api/docs`);
}
void bootstrap();
