import { useState } from 'react'
import { Card } from '../../components/common/Card'
import { FormRow, Input, Select } from '../../components/common/FormRow'
import { LoadingButton } from '../../components/common/LoadingButton'
import { ResultBox } from '../../components/common/ResultBox'
import { useToast } from '../../components/common/Toast'
import * as api from './calls'

export default function CallsPage() {
  const [result, setResult] = useState<unknown>(null)
  const { toast } = useToast()

  const handleAction = async (fn: () => Promise<any>, successMsg?: string) => {
    const res = await fn()
    setResult(res)
    if (!res.error && successMsg) toast(successMsg, 'success')
    return res
  }

  const [initChatId, setInitChatId] = useState('')
  const [callType, setCallType] = useState<'audio' | 'video'>('audio')
  const [respondCallId, setRespondCallId] = useState('')
  const [respondAction, setRespondAction] = useState('accept')
  const [endCallId, setEndCallId] = useState('')
  const [getCallId, setGetCallId] = useState('')
  const [historyChatId, setHistoryChatId] = useState('')

  return (
    <div className="page-content">
      <h1 className="page-title">Calls</h1>
      <p className="page-subtitle">Initiate, respond to, and manage calls</p>
      <Card title="Initiate Call" badge="POST">
        <FormRow>
          <Input placeholder="Chat ID" value={initChatId} onChange={e => setInitChatId(e.target.value)} />
          <Select value={callType} onChange={e => setCallType(e.target.value as 'audio' | 'video')}>
            <option value="audio">Audio</option>
            <option value="video">Video</option>
          </Select>
        </FormRow>
        <LoadingButton variant="success" onClick={() => handleAction(() => api.initiateCall({ chat_id: initChatId, type: callType }), 'Call initiated')}>Initiate Call</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="Respond to Call" badge="POST">
        <FormRow>
          <Input placeholder="Call ID" value={respondCallId} onChange={e => setRespondCallId(e.target.value)} />
          <Select value={respondAction} onChange={e => setRespondAction(e.target.value)}>
            <option value="accept">Accept</option>
            <option value="reject">Reject</option>
          </Select>
        </FormRow>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.respondCall(respondCallId, respondAction as 'accept' | 'reject' | 'ignore'), 'Call response sent')}>Respond</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="End Call" badge="POST">
        <FormRow>
          <Input placeholder="Call ID" value={endCallId} onChange={e => setEndCallId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="danger" onClick={() => handleAction(() => api.endCall(endCallId), 'Call ended')}>End Call</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="Get Call" badge="GET">
        <FormRow>
          <Input placeholder="Call ID" value={getCallId} onChange={e => setGetCallId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getCall(getCallId))}>Get Call</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="Call History" badge="GET">
        <FormRow>
          <Input placeholder="Chat ID" value={historyChatId} onChange={e => setHistoryChatId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getCallHistory(historyChatId))}>Get Call History</LoadingButton>
        <ResultBox data={result} />
      </Card>
    </div>
  )
}
