import { useState, useEffect } from 'react'
import {
  Sparkles,
  Check,
  Zap,
  HardDrive,
  Clock,
  Percent,
  Headphones,
  ShieldCheck,
  Loader2,
  Calendar,
  AlertCircle,
  ExternalLink,
} from 'lucide-react'
import { toast } from 'sonner'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '@/components/ui/dialog'
import { useAuth } from '@/store/auth'
import api from '@/lib/api'
import { formatBytes, formatRupiah } from '@/lib/utils'
import type { SubscriptionStatusResponse, Subscription } from '@/types'

export function SubscriptionPage() {
  const { user, refreshUser } = useAuth()
  const [loading, setLoading] = useState(true)
  const [subData, setSubData] = useState<SubscriptionStatusResponse | null>(null)
  const [upgrading, setUpgrading] = useState(false)
  const [simulating, setSimulating] = useState(false)
  const [activePaymentSub, setActivePaymentSub] = useState<Subscription | null>(null)
  const [qrModalOpen, setQrModalOpen] = useState(false)

  const fetchSubscription = async () => {
    try {
      setLoading(true)
      const res = await api.get<{ data: SubscriptionStatusResponse }>('/subscriptions/me')
      setSubData(res.data.data)
    } catch (err: any) {
      toast.error('Gagal memuat status langganan')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchSubscription()
  }, [])

  const handleStartUpgrade = async () => {
    try {
      setUpgrading(true)
      const res = await api.post<{
        data: {
          subscription: Subscription
          qr_code_url: string
          qr_code_string: string
          expires_at: string
        }
      }>('/subscriptions/upgrade')
      setActivePaymentSub(res.data.data.subscription)
      setQrModalOpen(true)
    } catch (err: any) {
      toast.error(err.response?.data?.message || 'Gagal memulai upgrade')
    } finally {
      setUpgrading(false)
    }
  }

  const handleSimulatePayment = async () => {
    if (!activePaymentSub) return
    try {
      setSimulating(true)
      await api.post(`/subscriptions/${activePaymentSub.id}/simulate`)
      toast.success('🎉 Upgrade Premium berhasil diaktifkan!')
      setQrModalOpen(false)
      setActivePaymentSub(null)
      await refreshUser()
      await fetchSubscription()
    } catch (err: any) {
      toast.error(err.response?.data?.message || 'Gagal simulasi pembayaran')
    } finally {
      setSimulating(false)
    }
  }

  const isPremium = user?.subscription_tier === 'premium'
  const storageLimit = isPremium ? 107374182400 : 3221225472
  const storageUsed = user?.storage_used_bytes || 0
  const storagePercent = Math.min((storageUsed / storageLimit) * 100, 100)

  return (
    <div className="space-y-8 max-w-5xl mx-auto pb-12">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Langganan & Paket</h1>
        <p className="text-sm text-muted-foreground mt-1">
          Pilih paket yang sesuai untuk mengembangkan bisnis fotografi event Anda.
        </p>
      </div>

      {/* Current Tier Banner */}
      <Card className="border-border/60 bg-gradient-to-br from-card to-muted/20">
        <CardContent className="p-6">
          <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
            <div className="space-y-1">
              <div className="flex items-center gap-2.5">
                <span className="text-sm font-medium text-muted-foreground">Paket Aktif:</span>
                <Badge
                  variant={isPremium ? 'brand' : 'secondary'}
                  className="px-2.5 py-0.5 text-xs font-semibold"
                >
                  {isPremium ? '✦ Premium Plan' : 'Basic Plan (Gratis)'}
                </Badge>
              </div>
              <p className="text-sm text-muted-foreground mt-1">
                {isPremium ? (
                  subData?.expires_at ? (
                    <span className="flex items-center gap-1.5 text-foreground">
                      <Calendar size={14} className="text-brand" />
                      Aktif hingga:{' '}
                      <strong className="font-medium">
                        {new Date(subData.expires_at).toLocaleDateString('id-ID', {
                          day: 'numeric',
                          month: 'long',
                          year: 'numeric',
                        })}
                      </strong>
                    </span>
                  ) : (
                    'Langganan aktif tanpa batas kedaluwarsa'
                  )
                ) : (
                  'Tingkatkan ke Premium untuk mendapatkan penyimpanan 100 GB dan potongan komisi hanya 4%.'
                )}
              </p>
            </div>

            {/* Storage Progress */}
            <div className="min-w-[240px] space-y-1.5">
              <div className="flex justify-between text-xs text-muted-foreground">
                <span>Kapasitas Digunakan</span>
                <span className="font-medium text-foreground">
                  {formatBytes(storageUsed)} / {formatBytes(storageLimit)}
                </span>
              </div>
              <div className="h-2 bg-muted rounded-full overflow-hidden">
                <div
                  className={`h-full rounded-full transition-all duration-300 ${
                    storagePercent > 80 ? 'bg-destructive' : 'bg-brand'
                  }`}
                  style={{ width: `${storagePercent}%` }}
                />
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Plan Comparison Grid */}
      <div className="grid md:grid-cols-2 gap-6 items-stretch">
        {/* Basic Plan */}
        <Card className={`relative flex flex-col justify-between ${!isPremium ? 'border-foreground/30 shadow-sm' : 'border-border'}`}>
          {!isPremium && (
            <div className="absolute -top-3 left-6">
              <span className="bg-muted text-muted-foreground text-[11px] font-semibold px-2.5 py-0.5 rounded-full border">
                Paket Anda Saat Ini
              </span>
            </div>
          )}
          <CardHeader className="pt-8">
            <CardTitle className="text-xl">Basic Plan</CardTitle>
            <CardDescription>Cocok untuk fotografer pemula atau event kecil berkala.</CardDescription>
            <div className="mt-4">
              <span className="text-3xl font-extrabold">Gratis</span>
              <span className="text-muted-foreground text-sm"> / selamanya</span>
            </div>
          </CardHeader>

          <CardContent className="space-y-4 flex-1">
            <div className="space-y-3 pt-2 text-sm">
              <div className="flex items-start gap-2.5">
                <HardDrive size={16} className="text-muted-foreground mt-0.5 flex-shrink-0" />
                <span>
                  Storage limit <strong>3 GB</strong>
                </span>
              </div>
              <div className="flex items-start gap-2.5">
                <Clock size={16} className="text-muted-foreground mt-0.5 flex-shrink-0" />
                <span>
                  Retensi data: Galeri & foto tersimpan <strong>14 hari</strong>
                </span>
              </div>
              <div className="flex items-start gap-2.5">
                <Percent size={16} className="text-muted-foreground mt-0.5 flex-shrink-0" />
                <span>
                  Komisi platform: <strong>20%</strong> dari harga jual foto
                </span>
              </div>
              <div className="flex items-start gap-2.5 text-muted-foreground">
                <Check size={16} className="mt-0.5 flex-shrink-0" />
                <span>Galeri publik dengan link unik & QRIS checkout</span>
              </div>
              <div className="flex items-start gap-2.5 text-muted-foreground">
                <Check size={16} className="mt-0.5 flex-shrink-0" />
                <span>Penarikan saldo dompet ke rekening bank</span>
              </div>
            </div>
          </CardContent>

          <div className="p-6 pt-0 mt-auto">
            <Button variant="outline" className="w-full" disabled={!isPremium}>
              {!isPremium ? 'Sedang Digunakan' : 'Paket Dasar'}
            </Button>
          </div>
        </Card>

        {/* Premium Plan */}
        <Card className={`relative flex flex-col justify-between border-2 ${isPremium ? 'border-brand shadow-md' : 'border-brand/80 shadow-lg'}`}>
          <div className="absolute -top-3 left-6">
            <span className="bg-brand text-brand-foreground text-[11px] font-semibold px-3 py-0.5 rounded-full flex items-center gap-1 shadow-sm">
              <Sparkles size={12} />
              {isPremium ? 'Langganan Aktif' : 'Sangat Direkomendasikan'}
            </span>
          </div>

          <CardHeader className="pt-8">
            <CardTitle className="text-xl flex items-center justify-between">
              <span>Premium Plan</span>
              <Zap size={20} className="text-brand" />
            </CardTitle>
            <CardDescription>Untuk fotografer profesional & pengelola event fotografi intensif.</CardDescription>
            <div className="mt-4">
              <span className="text-3xl font-extrabold">{formatRupiah(149000)}</span>
              <span className="text-muted-foreground text-sm"> / bulan</span>
            </div>
          </CardHeader>

          <CardContent className="space-y-4 flex-1">
            <div className="space-y-3 pt-2 text-sm">
              <div className="flex items-start gap-2.5">
                <HardDrive size={16} className="text-brand mt-0.5 flex-shrink-0" />
                <span>
                  Storage limit <strong>100 GB</strong> (33x lipat lebih besar!)
                </span>
              </div>
              <div className="flex items-start gap-2.5">
                <ShieldCheck size={16} className="text-brand mt-0.5 flex-shrink-0" />
                <span>
                  Retensi data: <strong>Permanen</strong> selama masa langganan aktif
                </span>
              </div>
              <div className="flex items-start gap-2.5">
                <Percent size={16} className="text-brand mt-0.5 flex-shrink-0" />
                <span>
                  Komisi platform <strong>hanya 4%</strong> (hemat 16% setiap penjualan foto!)
                </span>
              </div>
              <div className="flex items-start gap-2.5">
                <Check size={16} className="text-brand mt-0.5 flex-shrink-0" />
                <span>Unmetered bandwidth limit untuk download berkecepatan tinggi</span>
              </div>
              <div className="flex items-start gap-2.5">
                <Headphones size={16} className="text-brand mt-0.5 flex-shrink-0" />
                <span>Prioritas support WhatsApp & Email</span>
              </div>
            </div>
          </CardContent>

          <div className="p-6 pt-0 mt-auto">
            <Button
              className="w-full bg-brand text-white hover:bg-brand/90 font-medium"
              onClick={handleStartUpgrade}
              disabled={upgrading}
            >
              {upgrading ? (
                <>
                  <Loader2 size={16} className="animate-spin mr-2" />
                  Menyiapkan Pembayaran...
                </>
              ) : isPremium ? (
                'Perpanjang Langganan (+1 Bulan)'
              ) : (
                'Upgrade ke Premium Sekarang'
              )}
            </Button>
          </div>
        </Card>
      </div>

      {/* Upgrade Payment Dialog */}
      <Dialog open={qrModalOpen} onOpenChange={setQrModalOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <Sparkles size={18} className="text-brand" />
              Upgrade ke FotoIN Premium
            </DialogTitle>
            <DialogDescription>
              Scan QRIS di bawah ini menggunakan aplikasi e-wallet atau mobile banking favorit Anda.
            </DialogDescription>
          </DialogHeader>

          <div className="flex flex-col items-center py-4 space-y-4">
            {/* Amount display */}
            <div className="text-center">
              <div className="text-xs text-muted-foreground uppercase tracking-wider font-medium">Total Pembayaran</div>
              <div className="text-3xl font-bold tracking-tight text-foreground mt-0.5">{formatRupiah(149000)}</div>
              <div className="text-xs text-muted-foreground mt-1">Masa aktif: 1 Bulan (30 Hari)</div>
            </div>

            {/* QR Code display */}
            <div className="p-3 bg-white rounded-xl border-2 border-border shadow-sm flex flex-col items-center justify-center">
              <img
                src={
                  activePaymentSub?.qr_code_url ||
                  `https://api.qrserver.com/v1/create-qr-code/?size=220x220&data=${encodeURIComponent(
                    activePaymentSub?.qr_code_string || 'FOTOIN-PREMIUM-UPGRADE'
                  )}`
                }
                alt="QRIS Upgrade FotoIN"
                className="w-52 h-52 object-contain"
              />
              <div className="text-[10px] text-zinc-500 font-semibold tracking-wider uppercase mt-2">
                QRIS Standard Pembayaran Nasional
              </div>
            </div>

            <div className="w-full bg-muted/40 p-3 rounded-lg text-xs space-y-1.5 text-muted-foreground">
              <div className="flex items-center gap-1.5 text-foreground font-medium">
                <AlertCircle size={14} className="text-brand" />
                Petunjuk Pembayaran:
              </div>
              <p>1. Buka GoPay, OVO, DANA, BCA Mobile, atau Livin Mandiri.</p>
              <p>2. Arahkan kamera pemindai ke kode QRIS di atas.</p>
              <p>3. Konfirmasi pembayaran sebesar Rp 149.000.</p>
            </div>

            {/* Simulation Button for Testing/Demo */}
            <div className="w-full pt-2">
              <Button
                variant="default"
                className="w-full bg-emerald-600 hover:bg-emerald-700 text-white font-medium"
                onClick={handleSimulatePayment}
                disabled={simulating}
              >
                {simulating ? (
                  <>
                    <Loader2 size={16} className="animate-spin mr-2" />
                    Memproses Simulasi...
                  </>
                ) : (
                  '⚡ Simulasikan Pembayaran Sukses (Testing)'
                )}
              </Button>
              <p className="text-[11px] text-center text-muted-foreground mt-1.5">
                Tombol simulasi disediakan untuk pengujian instan alur transaksi.
              </p>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  )
}
export default SubscriptionPage
