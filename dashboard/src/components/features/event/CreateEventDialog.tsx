import { useState } from 'react';
import { eventApi } from '@/lib/event-api';
import type { CreateEventInput } from '@/lib/event-api';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { toast } from 'sonner';

interface CreateEventDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: () => void;
}

export function CreateEventDialog({
  open,
  onOpenChange,
  onSuccess,
}: CreateEventDialogProps) {
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState<CreateEventInput>({
    title: '',
    location: '',
    description: '',
    event_date: '',
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!formData.title.trim()) {
      toast.error('Nama event harus diisi');
      return;
    }

    try {
      setLoading(true);
      await eventApi.create(formData);
      toast.success('Event berhasil dibuat');
      onOpenChange(false);
      onSuccess();
      setFormData({ title: '', location: '', description: '', event_date: '' });
    } catch (error: any) {
      toast.error(error.response?.data?.message || 'Gagal membuat event');
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>Buat Event Baru</DialogTitle>
            <DialogDescription>
              Event membantu Anda mengorganisir galeri berdasarkan tempat atau acara
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <Label htmlFor="title">
                Nama Event <span className="text-red-500">*</span>
              </Label>
              <Input
                id="title"
                placeholder="contoh: Borobudur Marathon 2026"
                value={formData.title}
                onChange={(e) =>
                  setFormData({ ...formData, title: e.target.value })
                }
                required
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="event_date">Tanggal Event</Label>
              <Input
                id="event_date"
                type="date"
                value={formData.event_date}
                onChange={(e) =>
                  setFormData({ ...formData, event_date: e.target.value })
                }
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="location">Lokasi</Label>
              <Input
                id="location"
                placeholder="contoh: GOR Tri Lomba Juang Semarang"
                value={formData.location}
                onChange={(e) =>
                  setFormData({ ...formData, location: e.target.value })
                }
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="description">Deskripsi</Label>
              <Textarea
                id="description"
                placeholder="Deskripsi singkat tentang event..."
                value={formData.description}
                onChange={(e) =>
                  setFormData({ ...formData, description: e.target.value })
                }
                rows={3}
              />
            </div>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={loading}
            >
              Batal
            </Button>
            <Button type="submit" disabled={loading}>
              {loading ? 'Membuat...' : 'Buat Event'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
