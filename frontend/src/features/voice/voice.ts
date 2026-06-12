import { apiCall } from '../../api/client'

export async function createVoiceChat(chatId: string, data: { title?: string; scheduledInMins?: number }) {
  return apiCall('POST', `/api/chats/${chatId}/voice-chat`, data)
}

export async function getActiveVoiceChats(chatId: string) {
  return apiCall('GET', `/api/chats/${chatId}/voice-chats/active`)
}

export async function getVoiceChatHistory(chatId: string) {
  return apiCall('GET', `/api/chats/${chatId}/voice-chats/history`)
}

export async function getVoiceChat(id: string) {
  return apiCall('GET', `/api/voice-chats/${id}`)
}

export async function joinVoiceChat(id: string) {
  return apiCall('POST', `/api/voice-chats/${id}/join`)
}

export async function leaveVoiceChat(id: string) {
  return apiCall('POST', `/api/voice-chats/${id}/leave`)
}

export async function endVoiceChat(id: string) {
  return apiCall('POST', `/api/voice-chats/${id}/end`)
}

export async function muteParticipant(id: string, muted: boolean) {
  return apiCall('POST', `/api/voice-chats/${id}/mute`, { muted })
}
