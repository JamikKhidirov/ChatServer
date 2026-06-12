import { useState } from 'react'
import { Card } from '../../components/common/Card'
import { FormRow, Input, TextArea } from '../../components/common/FormRow'
import { LoadingButton } from '../../components/common/LoadingButton'
import { ResultBox } from '../../components/common/ResultBox'
import { useToast } from '../../components/common/Toast'
import * as api from './admin'

export default function AdminPage() {
  const [result, setResult] = useState<unknown>(null)
  const { toast } = useToast()

  const handleAction = async (fn: () => Promise<any>, successMsg?: string) => {
    const res = await fn()
    setResult(res)
    if (!res.error && successMsg) toast(successMsg, 'success')
    return res
  }

  const [banUserId, setBanUserId] = useState('')
  const [banReason, setBanReason] = useState('')
  const [unbanUserId, setUnbanUserId] = useState('')
  const [settingKey, setSettingKey] = useState('')
  const [settingValue, setSettingValue] = useState('')

  return (
    <div className="page-content">
      <h1 className="page-title">Admin</h1>
      <p className="page-subtitle">Administrative dashboard and controls</p>
      <Card title="Dashboard Stats" badge="GET">
        <LoadingButton variant="info" onClick={() => handleAction(() => api.dashboard())}>Dashboard Stats</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="List Users" badge="GET">
        <LoadingButton variant="info" onClick={() => handleAction(() => api.listUsers())}>List Users</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="List Messages" badge="GET">
        <LoadingButton variant="info" onClick={() => handleAction(() => api.listMessages())}>List Messages</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="Ban User" badge="POST">
        <FormRow>
          <Input placeholder="User ID" value={banUserId} onChange={e => setBanUserId(e.target.value)} />
          <Input placeholder="Reason" value={banReason} onChange={e => setBanReason(e.target.value)} />
        </FormRow>
        <LoadingButton variant="danger" onClick={() => handleAction(() => api.banUser({ user_id: banUserId, reason: banReason }), 'User banned')}>Ban User</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="Unban User" badge="POST">
        <FormRow>
          <Input placeholder="User ID" value={unbanUserId} onChange={e => setUnbanUserId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="warning" onClick={() => handleAction(() => api.unbanUser(unbanUserId), 'User unbanned')}>Unban User</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="Get Settings" badge="GET">
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getSettings())}>Get Settings</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="Update Setting" badge="PUT">
        <FormRow>
          <Input placeholder="Setting key" value={settingKey} onChange={e => setSettingKey(e.target.value)} />
          <Input placeholder="Setting value" value={settingValue} onChange={e => setSettingValue(e.target.value)} />
        </FormRow>
        <LoadingButton variant="warning" onClick={() => handleAction(() => api.updateSetting({ key: settingKey, value: settingValue }), 'Setting updated')}>Update Setting</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="Get Logs" badge="GET">
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getLogs())}>Get Logs</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="Get IP Blocks" badge="GET">
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getIPBlocks())}>Get IP Blocks</LoadingButton>
        <ResultBox data={result} />
      </Card>
    </div>
  )
}
