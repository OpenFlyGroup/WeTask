import {
  Entity,
  PrimaryGeneratedColumn,
  Column as TypeOrmColumn,
  ManyToOne,
  JoinColumn,
  Unique,
} from 'typeorm';
import { Team } from './team.entity';
import { User } from './user.entity';

@Entity('team_members')
@Unique(['teamId', 'userId'])
export class TeamMember {
  @PrimaryGeneratedColumn()
  id: number;

  @TypeOrmColumn()
  teamId: number;

  @TypeOrmColumn()
  userId: number;

  @TypeOrmColumn({ default: 'member' })
  role: string; // owner, admin, member

  @ManyToOne(() => Team, (team) => team.members, { onDelete: 'CASCADE' })
  @JoinColumn({ name: 'teamId' })
  team: Team;

  @ManyToOne(() => User, (user) => user.teams, { onDelete: 'CASCADE' })
  @JoinColumn({ name: 'userId' })
  user: User;
}
