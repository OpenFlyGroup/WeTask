import { http } from './http'

export type User = {
  id: number
  email: string
  name?: string | null
}

export type UpdateUserDto = Partial<Pick<User, 'name'>>

export function getMe(signal?: AbortSignal) {
  return http<User>('/users/me', { auth: true, signal })
}

export function getUserById(id: number, signal?: AbortSignal) {
  return http<User>(`/users/${id}`, { auth: true, signal })
}

export function updateUser(id: number, data: UpdateUserDto) {
  return http<User>(`/users/${id}`, { method: 'PATCH', body: data, auth: true })
}


