import { Controller } from '@nestjs/common';
import { MessagePattern, Payload } from '@nestjs/microservices';
import { RabbitMQPatterns } from '@libs/common';
import { UsersService } from './users.service';

@Controller()
export class UsersController {
  constructor(private readonly usersService: UsersService) {}

  @MessagePattern(RabbitMQPatterns.USERS_GET_BY_ID)
  async getById(@Payload() data: { id: number }) {
    return this.usersService.getById(data.id);
  }

  @MessagePattern(RabbitMQPatterns.USERS_GET_BY_EMAIL)
  async getByEmail(@Payload() data: { email: string }) {
    return this.usersService.getByEmail(data.email);
  }

  @MessagePattern(RabbitMQPatterns.USERS_GET_ME)
  async getMe(@Payload() data: { userId: number }) {
    return this.usersService.getById(data.userId);
  }

  @MessagePattern(RabbitMQPatterns.USERS_UPDATE)
  async update(@Payload() data: { id: number; name?: string; email?: string }) {
    return this.usersService.update(data.id, data);
  }
}
