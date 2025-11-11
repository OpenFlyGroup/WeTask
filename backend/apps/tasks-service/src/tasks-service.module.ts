import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { TypeOrmModule } from '@nestjs/typeorm';
import { MongooseModule } from '@nestjs/mongoose';
import { TasksController } from './tasks.controller';
import { TasksService } from './tasks.service';
import { DatabaseService } from '@libs/database';
import { User, RefreshToken, Team, TeamMember, Board, Column, Task } from '@libs/database/entities';
import { Comment, CommentSchema } from './schemas/comment.schema';
import { ActivityLog, ActivityLogSchema } from './schemas/activity-log.schema';

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
      entities: [__dirname + '/../../../../libs/database/src/entities/**/*.entity{.ts,.js}'],
      synchronize: process.env.NODE_ENV !== 'production',
    }),
    TypeOrmModule.forFeature([User, RefreshToken, Team, TeamMember, Board, Column, Task]),
    MongooseModule.forRoot(
      process.env.MONGODB_URI || 'mongodb://mongodb:27017/kanban',
    ),
    MongooseModule.forFeature([
      { name: Comment.name, schema: CommentSchema },
      { name: ActivityLog.name, schema: ActivityLogSchema },
    ]),
  ],
  controllers: [TasksController],
  providers: [TasksService, DatabaseService],
})
export class TasksServiceModule {}
