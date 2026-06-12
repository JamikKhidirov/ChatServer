import { useState } from 'react'
import { Card } from '../../components/common/Card'
import { FormRow, Input } from '../../components/common/FormRow'
import { LoadingButton } from '../../components/common/LoadingButton'
import { ResultBox } from '../../components/common/ResultBox'
import { useToast } from '../../components/common/Toast'
import * as api from './sessions'

export default function SessionsPage() {
  const [result, setResult] = useState<unknown>(null)
  const { toast } = useToast()

  const handleAction = async (fn: () => Promise<any>, successMsg?: string) => {
    const res = await fn()
    setResult(res)
    if (!res.error && successMsg) toast(successMsg, 'success')
    return res
  }

  const [sessionId, setSessionId] = useState('')

  return (
    <div className="page-content">
      <h1 className="page-title">Sessions</h1>
      <p className="page-subtitle">Manage your active sessions</p>
      <Card title="Get Sessions" badge="GET">
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getSessions())}>Get Sessions</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="Delete Session" badge="DELETE">
        <FormRow>
          <Input placeholder="Session ID" value={sessionId} onChange={e => setSessionId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="danger" onClick={() => handleAction(() => api.deleteSession(sessionId), 'Session deleted')}>Delete Session</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="Delete All Sessions" badge="DELETE">
        <LoadingButton variant="danger" onClick={() => handleAction(() => api.deleteAllSessions(), 'All sessions deleted')}>Delete All Sessions</LoadingButton>
        <ResultBox data={result} />
      </Card>
    </div>
  )
}
