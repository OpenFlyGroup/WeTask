import {
  Injectable,
  NotFoundException,
  ConflictException,
} from '@nestjs/common';
import { DatabaseService } from '@libs/database';

@Injectable()
export class TeamsService {
  constructor(private database: DatabaseService) {}

  async create(data: { name: string; ownerId: number }) {
    const team = await this.database.team.save({
      name: data.name,
    });

    const member = await this.database.teamMember.save({
      teamId: team.id,
      userId: data.ownerId,
      role: 'owner',
    });

    const teamWithMembers = await this.database.team.findOne({
      where: { id: team.id },
      relations: ['members', 'members.user'],
    });

    return {
      success: true,
      data: teamWithMembers,
    };
  }

  async getById(id: number) {
    const team = await this.database.team.findOne({
      where: { id },
      relations: ['members', 'members.user'],
    });

    if (!team) {
      return {
        success: false,
        error: 'Team not found',
        statusCode: 404,
      };
    }

    return {
      success: true,
      data: team,
    };
  }

  async getUserTeams(userId: number) {
    const teamMembers = await this.database.teamMember.find({
      where: { userId },
      relations: ['team', 'team.members', 'team.members.user'],
    });

    const teams = teamMembers.map((tm) => tm.team);

    return {
      success: true,
      data: teams,
    };
  }

  async addMember(data: { teamId: number; userId: number; role?: string }) {
    const team = await this.database.team.findOne({
      where: { id: data.teamId },
    });

    if (!team) {
      return {
        success: false,
        error: 'Team not found',
        statusCode: 404,
      };
    }

    const existingMember = await this.database.teamMember.findOne({
      where: {
        teamId: data.teamId,
        userId: data.userId,
      },
    });

    if (existingMember) {
      return {
        success: false,
        error: 'User is already a member of this team',
        statusCode: 409,
      };
    }

    const member = await this.database.teamMember.save({
      teamId: data.teamId,
      userId: data.userId,
      role: data.role || 'member',
    });

    const memberWithUser = await this.database.teamMember.findOne({
      where: { id: member.id },
      relations: ['user'],
    });

    return {
      success: true,
      data: memberWithUser,
    };
  }

  async removeMember(data: { teamId: number; userId: number }) {
    const member = await this.database.teamMember.findOne({
      where: {
        teamId: data.teamId,
        userId: data.userId,
      },
    });

    if (!member) {
      return {
        success: false,
        error: 'Member not found',
        statusCode: 404,
      };
    }

    await this.database.teamMember.remove(member);

    return {
      success: true,
      data: { message: 'Member removed successfully' },
    };
  }
}
