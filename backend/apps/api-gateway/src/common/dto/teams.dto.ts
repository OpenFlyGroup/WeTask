import { ApiProperty } from '@nestjs/swagger';
import { IsString, IsNotEmpty, IsOptional, IsEnum } from 'class-validator';

export class CreateTeamDto {
  @ApiProperty({
    description: 'Название команды',
    example: 'Development Team',
  })
  @IsString()
  @IsNotEmpty()
  name: string;
}

export class AddMemberDto {
  @ApiProperty({
    description: 'ID пользователя для добавления в команду',
    example: 2,
  })
  userId: number;

  @ApiProperty({
    description: 'Роль участника в команде',
    example: 'member',
    enum: ['owner', 'admin', 'member'],
    required: false,
  })
  @IsOptional()
  @IsEnum(['owner', 'admin', 'member'])
  role?: string;
}
