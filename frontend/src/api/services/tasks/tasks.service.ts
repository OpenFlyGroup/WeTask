import { instance } from '@/api/instance'
import { UrlBuilder } from '@relatecom/utils'
import {
  AddCommentDto,
  CreateTaskDto,
  MoveTaskDto,
  Task,
  UpdateTaskDto,
  Comment,
} from './tasks.interface'

const PATH = '/tasks'
const { buildUrl } = new UrlBuilder(PATH)

export const TasksService = {
  async getTaskById(id: number) {
    return (await instance.get<Task>(buildUrl(`/${id}`))).data
  },

  async getTasksByBoard(boardId: number) {
    return (await instance.get<Task[]>(buildUrl(`/board/${boardId}`))).data
  },

  async createTask(data: CreateTaskDto) {
    return (await instance.post<Task>(buildUrl('/'), data)).data
  },

  async updateTask(id: number, data: UpdateTaskDto) {
    return (await instance.put<Task>(buildUrl(`/${id}`), data)).data
  },

  async deleteTask(id: number) {
    return (await instance.delete<{ success: true }>(buildUrl(`/${id}`))).data
  },

  async moveTask(id: number, data: MoveTaskDto) {
    return (await instance.put<Task>(buildUrl(`/${id}/move`), data)).data
  },

  async addComment(taskId: number, data: AddCommentDto) {
    return (await instance.post<Comment>(buildUrl(`/${taskId}/comment`), data))
      .data
  },

  async getComments(taskId: number) {
    return (await instance.get<Comment[]>(buildUrl(`/${taskId}/comments`))).data
  },
}
