import { apiCall } from '../../api/client';

export async function getSessions(): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('GET', '/api/sessions');
}

export async function deleteSession(id: string): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('DELETE', `/api/sessions/${id}`);
}

export async function deleteAllSessions(): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('DELETE', '/api/sessions');
}
