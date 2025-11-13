export interface Column {
  id: number
  name: string
  order: number
  boardId: number
}

export interface CreateColumnDto {
  name: string
  boardId: number
}

export type UpdateColumnDto = Partial<CreateColumnDto> & { order?: number }
