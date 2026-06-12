import { useState } from 'react'
import { Card } from '../../components/common/Card'
import { FormRow, Input, TextArea } from '../../components/common/FormRow'
import { LoadingButton } from '../../components/common/LoadingButton'
import { ResultBox } from '../../components/common/ResultBox'
import { useToast } from '../../components/common/Toast'
import * as api from './security'

export default function SecurityPage() {
  const [result, setResult] = useState<unknown>(null)
  const { toast } = useToast()

  const handleAction = async (fn: () => Promise<any>, successMsg?: string) => {
    const res = await fn()
    setResult(res)
    if (!res.error && successMsg) toast(successMsg, 'success')
    return res
  }

  const [captchaToken, setCaptchaToken] = useState('')
  const [captchaSolution, setCaptchaSolution] = useState('')
  const [publicKey, setPublicKey] = useState('')
  const [privateKey, setPrivateKey] = useState('')
  const [e2eUserId, setE2eUserId] = useState('')
  const [email, setEmail] = useState('')
  const [emailCode, setEmailCode] = useState('')
  const [phone, setPhone] = useState('')
  const [phoneCode, setPhoneCode] = useState('')

  return (
    <div className="page-content">
      <h1 className="page-title">Security</h1>
      <p className="page-subtitle">Captcha, E2E encryption, and verification</p>
      <Card title="Generate Captcha" badge="GET">
        <LoadingButton variant="info" onClick={() => handleAction(() => api.generateCaptcha())}>Generate Captcha</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="Verify Captcha" badge="POST">
        <FormRow>
          <Input placeholder="Captcha token" value={captchaToken} onChange={e => setCaptchaToken(e.target.value)} />
          <Input placeholder="Solution" value={captchaSolution} onChange={e => setCaptchaSolution(e.target.value)} />
        </FormRow>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.verifyCaptcha({ token: captchaToken, solution: captchaSolution }), 'Captcha verified')}>Verify Captcha</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="Register E2E Key" badge="POST">
        <FormRow>
          <TextArea placeholder="Public key" value={publicKey} onChange={e => setPublicKey(e.target.value)} />
        </FormRow>
        <FormRow>
          <TextArea placeholder="Private key (encrypted)" value={privateKey} onChange={e => setPrivateKey(e.target.value)} />
        </FormRow>
        <LoadingButton variant="success" onClick={() => handleAction(() => api.registerE2EKey({ public_key: publicKey }), 'E2E key registered')}>Register E2E Key</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="Get E2E Public Key" badge="GET">
        <FormRow>
          <Input placeholder="User ID" value={e2eUserId} onChange={e => setE2eUserId(e.target.value)} />
        </FormRow>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getE2EPublicKey(e2eUserId))}>Get Public Key</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="Email Verification" badge="POST">
        <FormRow>
          <Input placeholder="Email address" value={email} onChange={e => setEmail(e.target.value)} />
          <LoadingButton variant="info" onClick={() => handleAction(() => api.sendEmailVerification(email), 'Verification email sent')}>Send Code</LoadingButton>
        </FormRow>
        <FormRow>
          <Input placeholder="Verification code" value={emailCode} onChange={e => setEmailCode(e.target.value)} />
          <LoadingButton variant="success" onClick={() => handleAction(() => api.verifyEmail(emailCode), 'Email verified')}>Verify Code</LoadingButton>
        </FormRow>
        <ResultBox data={result} />
      </Card>
      <Card title="Phone Verification" badge="POST">
        <FormRow>
          <Input placeholder="Phone number" value={phone} onChange={e => setPhone(e.target.value)} />
          <LoadingButton variant="info" onClick={() => handleAction(() => api.sendPhoneVerification(phone), 'Verification SMS sent')}>Send Code</LoadingButton>
        </FormRow>
        <FormRow>
          <Input placeholder="Verification code" value={phoneCode} onChange={e => setPhoneCode(e.target.value)} />
          <LoadingButton variant="success" onClick={() => handleAction(() => api.verifyPhone(phoneCode), 'Phone verified')}>Verify Code</LoadingButton>
        </FormRow>
        <ResultBox data={result} />
      </Card>
    </div>
  )
}
