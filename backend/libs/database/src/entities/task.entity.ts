import {
  Entity,
  PrimaryGeneratedColumn,
  Column as TypeOrmColumn,
  CreateDateColumn,
  UpdateDateColumn,
  ManyToOne,
  JoinColumn,
} from 'typeorm';
import { Column } from './column.entity';
import { User } from './user.entity';

@Entity('tasks')
export class Task {
  @PrimaryGeneratedColumn()
  id: number;

  @TypeOrmColumn()
  title: string;

  @TypeOrmColumn({ type: 'text', nullable: true })
  description: string | null;

  @TypeOrmColumn()
  status: string;

  @TypeOrmColumn({ nullable: true, default: 'medium' })
  priority: string | null;

  @TypeOrmColumn()
  columnId: number;

  @TypeOrmColumn({ nullable: true })
  assignedTo: number | null;

  @CreateDateColumn()
  createdAt: Date;

  @UpdateDateColumn()
  updatedAt: Date;

  @ManyToOne(() => Column, (column) => column.tasks, { onDelete: 'CASCADE' })
  @JoinColumn({ name: 'columnId' })
  column: Column;

  @ManyToOne(() => User, (user) => user.tasks, { onDelete: 'SET NULL' })
  @JoinColumn({ name: 'assignedTo' })
  user: User | null;
}
