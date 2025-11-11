import { Injectable, OnModuleInit } from '@nestjs/common';
import { DataSource, Repository } from 'typeorm';

import { User } from './entities/user.entity';
import { Team } from './entities/team.entity';
import { TeamMember } from './entities/team-member.entity';
import { Board } from './entities/board.entity';
import { Column } from './entities/column.entity';
import { Task } from './entities/task.entity';
import { RefreshToken } from './entities';

@Injectable()
export class DatabaseService implements OnModuleInit {
  public user!: Repository<User>;
  public refreshToken!: Repository<RefreshToken>;
  public team!: Repository<Team>;
  public teamMember!: Repository<TeamMember>;
  public board!: Repository<Board>;
  public column!: Repository<Column>;
  public task!: Repository<Task>;

  constructor(private readonly dataSource: DataSource) {}

  async onModuleInit(): Promise<void> {
    if (!this.dataSource.isInitialized) {
      await this.dataSource.initialize();
    }

    const entities = [
      User,
      RefreshToken,
      Team,
      TeamMember,
      Board,
      Column,
      Task,
    ];

    for (const entity of entities) {
      try {
        this.dataSource.getRepository(entity).target;
        this.dataSource.getRepository(entity).metadata;
      } catch (err) {
        console.error(`Failed to load metadata for ${entity.name}`, err);
        await new Promise((r) => setTimeout(r, 100));
        this.dataSource.getRepository(entity).target;
      }
    }

    this.user = this.dataSource.getRepository(User);
    this.refreshToken = this.dataSource.getRepository(RefreshToken);
    this.team = this.dataSource.getRepository(Team);
    this.teamMember = this.dataSource.getRepository(TeamMember);
    this.board = this.dataSource.getRepository(Board);
    this.column = this.dataSource.getRepository(Column);
    this.task = this.dataSource.getRepository(Task);

    console.log(
      'DatabaseService: метаданные 100% построены, репозитории готовы',
    );
  }
}
