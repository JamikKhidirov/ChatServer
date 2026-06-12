import { apiCall } from '../../api/client';

export async function createReport(body: {target_id: string; target_type: string; reason: string}): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('POST', '/api/reports', body);
}

export async function listReports(status?: string): Promise<{error: boolean; status: number; data: any}> {
  let path = '/api/reports';
  if (status) path += `?status=${encodeURIComponent(status)}`;
  return apiCall<any>('GET', path);
}

export async function resolveReport(id: string, status: string): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('PUT', `/api/reports/${id}`, { status });
}
