import { instance } from '@/api/instance'
import { UrlBuilder } from '@relatecom/utils'
import { User } from '@sentry/tanstackstart-react'

export type UpdateUserDto = Partial<Pick<User, 'name'>>

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
