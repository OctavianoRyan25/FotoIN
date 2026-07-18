import { useEffect, useState } from 'react'
import { TrendingUp, TrendingDown, Wallet, ArrowRight, Loader2, Banknote } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { formatRupiah, formatDateTime } from '@/lib/utils'
import api from '@/lib/api'
import type { WalletSummary, WithdrawalRequest } from '@/types'
import { toast } from 'sonner'
import { WithdrawalDialog } from '@/components/features/wallet/WithdrawalDialog'

export function WalletPage() {
  const [walletData, setWalletData] = useState<WalletSummary | null>(null)
  const [withdrawals, setWithdrawals] = useState<WithdrawalRequest[]>([])
  const [loading, setLoading] = useState(true)
  const [showWithdrawal, setShowWithdrawal] = useState(false)

  const fetchData = async () => {
    setLoading(true)
    try {
      const res = await api.get('/wallet?limit=20')
      setWalletData(res.data.data)
    } catch {
      toast.error('Gagal memuat data dompet')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchData() }, [])

  const getTransactionBadge = (type: string) => {
    return type === 'credit'
      ? { variant: 'success' as const, label: 'Masuk', icon: TrendingUp }
      : { variant: 'destructive' as const, label: 'Keluar', icon: TrendingDown }
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Dompet</h1>
        <p className="text-muted-foreground text-sm mt-1">Kelola pendapatan dan penarikan dana Anda</p>
      </div>

      {/* Balance Card */}
      <Card className="overflow-hidden">
        <div className="p-6 text-white" style={{ background: 'linear-gradient(135deg, #1F1B16 0%, #23474A 100%)' }}>
          <div className="flex items-start justify-between">
            <div>
              <p className="text-white/70 text-sm mb-1">Saldo Tersedia</p>
              {loading ? (
                <div className="h-10 w-48 bg-white/20 rounded-lg animate-pulse" />
              ) : (
                <p className="text-4xl font-bold">
                  {formatRupiah(walletData?.balance ?? 0)}
                </p>
              )}
            </div>
            <div className="bg-white/20 p-3 rounded-xl">
              <Wallet size={24} />
            </div>
          </div>
          <div className="mt-6 flex gap-3">
            <Button
              id="request-withdrawal"
              onClick={() => setShowWithdrawal(true)}
              className="bg-white text-[#1F1B16] hover:bg-[#F6F3EC] gap-2"
              disabled={loading || (walletData?.balance ?? 0) < 50000}
            >
              <Banknote size={16} />
              Tarik Dana
            </Button>
          </div>
          {(walletData?.balance ?? 0) < 50000 && !loading && (
            <p className="text-white/60 text-xs mt-2">
              Saldo minimum penarikan: Rp 50.000
            </p>
          )}
        </div>
      </Card>

      {/* Transactions */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">Riwayat Transaksi</CardTitle>
          <CardDescription>
            {walletData?.total ?? 0} transaksi
          </CardDescription>
        </CardHeader>
        <CardContent>
          {loading ? (
            <div className="space-y-3">
              {[1, 2, 3, 4, 5].map((i) => (
                <div key={i} className="h-16 bg-muted rounded-lg animate-pulse" />
              ))}
            </div>
          ) : !walletData?.transactions?.length ? (
            <div className="text-center py-12">
              <Wallet size={36} className="mx-auto text-muted-foreground mb-3" />
              <p className="text-muted-foreground text-sm">Belum ada transaksi</p>
              <p className="text-muted-foreground text-xs mt-1">
                Transaksi akan muncul setelah pembeli membayar foto Anda
              </p>
            </div>
          ) : (
            <div className="space-y-2">
              {walletData.transactions.map((tx) => {
                const badge = getTransactionBadge(tx.type)
                return (
                  <div
                    key={tx.id}
                    className="flex items-center gap-4 p-3 rounded-lg hover:bg-muted/50 transition-colors"
                  >
                    <div className={`p-2 rounded-lg ${
                      tx.type === 'credit'
                        ? 'bg-emerald-50 dark:bg-emerald-900/20'
                        : 'bg-red-50 dark:bg-red-900/20'
                    }`}>
                      {tx.type === 'credit'
                        ? <TrendingUp size={16} className="text-emerald-600" />
                        : <TrendingDown size={16} className="text-red-600" />
                      }
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium truncate">{tx.description}</p>
                      <p className="text-xs text-muted-foreground">{formatDateTime(tx.created_at)}</p>
                    </div>
                    <div className="text-right">
                      <p className={`text-sm font-semibold ${
                        tx.type === 'credit' ? 'text-emerald-600' : 'text-red-600'
                      }`}>
                        {tx.type === 'credit' ? '+' : '-'}{formatRupiah(tx.amount)}
                      </p>
                      <p className="text-xs text-muted-foreground">
                        Saldo: {formatRupiah(tx.balance_after)}
                      </p>
                    </div>
                  </div>
                )
              })}
            </div>
          )}
        </CardContent>
      </Card>

      <WithdrawalDialog
        open={showWithdrawal}
        onOpenChange={setShowWithdrawal}
        maxAmount={walletData?.balance ?? 0}
        onSubmitted={fetchData}
      />
    </div>
  )
}
