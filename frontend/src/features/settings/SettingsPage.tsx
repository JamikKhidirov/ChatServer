import { useState } from 'react'
import { Card } from '../../components/common/Card'
import { FormRow, Input, Select, TextArea, Checkbox } from '../../components/common/FormRow'
import { LoadingButton } from '../../components/common/LoadingButton'
import { ResultBox } from '../../components/common/ResultBox'
import { useToast } from '../../components/common/Toast'
import * as api from './settings'

export default function SettingsPage() {
  const [result, setResult] = useState<unknown>(null)
  const { toast } = useToast()

  const handleAction = async (fn: () => Promise<any>) => {
    const res = await fn()
    setResult(res)
    if (!res.error) toast('Success', 'success')
    return res
  }

  const [language, setLanguage] = useState('')
  const [theme, setTheme] = useState('')
  const [notifications, setNotifications] = useState(true)
  const [soundEnabled, setSoundEnabled] = useState(true)
  const [lastSeenMode, setLastSeenMode] = useState('')

  return (
    <div className="page-content">
      <h1 className="page-title">Settings</h1>
      <p className="page-subtitle">View and update your account settings</p>

      <Card title="View Settings" badge="GET">
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getSettings())}>View Settings</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Update Settings" badge="PUT">
        <FormRow>
          <Input placeholder="Language (e.g. en, ru)" value={language} onChange={e => setLanguage(e.target.value)} />
          <Input placeholder="Theme (e.g. dark, light)" value={theme} onChange={e => setTheme(e.target.value)} />
        </FormRow>
        <FormRow>
          <Select value={lastSeenMode} onChange={e => setLastSeenMode(e.target.value)}>
            <option value="">Last Seen Mode</option>
            <option value="everyone">Everyone</option>
            <option value="contacts">Contacts</option>
            <option value="nobody">Nobody</option>
          </Select>
        </FormRow>
        <FormRow>
          <Checkbox label="Notifications" checked={notifications} onChange={e => setNotifications(e.target.checked)} />
          <Checkbox label="Sound Enabled" checked={soundEnabled} onChange={e => setSoundEnabled(e.target.checked)} />
        </FormRow>
        <LoadingButton variant="warning" onClick={() => handleAction(() => api.updateSettings({
          language: language || undefined,
          theme: theme || undefined,
          notifications,
          soundEnabled,
          lastSeenMode: lastSeenMode || undefined,
        }))}>Update Settings</LoadingButton>
        <ResultBox data={result} />
      </Card>
    </div>
  )
}
