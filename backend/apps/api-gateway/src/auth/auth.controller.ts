import {
  Controller,
  Post,
  Body,
  HttpCode,
  HttpStatus,
  HttpException,
  ConflictException,
  UnauthorizedException,
} from '@nestjs/common';
import { Inject } from '@nestjs/common';
import { ClientProxy } from '@nestjs/microservices';
import { ApiTags, ApiOperation, ApiResponse, ApiBody } from '@nestjs/swagger';
import { RabbitMQPatterns } from '@libs/common';
import { firstValueFrom } from 'rxjs';
import {
  RegisterDto,
  LoginDto,
  RefreshTokenDto,
  AuthResponseDto,
} from '../common/dto/auth.dto';

@ApiTags('auth')
@Controller('auth')
export class AuthController {
  constructor(@Inject('AUTH_SERVICE') private authClient: ClientProxy) {}

  @Post('register')
  @HttpCode(HttpStatus.CREATED)
  @ApiOperation({ summary: 'Регистрация нового пользователя' })
  @ApiBody({ type: RegisterDto })
  @ApiResponse({
    status: 201,
    description: 'Пользователь успешно зарегистрирован',
    type: AuthResponseDto,
  })
  @ApiResponse({
    status: 409,
    description: 'Пользователь с таким email уже существует',
  })
  async register(@Body() data: RegisterDto) {
    // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
    const result = await firstValueFrom(
      this.authClient.send(RabbitMQPatterns.AUTH_REGISTER, data),
    );

    // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
    if (!result.success) {
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      const statusCode = result.statusCode || HttpStatus.CONFLICT;
      // eslint-disable-next-line @typescript-eslint/no-unsafe-argument, @typescript-eslint/no-unsafe-member-access
      throw new HttpException(result.error, statusCode);
    }

    // eslint-disable-next-line @typescript-eslint/no-unsafe-return, @typescript-eslint/no-unsafe-member-access
    return result.data;
  }

  @Post('login')
  @HttpCode(HttpStatus.OK)
  @ApiOperation({ summary: 'Вход в систему' })
  @ApiBody({ type: LoginDto })
  @ApiResponse({
    status: 200,
    description: 'Успешный вход',
    type: AuthResponseDto,
  })
  @ApiResponse({ status: 401, description: 'Неверные учетные данные' })
  async login(@Body() data: LoginDto) {
    // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
    const result = await firstValueFrom(
      this.authClient.send(RabbitMQPatterns.AUTH_LOGIN, data),
    );

    // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
    if (!result.success) {
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      const statusCode = result.statusCode || HttpStatus.UNAUTHORIZED;
      // eslint-disable-next-line @typescript-eslint/no-unsafe-argument, @typescript-eslint/no-unsafe-member-access
      throw new HttpException(result.error, statusCode);
    }

    // eslint-disable-next-line @typescript-eslint/no-unsafe-return, @typescript-eslint/no-unsafe-member-access
    return result.data;
  }

  @Post('refresh')
  @HttpCode(HttpStatus.OK)
  @ApiOperation({ summary: 'Обновление access токена' })
  @ApiBody({ type: RefreshTokenDto })
  @ApiResponse({
    status: 200,
    description: 'Токен успешно обновлен',
    schema: {
      type: 'object',
      properties: {
        accessToken: {
          type: 'string',
          example: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...',
        },
        refreshToken: { type: 'string', example: 'new-refresh-token' },
      },
    },
  })
  @ApiResponse({
    status: 401,
    description: 'Неверный или истекший refresh токен',
  })
  async refresh(@Body() data: RefreshTokenDto) {
    // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
    const result = await firstValueFrom(
      this.authClient.send(RabbitMQPatterns.AUTH_REFRESH, data),
    );

    // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
    if (!result.success) {
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      const statusCode = result.statusCode || HttpStatus.UNAUTHORIZED;
      // eslint-disable-next-line @typescript-eslint/no-unsafe-argument, @typescript-eslint/no-unsafe-member-access
      throw new HttpException(result.error, statusCode);
    }

    // eslint-disable-next-line @typescript-eslint/no-unsafe-return, @typescript-eslint/no-unsafe-member-access
    return result.data;
  }
}
