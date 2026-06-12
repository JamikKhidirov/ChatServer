import { useState } from 'react'
import { Card } from '../../components/common/Card'
import { FormRow, Input } from '../../components/common/FormRow'
import { LoadingButton } from '../../components/common/LoadingButton'
import { ResultBox } from '../../components/common/ResultBox'
import { useToast } from '../../components/common/Toast'
import * as api from './search'

export default function SearchPage() {
  const [result, setResult] = useState<unknown>(null)
  const { toast } = useToast()

  const handleAction = async (fn: () => Promise<any>, successMsg?: string) => {
    const res = await fn()
    setResult(res)
    if (!res.error && successMsg) toast(successMsg, 'success')
    return res
  }

  const [query, setQuery] = useState('')

  return (
    <div className="page-content">
      <h1 className="page-title">Search</h1>
      <p className="page-subtitle">Search for users</p>
      <Card title="Search Users" badge="GET">
        <FormRow>
          <Input placeholder="Search query" value={query} onChange={e => setQuery(e.target.value)} />
        </FormRow>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.searchUsers(query))}>Search</LoadingButton>
        <ResultBox data={result} />
      </Card>
    </div>
  )
}
