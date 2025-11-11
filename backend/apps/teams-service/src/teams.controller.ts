import { Controller } from '@nestjs/common';
import { MessagePattern, Payload, EventPattern } from '@nestjs/microservices';
import { RabbitMQPatterns } from '@libs/common';
import { TeamsService } from './teams.service';

@Controller()
export class TeamsController {
  constructor(private readonly teamsService: TeamsService) {}

  @MessagePattern(RabbitMQPatterns.TEAMS_CREATE)
  async create(@Payload() data: { name: string; ownerId: number }) {
    return this.teamsService.create(data);
  }

  @MessagePattern(RabbitMQPatterns.TEAMS_GET_ALL)
  async getAll(@Payload() data: { userId: number }) {
    return this.teamsService.getUserTeams(data.userId);
  }

  @MessagePattern(RabbitMQPatterns.TEAMS_GET_BY_ID)
  async getById(@Payload() data: { id: number }) {
    return this.teamsService.getById(data.id);
  }

  @MessagePattern(RabbitMQPatterns.TEAMS_GET_USER_TEAMS)
  async getUserTeams(@Payload() data: { userId: number }) {
    return this.teamsService.getUserTeams(data.userId);
  }

  @MessagePattern(RabbitMQPatterns.TEAMS_ADD_MEMBER)
  async addMember(
    @Payload() data: { teamId: number; userId: number; role?: string },
  ) {
    return this.teamsService.addMember(data);
  }

  @MessagePattern(RabbitMQPatterns.TEAMS_REMOVE_MEMBER)
  async removeMember(@Payload() data: { teamId: number; userId: number }) {
    return this.teamsService.removeMember(data);
  }
}
