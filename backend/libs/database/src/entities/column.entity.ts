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
import { Board } from './board.entity';
import { Task } from './task.entity';

@Entity('columns')
export class Column {
  @PrimaryGeneratedColumn()
  id: number;

  @TypeOrmColumn()
  title: string;

  @TypeOrmColumn()
  order: number;

  @TypeOrmColumn()
  boardId: number;

  @CreateDateColumn()
  createdAt: Date;

  @UpdateDateColumn()
  updatedAt: Date;

  @ManyToOne(() => Board, (board) => board.columns, { onDelete: 'CASCADE' })
  @JoinColumn({ name: 'boardId' })
  board: Board;

  @OneToMany(() => Task, (task) => task.column)
  tasks: Task[];
}
