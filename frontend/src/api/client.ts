const TOKEN_KEY = 'messenger_token';

function getBaseURL(): string {
  if (typeof window !== 'undefined' && window.location.host) {
    return window.location.protocol + '//' + window.location.host;
  }
  return '';
}

function getToken(): string {
  return localStorage.getItem(TOKEN_KEY) || '';
}

function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token);
}

function clearToken(): void {
  localStorage.removeItem(TOKEN_KEY);
}

function needsAuth(path: string): boolean {
  return !path.startsWith('/api/auth/register') &&
         !path.startsWith('/api/auth/login');
}

export async function apiCall<T = unknown>(
  method: string,
  path: string,
  body?: unknown,
  isFormData?: boolean
): Promise<ApiResponse<T>> {
  const fullURL = getBaseURL() + path;
  const token = getToken();

  if (needsAuth(path) && !token) {
    return { error: true, status: 401, data: 'Please login first' as T };
  }

  const opts: RequestInit = { method };
  const headers: Record<string, string> = {};

  if (token) headers['Authorization'] = 'Bearer ' + token;
  if (body && !isFormData) {
    headers['Content-Type'] = 'application/json';
    opts.body = JSON.stringify(body);
  } else if (body && isFormData) {
    opts.body = body as BodyInit;
  }

  if (Object.keys(headers).length > 0) opts.headers = headers;

  try {
    const res = await fetch(fullURL, opts);
    const ct = res.headers.get('content-type') || '';
    const textBody = await res.text();
    let data: T;

    if (ct.includes('application/json')) {
      try { data = JSON.parse(textBody) as T; } catch { data = textBody as T; }
    } else {
      data = textBody as T;
    }

    if (!res.ok) {
      let msg = 'Request failed';
      if (typeof data === 'object' && data !== null) {
        const d = data as Record<string, unknown>;
        msg = (d.error || d.message || d.detail || JSON.stringify(data)) as string;
      } else if (typeof data === 'string' && data) {
        msg = data;
      }
      return { error: true, status: res.status, data: msg as T };
    }

    return { error: false, status: res.status, data };
  } catch (e) {
    return { error: true, status: 0, data: (e as Error).message as T };
  }
}

export interface ApiResponse<T = unknown> {
  error: boolean;
  status: number;
  data: T;
}

export const api = {
  getToken,
  setToken,
  clearToken,
  call: apiCall,
};
