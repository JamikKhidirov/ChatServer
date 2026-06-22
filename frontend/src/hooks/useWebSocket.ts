import { useEffect, useRef, useCallback } from 'react'
import { api } from '../api/client'

type MessageHandler = (msg: { type: string; payload: any }) => void

interface UseWSReturn {
  send: (type: string, payload: any) => void
}

export default function useWebSocket(onMessage: MessageHandler, enabled: boolean): UseWSReturn {
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectRef = useRef<number>(0)

  const connect = useCallback(() => {
    const token = api.getToken()
    if (!token || !enabled) return

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.host
    const url = `${protocol}//${host}/ws?token=${token}`

    const ws = new WebSocket(url)
    wsRef.current = ws

    ws.onopen = () => { if (reconnectRef.current) { clearTimeout(reconnectRef.current); reconnectRef.current = 0 } }

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        onMessage(data)
      } catch { /* ignore */ }
    }

    ws.onclose = () => {
      wsRef.current = null
      reconnectRef.current = window.setTimeout(connect, 3000)
    }

    ws.onerror = () => { ws.close() }
  }, [onMessage, enabled])

  useEffect(() => {
    if (enabled) connect()
    return () => {
      if (reconnectRef.current) { clearTimeout(reconnectRef.current); reconnectRef.current = 0 }
      if (wsRef.current) { wsRef.current.close(); wsRef.current = null }
    }
  }, [enabled, connect])

  const send = useCallback((type: string, payload: any) => {
    const ws = wsRef.current
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ type, payload }))
    }
  }, [])

  return { send }
}
