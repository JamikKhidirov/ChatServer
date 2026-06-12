import { useState } from 'react'
import { Card } from '../../components/common/Card'
import { FormRow, Input, Select, TextArea, Checkbox } from '../../components/common/FormRow'
import { LoadingButton } from '../../components/common/LoadingButton'
import { ResultBox } from '../../components/common/ResultBox'
import { useToast } from '../../components/common/Toast'
import * as api from './chats'

export default function ChatsPage() {
  const [result, setResult] = useState<unknown>(null)
  const { toast } = useToast()

  const handleAction = async (fn: () => Promise<any>) => {
    const res = await fn()
    setResult(res)
    if (!res.error) toast('Success', 'success')
    return res
  }

  const [userId, setUserId] = useState('')
  const [groupName, setGroupName] = useState('')
  const [participantsJson, setParticipantsJson] = useState('')
  const [chatId, setChatId] = useState('')
  const [chatIdAction, setChatIdAction] = useState('')
  const [participantChatId, setParticipantChatId] = useState('')
  const [participantUserId, setParticipantUserId] = useState('')
  const [role, setRole] = useState('member')
  const [notifChatId, setNotifChatId] = useState('')
  const [muted, setMuted] = useState(false)
  const [pinChatId, setPinChatId] = useState('')
  const [archiveChatId, setArchiveChatId] = useState('')

  return (
    <div className="page-content">
      <h1 className="page-title">Chats</h1>
      <p className="page-subtitle">Create, list, and manage chats</p>

      <Card title="List Chats" badge="GET">
        <LoadingButton variant="info" onClick={() => handleAction(() => api.listChats())}>List Chats</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Create Private Chat" badge="POST">
        <FormRow>
          <Input placeholder="User ID" value={userId} onChange={e => setUserId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="success" onClick={() => handleAction(() => api.createChat({ type: 'private', participantIds: [userId] }))}>Create Private Chat</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Create Group Chat" badge="POST">
        <FormRow>
          <Input placeholder="Group Name" value={groupName} onChange={e => setGroupName(e.target.value)} />
          <Input placeholder="Participant IDs (JSON array)" value={participantsJson} onChange={e => setParticipantsJson(e.target.value)} />
        </FormRow>
        <LoadingButton variant="success" onClick={() => {
          let ids: string[] = []
          try { ids = JSON.parse(participantsJson) } catch { ids = participantsJson ? [participantsJson] : [] }
          return handleAction(() => api.createChat({ type: 'group', name: groupName, participantIds: ids }))
        }}>Create Group Chat</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Get Chat" badge="GET">
        <FormRow>
          <Input placeholder="Chat ID" value={chatId} onChange={e => setChatId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getChat(chatId))}>Get Chat</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Chat ID Actions" badge="POST/DELETE">
        <FormRow>
          <Input placeholder="Chat ID" value={chatIdAction} onChange={e => setChatIdAction(e.target.value)} />
        </FormRow>
        <FormRow>
          <LoadingButton variant="warning" onClick={() => handleAction(() => api.hideChat(chatIdAction))}>Hide</LoadingButton>
          <LoadingButton variant="danger" onClick={() => handleAction(() => api.deleteChat(chatIdAction))}>Delete</LoadingButton>
          <LoadingButton variant="warning" onClick={() => handleAction(() => api.leaveChat(chatIdAction))}>Leave</LoadingButton>
          <LoadingButton variant="info" onClick={() => handleAction(() => api.markRead(chatIdAction))}>Mark Read</LoadingButton>
        </FormRow>
        <ResultBox data={result} />
      </Card>

      <Card title="Participant Management" badge="POST/DELETE">
        <FormRow>
          <Input placeholder="Chat ID" value={participantChatId} onChange={e => setParticipantChatId(e.target.value)} />
          <Input placeholder="User ID" value={participantUserId} onChange={e => setParticipantUserId(e.target.value)} />
        </FormRow>
        <FormRow>
          <Select value={role} onChange={e => setRole(e.target.value)}>
            <option value="member">Member</option>
            <option value="admin">Admin</option>
          </Select>
        </FormRow>
        <FormRow>
          <LoadingButton variant="success" onClick={() => handleAction(() => api.addParticipant(participantChatId, participantUserId))}>Add</LoadingButton>
          <LoadingButton variant="danger" onClick={() => handleAction(() => api.removeParticipant(participantChatId, participantUserId))}>Remove</LoadingButton>
          <LoadingButton variant="warning" onClick={() => handleAction(() => api.setRole(participantChatId, participantUserId, role))}>Set Role</LoadingButton>
        </FormRow>
        <ResultBox data={result} />
      </Card>

      <Card title="Notifications" badge="PUT">
        <FormRow>
          <Input placeholder="Chat ID" value={notifChatId} onChange={e => setNotifChatId(e.target.value)} />
          <Checkbox label="Muted" checked={muted} onChange={e => setMuted(e.target.checked)} />
        </FormRow>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.setNotificationMuted(notifChatId, muted))}>Set Muted</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Pin / Unpin / Archive" badge="POST">
        <FormRow>
          <Input placeholder="Chat ID" value={pinChatId} onChange={e => setPinChatId(e.target.value)} />
        </FormRow>
        <FormRow>
          <LoadingButton variant="info" onClick={() => handleAction(() => api.pinChat(pinChatId))}>Pin</LoadingButton>
          <LoadingButton variant="warning" onClick={() => handleAction(() => api.unpinChat(pinChatId))}>Unpin</LoadingButton>
          <LoadingButton variant="info" onClick={() => handleAction(() => api.archiveChat(pinChatId))}>Archive</LoadingButton>
          <LoadingButton variant="warning" onClick={() => handleAction(() => api.unarchiveChat(pinChatId))}>Unarchive</LoadingButton>
        </FormRow>
        <ResultBox data={result} />
      </Card>
    </div>
  )
}
