export interface Board {
  id: number
  name: string
  description?: string | null
  teamId?: number | null
}

export interface CreateBoardDto {
  title: string
  teamId: number | undefined
}

export type UpdateBoardDto = Partial<CreateBoardDto>
