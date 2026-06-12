import { apiCall } from '../../api/client';

export async function getLinkPreview(url: string): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('GET', `/api/preview?url=${encodeURIComponent(url)}`);
}
