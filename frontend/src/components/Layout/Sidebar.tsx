import { api } from '../../api/client';

interface SidebarProps {
  activeTab: string;
  onTabChange: (tab: string) => void;
}

const TABS = [
  { id: 'auth', label: 'Auth', icon: '🔒' },
  { id: 'profile', label: 'Profile', icon: '👤' },
  { id: 'settings', label: 'Settings', icon: '⚙️' },
  { id: 'contacts', label: 'Contacts', icon: '📱' },
  { id: 'chats', label: 'Chats', icon: '💬' },
  { id: 'messages', label: 'Messages', icon: '✉️' },
  { id: 'location', label: 'Location', icon: '📍' },
  { id: 'saved', label: 'Saved', icon: '⭐' },
  { id: 'emojis', label: 'Emojis', icon: '🤩' },
  { id: 'voice', label: 'Voice', icon: '🎤' },
  { id: 'polls', label: 'Polls', icon: '📊' },
  { id: 'stickers', label: 'Stickers', icon: '📷' },
  { id: 'drafts', label: 'Drafts', icon: '📝' },
  { id: 'scheduled', label: 'Scheduled', icon: '⏰' },
  { id: 'sessions', label: 'Sessions', icon: '🖥️' },
  { id: 'bots', label: 'Bots', icon: '🤖' },
  { id: 'gifs', label: 'GIFs', icon: '👻' },
  { id: 'calls', label: 'Calls', icon: '📞' },
  { id: 'block', label: 'Block', icon: '⛔' },
  { id: 'security', label: 'Security', icon: '🛡️' },
  { id: 'bookmarks', label: 'Bookmarks', icon: '🔖' },
  { id: 'reports', label: 'Reports', icon: '⚠️' },
  { id: 'admin', label: 'Admin', icon: '🛠️' },
  { id: 'preview', label: 'Preview', icon: '🔗' },
  { id: 'search', label: 'Search', icon: '🔍' },
];

export function Sidebar({ activeTab, onTabChange }: SidebarProps) {
  const token = api.getToken();

  return (
    <div style={{
      width: 220, background: '#13132a',
      borderRight: '1px solid rgba(42,42,74,0.6)',
      display: 'flex', flexDirection: 'column',
      height: '100vh', position: 'sticky', top: 0,
      overflowY: 'auto', flexShrink: 0,
    }}>
      <div style={{
        padding: '24px 16px 16px',
        borderBottom: '1px solid rgba(255,255,255,0.06)',
        background: 'linear-gradient(180deg, rgba(233,69,96,0.06) 0%, transparent 100%)',
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 10, fontSize: 17, fontWeight: 800, color: '#e94560' }}>
          <span style={{fontSize: 22}}>📨</span>
          Go Messenger
        </div>
        <p style={{ fontSize: 11, color: '#5a5a7a', marginTop: 4, paddingLeft: 32 }}>API Tester v3.0</p>
      </div>

      <nav style={{ flex: 1, padding: '6px 0' }}>
        {TABS.map(tab => (
          <button
            key={tab.id}
            onClick={() => onTabChange(tab.id)}
            style={{
              display: 'flex', alignItems: 'center', gap: 10,
              width: '100%', padding: '9px 16px', background: 'none',
              border: 'none', borderLeft: '3px solid transparent',
              color: activeTab === tab.id ? '#e94560' : '#9a9ab8',
              cursor: 'pointer', textAlign: 'left', fontSize: 12.5,
              transition: 'all 0.25s ease',
              backgroundColor: activeTab === tab.id
                ? 'linear-gradient(90deg, rgba(233,69,96,0.1) 0%, transparent 100%)'
                : 'transparent',
              fontWeight: activeTab === tab.id ? 600 : 400,
              letterSpacing: 0.2,
            }}
            onMouseEnter={e => { if (activeTab !== tab.id) e.currentTarget.style.background = 'rgba(15,52,96,0.4)'; }}
            onMouseLeave={e => { if (activeTab !== tab.id) e.currentTarget.style.background = 'none'; }}
          >
            <span style={{ fontSize: 15, width: 20, textAlign: 'center', opacity: 0.85 }}>{tab.icon}</span>
            {tab.label}
          </button>
        ))}
      </nav>

      <div style={{
        padding: '12px 16px', borderTop: '1px solid rgba(255,255,255,0.06)',
        fontSize: 11, color: '#5a5a7a', display: 'flex', alignItems: 'center',
        background: 'rgba(0,0,0,0.15)',
      }}>
        <span style={{
          display: 'inline-block', width: 8, height: 8, borderRadius: '50%',
          marginRight: 8, background: token ? '#2ecc71' : '#5a5a7a',
          boxShadow: token ? '0 0 8px rgba(46,204,113,0.4)' : 'none',
        }} />
        {token ? 'Authenticated' : 'Not connected'}
      </div>
    </div>
  );
}
