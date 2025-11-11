import { http } from './http'

export type Column = {
  id: number
  name: string
  order: number
  boardId: number
}

export type CreateColumnDto = {
  name: string
  boardId: number
}

export type UpdateColumnDto = Partial<CreateColumnDto> & { order?: number }

export function getColumnsByBoard(boardId: number, signal?: AbortSignal) {
  return http<Column[]>(`/columns/board/${boardId}`, { auth: true, signal })
}

export function createColumn(data: CreateColumnDto) {
  return http<Column>('/columns', { method: 'POST', body: data, auth: true })
}

export function updateColumn(id: number, data: UpdateColumnDto) {
  return http<Column>(`/columns/${id}`, { method: 'PUT', body: data, auth: true })
}

export function deleteColumn(id: number) {
  return http<{ success: true }>(`/columns/${id}`, { method: 'DELETE', auth: true })
}


