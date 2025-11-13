export interface User {
  id: number
  email: string
  name?: string | null
}

export type UpdateUserDto = Partial<Pick<User, 'name'>>
