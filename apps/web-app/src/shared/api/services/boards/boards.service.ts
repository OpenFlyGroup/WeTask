import { UrlBuilder } from '@relatecom/utils'
import { Board, CreateBoardDto, UpdateBoardDto } from './boards.interface'
import { instance } from '../../instance'

const PATH = '/boards'
const { buildUrl } = new UrlBuilder(PATH)

export const BoardsService = {
  async getBoards() {
    return (await instance.get<Board[]>(buildUrl('/'))).data
  },

  async getBoardById(id: number) {
    return (await instance.get<Board>(buildUrl(`/${id}`))).data
  },

  async createBoard(data: CreateBoardDto) {
    return (await instance.post<Board>(buildUrl('/'), data)).data
  },

  async updateBoard(id: number, data: UpdateBoardDto) {
    return (await instance.put<Board>(buildUrl(`/${id}`), data)).data
  },

  async deleteBoard(id: number) {
    return (await instance.delete<{ success: true }>(buildUrl(`/${id}`))).data
  },
}
