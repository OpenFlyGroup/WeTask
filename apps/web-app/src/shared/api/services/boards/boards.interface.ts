export interface Board {
  id: number
  name: string
  description?: string | null
  teamId?: number | null
}

export interface CreateBoardDto {
  name: string
  description?: string
  teamId?: number | null
}

export type UpdateBoardDto = Partial<CreateBoardDto>
