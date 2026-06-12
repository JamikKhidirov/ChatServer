import { useState } from 'react';

interface LoadingButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'default' | 'success' | 'danger' | 'info' | 'warning';
  glass?: boolean;
}

const variantStyles: Record<string, React.CSSProperties> = {
  success: { borderColor: 'rgba(46,204,113,0.4)' },
  danger: { borderColor: 'rgba(231,76,60,0.4)' },
  info: { borderColor: 'rgba(52,152,219,0.4)' },
  warning: { borderColor: 'rgba(243,156,18,0.4)' },
};

export function LoadingButton({ children, variant, glass, style, onClick, disabled, ...props }: LoadingButtonProps) {
  const [loading, setLoading] = useState(false);

  const handleClick = async (e: React.MouseEvent<HTMLButtonElement>) => {
    if (loading || disabled) return;
    setLoading(true);
    try {
      await onClick?.(e);
    } finally {
      setLoading(false);
    }
  };

  return (
    <button
      {...props}
      onClick={handleClick}
      disabled={disabled || loading}
      style={{
        display: 'inline-flex', alignItems: 'center', gap: 6,
        background: glass ? 'rgba(233,69,96,0.08)' : 'rgba(15,52,96,0.6)',
        color: '#e8e8f0', border: '1px solid rgba(255,255,255,0.06)',
        padding: '9px 20px', borderRadius: 6, cursor: loading ? 'not-allowed' : 'pointer',
        fontSize: 13, fontWeight: 500, transition: 'all 0.25s ease',
        whiteSpace: 'nowrap', lineHeight: 1, backdropFilter: 'blur(8px)',
        letterSpacing: 0.2, opacity: loading ? 0.6 : 1,
        ...variantStyles[variant || ''],
        ...style,
      }}
      onMouseEnter={e => {
        if (loading || disabled) return;
        const colors: Record<string, string> = {
          success: '#2ecc71', danger: '#e74c3c', info: '#3498db', warning: '#f39c12',
        };
        e.currentTarget.style.background = colors[variant || ''] || '#e94560';
        e.currentTarget.style.color = '#fff';
        e.currentTarget.style.borderColor = colors[variant || ''] || '#e94560';
        e.currentTarget.style.transform = 'translateY(-2px)';
        e.currentTarget.style.boxShadow = `0 8px 24px rgba(${variant === 'success' ? '46,204,113' : variant === 'danger' ? '231,76,60' : variant === 'info' ? '52,152,219' : '233,69,96'},0.3)`;
      }}
      onMouseLeave={e => {
        if (loading || disabled) return;
        e.currentTarget.style.background = glass ? 'rgba(233,69,96,0.08)' : 'rgba(15,52,96,0.6)';
        e.currentTarget.style.color = '#e8e8f0';
        e.currentTarget.style.borderColor = 'rgba(255,255,255,0.06)';
        e.currentTarget.style.transform = 'none';
        e.currentTarget.style.boxShadow = 'none';
      }}
    >
      {loading && <span style={{ display: 'inline-block', width: 14, height: 14, border: '2px solid rgba(255,255,255,0.15)', borderTopColor: '#e94560', borderRadius: '50%', animation: 'spin 0.6s linear infinite' }} />}
      {!loading && children}
    </button>
  );
}
