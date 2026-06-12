import { apiCall } from '../../api/client';

export async function blockUser(body: {user_id: string}): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('POST', '/api/block', body);
}

export async function unblockUser(userId: string): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('DELETE', `/api/block/${userId}`);
}

export async function listBlocked(): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('GET', '/api/block');
}
