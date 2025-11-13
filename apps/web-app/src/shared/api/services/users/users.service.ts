import { UrlBuilder } from '@relatecom/utils'
import { UpdateUserDto, User } from './users.interface'
import { instance } from '../../instance'

const PATH = '/users'
const { buildUrl } = new UrlBuilder(PATH)

export const UsersService = {
  async getMe() {
    return (await instance.get<User>(buildUrl('/me'))).data
  },

  async getUserById(id: number) {
    return (await instance.get<User>(buildUrl(`/${id}`))).data
  },

  async updateUser(id: number, data: UpdateUserDto) {
    return (await instance.patch<User>(buildUrl(`/${id}`), data)).data
  },
}
