import { apiCall } from '../../api/client';

export async function bookmarkMessage(body: {message_id: string; chat_id: string}): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('POST', '/api/bookmarks', body);
}

export async function getBookmarks(): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('GET', '/api/bookmarks');
}

export async function removeBookmark(msgId: string): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('DELETE', `/api/bookmarks/${msgId}`);
}
