import { Injectable } from '@nestjs/common';
import { DatabaseService } from '@libs/database';

@Injectable()
export class ColumnsService {
  constructor(private database: DatabaseService) {}

  async create(data: { title: string; boardId: number; order: number }) {
    const column = await this.database.column.save({
      title: data.title,
      boardId: data.boardId,
      order: data.order,
    });

    return {
      success: true,
      data: column,
    };
  }

  async getByBoard(boardId: number) {
    const columns = await this.database.column.find({
      where: { boardId },
      relations: ['tasks', 'tasks.user'],
    });

    columns.sort((a, b) => a.order - b.order);

    return {
      success: true,
      data: columns,
    };
  }

  async update(id: number, data: { title?: string; order?: number }) {
    const column = await this.database.column.findOne({
      where: { id },
    });

    if (!column) {
      return {
        success: false,
        error: 'Column not found',
        statusCode: 404,
      };
    }

    if (data.title) column.title = data.title;
    if (data.order !== undefined) column.order = data.order;

    const updatedColumn = await this.database.column.save(column);

    return {
      success: true,
      data: updatedColumn,
    };
  }

  async delete(id: number) {
    const column = await this.database.column.findOne({
      where: { id },
    });

    if (!column) {
      return {
        success: false,
        error: 'Column not found',
        statusCode: 404,
      };
    }

    await this.database.column.remove(column);

    return {
      success: true,
      data: { message: 'Column deleted successfully' },
    };
  }
}
