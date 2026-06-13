import { api } from '../api/client'

interface Props { message: any }

export default function MessageBubble({ message }: Props) {
  const msg = message
  const senderId = msg.senderId || msg.sender?.id || msg.SenderID || msg.sender?.ID || ''
  const content = msg.content || msg.Content || ''
  const createdAt = msg.createdAt || msg.CreatedAt || ''
  const type = msg.type || msg.Type || 'text'
  const time = createdAt ? new Date(createdAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) : ''

  const token = api.getToken()
  const userId = token ? JSON.parse(atob(token.split('.')[1])).sub : ''
  const isOwn = senderId === userId

  return (
    <div className={`message-row ${isOwn ? 'own' : 'other'}`}>
      <div className={`message-bubble ${isOwn ? 'own' : 'other'}`}>
        <div className="msg-content">{content}</div>
        <div className="msg-time">{time}</div>
      </div>
    </div>
  )
}