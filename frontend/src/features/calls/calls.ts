import { apiCall } from '../../api/client';

export async function initiateCall(body: {chat_id: string; type: 'audio' | 'video'}): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('POST', '/api/calls', body);
}

export async function respondCall(callId: string, action: 'accept' | 'reject' | 'ignore'): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('POST', `/api/calls/${callId}/respond`, { action });
}

export async function endCall(callId: string): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('POST', `/api/calls/${callId}/end`);
}

export async function getCall(id: string): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('GET', `/api/calls/${id}`);
}

export async function getCallHistory(chatId: string): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('GET', `/api/chats/${chatId}/calls`);
}
