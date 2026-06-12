import { apiCall } from '../../api/client';

export async function searchUsers(query: string): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('GET', `/api/search/users?q=${encodeURIComponent(query)}`);
}
