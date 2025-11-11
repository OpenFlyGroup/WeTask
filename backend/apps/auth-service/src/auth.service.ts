import { Injectable } from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import * as bcrypt from 'bcryptjs';
import { DatabaseService } from '@libs/database';

@Injectable()
export class AuthService {
  constructor(
    private database: DatabaseService,
    private jwtService: JwtService,
  ) {}

  async register(data: { email: string; password: string; name: string }) {
    const existingUser = await this.database.user.findOne({
      where: { email: data.email },
    });

    if (existingUser) {
      return {
        success: false,
        error: 'User with this email already exists',
        statusCode: 409,
      };
    }

    const hashedPassword = await bcrypt.hash(data.password, 10);

    const user = await this.database.user.save({
      email: data.email,
      password: hashedPassword,
      name: data.name,
    });

    const tokens = await this.generateTokens(user.id);

    return {
      success: true,
      data: {
        user: {
          id: user.id,
          email: user.email,
          name: user.name,
          createdAt: user.createdAt,
        },
        ...tokens,
      },
    };
  }

  async login(data: { email: string; password: string }) {
    const user = await this.database.user.findOne({
      where: { email: data.email },
    });

    if (!user) {
      return {
        success: false,
        error: 'Invalid credentials',
        statusCode: 401,
      };
    }

    const isPasswordValid = await bcrypt.compare(data.password, user.password);

    if (!isPasswordValid) {
      return {
        success: false,
        error: 'Invalid credentials',
        statusCode: 401,
      };
    }

    const tokens = await this.generateTokens(user.id);

    return {
      success: true,
      data: {
        user: {
          id: user.id,
          email: user.email,
          name: user.name,
          createdAt: user.createdAt,
        },
        ...tokens,
      },
    };
  }

  async refresh(refreshToken: string) {
    const token = await this.database.refreshToken.findOne({
      where: { token: refreshToken },
      relations: ['user'],
    });

    if (!token || token.expiresAt < new Date()) {
      return {
        success: false,
        error: 'Invalid or expired refresh token',
        statusCode: 401,
      };
    }

    await this.database.refreshToken.remove(token);

    const tokens = await this.generateTokens(token.userId);

    return {
      success: true,
      data: tokens,
    };
  }

  async validateToken(tokenOrUserId: string | number) {
    try {
      // Если передан числовой ID пользователя (из JWT payload)
      const userId = typeof tokenOrUserId === 'number' 
        ? tokenOrUserId 
        : parseInt(tokenOrUserId, 10);
      
      if (!isNaN(userId)) {
        const user = await this.database.user.findOne({
          where: { id: userId },
        });

        if (!user) {
          return {
            success: false,
            error: 'User not found',
            statusCode: 401,
          };
        }

        return {
          success: true,
          data: {
            id: user.id,
            email: user.email,
            name: user.name,
          },
        };
      }

      // Если передан JWT токен (строка)
      const payload = this.jwtService.verify(tokenOrUserId as string);
      const user = await this.database.user.findOne({
        where: { id: payload.sub },
      });

      if (!user) {
        return {
          success: false,
          error: 'User not found',
          statusCode: 401,
        };
      }

      return {
        success: true,
        data: {
          id: user.id,
          email: user.email,
          name: user.name,
        },
      };
    } catch (error) {
      return {
        success: false,
        error: 'Invalid token',
        statusCode: 401,
      };
    }
  }

  private async generateTokens(userId: number) {
    const payload = { sub: userId };
    const accessToken = this.jwtService.sign(payload);
    const refreshToken = this.generateRefreshToken();

    const expiresAt = new Date();
    expiresAt.setDate(expiresAt.getDate() + 7); // 7 days

    await this.database.refreshToken.save({
      token: refreshToken,
      userId,
      expiresAt,
    });

    return {
      accessToken,
      refreshToken,
    };
  }

  private generateRefreshToken(): string {
    return require('crypto').randomBytes(64).toString('hex');
  }
}
