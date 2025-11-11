import { Injectable, NotFoundException } from '@nestjs/common';
import { DatabaseService } from '@libs/database';

@Injectable()
export class UsersService {
  constructor(private database: DatabaseService) {}

  async getById(id: number) {
    const user = await this.database.user.findOne({
      where: { id },
    });

    if (!user) {
      return {
        success: false,
        error: 'User not found',
        statusCode: 404,
      };
    }

    return {
      success: true,
      data: {
        id: user.id,
        email: user.email,
        name: user.name,
        createdAt: user.createdAt,
        updatedAt: user.updatedAt,
      },
    };
  }

  async getByEmail(email: string) {
    const user = await this.database.user.findOne({
      where: { email },
    });

    if (!user) {
      return {
        success: false,
        error: 'User not found',
        statusCode: 404,
      };
    }

    return {
      success: true,
      data: {
        id: user.id,
        email: user.email,
        name: user.name,
        createdAt: user.createdAt,
        updatedAt: user.updatedAt,
      },
    };
  }

  async update(id: number, data: { name?: string; email?: string }) {
    const user = await this.database.user.findOne({
      where: { id },
    });

    if (!user) {
      return {
        success: false,
        error: 'User not found',
        statusCode: 404,
      };
    }

    if (data.name) user.name = data.name;
    if (data.email) user.email = data.email;

    const updatedUser = await this.database.user.save(user);

    return {
      success: true,
      data: {
        id: updatedUser.id,
        email: updatedUser.email,
        name: updatedUser.name,
        createdAt: updatedUser.createdAt,
        updatedAt: updatedUser.updatedAt,
      },
    };
  }
}
