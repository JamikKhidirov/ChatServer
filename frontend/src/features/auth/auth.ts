import { apiCall, api } from '../../api/client'

export function getToken() {
  return api.getToken()
}

function storeToken(token: string) {
  api.setToken(token)
}

export async function register(data: { username: string; email: string; password: string; displayName: string }) {
  return apiCall('POST', '/api/auth/register', data)
}

export async function login(data: { email: string; password: string }) {
  const res = await apiCall('POST', '/api/auth/login', data)
  if (!res.error && res.data && typeof res.data === 'object' && 'token' in (res.data as any)) {
    storeToken((res.data as any).token)
  }
  return res
}

export async function refreshToken() {
  return apiCall('GET', '/api/auth/refresh')
}

export async function changePassword(oldPassword: string, newPassword: string) {
  return apiCall('PUT', '/api/auth/change-password', { oldPassword, newPassword })
}

export async function sendEmailLoginCode(email: string) {
  return apiCall('POST', '/api/auth/login/email', { email })
}

export async function verifyEmailLoginCode(email: string, code: string) {
  const res = await apiCall('POST', '/api/auth/login/email/verify', { email, code })
  if (!res.error && res.data && typeof res.data === 'object' && 'token' in (res.data as any)) {
    storeToken((res.data as any).token)
  }
  return res
}

export async function sendPhoneLoginCode(phone: string) {
  return apiCall('POST', '/api/auth/login/phone', { phone })
}

export async function verifyPhoneLoginCode(phone: string, code: string) {
  const res = await apiCall('POST', '/api/auth/login/phone/verify', { phone, code })
  if (!res.error && res.data && typeof res.data === 'object' && 'token' in (res.data as any)) {
    storeToken((res.data as any).token)
  }
  return res
}
