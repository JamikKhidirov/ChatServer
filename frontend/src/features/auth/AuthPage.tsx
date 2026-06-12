import { useState } from 'react'
import { Card } from '../../components/common/Card'
import { FormRow, Input, Select, TextArea, Checkbox } from '../../components/common/FormRow'
import { LoadingButton } from '../../components/common/LoadingButton'
import { ResultBox } from '../../components/common/ResultBox'
import { useToast } from '../../components/common/Toast'
import * as api from './auth'

export default function AuthPage() {
  const [result, setResult] = useState<unknown>(null)
  const { toast } = useToast()

  const handleAction = async (fn: () => Promise<any>) => {
    const res = await fn()
    setResult(res)
    if (!res.error) toast('Success', 'success')
    return res
  }

  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [displayName, setDisplayName] = useState('')
  const [oldPassword, setOldPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [loginEmail, setLoginEmail] = useState('')
  const [loginPassword, setLoginPassword] = useState('')
  const [codeEmail, setCodeEmail] = useState('')
  const [codeEmailCode, setCodeEmailCode] = useState('')
  const [codePhone, setCodePhone] = useState('')
  const [codePhoneCode, setCodePhoneCode] = useState('')

  const token = api.getToken()

  return (
    <div className="page-content">
      <h1 className="page-title">Auth</h1>
      <p className="page-subtitle">Register, login, and manage your JWT token</p>

      <Card title="Register New Account" badge="POST">
        <FormRow>
          <Input placeholder="Username *" value={username} onChange={e => setUsername(e.target.value)} />
          <Input placeholder="Email *" value={email} onChange={e => setEmail(e.target.value)} />
        </FormRow>
        <FormRow>
          <Input placeholder="Password *" type="password" value={password} onChange={e => setPassword(e.target.value)} />
          <Input placeholder="Display Name *" value={displayName} onChange={e => setDisplayName(e.target.value)} />
        </FormRow>
        <LoadingButton variant="success" onClick={() => handleAction(() => api.register({ username, email, password, displayName }))}>Register</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Login" badge="POST">
        <FormRow>
          <Input placeholder="Email *" value={loginEmail} onChange={e => setLoginEmail(e.target.value)} />
          <Input placeholder="Password *" type="password" value={loginPassword} onChange={e => setLoginPassword(e.target.value)} />
        </FormRow>
        <LoadingButton variant="success" onClick={() => handleAction(() => api.login({ email: loginEmail, password: loginPassword }))}>Login</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Refresh Token" badge="GET">
        <LoadingButton variant="info" onClick={() => handleAction(() => api.refreshToken())}>Refresh</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Change Password" badge="PUT">
        <FormRow>
          <Input placeholder="Old Password *" type="password" value={oldPassword} onChange={e => setOldPassword(e.target.value)} />
          <Input placeholder="New Password *" type="password" value={newPassword} onChange={e => setNewPassword(e.target.value)} />
        </FormRow>
        <LoadingButton variant="warning" onClick={() => handleAction(() => api.changePassword(oldPassword, newPassword))}>Change Password</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Email Login Code" badge="POST">
        <FormRow>
          <Input placeholder="Email" value={codeEmail} onChange={e => setCodeEmail(e.target.value)} />
        </FormRow>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.sendEmailLoginCode(codeEmail))}>Send Code</LoadingButton>
        <FormRow>
          <Input placeholder="Code" value={codeEmailCode} onChange={e => setCodeEmailCode(e.target.value)} />
        </FormRow>
        <LoadingButton variant="success" onClick={() => handleAction(() => api.verifyEmailLoginCode(codeEmail, codeEmailCode))}>Verify Code</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Phone Login Code" badge="POST">
        <FormRow>
          <Input placeholder="Phone" value={codePhone} onChange={e => setCodePhone(e.target.value)} />
        </FormRow>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.sendPhoneLoginCode(codePhone))}>Send Code</LoadingButton>
        <FormRow>
          <Input placeholder="Code" value={codePhoneCode} onChange={e => setCodePhoneCode(e.target.value)} />
        </FormRow>
        <LoadingButton variant="success" onClick={() => handleAction(() => api.verifyPhoneLoginCode(codePhone, codePhoneCode))}>Verify Code</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Current JWT Token">
        <div style={{ fontFamily: "'JetBrains Mono',monospace", fontSize: 11, color: '#9a9ab8', wordBreak: 'break-all', background: 'rgba(0,0,0,0.3)', padding: 12, borderRadius: 6 }}>
          {token || 'No token'}
        </div>
      </Card>
    </div>
  )
}
