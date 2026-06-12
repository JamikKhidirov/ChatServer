import { useState } from 'react'
import { Card } from '../../components/common/Card'
import { FormRow, Input } from '../../components/common/FormRow'
import { LoadingButton } from '../../components/common/LoadingButton'
import { ResultBox } from '../../components/common/ResultBox'
import { useToast } from '../../components/common/Toast'
import * as api from './preview'

export default function PreviewPage() {
  const [result, setResult] = useState<unknown>(null)
  const { toast } = useToast()

  const handleAction = async (fn: () => Promise<any>, successMsg?: string) => {
    const res = await fn()
    setResult(res)
    if (!res.error && successMsg) toast(successMsg, 'success')
    return res
  }

  const [url, setUrl] = useState('')

  return (
    <div className="page-content">
      <h1 className="page-title">Link Preview</h1>
      <p className="page-subtitle">Fetch metadata preview for a URL</p>
      <Card title="Get Link Preview" badge="GET">
        <FormRow>
          <Input placeholder="URL to preview" value={url} onChange={e => setUrl(e.target.value)} />
        </FormRow>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getLinkPreview(url))}>Get Preview</LoadingButton>
        <ResultBox data={result} />
      </Card>
    </div>
  )
}
