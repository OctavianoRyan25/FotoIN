import { useState } from 'react'
import * as Dialog from '@radix-ui/react-dialog'
import { X, Loader2, Calendar, MapPin, DollarSign } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'
import api from '@/lib/api'
import { toast } from 'sonner'

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
  onCreated: () => void
}

export function CreateGalleryDialog({ open, onOpenChange, onCreated }: Props) {
  const [title, setTitle] = useState('')
  const [eventDate, setEventDate] = useState('')
  const [location, setLocation] = useState('')
  const [description, setDescription] = useState('')
  const [pricePerPhoto, setPricePerPhoto] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    const price = parseFloat(pricePerPhoto)
    if (price < 10000) {
      toast.error('Harga minimum per foto adalah Rp 10.000')
      return
    }

    setLoading(true)
    try {
      await api.post('/galleries', {
        title,
        event_date: eventDate,
        location,
        description,
        price_per_photo: price,
      })
      toast.success('Galeri berhasil dibuat!')
      onCreated()
      onOpenChange(false)
      // Reset form
      setTitle(''); setEventDate(''); setLocation(''); setDescription(''); setPricePerPhoto('')
    } catch (err: unknown) {
      const message = (err as { response?: { data?: { message?: string } } })?.response?.data?.message ?? 'Gagal membuat galeri'
      toast.error(message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <Dialog.Root open={open} onOpenChange={onOpenChange}>
      <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 animate-in fade-in" />
        <Dialog.Content className="fixed left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 z-50 w-full max-w-md bg-background rounded-xl shadow-xl p-6 animate-in fade-in slide-in-from-bottom-4">
          <div className="flex items-center justify-between mb-6">
            <Dialog.Title className="text-lg font-semibold">Buat Galeri Baru</Dialog.Title>
            <Dialog.Close asChild>
              <Button variant="ghost" size="icon" className="h-8 w-8">
                <X size={16} />
              </Button>
            </Dialog.Close>
          </div>

          <form onSubmit={handleSubmit} id="create-gallery-form" className="space-y-4">
            <div>
              <label className="text-sm font-medium mb-1.5 block">Judul Event *</label>
              <Input
                id="gallery-title"
                placeholder="Sunmori Semarang - Juli 2024"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                required
              />
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div>
                <label className="text-sm font-medium mb-1.5 flex items-center gap-1.5">
                  <Calendar size={13} /> Tanggal Event *
                </label>
                <Input
                  id="gallery-event-date"
                  type="date"
                  value={eventDate}
                  onChange={(e) => setEventDate(e.target.value)}
                  required
                />
              </div>
              <div>
                <label className="text-sm font-medium mb-1.5 flex items-center gap-1.5">
                  <MapPin size={13} /> Lokasi
                </label>
                <Input
                  id="gallery-location"
                  placeholder="Semarang, Jawa Tengah"
                  value={location}
                  onChange={(e) => setLocation(e.target.value)}
                />
              </div>
            </div>

            <div>
              <label className="text-sm font-medium mb-1.5 flex items-center gap-1.5">
                <DollarSign size={13} /> Harga per Foto (IDR) *
              </label>
              <div className="relative">
                <span className="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-muted-foreground">Rp</span>
                <Input
                  id="gallery-price"
                  type="number"
                  min={10000}
                  step={1000}
                  placeholder="25000"
                  value={pricePerPhoto}
                  onChange={(e) => setPricePerPhoto(e.target.value)}
                  required
                  className="pl-9"
                />
              </div>
              <p className="text-xs text-muted-foreground mt-1">Minimum Rp 10.000 per foto</p>
            </div>

            <div>
              <label className="text-sm font-medium mb-1.5 block">Deskripsi (opsional)</label>
              <textarea
                id="gallery-description"
                placeholder="Ceritakan tentang event ini..."
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                rows={3}
                className={cn(
                  'flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm',
                  'placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2',
                  'focus-visible:ring-ring resize-none'
                )}
              />
            </div>

            <div className="flex gap-3 pt-2">
              <Dialog.Close asChild>
                <Button type="button" variant="outline" className="flex-1">Batal</Button>
              </Dialog.Close>
              <Button
                id="submit-create-gallery"
                type="submit"
                disabled={loading}
                className="flex-1 bg-violet-600 hover:bg-violet-700 text-white"
              >
                {loading && <Loader2 size={14} className="animate-spin" />}
                Buat Galeri
              </Button>
            </div>
          </form>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  )
}
