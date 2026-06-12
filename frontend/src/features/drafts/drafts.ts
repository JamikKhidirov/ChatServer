import { apiCall } from '../../api/client'

export async function saveDraft(data: { chatId: string; content: string }) {
  return apiCall('POST', '/api/drafts', data)
}

export async function getDraft(chatId: string) {
  return apiCall('GET', `/api/drafts?chatId=${encodeURIComponent(chatId)}`)
}
