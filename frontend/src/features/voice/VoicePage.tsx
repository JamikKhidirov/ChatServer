import { useState } from 'react'
import { Card } from '../../components/common/Card'
import { FormRow, Input, Select, TextArea, Checkbox } from '../../components/common/FormRow'
import { LoadingButton } from '../../components/common/LoadingButton'
import { ResultBox } from '../../components/common/ResultBox'
import { useToast } from '../../components/common/Toast'
import * as api from './voice'

export default function VoicePage() {
  const [result, setResult] = useState<unknown>(null)
  const { toast } = useToast()

  const handleAction = async (fn: () => Promise<any>) => {
    const res = await fn()
    setResult(res)
    if (!res.error) toast('Success', 'success')
    return res
  }

  const [chatId, setChatId] = useState('')
  const [title, setTitle] = useState('')
  const [scheduledMins, setScheduledMins] = useState('0')
  const [voiceChatId, setVoiceChatId] = useState('')
  const [muted, setMuted] = useState(false)

  return (
    <div className="page-content">
      <h1 className="page-title">Voice Chats</h1>
      <p className="page-subtitle">Create, join, leave, and manage voice chat rooms</p>

      <Card title="Create Voice Chat" badge="POST">
        <FormRow>
          <Input placeholder="Chat ID *" value={chatId} onChange={e => setChatId(e.target.value)} />
          <Input placeholder="Title" value={title} onChange={e => setTitle(e.target.value)} />
          <Input placeholder="Scheduled In Minutes" type="number" value={scheduledMins} onChange={e => setScheduledMins(e.target.value)} />
        </FormRow>
        <LoadingButton variant="success" onClick={() => handleAction(() => api.createVoiceChat(chatId, { title: title || undefined, scheduledInMins: Number(scheduledMins) || undefined }))}>Create</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Active Voice Chats" badge="GET">
        <FormRow>
          <Input placeholder="Chat ID *" value={chatId} onChange={e => setChatId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getActiveVoiceChats(chatId))}>Get Active</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Voice Chat History" badge="GET">
        <FormRow>
          <Input placeholder="Chat ID *" value={chatId} onChange={e => setChatId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getVoiceChatHistory(chatId))}>Get History</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Get Voice Chat" badge="GET">
        <FormRow>
          <Input placeholder="Voice Chat ID *" value={voiceChatId} onChange={e => setVoiceChatId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getVoiceChat(voiceChatId))}>Get Voice Chat</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Join / Leave / End" badge="POST">
        <FormRow>
          <Input placeholder="Voice Chat ID *" value={voiceChatId} onChange={e => setVoiceChatId(e.target.value)} />
        </FormRow>
        <FormRow>
          <LoadingButton variant="success" onClick={() => handleAction(() => api.joinVoiceChat(voiceChatId))}>Join</LoadingButton>
          <LoadingButton variant="warning" onClick={() => handleAction(() => api.leaveVoiceChat(voiceChatId))}>Leave</LoadingButton>
          <LoadingButton variant="danger" onClick={() => handleAction(() => api.endVoiceChat(voiceChatId))}>End</LoadingButton>
        </FormRow>
        <ResultBox data={result} />
      </Card>

      <Card title="Mute / Unmute" badge="POST">
        <FormRow>
          <Input placeholder="Voice Chat ID *" value={voiceChatId} onChange={e => setVoiceChatId(e.target.value)} />
          <Checkbox label="Muted" checked={muted} onChange={e => setMuted(e.target.checked)} />
        </FormRow>
        <LoadingButton variant="warning" onClick={() => handleAction(() => api.muteParticipant(voiceChatId, muted))}>Set Mute</LoadingButton>
        <ResultBox data={result} />
      </Card>
    </div>
  )
}
