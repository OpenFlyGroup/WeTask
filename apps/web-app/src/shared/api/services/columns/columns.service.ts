import { UrlBuilder } from '@relatecom/utils'
import { Column, CreateColumnDto } from './columns.interface'
import { instance } from '../../instance'

export type UpdateColumnDto = Partial<CreateColumnDto> & { order?: number }

const PATH = '/columns'
const { buildUrl } = new UrlBuilder(PATH)

export const ColumnsService = {
  async getColumnsByBoard(boardId: number) {
    return (await instance.get<Column[]>(buildUrl(`/board/${boardId}`))).data
  },

  async createColumn(data: CreateColumnDto) {
    return (await instance.post<Column>(buildUrl(''), data)).data
  },

  async updateColumn(id: number, data: UpdateColumnDto) {
    return (await instance.put<Column>(buildUrl(`/${id}`), data)).data
  },

  async deleteColumn(id: number) {
    return instance.delete<{ success: true }>(buildUrl(`/${id}`))
  },
}
