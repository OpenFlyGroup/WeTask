const API_BASE_URL = 'ws://localhost:3000/ws'
let socket: any | null = null
let loadingPromise: Promise<any> | null = null

export async function getSocket(): Promise<any> {
  if (socket) return socket
  if (loadingPromise) return loadingPromise
  loadingPromise = (async () => {
    try {
      const { io } = await import('socket.io-client')
      socket = io(API_BASE_URL, {
        transports: ['websocket'],
        withCredentials: true,
        autoConnect: true,
      })
      return socket
    } catch (e) {
      // socket.io-client not installed or failed to load - degrade gracefully
      console.warn('Realtime disabled: socket.io-client not available', e)
      socket = null
      return null
    } finally {
      loadingPromise = null
    }
  })()
  return loadingPromise
}

export async function disconnectSocket() {
  if (socket?.disconnect) {
    try {
      socket.disconnect()
    } catch {
      // ignore
    }
  }
  socket = null
}
