import { useState } from 'react'
import { Card } from '../../components/common/Card'
import { FormRow, Input } from '../../components/common/FormRow'
import { LoadingButton } from '../../components/common/LoadingButton'
import { ResultBox } from '../../components/common/ResultBox'
import { useToast } from '../../components/common/Toast'
import * as api from './gifs'

export default function GifsPage() {
  const [result, setResult] = useState<unknown>(null)
  const { toast } = useToast()

  const handleAction = async (fn: () => Promise<any>, successMsg?: string) => {
    const res = await fn()
    setResult(res)
    if (!res.error && successMsg) toast(successMsg, 'success')
    return res
  }

  const [gifUrl, setGifUrl] = useState('')
  const [deleteUrl, setDeleteUrl] = useState('')

  return (
    <div className="page-content">
      <h1 className="page-title">GIFs</h1>
      <p className="page-subtitle">Manage your saved GIFs</p>
      <Card title="Save GIF" badge="POST">
        <FormRow>
          <Input placeholder="GIF URL" value={gifUrl} onChange={e => setGifUrl(e.target.value)} />
        </FormRow>
        <LoadingButton variant="success" onClick={() => handleAction(() => api.saveGif(gifUrl), 'GIF saved')}>Save GIF</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="Get Saved GIFs" badge="GET">
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getSavedGifs())}>Get Saved GIFs</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="Delete GIF" badge="DELETE">
        <FormRow>
          <Input placeholder="GIF URL" value={deleteUrl} onChange={e => setDeleteUrl(e.target.value)} />
        </FormRow>
        <LoadingButton variant="danger" onClick={() => handleAction(() => api.deleteGif(deleteUrl), 'GIF deleted')}>Delete GIF</LoadingButton>
        <ResultBox data={result} />
      </Card>
    </div>
  )
}
