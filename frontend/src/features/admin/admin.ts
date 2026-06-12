import { apiCall } from '../../api/client';

export async function dashboard(): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('GET', '/api/admin/dashboard');
}

export async function listUsers(): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('GET', '/api/admin/users');
}

export async function listMessages(): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('GET', '/api/admin/messages');
}

export async function readMessage(id: string): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('GET', `/api/admin/messages/${id}`);
}

export async function banUser(body: {user_id: string; reason?: string}): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('POST', '/api/admin/ban', body);
}

export async function unbanUser(userId: string): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('POST', '/api/admin/unban', { user_id: userId });
}

export async function getSettings(): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('GET', '/api/admin/settings');
}

export async function updateSetting(body: {key: string; value: unknown}): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('PUT', '/api/admin/settings', body);
}

export async function getLogs(): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('GET', '/api/admin/logs');
}

export async function getIPBlocks(): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('GET', '/api/admin/ip-blocks');
}

export async function unblockIP(ip: string): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('DELETE', `/api/admin/ip-blocks/${encodeURIComponent(ip)}`);
}
