import { useState } from 'react'
import { Card } from '../../components/common/Card'
import { FormRow, Input, Select, TextArea, Checkbox } from '../../components/common/FormRow'
import { LoadingButton } from '../../components/common/LoadingButton'
import { ResultBox } from '../../components/common/ResultBox'
import { useToast } from '../../components/common/Toast'
import * as api from './emojis'

export default function EmojisPage() {
  const [result, setResult] = useState<unknown>(null)
  const { toast } = useToast()

  const handleAction = async (fn: () => Promise<any>) => {
    const res = await fn()
    setResult(res)
    if (!res.error) toast('Success', 'success')
    return res
  }

  const [shortcode, setShortcode] = useState('')
  const [deleteId, setDeleteId] = useState('')

  return (
    <div className="page-content">
      <h1 className="page-title">Custom Emojis</h1>
      <p className="page-subtitle">Upload, view, and delete custom emojis</p>

      <Card title="Upload Emoji" badge="POST">
        <FormRow>
          <Input placeholder="Shortcode *" value={shortcode} onChange={e => setShortcode(e.target.value)} />
          <input type="file" accept="image/*" style={{ color: '#9a9ab8', fontSize: 13, flex: 1 }} id="emojiFile" />
        </FormRow>
        <LoadingButton variant="success" onClick={() => {
          const fileInput = document.getElementById('emojiFile') as HTMLInputElement
          return handleAction(() => api.createEmoji(shortcode, fileInput?.files?.[0]!))
        }}>Upload</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="My Emojis" badge="GET">
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getMyEmojis())}>My Emojis</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="All Emojis" badge="GET">
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getAllEmojis())}>All Emojis</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Delete Emoji" badge="DELETE">
        <FormRow>
          <Input placeholder="Emoji ID *" value={deleteId} onChange={e => setDeleteId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="danger" onClick={() => handleAction(() => api.deleteEmoji(deleteId))}>Delete</LoadingButton>
        <ResultBox data={result} />
      </Card>
    </div>
  )
}
