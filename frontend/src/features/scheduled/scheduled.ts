import { apiCall } from '../../api/client';

export async function scheduleMessage(body: {chat_id: string; content: string; scheduled_at: string}): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('POST', '/api/scheduled', body);
}

export async function getScheduled(): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('GET', '/api/scheduled');
}

export async function cancelScheduled(id: string): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('DELETE', `/api/scheduled/${id}`);
}
