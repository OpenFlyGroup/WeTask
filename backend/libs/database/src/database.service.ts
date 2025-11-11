import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository, DataSource } from 'typeorm';
import { User } from './entities/user.entity';
import { RefreshToken } from './entities/refresh-token.entity';
import { Team } from './entities/team.entity';
import { TeamMember } from './entities/team-member.entity';
import { Board } from './entities/board.entity';
import { Column } from './entities/column.entity';
import { Task } from './entities/task.entity';

@Injectable()
export class DatabaseService {
  constructor(
    private dataSource: DataSource,
    @InjectRepository(User)
    public user: Repository<User>,
    @InjectRepository(RefreshToken)
    public refreshToken: Repository<RefreshToken>,
    @InjectRepository(Team)
    public team: Repository<Team>,
    @InjectRepository(TeamMember)
    public teamMember: Repository<TeamMember>,
    @InjectRepository(Board)
    public board: Repository<Board>,
    @InjectRepository(Column)
    public column: Repository<Column>,
    @InjectRepository(Task)
    public task: Repository<Task>,
  ) {}
}
