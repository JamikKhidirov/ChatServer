import { useState, useEffect, useCallback, useRef } from 'react'
import { api } from '../api/client'
import ChatList from '../components/ChatList'
import MessageArea from '../components/MessageArea'
import CreateGroupModal from '../components/CreateGroupModal'
import useWebSocket from '../hooks/useWebSocket'
import * as authApi from '../features/auth/auth'

interface Props { onLogout: () => void }

export default function MainScreen({ onLogout }: Props) {
  const [chats, setChats] = useState<any[]>([])
  const [activeChat, setActiveChat] = useState<string | null>(null)
  const [searchQuery, setSearchQuery] = useState('')
  const [showSearch, setShowSearch] = useState(false)
  const [searchResults, setSearchResults] = useState<any[]>([])
  const [activeChatInfo, setActiveChatInfo] = useState<any>(null)
  const [wsMessage, setWsMessage] = useState<any>(null)
  const [showCreateGroup, setShowCreateGroup] = useState(false)
  const [typingUsers, setTypingUsers] = useState<Record<string, { userId: string; name: string }[]>>({})

  const token = api.getToken()
  const chatRef = useRef(activeChat)
  chatRef.current = activeChat

  const { send: wsSend } = useWebSocket((msg: any) => {
    if (msg.type === 'message:new' && msg.payload?.chatId === chatRef.current) {
      setWsMessage(msg.payload)
    }
    if (msg.type === 'chat:created' || msg.type === 'chat:deleted') {
      loadChats()
    }
    if (msg.type === 'user:typing' && msg.payload?.chatId === chatRef.current) {
      setTypingUsers(prev => {
        const chatId = msg.payload.chatId
        const existing = prev[chatId] || []
        if (existing.some(u => u.userId === msg.payload.userId)) return prev
        return { ...prev, [chatId]: [...existing, { userId: msg.payload.userId, name: msg.payload.userId }] }
      })
    }
    if (msg.type === 'user:stop_typing' && msg.payload?.chatId === chatRef.current) {
      setTypingUsers(prev => {
        const chatId = msg.payload.chatId
        const existing = (prev[chatId] || []).filter(u => u.userId !== msg.payload.userId)
        if (existing.length === 0) {
          const { [chatId]: _, ...rest } = prev
          return rest
        }
        return { ...prev, [chatId]: existing }
      })
    }
  }, true)

  const loadChats = useCallback(async () => {
    const res = await api.call('GET', '/api/chats')
    const d = res.data as any
    if (!res.error && Array.isArray(d?.data)) {
      setChats(d.data)
    } else if (!res.error && Array.isArray(d)) {
      setChats(d)
    }
  }, [])

  useEffect(() => { loadChats() }, [loadChats])

  const handleSearch = async () => {
    if (!searchQuery.trim()) return
    const res = await api.call('GET', '/api/users/search?q=' + encodeURIComponent(searchQuery))
    const d = res.data as any
    if (!res.error && Array.isArray(d?.data)) {
      setSearchResults(d.data)
    }
  }

  const startChat = async (userId: string) => {
    const res = await api.call('POST', '/api/chats/start/' + userId)
    if (!res.error) {
      const d = res.data as any
      const chat = d?.data || d
      setActiveChat(chat.id || chat.ID)
      setActiveChatInfo(chat)
      setShowSearch(false)
      setSearchQuery('')
      loadChats()
    }
  }

  const selectChat = async (chat: any) => {
    setActiveChat(chat.id || chat.ID)
    setActiveChatInfo(chat)
  }

  const currentTyping = activeChat ? typingUsers[activeChat] || [] : []
  const typingText = currentTyping.length > 0
    ? currentTyping.map(u => u.name).join(', ') + (currentTyping.length === 1 ? ' is typing...' : ' are typing...')
    : ''

  return (
    <div className="main-screen">
      {showCreateGroup && (
        <CreateGroupModal
          onClose={() => setShowCreateGroup(false)}
          onCreated={() => { setActiveChat(null); loadChats() }}
        />
      )}
      <div className="sidebar">
        <div className="sidebar-header">
          <div className="sidebar-brand">
            <span className="brand-icon">📨</span>
            <span className="brand-name">Messages</span>
          </div>
          <div className="sidebar-actions">
            <button className="icon-btn" onClick={() => setShowSearch(!showSearch)} title="Search users">🔍</button>
            <button className="icon-btn" onClick={() => setShowCreateGroup(true)} title="New group">➕</button>
            <button className="icon-btn" onClick={onLogout} title="Logout">🚪</button>
          </div>
        </div>

        {showSearch && (
          <div className="search-box">
            <input
              type="text" value={searchQuery}
              onChange={e => setSearchQuery(e.target.value)}
              onKeyDown={e => e.key === 'Enter' && handleSearch()}
              placeholder="Search users..." autoFocus
            />
            {searchResults.length > 0 && (
              <div className="search-results">
                {searchResults.map((u: any) => (
                  <div key={u.id || u.ID} className="search-item" onClick={() => startChat(u.id || u.ID)}>
                    <div className="search-avatar">{u.displayName?.charAt(0) || '?'}</div>
                    <div>
                      <div className="search-name">{u.displayName || u.Username}</div>
                      <div className="search-status">@{u.username || u.Username}</div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}

        <ChatList
          chats={chats}
          activeChat={activeChat}
          onSelect={selectChat}
        />
      </div>

      <div style={{ flex: 1, display: 'flex', flexDirection: 'column', position: 'relative' }}>
        <MessageArea
          chatId={activeChat}
          chatInfo={activeChatInfo}
          token={token}
          wsMessage={wsMessage}
          onWsConsumed={() => setWsMessage(null)}
          wsSend={wsSend}
        />
        {typingText && (
          <div className="typing-indicator">{typingText}</div>
        )}
      </div>
    </div>
  )
}
