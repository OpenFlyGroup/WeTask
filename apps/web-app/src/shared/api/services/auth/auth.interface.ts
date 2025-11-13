export interface ISignInDto {
  email: string
  password: string
}

export interface ISignUpDto {
  email: string
  password: string
  name?: string
}

export interface IAuthResponse {
  accessToken: string
  refreshToken: string
  user?: {
    id: string
    email: string
    name: string
    createdAt: string
    updatedAt: string
  }
}
