import { apiCall } from '../../api/client'

export async function getSettings() {
  return apiCall('GET', '/api/account/settings')
}

export async function updateSettings(data: { language?: string; theme?: string; notifications?: boolean; soundEnabled?: boolean; lastSeenMode?: string }) {
  return apiCall('PUT', '/api/account/settings', data)
}
