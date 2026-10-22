import { useState, useEffect } from 'react'
import {
  Wallet,
  CheckCircle2,
  XCircle,
  Clock,
  RefreshCw,
  Loader2,
  AlertCircle,
  Building2,
  FileText,
} from 'lucide-react'
import { toast } from 'sonner'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '@/components/ui/dialog'
import api from '@/lib/api'
import { formatRupiah, formatDate } from '@/lib/utils'
import type { WithdrawalRequest, PaginatedResponse } from '@/types'

export function WithdrawalApprovalPage() {
  const [withdrawals, setWithdrawals] = useState<WithdrawalRequest[]>([])
  const [loading, setLoading] = useState(true)
  const [statusFilter, setStatusFilter] = useState<string>('pending')
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [totalCount, setTotalCount] = useState(0)

  // Action Dialog state
  const [selectedReq, setSelectedReq] = useState<WithdrawalRequest | null>(null)
  const [actionType, setActionType] = useState<'approve' | 'reject' | null>(null)
  const [actionNotes, setActionNotes] = useState('')
  const [actionLoading, setActionLoading] = useState(false)
  const [dialogOpen, setDialogOpen] = useState(false)

  const fetchWithdrawals = async () => {
    try {
      setLoading(true)
      const res = await api.get<PaginatedResponse<WithdrawalRequest>>('/admin/withdrawals', {
        params: {
          page,
          limit: 15,
          status: statusFilter !== 'all' ? statusFilter : undefined,
        },
      })
      setWithdrawals(res.data.data || [])
      setTotalPages(res.data.meta?.total_pages || 1)
      setTotalCount(res.data.meta?.total || 0)
    } catch (err: any) {
      toast.error('Gagal memuat antrean penarikan dana')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchWithdrawals()
  }, [page, statusFilter])

  const openActionDialog = (req: WithdrawalRequest, action: 'approve' | 'reject') => {
    setSelectedReq(req)
    setActionType(action)
    setActionNotes(action === 'approve' ? 'Transfer berhasil diproses manual.' : '')
    setDialogOpen(true)
  }

  const handleProcessAction = async () => {
    if (!selectedReq || !actionType) return

    if (actionType === 'reject' && !actionNotes.trim()) {
      toast.error('Mohon berikan alasan penolakan pada kolom catatan')
      return
    }

    try {
      setActionLoading(true)
      await api.patch(`/admin/withdrawals/${selectedReq.id}`, {
        action: actionType,
        notes: actionNotes.trim(),
      })
      toast.success(
        actionType === 'approve'
          ? `Penarikan dana sebesar ${formatRupiah(selectedReq.amount)} telah disetujui`
          : `Penarikan dana telah ditolak dan saldo dikembalikan ke fotografer`
      )
      setDialogOpen(false)
      setSelectedReq(null)
      fetchWithdrawals()
    } catch (err: any) {
      toast.error(err.response?.data?.message || 'Gagal memproses penarikan')
    } finally {
      setActionLoading(false)
    }
  }

  return (
    <div className="space-y-6 max-w-7xl mx-auto pb-12">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Antrean Penarikan Dana (Withdrawal)</h1>
          <p className="text-sm text-muted-foreground mt-1">
            Tinjau dan proses permohonan penarikan saldo pendapatan fotografer ke rekening bank.
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={fetchWithdrawals} disabled={loading} className="gap-2 self-start">
          <RefreshCw size={14} className={loading ? 'animate-spin' : ''} />
          Segarkan Data
        </Button>
      </div>

      {/* Filter Tabs */}
      <Card>
        <CardContent className="p-3">
          <div className="flex items-center gap-1.5 overflow-x-auto">
            {[
              { id: 'pending', label: 'Menunggu Approval (Pending)' },
              { id: 'approved', label: 'Disetujui (Approved)' },
              { id: 'rejected', label: 'Ditolak (Rejected)' },
              { id: 'all', label: 'Semua Status' },
            ].map((tab) => (
              <Button
                key={tab.id}
                variant={statusFilter === tab.id ? 'default' : 'ghost'}
                size="sm"
                onClick={() => {
                  setStatusFilter(tab.id)
                  setPage(1)
                }}
                className={statusFilter === tab.id ? 'bg-brand text-white' : 'text-muted-foreground'}
              >
                {tab.label}
              </Button>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Withdrawals Table */}
      <Card>
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <table className="w-full text-sm text-left">
              <thead className="bg-muted/40 border-b text-xs text-muted-foreground uppercase tracking-wider">
                <tr>
                  <th className="py-3 px-4 font-medium">Fotografer</th>
                  <th className="py-3 px-4 font-medium">Tujuan Rekening Bank</th>
                  <th className="py-3 px-4 font-medium">Nominal Penarikan</th>
                  <th className="py-3 px-4 font-medium">Waktu Pengajuan</th>
                  <th className="py-3 px-4 font-medium">Status</th>
                  <th className="py-3 px-4 font-medium">Catatan</th>
                  <th className="py-3 px-4 font-medium text-right">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border/60">
                {loading ? (
                  <tr>
                    <td colSpan={7} className="py-12 text-center text-muted-foreground">
                      <div className="flex items-center justify-center gap-2">
                        <Loader2 size={18} className="animate-spin text-brand" />
                        <span>Memuat daftar penarikan...</span>
                      </div>
                    </td>
                  </tr>
                ) : withdrawals.length === 0 ? (
                  <tr>
                    <td colSpan={7} className="py-12 text-center text-muted-foreground">
                      Tidak ada permohonan penarikan dana pada kategori ini.
                    </td>
                  </tr>
                ) : (
                  withdrawals.map((w) => (
                    <tr key={w.id} className="hover:bg-muted/30 transition-colors">
                      <td className="py-3.5 px-4">
                        <div className="font-medium text-foreground">
                          {w.photographer?.name || 'Fotografer'}
                        </div>
                        <div className="text-xs text-muted-foreground">{w.photographer?.email || '-'}</div>
                      </td>

                      <td className="py-3.5 px-4">
                        <div className="flex items-center gap-1.5 font-medium text-foreground">
                          <Building2 size={13} className="text-muted-foreground" />
                          <span>{w.bank_name}</span>
                          <span className="font-mono text-xs">({w.account_number})</span>
                        </div>
                        <div className="text-xs text-muted-foreground">a.n. {w.account_holder}</div>
                      </td>

                      <td className="py-3.5 px-4 font-bold text-foreground">
                        {formatRupiah(w.amount)}
                      </td>

                      <td className="py-3.5 px-4 text-xs text-muted-foreground">
                        {formatDate(w.created_at)}
                      </td>

                      <td className="py-3.5 px-4">
                        <Badge
                          variant={
                            w.status === 'approved'
                              ? 'default'
                              : w.status === 'rejected'
                              ? 'destructive'
                              : 'secondary'
                          }
                          className="text-xs capitalize font-medium"
                        >
                          {w.status === 'approved' && '✓ Disetujui'}
                          {w.status === 'rejected' && '✗ Ditolak'}
                          {w.status === 'pending' && '⏳ Pending'}
                        </Badge>
                      </td>

                      <td className="py-3.5 px-4 text-xs text-muted-foreground max-w-[200px] truncate">
                        {w.notes || '-'}
                      </td>

                      <td className="py-3.5 px-4 text-right">
                        {w.status === 'pending' ? (
                          <div className="flex items-center justify-end gap-1.5">
                            <Button
                              size="sm"
                              className="h-8 px-2.5 text-xs bg-emerald-600 hover:bg-emerald-700 text-white"
                              onClick={() => openActionDialog(w, 'approve')}
                            >
                              <CheckCircle2 size={13} className="mr-1" />
                              Setujui
                            </Button>
                            <Button
                              size="sm"
                              variant="outline"
                              className="h-8 px-2.5 text-xs text-destructive hover:bg-destructive/10 border-destructive/30"
                              onClick={() => openActionDialog(w, 'reject')}
                            >
                              <XCircle size={13} className="mr-1" />
                              Tolak
                            </Button>
                          </div>
                        ) : (
                          <span className="text-xs text-muted-foreground italic">Selesai</span>
                        )}
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex items-center justify-between p-4 border-t text-xs text-muted-foreground">
              <div>
                Menampilkan halaman {page} dari {totalPages} (Total {totalCount} permohonan)
              </div>
              <div className="flex gap-1.5">
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => setPage((p) => Math.max(p - 1, 1))}
                  disabled={page <= 1}
                >
                  Sebelumnya
                </Button>
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => setPage((p) => Math.min(p + 1, totalPages))}
                  disabled={page >= totalPages}
                >
                  Berikutnya
                </Button>
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Action Dialog (Approve or Reject) */}
      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              {actionType === 'approve' ? (
                <>
                  <CheckCircle2 size={18} className="text-emerald-600" />
                  Setujui Penarikan Dana
                </>
              ) : (
                <>
                  <XCircle size={18} className="text-destructive" />
                  Tolak Penarikan Dana
                </>
              )}
            </DialogTitle>
            <DialogDescription>
              {actionType === 'approve'
                ? 'Konfirmasi bahwa transfer dana ke rekening fotografer telah berhasil diproses.'
                : 'Menolak permohonan ini akan mengembalikan saldo nominal ke dompet fotografer.'}
            </DialogDescription>
          </DialogHeader>

          {selectedReq && (
            <div className="space-y-4 py-2 text-sm">
              <div className="p-3 bg-muted/40 rounded-lg space-y-1.5 text-xs">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Fotografer:</span>
                  <span className="font-medium text-foreground">{selectedReq.photographer?.name}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Nominal Transfer:</span>
                  <span className="font-bold text-foreground text-sm">{formatRupiah(selectedReq.amount)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Bank & Rekening:</span>
                  <span className="font-mono text-foreground">
                    {selectedReq.bank_name} - {selectedReq.account_number}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Nama Pemilik:</span>
                  <span className="font-medium text-foreground">{selectedReq.account_holder}</span>
                </div>
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="notes" className="text-xs">
                  Catatan {actionType === 'reject' ? '(Wajib - alasan penolakan)' : '(Opsional)'}:
                </Label>
                <Textarea
                  id="notes"
                  value={actionNotes}
                  onChange={(e) => setActionNotes(e.target.value)}
                  placeholder={
                    actionType === 'reject'
                      ? 'Contoh: Nomor rekening tidak cocok dengan nama akun bank'
                      : 'Catatan bukti transfer atau ID referensi bank'
                  }
                  rows={3}
                />
              </div>
            </div>
          )}

          <div className="flex justify-end gap-2 pt-2">
            <Button variant="outline" onClick={() => setDialogOpen(false)} disabled={actionLoading}>
              Batal
            </Button>
            <Button
              className={
                actionType === 'approve'
                  ? 'bg-emerald-600 hover:bg-emerald-700 text-white'
                  : 'bg-destructive hover:bg-destructive/90 text-white'
              }
              onClick={handleProcessAction}
              disabled={actionLoading}
            >
              {actionLoading ? (
                <>
                  <Loader2 size={14} className="animate-spin mr-1.5" /> Memproses...
                </>
              ) : actionType === 'approve' ? (
                'Konfirmasi Setujui'
              ) : (
                'Konfirmasi Tolak'
              )}
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  )
}
export default WithdrawalApprovalPage
