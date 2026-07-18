import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { eventApi } from '@/lib/event-api';
import type { Event } from '@/lib/event-api';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardFooter, CardHeader } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Plus, Search, Calendar, MapPin, Trash2 } from 'lucide-react';
import { toast } from 'sonner';
import { CreateEventDialog } from '@/components/features/event/CreateEventDialog';

export default function EventListPage() {
  const [events, setEvents] = useState<Event[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');
  const [showCreateDialog, setShowCreateDialog] = useState(false);
  const navigate = useNavigate();

  const fetchEvents = async () => {
    try {
      setLoading(true);
      const response = await eventApi.list(1, 100);
      setEvents(response.data || []);
    } catch (error) {
      toast.error('Gagal memuat event');
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchEvents();
  }, []);

  const handleDelete = async (id: string, title: string) => {
    if (!confirm(`Hapus event "${title}"?`)) return;

    try {
      await eventApi.delete(id);
      toast.success('Event berhasil dihapus');
      fetchEvents();
    } catch (error) {
      toast.error('Gagal menghapus event');
      console.error(error);
    }
  };

  const filteredEvents = events.filter(
    (event) =>
      event.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
      event.location?.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Event</h1>
          <p className="text-muted-foreground mt-1">
            Kelola event dan lokasi pengambilan foto
          </p>
        </div>
        <Button onClick={() => setShowCreateDialog(true)}>
          <Plus className="mr-2 h-4 w-4" />
          Buat Event
        </Button>
      </div>

      <div className="flex items-center gap-4">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            placeholder="Cari event atau lokasi..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-9"
          />
        </div>
      </div>

      {loading ? (
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {[1, 2, 3].map((i) => (
            <Card key={i} className="animate-pulse">
              <CardHeader className="h-48 bg-muted" />
              <CardContent className="pt-6">
                <div className="h-6 bg-muted rounded mb-2" />
                <div className="h-4 bg-muted rounded w-2/3" />
              </CardContent>
            </Card>
          ))}
        </div>
      ) : filteredEvents.length === 0 ? (
        <Card className="p-12 text-center">
          <Calendar className="mx-auto h-12 w-12 text-muted-foreground mb-4" />
          <h3 className="text-lg font-semibold mb-2">Belum ada event</h3>
          <p className="text-muted-foreground mb-4">
            {searchQuery
              ? 'Tidak ada event yang cocok dengan pencarian'
              : 'Buat event pertama Anda untuk mulai mengorganisir galeri foto'}
          </p>
          {!searchQuery && (
            <Button onClick={() => setShowCreateDialog(true)}>
              <Plus className="mr-2 h-4 w-4" />
              Buat Event Pertama
            </Button>
          )}
        </Card>
      ) : (
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {filteredEvents.map((event) => (
            <Card
              key={event.id}
              className="overflow-hidden hover:shadow-lg transition-shadow cursor-pointer"
            >
              <div
                onClick={() => navigate(`/dashboard/events/${event.id}`)}
                className="h-48 relative overflow-hidden flex items-center justify-center" style={{ background: 'linear-gradient(135deg, #23474A, #1F1B16)' }}
              >
                {event.cover_photo_url ? (
                  <img
                    src={event.cover_photo_url}
                    alt={event.title}
                    className="w-full h-full object-cover"
                  />
                ) : (
                  <Calendar className="h-16 w-16 text-white opacity-50" />
                )}
              </div>

              <CardHeader className="pb-3">
                <div className="flex items-start justify-between gap-2">
                  <h3 className="font-semibold text-lg line-clamp-2 flex-1">
                    {event.title}
                  </h3>
                  <Badge variant={event.is_active ? 'brand' : 'secondary'}>
                    {event.is_active ? 'Aktif' : 'Nonaktif'}
                  </Badge>
                </div>
              </CardHeader>

              <CardContent className="pb-3 space-y-2 text-sm text-muted-foreground">
                {event.location && (
                  <div className="flex items-center gap-2">
                    <MapPin className="h-4 w-4 flex-shrink-0" />
                    <span className="line-clamp-1">{event.location}</span>
                  </div>
                )}
                {event.event_date && (
                  <div className="flex items-center gap-2">
                    <Calendar className="h-4 w-4 flex-shrink-0" />
                    <span>
                      {new Date(event.event_date).toLocaleDateString('id-ID', {
                        day: 'numeric',
                        month: 'long',
                        year: 'numeric',
                      })}
                    </span>
                  </div>
                )}
              </CardContent>

              <CardFooter className="pt-3 border-t">
                <Button
                  variant="ghost"
                  size="sm"
                  className="text-red-600 hover:text-red-700 hover:bg-red-50"
                  onClick={(e) => {
                    e.stopPropagation();
                    handleDelete(event.id, event.title);
                  }}
                >
                  <Trash2 className="mr-2 h-4 w-4" />
                  Hapus
                </Button>
              </CardFooter>
            </Card>
          ))}
        </div>
      )}

      <CreateEventDialog
        open={showCreateDialog}
        onOpenChange={setShowCreateDialog}
        onSuccess={fetchEvents}
      />
    </div>
  );
}
