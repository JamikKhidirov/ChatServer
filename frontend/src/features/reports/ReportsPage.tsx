import { useState } from 'react'
import { Card } from '../../components/common/Card'
import { FormRow, Input, TextArea, Select } from '../../components/common/FormRow'
import { LoadingButton } from '../../components/common/LoadingButton'
import { ResultBox } from '../../components/common/ResultBox'
import { useToast } from '../../components/common/Toast'
import * as api from './reports'

export default function ReportsPage() {
  const [result, setResult] = useState<unknown>(null)
  const { toast } = useToast()

  const handleAction = async (fn: () => Promise<any>, successMsg?: string) => {
    const res = await fn()
    setResult(res)
    if (!res.error && successMsg) toast(successMsg, 'success')
    return res
  }

  const [reportMessageId, setReportMessageId] = useState('')
  const [reportReason, setReportReason] = useState('')
  const [reportDescription, setReportDescription] = useState('')
  const [listStatus, setListStatus] = useState('')
  const [resolveReportId, setResolveReportId] = useState('')
  const [resolveStatus, setResolveStatus] = useState('resolved')

  return (
    <div className="page-content">
      <h1 className="page-title">Reports</h1>
      <p className="page-subtitle">Report messages and manage report resolutions</p>
      <Card title="Create Report" badge="POST">
        <FormRow>
          <Input placeholder="Message ID" value={reportMessageId} onChange={e => setReportMessageId(e.target.value)} />
        </FormRow>
        <FormRow>
          <Input placeholder="Reason" value={reportReason} onChange={e => setReportReason(e.target.value)} />
        </FormRow>
        <FormRow>
          <TextArea placeholder="Description" value={reportDescription} onChange={e => setReportDescription(e.target.value)} />
        </FormRow>
        <LoadingButton variant="danger" onClick={() => handleAction(() => api.createReport({ target_id: reportMessageId, target_type: 'message', reason: reportReason }), 'Report created')}>Create Report</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="List Reports" badge="GET">
        <FormRow>
          <Input placeholder="Status filter (optional)" value={listStatus} onChange={e => setListStatus(e.target.value)} />
        </FormRow>
        <LoadingButton variant="info" onClick={() => handleAction(() => api.listReports(listStatus))}>List Reports</LoadingButton>
        <ResultBox data={result} />
      </Card>
      <Card title="Resolve Report" badge="PATCH">
        <FormRow>
          <Input placeholder="Report ID" value={resolveReportId} onChange={e => setResolveReportId(e.target.value)} />
          <Select value={resolveStatus} onChange={e => setResolveStatus(e.target.value)}>
            <option value="resolved">Resolved</option>
            <option value="dismissed">Dismissed</option>
            <option value="pending">Pending</option>
          </Select>
        </FormRow>
        <LoadingButton variant="success" onClick={() => handleAction(() => api.resolveReport(resolveReportId, resolveStatus), 'Report resolved')}>Resolve Report</LoadingButton>
        <ResultBox data={result} />
      </Card>
    </div>
  )
}
