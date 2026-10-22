import { useState, useEffect } from 'react'
import {
  Sparkles,
  Search,
  Calendar,
  CheckCircle2,
  Clock,
  RefreshCw,
  Loader2,
  Users,
  DollarSign,
} from 'lucide-react'
import { toast } from 'sonner'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import api from '@/lib/api'
import { formatRupiah, formatDate } from '@/lib/utils'
import type { Subscription, PaginatedResponse } from '@/types'

export function SubscriptionManagementPage() {
  const [subscriptions, setSubscriptions] = useState<Subscription[]>([])
  const [loading, setLoading] = useState(true)
  const [statusFilter, setStatusFilter] = useState('all')
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [totalCount, setTotalCount] = useState(0)

  const fetchSubscriptions = async () => {
    try {
      setLoading(true)
      const res = await api.get<PaginatedResponse<Subscription>>('/admin/subscriptions', {
        params: {
          page,
          limit: 15,
          status: statusFilter !== 'all' ? statusFilter : undefined,
        },
      })
      setSubscriptions(res.data.data || [])
      setTotalPages(res.data.meta?.total_pages || 1)
      setTotalCount(res.data.meta?.total || 0)
    } catch (err: any) {
      toast.error('Gagal memuat daftar langganan')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchSubscriptions()
  }, [page, statusFilter])

  return (
    <div className="space-y-6 max-w-7xl mx-auto pb-12">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Manajemen Langganan Premium</h1>
          <p className="text-sm text-muted-foreground mt-1">
            Pantau status langganan aktif dan riwayat transaksi pembayaran paket Premium para fotografer.
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={fetchSubscriptions} disabled={loading} className="gap-2 self-start">
          <RefreshCw size={14} className={loading ? 'animate-spin' : ''} />
          Segarkan Data
        </Button>
      </div>

      {/* Filter Tabs */}
      <Card>
        <CardContent className="p-3">
          <div className="flex items-center gap-1.5 overflow-x-auto">
            {[
              { id: 'all', label: 'Semua Transaksi' },
              { id: 'paid', label: 'Aktif / Lunas (Paid)' },
              { id: 'pending', label: 'Menunggu Pembayaran (Pending)' },
              { id: 'expired', label: 'Kedaluwarsa (Expired)' },
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

      {/* Subscriptions Table */}
      <Card>
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <table className="w-full text-sm text-left">
              <thead className="bg-muted/40 border-b text-xs text-muted-foreground uppercase tracking-wider">
                <tr>
                  <th className="py-3 px-4 font-medium">Fotografer</th>
                  <th className="py-3 px-4 font-medium">Paket</th>
                  <th className="py-3 px-4 font-medium">Nominal</th>
                  <th className="py-3 px-4 font-medium">Durasi</th>
                  <th className="py-3 px-4 font-medium">Status</th>
                  <th className="py-3 px-4 font-medium">Dibayar Pada</th>
                  <th className="py-3 px-4 font-medium">Masa Aktif Hingga</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border/60">
                {loading ? (
                  <tr>
                    <td colSpan={7} className="py-12 text-center text-muted-foreground">
                      <div className="flex items-center justify-center gap-2">
                        <Loader2 size={18} className="animate-spin text-brand" />
                        <span>Memuat data langganan...</span>
                      </div>
                    </td>
                  </tr>
                ) : subscriptions.length === 0 ? (
                  <tr>
                    <td colSpan={7} className="py-12 text-center text-muted-foreground">
                      Belum ada riwayat langganan tercatat.
                    </td>
                  </tr>
                ) : (
                  subscriptions.map((sub) => (
                    <tr key={sub.id} className="hover:bg-muted/30 transition-colors">
                      <td className="py-3.5 px-4">
                        <div className="font-medium text-foreground">{sub.user?.name || 'Fotografer'}</div>
                        <div className="text-xs text-muted-foreground">{sub.user?.email || '-'}</div>
                      </td>

                      <td className="py-3.5 px-4">
                        <Badge variant="brand" className="text-xs font-medium">
                          ✦ Premium
                        </Badge>
                      </td>

                      <td className="py-3.5 px-4 font-bold text-foreground">
                        {sub.amount > 0 ? formatRupiah(sub.amount) : 'Gratis (Admin Grant)'}
                      </td>

                      <td className="py-3.5 px-4 text-xs font-medium text-muted-foreground">
                        {sub.duration_months} Bulan
                      </td>

                      <td className="py-3.5 px-4">
                        <Badge
                          variant={
                            sub.status === 'paid'
                              ? 'default'
                              : sub.status === 'pending'
                              ? 'secondary'
                              : 'outline'
                          }
                          className="text-xs font-medium"
                        >
                          {sub.status === 'paid' && '✓ Aktif (Paid)'}
                          {sub.status === 'pending' && '⏳ Pending'}
                          {sub.status === 'expired' && 'Expired'}
                        </Badge>
                      </td>

                      <td className="py-3.5 px-4 text-xs text-muted-foreground">
                        {sub.paid_at ? formatDate(sub.paid_at) : '-'}
                      </td>

                      <td className="py-3.5 px-4 text-xs font-medium text-foreground">
                        {sub.expires_at ? formatDate(sub.expires_at) : '-'}
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
                Menampilkan halaman {page} dari {totalPages} (Total {totalCount} data langganan)
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
    </div>
  )
}
export default SubscriptionManagementPage
