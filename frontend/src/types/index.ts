export interface ApiResponse<T = unknown> {
  error: boolean;
  status: number;
  data: T;
}

export interface ToastMessage {
  id: number;
  text: string;
  type: 'success' | 'error' | 'info';
}

export interface TabConfig {
  id: string;
  label: string;
  icon: string;
}
