import { useState, useEffect } from 'react'
import {
  TrendingUp,
  DollarSign,
  ShoppingBag,
  Users,
  Sparkles,
  Percent,
  RefreshCw,
  Loader2,
  PieChart as PieChartIcon,
  ShieldCheck,
} from 'lucide-react'
import { toast } from 'sonner'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import api from '@/lib/api'
import { formatRupiah } from '@/lib/utils'
import type { RevenueOverview } from '@/types'

export function RevenuePage() {
  const [data, setData] = useState<RevenueOverview | null>(null)
  const [loading, setLoading] = useState(true)

  const fetchRevenue = async () => {
    try {
      setLoading(true)
      const res = await api.get<{ data: RevenueOverview }>('/admin/revenue')
      setData(res.data.data)
    } catch (err: any) {
      toast.error('Gagal memuat ringkasan pendapatan')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchRevenue()
  }, [])

  const totalRevenue = data?.total_revenue || 0
  const basicRev = data?.basic_tier_revenue || 0
  const premiumRev = data?.premium_tier_revenue || 0
  const subRev = data?.subscription_revenue || 0
  const orderComm = data?.order_commission || 0

  const basicPercent = orderComm > 0 ? ((basicRev / orderComm) * 100).toFixed(1) : '0'
  const premiumPercent = orderComm > 0 ? ((premiumRev / orderComm) * 100).toFixed(1) : '0'

  return (
    <div className="space-y-6 max-w-7xl mx-auto pb-12">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Platform Revenue Overview</h1>
          <p className="text-sm text-muted-foreground mt-1">
            Ringkasan pendapatan platform, komisi penjualan foto, dan pemasukan langganan Premium.
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={fetchRevenue} disabled={loading} className="gap-2 self-start">
          <RefreshCw size={14} className={loading ? 'animate-spin' : ''} />
          Segarkan Data
        </Button>
      </div>

      {/* KPI Cards */}
      <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* Total Platform Revenue */}
        <Card className="border-brand/40 bg-gradient-to-br from-card to-brand/5">
          <CardHeader className="pb-2">
            <div className="flex items-center justify-between">
              <span className="text-xs font-medium text-muted-foreground">Total Pendapatan Bersih Platform</span>
              <div className="w-8 h-8 rounded-full bg-brand/15 text-brand flex items-center justify-center">
                <DollarSign size={16} />
              </div>
            </div>
            <CardTitle className="text-2xl font-extrabold mt-1 text-foreground">
              {loading ? <Loader2 size={20} className="animate-spin text-brand" /> : formatRupiah(totalRevenue)}
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <p className="text-xs text-muted-foreground">Komisi penjualan foto + Langganan Premium</p>
          </CardContent>
        </Card>

        {/* Order Commission */}
        <Card>
          <CardHeader className="pb-2">
            <div className="flex items-center justify-between">
              <span className="text-xs font-medium text-muted-foreground">Total Komisi Penjualan Foto</span>
              <div className="w-8 h-8 rounded-full bg-emerald-500/15 text-emerald-600 flex items-center justify-center">
                <Percent size={16} />
              </div>
            </div>
            <CardTitle className="text-2xl font-bold mt-1 text-foreground">
              {loading ? <Loader2 size={20} className="animate-spin" /> : formatRupiah(orderComm)}
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <p className="text-xs text-muted-foreground">Dari {data?.total_orders || 0} order foto berhasil</p>
          </CardContent>
        </Card>

        {/* Subscription Revenue */}
        <Card>
          <CardHeader className="pb-2">
            <div className="flex items-center justify-between">
              <span className="text-xs font-medium text-muted-foreground">Pendapatan Langganan Premium</span>
              <div className="w-8 h-8 rounded-full bg-violet-500/15 text-violet-600 flex items-center justify-center">
                <Sparkles size={16} />
              </div>
            </div>
            <CardTitle className="text-2xl font-bold mt-1 text-foreground">
              {loading ? <Loader2 size={20} className="animate-spin" /> : formatRupiah(subRev)}
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <p className="text-xs text-muted-foreground">{data?.active_subscribers || 0} fotografer Premium aktif</p>
          </CardContent>
        </Card>

        {/* Total Photographers */}
        <Card>
          <CardHeader className="pb-2">
            <div className="flex items-center justify-between">
              <span className="text-xs font-medium text-muted-foreground">Total Fotografer Terdaftar</span>
              <div className="w-8 h-8 rounded-full bg-blue-500/15 text-blue-600 flex items-center justify-center">
                <Users size={16} />
              </div>
            </div>
            <CardTitle className="text-2xl font-bold mt-1 text-foreground">
              {loading ? <Loader2 size={20} className="animate-spin" /> : data?.total_photographers || 0}
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <p className="text-xs text-muted-foreground">Mitra kreator aktif di platform</p>
          </CardContent>
        </Card>
      </div>

      {/* Tier Revenue Breakdown (Per Tier Analysis) */}
      <div className="grid md:grid-cols-2 gap-6">
        {/* Basic Tier (20%) */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="text-lg">Komisi Basic Plan (20%)</CardTitle>
                <CardDescription>Pendapatan platform dari fotografer tier gratis (Basic)</CardDescription>
              </div>
              <Badge variant="secondary">20% Fee</Badge>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="p-4 rounded-xl bg-muted/30 border space-y-2">
              <div className="text-xs text-muted-foreground">Total Komisi Terkumpul:</div>
              <div className="text-3xl font-extrabold text-foreground">{formatRupiah(basicRev)}</div>
              <div className="text-xs text-muted-foreground">
                Kontribusi <strong>{basicPercent}%</strong> dari seluruh komisi order foto.
              </div>
            </div>

            <div className="text-xs text-muted-foreground space-y-1.5 pt-1">
              <p>• Fotografer Basic dipotong komisi 20% per foto terjual.</p>
              <p>• Kuota penyimpanan 3 GB dengan masa retensi galeri 14 hari.</p>
              <p>• Margin platform tinggi namun kapasitas S3 tetap terbatas.</p>
            </div>
          </CardContent>
        </Card>

        {/* Premium Tier (4%) */}
        <Card className="border-brand/30">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="text-lg flex items-center gap-2">
                  <Sparkles size={18} className="text-brand" />
                  Komisi Premium Plan (4%)
                </CardTitle>
                <CardDescription>Pendapatan komisi transaksi dari fotografer Premium</CardDescription>
              </div>
              <Badge variant="brand">4% Fee</Badge>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="p-4 rounded-xl bg-brand/5 border border-brand/20 space-y-2">
              <div className="text-xs text-muted-foreground">Total Komisi Terkumpul:</div>
              <div className="text-3xl font-extrabold text-brand">{formatRupiah(premiumRev)}</div>
              <div className="text-xs text-muted-foreground">
                Kontribusi <strong>{premiumPercent}%</strong> dari seluruh komisi order foto.
              </div>
            </div>

            <div className="text-xs text-muted-foreground space-y-1.5 pt-1">
              <p>• Komisi 4% mencakup ~2% payment gateway fee + ~2% platform margin.</p>
              <p>• Tambahan pendapatan tetap dari biaya langganan bulanan Rp 149.000.</p>
              <p>• Retensi data permanen dan penyimpanan 100 GB.</p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
export default RevenuePage
