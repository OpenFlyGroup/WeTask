import { Controller } from '@nestjs/common';
import { MessagePattern, Payload } from '@nestjs/microservices';
import { RabbitMQPatterns } from '@libs/common';
import { BoardsService } from './boards.service';

@Controller()
export class BoardsController {
  constructor(private readonly boardsService: BoardsService) {}

  @MessagePattern(RabbitMQPatterns.BOARDS_CREATE)
  async create(@Payload() data: { title: string; teamId: number }) {
    return this.boardsService.create(data);
  }

  @MessagePattern(RabbitMQPatterns.BOARDS_GET_ALL)
  async getAll(@Payload() data: { userId: number }) {
    return this.boardsService.getAllByUser(data.userId);
  }

  @MessagePattern(RabbitMQPatterns.BOARDS_GET_BY_ID)
  async getById(@Payload() data: { id: number }) {
    return this.boardsService.getById(data.id);
  }

  @MessagePattern(RabbitMQPatterns.BOARDS_GET_BY_TEAM)
  async getByTeam(@Payload() data: { teamId: number }) {
    return this.boardsService.getByTeam(data.teamId);
  }

  @MessagePattern(RabbitMQPatterns.BOARDS_UPDATE)
  async update(@Payload() data: { id: number; title?: string }) {
    return this.boardsService.update(data.id, data);
  }

  @MessagePattern(RabbitMQPatterns.BOARDS_DELETE)
  async delete(@Payload() data: { id: number }) {
    return this.boardsService.delete(data.id);
  }
}
