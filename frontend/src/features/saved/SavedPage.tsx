import { useState } from 'react'
import { Card } from '../../components/common/Card'
import { FormRow, Input, Select, TextArea, Checkbox } from '../../components/common/FormRow'
import { LoadingButton } from '../../components/common/LoadingButton'
import { ResultBox } from '../../components/common/ResultBox'
import { useToast } from '../../components/common/Toast'
import * as api from './saved'

export default function SavedPage() {
  const [result, setResult] = useState<unknown>(null)
  const { toast } = useToast()

  const handleAction = async (fn: () => Promise<any>) => {
    const res = await fn()
    setResult(res)
    if (!res.error) toast('Success', 'success')
    return res
  }

  const [messageId, setMessageId] = useState('')
  const [chatId, setChatId] = useState('')
  const [limit, setLimit] = useState('50')
  const [offset, setOffset] = useState('0')
  const [deleteId, setDeleteId] = useState('')

  return (
    <div className="page-content">
      <h1 className="page-title">Saved Messages</h1>
      <p className="page-subtitle">Save, view, and delete saved messages</p>

      <Card title="Save Message" badge="POST">
        <FormRow>
          <Input placeholder="Message ID *" value={messageId} onChange={e => setMessageId(e.target.value)} />
          <Input placeholder="Chat ID *" value={chatId} onChange={e => setChatId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="success" onClick={() => handleAction(() => api.saveMessage(messageId, chatId))}>Save Message</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Get Saved Messages" badge="GET">
        <FormRow>
          <Input placeholder="Limit" value={limit} onChange={e => setLimit(e.target.value)} />
          <Input placeholder="Offset" value={offset} onChange={e => setOffset(e.target.value)} />
        </FormRow>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getSavedMessages(Number(limit) || undefined, Number(offset) || undefined))}>Get Saved</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Delete Saved Message" badge="DELETE">
        <FormRow>
          <Input placeholder="Saved Message ID *" value={deleteId} onChange={e => setDeleteId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="danger" onClick={() => handleAction(() => api.deleteSavedMessage(deleteId))}>Delete</LoadingButton>
        <ResultBox data={result} />
      </Card>
    </div>
  )
}
