import { apiCall } from '../../api/client'

export async function saveMessage(messageId: string, chatId: string) {
  return apiCall('POST', `/api/messages/${messageId}/save?chatId=${encodeURIComponent(chatId)}`)
}

export async function getSavedMessages(limit?: number, offset?: number) {
  let path = '/api/saved-messages'
  const params: string[] = []
  if (limit !== undefined) params.push(`limit=${limit}`)
  if (offset !== undefined) params.push(`offset=${offset}`)
  if (params.length) path += '?' + params.join('&')
  return apiCall('GET', path)
}

export async function deleteSavedMessage(id: string) {
  return apiCall('DELETE', `/api/saved-messages/${id}`)
}
