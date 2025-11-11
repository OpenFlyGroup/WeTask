import { http } from './http'

export type Board = {
  id: number
  name: string
  description?: string | null
  teamId?: number | null
}

export type CreateBoardDto = {
  name: string
  description?: string
  teamId?: number | null
}

export type UpdateBoardDto = Partial<CreateBoardDto>

export function getBoards(signal?: AbortSignal) {
  return http<Board[]>('/boards', { auth: true, signal })
}

export function getBoardById(id: number, signal?: AbortSignal) {
  return http<Board>(`/boards/${id}`, { auth: true, signal })
}

export function createBoard(data: CreateBoardDto) {
  return http<Board>('/boards', { method: 'POST', body: data, auth: true })
}

export function updateBoard(id: number, data: UpdateBoardDto) {
  return http<Board>(`/boards/${id}`, { method: 'PUT', body: data, auth: true })
}

export function deleteBoard(id: number) {
  return http<{ success: true }>(`/boards/${id}`, { method: 'DELETE', auth: true })
}


