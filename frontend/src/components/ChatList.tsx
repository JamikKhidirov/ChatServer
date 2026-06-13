interface Props {
  chats: any[]
  activeChat: string | null
  onSelect: (chat: any) => void
}

export default function ChatList({ chats, activeChat, onSelect }: Props) {
  if (!chats.length) {
    return (
      <div className="chat-list-empty">
        <div className="empty-icon">💬</div>
        <div className="empty-text">No conversations yet</div>
        <div className="empty-hint">Search for users to start chatting</div>
      </div>
    )
  }

  return (
    <div className="chat-list">
      {chats.map(chat => {
        const id = chat.id || chat.ID
        const name = chat.name || chat.Name || 'Chat'
        const lastMsg = chat.lastMessage?.content || chat.LastMessage?.Content || ''
        const lastTime = chat.lastMessage?.createdAt || chat.LastMessage?.CreatedAt || ''
        const avatar = chat.avatarUrl || chat.AvatarURL || ''
        const unread = chat.unreadCount || chat.UnreadCount || 0
        const time = lastTime ? new Date(lastTime).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) : ''

        return (
          <div
            key={id}
            className={`chat-card ${activeChat === id ? 'active' : ''}`}
            onClick={() => onSelect(chat)}
          >
            <div className="chat-avatar">
              {avatar ? <img src={avatar} alt="" /> : <span>{name.charAt(0).toUpperCase()}</span>}
            </div>
            <div className="chat-info">
              <div className="chat-top">
                <span className="chat-name">{name}</span>
                {time && <span className="chat-time">{time}</span>}
              </div>
              <div className="chat-bottom">
                <span className="chat-preview">{lastMsg || 'No messages yet'}</span>
                {unread > 0 && <span className="chat-unread">{unread}</span>}
              </div>
            </div>
          </div>
        )
      })}
    </div>
  )
}