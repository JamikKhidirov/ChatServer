import { apiCall } from '../../api/client';

export async function createBot(body: {name: string; description?: string; avatar?: string}): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('POST', '/api/bots', body);
}

export async function getMyBots(): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('GET', '/api/bots/my');
}

export async function updateBot(id: string, body: {name?: string; description?: string; avatar?: string}): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('PUT', `/api/bots/${id}`, body);
}

export async function deleteBot(id: string): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('DELETE', `/api/bots/${id}`);
}

export async function regenerateBotToken(id: string): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('POST', `/api/bots/${id}/regenerate-token`);
}
