import { http } from './http'

export type Task = {
  id: number
  title: string
  description?: string | null
  boardId: number
  columnId: number
}

export type CreateTaskDto = {
  title: string
  description?: string
  boardId: number
  columnId: number
}

export type UpdateTaskDto = Partial<CreateTaskDto>

export type MoveTaskDto = {
  columnId: number
}

export type Comment = {
  id: number
  taskId: number
  userId: number
  message: string
  createdAt: string
}

export type AddCommentDto = {
  message: string
}

export function getTaskById(id: number, signal?: AbortSignal) {
  return http<Task>(`/tasks/${id}`, { auth: true, signal })
}

export function getTasksByBoard(boardId: number, signal?: AbortSignal) {
  return http<Task[]>(`/tasks/board/${boardId}`, { auth: true, signal })
}

export function createTask(data: CreateTaskDto) {
  return http<Task>('/tasks', { method: 'POST', body: data, auth: true })
}

export function updateTask(id: number, data: UpdateTaskDto) {
  return http<Task>(`/tasks/${id}`, { method: 'PUT', body: data, auth: true })
}

export function deleteTask(id: number) {
  return http<{ success: true }>(`/tasks/${id}`, { method: 'DELETE', auth: true })
}

export function moveTask(id: number, data: MoveTaskDto) {
  return http<Task>(`/tasks/${id}/move`, { method: 'PUT', body: data, auth: true })
}

export function addComment(taskId: number, data: AddCommentDto) {
  return http<Comment>(`/tasks/${taskId}/comment`, { method: 'POST', body: data, auth: true })
}

export function getComments(taskId: number, signal?: AbortSignal) {
  return http<Comment[]>(`/tasks/${taskId}/comments`, { auth: true, signal })
}


