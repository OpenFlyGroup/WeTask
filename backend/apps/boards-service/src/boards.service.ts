import { Injectable } from '@nestjs/common';
import { DatabaseService } from '@libs/database';

@Injectable()
export class BoardsService {
  constructor(private database: DatabaseService) {}

  async create(data: { title: string; teamId: number }) {
    const board = await this.database.board.save({
      title: data.title,
      teamId: data.teamId,
    });

    const boardWithRelations = await this.database.board.findOne({
      where: { id: board.id },
      relations: ['team', 'columns'],
    });

    if (boardWithRelations?.columns) {
      boardWithRelations.columns.sort((a, b) => a.order - b.order);
    }

    return {
      success: true,
      data: boardWithRelations,
    };
  }

  async getAllByUser(userId: number) {
    const teamMembers = await this.database.teamMember.find({
      where: { userId },
      relations: ['team', 'team.boards', 'team.boards.columns'],
    });

    const teams = teamMembers.map((tm) => tm.team);
    const boards = teams.flatMap((team) => team.boards || []);

    // Sort columns
    boards.forEach((board) => {
      if (board.columns) {
        board.columns.sort((a, b) => a.order - b.order);
      }
    });

    return {
      success: true,
      data: boards,
    };
  }

  async getById(id: number) {
    const board = await this.database.board.findOne({
      where: { id },
      relations: [
        'team',
        'team.members',
        'team.members.user',
        'columns',
        'columns.tasks',
        'columns.tasks.user',
      ],
    });

    if (!board) {
      return {
        success: false,
        error: 'Board not found',
        statusCode: 404,
      };
    }

    if (board.columns) {
      board.columns.sort((a, b) => a.order - b.order);
      board.columns.forEach((col) => {
        if (col.tasks) {
          col.tasks.sort((a, b) => a.createdAt.getTime() - b.createdAt.getTime());
        }
      });
    }

    return {
      success: true,
      data: board,
    };
  }

  async getByTeam(teamId: number) {
    const boards = await this.database.board.find({
      where: { teamId },
      relations: ['columns'],
    });

    boards.forEach((board) => {
      if (board.columns) {
        board.columns.sort((a, b) => a.order - b.order);
      }
    });

    return {
      success: true,
      data: boards,
    };
  }

  async update(id: number, data: { title?: string }) {
    const board = await this.database.board.findOne({
      where: { id },
    });

    if (!board) {
      return {
        success: false,
        error: 'Board not found',
        statusCode: 404,
      };
    }

    if (data.title) board.title = data.title;

    const updatedBoard = await this.database.board.save(board);

    const boardWithColumns = await this.database.board.findOne({
      where: { id: updatedBoard.id },
      relations: ['columns'],
    });

    if (boardWithColumns?.columns) {
      boardWithColumns.columns.sort((a, b) => a.order - b.order);
    }

    return {
      success: true,
      data: boardWithColumns,
    };
  }

  async delete(id: number) {
    const board = await this.database.board.findOne({
      where: { id },
    });

    if (!board) {
      return {
        success: false,
        error: 'Board not found',
        statusCode: 404,
      };
    }

    await this.database.board.remove(board);

    return {
      success: true,
      data: { message: 'Board deleted successfully' },
    };
  }
}
