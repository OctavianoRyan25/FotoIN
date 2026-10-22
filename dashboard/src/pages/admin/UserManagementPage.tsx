import { useState, useEffect } from 'react'
import {
  Users,
  Search,
  ShieldAlert,
  ShieldCheck,
  Sparkles,
  Loader2,
  RefreshCw,
  MoreVertical,
  Filter,
} from 'lucide-react'
import { toast } from 'sonner'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '@/components/ui/dialog'
import api from '@/lib/api'
import { formatBytes, formatRupiah, formatDate } from '@/lib/utils'
import type { User, PaginatedResponse } from '@/types'

export function UserManagementPage() {
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')
  const [roleFilter, setRoleFilter] = useState('all')
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [totalCount, setTotalCount] = useState(0)

  // Suspend action modal
  const [selectedUser, setSelectedUser] = useState<User | null>(null)
  const [suspendModalOpen, setSuspendModalOpen] = useState(false)
  const [suspendLoading, setSuspendLoading] = useState(false)

  // Grant premium modal
  const [grantModalOpen, setGrantModalOpen] = useState(false)
  const [grantMonths, setGrantMonths] = useState(1)
  const [grantLoading, setGrantLoading] = useState(false)

  const fetchUsers = async () => {
    try {
      setLoading(true)
      const res = await api.get<PaginatedResponse<User>>('/admin/users', {
        params: {
          page,
          limit: 15,
          search: search || undefined,
          role: roleFilter !== 'all' ? roleFilter : undefined,
        },
      })
      setUsers(res.data.data || [])
      setTotalPages(res.data.meta?.total_pages || 1)
      setTotalCount(res.data.meta?.total || 0)
    } catch (err: any) {
      toast.error('Gagal memuat daftar pengguna')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchUsers()
  }, [page, roleFilter])

  const handleSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setPage(1)
    fetchUsers()
  }

  const handleToggleSuspend = async () => {
    if (!selectedUser) return
    try {
      setSuspendLoading(true)
      const nextSuspendedState = !selectedUser.is_suspended
      await api.patch(`/admin/users/${selectedUser.id}/suspend`, {
        suspended: nextSuspendedState,
      })
      toast.success(
        nextSuspendedState
          ? `Akun ${selectedUser.name} berhasil ditangguhkan`
          : `Akun ${selectedUser.name} berhasil diaktifkan kembali`
      )
      setSuspendModalOpen(false)
      setSelectedUser(null)
      fetchUsers()
    } catch (err: any) {
      toast.error(err.response?.data?.message || 'Gagal mengubah status penangguhan akun')
    } finally {
      setSuspendLoading(false)
    }
  }

  const handleGrantPremium = async () => {
    if (!selectedUser) return
    try {
      setGrantLoading(true)
      await api.post('/admin/subscriptions/grant', {
        user_id: selectedUser.id,
        months: grantMonths,
      })
      toast.success(`Berhasil memberikan paket Premium (${grantMonths} bulan) untuk ${selectedUser.name}`)
      setGrantModalOpen(false)
      setSelectedUser(null)
      fetchUsers()
    } catch (err: any) {
      toast.error(err.response?.data?.message || 'Gagal memberikan paket Premium')
    } finally {
      setGrantLoading(false)
    }
  }

  return (
    <div className="space-y-6 max-w-7xl mx-auto pb-12">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Manajemen Pengguna</h1>
          <p className="text-sm text-muted-foreground mt-1">
            Pantau dan kelola akun fotografer, pembeli, dan administrator platform.
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={fetchUsers} disabled={loading} className="gap-2 self-start">
          <RefreshCw size={14} className={loading ? 'animate-spin' : ''} />
          Segarkan Data
        </Button>
      </div>

      {/* Filter and Search Bar */}
      <Card>
        <CardContent className="p-4">
          <div className="flex flex-col md:flex-row gap-3 items-center justify-between">
            {/* Role Filter Tabs */}
            <div className="flex items-center gap-1.5 overflow-x-auto w-full md:w-auto pb-2 md:pb-0">
              {[
                { id: 'all', label: 'Semua Akun' },
                { id: 'photographer', label: 'Fotografer' },
                { id: 'client', label: 'Pembeli (Klien)' },
                { id: 'super_admin', label: 'Super Admin' },
              ].map((tab) => (
                <Button
                  key={tab.id}
                  variant={roleFilter === tab.id ? 'default' : 'ghost'}
                  size="sm"
                  onClick={() => {
                    setRoleFilter(tab.id)
                    setPage(1)
                  }}
                  className={roleFilter === tab.id ? 'bg-brand text-white' : 'text-muted-foreground'}
                >
                  {tab.label}
                </Button>
              ))}
            </div>

            {/* Search Input */}
            <form onSubmit={handleSearchSubmit} className="flex items-center gap-2 w-full md:w-80">
              <div className="relative flex-1">
                <Search size={15} className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" />
                <Input
                  placeholder="Cari nama atau email..."
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  className="pl-9 text-sm"
                />
              </div>
              <Button type="submit" size="sm" variant="secondary">
                Cari
              </Button>
            </form>
          </div>
        </CardContent>
      </Card>

      {/* Users DataTable */}
      <Card>
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <table className="w-full text-sm text-left">
              <thead className="bg-muted/40 border-b text-xs text-muted-foreground uppercase tracking-wider">
                <tr>
                  <th className="py-3 px-4 font-medium">Pengguna</th>
                  <th className="py-3 px-4 font-medium">Peran</th>
                  <th className="py-3 px-4 font-medium">Paket Langganan</th>
                  <th className="py-3 px-4 font-medium">Storage S3</th>
                  <th className="py-3 px-4 font-medium">Saldo Dompet</th>
                  <th className="py-3 px-4 font-medium">Status</th>
                  <th className="py-3 px-4 font-medium">Terdaftar</th>
                  <th className="py-3 px-4 font-medium text-right">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border/60">
                {loading ? (
                  <tr>
                    <td colSpan={8} className="py-12 text-center text-muted-foreground">
                      <div className="flex items-center justify-center gap-2">
                        <Loader2 size={18} className="animate-spin text-brand" />
                        <span>Memuat data pengguna...</span>
                      </div>
                    </td>
                  </tr>
                ) : users.length === 0 ? (
                  <tr>
                    <td colSpan={8} className="py-12 text-center text-muted-foreground">
                      Tidak ada pengguna ditemukan.
                    </td>
                  </tr>
                ) : (
                  users.map((u) => {
                    const isPrem = u.subscription_tier === 'premium'
                    return (
                      <tr key={u.id} className="hover:bg-muted/30 transition-colors">
                        <td className="py-3 px-4">
                          <div className="flex items-center gap-3">
                            <div className="w-8 h-8 rounded-full bg-brand/15 text-brand flex items-center justify-center font-bold text-xs uppercase">
                              {u.name.charAt(0)}
                            </div>
                            <div>
                              <div className="font-medium text-foreground">{u.name}</div>
                              <div className="text-xs text-muted-foreground">{u.email}</div>
                            </div>
                          </div>
                        </td>

                        <td className="py-3 px-4">
                          <Badge variant="outline" className="capitalize text-xs">
                            {u.role === 'super_admin' ? 'Super Admin' : u.role === 'client' ? 'Klien' : 'Fotografer'}
                          </Badge>
                        </td>

                        <td className="py-3 px-4">
                          <Badge
                            variant={isPrem ? 'brand' : 'secondary'}
                            className="text-xs font-medium"
                          >
                            {isPrem ? '✦ Premium' : 'Basic'}
                          </Badge>
                        </td>

                        <td className="py-3 px-4 text-xs font-mono text-muted-foreground">
                          {formatBytes(u.storage_used_bytes || 0)}
                        </td>

                        <td className="py-3 px-4 text-xs font-medium">
                          {formatRupiah(u.wallet_balance || 0)}
                        </td>

                        <td className="py-3 px-4">
                          <Badge
                            variant={u.is_suspended ? 'destructive' : 'outline'}
                            className="text-xs"
                          >
                            {u.is_suspended ? 'Ditangguhkan' : 'Aktif'}
                          </Badge>
                        </td>

                        <td className="py-3 px-4 text-xs text-muted-foreground">
                          {formatDate(u.created_at)}
                        </td>

                        <td className="py-3 px-4 text-right">
                          <div className="flex items-center justify-end gap-1.5">
                            {/* Grant Premium Button */}
                            {u.role === 'photographer' && (
                              <Button
                                size="sm"
                                variant="ghost"
                                className="h-8 px-2 text-xs text-brand hover:text-brand-dark"
                                title="Beri Akses Premium"
                                onClick={() => {
                                  setSelectedUser(u)
                                  setGrantModalOpen(true)
                                }}
                              >
                                <Sparkles size={14} className="mr-1" />
                                Premium
                              </Button>
                            )}

                            {/* Suspend/Unsuspend Button */}
                            <Button
                              size="sm"
                              variant="ghost"
                              className={`h-8 px-2 text-xs ${
                                u.is_suspended ? 'text-emerald-600 hover:text-emerald-700' : 'text-destructive hover:text-destructive'
                              }`}
                              title={u.is_suspended ? 'Aktifkan Akun' : 'Tangguhkan Akun'}
                              onClick={() => {
                                setSelectedUser(u)
                                setSuspendModalOpen(true)
                              }}
                            >
                              {u.is_suspended ? (
                                <>
                                  <ShieldCheck size={14} className="mr-1" />
                                  Aktifkan
                                </>
                              ) : (
                                <>
                                  <ShieldAlert size={14} className="mr-1" />
                                  Suspend
                                </>
                              )}
                            </Button>
                          </div>
                        </td>
                      </tr>
                    )
                  })
                )}
              </tbody>
            </table>
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex items-center justify-between p-4 border-t text-xs text-muted-foreground">
              <div>
                Menampilkan halaman {page} dari {totalPages} (Total {totalCount} pengguna)
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

      {/* Suspend / Unsuspend Confirmation Dialog */}
      <Dialog open={suspendModalOpen} onOpenChange={setSuspendModalOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <ShieldAlert size={18} className={selectedUser?.is_suspended ? 'text-emerald-600' : 'text-destructive'} />
              {selectedUser?.is_suspended ? 'Aktifkan Kembali Pengguna' : 'Tangguhkan Pengguna'}
            </DialogTitle>
            <DialogDescription>
              {selectedUser?.is_suspended
                ? `Apakah Anda yakin ingin membatalkan penangguhan akun "${selectedUser?.name}" (${selectedUser?.email})? Pengguna akan dapat login kembali.`
                : `Apakah Anda yakin ingin menangguhkan akun "${selectedUser?.name}" (${selectedUser?.email})? Pengguna tidak akan dapat mengakses sistem selama ditangguhkan.`}
            </DialogDescription>
          </DialogHeader>

          <div className="flex justify-end gap-2 pt-4">
            <Button variant="outline" onClick={() => setSuspendModalOpen(false)} disabled={suspendLoading}>
              Batal
            </Button>
            <Button
              variant={selectedUser?.is_suspended ? 'default' : 'destructive'}
              onClick={handleToggleSuspend}
              disabled={suspendLoading}
            >
              {suspendLoading ? (
                <>
                  <Loader2 size={14} className="animate-spin mr-1.5" /> Memproses...
                </>
              ) : selectedUser?.is_suspended ? (
                'Aktifkan Akun'
              ) : (
                'Tangguhkan Akun'
              )}
            </Button>
          </div>
        </DialogContent>
      </Dialog>

      {/* Grant Premium Dialog */}
      <Dialog open={grantModalOpen} onOpenChange={setGrantModalOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <Sparkles size={18} className="text-brand" />
              Beri Akses Paket Premium
            </DialogTitle>
            <DialogDescription>
              Tingkatkan status akun fotografer <strong>{selectedUser?.name}</strong> ke paket Premium (100 GB Storage, 4% Komisi).
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-2">
            <div>
              <label className="text-xs font-medium text-foreground">Durasi Langganan:</label>
              <div className="grid grid-cols-4 gap-2 mt-2">
                {[1, 3, 6, 12].map((m) => (
                  <Button
                    key={m}
                    type="button"
                    variant={grantMonths === m ? 'default' : 'outline'}
                    size="sm"
                    className={grantMonths === m ? 'bg-brand text-white' : ''}
                    onClick={() => setGrantMonths(m)}
                  >
                    {m} Bulan
                  </Button>
                ))}
              </div>
            </div>

            <div className="p-3 bg-muted/40 rounded-lg text-xs text-muted-foreground space-y-1">
              <p>• Kuota penyimpanan akan ditingkatkan menjadi 100 GB.</p>
              <p>• Komisi penjualan foto akan otomatis dipotong hanya 4%.</p>
              <p>• Galeri tidak akan kedaluwarsa selama durasi aktif.</p>
            </div>
          </div>

          <div className="flex justify-end gap-2 pt-2">
            <Button variant="outline" onClick={() => setGrantModalOpen(false)} disabled={grantLoading}>
              Batal
            </Button>
            <Button
              className="bg-brand text-white hover:bg-brand/90"
              onClick={handleGrantPremium}
              disabled={grantLoading}
            >
              {grantLoading ? (
                <>
                  <Loader2 size={14} className="animate-spin mr-1.5" /> Menyimpan...
                </>
              ) : (
                'Konfirmasi Berikan Premium'
              )}
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  )
}
export default UserManagementPage
