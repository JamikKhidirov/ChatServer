import { useState, useEffect, useRef, useCallback } from 'react'
import { api } from '../api/client'
import MessageBubble from './MessageBubble'

interface Props {
  chatId: string | null
  chatInfo: any
  token: string
  wsMessage?: any
  onWsConsumed?: () => void
  wsSend?: (type: string, payload: any) => void
}

export default function MessageArea({ chatId, chatInfo, wsMessage, onWsConsumed, wsSend }: Props) {
  const [messages, setMessages] = useState<any[]>([])
  const [text, setText] = useState('')
  const [sending, setSending] = useState(false)
  const bottomRef = useRef<HTMLDivElement>(null)
  const [chatName, setChatName] = useState('')
  const [uploading, setUploading] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)
  const typingTimerRef = useRef<number>(0)
  const isTypingRef = useRef(false)

  useEffect(() => {
    if (chatInfo?.name || chatInfo?.Name) {
      setChatName(chatInfo.name || chatInfo.Name)
    } else if (chatInfo?.participants) {
      const others = (chatInfo.participants || []).filter((p: any) => p.id !== api.getToken())
      setChatName(others.map((p: any) => p.displayName || p.username).join(', ') || 'Chat')
    }
  }, [chatInfo])

  useEffect(() => {
    if (wsMessage && chatId === wsMessage.chatId) {
      setMessages(prev => {
        if (prev.some(m => (m.id || m.ID) === (wsMessage.id || wsMessage.ID))) return prev
        return [...prev, wsMessage]
      })
      onWsConsumed?.()
    }
  }, [wsMessage])

  const loadMessages = useCallback(async () => {
    if (!chatId) return
    const res = await api.call('GET', `/api/chats/${chatId}/messages`)
    if (!res.error) {
      const d = res.data as any
      const msgs = d?.data || d || []
      setMessages(Array.isArray(msgs) ? msgs : [])
    }
  }, [chatId])

  useEffect(() => { loadMessages(); const iv = setInterval(loadMessages, 3000); return () => clearInterval(iv) }, [loadMessages])

  useEffect(() => { bottomRef.current?.scrollIntoView({ behavior: 'smooth' }) }, [messages])

  const sendTyping = (isTyping: boolean) => {
    if (!chatId || !wsSend) return
    wsSend(isTyping ? 'user:typing' : 'user:stop_typing', { chatId })
  }

  const handleTextChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setText(e.target.value)
    if (!chatId || !wsSend) return
    if (!isTypingRef.current) {
      isTypingRef.current = true
      sendTyping(true)
    }
    if (typingTimerRef.current) clearTimeout(typingTimerRef.current)
    typingTimerRef.current = window.setTimeout(() => {
      isTypingRef.current = false
      sendTyping(false)
    }, 3000)
  }

  const sendMessage = async () => {
    if (!text.trim() || !chatId || sending) return
    setSending(true)
    if (isTypingRef.current) { isTypingRef.current = false; sendTyping(false) }
    if (typingTimerRef.current) { clearTimeout(typingTimerRef.current) }
    const res = await api.call('POST', `/api/chats/${chatId}/messages`, {
      content: text.trim(),
      type: 'text',
    })
    setSending(false)
    if (!res.error) {
      setText('')
      loadMessages()
    }
  }

  const handleFileSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file || !chatId) return
    setUploading(true)
    const fd = new FormData()
    fd.append('file', file)
    const res = await api.call('POST', `/api/chats/${chatId}/messages/file`, fd, true)
    setUploading(false)
    if (!res.error) loadMessages()
    if (fileInputRef.current) fileInputRef.current.value = ''
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); sendMessage() }
  }

  if (!chatId) {
    return (
      <div className="no-chat-selected">
        <div className="no-chat-icon">💬</div>
        <div className="no-chat-text">Select a chat to start messaging</div>
        <div className="no-chat-hint">Choose a conversation from the left sidebar</div>
      </div>
    )
  }

  return (
    <div className="message-area">
      <div className="message-header">
        <div className="msg-header-avatar">{chatName.charAt(0).toUpperCase()}</div>
        <div className="msg-header-info">
          <div className="msg-header-name">{chatName}</div>
          <div className="msg-header-status">{messages.length} messages</div>
        </div>
      </div>

      <div className="messages-container">
        {messages.length === 0 && (
          <div className="no-messages">
            <div className="no-msg-icon">👋</div>
            <div className="no-msg-text">No messages yet</div>
            <div className="no-msg-hint">Send a message to start the conversation</div>
          </div>
        )}
        {messages.map((msg: any) => (
          <MessageBubble key={msg.id || msg.ID} message={msg} />
        ))}
        <div ref={bottomRef} />
      </div>

      <div className="message-input-area">
        <input
          ref={fileInputRef}
          type="file"
          id="file-upload-input"
          style={{ display: 'none' }}
          onChange={handleFileSelect}
        />
        <button
          className="attach-btn"
          onClick={() => fileInputRef.current?.click()}
          disabled={uploading}
          title="Attach file"
        >
          {uploading ? <span className="spinner-sm" /> : '📎'}
        </button>
        <input
          className="message-input"
          type="text"
          value={text}
          onChange={handleTextChange}
          onKeyDown={handleKeyDown}
          placeholder="Type a message..."
          disabled={sending || uploading}
        />
        <button className="send-btn" onClick={sendMessage} disabled={sending || !text.trim() || uploading}>
          {sending ? <span className="spinner-sm" /> : '➤'}
        </button>
      </div>
    </div>
  )
}
