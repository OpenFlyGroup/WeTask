import {
  Entity,
  PrimaryGeneratedColumn,
  Column as TypeOrmColumn,
  CreateDateColumn,
  UpdateDateColumn,
  ManyToOne,
  OneToMany,
  JoinColumn,
} from 'typeorm';
import { Team } from './team.entity';
import { Column as ColumnEntity } from './column.entity';

@Entity('boards')
export class Board {
  @PrimaryGeneratedColumn()
  id: number;

  @TypeOrmColumn()
  title: string;

  @TypeOrmColumn()
  teamId: number;

  @CreateDateColumn()
  createdAt: Date;

  @UpdateDateColumn()
  updatedAt: Date;

  @ManyToOne(() => Team, (team) => team.boards, { onDelete: 'CASCADE' })
  @JoinColumn({ name: 'teamId' })
  team: Team;

  @OneToMany(() => ColumnEntity, (column) => column.board)
  columns: ColumnEntity[];
}
