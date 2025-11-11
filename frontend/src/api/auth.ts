import { http, authStorage, type AuthTokens } from './http'

export type LoginDto = { email: string; password: string }
export type RegisterDto = { email: string; password: string; name: string }

export async function login(data: LoginDto) {
  const res = await http<AuthTokens>('/auth/login', { method: 'POST', body: data })
  authStorage.setTokens(res)
  return res
}

export async function register(data: RegisterDto) {
  const res = await http<AuthTokens>('/auth/register', { method: 'POST', body: data })
  authStorage.setTokens(res)
  return res
}

export function logout() {
  authStorage.clear()
}


