export interface Task {
  id: number
  title: string
  description?: string | null
  boardId: number
  columnId: number
}

export interface CreateTaskDto {
  title: string
  description?: string
  boardId: number
  columnId: number
}

export type UpdateTaskDto = Partial<CreateTaskDto>

export interface MoveTaskDto {
  columnId: number
}

export interface Comment {
  id: number
  taskId: number
  userId: number
  message: string
  createdAt: string
}

export interface AddCommentDto {
  message: string
}
