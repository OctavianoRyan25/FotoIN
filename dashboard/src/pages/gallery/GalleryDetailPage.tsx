import { useEffect, useState, useRef } from "react";
import { useParams, Link } from "react-router-dom";
import {
  ArrowLeft,
  Upload,
  Trash2,
  Download,
  ImageIcon,
  Loader2,
  Plus,
  Check,
  X,
  ChevronLeft,
  ChevronRight,
  MapPin,
  Calendar,
  Wallet,
  Images,
  ExternalLink,
  FileImage,
} from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { formatRupiah, formatBytes, formatDate } from "@/lib/utils";
import api from "@/lib/api";
import type { Gallery, Photo } from "@/types";
import { toast } from "sonner";

export function GalleryDetailPage() {
  const { id } = useParams<{ id: string }>();
  const [gallery, setGallery] = useState<Gallery | null>(null);
  const [photos, setPhotos] = useState<Photo[]>([]);
  const [loading, setLoading] = useState(true);
  const [uploading, setUploading] = useState(false);
  const [uploadQueue, setUploadQueue] = useState<File[]>([]);
  const [uploadProgress, setUploadProgress] = useState<
    Record<string, "uploading" | "done" | "error">
  >({});
  const [selectedPhoto, setSelectedPhoto] = useState<number | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const fetchData = async () => {
    if (!id) return;
    setLoading(true);
    try {
      const [galRes, photoRes] = await Promise.all([
        api.get(`/galleries/${id}`),
        api.get(`/galleries/${id}/photos?limit=100`),
      ]);
      setGallery(galRes.data.data);
      setPhotos(photoRes.data.data ?? []);
    } catch {
      toast.error("Gagal memuat data galeri");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [id]);

  const handleFileSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files ?? []);
    if (!files.length || !id) return;

    // Reset input so same file can be re-selected
    e.target.value = "";

    setUploading(true);
    const progress: Record<string, "uploading" | "done" | "error"> = {};
    files.forEach((f) => (progress[f.name] = "uploading"));
    setUploadQueue(files);
    setUploadProgress({ ...progress });

    for (const file of files) {
      try {
        // Guard: empty MIME causes S3 SignatureDoesNotMatch
        const contentType = file.type || "application/octet-stream";

        // Step 1 — Get presigned PUT URL from backend
        const urlRes = await api.post(`/galleries/${id}/photos/upload-url`, {
          filename: file.name,
          content_type: contentType,
        });
        const { upload_url, s3_key } = urlRes.data.data;

        // Step 2 — PUT file directly to S3
        // Content-Type header MUST exactly match what was signed
        const uploadResponse = await fetch(upload_url, {
          method: "PUT",
          body: file,
          headers: { "Content-Type": contentType },
        });

        if (!uploadResponse.ok) {
          // Parse AWS XML error body for a meaningful message
          const errText = await uploadResponse.text();
          const code =
            errText.match(/<Code>(.*?)<\/Code>/)?.[1] ?? uploadResponse.status;
          const msg =
            errText.match(/<Message>(.*?)<\/Message>/)?.[1] ??
            uploadResponse.statusText;
          throw new Error(`S3 [${code}]: ${msg}`);
        }

        // Step 3 — Register in DB and update storage quota
        await api.post(`/galleries/${id}/photos/confirm`, {
          gallery_id: id,
          s3_key,
          filename: file.name,
          size_bytes: file.size,
        });

        progress[file.name] = "done";
      } catch (err: unknown) {
        progress[file.name] = "error";
        const msg = err instanceof Error ? err.message : "Upload gagal";
        toast.error(`${file.name}: ${msg}`);
        console.error("[S3 Upload]", file.name, err);
      }
      setUploadProgress({ ...progress });
    }

    setUploading(false);
    setTimeout(() => {
      setUploadQueue([]);
      setUploadProgress({});
      fetchData();
    }, 2000);
  };

  const handleDeletePhoto = async (photoId: string) => {
    try {
      await api.delete(`/photos/${photoId}`);
      toast.success("Foto berhasil dihapus");
      setPhotos((prev) => prev.filter((p) => p.id !== photoId));
    } catch {
      toast.error("Gagal menghapus foto");
    }
  };

  const doneCount = Object.values(uploadProgress).filter(
    (s) => s === "done",
  ).length;
  const errorCount = Object.values(uploadProgress).filter(
    (s) => s === "error",
  ).length;

  if (loading)
    return (
      <div className="flex flex-col items-center justify-center h-64 gap-3">
        <Loader2 className="animate-spin text-brand" size={32} />
        <p className="text-sm text-muted-foreground">Memuat galeri...</p>
      </div>
    );

  if (!gallery)
    return (
      <div className="text-center py-20">
        <ImageIcon
          size={40}
          className="mx-auto text-muted-foreground/40 mb-3"
        />
        <p className="text-muted-foreground">Galeri tidak ditemukan</p>
        <Link to="/dashboard/galleries">
          <Button className="mt-4 gap-2">
            <ArrowLeft size={14} /> Kembali
          </Button>
        </Link>
      </div>
    );

  return (
    <div className="space-y-6">
      {/* Breadcrumb */}
      <Link
        to="/dashboard/galleries"
        className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground group"
      >
        <ArrowLeft
          size={16}
          className="group-hover:-translate-x-0.5 transition-transform"
        />{" "}
        Kembali ke Galeri
      </Link>

      {/* Header */}
      <div className="flex items-start justify-between">
        <div>
          <div className="flex items-center gap-3 mb-1">
            <h1 className="text-2xl font-bold">{gallery.title}</h1>
            <Badge
              variant={gallery.is_active ? "success" : "secondary"}
              className="gap-1"
            >
              <span
                className={`h-1.5 w-1.5 rounded-full ${gallery.is_active ? "bg-emerald-500" : "bg-muted-foreground/50"}`}
              />
              {gallery.is_active ? "Aktif" : "Nonaktif"}
            </Badge>
          </div>
          <div className="flex flex-wrap items-center gap-x-4 gap-y-1.5 text-sm text-muted-foreground">
            <span className="inline-flex items-center gap-1.5">
              <MapPin size={14} /> {gallery.location}
            </span>
            <span className="inline-flex items-center gap-1.5">
              <Calendar size={14} /> {formatDate(gallery.event_date)}
            </span>
            <span className="inline-flex items-center gap-1.5">
              <Wallet size={14} />
              <span className="font-medium text-foreground">
                {formatRupiah(gallery.price_per_photo)}
              </span>
              /foto
            </span>
            <span className="inline-flex items-center gap-1.5">
              <Images size={14} /> {photos.length} foto
            </span>
          </div>
        </div>

        <div className="flex gap-2">
          <a
            href={`${window.location.origin.replace("5173", "3000")}/gallery/${gallery.slug}`}
            target="_blank"
            rel="noreferrer"
          >
            <Button
              variant="outline"
              size="sm"
              id="view-public-gallery"
              className="gap-2"
            >
              <ExternalLink size={14} /> Lihat Halaman Publik
            </Button>
          </a>
          <Button
            id="upload-photos-btn"
            size="sm"
            className="gap-2"
            onClick={() => fileInputRef.current?.click()}
            disabled={uploading}
          >
            {uploading ? (
              <Loader2 size={14} className="animate-spin" />
            ) : (
              <Upload size={14} />
            )}
            {uploading ? "Mengupload..." : "Upload Foto"}
          </Button>
          <input
            ref={fileInputRef}
            type="file"
            multiple
            accept="image/*"
            className="hidden"
            onChange={handleFileSelect}
            id="file-upload-input"
          />
        </div>
      </div>

      {/* Upload Progress */}
      {uploadQueue.length > 0 && (
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between mb-3">
              <p className="text-sm font-medium flex items-center gap-2">
                <Upload size={14} className="text-brand" />
                Mengupload {uploadQueue.length} foto...
              </p>
              <p className="text-xs text-muted-foreground">
                {doneCount}/{uploadQueue.length} selesai
                {errorCount > 0 && (
                  <span className="text-red-500 ml-1.5">
                    · {errorCount} gagal
                  </span>
                )}
              </p>
            </div>
            {/* Overall progress bar */}
            <div className="h-1.5 w-full rounded-full bg-muted overflow-hidden mb-3">
              <div
                className="h-full bg-brand transition-all duration-300"
                style={{ width: `${(doneCount / uploadQueue.length) * 100}%` }}
              />
            </div>
            <div className="space-y-2">
              {uploadQueue.map((file) => (
                <div
                  key={file.name}
                  className="flex items-center gap-3 text-sm"
                >
                  <div
                    className={`w-4 h-4 flex-shrink-0 ${
                      uploadProgress[file.name] === "done"
                        ? "text-emerald-500"
                        : uploadProgress[file.name] === "error"
                          ? "text-red-500"
                          : "text-violet-600"
                    }`}
                  >
                    {uploadProgress[file.name] === "done" ? (
                      <Check size={16} />
                    ) : uploadProgress[file.name] === "error" ? (
                      <X size={16} />
                    ) : (
                      <Loader2 size={16} className="animate-spin" />
                    )}
                  </div>
                  <FileImage
                    size={14}
                    className="text-muted-foreground flex-shrink-0"
                  />
                  <span className="truncate flex-1">{file.name}</span>
                  <span className="text-muted-foreground text-xs">
                    {formatBytes(file.size)}
                  </span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Photos Grid */}
      {photos.length === 0 ? (
        <div
          className="border-2 border-dashed border-muted rounded-xl p-16 text-center cursor-pointer hover:border-brand/40 hover:bg-brand/5 transition-colors"
          onClick={() => fileInputRef.current?.click()}
          id="upload-drop-zone"
        >
          <ImageIcon size={48} className="mx-auto text-muted-foreground mb-4" />
          <h3 className="text-lg font-medium mb-2">Belum ada foto</h3>
          <p className="text-muted-foreground text-sm mb-4">
            Klik atau seret foto ke sini untuk mengupload
          </p>
          <Button className="gap-2">
            <Plus size={16} /> Pilih Foto
          </Button>
        </div>
      ) : (
        <div>
          <p className="text-sm text-muted-foreground mb-4 flex items-center gap-1.5">
            <Images size={14} /> {photos.length} foto
          </p>
          <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-3">
            {photos.map((photo, idx) => (
              <div
                key={photo.id}
                className="group relative aspect-square rounded-lg overflow-hidden bg-muted cursor-pointer ring-1 ring-transparent hover:ring-brand/40 transition-all"
                onClick={() => setSelectedPhoto(idx)}
              >
                {photo.preview_url ? (
                  <img
                    src={photo.preview_url}
                    alt={photo.filename}
                    className="w-full h-full object-cover transition-transform group-hover:scale-105"
                    loading="lazy"
                  />
                ) : (
                  <div className="w-full h-full flex items-center justify-center">
                    <ImageIcon size={24} className="text-muted-foreground" />
                  </div>
                )}
                {/* Watermark overlay */}
                <div className="absolute inset-0 flex items-center justify-center pointer-events-none opacity-30">
                  <span className="text-white font-bold text-xs rotate-[-30deg] select-none">
                    FotoIN
                  </span>
                </div>
                {/* Hover actions */}
                <div className="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-2">
                  <Button
                    size="icon"
                    variant="secondary"
                    className="h-8 w-8"
                    title="Perbesar"
                    onClick={(e) => {
                      e.stopPropagation();
                      setSelectedPhoto(idx);
                    }}
                  >
                    <ImageIcon size={14} />
                  </Button>
                  <Button
                    size="icon"
                    variant="destructive"
                    className="h-8 w-8"
                    title="Hapus foto"
                    onClick={(e) => {
                      e.stopPropagation();
                      handleDeletePhoto(photo.id);
                    }}
                  >
                    <Trash2 size={14} />
                  </Button>
                </div>
                <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/60 to-transparent p-2 opacity-0 group-hover:opacity-100 transition-opacity">
                  <p className="text-white text-xs truncate flex items-center gap-1">
                    <FileImage size={10} className="flex-shrink-0" />
                    {photo.filename}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Lightbox */}
      {selectedPhoto !== null && (
        <div
          className="fixed inset-0 bg-black/90 z-50 flex items-center justify-center"
          onClick={() => setSelectedPhoto(null)}
        >
          <button
            className="absolute top-4 right-4 text-white hover:text-gray-300"
            onClick={() => setSelectedPhoto(null)}
          >
            <X size={24} />
          </button>
          <button
            className="absolute left-4 top-1/2 -translate-y-1/2 text-white hover:text-gray-300 disabled:opacity-30"
            onClick={(e) => {
              e.stopPropagation();
              setSelectedPhoto((p) => (p! > 0 ? p! - 1 : p));
            }}
            disabled={selectedPhoto === 0}
          >
            <ChevronLeft size={32} />
          </button>
          <img
            src={photos[selectedPhoto]?.preview_url}
            alt={photos[selectedPhoto]?.filename}
            className="max-w-[90vw] max-h-[90vh] object-contain"
            onClick={(e) => e.stopPropagation()}
          />
          <button
            className="absolute right-4 top-1/2 -translate-y-1/2 text-white hover:text-gray-300 disabled:opacity-30"
            onClick={(e) => {
              e.stopPropagation();
              setSelectedPhoto((p) => (p! < photos.length - 1 ? p! + 1 : p));
            }}
            disabled={selectedPhoto === photos.length - 1}
          >
            <ChevronRight size={32} />
          </button>
          <div className="absolute bottom-4 flex flex-col items-center gap-1">
            <p className="text-white/90 text-xs flex items-center gap-1.5">
              <FileImage size={12} /> {photos[selectedPhoto]?.filename}
            </p>
            <p className="text-white/60 text-sm">
              {selectedPhoto + 1} / {photos.length}
            </p>
          </div>
        </div>
      )}
    </div>
  );
}
