import { Injectable, UnauthorizedException } from '@nestjs/common';
import { PassportStrategy } from '@nestjs/passport';
import { ExtractJwt, Strategy } from 'passport-jwt';
import { Inject } from '@nestjs/common';
import { ClientProxy } from '@nestjs/microservices';
import { RabbitMQPatterns } from '@libs/common';
import { firstValueFrom } from 'rxjs';

@Injectable()
export class JwtStrategy extends PassportStrategy(Strategy) {
  constructor(@Inject('AUTH_SERVICE') private authClient: ClientProxy) {
    super({
      jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
      ignoreExpiration: false,
      secretOrKey: process.env.JWT_SECRET || 'your-secret-key',
    });
  }

  async validate(payload: any) {
    // JWT уже валидирован на уровне стратегии, просто возвращаем payload
    // Можно дополнительно проверить пользователя через сервис
    try {
      // Проверяем пользователя по ID из payload (JWT уже валидирован)
      const result = await firstValueFrom(
        this.authClient.send(RabbitMQPatterns.AUTH_VALIDATE, {
          token: payload.sub || payload.id,
        }),
      );

      if (!result.success) {
        throw new UnauthorizedException();
      }

      return result.data;
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
    } catch (error) {
      // Если сервис недоступен, используем payload напрямую
      return {
        id: payload.sub,
        email: payload.email,
        name: payload.name,
      };
    }
  }
}
