import { useState, useEffect, useCallback } from 'react'
import { api } from '../api/client'
import ChatList from '../components/ChatList'
import MessageArea from '../components/MessageArea'
import * as authApi from '../features/auth/auth'

interface Props { onLogout: () => void }

export default function MainScreen({ onLogout }: Props) {
  const [chats, setChats] = useState<any[]>([])
  const [activeChat, setActiveChat] = useState<string | null>(null)
  const [recipient, setRecipient] = useState('')
  const [searchQuery, setSearchQuery] = useState('')
  const [showSearch, setShowSearch] = useState(false)
  const [searchResults, setSearchResults] = useState<any[]>([])
  const [activeChatInfo, setActiveChatInfo] = useState<any>(null)

  const token = api.getToken()

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

  const createGroup = async () => {
    const name = prompt('Group name:')
    if (!name) return
    const members = prompt('Participant IDs (comma-separated):')
    if (!members) return
    const ids = members.split(',').map(s => s.trim()).filter(Boolean)
    const res = await api.call('POST', '/api/chats', {
      type: 'group', name, participantIds: ids,
    })
    if (!res.error) {
      setActiveChat(null)
      loadChats()
    }
  }

  return (
    <div className="main-screen">
      <div className="sidebar">
        <div className="sidebar-header">
          <div className="sidebar-brand">
            <span className="brand-icon">📨</span>
            <span className="brand-name">Messages</span>
          </div>
          <div className="sidebar-actions">
            <button className="icon-btn" onClick={() => setShowSearch(!showSearch)} title="Search users">🔍</button>
            <button className="icon-btn" onClick={createGroup} title="New group">➕</button>
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

      <MessageArea
        chatId={activeChat}
        chatInfo={activeChatInfo}
        token={token}
      />
    </div>
  )
}