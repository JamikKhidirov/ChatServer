import { useState, lazy, Suspense } from 'react'
import { Sidebar } from './components/Layout/Sidebar'
import { useToast } from './components/common/Toast'

const pages: Record<string, React.LazyExoticComponent<React.ComponentType>> = {
  auth: lazy(() => import('./features/auth/AuthPage')),
  profile: lazy(() => import('./features/profile/ProfilePage')),
  settings: lazy(() => import('./features/settings/SettingsPage')),
  contacts: lazy(() => import('./features/contacts/ContactsPage')),
  chats: lazy(() => import('./features/chats/ChatsPage')),
  messages: lazy(() => import('./features/messages/MessagesPage')),
  location: lazy(() => import('./features/location/LocationPage')),
  saved: lazy(() => import('./features/saved/SavedPage')),
  emojis: lazy(() => import('./features/emojis/EmojisPage')),
  voice: lazy(() => import('./features/voice/VoicePage')),
  polls: lazy(() => import('./features/polls/PollsPage')),
  stickers: lazy(() => import('./features/stickers/StickersPage')),
  drafts: lazy(() => import('./features/drafts/DraftsPage')),
  scheduled: lazy(() => import('./features/scheduled/ScheduledPage')),
  sessions: lazy(() => import('./features/sessions/SessionsPage')),
  bots: lazy(() => import('./features/bots/BotsPage')),
  gifs: lazy(() => import('./features/gifs/GifsPage')),
  calls: lazy(() => import('./features/calls/CallsPage')),
  block: lazy(() => import('./features/block/BlockPage')),
  security: lazy(() => import('./features/security/SecurityPage')),
  bookmarks: lazy(() => import('./features/bookmarks/BookmarksPage')),
  reports: lazy(() => import('./features/reports/ReportsPage')),
  admin: lazy(() => import('./features/admin/AdminPage')),
  preview: lazy(() => import('./features/preview/PreviewPage')),
  search: lazy(() => import('./features/search/SearchPage')),
}

function LoadingPage() {
  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '60vh', color: '#5a5a7a', fontSize: 14 }}>
      <div style={{ textAlign: 'center' }}>
        <div style={{ width: 32, height: 32, border: '3px solid rgba(255,255,255,0.1)', borderTopColor: '#e94560', borderRadius: '50%', animation: 'spin 0.6s linear infinite', margin: '0 auto 12px' }} />
        Loading...
      </div>
    </div>
  )
}

export default function App() {
  const [activeTab, setActiveTab] = useState('auth')
  const { toast } = useToast()

  const PageComponent = pages[activeTab]

  return (
    <div style={{ display: 'flex', minHeight: '100vh' }}>
      <Sidebar activeTab={activeTab} onTabChange={setActiveTab} />
      <main style={{ flex: 1, padding: '32px 36px', overflowY: 'auto', maxHeight: '100vh', minWidth: 0 }}>
        <Suspense fallback={<LoadingPage />}>
          {PageComponent && <PageComponent />}
        </Suspense>
      </main>
    </div>
  )
}
