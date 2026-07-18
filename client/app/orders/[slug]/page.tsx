"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { MapPin, Calendar, ArrowLeft, ImageIcon } from "lucide-react";
import {
  FONT_CLASSES,
  Header,
  Footer,
  FilmStrip,
  formatDate,
  Event,
} from "../../components/site-ui";

const API_URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api/v1";

// Demo galleries for an event — slugs map directly onto /gallery/[slug], which
// falls back to its own demo mode for any slug the backend doesn't recognize.
function demoGalleriesFor(event: Event) {
  return [
    {
      slug: event.slug,
      label: "Galeri Utama",
      desc: "Semua momen dari sesi utama event ini.",
    },
    {
      slug: `${event.slug}-b`,
      label: "Sesi Sore",
      desc: "Dokumentasi tambahan menjelang penutupan.",
    },
  ];
}

export default function EventDetailPage() {
  const { slug } = useParams<{ slug: string }>();
  const [event, setEvent] = useState<Event | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchEvent = async () => {
      setLoading(true);
      try {
        const res = await fetch(`${API_URL}/public/events/${slug}`);
        if (!res.ok) throw new Error("not found");
        const data = await res.json();
        setEvent(data.data);
      } catch {
        setEvent(null);
      } finally {
        setLoading(false);
      }
    };
    if (slug) fetchEvent();
  }, [slug]);

  if (loading || !event) {
    return (
      <main className={`${FONT_CLASSES} min-h-screen bg-[#F6F3EC]`}>
        <Header active="events" />
        <div className="max-w-5xl mx-auto px-4 py-14">
          <div className="h-64 w-full rounded-3xl bg-[#EFE9DA] animate-pulse mb-8" />
          <div className="h-8 w-72 bg-[#EFE9DA] rounded animate-pulse mb-3" />
          <div className="h-4 w-48 bg-[#EFE9DA] rounded animate-pulse" />
        </div>
      </main>
    );
  }

  const galleries = event.galleries?.length
    ? event.galleries.map((g) => ({
        slug: g.slug,
        label: g.title,
        desc: g.description,
      }))
    : [];

  return (
    <main className={`${FONT_CLASSES} min-h-screen flex flex-col bg-[#F6F3EC]`}>
      <Header active="events" />

      <div className="max-w-5xl mx-auto px-4 pt-8 w-full">
        <Link
          href="/events"
          className="inline-flex items-center gap-1.5 text-xs font-[family-name:var(--font-mono)] text-[#4A443C]/70 hover:text-[#B8863C] transition-colors mb-6"
        >
          <ArrowLeft size={12} /> Semua Event
        </Link>
      </div>

      {/* Cover */}
      <div className="max-w-5xl mx-auto px-4 w-full mb-8">
        <div className="relative aspect-[21/9] rounded-3xl overflow-hidden bg-[#EFE9DA]">
          <img
            src={
              event.cover_photo_url ||
              `https://picsum.photos/seed/${event.slug}/1200/500`
            }
            alt={event.title}
            className="w-full h-full object-cover"
          />
          <div className="absolute inset-0 bg-gradient-to-t from-[#1F1B16]/80 via-[#1F1B16]/10 to-transparent" />
          <div className="absolute bottom-0 left-0 right-0 p-6 sm:p-8">
            <h1 className="text-2xl sm:text-4xl font-[family-name:var(--font-display)] text-white mb-3 max-w-2xl">
              {event.title}
            </h1>
            <div className="flex flex-wrap items-center gap-x-4 gap-y-1.5 text-sm text-white/85 font-[family-name:var(--font-mono)]">
              {event.location && (
                <span className="flex items-center gap-1.5">
                  <MapPin size={14} className="text-[#D8AE6E]" />
                  {event.location}
                </span>
              )}
              {event.event_date && (
                <span className="flex items-center gap-1.5">
                  <Calendar size={14} className="text-[#D8AE6E]" />
                  {formatDate(event.event_date)}
                </span>
              )}
              {event.photographer?.name && (
                <span>oleh {event.photographer.name}</span>
              )}
            </div>
          </div>
        </div>
        {event.description && (
          <p className="text-[#4A443C] text-sm sm:text-base max-w-2xl mt-6 leading-relaxed">
            {event.description}
          </p>
        )}
      </div>

      <FilmStrip tone="ink" />

      {/* Galleries */}
      <section className="py-14 px-4">
        <div className="max-w-5xl mx-auto">
          <h2 className="text-2xl font-[family-name:var(--font-display)] text-[#1F1B16] mb-8">
            Galeri di Event Ini
          </h2>
          {galleries.length === 0 ? (
            <div className="rounded-2xl border border-[#1F1B16]/10 bg-white p-8 text-center text-[#4A443C]">
              Belum ada galeri yang tersedia untuk event ini dari backend.
            </div>
          ) : (
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
              {galleries.map((g) => (
                <Link
                  key={g.slug}
                  href={`/gallery/${g.slug}`}
                  className="group flex items-center gap-4 p-4 rounded-2xl border border-[#1F1B16]/8 bg-white hover:border-[#B8863C]/40 hover:-translate-y-0.5 transition-all duration-300"
                >
                  <div className="w-20 h-20 rounded-xl overflow-hidden shrink-0 bg-[#EFE9DA]">
                    <img
                      src={`https://picsum.photos/seed/${g.slug}-0/200/200`}
                      alt={g.label}
                      className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
                    />
                  </div>
                  <div className="min-w-0">
                    <h3 className="font-[family-name:var(--font-display)] text-[#1F1B16] text-lg truncate">
                      {g.label}
                    </h3>
                    <p className="text-xs text-[#4A443C] mt-1 line-clamp-2">
                      {g.desc || "Lihat dan pilih foto dari sesi ini."}
                    </p>
                    <span className="inline-flex items-center gap-1 text-[11px] font-[family-name:var(--font-mono)] text-[#B8863C] mt-2">
                      <ImageIcon size={11} /> Buka galeri
                    </span>
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
