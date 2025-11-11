import { ApiProperty } from '@nestjs/swagger';
import { IsString, IsNotEmpty, IsInt, IsOptional } from 'class-validator';

export class CreateBoardDto {
  @ApiProperty({
    description: 'Название доски',
    example: 'Project Board',
  })
  @IsString()
  @IsNotEmpty()
  title: string;

  @ApiProperty({
    description: 'ID команды, к которой принадлежит доска',
    example: 1,
  })
  @IsInt()
  teamId: number;
}

export class UpdateBoardDto {
  @ApiProperty({
    description: 'Название доски',
    example: 'Updated Project Board',
    required: false,
  })
  @IsOptional()
  @IsString()
  title?: string;
}
