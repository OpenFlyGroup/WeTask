import { isBrowser } from '@/utils/browser'

interface Tokens {
  accessToken: string
  refreshToken: string
}

class AuthStorage {
  private static ACCESS_TOKEN_KEY = 'access_token'
  private static REFRESH_TOKEN_KEY = 'refresh_token'

  private static canUseLocalStorage(): boolean {
    if (!isBrowser) return false
    try {
      const test = '__wetask_ls_test__'
      localStorage.setItem(test, test)
      localStorage.removeItem(test)
      return true
    } catch {
      return false
    }
  }

  static setTokens(tokens: Tokens): void {
    if (!tokens?.accessToken || !tokens?.refreshToken) {
      console.warn('Invalid tokens:', tokens)
      return
    }

    try {
      localStorage.setItem(this.ACCESS_TOKEN_KEY, tokens.accessToken)
      localStorage.setItem(this.REFRESH_TOKEN_KEY, tokens.refreshToken)
    } catch (error) {
      console.error('Failed to save tokens', error)
    }
  }

  static getTokens(): Tokens | null {
    if (!this.canUseLocalStorage()) return null
    try {
      const accessToken = localStorage.getItem(this.ACCESS_TOKEN_KEY)
      const refreshToken = localStorage.getItem(this.REFRESH_TOKEN_KEY)
      if (!accessToken || !refreshToken) return null
      return { accessToken, refreshToken }
    } catch (error) {
      console.warn('Failed to read tokens from localStorage', error)
      return null
    }
  }

  static clearTokens(): void {
    if (!this.canUseLocalStorage()) return
    try {
      localStorage.removeItem(this.ACCESS_TOKEN_KEY)
      localStorage.removeItem(this.REFRESH_TOKEN_KEY)
    } catch (error) {
      console.warn('Failed to clear tokens', error)
    }
  }

  static isAuthenticated(): boolean {
    return this.getTokens() !== null
  }
}

export default AuthStorage
