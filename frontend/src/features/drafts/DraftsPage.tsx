import { useState } from 'react'
import { Card } from '../../components/common/Card'
import { FormRow, Input, Select, TextArea, Checkbox } from '../../components/common/FormRow'
import { LoadingButton } from '../../components/common/LoadingButton'
import { ResultBox } from '../../components/common/ResultBox'
import { useToast } from '../../components/common/Toast'
import * as api from './drafts'

export default function DraftsPage() {
  const [result, setResult] = useState<unknown>(null)
  const { toast } = useToast()

  const handleAction = async (fn: () => Promise<any>) => {
    const res = await fn()
    setResult(res)
    if (!res.error) toast('Success', 'success')
    return res
  }

  const [chatId, setChatId] = useState('')
  const [content, setContent] = useState('')

  return (
    <div className="page-content">
      <h1 className="page-title">Drafts</h1>
      <p className="page-subtitle">Save and retrieve message drafts</p>

      <Card title="Save Draft" badge="POST">
        <FormRow>
          <Input placeholder="Chat ID *" value={chatId} onChange={e => setChatId(e.target.value)} />
        </FormRow>
        <FormRow>
          <TextArea placeholder="Draft content" value={content} onChange={e => setContent(e.target.value)} />
        </FormRow>
        <LoadingButton variant="success" onClick={() => handleAction(() => api.saveDraft({ chatId, content }))}>Save Draft</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Get Draft" badge="GET">
        <FormRow>
          <Input placeholder="Chat ID *" value={chatId} onChange={e => setChatId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getDraft(chatId))}>Get Draft</LoadingButton>
        <ResultBox data={result} />
      </Card>
    </div>
  )
}
