import { apiCall } from '../../api/client'

export async function syncContacts(data: { contacts: Array<{ phone: string; name: string }> }) {
  return apiCall('POST', '/api/contacts/sync', data)
}

export async function getContacts() {
  return apiCall('GET', '/api/contacts')
}
