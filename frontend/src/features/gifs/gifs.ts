import { apiCall } from '../../api/client';

export async function saveGif(url: string): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('POST', '/api/gifs', { url });
}

export async function getSavedGifs(): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('GET', '/api/gifs');
}

export async function deleteGif(url: string): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('DELETE', `/api/gifs?url=${encodeURIComponent(url)}`);
}
