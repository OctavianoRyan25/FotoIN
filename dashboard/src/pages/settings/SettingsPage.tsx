import { useState } from 'react'
import {
  User as UserIcon,
  Lock,
  Shield,
  HardDrive,
  Save,
  Loader2,
  CheckCircle2,
} from 'lucide-react'
import { toast } from 'sonner'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { useAuth } from '@/store/auth'
import api from '@/lib/api'
import { formatBytes } from '@/lib/utils'

export function SettingsPage() {
  const { user, refreshUser } = useAuth()
  const [name, setName] = useState(user?.name || '')
  const [profileSaving, setProfileSaving] = useState(false)

  // Password state
  const [oldPassword, setOldPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [passwordSaving, setPasswordSaving] = useState(false)

  const handleSaveProfile = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!name.trim()) {
      toast.error('Nama lengkap tidak boleh kosong')
      return
    }

    try {
      setProfileSaving(true)
      await api.put('/me/profile', { name })
      toast.success('Profil berhasil diperbarui')
      await refreshUser()
    } catch (err: any) {
      toast.error(err.response?.data?.message || 'Gagal memperbarui profil')
    } finally {
      setProfileSaving(false)
    }
  }

  const handleUpdatePassword = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!newPassword || newPassword.length < 6) {
      toast.error('Password baru minimal 6 karakter')
      return
    }
    if (newPassword !== confirmPassword) {
      toast.error('Konfirmasi password tidak cocok')
      return
    }

    try {
      setPasswordSaving(true)
      await api.put('/me/password', {
        old_password: oldPassword,
        new_password: newPassword,
      })
      toast.success('Password berhasil diperbarui')
      setOldPassword('')
      setNewPassword('')
      setConfirmPassword('')
    } catch (err: any) {
      toast.error(err.response?.data?.message || 'Gagal memperbarui password')
    } finally {
      setPasswordSaving(false)
    }
  }

  const isPremium = user?.subscription_tier === 'premium'
  const storageLimit = isPremium ? 107374182400 : 3221225472
  const storageUsed = user?.storage_used_bytes || 0
  const storagePercent = Math.min((storageUsed / storageLimit) * 100, 100)

  return (
    <div className="space-y-6 max-w-4xl mx-auto pb-12">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Pengaturan Akun</h1>
        <p className="text-sm text-muted-foreground mt-1">
          Kelola informasi profil, preferensi, dan keamanan akun Anda.
        </p>
      </div>

      <div className="grid gap-6">
        {/* Profile Information */}
        <Card>
          <CardHeader>
            <CardTitle className="text-lg flex items-center gap-2">
              <UserIcon size={18} className="text-brand" />
              Informasi Pribadi
            </CardTitle>
            <CardDescription>Perbarui nama akun Anda yang ditampilkan ke publik.</CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSaveProfile} className="space-y-4 max-w-md">
              <div className="space-y-2">
                <Label htmlFor="email">Email</Label>
                <Input id="email" value={user?.email || ''} disabled className="bg-muted text-muted-foreground" />
                <p className="text-[11px] text-muted-foreground">Email akun digunakan untuk login dan notifikasi.</p>
              </div>

              <div className="space-y-2">
                <Label htmlFor="name">Nama Lengkap</Label>
                <Input
                  id="name"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="Masukkan nama Anda"
                  required
                />
              </div>

              <Button type="submit" disabled={profileSaving} className="bg-brand text-white hover:bg-brand/90">
                {profileSaving ? (
                  <>
                    <Loader2 size={16} className="animate-spin mr-2" />
                    Menyimpan...
                  </>
                ) : (
                  <>
                    <Save size={16} className="mr-2" />
                    Simpan Perubahan
                  </>
                )}
              </Button>
            </form>
          </CardContent>
        </Card>

        {/* Change Password */}
        <Card>
          <CardHeader>
            <CardTitle className="text-lg flex items-center gap-2">
              <Lock size={18} className="text-brand" />
              Keamanan & Password
            </CardTitle>
            <CardDescription>Pastikan akun Anda terlindungi dengan password yang kuat.</CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleUpdatePassword} className="space-y-4 max-w-md">
              <div className="space-y-2">
                <Label htmlFor="old-password">Password Saat Ini</Label>
                <Input
                  id="old-password"
                  type="password"
                  value={oldPassword}
                  onChange={(e) => setOldPassword(e.target.value)}
                  placeholder="••••••••"
                  required
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="new-password">Password Baru</Label>
                <Input
                  id="new-password"
                  type="password"
                  value={newPassword}
                  onChange={(e) => setNewPassword(e.target.value)}
                  placeholder="Minimal 6 karakter"
                  required
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="confirm-password">Konfirmasi Password Baru</Label>
                <Input
                  id="confirm-password"
                  type="password"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  placeholder="Ulangi password baru"
                  required
                />
              </div>

              <Button type="submit" disabled={passwordSaving} variant="outline">
                {passwordSaving ? (
                  <>
                    <Loader2 size={16} className="animate-spin mr-2" />
                    Mengubah...
                  </>
                ) : (
                  'Perbarui Password'
                )}
              </Button>
            </form>
          </CardContent>
        </Card>

        {/* Account Details & Limits */}
        <Card>
          <CardHeader>
            <CardTitle className="text-lg flex items-center gap-2">
              <Shield size={18} className="text-brand" />
              Detail Akun & Kuota
            </CardTitle>
            <CardDescription>Status tingkatan paket langganan dan kapasitas penyimpanan.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid sm:grid-cols-2 gap-4">
              <div className="p-4 rounded-lg bg-muted/30 border border-border/50 space-y-1">
                <div className="text-xs text-muted-foreground">Tingkatan Akun</div>
                <div className="flex items-center gap-2 mt-1">
                  <Badge variant={isPremium ? 'brand' : 'secondary'}>
                    {isPremium ? '✦ Premium' : 'Basic (Gratis)'}
                  </Badge>
                  {isPremium && (
                    <span className="text-xs text-muted-foreground flex items-center gap-1">
                      <CheckCircle2 size={13} className="text-brand" /> Aktif
                    </span>
                  )}
                </div>
              </div>

              <div className="p-4 rounded-lg bg-muted/30 border border-border/50 space-y-1">
                <div className="text-xs text-muted-foreground">Komisi Platform</div>
                <div className="text-lg font-bold">{isPremium ? '4%' : '20%'}</div>
                <div className="text-[11px] text-muted-foreground">
                  {isPremium ? 'Tier Premium: komisi terendah' : 'Upgrade ke Premium untuk komisi 4%'}
                </div>
              </div>
            </div>

            <div className="p-4 rounded-lg bg-muted/30 border border-border/50 space-y-2">
              <div className="flex justify-between text-xs text-muted-foreground">
                <span className="flex items-center gap-1.5 font-medium text-foreground">
                  <HardDrive size={14} className="text-brand" />
                  Kapasitas Penyimpanan S3
                </span>
                <span>
                  {formatBytes(storageUsed)} dari {formatBytes(storageLimit)} ({storagePercent.toFixed(1)}%)
                </span>
              </div>
              <div className="h-2 bg-muted rounded-full overflow-hidden">
                <div
                  className={`h-full rounded-full transition-all ${
                    storagePercent > 80 ? 'bg-destructive' : 'bg-brand'
                  }`}
                  style={{ width: `${storagePercent}%` }}
                />
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
export default SettingsPage
