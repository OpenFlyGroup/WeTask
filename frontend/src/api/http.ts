export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'

export const API_BASE_URL =
  (import.meta as any).env?.VITE_API_URL || (globalThis as any).VITE_API_URL || 'http://localhost:3000'
const API_URL = API_BASE_URL

type RequestOptions<TBody = unknown> = {
  method?: HttpMethod
  body?: TBody
  headers?: Record<string, string>
  signal?: AbortSignal
  auth?: boolean
}

export type AuthTokens = {
  accessToken: string
  refreshToken: string
}

const ACCESS_TOKEN_KEY = 'wetask.accessToken'
const REFRESH_TOKEN_KEY = 'wetask.refreshToken'

let memoryTokens: AuthTokens | null = null

function safeLocalStorageGet(key: string): string | null {
  try {
    if (typeof window === 'undefined') return null
    return window.localStorage?.getItem(key) ?? null
  } catch {
    return null
  }
}
function safeLocalStorageSet(key: string, value: string): void {
  try {
    if (typeof window === 'undefined') return
    window.localStorage?.setItem(key, value)
  } catch {
    // ignore
  }
}
function safeLocalStorageRemove(key: string): void {
  try {
    if (typeof window === 'undefined') return
    window.localStorage?.removeItem(key)
  } catch {
    // ignore
  }
}

export const authStorage = {
  getTokens(): AuthTokens | null {
    // Prefer persisted tokens, but tolerate environments where localStorage is denied
    const accessToken = safeLocalStorageGet(ACCESS_TOKEN_KEY)
    const refreshToken = safeLocalStorageGet(REFRESH_TOKEN_KEY)
    if (accessToken && refreshToken) {
      memoryTokens = { accessToken, refreshToken }
      return memoryTokens
    }
    // fallback to memory
    return memoryTokens
  },
  setTokens(tokens: AuthTokens) {
    memoryTokens = tokens
    safeLocalStorageSet(ACCESS_TOKEN_KEY, tokens.accessToken)
    safeLocalStorageSet(REFRESH_TOKEN_KEY, tokens.refreshToken)
  },
  clear() {
    memoryTokens = null
    safeLocalStorageRemove(ACCESS_TOKEN_KEY)
    safeLocalStorageRemove(REFRESH_TOKEN_KEY)
  },
}

async function refreshTokenIfNeeded(response: Response): Promise<boolean> {
  if (response.status !== 401) return false
  const tokens = authStorage.getTokens()
  if (!tokens) return false
  try {
    const res = await fetch(`${API_URL}/auth/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refreshToken: tokens.refreshToken }),
    })
    if (!res.ok) return false
    const data = (await res.json()) as AuthTokens
    if (data?.accessToken && data?.refreshToken) {
      authStorage.setTokens({ accessToken: data.accessToken, refreshToken: data.refreshToken })
      return true
    }
    return false
  } catch {
    return false
  }
}

export async function http<TResponse, TBody = unknown>(
  path: string,
  opts: RequestOptions<TBody> = {},
): Promise<TResponse> {
  const url = path.startsWith('http') ? path : `${API_URL}${path}`
  const method = opts.method ?? 'GET'
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(opts.headers ?? {}),
  }
  if (opts.auth) {
    const tokens = authStorage.getTokens()
    if (tokens?.accessToken) {
      headers.Authorization = `Bearer ${tokens.accessToken}`
    }
  }

  const doFetch = async () =>
    fetch(url, {
      method,
      headers,
      body: opts.body ? JSON.stringify(opts.body) : undefined,
      signal: opts.signal,
    })

  let res = await doFetch()
  if (opts.auth && res.status === 401) {
    const refreshed = await refreshTokenIfNeeded(res)
    if (refreshed) {
      // retry once with new token
      const tokens = authStorage.getTokens()
      if (tokens?.accessToken) {
        headers.Authorization = `Bearer ${tokens.accessToken}`
      }
      res = await doFetch()
    } else {
      authStorage.clear()
      // Proactively route to login on client when unauthorized
      try {
        if (typeof window !== 'undefined' && window.location?.pathname !== '/auth/login') {
          window.location.href = '/auth/login'
        }
      } catch {
        // ignore navigation errors in non-browser environments
      }
    }
  }

  if (!res.ok) {
    let message = 'Request failed'
    try {
      const data = (await res.json()) as any
      message = data?.message || data?.error || message
    } catch {
      // ignore
    }
    throw new Error(message)
  }

  if (res.status === 204) {
    return undefined as unknown as TResponse
  }

  return (await res.json()) as TResponse
}


