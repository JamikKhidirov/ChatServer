import { useState, useEffect, useRef, useCallback } from 'react'
import { api } from '../api/client'
import MessageBubble from './MessageBubble'

interface Props {
  chatId: string | null
  chatInfo: any
  token: string
}

export default function MessageArea({ chatId, chatInfo }: Props) {
  const [messages, setMessages] = useState<any[]>([])
  const [text, setText] = useState('')
  const [sending, setSending] = useState(false)
  const bottomRef = useRef<HTMLDivElement>(null)
  const [chatName, setChatName] = useState('')

  useEffect(() => {
    if (chatInfo?.name || chatInfo?.Name) {
      setChatName(chatInfo.name || chatInfo.Name)
    } else if (chatInfo?.participants) {
      const others = (chatInfo.participants || []).filter((p: any) => p.id !== api.getToken())
      setChatName(others.map((p: any) => p.displayName || p.username).join(', ') || 'Chat')
    }
  }, [chatInfo])

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

  const sendMessage = async () => {
    if (!text.trim() || !chatId || sending) return
    setSending(true)
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
          className="message-input"
          type="text"
          value={text}
          onChange={e => setText(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="Type a message..."
          disabled={sending}
        />
        <button className="send-btn" onClick={sendMessage} disabled={sending || !text.trim()}>
          {sending ? <span className="spinner-sm" /> : '➤'}
        </button>
      </div>
    </div>
  )
}