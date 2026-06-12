import { apiCall } from '../../api/client'

export async function createEmoji(shortcode: string, file: File) {
  const fd = new FormData()
  fd.append('shortcode', shortcode)
  fd.append('emoji', file)
  return apiCall('POST', '/api/emojis', fd, true)
}

export async function getMyEmojis() {
  return apiCall('GET', '/api/emojis/my')
}

export async function getAllEmojis() {
  return apiCall('GET', '/api/emojis')
}

export async function deleteEmoji(id: string) {
  return apiCall('DELETE', `/api/emojis/${id}`)
}
