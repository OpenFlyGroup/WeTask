/* eslint-disable @typescript-eslint/no-unsafe-argument */
/* eslint-disable @typescript-eslint/no-unsafe-return */
/* eslint-disable @typescript-eslint/no-unsafe-member-access */
/* eslint-disable @typescript-eslint/no-unsafe-assignment */
import {
  Controller,
  Get,
  Post,
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
import { CreateTeamDto, AddMemberDto } from '../common/dto/teams.dto';
import { handleRPCError } from '../common/utils/error-handler';

@ApiTags('teams')
@ApiBearerAuth('JWT-auth')
@Controller('teams')
@UseGuards(AuthGuard('jwt'))
export class TeamsController {
  constructor(@Inject('TEAMS_SERVICE') private teamsClient: ClientProxy) {}

  @Get()
  @ApiOperation({ summary: 'Получить все команды текущего пользователя' })
  @ApiResponse({ status: 200, description: 'Список команд' })
  async getAll(@Request() req) {
    const result = await firstValueFrom(
      this.teamsClient.send(RabbitMQPatterns.TEAMS_GET_USER_TEAMS, {
        userId: req.user.id,
      }),
    );

    if (!result.success) {
      // eslint-disable-next-line @typescript-eslint/no-unsafe-argument
      handleRPCError(result);
    }

    return result.data;
  }

  @Post()
  @ApiOperation({ summary: 'Создать новую команду' })
  @ApiBody({ type: CreateTeamDto })
  @ApiResponse({ status: 201, description: 'Команда успешно создана' })
  async create(@Request() req, @Body() data: CreateTeamDto) {
    const result = await firstValueFrom(
      this.teamsClient.send(RabbitMQPatterns.TEAMS_CREATE, {
        name: data.name,
        ownerId: req.user.id,
      }),
    );

    if (!result.success) {
      handleRPCError(result);
    }

    return result.data;
  }

  @Get(':id')
  @ApiOperation({ summary: 'Получить команду по ID' })
  @ApiParam({ name: 'id', description: 'ID команды', type: 'number' })
  @ApiResponse({ status: 200, description: 'Информация о команде' })
  @ApiResponse({ status: 404, description: 'Команда не найдена' })
  async getById(@Param('id') id: string) {
    const result = await firstValueFrom(
      this.teamsClient.send(RabbitMQPatterns.TEAMS_GET_BY_ID, {
        id: parseInt(id),
      }),
    );

    if (!result.success) {
      handleRPCError(result);
    }

    return result.data;
  }

  @Post(':id/members')
  @ApiOperation({ summary: 'Добавить участника в команду' })
  @ApiParam({ name: 'id', description: 'ID команды', type: 'number' })
  @ApiBody({ type: AddMemberDto })
  @ApiResponse({ status: 201, description: 'Участник успешно добавлен' })
  @ApiResponse({ status: 404, description: 'Команда не найдена' })
  @ApiResponse({
    status: 409,
    description: 'Пользователь уже является участником команды',
  })
  async addMember(@Param('id') id: string, @Body() data: AddMemberDto) {
    const result = await firstValueFrom(
      this.teamsClient.send(RabbitMQPatterns.TEAMS_ADD_MEMBER, {
        teamId: parseInt(id),
        ...data,
      }),
    );

    if (!result.success) {
      handleRPCError(result);
    }

    return result.data;
  }

  @Delete(':id/members/:userId')
  @ApiOperation({ summary: 'Удалить участника из команды' })
  @ApiParam({ name: 'id', description: 'ID команды', type: 'number' })
  @ApiParam({ name: 'userId', description: 'ID пользователя', type: 'number' })
  @ApiResponse({ status: 200, description: 'Участник успешно удален' })
  @ApiResponse({ status: 404, description: 'Участник не найден' })
  async removeMember(@Param('id') id: string, @Param('userId') userId: string) {
    const result = await firstValueFrom(
      this.teamsClient.send(RabbitMQPatterns.TEAMS_REMOVE_MEMBER, {
        teamId: parseInt(id),
        userId: parseInt(userId),
      }),
    );

    if (!result.success) {
      handleRPCError(result);
    }

    return result.data;
  }
}
