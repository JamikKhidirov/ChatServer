import { useState } from 'react'
import { Card } from '../../components/common/Card'
import { FormRow, Input } from '../../components/common/FormRow'
import { LoadingButton } from '../../components/common/LoadingButton'
import { ResultBox } from '../../components/common/ResultBox'
import { useToast } from '../../components/common/Toast'
import * as api from './bots'

export default function BotsPage() {
  const [result, setResult] = useState<unknown>(null)
  const { toast } = useToast()

  const handleAction = async (fn: () => Promise<any>, successMsg?: string) => {
    const res = await fn()
    setResult(res)
    if (!res.error && successMsg) toast(successMsg, 'success')
    return res
  }

  const [name, setName] = useState('')
  const [webhookUrl, setWebhookUrl] = useState('')
  const [updateBotId, setUpdateBotId] = useState('')
  const [updateName, setUpdateName] = useState('')
  const [updateWebhookUrl, setUpdateWebhookUrl] = useState('')
  const [deleteBotId, setDeleteBotId] = useState('')
  const [regenerateBotId, setRegenerateBotId] = useState('')

  return (
    <div className="page-content">
      <h1 className="page-title">Bots</h1>
      <p className="page-subtitle">Create and manage your chat bots</p>
      <Card title="Create Bot" badge="POST">
        <FormRow>
          <Input placeholder="Bot name" value={name} onChange={e => setName(e.target.value)} />
          <Input placeholder="Webhook URL" value={webhookUrl} onChange={e => setWebhookUrl(e.target.value)} />
        </FormRow>
        <LoadingButton variant="success" onClick={() => handleAction(() => api.createBot({ name }), 'Bot created')}>Create Bot</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="Get My Bots" badge="GET">
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getMyBots())}>Get My Bots</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="Update Bot" badge="PUT">
        <FormRow>
          <Input placeholder="Bot ID" value={updateBotId} onChange={e => setUpdateBotId(e.target.value)} />
        </FormRow>
        <FormRow>
          <Input placeholder="New name" value={updateName} onChange={e => setUpdateName(e.target.value)} />
          <Input placeholder="New webhook URL" value={updateWebhookUrl} onChange={e => setUpdateWebhookUrl(e.target.value)} />
        </FormRow>
        <LoadingButton variant="warning" onClick={() => handleAction(() => api.updateBot(updateBotId, { name: updateName }), 'Bot updated')}>Update Bot</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="Delete Bot" badge="DELETE">
        <FormRow>
          <Input placeholder="Bot ID" value={deleteBotId} onChange={e => setDeleteBotId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="danger" onClick={() => handleAction(() => api.deleteBot(deleteBotId), 'Bot deleted')}>Delete Bot</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="Regenerate Token" badge="POST">
        <FormRow>
          <Input placeholder="Bot ID" value={regenerateBotId} onChange={e => setRegenerateBotId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="warning" onClick={() => handleAction(() => api.regenerateBotToken(regenerateBotId), 'Token regenerated')}>Regenerate Token</LoadingButton>
        <ResultBox data={result} />
      </Card>
    </div>
  )
}
