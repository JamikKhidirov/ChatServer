import { useState } from 'react'
import { Card } from '../../components/common/Card'
import { FormRow, Input, TextArea, Select } from '../../components/common/FormRow'
import { LoadingButton } from '../../components/common/LoadingButton'
import { ResultBox } from '../../components/common/ResultBox'
import { useToast } from '../../components/common/Toast'
import * as api from './scheduled'

export default function ScheduledPage() {
  const [result, setResult] = useState<unknown>(null)
  const { toast } = useToast()

  const handleAction = async (fn: () => Promise<any>, successMsg?: string) => {
    const res = await fn()
    setResult(res)
    if (!res.error && successMsg) toast(successMsg, 'success')
    return res
  }

  const [chatId, setChatId] = useState('')
  const [content, setContent] = useState('')
  const [type, setType] = useState('text')
  const [scheduledAt, setScheduledAt] = useState('')
  const [cancelId, setCancelId] = useState('')

  return (
    <div className="page-content">
      <h1 className="page-title">Scheduled Messages</h1>
      <p className="page-subtitle">Schedule, view, and cancel messages</p>
      <Card title="Schedule Message" badge="POST">
        <FormRow>
          <Input placeholder="Chat ID" value={chatId} onChange={e => setChatId(e.target.value)} />
        </FormRow>
        <FormRow>
          <TextArea placeholder="Message content" value={content} onChange={e => setContent(e.target.value)} />
        </FormRow>
        <FormRow>
          <Select value={type} onChange={e => setType(e.target.value)}>
            <option value="text">Text</option>
            <option value="image">Image</option>
            <option value="video">Video</option>
            <option value="audio">Audio</option>
            <option value="file">File</option>
          </Select>
          <Input placeholder="Scheduled At (RFC3339)" value={scheduledAt} onChange={e => setScheduledAt(e.target.value)} />
        </FormRow>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.scheduleMessage({ chat_id: chatId, content, scheduled_at: scheduledAt }), 'Message scheduled')}>Schedule</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="Get Scheduled" badge="GET">
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getScheduled())}>Get Scheduled Messages</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="Cancel Scheduled" badge="DELETE">
        <FormRow>
          <Input placeholder="Scheduled Message ID" value={cancelId} onChange={e => setCancelId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="danger" onClick={() => handleAction(() => api.cancelScheduled(cancelId), 'Scheduled message cancelled')}>Cancel</LoadingButton>
        <ResultBox data={result} />
      </Card>
    </div>
  )
}
