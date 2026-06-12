import { useState } from 'react'
import { Card } from '../../components/common/Card'
import { FormRow, Input, Select, TextArea, Checkbox } from '../../components/common/FormRow'
import { LoadingButton } from '../../components/common/LoadingButton'
import { ResultBox } from '../../components/common/ResultBox'
import { useToast } from '../../components/common/Toast'
import * as api from './location'

export default function LocationPage() {
  const [result, setResult] = useState<unknown>(null)
  const { toast } = useToast()

  const handleAction = async (fn: () => Promise<any>) => {
    const res = await fn()
    setResult(res)
    if (!res.error) toast('Success', 'success')
    return res
  }

  const [chatId, setChatId] = useState('')
  const [latitude, setLatitude] = useState('')
  const [longitude, setLongitude] = useState('')
  const [title, setTitle] = useState('')
  const [replyToId, setReplyToId] = useState('')
  const [effect, setEffect] = useState('')

  return (
    <div className="page-content">
      <h1 className="page-title">Location</h1>
      <p className="page-subtitle">Send location messages to chats</p>

      <Card title="Send Location" badge="POST">
        <FormRow>
          <Input placeholder="Chat ID *" value={chatId} onChange={e => setChatId(e.target.value)} />
          <Input placeholder="Latitude *" type="number" value={latitude} onChange={e => setLatitude(e.target.value)} />
          <Input placeholder="Longitude *" type="number" value={longitude} onChange={e => setLongitude(e.target.value)} />
        </FormRow>
        <FormRow>
          <Input placeholder="Title" value={title} onChange={e => setTitle(e.target.value)} />
          <Input placeholder="Reply To ID" value={replyToId} onChange={e => setReplyToId(e.target.value)} />
          <Select value={effect} onChange={e => setEffect(e.target.value)}>
            <option value="">No Effect</option>
            <option value="confetti">Confetti</option>
            <option value="fireworks">Fireworks</option>
            <option value="hearts">Hearts</option>
            <option value="balloons">Balloons</option>
            <option value="stars">Stars</option>
          </Select>
        </FormRow>
        <LoadingButton variant="success" onClick={() => handleAction(() => api.sendLocation(chatId, {
          latitude: Number(latitude),
          longitude: Number(longitude),
          title: title || undefined,
          replyToId: replyToId || undefined,
          effect: effect || undefined,
        }))}>Send Location</LoadingButton>
        <ResultBox data={result} />
      </Card>
    </div>
  )
}
