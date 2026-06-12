import { useState } from 'react'
import { Card } from '../../components/common/Card'
import { FormRow, Input } from '../../components/common/FormRow'
import { LoadingButton } from '../../components/common/LoadingButton'
import { ResultBox } from '../../components/common/ResultBox'
import { useToast } from '../../components/common/Toast'
import * as api from './bookmarks'

export default function BookmarksPage() {
  const [result, setResult] = useState<unknown>(null)
  const { toast } = useToast()

  const handleAction = async (fn: () => Promise<any>, successMsg?: string) => {
    const res = await fn()
    setResult(res)
    if (!res.error && successMsg) toast(successMsg, 'success')
    return res
  }

  const [messageId, setMessageId] = useState('')
  const [chatId, setChatId] = useState('')
  const [removeMessageId, setRemoveMessageId] = useState('')

  return (
    <div className="page-content">
      <h1 className="page-title">Bookmarks</h1>
      <p className="page-subtitle">Bookmark and manage your favorite messages</p>
      <Card title="Bookmark Message" badge="POST">
        <FormRow>
          <Input placeholder="Message ID" value={messageId} onChange={e => setMessageId(e.target.value)} />
          <Input placeholder="Chat ID" value={chatId} onChange={e => setChatId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="success" onClick={() => handleAction(() => api.bookmarkMessage({ message_id: messageId, chat_id: chatId }), 'Message bookmarked')}>Bookmark</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="Get Bookmarks" badge="GET">
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getBookmarks())}>Get Bookmarks</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="Remove Bookmark" badge="DELETE">
        <FormRow>
          <Input placeholder="Message ID" value={removeMessageId} onChange={e => setRemoveMessageId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="danger" onClick={() => handleAction(() => api.removeBookmark(removeMessageId), 'Bookmark removed')}>Remove Bookmark</LoadingButton>
        <ResultBox data={result} />
      </Card>
    </div>
  )
}
