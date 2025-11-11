/* eslint-disable @typescript-eslint/no-unsafe-return */
/* eslint-disable @typescript-eslint/no-unsafe-argument */
/* eslint-disable @typescript-eslint/no-unsafe-member-access */
/* eslint-disable @typescript-eslint/no-unsafe-assignment */
import {
  Controller,
  Get,
  Post,
  Put,
  Delete,
  Param,
  Body,
  UseGuards,
  Request,
} from '@nestjs/common';
import { Inject } from '@nestjs/common';
import { ClientProxy } from '@nestjs/microservices';
import { AuthGuard } from '@nestjs/passport';
import {
  ApiTags,
  ApiOperation,
  ApiResponse,
  ApiBearerAuth,
  ApiParam,
  ApiBody,
} from '@nestjs/swagger';
import { RabbitMQPatterns } from '@libs/common';
import { firstValueFrom } from 'rxjs';
import {
  CreateTaskDto,
  UpdateTaskDto,
  MoveTaskDto,
  AddCommentDto,
} from '../common/dto/tasks.dto';
import { handleRPCError } from '../common/utils/error-handler';

@ApiTags('tasks')
@ApiBearerAuth('JWT-auth')
@Controller('tasks')
@UseGuards(AuthGuard('jwt'))
export class TasksController {
  constructor(@Inject('TASKS_SERVICE') private tasksClient: ClientProxy) {}

  @Post()
  @ApiOperation({ summary: 'Создать новую задачу' })
  @ApiBody({ type: CreateTaskDto })
  @ApiResponse({ status: 201, description: 'Задача успешно создана' })
  async create(@Body() data: CreateTaskDto) {
    const result = await firstValueFrom(
      this.tasksClient.send(RabbitMQPatterns.TASKS_CREATE, data),
    );

    if (!result.success) {
      handleRPCError(result);
    }

    return result.data;
  }

  @Get(':id')
  @ApiOperation({ summary: 'Получить задачу по ID' })
  @ApiParam({ name: 'id', description: 'ID задачи', type: 'number' })
  @ApiResponse({ status: 200, description: 'Информация о задаче' })
  @ApiResponse({ status: 404, description: 'Задача не найдена' })
  async getById(@Param('id') id: string) {
    const result = await firstValueFrom(
      this.tasksClient.send(RabbitMQPatterns.TASKS_GET_BY_ID, {
        id: parseInt(id),
      }),
    );

    if (!result.success) {
      handleRPCError(result);
    }

    return result.data;
  }

  @Get('board/:boardId')
  @ApiOperation({ summary: 'Получить все задачи доски' })
  @ApiParam({ name: 'boardId', description: 'ID доски', type: 'number' })
  @ApiResponse({ status: 200, description: 'Список задач доски' })
  @ApiResponse({ status: 404, description: 'Доска не найдена' })
  async getByBoard(@Param('boardId') boardId: string) {
    const result = await firstValueFrom(
      this.tasksClient.send(RabbitMQPatterns.TASKS_GET_BY_BOARD, {
        boardId: parseInt(boardId),
      }),
    );

    if (!result.success) {
      handleRPCError(result);
    }

    return result.data;
  }

  @Put(':id')
  @ApiOperation({ summary: 'Обновить задачу' })
  @ApiParam({ name: 'id', description: 'ID задачи', type: 'number' })
  @ApiBody({ type: UpdateTaskDto })
  @ApiResponse({ status: 200, description: 'Задача успешно обновлена' })
  @ApiResponse({ status: 404, description: 'Задача не найдена' })
  async update(@Param('id') id: string, @Body() data: UpdateTaskDto) {
    const result = await firstValueFrom(
      this.tasksClient.send(RabbitMQPatterns.TASKS_UPDATE, {
        id: parseInt(id),
        ...data,
      }),
    );

    if (!result.success) {
      handleRPCError(result);
    }

    return result.data;
  }

  @Delete(':id')
  @ApiOperation({ summary: 'Удалить задачу' })
  @ApiParam({ name: 'id', description: 'ID задачи', type: 'number' })
  @ApiResponse({ status: 200, description: 'Задача успешно удалена' })
  @ApiResponse({ status: 404, description: 'Задача не найдена' })
  async delete(@Param('id') id: string) {
    const result = await firstValueFrom(
      this.tasksClient.send(RabbitMQPatterns.TASKS_DELETE, {
        id: parseInt(id),
      }),
    );

    if (!result.success) {
      handleRPCError(result);
    }

    return result.data;
  }

  @Put(':id/move')
  @ApiOperation({ summary: 'Переместить задачу в другую колонку' })
  @ApiParam({ name: 'id', description: 'ID задачи', type: 'number' })
  @ApiBody({ type: MoveTaskDto })
  @ApiResponse({ status: 200, description: 'Задача успешно перемещена' })
  @ApiResponse({ status: 404, description: 'Задача не найдена' })
  async move(@Param('id') id: string, @Body() data: MoveTaskDto) {
    const result = await firstValueFrom(
      this.tasksClient.send(RabbitMQPatterns.TASKS_MOVE, {
        id: parseInt(id),
        columnId: data.columnId,
      }),
    );

    if (!result.success) {
      handleRPCError(result);
    }

    return result.data;
  }

  @Post(':id/comment')
  @ApiOperation({ summary: 'Добавить комментарий к задаче' })
  @ApiParam({ name: 'id', description: 'ID задачи', type: 'number' })
  @ApiBody({ type: AddCommentDto })
  @ApiResponse({ status: 201, description: 'Комментарий успешно добавлен' })
  async addComment(
    @Param('id') id: string,
    @Body() data: AddCommentDto,
    @Request() req,
  ) {
    const result = await firstValueFrom(
      this.tasksClient.send(RabbitMQPatterns.TASKS_ADD_COMMENT, {
        taskId: parseInt(id),
        userId: req.user.id,
        message: data.message,
      }),
    );

    if (!result.success) {
      handleRPCError(result);
    }

    return result.data;
  }

  @Get(':id/comments')
  @ApiOperation({ summary: 'Получить все комментарии задачи' })
  @ApiParam({ name: 'id', description: 'ID задачи', type: 'number' })
  @ApiResponse({ status: 200, description: 'Список комментариев' })
  async getComments(@Param('id') id: string) {
    const result = await firstValueFrom(
      this.tasksClient.send(RabbitMQPatterns.TASKS_GET_COMMENTS, {
        taskId: parseInt(id),
      }),
    );

    if (!result.success) {
      handleRPCError(result);
    }

    return result.data;
  }
}
