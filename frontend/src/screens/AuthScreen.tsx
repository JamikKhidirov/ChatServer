import { useState } from 'react'
import { api } from '../api/client'
import * as authApi from '../features/auth/auth'

interface Props { onLogin: () => void }

export default function AuthScreen({ onLogin }: Props) {
  const [tab, setTab] = useState<'login' | 'register'>('login')
  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [displayName, setDisplayName] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleSubmit = async () => {
    setLoading(true); setError('')
    try {
      let res
      if (tab === 'register') {
        res = await authApi.register({ username, email, password, displayName })
      } else {
        res = await authApi.login({ email, password })
      }
      if (res.error) { setError(typeof res.data === 'string' ? res.data : 'Auth failed'); return }
      onLogin()
    } catch { setError('Connection error') }
    finally { setLoading(false) }
  }

  return (
    <div className="auth-screen">
      <div className="auth-card">
        <div className="auth-logo">📨</div>
        <h1 className="auth-title">Go Messenger</h1>
        <p className="auth-subtitle">Connect with anyone, anywhere</p>

        <div className="auth-tabs">
          <button className={`auth-tab ${tab === 'login' ? 'active' : ''}`} onClick={() => setTab('login')}>Sign In</button>
          <button className={`auth-tab ${tab === 'register' ? 'active' : ''}`} onClick={() => setTab('register')}>Register</button>
        </div>

        {tab === 'register' && (
          <div className="auth-field">
            <label>Username</label>
            <input type="text" value={username} onChange={e => setUsername(e.target.value)} placeholder="your_username" />
          </div>
        )}
        {tab === 'register' && (
          <div className="auth-field">
            <label>Display Name</label>
            <input type="text" value={displayName} onChange={e => setDisplayName(e.target.value)} placeholder="Your Name" />
          </div>
        )}
        <div className="auth-field">
          <label>Email</label>
          <input type="email" value={email} onChange={e => setEmail(e.target.value)} placeholder="you@example.com" />
        </div>
        <div className="auth-field">
          <label>Password</label>
          <input type="password" value={password} onChange={e => setPassword(e.target.value)} placeholder="••••••••" />
        </div>

        {error && <div className="auth-error">{error}</div>}

        <button className="auth-btn" onClick={handleSubmit} disabled={loading}>
          {loading ? <span className="spinner" /> : (tab === 'login' ? 'Sign In' : 'Create Account')}
        </button>
      </div>
    </div>
  )
}