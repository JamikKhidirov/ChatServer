import { useState } from 'react'
import { Card } from '../../components/common/Card'
import { FormRow, Input, Select, TextArea, Checkbox } from '../../components/common/FormRow'
import { LoadingButton } from '../../components/common/LoadingButton'
import { ResultBox } from '../../components/common/ResultBox'
import { useToast } from '../../components/common/Toast'
import * as api from './polls'

export default function PollsPage() {
  const [result, setResult] = useState<unknown>(null)
  const { toast } = useToast()

  const handleAction = async (fn: () => Promise<any>) => {
    const res = await fn()
    setResult(res)
    if (!res.error) toast('Success', 'success')
    return res
  }

  const [chatId, setChatId] = useState('')
  const [question, setQuestion] = useState('')
  const [options, setOptions] = useState('')
  const [anonymous, setAnonymous] = useState(false)
  const [multiple, setMultiple] = useState(false)
  const [pollId, setPollId] = useState('')
  const [optionIndex, setOptionIndex] = useState('0')

  return (
    <div className="page-content">
      <h1 className="page-title">Polls</h1>
      <p className="page-subtitle">Create, vote, and close polls</p>

      <Card title="Create Poll" badge="POST">
        <FormRow>
          <Input placeholder="Chat ID *" value={chatId} onChange={e => setChatId(e.target.value)} />
          <Input placeholder="Question *" value={question} onChange={e => setQuestion(e.target.value)} />
        </FormRow>
        <FormRow>
          <TextArea placeholder='Options as JSON array, e.g. ["Yes","No","Maybe"]' value={options} onChange={e => setOptions(e.target.value)} />
        </FormRow>
        <FormRow>
          <Checkbox label="Anonymous" checked={anonymous} onChange={e => setAnonymous(e.target.checked)} />
          <Checkbox label="Multiple Choice" checked={multiple} onChange={e => setMultiple(e.target.checked)} />
        </FormRow>
        <LoadingButton variant="success" onClick={() => {
          let opts: string[] = []
          try { opts = JSON.parse(options) } catch { opts = options ? [options] : [] }
          return handleAction(() => api.createPoll({ chatId, question, options: opts, isAnonymous: anonymous, multipleChoice: multiple }))
        }}>Create Poll</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Vote" badge="POST">
        <FormRow>
          <Input placeholder="Poll ID *" value={pollId} onChange={e => setPollId(e.target.value)} />
          <Input placeholder="Option Index" type="number" value={optionIndex} onChange={e => setOptionIndex(e.target.value)} />
        </FormRow>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.votePoll(pollId, Number(optionIndex)))}>Vote</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Close Poll" badge="POST">
        <FormRow>
          <Input placeholder="Poll ID *" value={pollId} onChange={e => setPollId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="danger" onClick={() => handleAction(() => api.closePoll(pollId))}>Close Poll</LoadingButton>
        <ResultBox data={result} />
      </Card>
    </div>
  )
}
