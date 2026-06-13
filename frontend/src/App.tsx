import { useState } from 'react'
import AuthScreen from './screens/AuthScreen'
import MainScreen from './screens/MainScreen'
import { api } from './api/client'

export default function App() {
  const [screen, setScreen] = useState<'auth' | 'main'>(api.getToken() ? 'main' : 'auth')

  return (
    <div className="app">
      {screen === 'auth' ? (
        <AuthScreen onLogin={() => setScreen('main')} />
      ) : (
        <MainScreen onLogout={() => { api.clearToken(); setScreen('auth') }} />
      )}
    </div>
  )
}