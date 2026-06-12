import { apiCall } from '../../api/client'

export async function sendLocation(chatId: string, data: { latitude: number; longitude: number; title?: string; replyToId?: string; effect?: string }) {
  return apiCall('POST', `/api/chats/${chatId}/messages/location`, data)
}
