import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import {
  Plus,
  Images,
  MapPin,
  Calendar,
  Trash2,
  Edit,
  ExternalLink,
  Loader2,
  Search,
  SlidersHorizontal,
  Camera,
  Eye,
  ImageOff,
  ArrowUpRight,
  Sparkles,
} from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { formatRupiah, formatDate } from "@/lib/utils";
import api from "@/lib/api";
import type { Gallery } from "@/types";
import { toast } from "sonner";
import { CreateGalleryDialog } from "@/components/features/gallery/CreateGalleryDialog";

export function GalleryListPage() {
  const [galleries, setGalleries] = useState<Gallery[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [showCreate, setShowCreate] = useState(false);
  const [total, setTotal] = useState(0);

  const fetchGalleries = async () => {
    setLoading(true);
    try {
      const res = await api.get("/galleries?limit=50");
      setGalleries(res.data.data ?? []);
      setTotal(res.data.meta?.total ?? 0);
    } catch {
      toast.error("Gagal memuat galeri");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchGalleries();
  }, []);

  const handleDelete = async (id: string, title: string) => {
    if (!confirm(`Hapus galeri "${title}"? Semua foto akan ikut terhapus.`))
      return;
    try {
      await api.delete(`/galleries/${id}`);
      toast.success("Galeri berhasil dihapus");
      fetchGalleries();
    } catch {
      toast.error("Gagal menghapus galeri");
    }
  };

  const filtered = galleries.filter(
    (g) =>
      g.title.toLowerCase().includes(search.toLowerCase()) ||
      g.location.toLowerCase().includes(search.toLowerCase()),
  );

  const activeCount = galleries.filter((g) => g.is_active).length;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <div className="flex items-center gap-2">
            <Camera size={20} className="text-brand" />
            <h1 className="text-2xl font-bold">Galeri</h1>
          </div>
          <p className="text-muted-foreground text-sm mt-1 flex items-center gap-1.5">
            <span>{total} galeri ditemukan</span>
            {activeCount > 0 && (
              <>
                <span className="text-muted-foreground/40">•</span>
                <span className="inline-flex items-center gap-1 text-emerald-600">
                  <span className="h-1.5 w-1.5 rounded-full bg-emerald-500" />
                  {activeCount} aktif
                </span>
              </>
            )}
          </p>
        </div>
        <Button
          id="open-create-gallery"
          onClick={() => setShowCreate(true)}
          className="gap-2"
        >
          <Plus size={16} /> Buat Galeri
        </Button>
      </div>

      {/* Search & Filter */}
      <div className="flex gap-3">
        <div className="relative flex-1 max-w-sm">
          <Search
            size={16}
            className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
          />
          <Input
            id="gallery-search"
            placeholder="Cari galeri..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="pl-9"
          />
        </div>
        <Button variant="outline" className="gap-2">
          <SlidersHorizontal size={16} /> Filter
        </Button>
      </div>

      {/* Grid */}
      {loading ? (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {[1, 2, 3, 4, 5, 6].map((i) => (
            <div
              key={i}
              className="h-48 bg-muted rounded-xl animate-pulse flex items-center justify-center"
            >
              <Loader2
                size={20}
                className="text-muted-foreground/30 animate-spin"
              />
            </div>
          ))}
        </div>
      ) : filtered.length === 0 ? (
        <div className="text-center py-20">
          {search ? (
            <ImageOff
              size={48}
              className="mx-auto text-muted-foreground mb-4"
            />
          ) : (
            <div className="relative inline-block mb-4">
              <Images size={48} className="mx-auto text-muted-foreground" />
              <Sparkles
                size={16}
                className="absolute -top-1 -right-2 text-brand"
              />
            </div>
          )}
          <h3 className="text-lg font-medium mb-2">
            {search ? "Galeri tidak ditemukan" : "Belum ada galeri"}
          </h3>
          <p className="text-muted-foreground text-sm mb-6">
            {search
              ? "Coba kata kunci yang berbeda"
              : "Buat galeri pertama Anda untuk mulai menjual foto"}
          </p>
          {!search && (
            <Button
              id="create-gallery-empty"
              onClick={() => setShowCreate(true)}
              className="gap-2"
            >
              <Plus size={16} /> Buat Galeri
            </Button>
          )}
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {filtered.map((gallery) => (
            <Card
              key={gallery.id}
              id={`gallery-card-${gallery.id}`}
              className="group hover:shadow-lg hover:-translate-y-0.5 transition-all duration-200 overflow-hidden"
            >
              {/* Cover */}
              <div
                className="h-32 flex items-center justify-center relative overflow-hidden"
                style={{
                  background: "linear-gradient(135deg, #EFE9DA, #F6F3EC)",
                }}
              >
                {gallery.cover_photo_url ? (
                  <img
                    src={gallery.cover_photo_url}
                    alt={gallery.title}
                    className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                  />
                ) : (
                  <Camera
                    size={32}
                    className="text-muted-foreground/40"
                    strokeWidth={1.5}
                  />
                )}
                <div className="absolute inset-0 bg-gradient-to-t from-black/25 via-black/0 to-black/0 opacity-0 group-hover:opacity-100 transition-opacity" />

                {/* Photo count chip on cover */}
                <div className="absolute bottom-2 left-2 flex items-center gap-1 rounded-full bg-black/50 backdrop-blur-sm px-2 py-0.5 text-[11px] font-medium text-white opacity-0 group-hover:opacity-100 transition-opacity">
                  <Images size={11} /> {gallery.photo_count}
                </div>

                <div className="absolute top-2 right-2 flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                  <Link to={`/dashboard/galleries/${gallery.id}`}>
                    <Button
                      size="icon"
                      variant="secondary"
                      className="h-7 w-7"
                      title="Lihat galeri"
                    >
                      <Eye size={12} />
                    </Button>
                  </Link>
                  <Link to={`/dashboard/galleries/${gallery.id}/edit`}>
                    <Button
                      size="icon"
                      variant="secondary"
                      className="h-7 w-7"
                      title="Edit galeri"
                    >
                      <Edit size={12} />
                    </Button>
                  </Link>
                  <Button
                    size="icon"
                    variant="destructive"
                    className="h-7 w-7"
                    title="Hapus galeri"
                    onClick={() => handleDelete(gallery.id, gallery.title)}
                  >
                    <Trash2 size={12} />
                  </Button>
                </div>
              </div>

              <CardContent className="p-4">
                <div className="flex items-start justify-between mb-2">
                  <h3 className="font-semibold text-sm truncate flex-1">
                    {gallery.title}
                  </h3>
                  <Badge
                    variant={gallery.is_active ? "success" : "secondary"}
                    className="ml-2 flex-shrink-0 gap-1"
                  >
                    <span
                      className={`h-1.5 w-1.5 rounded-full ${gallery.is_active ? "bg-emerald-500" : "bg-muted-foreground/50"}`}
                    />
                    {gallery.is_active ? "Aktif" : "Nonaktif"}
                  </Badge>
                </div>

                <div className="space-y-1 text-xs text-muted-foreground">
                  <div className="flex items-center gap-1.5">
                    <MapPin size={11} />
                    <span className="truncate">{gallery.location || "—"}</span>
                  </div>
                  <div className="flex items-center gap-1.5">
                    <Calendar size={11} />
                    <span>{formatDate(gallery.event_date)}</span>
                  </div>
                </div>

                <div className="flex items-center justify-between mt-3 pt-3 border-t">
                  <div>
                    <p className="text-xs text-muted-foreground">Harga/foto</p>
                    <p className="text-sm font-semibold text-brand">
                      {formatRupiah(gallery.price_per_photo)}
                    </p>
                  </div>
                  <div className="text-right flex items-center gap-2">
                    <div>
                      <p className="text-xs text-muted-foreground">Foto</p>
                      <p className="text-sm font-medium flex items-center gap-1 justify-end">
                        <Images size={12} className="text-muted-foreground" />
                        {gallery.photo_count}
                      </p>
                    </div>
                    <ArrowUpRight
                      size={14}
                      className="text-muted-foreground/30 group-hover:text-brand group-hover:translate-x-0.5 group-hover:-translate-y-0.5 transition-all"
                    />
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      <CreateGalleryDialog
        open={showCreate}
        onOpenChange={setShowCreate}
        onCreated={fetchGalleries}
      />
    </div>
  );
}
