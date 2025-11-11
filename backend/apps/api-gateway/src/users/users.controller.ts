/* eslint-disable @typescript-eslint/no-unsafe-return */
/* eslint-disable @typescript-eslint/no-unsafe-argument */
/* eslint-disable @typescript-eslint/no-unsafe-member-access */
/* eslint-disable @typescript-eslint/no-unsafe-assignment */
import {
  Controller,
  Get,
  Patch,
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
import { UpdateUserDto } from '../common/dto/tasks.dto';
import { handleRPCError } from '../common/utils/error-handler';

@ApiTags('users')
@ApiBearerAuth('JWT-auth')
@Controller('users')
@UseGuards(AuthGuard('jwt'))
export class UsersController {
  constructor(@Inject('USERS_SERVICE') private usersClient: ClientProxy) {}

  @Get('me')
  @ApiOperation({ summary: 'Получить информацию о текущем пользователе' })
  @ApiResponse({ status: 200, description: 'Информация о пользователе' })
  @ApiResponse({ status: 401, description: 'Не авторизован' })
  async getMe(@Request() req) {
    const result = await firstValueFrom(
      this.usersClient.send(RabbitMQPatterns.USERS_GET_ME, {
        userId: req.user.id,
      }),
    );

    if (!result.success) {
      handleRPCError(result);
    }

    return result.data;
  }

  @Get(':id')
  @ApiOperation({ summary: 'Получить пользователя по ID' })
  @ApiParam({ name: 'id', description: 'ID пользователя', type: 'number' })
  @ApiResponse({ status: 200, description: 'Информация о пользователе' })
  @ApiResponse({ status: 404, description: 'Пользователь не найден' })
  async getById(@Param('id') id: string) {
    const result = await firstValueFrom(
      this.usersClient.send(RabbitMQPatterns.USERS_GET_BY_ID, {
        id: parseInt(id),
      }),
    );

    if (!result.success) {
      handleRPCError(result);
    }

    return result.data;
  }

  @Patch(':id')
  @ApiOperation({ summary: 'Обновить информацию о пользователе' })
  @ApiParam({ name: 'id', description: 'ID пользователя', type: 'number' })
  @ApiBody({ type: UpdateUserDto })
  @ApiResponse({ status: 200, description: 'Пользователь успешно обновлен' })
  @ApiResponse({ status: 404, description: 'Пользователь не найден' })
  async update(@Param('id') id: string, @Body() data: UpdateUserDto) {
    const result = await firstValueFrom(
      this.usersClient.send(RabbitMQPatterns.USERS_UPDATE, {
        id: parseInt(id),
        ...data,
      }),
    );

    if (!result.success) {
      handleRPCError(result);
    }

    return result.data;
  }
}
