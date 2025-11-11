import { Injectable, Module, OnApplicationBootstrap } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { TypeOrmModule, InjectConnection } from '@nestjs/typeorm';
import { JwtModule } from '@nestjs/jwt';
import { Connection } from 'typeorm';

import { AuthController } from './auth.controller';
import { AuthService } from './auth.service';
import {
  Board,
  Column,
  DatabaseService,
  RefreshToken,
  Task,
  Team,
  TeamMember,
  User,
} from '@libs/database';

@Injectable()
class TypeormReadyService implements OnApplicationBootstrap {
  constructor(@InjectConnection() private readonly connection: Connection) {}

  async onApplicationBootstrap() {
    if (!this.connection.isInitialized) {
      await this.connection.initialize();
    }

    this.connection.entityMetadatas.forEach((meta) => {
      this.connection.getMetadata(meta.name);
    });

    console.log(
      'TypeORM полностью готов в auth-service – принимаем запросы от gateway',
    );
  }
}

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      envFilePath: '.env',
    }),

    TypeOrmModule.forRootAsync({
      useFactory: () => ({
        type: 'postgres' as const,
        host: process.env.DB_HOST || 'postgres',
        port: parseInt(process.env.DB_PORT || '5432', 10),
        username: process.env.DB_USER || 'kanban',
        password: process.env.DB_PASSWORD || 'kanban123',
        database: process.env.DB_NAME || 'kanban',
        entities: [User, RefreshToken, Team, TeamMember, Board, Column, Task],
        synchronize: process.env.NODE_ENV !== 'production',
        logging: process.env.NODE_ENV === 'development',
        retryAttempts: 5,
        retryDelay: 3000,
      }),
    }),

    JwtModule.register({
      secret: process.env.JWT_SECRET || 'fallback-super-secret-2025',
      signOptions: { expiresIn: '15m' },
    }),
  ],
  controllers: [AuthController],
  providers: [AuthService, DatabaseService, TypeormReadyService],
})
export class AuthServiceModule {}
