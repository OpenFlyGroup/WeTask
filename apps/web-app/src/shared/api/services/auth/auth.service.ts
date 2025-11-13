import { UrlBuilder } from '@relatecom/utils'
import { IAuthResponse, ISignInDto, ISignUpDto } from './auth.interface'
import { instance } from '../../instance'
import AuthStorage from '@/shared/store/authStore'

const PATH = '/auth'
const { buildUrl } = new UrlBuilder(PATH)

export const AuthService = {
  async emailConfirmation({
    code,
    uid,
  }: {
    uid: string | undefined
    code: string | undefined
  }) {
    if (!uid || !code) {
      throw new Error('UID и code обязательны для подтверждения email')
    }

    const response = await instance.get(
      buildUrl(`/confirm-email/${uid}/${code}`),
    )
    return response.data
  },

  async getNewTokens(refreshToken: string): Promise<IAuthResponse> {
    const response = await instance.post<IAuthResponse>(buildUrl('/refresh'), {
      refreshToken,
    })

    return response.data
  },

  async logout(refreshToken: string | null) {
    if (!refreshToken) {
      console.warn('Refresh token отсутствует при логауте')
      return
    }

    try {
      await instance.post(buildUrl('/logout'), { refreshToken })
    } catch (error) {
      console.error('Ошибка при логауте', error)
    } finally {
      AuthStorage.clearTokens()
    }
  },

  async signIn(data: ISignInDto): Promise<IAuthResponse> {
    const response = await instance.post<IAuthResponse>(
      buildUrl('/login'),
      data,
    )
    const { accessToken, refreshToken } = response.data

    if (accessToken && refreshToken) {
      AuthStorage.setTokens({ accessToken, refreshToken })
    } else {
      console.error('No tokens in response:', response.data)
    }

    return response.data
  },

  async signUp(data: ISignUpDto): Promise<IAuthResponse> {
    const response = await instance.post<IAuthResponse>(
      buildUrl(`/register`),
      data,
    )
    const { accessToken, refreshToken } = response.data

    if (accessToken && refreshToken) {
      AuthStorage.setTokens({ accessToken, refreshToken })
    } else {
      console.error('No tokens in response:', response.data)
    }
    return response.data
  },

  async getCurrentUser() {
    const response = await instance.get(buildUrl('/me'))
    return response.data
  },
}
