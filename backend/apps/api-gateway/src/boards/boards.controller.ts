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
import { CreateBoardDto, UpdateBoardDto } from '../common/dto/boards.dto';
import { handleRPCError } from '../common/utils/error-handler';

@ApiTags('boards')
@ApiBearerAuth('JWT-auth')
@Controller('boards')
@UseGuards(AuthGuard('jwt'))
export class BoardsController {
  constructor(@Inject('BOARDS_SERVICE') private boardsClient: ClientProxy) {}

  @Get()
  @ApiOperation({ summary: 'Получить все доски текущего пользователя' })
  @ApiResponse({ status: 200, description: 'Список досок' })
  async getAll(@Request() req) {
    const result = await firstValueFrom(
      this.boardsClient.send(RabbitMQPatterns.BOARDS_GET_ALL, {
        userId: req.user.id,
      }),
    );

    if (!result.success) {
      handleRPCError(result);
    }

    return result.data;
  }

  @Post()
  @ApiOperation({ summary: 'Создать новую доску' })
  @ApiBody({ type: CreateBoardDto })
  @ApiResponse({ status: 201, description: 'Доска успешно создана' })
  async create(@Body() data: CreateBoardDto) {
    const result = await firstValueFrom(
      this.boardsClient.send(RabbitMQPatterns.BOARDS_CREATE, data),
    );

    if (!result.success) {
      handleRPCError(result);
    }

    return result.data;
  }

  @Get(':id')
  @ApiOperation({ summary: 'Получить доску по ID с колонками и задачами' })
  @ApiParam({ name: 'id', description: 'ID доски', type: 'number' })
  @ApiResponse({ status: 200, description: 'Информация о доске' })
  @ApiResponse({ status: 404, description: 'Доска не найдена' })
  async getById(@Param('id') id: string) {
    const result = await firstValueFrom(
      this.boardsClient.send(RabbitMQPatterns.BOARDS_GET_BY_ID, {
        id: parseInt(id),
      }),
    );

    if (!result.success) {
      handleRPCError(result);
    }

    return result.data;
  }

  @Put(':id')
  @ApiOperation({ summary: 'Обновить доску' })
  @ApiParam({ name: 'id', description: 'ID доски', type: 'number' })
  @ApiBody({ type: UpdateBoardDto })
  @ApiResponse({ status: 200, description: 'Доска успешно обновлена' })
  @ApiResponse({ status: 404, description: 'Доска не найдена' })
  async update(@Param('id') id: string, @Body() data: UpdateBoardDto) {
    const result = await firstValueFrom(
      this.boardsClient.send(RabbitMQPatterns.BOARDS_UPDATE, {
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
  @ApiOperation({ summary: 'Удалить доску' })
  @ApiParam({ name: 'id', description: 'ID доски', type: 'number' })
  @ApiResponse({ status: 200, description: 'Доска успешно удалена' })
  @ApiResponse({ status: 404, description: 'Доска не найдена' })
  async delete(@Param('id') id: string) {
    const result = await firstValueFrom(
      this.boardsClient.send(RabbitMQPatterns.BOARDS_DELETE, {
        id: parseInt(id),
      }),
    );

    if (!result.success) {
      handleRPCError(result);
    }

    return result.data;
  }
}
