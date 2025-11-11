import {
  Entity,
  PrimaryGeneratedColumn,
  Column as TypeOrmColumn,
  CreateDateColumn,
  UpdateDateColumn,
  OneToMany,
} from 'typeorm';
import { RefreshToken } from './refresh-token.entity';
import { TeamMember } from './team-member.entity';
import { Task } from './task.entity';

@Entity('users')
export class User {
  @PrimaryGeneratedColumn()
  id: number;

  @TypeOrmColumn({ unique: true })
  email: string;

  @TypeOrmColumn()
  password: string;

  @TypeOrmColumn()
  name: string;

  @CreateDateColumn()
  createdAt: Date;

  @UpdateDateColumn()
  updatedAt: Date;

  @OneToMany(() => RefreshToken, (refreshToken) => refreshToken.user)
  refreshTokens: RefreshToken[];

  @OneToMany(() => TeamMember, (teamMember) => teamMember.user)
  teams: TeamMember[];

  @OneToMany(() => Task, (task) => task.user)
  tasks: Task[];
}
