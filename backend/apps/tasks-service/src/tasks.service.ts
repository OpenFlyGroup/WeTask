import { Injectable } from '@nestjs/common';
import { InjectModel } from '@nestjs/mongoose';
import { Model } from 'mongoose';
import { DatabaseService } from '@libs/database';
import { Comment } from './schemas/comment.schema';
import { ActivityLog } from './schemas/activity-log.schema';

@Injectable()
export class TasksService {
  constructor(
    private database: DatabaseService,
    @InjectModel(Comment.name) private commentModel: Model<Comment>,
    @InjectModel(ActivityLog.name) private activityLogModel: Model<ActivityLog>,
  ) {}

  async create(data: {
    title: string;
    description?: string;
    columnId: number;
    assignedTo?: number;
    priority?: string;
  }) {
    const task = await this.database.task.save({
      title: data.title,
      description: data.description,
      columnId: data.columnId,
      assignedTo: data.assignedTo,
      priority: data.priority || 'medium',
      status: 'todo',
    });

    const taskWithRelations = await this.database.task.findOne({
      where: { id: task.id },
      relations: ['user', 'column', 'column.board'],
    });

    // Log activity
    await this.activityLogModel.create({
      entityType: 'task',
      entityId: task.id,
      action: 'created',
      userId: data.assignedTo || 0,
      metadata: { title: task.title },
    });

    return {
      success: true,
      data: taskWithRelations,
    };
  }

  async getById(id: number) {
    const task = await this.database.task.findOne({
      where: { id },
      relations: ['user', 'column', 'column.board', 'column.board.team'],
    });

    if (!task) {
      return {
        success: false,
        error: 'Task not found',
        statusCode: 404,
      };
    }

    return {
      success: true,
      data: task,
    };
  }

  async getByBoard(boardId: number) {
    const board = await this.database.board.findOne({
      where: { id: boardId },
      relations: ['columns', 'columns.tasks', 'columns.tasks.user'],
    });

    if (board?.columns) {
      board.columns.sort((a, b) => a.order - b.order);
    }

    if (!board) {
      return {
        success: false,
        error: 'Board not found',
        statusCode: 404,
      };
    }

    return {
      success: true,
      data: board.columns.flatMap((col) => col.tasks),
    };
  }

  async update(
    id: number,
    data: {
      title?: string;
      description?: string;
      priority?: string;
      assignedTo?: number;
    },
  ) {
    const task = await this.database.task.findOne({
      where: { id },
    });

    if (!task) {
      return {
        success: false,
        error: 'Task not found',
        statusCode: 404,
      };
    }

    if (data.title) task.title = data.title;
    if (data.description !== undefined) task.description = data.description;
    if (data.priority) task.priority = data.priority;
    if (data.assignedTo !== undefined) task.assignedTo = data.assignedTo;

    const updatedTask = await this.database.task.save(task);

    const taskWithUser = await this.database.task.findOne({
      where: { id: updatedTask.id },
      relations: ['user'],
    });

    // Log activity
    await this.activityLogModel.create({
      entityType: 'task',
      entityId: task.id,
      action: 'updated',
      userId: data.assignedTo || task.assignedTo || 0,
      metadata: {
        title: data.title,
        description: data.description,
        priority: data.priority,
        assignedTo: data.assignedTo,
      },
    });

    return {
      success: true,
      data: taskWithUser,
    };
  }

  async delete(id: number) {
    const task = await this.database.task.findOne({
      where: { id },
    });

    if (!task) {
      return {
        success: false,
        error: 'Task not found',
        statusCode: 404,
      };
    }

    await this.database.task.remove(task);

    // Log activity
    await this.activityLogModel.create({
      entityType: 'task',
      entityId: id,
      action: 'deleted',
      userId: task.assignedTo || 0,
    });

    return {
      success: true,
      data: { message: 'Task deleted successfully' },
    };
  }

  async move(id: number, columnId: number) {
    const task = await this.database.task.findOne({
      where: { id },
    });

    if (!task) {
      return {
        success: false,
        error: 'Task not found',
        statusCode: 404,
      };
    }

    task.columnId = columnId;
    const updatedTask = await this.database.task.save(task);

    const taskWithRelations = await this.database.task.findOne({
      where: { id: updatedTask.id },
      relations: ['user', 'column', 'column.board'],
    });

    // Log activity
    await this.activityLogModel.create({
      entityType: 'task',
      entityId: task.id,
      action: 'moved',
      userId: task.assignedTo || 0,
      metadata: { fromColumnId: task.columnId, toColumnId: columnId },
    });

    return {
      success: true,
      data: taskWithRelations,
    };
  }

  async addComment(data: { taskId: number; userId: number; message: string }) {
    const comment = await this.commentModel.create({
      taskId: data.taskId,
      userId: data.userId,
      message: data.message,
    });

    // Log activity
    await this.activityLogModel.create({
      entityType: 'task',
      entityId: data.taskId,
      action: 'comment_added',
      userId: data.userId,
    });

    return {
      success: true,
      data: comment,
    };
  }

  async getComments(taskId: number) {
    const comments = await this.commentModel
      .find({ taskId })
      .sort({ createdAt: -1 })
      .exec();

    return {
      success: true,
      data: comments,
    };
  }
}
