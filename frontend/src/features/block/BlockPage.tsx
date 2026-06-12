import { useState } from 'react'
import { Card } from '../../components/common/Card'
import { FormRow, Input } from '../../components/common/FormRow'
import { LoadingButton } from '../../components/common/LoadingButton'
import { ResultBox } from '../../components/common/ResultBox'
import { useToast } from '../../components/common/Toast'
import * as api from './block'

export default function BlockPage() {
  const [result, setResult] = useState<unknown>(null)
  const { toast } = useToast()

  const handleAction = async (fn: () => Promise<any>, successMsg?: string) => {
    const res = await fn()
    setResult(res)
    if (!res.error && successMsg) toast(successMsg, 'success')
    return res
  }

  const [blockUserId, setBlockUserId] = useState('')
  const [unblockUserId, setUnblockUserId] = useState('')

  return (
    <div className="page-content">
      <h1 className="page-title">Block Management</h1>
      <p className="page-subtitle">Block, unblock, and list blocked users</p>
      <Card title="Block User" badge="POST">
        <FormRow>
          <Input placeholder="User ID to block" value={blockUserId} onChange={e => setBlockUserId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="danger" onClick={() => handleAction(() => api.blockUser({ user_id: blockUserId }), 'User blocked')}>Block User</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="Unblock User" badge="POST">
        <FormRow>
          <Input placeholder="User ID to unblock" value={unblockUserId} onChange={e => setUnblockUserId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="warning" onClick={() => handleAction(() => api.unblockUser(unblockUserId), 'User unblocked')}>Unblock User</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="List Blocked" badge="GET">
        <LoadingButton variant="info" onClick={() => handleAction(() => api.listBlocked())}>List Blocked Users</LoadingButton>
        <ResultBox data={result} />
      </Card>
    </div>
  )
}
