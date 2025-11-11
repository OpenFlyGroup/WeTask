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
import { CreateColumnDto, UpdateColumnDto } from '../common/dto/columns.dto';
import { handleRPCError } from '../common/utils/error-handler';

@ApiTags('columns')
@ApiBearerAuth('JWT-auth')
@Controller('columns')
@UseGuards(AuthGuard('jwt'))
export class ColumnsController {
  constructor(@Inject('BOARDS_SERVICE') private boardsClient: ClientProxy) {}

  @Post()
  @ApiOperation({ summary: 'Создать новую колонку' })
  @ApiBody({ type: CreateColumnDto })
  @ApiResponse({ status: 201, description: 'Колонка успешно создана' })
  async create(@Body() data: CreateColumnDto) {
    const result = await firstValueFrom(
      this.boardsClient.send(RabbitMQPatterns.COLUMNS_CREATE, data),
    );

    if (!result.success) {
      handleRPCError(result);
    }

    return result.data;
  }

  @Get('board/:boardId')
  @ApiOperation({ summary: 'Получить все колонки доски' })
  @ApiParam({ name: 'boardId', description: 'ID доски', type: 'number' })
  @ApiResponse({ status: 200, description: 'Список колонок с задачами' })
  async getByBoard(@Param('boardId') boardId: string) {
    const result = await firstValueFrom(
      this.boardsClient.send(RabbitMQPatterns.COLUMNS_GET_BY_BOARD, {
        boardId: parseInt(boardId),
      }),
    );

    if (!result.success) {
      handleRPCError(result);
    }

    return result.data;
  }

  @Put(':id')
  @ApiOperation({ summary: 'Обновить колонку' })
  @ApiParam({ name: 'id', description: 'ID колонки', type: 'number' })
  @ApiBody({ type: UpdateColumnDto })
  @ApiResponse({ status: 200, description: 'Колонка успешно обновлена' })
  @ApiResponse({ status: 404, description: 'Колонка не найдена' })
  async update(@Param('id') id: string, @Body() data: UpdateColumnDto) {
    const result = await firstValueFrom(
      this.boardsClient.send(RabbitMQPatterns.COLUMNS_UPDATE, {
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
  @ApiOperation({ summary: 'Удалить колонку' })
  @ApiParam({ name: 'id', description: 'ID колонки', type: 'number' })
  @ApiResponse({ status: 200, description: 'Колонка успешно удалена' })
  @ApiResponse({ status: 404, description: 'Колонка не найдена' })
  async delete(@Param('id') id: string) {
    const result = await firstValueFrom(
      this.boardsClient.send(RabbitMQPatterns.COLUMNS_DELETE, {
        id: parseInt(id),
      }),
    );

    if (!result.success) {
      handleRPCError(result);
    }

    return result.data;
  }
}
