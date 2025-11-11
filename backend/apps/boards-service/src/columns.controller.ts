import { Controller } from '@nestjs/common';
import { MessagePattern, Payload } from '@nestjs/microservices';
import { RabbitMQPatterns } from '@libs/common';
import { ColumnsService } from './columns.service';

@Controller()
export class ColumnsController {
  constructor(private readonly columnsService: ColumnsService) {}

  @MessagePattern(RabbitMQPatterns.COLUMNS_CREATE)
  async create(
    @Payload() data: { title: string; boardId: number; order: number },
  ) {
    return this.columnsService.create(data);
  }

  @MessagePattern(RabbitMQPatterns.COLUMNS_GET_BY_BOARD)
  async getByBoard(@Payload() data: { boardId: number }) {
    return this.columnsService.getByBoard(data.boardId);
  }

  @MessagePattern(RabbitMQPatterns.COLUMNS_UPDATE)
  async update(
    @Payload() data: { id: number; title?: string; order?: number },
  ) {
    return this.columnsService.update(data.id, data);
  }

  @MessagePattern(RabbitMQPatterns.COLUMNS_DELETE)
  async delete(@Payload() data: { id: number }) {
    return this.columnsService.delete(data.id);
  }
}
