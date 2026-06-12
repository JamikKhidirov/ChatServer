import { apiCall } from '../../api/client'

export async function createPoll(data: { chatId: string; question: string; options: string[]; isAnonymous?: boolean; multipleChoice?: boolean }) {
  return apiCall('POST', `/api/chats/${data.chatId}/polls`, data)
}

export async function votePoll(pollId: string, optionIndex: number) {
  return apiCall('POST', `/api/polls/${pollId}/vote`, { optionIndex })
}

export async function closePoll(pollId: string) {
  return apiCall('POST', `/api/polls/${pollId}/close`)
}
