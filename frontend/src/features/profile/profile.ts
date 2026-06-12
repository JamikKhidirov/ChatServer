import { apiCall } from '../../api/client'

export async function getProfile() {
  return apiCall('GET', '/api/users/profile')
}

export async function updateProfile(data: { displayName?: string; bio?: string; phone?: string; gender?: string; dateOfBirth?: string }) {
  return apiCall('PUT', '/api/users/profile', data)
}

export async function uploadAvatar(file: File) {
  const fd = new FormData()
  fd.append('avatar', file)
  return apiCall('POST', '/api/users/avatar', fd, true)
}

export async function updateStatus(data: { status: string }) {
  return apiCall('PUT', '/api/users/status', data)
}

export async function savePushToken(data: { token: string; provider: string }) {
  return apiCall('PUT', '/api/users/push-token', data)
}

export async function testPush() {
  return apiCall('POST', '/api/users/push-test', { title: 'Test', body: 'Hello from API Tester' })
}

export async function deleteAccount() {
  return apiCall('DELETE', '/api/users/account')
}
