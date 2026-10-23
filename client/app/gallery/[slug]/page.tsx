"use client";

import { useEffect, useState, useCallback } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import {
  ShoppingCart,
  X,
  MapPin,
  Calendar,
  CheckCircle,
  ArrowLeft,
} from "lucide-react";
import { toast } from "sonner";
import QRPaymentBlock from "../../../components/QRPaymentBlock";
import {
  FONT_CLASSES,
  FilmStrip,
  Header,
  Footer,
  Photo,
  Gallery,
  formatRupiah,
  formatDate,
} from "../../components/site-ui";

const API_URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api/v1";

function WatermarkOverlay() {
  return (
    <div className="absolute inset-0 pointer-events-none select-none overflow-hidden flex items-center justify-center">
      <div className="absolute inset-0 opacity-[0.16] flex flex-wrap items-center justify-around gap-10 rotate-[-25deg] scale-125">
        {Array.from({ length: 12 }).map((_, i) => (
          <span
            key={i}
            className="text-white text-xs sm:text-sm font-extrabold uppercase tracking-widest drop-shadow-[0_1px_2px_rgba(0,0,0,0.8)]"
          >
            FotoIN • Preview
          </span>
        ))}
      </div>
      <div className="opacity-30 px-2.5 py-0.5 rounded border border-white/50 bg-black/30 backdrop-blur-[1px] rotate-[-12deg]">
        <span className="text-white text-[11px] font-bold tracking-wider uppercase drop-shadow">
          FotoIN
        </span>
      </div>
    </div>
  );
}

