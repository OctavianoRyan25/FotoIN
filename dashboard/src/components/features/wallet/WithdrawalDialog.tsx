import { useState } from 'react'
import * as Dialog from '@radix-ui/react-dialog'
import { X, Loader2, Banknote } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { formatRupiah } from '@/lib/utils'
import api from '@/lib/api'
import { toast } from 'sonner'

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
  maxAmount: number
  onSubmitted: () => void
}

const BANK_LIST = ['BCA', 'BNI', 'BRI', 'Mandiri', 'CIMB Niaga', 'Permata', 'Danamon', 'GoPay', 'OVO', 'DANA']

export function WithdrawalDialog({ open, onOpenChange, maxAmount, onSubmitted }: Props) {
  const [amount, setAmount] = useState('')
  const [bankName, setBankName] = useState('')
  const [accountNumber, setAccountNumber] = useState('')
  const [accountHolder, setAccountHolder] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const amountNum = parseFloat(amount)
    if (amountNum < 50000) {
      toast.error('Minimum penarikan Rp 50.000')
      return
    }
    if (amountNum > maxAmount) {
      toast.error(`Saldo tidak cukup (maks: ${formatRupiah(maxAmount)})`)
      return
    }

    setLoading(true)
    try {
      await api.post('/wallet/withdrawal', {
        amount: amountNum,
        bank_name: bankName,
        account_number: accountNumber,
        account_holder: accountHolder,
      })
      toast.success('Permintaan penarikan berhasil diajukan! Tim kami akan memprosesnya.')
      onOpenChange(false)
      onSubmitted()
      setAmount(''); setBankName(''); setAccountNumber(''); setAccountHolder('')
    } catch (err: unknown) {
      const message = (err as { response?: { data?: { message?: string } } })?.response?.data?.message ?? 'Gagal mengajukan penarikan'
      toast.error(message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <Dialog.Root open={open} onOpenChange={onOpenChange}>
      <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50" />
        <Dialog.Content className="fixed left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 z-50 w-full max-w-md bg-background rounded-xl shadow-xl p-6">
          <div className="flex items-center justify-between mb-6">
            <Dialog.Title className="text-lg font-semibold flex items-center gap-2">
              <Banknote size={20} className="text-violet-600" />
              Tarik Dana
            </Dialog.Title>
            <Dialog.Close asChild>
              <Button variant="ghost" size="icon" className="h-8 w-8"><X size={16} /></Button>
            </Dialog.Close>
          </div>

          <div className="bg-muted/50 rounded-lg p-3 mb-4">
            <p className="text-xs text-muted-foreground">Saldo tersedia</p>
            <p className="text-lg font-bold text-violet-600">{formatRupiah(maxAmount)}</p>
          </div>

          <form onSubmit={handleSubmit} id="withdrawal-form" className="space-y-4">
            <div>
              <label className="text-sm font-medium mb-1.5 block">Jumlah Penarikan *</label>
              <div className="relative">
                <span className="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-muted-foreground">Rp</span>
                <Input
                  id="withdrawal-amount"
                  type="number"
                  min={50000}
                  max={maxAmount}
                  step={10000}
                  placeholder="50000"
                  value={amount}
                  onChange={(e) => setAmount(e.target.value)}
                  required
                  className="pl-9"
                />
              </div>
              <p className="text-xs text-muted-foreground mt-1">Minimum Rp 50.000</p>
            </div>

            <div>
              <label className="text-sm font-medium mb-1.5 block">Bank / E-Wallet *</label>
              <select
                id="withdrawal-bank"
                value={bankName}
                onChange={(e) => setBankName(e.target.value)}
                required
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
              >
                <option value="">Pilih bank...</option>
                {BANK_LIST.map((bank) => (
                  <option key={bank} value={bank}>{bank}</option>
                ))}
              </select>
            </div>

            <div>
              <label className="text-sm font-medium mb-1.5 block">Nomor Rekening *</label>
              <Input
                id="withdrawal-account-number"
                placeholder="1234567890"
                value={accountNumber}
                onChange={(e) => setAccountNumber(e.target.value)}
                required
              />
            </div>

            <div>
              <label className="text-sm font-medium mb-1.5 block">Nama Pemilik Rekening *</label>
              <Input
                id="withdrawal-account-holder"
                placeholder="Sesuai nama di buku tabungan"
                value={accountHolder}
                onChange={(e) => setAccountHolder(e.target.value)}
                required
              />
            </div>

            <div className="bg-amber-50 dark:bg-amber-900/20 rounded-lg p-3">
              <p className="text-xs text-amber-800 dark:text-amber-400">
                ⚠️ Penarikan akan diproses manual dalam 1-3 hari kerja. Pastikan data rekening benar.
              </p>
            </div>

            <div className="flex gap-3 pt-2">
              <Dialog.Close asChild>
                <Button type="button" variant="outline" className="flex-1">Batal</Button>
              </Dialog.Close>
              <Button
                id="submit-withdrawal"
                type="submit"
                disabled={loading}
                className="flex-1 bg-violet-600 hover:bg-violet-700 text-white"
              >
                {loading && <Loader2 size={14} className="animate-spin" />}
                Ajukan Penarikan
              </Button>
            </div>
          </form>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  )
}
