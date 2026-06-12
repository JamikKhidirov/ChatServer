import { useState } from 'react'
import { Card } from '../../components/common/Card'
import { FormRow, Input, Select, TextArea, Checkbox } from '../../components/common/FormRow'
import { LoadingButton } from '../../components/common/LoadingButton'
import { ResultBox } from '../../components/common/ResultBox'
import { useToast } from '../../components/common/Toast'
import * as api from './contacts'

export default function ContactsPage() {
  const [result, setResult] = useState<unknown>(null)
  const { toast } = useToast()

  const handleAction = async (fn: () => Promise<any>) => {
    const res = await fn()
    setResult(res)
    if (!res.error) toast('Success', 'success')
    return res
  }

  const [contactsJson, setContactsJson] = useState('')

  return (
    <div className="page-content">
      <h1 className="page-title">Contacts</h1>
      <p className="page-subtitle">Sync and manage your phone contacts</p>

      <Card title="Sync Contacts" badge="POST">
        <FormRow>
          <TextArea placeholder='[{"phone":"+1234567890","name":"John Doe"}]' value={contactsJson} onChange={e => setContactsJson(e.target.value)} />
        </FormRow>
        <LoadingButton variant="success" onClick={() => {
          let data
          try { data = JSON.parse(contactsJson) } catch { data = { contacts: [] } }
          return handleAction(() => api.syncContacts({ contacts: Array.isArray(data) ? data : data.contacts || [] }))
        }}>Sync Contacts</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="View Contacts" badge="GET">
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getContacts())}>View Contacts</LoadingButton>
        <ResultBox data={result} />
      </Card>
    </div>
  )
}