export default function GalleryPage() {
  const { slug } = useParams<{ slug: string }>();
  const [gallery, setGallery] = useState<Gallery | null>(null);
  const [photos, setPhotos] = useState<Photo[]>([]);
  const [loading, setLoading] = useState(true);
  const [cart, setCart] = useState<Set<string>>(new Set());
  const [showCart, setShowCart] = useState(false);
  const [checkoutEmail, setCheckoutEmail] = useState("");
  const [checkoutLoading, setCheckoutLoading] = useState(false);
  const [selectedPhoto, setSelectedPhoto] = useState<Photo | null>(null);
  const [qrData, setQrData] = useState<Record<string, unknown> | null>(null);

  // useEffect(() => {
  //   const fetchGallery = async () => {
  //     setLoading(true);
  //     try {
  //       const [galRes, photoRes] = await Promise.all([
  //         fetch(`${API_URL}/public/galleries/${slug}`),
  //         fetch(`${API_URL}/public/galleries/${gallery?.id}/photos?limit=100`),
  //       ]);
  //       if (!galRes.ok) throw new Error("gallery not found");
  //       const galData = await galRes.json();
  //       const photoData = await photoRes.json();
  //       setGallery(galData.data ?? null);
  //       setPhotos(Array.isArray(photoData.data) ? photoData.data : []);
  //     } catch {
  //       setGallery(null);
  //       setPhotos([]);
  //     } finally {
  //       setLoading(false);
  //     }
  //   };
  //   if (slug) fetchGallery();
  // }, [slug]);

  useEffect(() => {
    const fetchGallery = async () => {
      setLoading(true);
      try {
        const galRes = await fetch(`${API_URL}/public/galleries/${slug}`);
        if (!galRes.ok) throw new Error("gallery not found");
        const galData = await galRes.json();
        const galleryData = galData.data ?? null;
        setGallery(galleryData);

        if (galleryData?.id) {
          const photoRes = await fetch(
            `${API_URL}/public/galleries/${galleryData.id}/photos?limit=100`,
          );
          const photoData = await photoRes.json();
          setPhotos(Array.isArray(photoData.data) ? photoData.data : []);
        } else {
          setPhotos([]);
        }
      } catch {
        setGallery(null);
        setPhotos([]);
      } finally {
        setLoading(false);
      }
    };
    if (slug) fetchGallery();
  }, [slug]);

  const toggleCart = useCallback((photoId: string) => {
    setCart((prev) => {
      const next = new Set(prev);
      if (next.has(photoId)) {
        next.delete(photoId);
      } else {
        next.add(photoId);
      }
      return next;
    });
  }, []);

  const totalPrice = cart.size * (gallery?.price_per_photo ?? 0);

  const handleCheckout = async () => {
    const trimmedEmail = checkoutEmail.trim();
    if (!gallery) {
      toast.error("Galeri tidak tersedia saat ini");
      return;
    }
    if (cart.size === 0) {
      toast.error("Pilih minimal 1 foto");
      return;
    }
    if (!trimmedEmail) {
      toast.error("Masukkan email untuk melanjutkan pembayaran");
      return;
    }

    setCheckoutLoading(true);
    try {
      const res = await fetch(`${API_URL}/public/orders`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          buyer_email: trimmedEmail,
          gallery_id: gallery.id,
          photo_ids: Array.from(cart),
        }),
      });
      const data = await res.json();

      if (data.data?.qr_code_payment) {
        // Show QR code inline — no redirect!
        setQrData({
          ...data.data.qr_code_payment,
          order_id: data.data.order?.id ?? data.data.qr_code_payment.order_id,
        });
      } else {
        toast.success(
          "Pesanan dibuat! Cek email Anda untuk instruksi pembayaran.",
        );
        setCart(new Set());
        setShowCart(false);
      }
    } catch {
      toast.error("Gagal membuat order. Coba lagi.");
    } finally {
      setCheckoutLoading(false);
    }
  };

  // Loading skeleton
  if (loading) {
    return (
      <div className={`${FONT_CLASSES} min-h-screen bg-[#F6F3EC]`}>
        <Header active="gallery" />
        <div className="max-w-6xl mx-auto px-4 py-10">
          <div className="h-8 w-64 mb-3 rounded-lg bg-[#EFE9DA] animate-pulse" />
          <div className="h-4 w-48 mb-10 rounded-lg bg-[#EFE9DA] animate-pulse" />
          <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3">
            {Array.from({ length: 12 }).map((_, i) => (
              <div
                key={i}
                className="aspect-square rounded-xl bg-[#EFE9DA] animate-pulse"
              />
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (!gallery) {
    return <GalleryUnavailableState slug={String(slug)} />;
  }

  return (
    <GalleryView
      gallery={gallery}
      photos={photos}
      cart={cart}
      toggleCart={toggleCart}
      selectedPhoto={selectedPhoto}
      setSelectedPhoto={setSelectedPhoto}
      showCart={showCart}
      setShowCart={setShowCart}
      checkoutEmail={checkoutEmail}
      setCheckoutEmail={setCheckoutEmail}
      handleCheckout={handleCheckout}
      checkoutLoading={checkoutLoading}
      totalPrice={totalPrice}
      qrData={qrData}
      onQrPaid={() => {
        setCart(new Set());
        setShowCart(false);
        setQrData(null);
      }}
    />
  );
}

function GalleryUnavailableState({ slug }: { slug: string }) {
  return (
    <div className={`${FONT_CLASSES} min-h-screen flex flex-col bg-[#F6F3EC]`}>
      <Header active="gallery" />
      <div className="max-w-4xl mx-auto px-4 py-16 flex-1 w-full">
        <div className="rounded-3xl border border-[#1F1B16]/10 bg-white p-8 text-center shadow-sm">
          <p className="text-sm uppercase tracking-[0.25em] text-[#23474A] mb-3">
            Galeri
          </p>
          <h1 className="text-2xl sm:text-3xl font-[family-name:var(--font-display)] text-[#1F1B16] mb-3">
            Galeri belum tersedia
          </h1>
          <p className="text-[#4A443C] max-w-xl mx-auto">
            Belum ada data galeri real dari backend untuk slug {slug}. Pastikan
            API berjalan dan galeri sudah dibuat.
          </p>
          <Link
            href="/events"
            className="inline-flex items-center gap-2 mt-6 rounded-full bg-[#1F1B16] px-5 py-2.5 text-sm font-medium text-[#F6F3EC] hover:bg-[#B8863C] transition-colors"
          >
            <ArrowLeft size={14} /> Kembali ke event
          </Link>
        </div>
      </div>
      <Footer />
    </div>
  );
}

interface GalleryViewProps {
  gallery: Gallery;
  photos: Photo[];
  cart: Set<string>;
  toggleCart: (id: string) => void;
  selectedPhoto: Photo | null;
  setSelectedPhoto: (p: Photo | null) => void;
  showCart: boolean;
  setShowCart: (v: boolean) => void;
  checkoutEmail: string;
  setCheckoutEmail: (v: string) => void;
  handleCheckout: () => void;
  checkoutLoading: boolean;
  totalPrice: number;
  isDemoMode?: boolean;
  qrData?: Record<string, unknown> | null;
  onQrPaid?: () => void;
}

function GalleryView({
  gallery,
  photos,
  cart,
  toggleCart,
  selectedPhoto,
  setSelectedPhoto,
  showCart,
  setShowCart,
  checkoutEmail,
  setCheckoutEmail,
  handleCheckout,
  checkoutLoading,
  totalPrice,
  qrData,
  onQrPaid,
}: GalleryViewProps) {
  const cartButton = (
    <button
      id="open-cart"
      onClick={() => setShowCart(true)}
      className="relative flex items-center gap-2 bg-[#1F1B16] text-[#F6F3EC] px-4 py-2 rounded-full font-medium text-sm hover:bg-[#B8863C] transition-colors"
    >
      <ShoppingCart size={16} />
      <span className="hidden sm:inline">Keranjang</span>
      {cart.size > 0 && (
        <span className="absolute -top-2 -right-2 w-5 h-5 bg-[#A8433F] text-white text-xs rounded-full flex items-center justify-center font-bold">
          {cart.size}
        </span>
      )}
    </button>
  );

  return (
    <div className={`${FONT_CLASSES} min-h-screen flex flex-col bg-[#F6F3EC]`}>
      <Header active="gallery" right={cartButton} />

      {/* Gallery Info */}
      <div className="max-w-6xl mx-auto px-4 pt-10 pb-8 w-full">
        {gallery.event && (
          <Link
            href={`/events/${gallery.event.slug}`}
            className="inline-flex items-center gap-1.5 text-xs font-[family-name:var(--font-mono)] text-[#4A443C]/70 hover:text-[#B8863C] transition-colors mb-5"
          >
            <ArrowLeft size={12} /> {gallery.event.title}
          </Link>
        )}
        <h1 className="text-3xl sm:text-4xl font-[family-name:var(--font-display)] text-[#1F1B16] mb-4 max-w-2xl">
          {gallery.title}
        </h1>
        <div className="flex flex-wrap items-center gap-x-4 gap-y-2 text-sm text-[#4A443C] font-[family-name:var(--font-mono)] mb-4">
          <span className="flex items-center gap-1.5">
            <MapPin size={14} className="text-[#B8863C]" />
            {gallery.location}
          </span>
          <span className="text-[#1F1B16]/20 hidden sm:inline">·</span>
          <span className="flex items-center gap-1.5">
            <Calendar size={14} className="text-[#B8863C]" />
            {formatDate(gallery.event_date)}
          </span>
          <span className="text-[#1F1B16]/20 hidden sm:inline">·</span>
          <span>oleh {gallery.photographer?.name}</span>
          <span className="text-[#1F1B16]/20 hidden sm:inline">·</span>
          <span className="font-semibold text-[#B8863C]">
            {formatRupiah(gallery.price_per_photo)}/foto
          </span>
        </div>
        {gallery.description && (
          <p className="text-[#4A443C] text-sm max-w-2xl leading-relaxed">
            {gallery.description}
          </p>
        )}
      </div>

      <FilmStrip tone="ink" />

      {/* Photo Grid */}
      <div className="max-w-6xl mx-auto px-4 pt-8 pb-28 flex-1 w-full">
        <p className="text-sm text-[#4A443C]/60 font-[family-name:var(--font-mono)] mb-5">
          {photos.length} foto tersedia
        </p>
        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3">
          {photos.map((photo, i) => {
            const isSelected = cart.has(photo.id);
            return (
              <div
                key={photo.id}
                className={`relative group aspect-square rounded-xl overflow-hidden cursor-pointer ring-2 transition-all ${
                  isSelected
                    ? "ring-[#B8863C] ring-offset-2 ring-offset-[#F6F3EC]"
                    : "ring-transparent hover:ring-[#1F1B16]/15"
                }`}
                onClick={() => setSelectedPhoto(photo)}
              >
                <img
                  src={photo.preview_url}
                  alt={photo.filename}
                  draggable={false}
                  onContextMenu={(e) => e.preventDefault()}
                  className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105 pointer-events-none select-none"
                  loading="lazy"
                />
                <WatermarkOverlay />
                <div className="absolute bottom-2 left-2 px-1.5 py-0.5 rounded bg-[#1F1B16]/55 backdrop-blur-sm">
                  <span className="text-[10px] font-[family-name:var(--font-mono)] text-[#F6F3EC] tracking-wide">
                    No. {String(i + 1).padStart(3, "0")}
                  </span>
                </div>
                {isSelected && (
                  <div className="absolute top-2 right-2 w-6 h-6 bg-[#B8863C] rounded-full flex items-center justify-center">
                    <CheckCircle size={14} className="text-white" />
                  </div>
                )}
                <div
                  className="absolute inset-0 bg-[#1F1B16]/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-end p-3"
                  onClick={(e) => {
                    e.stopPropagation();
                    toggleCart(photo.id);
                  }}
                >
                  <button
                    id={`select-photo-${photo.id}`}
                    className={`w-full py-1.5 rounded-lg text-sm font-medium transition-colors ${
                      isSelected
                        ? "bg-[#A8433F] text-white"
                        : "bg-[#F6F3EC] text-[#1F1B16]"
                    }`}
                  >
                    {isSelected ? "Batal Pilih" : "+ Pilih Foto"}
                  </button>
                </div>
              </div>
            );
          })}
        </div>
      </div>

      {!showCart && <Footer />}

      {/* Sticky Bottom Cart Bar */}
      {cart.size > 0 && !showCart && (
        <div className="fixed bottom-0 left-0 right-0 z-30 bg-[#1F1B16] px-4 py-4">
          <div className="max-w-6xl mx-auto flex items-center justify-between gap-4">
            <div>
              <p className="text-sm text-[#F6F3EC]/60 font-[family-name:var(--font-mono)]">
                {cart.size} foto dipilih
              </p>
              <p className="text-lg font-[family-name:var(--font-display)] text-[#D8AE6E]">
                {formatRupiah(totalPrice)}
              </p>
            </div>
            <button
              id="checkout-cta"
              onClick={() => setShowCart(true)}
              className="px-6 py-3 bg-[#B8863C] text-[#1F1B16] rounded-full font-semibold hover:bg-[#D8AE6E] transition-colors"
            >
              Checkout →
            </button>
          </div>
        </div>
      )}

      {/* Cart/Checkout Sheet */}
      {showCart && (
        <div className="fixed inset-0 z-50 flex justify-end">
          <div
            className="absolute inset-0 bg-[#1F1B16]/60"
            onClick={() => {
              if (!qrData) setShowCart(false);
            }}
          />
          <div className="relative w-full max-w-md bg-[#F6F3EC] h-full flex flex-col shadow-xl">
            <div className="flex items-center justify-between p-4 border-b border-[#1F1B16]/10">
              <h2 className="text-lg font-[family-name:var(--font-display)] text-[#1F1B16]">
                {qrData ? "Scan QR untuk Bayar" : "Checkout"}
              </h2>
              {!qrData && (
                <button
                  onClick={() => setShowCart(false)}
                  className="p-2 hover:bg-[#EFE9DA] rounded-full transition-colors"
                >
                  <X size={20} className="text-[#1F1B16]" />
                </button>
              )}
            </div>

            {/* QR Payment view — shown after order creation */}
            {qrData ? (
              <div className="flex-1 overflow-y-auto p-6">
                <QRPaymentBlock
                  qr={
                    qrData as unknown as Parameters<
                      typeof QRPaymentBlock
                    >[0]["qr"]
                  }
                  buyerEmail={checkoutEmail}
                  onPaid={onQrPaid}
                />
              </div>
            ) : (
              <>
                <div className="flex-1 overflow-y-auto p-4 space-y-3">
                  <div className="bg-[#EFE9DA] rounded-xl p-4 mb-4 border border-[#B8863C]/20">
                    <p className="text-sm text-[#7A5A22] font-[family-name:var(--font-mono)] mb-1">
                      {gallery.title}
                    </p>
                    <p className="text-2xl font-[family-name:var(--font-display)] text-[#1F1B16]">
                      {formatRupiah(totalPrice)}
                    </p>
                    <p className="text-xs text-[#4A443C] font-[family-name:var(--font-mono)] mt-1">
                      {cart.size} foto × {formatRupiah(gallery.price_per_photo)}
                    </p>
                  </div>

                  {photos
                    .filter((p) => cart.has(p.id))
                    .map((photo) => (
                      <div key={photo.id} className="flex items-center gap-3">
                        <img
                          src={photo.preview_url}
                          alt={photo.filename}
                          className="w-14 h-14 object-cover rounded-lg"
                        />
                        <div className="flex-1 min-w-0">
                          <p className="text-sm font-medium text-[#1F1B16] truncate">
                            {photo.filename}
                          </p>
                          <p className="text-xs text-[#4A443C] font-[family-name:var(--font-mono)]">
                            {formatRupiah(gallery.price_per_photo)}
                          </p>
                        </div>
                        <button
                          onClick={() => toggleCart(photo.id)}
                          className="text-[#4A443C]/50 hover:text-[#A8433F] p-1 transition-colors"
                        >
                          <X size={16} />
                        </button>
                      </div>
                    ))}
                </div>

                <div className="p-4 border-t border-[#1F1B16]/10 space-y-3">
                  <div>
                    <label className="text-sm font-medium mb-1.5 block text-[#1F1B16]">
                      Email untuk link download *
                    </label>
                    <input
                      id="checkout-email"
                      type="email"
                      placeholder="email@anda.com"
                      value={checkoutEmail}
                      onChange={(e) => setCheckoutEmail(e.target.value)}
                      className="w-full border border-[#1F1B16]/15 bg-white rounded-xl px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-[#B8863C]"
                    />
                    <p className="text-xs text-[#4A443C]/70 mt-1">
                      QR QRIS akan muncul setelah konfirmasi. Link download
                      dikirim ke email ini.
                    </p>
                  </div>
                  <button
                    id="pay-now-btn"
                    type="button"
                    onClick={handleCheckout}
                    disabled={
                      checkoutLoading ||
                      cart.size === 0 ||
                      !checkoutEmail.trim()
                    }
                    className="w-full py-3.5 bg-[#1F1B16] text-[#F6F3EC] rounded-xl font-semibold hover:bg-[#B8863C] disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
                  >
                    {checkoutLoading
                      ? "Membuat pesanan..."
                      : `Lanjut ke Pembayaran QRIS — ${formatRupiah(totalPrice)}`}
                  </button>
                </div>
              </>
            )}
          </div>
        </div>
      )}

      {/* Photo Lightbox */}
      {selectedPhoto && (
        <div
          className="fixed inset-0 z-50 bg-[#1F1B16]/95 flex items-center justify-center"
          onClick={() => setSelectedPhoto(null)}
        >
          <button
            className="absolute top-4 right-4 text-[#F6F3EC] hover:text-[#B8863C] transition-colors"
            onClick={() => setSelectedPhoto(null)}
          >
            <X size={28} />
          </button>
          <div
            className="max-w-4xl max-h-[90vh] p-4"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="relative">
              <img
                src={selectedPhoto.preview_url}
                alt={selectedPhoto.filename}
                draggable={false}
                onContextMenu={(e) => e.preventDefault()}
                className="max-w-full max-h-[80vh] object-contain rounded-lg select-none"
              />
              <WatermarkOverlay />
            </div>
            <div className="mt-4 flex items-center justify-between">
              <p className="text-[#F6F3EC]/70 text-sm font-[family-name:var(--font-mono)]">
                {selectedPhoto.filename}
              </p>
              <button
                onClick={() => {
                  toggleCart(selectedPhoto.id);
                  setSelectedPhoto(null);
                }}
                className={`px-6 py-2 rounded-full font-medium text-sm transition-colors ${
                  cart.has(selectedPhoto.id)
                    ? "bg-[#A8433F] text-white"
                    : "bg-[#B8863C] text-[#1F1B16] hover:bg-[#D8AE6E]"
                }`}
              >
                {cart.has(selectedPhoto.id)
                  ? "Batalkan Pilihan"
                  : `+ Pilih — ${formatRupiah(gallery.price_per_photo)}`}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
