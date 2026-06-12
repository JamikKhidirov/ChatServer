export function Card({ title, badge, children, style }: {
  title?: string;
  badge?: string;
  children: React.ReactNode;
  style?: React.CSSProperties;
}) {
  return (
    <div style={{
      background: 'rgba(22,33,62,0.85)',
      border: '1px solid rgba(255,255,255,0.06)',
      borderRadius: 12,
      padding: '18px 22px',
      marginBottom: 18,
      boxShadow: '0 8px 32px rgba(0,0,0,0.4)',
      backdropFilter: 'blur(12px)',
      position: 'relative',
      overflow: 'hidden',
      ...style,
    }}>
      {title && (
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 14 }}>
          <h4 style={{ fontSize: 14, fontWeight: 600, color: '#e94560', display: 'flex', alignItems: 'center', gap: 8, margin: 0 }}>
            {title}
            {badge && (
              <span style={{
                fontSize: 9, background: 'linear-gradient(135deg,#e94560,#ff6b81)',
                color: '#fff', padding: '2px 9px', borderRadius: 12,
                fontWeight: 600, letterSpacing: 0.5, textTransform: 'uppercase',
                boxShadow: '0 2px 8px rgba(233,69,96,0.3)'
              }}>{badge}</span>
            )}
          </h4>
        </div>
      )}
      {children}
    </div>
  );
}
