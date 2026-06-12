import { apiCall } from '../../api/client'

export async function listChats() {
  return apiCall('GET', '/api/chats')
}

export async function createChat(data: { type: string; name?: string; participantIds: string[] }) {
  return apiCall('POST', '/api/chats', data)
}

export async function getChat(id: string) {
  return apiCall('GET', `/api/chats/${id}`)
}

export async function hideChat(id: string) {
  return apiCall('POST', `/api/chats/${id}/hide`)
}

export async function deleteChat(id: string) {
  return apiCall('DELETE', `/api/chats/${id}`)
}

export async function leaveChat(id: string) {
  return apiCall('POST', `/api/chats/${id}/leave`)
}

export async function markRead(id: string) {
  return apiCall('POST', `/api/chats/${id}/read`)
}

export async function addParticipant(chatId: string, userId: string) {
  return apiCall('POST', `/api/chats/${chatId}/participants`, { userId })
}

export async function removeParticipant(chatId: string, userId: string) {
  return apiCall('DELETE', `/api/chats/${chatId}/participants/${userId}`)
}

export async function setRole(chatId: string, userId: string, role: string) {
  return apiCall('PUT', `/api/chats/${chatId}/participants/${userId}/role`, { role })
}

export async function setNotificationMuted(chatId: string, muted: boolean) {
  return apiCall('PUT', `/api/chats/${chatId}/notifications`, { muted })
}

export async function pinChat(id: string) {
  return apiCall('POST', `/api/chats/${id}/pin`)
}

export async function unpinChat(id: string) {
  return apiCall('DELETE', `/api/chats/${id}/pin`)
}

export async function archiveChat(id: string) {
  return apiCall('POST', `/api/chats/${id}/archive`)
}

export async function unarchiveChat(id: string) {
  return apiCall('POST', `/api/chats/${id}/unarchive`)
}
