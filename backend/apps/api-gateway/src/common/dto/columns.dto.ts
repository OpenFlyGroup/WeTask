import { ApiProperty } from '@nestjs/swagger';
import { IsString, IsNotEmpty, IsInt, IsOptional } from 'class-validator';

export class CreateColumnDto {
  @ApiProperty({
    description: 'Название колонки',
    example: 'To Do',
  })
  @IsString()
  @IsNotEmpty()
  title: string;

  @ApiProperty({
    description: 'ID доски, к которой принадлежит колонка',
    example: 1,
  })
  @IsInt()
  boardId: number;

  @ApiProperty({
    description: 'Порядок колонки в доске',
    example: 0,
  })
  @IsInt()
  order: number;
}

export class UpdateColumnDto {
  @ApiProperty({
    description: 'Название колонки',
    example: 'In Progress',
    required: false,
  })
  @IsOptional()
  @IsString()
  title?: string;

  @ApiProperty({
    description: 'Порядок колонки в доске',
    example: 1,
    required: false,
  })
  @IsOptional()
  @IsInt()
  order?: number;
}
