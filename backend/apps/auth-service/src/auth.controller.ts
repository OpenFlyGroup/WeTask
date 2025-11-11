import { Controller } from '@nestjs/common';
import { MessagePattern, Payload } from '@nestjs/microservices';
import { RabbitMQPatterns } from '@libs/common';
import { AuthService } from './auth.service';

@Controller()
export class AuthController {
  constructor(private readonly authService: AuthService) {}

  @MessagePattern(RabbitMQPatterns.AUTH_REGISTER)
  async register(
    @Payload() data: { email: string; password: string; name: string },
  ) {
    return this.authService.register(data);
  }

  @MessagePattern(RabbitMQPatterns.AUTH_LOGIN)
  async login(@Payload() data: { email: string; password: string }) {
    return this.authService.login(data);
  }

  @MessagePattern(RabbitMQPatterns.AUTH_REFRESH)
  async refresh(@Payload() data: { refreshToken: string }) {
    return this.authService.refresh(data.refreshToken);
  }

  @MessagePattern(RabbitMQPatterns.AUTH_VALIDATE)
  async validate(@Payload() data: { token: string }) {
    return this.authService.validateToken(data.token);
  }
}
