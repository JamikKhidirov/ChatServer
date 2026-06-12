import { useState } from 'react'
import { Card } from '../../components/common/Card'
import { FormRow, Input, Select, TextArea, Checkbox } from '../../components/common/FormRow'
import { LoadingButton } from '../../components/common/LoadingButton'
import { ResultBox } from '../../components/common/ResultBox'
import { useToast } from '../../components/common/Toast'
import * as api from './messages'

const EMOJIS = ['👍', '❤️', '😂', '😮', '😢', '🙏', '🎉', '🔥', '💯', '✅']

export default function MessagesPage() {
  const [result, setResult] = useState<unknown>(null)
  const { toast } = useToast()

  const handleAction = async (fn: () => Promise<any>) => {
    const res = await fn()
    setResult(res)
    if (!res.error) toast('Success', 'success')
    return res
  }

  const [listChatId, setListChatId] = useState('')
  const [limit, setLimit] = useState('50')
  const [offset, setOffset] = useState('0')
  const [sendChatId, setSendChatId] = useState('')
  const [content, setContent] = useState('')
  const [msgType, setMsgType] = useState('text')
  const [effect, setEffect] = useState('')
  const [fileChatId, setFileChatId] = useState('')
  const [videoChatId, setVideoChatId] = useState('')
  const [caption, setCaption] = useState('')
  const [editMsgId, setEditMsgId] = useState('')
  const [editContent, setEditContent] = useState('')
  const [deleteMsgId, setDeleteMsgId] = useState('')
  const [reactionMsgId, setReactionMsgId] = useState('')
  const [pinMsgId, setPinMsgId] = useState('')
  const [pinValue, setPinValue] = useState(true)
  const [forwardMsgId, setForwardMsgId] = useState('')
  const [forwardFromChat, setForwardFromChat] = useState('')
  const [forwardToChat, setForwardToChat] = useState('')
  const [starMsgId, setStarMsgId] = useState('')
  const [selfDestructMsgId, setSelfDestructMsgId] = useState('')
  const [selfDestructSec, setSelfDestructSec] = useState('5')
  const [historyMsgId, setHistoryMsgId] = useState('')
  const [mediaChatId, setMediaChatId] = useState('')
  const [mediaType, setMediaType] = useState('')
  const [exportChatId, setExportChatId] = useState('')
  const [forMeMsgId, setForMeMsgId] = useState('')
  const [searchQuery, setSearchQuery] = useState('')

  return (
    <div className="page-content">
      <h1 className="page-title">Messages</h1>
      <p className="page-subtitle">Send, edit, delete, and manage messages</p>

      <Card title="List Messages" badge="GET">
        <FormRow>
          <Input placeholder="Chat ID *" value={listChatId} onChange={e => setListChatId(e.target.value)} />
          <Input placeholder="Limit" value={limit} onChange={e => setLimit(e.target.value)} />
          <Input placeholder="Offset" value={offset} onChange={e => setOffset(e.target.value)} />
        </FormRow>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.listMessages(listChatId, Number(limit) || undefined, Number(offset) || undefined))}>List Messages</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Send Message" badge="POST">
        <FormRow>
          <Input placeholder="Chat ID *" value={sendChatId} onChange={e => setSendChatId(e.target.value)} />
          <Select value={msgType} onChange={e => setMsgType(e.target.value)}>
            <option value="text">Text</option>
            <option value="image">Image</option>
            <option value="gif">GIF</option>
            <option value="voice">Voice</option>
            <option value="video">Video</option>
            <option value="audio">Audio</option>
          </Select>
        </FormRow>
        <FormRow>
          <TextArea placeholder="Content *" value={content} onChange={e => setContent(e.target.value)} />
          <Select value={effect} onChange={e => setEffect(e.target.value)}>
            <option value="">No Effect</option>
            <option value="confetti">Confetti</option>
            <option value="fireworks">Fireworks</option>
            <option value="hearts">Hearts</option>
            <option value="balloons">Balloons</option>
            <option value="stars">Stars</option>
          </Select>
        </FormRow>
        <LoadingButton variant="success" onClick={() => handleAction(() => api.sendMessage(sendChatId, { content, type: msgType, effect: effect || undefined }))}>Send</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Upload File" badge="POST">
        <FormRow>
          <Input placeholder="Chat ID *" value={fileChatId} onChange={e => setFileChatId(e.target.value)} />
          <input type="file" style={{ color: '#9a9ab8', fontSize: 13, flex: 1 }} id="fileInput" />
        </FormRow>
        <LoadingButton variant="info" onClick={() => {
          const fileInput = document.getElementById('fileInput') as HTMLInputElement
          return handleAction(() => api.uploadFile(fileChatId, fileInput?.files?.[0]!))
        }}>Upload</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Video Circle" badge="POST">
        <FormRow>
          <Input placeholder="Chat ID *" value={videoChatId} onChange={e => setVideoChatId(e.target.value)} />
          <Input placeholder="Caption" value={caption} onChange={e => setCaption(e.target.value)} />
        </FormRow>
        <FormRow>
          <input type="file" accept="video/mp4" style={{ color: '#9a9ab8', fontSize: 13, flex: 1 }} id="videoInput" />
        </FormRow>
        <LoadingButton variant="info" onClick={() => {
          const vi = document.getElementById('videoInput') as HTMLInputElement
          return handleAction(() => api.uploadVideoCircle(videoChatId, vi?.files?.[0]!, caption))
        }}>Upload Video Circle</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Edit / Delete Message" badge="PUT/DELETE">
        <FormRow>
          <Input placeholder="Message ID" value={editMsgId} onChange={e => setEditMsgId(e.target.value)} />
          <Input placeholder="New Content" value={editContent} onChange={e => setEditContent(e.target.value)} />
        </FormRow>
        <FormRow>
          <LoadingButton variant="warning" onClick={() => handleAction(() => api.editMessage(editMsgId, editContent))}>Edit</LoadingButton>
          <LoadingButton variant="danger" onClick={() => handleAction(() => api.deleteMessage(deleteMsgId))}>Delete</LoadingButton>
          <Input placeholder="Message ID to delete" value={deleteMsgId} onChange={e => setDeleteMsgId(e.target.value)} />
        </FormRow>
        <ResultBox data={result} />
      </Card>

      <Card title="Reactions" badge="POST/DELETE">
        <FormRow>
          <Input placeholder="Message ID" value={reactionMsgId} onChange={e => setReactionMsgId(e.target.value)} />
        </FormRow>
        <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8, marginBottom: 10 }}>
          {EMOJIS.map(emoji => (
            <span key={emoji} onClick={() => handleAction(() => api.addReaction(reactionMsgId, emoji))} style={{ fontSize: 22, cursor: 'pointer', padding: '4px 6px', borderRadius: 6, background: 'rgba(255,255,255,0.04)', transition: 'all 0.2s' }}>{emoji}</span>
          ))}
        </div>
        <FormRow>
          <Input placeholder="Emoji to remove" id="removeEmoji" style={{ minWidth: 80, flex: '0 1 120px' }} />
          <LoadingButton variant="danger" onClick={() => {
            const em = (document.getElementById('removeEmoji') as HTMLInputElement)?.value
            return handleAction(() => api.removeReaction(reactionMsgId, em))
          }}>Remove Reaction</LoadingButton>
        </FormRow>
        <ResultBox data={result} />
      </Card>

      <Card title="Pin / Unpin / Forward" badge="PUT/POST">
        <FormRow>
          <Input placeholder="Message ID" value={pinMsgId} onChange={e => setPinMsgId(e.target.value)} />
          <Checkbox label="Pin" checked={pinValue} onChange={e => setPinValue(e.target.checked)} />
        </FormRow>
        <FormRow>
          <LoadingButton variant="info" onClick={() => handleAction(() => api.togglePin(pinMsgId, pinValue))}>Toggle Pin</LoadingButton>
        </FormRow>
        <FormRow>
          <Input placeholder="Msg ID" value={forwardMsgId} onChange={e => setForwardMsgId(e.target.value)} />
          <Input placeholder="From Chat ID" value={forwardFromChat} onChange={e => setForwardFromChat(e.target.value)} />
          <Input placeholder="To Chat ID" value={forwardToChat} onChange={e => setForwardToChat(e.target.value)} />
        </FormRow>
        <LoadingButton variant="success" onClick={() => handleAction(() => api.forwardMessage({ messageId: forwardMsgId, fromChatId: forwardFromChat, toChatId: forwardToChat }))}>Forward</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Star / Unstar / List Starred" badge="POST/GET">
        <FormRow>
          <Input placeholder="Message ID" value={starMsgId} onChange={e => setStarMsgId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="warning" onClick={() => handleAction(() => api.starMessage(starMsgId))}>Star</LoadingButton>
        <LoadingButton variant="warning" onClick={() => handleAction(() => api.unstarMessage(starMsgId))}>Unstar</LoadingButton>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getStarredMessages())}>List Starred</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Self-Destruct" badge="POST">
        <FormRow>
          <Input placeholder="Message ID" value={selfDestructMsgId} onChange={e => setSelfDestructMsgId(e.target.value)} />
          <Input placeholder="Seconds" type="number" value={selfDestructSec} onChange={e => setSelfDestructSec(e.target.value)} />
        </FormRow>
        <LoadingButton variant="danger" onClick={() => handleAction(() => api.selfDestruct(selfDestructMsgId, Number(selfDestructSec)))}>Set Self-Destruct</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Edit History" badge="GET">
        <FormRow>
          <Input placeholder="Message ID" value={historyMsgId} onChange={e => setHistoryMsgId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getMessageHistory(historyMsgId))}>Get History</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Chat Media" badge="GET">
        <FormRow>
          <Input placeholder="Chat ID" value={mediaChatId} onChange={e => setMediaChatId(e.target.value)} />
          <Select value={mediaType} onChange={e => setMediaType(e.target.value)}>
            <option value="">All</option>
            <option value="photo">Photo</option>
            <option value="video">Video</option>
            <option value="audio">Audio</option>
            <option value="document">Document</option>
          </Select>
        </FormRow>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getChatMedia(mediaChatId, mediaType || undefined))}>Get Media</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Export Chat" badge="GET">
        <FormRow>
          <Input placeholder="Chat ID" value={exportChatId} onChange={e => setExportChatId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.exportChat(exportChatId))}>Export</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Delete For Me" badge="DELETE">
        <FormRow>
          <Input placeholder="Message ID" value={forMeMsgId} onChange={e => setForMeMsgId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="danger" onClick={() => handleAction(() => api.deleteForMe(forMeMsgId))}>Delete For Me</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Search All Messages" badge="GET">
        <FormRow>
          <Input placeholder="Search query" value={searchQuery} onChange={e => setSearchQuery(e.target.value)} />
        </FormRow>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.searchAllMessages(searchQuery))}>Search</LoadingButton>
        <ResultBox data={result} />
      </Card>
    </div>
  )
}
