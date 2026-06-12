import { useState, useRef } from 'react'
import { Card } from '../../components/common/Card'
import { FormRow, Input, Select, TextArea, Checkbox } from '../../components/common/FormRow'
import { LoadingButton } from '../../components/common/LoadingButton'
import { ResultBox } from '../../components/common/ResultBox'
import { useToast } from '../../components/common/Toast'
import * as api from './profile'

export default function ProfilePage() {
  const [result, setResult] = useState<unknown>(null)
  const { toast } = useToast()
  const fileRef = useRef<HTMLInputElement>(null)

  const handleAction = async (fn: () => Promise<any>) => {
    const res = await fn()
    setResult(res)
    if (!res.error) toast('Success', 'success')
    return res
  }

  const [displayName, setDisplayName] = useState('')
  const [bio, setBio] = useState('')
  const [phone, setPhone] = useState('')
  const [gender, setGender] = useState('')
  const [dob, setDob] = useState('')
  const [statusText, setStatusText] = useState('')
  const [statusType, setStatusType] = useState('online')
  const [pushToken, setPushToken] = useState('')
  const [pushPlatform, setPushPlatform] = useState('fcm')

  return (
    <div className="page-content">
      <h1 className="page-title">Profile</h1>
      <p className="page-subtitle">View and manage your user profile</p>

      <Card title="View Profile" badge="GET">
        <LoadingButton variant="info" onClick={() => handleAction(() => api.getProfile())}>View Profile</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Update Profile" badge="PUT">
        <FormRow>
          <Input placeholder="Display Name" value={displayName} onChange={e => setDisplayName(e.target.value)} />
          <Input placeholder="Phone" value={phone} onChange={e => setPhone(e.target.value)} />
        </FormRow>
        <FormRow>
          <TextArea placeholder="Bio" value={bio} onChange={e => setBio(e.target.value)} />
          <Select value={gender} onChange={e => setGender(e.target.value)}>
            <option value="">Gender</option>
            <option value="male">Male</option>
            <option value="female">Female</option>
            <option value="other">Other</option>
          </Select>
        </FormRow>
        <FormRow>
          <Input placeholder="Date of Birth" value={dob} onChange={e => setDob(e.target.value)} />
        </FormRow>
        <LoadingButton variant="warning" onClick={() => handleAction(() => api.updateProfile({ displayName, bio, phone, gender, dateOfBirth: dob }))}>Update Profile</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Upload Avatar" badge="POST">
        <FormRow>
          <input ref={fileRef} type="file" accept="image/*" style={{ color: '#9a9ab8', fontSize: 13, flex: 1 }} />
        </FormRow>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.uploadAvatar(fileRef.current?.files?.[0]!))}>Upload</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Update Status" badge="PUT">
        <FormRow>
          <Input placeholder="Status text" value={statusText} onChange={e => setStatusText(e.target.value)} />
          <Select value={statusType} onChange={e => setStatusType(e.target.value)}>
            <option value="online">Online</option>
            <option value="offline">Offline</option>
            <option value="busy">Busy</option>
            <option value="away">Away</option>
          </Select>
        </FormRow>
        <LoadingButton variant="warning" onClick={() => handleAction(() => api.updateStatus({ status: statusText || statusType }))}>Update Status</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Push Notifications" badge="PUT">
        <FormRow>
          <Input placeholder="Push Token" value={pushToken} onChange={e => setPushToken(e.target.value)} />
          <Select value={pushPlatform} onChange={e => setPushPlatform(e.target.value)}>
            <option value="fcm">FCM</option>
            <option value="apns">APNS</option>
          </Select>
        </FormRow>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.savePushToken({ token: pushToken, provider: pushPlatform }))}>Save Token</LoadingButton>
        <LoadingButton variant="warning" onClick={() => handleAction(() => api.testPush())}>Test Push</LoadingButton>
        <ResultBox data={result} />
      </Card>

      <Card title="Delete Account" badge="DELETE">
        <LoadingButton variant="danger" onClick={() => handleAction(() => api.deleteAccount())}>Delete Account</LoadingButton>
        <ResultBox data={result} />
      </Card>
    </div>
  )
}
