"use client";

import { useEffect, useState, FormEvent } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Search, ArrowRight, MapPin, Calendar, ImageIcon } from "lucide-react";
import {
  FONT_CLASSES,
  FilmStrip,
  Header,
  Footer,
  Gallery,
  formatDate,
} from "../components/site-ui";

const API_URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api/v1";

export default function GalleryListPage() {
  const router = useRouter();
  const [query, setQuery] = useState("");
  const [galleries, setGalleries] = useState<Gallery[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchGalleries = async () => {
      setLoading(true);
      try {
        const res = await fetch(`${API_URL}/public/galleries?limit=20`);
        if (!res.ok) throw new Error("failed");
        const data = await res.json();
        const list: Gallery[] = Array.isArray(data.data) ? data.data : [];
        setGalleries(list);
      } catch {
        setGalleries([]);
      } finally {
        setLoading(false);
      }
    };
    fetchGalleries();
  }, []);

  const handleSearch = (e: FormEvent) => {
    e.preventDefault();
    if (query.trim()) {
      // client-side filtering or you can implement backend search.
      // For now, let's just alert or implement basic client filter.
    }
  };

  const filteredGalleries = query.trim()
    ? galleries.filter(
        (g) =>
          g.title.toLowerCase().includes(query.toLowerCase()) ||
          g.location?.toLowerCase().includes(query.toLowerCase())
      )
    : galleries;

  return (
    <main className={`${FONT_CLASSES} min-h-screen flex flex-col bg-[#F6F3EC]`}>
      <Header active="gallery" />

      {/* Hero section */}
      <section className="px-4 pt-16 pb-12 sm:pt-20 sm:pb-16 text-center">
        <div className="max-w-3xl mx-auto">
          <p className="uppercase tracking-[0.25em] text-xs font-[family-name:var(--font-mono)] text-[#23474A] mb-4">
            Eksplorasi Galeri
          </p>
          <h1 className="text-4xl sm:text-5xl font-[family-name:var(--font-display)] font-medium text-[#1F1B16] leading-[1.1] mb-5">
            Temukan Galeri <em className="italic text-[#B8863C]">Fotografi</em>
          </h1>
          <p className="text-[#4A443C] max-w-xl mx-auto mb-9 leading-relaxed">
            Jelajahi ribuan foto dari berbagai event dan momen spesial, temukan dirimu dan dapatkan foto terbaikmu.
          </p>

          {/* Search */}
          <form
            onSubmit={handleSearch}
            className="flex flex-col sm:flex-row gap-2.5 max-w-xl mx-auto"
          >
            <div className="relative flex-1">
              <Search
                size={17}
                className="absolute left-4 top-1/2 -translate-y-1/2 text-[#4A443C]/50"
              />
              <input
                id="gallery-search"
                type="text"
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                placeholder="Cari nama galeri atau lokasi..."
                className="w-full bg-white border border-[#1F1B16]/12 rounded-full pl-11 pr-4 py-3.5 text-sm text-[#1F1B16] placeholder:text-[#4A443C]/40 focus:outline-none focus:ring-2 focus:ring-[#B8863C] shadow-sm"
              />
            </div>
            <button
              id="search-gallery-cta"
              type="submit"
              className="px-7 py-3.5 bg-[#1F1B16] text-[#F6F3EC] rounded-full font-medium text-sm hover:bg-[#B8863C] transition-colors duration-300 shrink-0"
            >
              Cari Galeri
            </button>
          </form>
        </div>
      </section>

      <FilmStrip tone="ink" />

      {/* Gallery List */}
      <section className="py-14 px-4 flex-1">
        <div className="max-w-6xl mx-auto">
          {loading ? (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
              {Array.from({ length: 6 }).map((_, i) => (
                <div
                  key={i}
                  className="rounded-2xl overflow-hidden border border-[#1F1B16]/8 bg-white"
                >
                  <div className="aspect-[4/3] bg-[#EFE9DA] animate-pulse" />
                  <div className="p-4 space-y-2">
                    <div className="h-4 w-3/4 bg-[#EFE9DA] rounded animate-pulse" />
                    <div className="h-3 w-1/2 bg-[#EFE9DA] rounded animate-pulse mt-2" />
                  </div>
                </div>
              ))}
            </div>
          ) : filteredGalleries.length === 0 ? (
            <div className="rounded-2xl border border-[#1F1B16]/10 bg-white p-8 text-center text-[#4A443C]">
              {galleries.length === 0 
                ? "Belum ada galeri yang tersedia dari backend." 
                : "Tidak ada galeri yang cocok dengan pencarian."}
            </div>
          ) : (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
              {filteredGalleries.map((gallery) => (
                <Link
                  key={gallery.slug}
                  href={`/gallery/${gallery.slug}`}
                  className="group block rounded-2xl overflow-hidden border border-[#1F1B16]/8 bg-white hover:border-[#B8863C]/40 hover:-translate-y-1 transition-all duration-300"
                >
                  <div className="relative aspect-[4/3] overflow-hidden bg-[#EFE9DA]">
                    <img
                      src={
                        gallery.cover_photo_url ||
                        `https://picsum.photos/seed/${gallery.slug}/600/450`
                      }
                      alt={gallery.title}
                      className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
                      loading="lazy"
                    />
                    <div className="absolute inset-0 bg-gradient-to-t from-[#1F1B16]/75 via-[#1F1B16]/5 to-transparent" />
                    
                    {gallery.photographer?.name && (
                      <span className="absolute top-3 left-3 px-2 py-1 rounded-full bg-[#1F1B16]/70 backdrop-blur-sm text-[10px] font-[family-name:var(--font-mono)] text-[#F6F3EC]">
                        oleh {gallery.photographer.name}
                      </span>
                    )}
                    
                    <div className="absolute bottom-0 left-0 right-0 p-4 flex justify-between items-end">
                      <h3 className="text-white font-[family-name:var(--font-display)] text-lg leading-snug flex-1">
                        {gallery.title}
                      </h3>
                    </div>
                  </div>
                  <div className="p-4 flex flex-col gap-2">
                    <div className="flex items-center gap-1.5 text-xs text-[#4A443C] font-[family-name:var(--font-mono)]">
                      <MapPin size={12} className="text-[#B8863C] shrink-0" />
                      {gallery.location || "Lokasi belum diisi"}
                    </div>
                    {gallery.event_date && (
                      <div className="flex items-center gap-1.5 text-xs text-[#4A443C] font-[family-name:var(--font-mono)]">
                        <Calendar size={12} className="text-[#B8863C] shrink-0" />
                        {formatDate(gallery.event_date)}
                      </div>
                    )}
                    <div className="flex items-center gap-1.5 text-xs text-[#4A443C] font-[family-name:var(--font-mono)]">
                        <ImageIcon size={12} className="text-[#B8863C] shrink-0" />
                        {gallery.photo_count || 0} foto tersedia
                    </div>
                  </div>
                </Link>
              ))}
            </div>
          )}
        </div>
      </section>

      <Footer />
    </main>
  );
}
