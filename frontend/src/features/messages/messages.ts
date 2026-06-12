import { apiCall } from '../../api/client'

export async function listMessages(chatId: string, limit?: number, offset?: number) {
  let path = `/api/chats/${chatId}/messages`
  const params: string[] = []
  if (limit !== undefined) params.push(`limit=${limit}`)
  if (offset !== undefined) params.push(`offset=${offset}`)
  if (params.length) path += '?' + params.join('&')
  return apiCall('GET', path)
}

export async function sendMessage(chatId: string, body: { content: string; type?: string; replyToId?: string; effect?: string }) {
  return apiCall('POST', `/api/chats/${chatId}/messages`, body)
}

export async function uploadFile(chatId: string, file: File) {
  const fd = new FormData()
  fd.append('file', file)
  return apiCall('POST', `/api/chats/${chatId}/messages/file`, fd, true)
}

export async function uploadVideoCircle(chatId: string, file: File, caption?: string) {
  const fd = new FormData()
  fd.append('video', file)
  if (caption) fd.append('caption', caption)
  return apiCall('POST', `/api/chats/${chatId}/messages/video-circle`, fd, true)
}

export async function editMessage(messageId: string, content: string) {
  return apiCall('PUT', `/api/messages/${messageId}`, { content })
}

export async function deleteMessage(messageId: string) {
  return apiCall('DELETE', `/api/messages/${messageId}`)
}

export async function addReaction(messageId: string, emoji: string) {
  return apiCall('POST', `/api/messages/${messageId}/reactions`, { emoji })
}

export async function removeReaction(messageId: string, emoji: string) {
  return apiCall('DELETE', `/api/messages/${messageId}/reactions?emoji=${encodeURIComponent(emoji)}`)
}

export async function togglePin(messageId: string, pin: boolean) {
  return apiCall('PUT', `/api/messages/${messageId}/pin`, { pin })
}

export async function forwardMessage(data: { messageId: string; fromChatId: string; toChatId: string }) {
  return apiCall('POST', '/api/messages/forward', data)
}

export async function starMessage(messageId: string) {
  return apiCall('POST', `/api/messages/${messageId}/star`)
}

export async function unstarMessage(messageId: string) {
  return apiCall('DELETE', `/api/messages/${messageId}/star`)
}

export async function getStarredMessages() {
  return apiCall('GET', '/api/messages/starred')
}

export async function selfDestruct(messageId: string, seconds: number) {
  return apiCall('POST', `/api/messages/${messageId}/self-destruct`, { seconds })
}

export async function getMessageHistory(messageId: string) {
  return apiCall('GET', `/api/messages/${messageId}/history`)
}

export async function getChatMedia(chatId: string, type?: string) {
  let path = `/api/chats/${chatId}/media`
  if (type) path += `?type=${encodeURIComponent(type)}`
  return apiCall('GET', path)
}

export async function exportChat(chatId: string) {
  return apiCall('GET', `/api/chats/${chatId}/export`)
}

export async function deleteForMe(messageId: string) {
  return apiCall('DELETE', `/api/messages/${messageId}/for-me`)
}

export async function searchAllMessages(query: string) {
  return apiCall('GET', `/api/messages/search?q=${encodeURIComponent(query)}`)
}
