import { apiCall } from '../../api/client';

export async function generateCaptcha(): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('GET', '/api/security/captcha');
}

export async function verifyCaptcha(body: {token: string; solution: string}): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('POST', '/api/security/captcha/verify', body);
}

export async function registerE2EKey(body: {public_key: string}): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('POST', '/api/security/e2e/keys', body);
}

export async function getE2EPublicKey(userId: string): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('GET', `/api/security/e2e/keys/${userId}`);
}

export async function sendEmailVerification(email: string): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('POST', '/api/security/verify/email', { email });
}

export async function verifyEmail(code: string): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('POST', '/api/security/verify/email/confirm', { code });
}

export async function sendPhoneVerification(phone: string): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('POST', '/api/security/verify/phone', { phone });
}

export async function verifyPhone(code: string): Promise<{error: boolean; status: number; data: any}> {
  return apiCall<any>('POST', '/api/security/verify/phone/confirm', { code });
}
