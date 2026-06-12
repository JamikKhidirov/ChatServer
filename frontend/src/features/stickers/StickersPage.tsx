import { useState } from 'react'
import { Card } from '../../components/common/Card'
import { FormRow, Input, Select, TextArea, Checkbox } from '../../components/common/FormRow'
import { LoadingButton } from '../../components/common/LoadingButton'
import { ResultBox } from '../../components/common/ResultBox'
import { useToast } from '../../components/common/Toast'
import * as api from './stickers'

export default function StickersPage() {
  const [result, setResult] = useState<unknown>(null)
  const { toast } = useToast()

  const handleAction = async (fn: () => Promise<any>) => {
    const res = await fn()
    setResult(res)
    if (!res.error) toast('Success', 'success')
    return res
  }

  const [packName, setPackName] = useState('')
  const [packId, setPackId] = useState('')
  const [emoji, setEmoji] = useState('')
  const [imageUrl, setImageUrl] = useState('')
  const [deletePackId, setDeletePackId] = useState('')
  const [libraryStickerId, setLibraryStickerId] = useState('')

  return (
    <div className="page-content">
      <h1 className="page-title">Stickers</h1>
      <p className="page-subtitle">Create sticker packs, add stickers, and manage library</p>

      <Card title="Create Pack" badge="POST">
        <FormRow>
          <Input placeholder="Pack Name *" value={packName} onChange={e => setPackName(e.target.value)} />
        </FormRow>
        <LoadingButton variant="success" onClick={() => handleAction(() => api.createStickerPack({ name: packName }))}>Create Pack</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="List All Packs" badge="GET">
        <LoadingButton variant="info" onClick={() => handleAction(() => api.listStickerPacks())}>List All Packs</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="My Packs" badge="GET">
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getMyStickerPacks())}>My Packs</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Get Pack" badge="GET">
        <FormRow>
          <Input placeholder="Pack ID *" value={packId} onChange={e => setPackId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getStickerPack(packId))}>Get Pack</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Add Sticker" badge="POST">
        <FormRow>
          <Input placeholder="Pack ID *" value={packId} onChange={e => setPackId(e.target.value)} />
          <Input placeholder="Emoji *" value={emoji} onChange={e => setEmoji(e.target.value)} />
          <Input placeholder="Image URL" value={imageUrl} onChange={e => setImageUrl(e.target.value)} />
        </FormRow>
        <LoadingButton variant="success" onClick={() => handleAction(() => api.addSticker(packId, { emoji, imageUrl: imageUrl || undefined }))}>Add Sticker</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Delete Pack" badge="DELETE">
        <FormRow>
          <Input placeholder="Pack ID *" value={deletePackId} onChange={e => setDeletePackId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="danger" onClick={() => handleAction(() => api.deleteStickerPack(deletePackId))}>Delete Pack</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Sticker Library" badge="GET/POST">
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getStickerLibrary())}>Get Library</LoadingButton>
        <FormRow>
          <Input placeholder="Sticker ID" value={libraryStickerId} onChange={e => setLibraryStickerId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="success" onClick={() => handleAction(() => api.addStickerToLibrary(libraryStickerId))}>Add to Library</LoadingButton>
        <ResultBox data={result} />
      </Card>
    </div>
  )
}
