import { useState, useEffect, useCallback, createContext, useContext } from 'react';

interface Toast {
  id: number;
  text: string;
  type: 'success' | 'error' | 'info';
}

interface ToastContextType {
  toast: (text: string, type?: 'success' | 'error' | 'info') => void;
}

const ToastContext = createContext<ToastContextType>({ toast: () => {} });

export const useToast = () => useContext(ToastContext);

let nextId = 0;

export function ToastProvider({ children }: { children: React.ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([]);

  const addToast = useCallback((text: string, type: 'success' | 'error' | 'info' = 'info') => {
    const id = nextId++;
    setToasts(prev => [...prev, { id, text, type }]);
    setTimeout(() => {
      setToasts(prev => prev.filter(t => t.id !== id));
    }, 3500);
  }, []);

  const icons: Record<string, string> = { error: '⚠️', success: '✅', info: '💬' };

  return (
    <ToastContext.Provider value={{ toast: addToast }}>
      {children}
      <div style={{
        position: 'fixed', top: 16, right: 16, zIndex: 9999,
        display: 'flex', flexDirection: 'column', gap: 8, maxWidth: 400,
        pointerEvents: 'none'
      }}>
        {toasts.map(t => (
          <div key={t.id} style={{
            pointerEvents: 'auto',
            background: 'rgba(22,33,62,0.95)',
            border: '1px solid rgba(255,255,255,0.06)',
            borderLeft: `4px solid ${t.type === 'error' ? '#e74c3c' : t.type === 'success' ? '#2ecc71' : '#f39c12'}`,
            borderRadius: 12,
            padding: '14px 22px',
            color: '#e8e8f0',
            fontSize: 13,
            boxShadow: '0 12px 48px rgba(0,0,0,0.6)',
            backdropFilter: 'blur(20px)',
            display: 'flex',
            alignItems: 'flex-start',
            gap: 10,
            animation: 'toastIn 0.35s ease',
          }}>
            <span style={{fontSize: 16}}>{icons[t.type]}</span>
            <span>{t.text}</span>
          </div>
        ))}
      </div>
    </ToastContext.Provider>
  );
}
