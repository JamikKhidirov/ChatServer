export function FormRow({ children, style }: { children: React.ReactNode; style?: React.CSSProperties }) {
  return (
    <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8, marginBottom: 10, ...style }}>
      {children}
    </div>
  );
}

export function Input(props: React.InputHTMLAttributes<HTMLInputElement>) {
  return (
    <input
      {...props}
      style={{
        background: 'rgba(13,13,26,0.7)', border: '1px solid rgba(42,42,74,0.6)',
        color: '#e8e8f0', padding: '9px 14px', borderRadius: 6, fontSize: 13,
        outline: 'none', transition: 'all 0.25s ease', fontFamily: 'inherit',
        backdropFilter: 'blur(8px)', flex: 1, minWidth: 100,
        ...(props.style || {}),
      }}
      onFocus={e => {
        e.currentTarget.style.borderColor = '#e94560';
        e.currentTarget.style.boxShadow = '0 0 0 3px rgba(233,69,96,0.15), 0 0 20px rgba(233,69,96,0.15)';
      }}
      onBlur={e => {
        e.currentTarget.style.borderColor = 'rgba(42,42,74,0.6)';
        e.currentTarget.style.boxShadow = 'none';
      }}
    />
  );
}

export function Select(props: React.SelectHTMLAttributes<HTMLSelectElement>) {
  return (
    <select
      {...props}
      style={{
        background: 'rgba(13,13,26,0.7)', border: '1px solid rgba(42,42,74,0.6)',
        color: '#e8e8f0', padding: '9px 14px', borderRadius: 6, fontSize: 13,
        outline: 'none', transition: 'all 0.25s ease', fontFamily: 'inherit',
        backdropFilter: 'blur(8px)', flex: 1, minWidth: 100, cursor: 'pointer',
        ...(props.style || {}),
      }}
    >
      {props.children}
    </select>
  );
}

export function TextArea(props: React.TextareaHTMLAttributes<HTMLTextAreaElement>) {
  return (
    <textarea
      {...props}
      style={{
        background: 'rgba(13,13,26,0.7)', border: '1px solid rgba(42,42,74,0.6)',
        color: '#e8e8f0', padding: '9px 14px', borderRadius: 6, fontSize: 13,
        outline: 'none', transition: 'all 0.25s ease', fontFamily: "'JetBrains Mono','Cascadia Code',monospace",
        backdropFilter: 'blur(8px)', flex: 1, minWidth: 100, minHeight: 56, resize: 'vertical',
        lineHeight: 1.5,
        ...(props.style || {}),
      }}
    />
  );
}

export function Checkbox({ label, ...props }: { label: string } & React.InputHTMLAttributes<HTMLInputElement>) {
  return (
    <label style={{ display: 'flex', alignItems: 'center', gap: 6, color: '#9a9ab8', fontSize: 13, cursor: 'pointer', flex: '0 0 auto' }}>
      <input type="checkbox" {...props} style={{ width: 16, height: 16, accentColor: '#e94560', cursor: 'pointer', flex: '0 0 auto' }} />
      {label}
    </label>
  );
}
