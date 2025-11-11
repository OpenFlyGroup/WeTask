import { ApiProperty } from '@nestjs/swagger';
import {
  IsString,
  IsNotEmpty,
  IsInt,
  IsOptional,
  IsEnum,
} from 'class-validator';

export class CreateTaskDto {
  @ApiProperty({
    description: 'Название задачи',
    example: 'Implement user authentication',
  })
  @IsString()
  @IsNotEmpty()
  title: string;

  @ApiProperty({
    description: 'Описание задачи',
    example: 'Implement JWT-based authentication system',
    required: false,
  })
  @IsOptional()
  @IsString()
  description?: string;

  @ApiProperty({
    description: 'ID колонки, в которой находится задача',
    example: 1,
  })
  @IsInt()
  columnId: number;

  @ApiProperty({
    description: 'ID пользователя, которому назначена задача',
    example: 1,
    required: false,
  })
  @IsOptional()
  @IsInt()
  assignedTo?: number;

  @ApiProperty({
    description: 'Приоритет задачи',
    example: 'high',
    enum: ['low', 'medium', 'high'],
    required: false,
  })
  @IsOptional()
  @IsEnum(['low', 'medium', 'high'])
  priority?: string;
}

export class UpdateTaskDto {
  @ApiProperty({
    description: 'Название задачи',
    example: 'Updated task title',
    required: false,
  })
  @IsOptional()
  @IsString()
  title?: string;

  @ApiProperty({
    description: 'Описание задачи',
    example: 'Updated task description',
    required: false,
  })
  @IsOptional()
  @IsString()
  description?: string;

  @ApiProperty({
    description: 'Приоритет задачи',
    example: 'high',
    enum: ['low', 'medium', 'high'],
    required: false,
  })
  @IsOptional()
  @IsEnum(['low', 'medium', 'high'])
  priority?: string;

  @ApiProperty({
    description: 'ID пользователя, которому назначена задача',
    example: 2,
    required: false,
  })
  @IsOptional()
  @IsInt()
  assignedTo?: number;
}

export class MoveTaskDto {
  @ApiProperty({
    description: 'ID колонки, в которую перемещается задача',
    example: 2,
  })
  @IsInt()
  columnId: number;
}

export class AddCommentDto {
  @ApiProperty({
    description: 'Текст комментария',
    example: 'This task needs more work',
  })
  @IsString()
  @IsNotEmpty()
  message: string;
}

export class UpdateUserDto {
  @ApiProperty({
    description: 'Имя пользователя',
    example: 'John Doe',
    required: false,
  })
  @IsOptional()
  @IsString()
  name?: string;

  @ApiProperty({
    description: 'Email пользователя',
    example: 'newemail@example.com',
    required: false,
  })
  @IsOptional()
  @IsString()
  email?: string;
}
