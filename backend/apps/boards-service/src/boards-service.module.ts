import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { TypeOrmModule } from '@nestjs/typeorm';
import { BoardsController } from './boards.controller';
import { BoardsService } from './boards.service';
import { ColumnsController } from './columns.controller';
import { ColumnsService } from './columns.service';
import { DatabaseService } from '@libs/database';

@Module({
  imports: [
    ConfigModule.forRoot({ isGlobal: true }),
    TypeOrmModule.forRoot({
      type: 'postgres',
      host: process.env.DB_HOST || 'postgres',
      port: parseInt(process.env.DB_PORT || '5432'),
      username: process.env.DB_USER || 'kanban',
      password: process.env.DB_PASSWORD || 'kanban123',
      database: process.env.DB_NAME || 'kanban',
      autoLoadEntities: true,
      entities: [
        __dirname +
          '/../../../../libs/database/src/entities/**/*.entity{.ts,.js}',
      ],
      synchronize: process.env.NODE_ENV !== 'production',
    }),
  ],
  controllers: [BoardsController, ColumnsController],
  providers: [BoardsService, ColumnsService, DatabaseService],
})
export class BoardsServiceModule {}
