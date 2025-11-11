import { Controller } from '@nestjs/common';
import { MessagePattern, Payload, EventPattern } from '@nestjs/microservices';
import { RabbitMQPatterns } from '@libs/common';
import { TasksService } from './tasks.service';

@Controller()
export class TasksController {
  constructor(private readonly tasksService: TasksService) {}

  @MessagePattern(RabbitMQPatterns.TASKS_CREATE)
  async create(
    @Payload()
    data: {
      title: string;
      description?: string;
      columnId: number;
      assignedTo?: number;
      priority?: string;
    },
  ) {
    return this.tasksService.create(data);
  }

  @MessagePattern(RabbitMQPatterns.TASKS_GET_BY_ID)
  async getById(@Payload() data: { id: number }) {
    return this.tasksService.getById(data.id);
  }

  @MessagePattern(RabbitMQPatterns.TASKS_GET_BY_BOARD)
  async getByBoard(@Payload() data: { boardId: number }) {
    return this.tasksService.getByBoard(data.boardId);
  }

  @MessagePattern(RabbitMQPatterns.TASKS_UPDATE)
  async update(
    @Payload()
    data: {
      id: number;
      title?: string;
      description?: string;
      priority?: string;
      assignedTo?: number;
    },
  ) {
    return this.tasksService.update(data.id, data);
  }

  @MessagePattern(RabbitMQPatterns.TASKS_DELETE)
  async delete(@Payload() data: { id: number }) {
    return this.tasksService.delete(data.id);
  }

  @MessagePattern(RabbitMQPatterns.TASKS_MOVE)
  async move(@Payload() data: { id: number; columnId: number }) {
    return this.tasksService.move(data.id, data.columnId);
  }

  @MessagePattern(RabbitMQPatterns.TASKS_ADD_COMMENT)
  async addComment(
    @Payload() data: { taskId: number; userId: number; message: string },
  ) {
    return this.tasksService.addComment(data);
  }

  @MessagePattern(RabbitMQPatterns.TASKS_GET_COMMENTS)
  async getComments(@Payload() data: { taskId: number }) {
    return this.tasksService.getComments(data.taskId);
  }
}
