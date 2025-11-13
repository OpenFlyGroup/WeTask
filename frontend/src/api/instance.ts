import AuthStorage from '@/store/auth'
import axios, { AxiosError, AxiosRequestConfig, AxiosResponse } from 'axios'
import toast from 'react-hot-toast'
import { AuthService } from './services/auth/auth.service'

export const instance = axios.create({
  baseURL: import.meta.env.VITE_BASE_URL,
})

instance.interceptors.request.use(
  (config) => {
    const tokens = AuthStorage.getTokens()
    if (tokens?.accessToken) {
      config.headers = config.headers ?? {}
      config.headers.Authorization = `Bearer ${tokens.accessToken}`
    }
    return config
  },
  (error) => Promise.reject(error),
)

interface FailedRequest {
  resolve: (value: AxiosResponse) => void
  reject: (value: AxiosError) => void
  config: AxiosRequestConfig
}

let failedRequests: FailedRequest[] = []
let isTokenRefreshing = false

instance.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const status = error.response?.status
    const originalRequestConfig = error.config as AxiosRequestConfig<any>

    if (status !== 401) {
      return Promise.reject(error)
    }

    if (isTokenRefreshing) {
      return new Promise<AxiosResponse>((resolve, reject) => {
        failedRequests.push({
          resolve,
          reject,
          config: originalRequestConfig,
        })
      })
    }

    isTokenRefreshing = true

    try {
      const tokens = AuthStorage.getTokens()

      if (!tokens?.refreshToken) {
        throw new Error('Refresh token отсутствует')
      }

      const response = await AuthService.getNewTokens(tokens.refreshToken)

      if (!response?.accessToken || !response?.refreshToken) {
        throw new Error('Некорректный ответ при обновлении токенов')
      }
      const { accessToken, refreshToken } = response

      AuthStorage.setTokens({ accessToken, refreshToken })

      const retryRequests = failedRequests.map(({ resolve, reject, config }) =>
        instance(config).then(resolve).catch(reject),
      )
      await Promise.all(retryRequests)
    } catch (refreshError) {
      console.error('Ошибка обновления токенов:', refreshError)

      failedRequests.forEach(({ reject }) => reject(refreshError as AxiosError))

      AuthStorage.clearTokens()
      toast.error('Сессия истекла. Пожалуйста, войдите заново.')

      return Promise.reject(refreshError)
    } finally {
      failedRequests = []
      isTokenRefreshing = false
    }

    return instance(originalRequestConfig)
  },
)
