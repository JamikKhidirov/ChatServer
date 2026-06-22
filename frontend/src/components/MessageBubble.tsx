import { api } from '../api/client'

interface Props { message: any }

function formatFileSize(bytes: number): string {
  if (!bytes) return ''
  const units = ['B', 'KB', 'MB', 'GB']
  let i = 0; let size = bytes
  while (size >= 1024 && i < units.length - 1) { size /= 1024; i++ }
  return `${size.toFixed(1)} ${units[i]}`
}

function getFileUrl(path?: string): string {
  if (!path) return ''
  if (path.startsWith('http')) return path
  const base = window.location.protocol + '//' + window.location.host
  return base + (path.startsWith('/') ? '' : '/') + path
}

export default function MessageBubble({ message }: Props) {
  const msg = message
  const senderId = msg.senderId || msg.sender?.id || msg.SenderID || msg.sender?.ID || ''
  const content = msg.content || msg.Content || ''
  const caption = msg.caption || msg.Caption || ''
  const createdAt = msg.createdAt || msg.CreatedAt || ''
  const type = msg.type || msg.Type || 'text'
  const fileName = msg.fileName || msg.FileName || ''
  const fileUrl = msg.fileUrl || msg.FileUrl || msg.FileURL || ''
  const fileSize = msg.fileSize || msg.FileSize || 0
  const mimeType = msg.mimeType || msg.MimeType || ''
  const latitude = msg.latitude || msg.Latitude || 0
  const longitude = msg.longitude || msg.Longitude || 0
  const locationTitle = msg.locationTitle || msg.LocationTitle || ''
  const effect = msg.effect || msg.Effect || ''
  const width = msg.width || msg.Width || 0
  const height = msg.height || msg.Height || 0

  const time = createdAt ? new Date(createdAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) : ''

  const token = api.getToken()
  const userId = token ? JSON.parse(atob(token.split('.')[1])).sub : ''
  const isOwn = senderId === userId

  const renderContent = () => {
    switch (type) {
      case 'image': {
        const src = getFileUrl(fileUrl)
        return (
          <div>
            <img src={src} alt={fileName || 'image'} style={{ maxWidth: '100%', maxHeight: 300, borderRadius: 12, display: 'block', cursor: 'pointer' }} onClick={() => window.open(src, '_blank')} />
            {caption && <div style={{ marginTop: 6, fontSize: 13 }}>{caption}</div>}
          </div>
        )
      }
      case 'gif': {
        const src = getFileUrl(fileUrl)
        return (
          <div>
            <img src={src} alt="gif" style={{ maxWidth: '100%', maxHeight: 250, borderRadius: 12, display: 'block' }} />
            {caption && <div style={{ marginTop: 6, fontSize: 13 }}>{caption}</div>}
          </div>
        )
      }
      case 'voice': {
        const src = getFileUrl(fileUrl)
        return (
          <div style={{ minWidth: 200 }}>
            <audio controls src={src} style={{ width: '100%', height: 40 }} />
            {caption && <div style={{ marginTop: 4, fontSize: 12, color: isOwn ? 'rgba(255,255,255,0.7)' : '#aaa' }}>{caption}</div>}
          </div>
        )
      }
      case 'video':
      case 'video_circle': {
        const src = getFileUrl(fileUrl)
        return (
          <div>
            <video controls src={src} style={{ maxWidth: '100%', maxHeight: 350, borderRadius: 12, display: 'block' }} />
            {caption && <div style={{ marginTop: 6, fontSize: 13 }}>{caption}</div>}
          </div>
        )
      }
      case 'audio': {
        const src = getFileUrl(fileUrl)
        return (
          <div style={{ minWidth: 220 }}>
            <div style={{ fontSize: 12, marginBottom: 4, color: isOwn ? 'rgba(255,255,255,0.8)' : '#bbb' }}>{fileName}</div>
            <audio controls src={src} style={{ width: '100%', height: 40 }} />
          </div>
        )
      }
      case 'location': {
        const mapsUrl = `https://www.google.com/maps?q=${latitude},${longitude}`
        return (
          <div>
            <div style={{ fontSize: 13, marginBottom: 4 }}>{locationTitle || 'Location'}</div>
            <div style={{ fontSize: 11, color: isOwn ? 'rgba(255,255,255,0.6)' : '#888' }}>{latitude.toFixed(4)}, {longitude.toFixed(4)}</div>
            <a href={mapsUrl} target="_blank" rel="noopener noreferrer" style={{ fontSize: 12, color: '#5865f2', textDecoration: 'none' }}>Open in Maps</a>
          </div>
        )
      }
      case 'file': {
        const src = getFileUrl(fileUrl)
        return (
          <div>
            <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
              <span style={{ fontSize: 28 }}>📎</span>
              <div>
                <div style={{ fontSize: 13, fontWeight: 600 }}>{fileName || 'File'}</div>
                {fileSize > 0 && <div style={{ fontSize: 11, color: isOwn ? 'rgba(255,255,255,0.6)' : '#888' }}>{formatFileSize(fileSize)}</div>}
              </div>
            </div>
            {caption && <div style={{ marginTop: 6, fontSize: 13 }}>{caption}</div>}
            <a href={src} target="_blank" rel="noopener noreferrer" download style={{ fontSize: 12, color: '#5865f2', textDecoration: 'none', display: 'inline-block', marginTop: 6 }}>Download</a>
          </div>
        )
      }
      default:
        return <div className="msg-content">{content}</div>
    }
  }

  return (
    <div className={`message-row ${isOwn ? 'own' : 'other'}`}>
      <div className={`message-bubble ${isOwn ? 'own' : 'other'} ${effect ? 'msg-effect-' + effect : ''}`}>
        {renderContent()}
        <div className="msg-time">{time}</div>
      </div>
    </div>
  )
}
