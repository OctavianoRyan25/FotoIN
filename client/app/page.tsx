"use client";

import { useEffect, useState, FormEvent } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Search, ArrowRight } from "lucide-react";
import {
  FONT_CLASSES,
  FilmStrip,
  Header,
  Footer,
  EventCard,
  Event,
} from "./components/site-ui";

const API_URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api/v1";
const QUICK_TAGS = ["Sunmori", "Lari 10K", "Motocross", "Triathlon"];

export default function Home() {
  const router = useRouter();
  const [query, setQuery] = useState("");
  const [events, setEvents] = useState<Event[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchEvents = async () => {
      setLoading(true);
      try {
        const res = await fetch(`${API_URL}/public/events?limit=6`);
        if (!res.ok) throw new Error("failed");
        const data = await res.json();
        const list: Event[] = Array.isArray(data.data) ? data.data : [];
        setEvents(list);
      } catch {
        setEvents([]);
      } finally {
        setLoading(false);
      }
    };
    fetchEvents();
  }, []);

  const handleSearch = (e: FormEvent) => {
    e.preventDefault();
    router.push(
      query.trim()
        ? `/events?q=${encodeURIComponent(query.trim())}`
        : "/events",
    );
  };

  return (
    <main className={`${FONT_CLASSES} min-h-screen flex flex-col bg-[#F6F3EC]`}>
      <Header active="home" />

      {/* Hero */}
      <section className="relative overflow-hidden px-4 pt-16 pb-20 sm:pt-20 sm:pb-24 text-center">
        <div
          aria-hidden
          className="pointer-events-none absolute -right-24 -top-24 w-[420px] h-[420px] rounded-full border border-[#B8863C]/20"
        />
        <div
          aria-hidden
          className="pointer-events-none absolute -left-32 bottom-0 w-[320px] h-[320px] rounded-full border border-[#23474A]/15"
        />

        <div className="relative max-w-3xl mx-auto">
          <p className="uppercase tracking-[0.25em] text-xs font-[family-name:var(--font-mono)] text-[#23474A] mb-4">
            Galeri Foto Event Indonesia
          </p>

          <h1 className="text-4xl sm:text-6xl font-[family-name:var(--font-display)] font-medium text-[#1F1B16] leading-[1.1] mb-5">
            Setiap momen, tersimpan{" "}
            <em className="italic text-[#B8863C]">selamanya</em>
          </h1>
          <p className="text-base sm:text-lg text-[#4A443C] max-w-xl mx-auto mb-9 leading-relaxed">
            Temukan dirimu di ribuan foto Sunmori, race day, dan event olahraga
            favorit — diambil oleh fotografer terbaik di kotamu.
          </p>

          {/* Search */}
          <form
            onSubmit={handleSearch}
            className="flex flex-col sm:flex-row gap-2.5 max-w-xl mx-auto mb-4"
          >
            <div className="relative flex-1">
              <Search
                size={17}
                className="absolute left-4 top-1/2 -translate-y-1/2 text-[#4A443C]/50"
              />
              <input
                id="event-search"
                type="text"
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                placeholder="Cari event, lokasi, atau fotografer..."
                className="w-full bg-white border border-[#1F1B16]/12 rounded-full pl-11 pr-4 py-3.5 text-sm text-[#1F1B16] placeholder:text-[#4A443C]/40 focus:outline-none focus:ring-2 focus:ring-[#B8863C] shadow-sm"
              />
            </div>
            <button
              id="search-events-cta"
              type="submit"
              className="px-7 py-3.5 bg-[#1F1B16] text-[#F6F3EC] rounded-full font-medium text-sm hover:bg-[#B8863C] transition-colors duration-300 shrink-0"
            >
              Cari Event
            </button>
          </form>

          <div className="flex flex-wrap items-center justify-center gap-2 mb-14">
            <span className="text-xs text-[#4A443C]/50 font-[family-name:var(--font-mono)] mr-1">
              Populer:
            </span>
            {QUICK_TAGS.map((tag) => (
              <button
                key={tag}
                onClick={() =>
                  router.push(`/events?q=${encodeURIComponent(tag)}`)
                }
                className="px-3 py-1.5 rounded-full text-xs font-medium bg-[#EFE9DA] text-[#4A443C] hover:bg-[#B8863C] hover:text-white transition-colors"
              >
                {tag}
              </button>
            ))}
          </div>

          {/* Stats */}
          <div className="inline-flex items-center gap-8 sm:gap-12 rounded-2xl bg-white/60 border border-[#1F1B16]/8 px-8 sm:px-12 py-6">
            {[
              { value: "10K+", label: "Foto Terjual" },
              { value: "200+", label: "Fotografer" },
              { value: "50+", label: "Event" },
            ].map((stat, i) => (
              <div
                key={stat.label}
                className="flex items-center gap-8 sm:gap-12"
              >
                <div>
                  <p className="text-2xl sm:text-3xl font-[family-name:var(--font-display)] text-[#B8863C]">
                    {stat.value}
                  </p>
                  <p className="text-[#4A443C]/70 text-[11px] font-[family-name:var(--font-mono)] uppercase tracking-wider mt-1">
                    {stat.label}
                  </p>
                </div>
                {i < 2 && <div className="w-px h-9 bg-[#1F1B16]/10" />}
              </div>
            ))}
          </div>
        </div>
      </section>

      <FilmStrip tone="ink" />

      {/* Popular Events */}
      <section className="py-20 px-4 bg-[#F6F3EC]">
        <div className="max-w-6xl mx-auto">
          <div className="flex items-end justify-between mb-10 gap-4">
            <div>
              <p className="uppercase tracking-[0.25em] text-xs font-[family-name:var(--font-mono)] text-[#23474A] mb-3">
                Trending Minggu Ini
              </p>
              <h2 className="text-3xl sm:text-4xl font-[family-name:var(--font-display)] text-[#1F1B16]">
                Event Populer
              </h2>
            </div>
            <Link
              href="/events"
              className="hidden sm:flex items-center gap-1.5 text-sm font-medium text-[#1F1B16] hover:text-[#B8863C] transition-colors shrink-0"
            >
              Lihat semua event <ArrowRight size={15} />
            </Link>
          </div>

          {loading ? (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
              {Array.from({ length: 3 }).map((_, i) => (
                <div
                  key={i}
                  className="rounded-2xl overflow-hidden border border-[#1F1B16]/8"
                >
                  <div className="aspect-[4/3] bg-[#EFE9DA] animate-pulse" />
                  <div className="p-4 space-y-2">
                    <div className="h-3 w-2/3 bg-[#EFE9DA] rounded animate-pulse" />
                    <div className="h-3 w-1/2 bg-[#EFE9DA] rounded animate-pulse" />
                  </div>
                </div>
              ))}
            </div>
          ) : events.length === 0 ? (
            <div className="rounded-2xl border border-[#1F1B16]/10 bg-white p-8 text-center text-[#4A443C]">
              Event belum tersedia dari backend. Pastikan server API sedang
              berjalan dan data event sudah dibuat.
            </div>
          ) : (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
              {events.map((event) => (
                <EventCard key={event.id} event={event} />
              ))}
            </div>
          )}

          <Link
            href="/events"
            className="sm:hidden mt-8 flex items-center justify-center gap-1.5 text-sm font-medium text-[#1F1B16]"
          >
            Lihat semua event <ArrowRight size={15} />
          </Link>
        </div>
      </section>

      <FilmStrip tone="ink" />

      {/* Cara Kerja */}
      <section className="py-20 px-4 bg-[#EFE9DA]">
        <div className="max-w-4xl mx-auto">
          <p className="text-center uppercase tracking-[0.25em] text-xs font-[family-name:var(--font-mono)] text-[#23474A] mb-3">
            Cara Kerja
          </p>
          <h2 className="text-3xl sm:text-4xl font-[family-name:var(--font-display)] text-center text-[#1F1B16] mb-14">
            Tiga langkah menuju fotomu
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {[
              {
                step: "01",
                title: "Temukan Galeri",
                desc: "Cari event lewat kolom pencarian, tautan dari fotografer, atau pindai kode QR di lokasi.",
              },
              {
                step: "02",
                title: "Pilih Foto Anda",
                desc: "Lihat pratinjau berwatermark, tambahkan foto favorit ke keranjang.",
              },
              {
                step: "03",
                title: "Bayar & Unduh",
                desc: "Bayar lewat transfer atau kartu, tautan unduh terkirim langsung ke email Anda.",
              },
            ].map((item) => (
              <div
                key={item.step}
                className="relative p-7 rounded-2xl bg-[#F6F3EC] border border-[#1F1B16]/8 hover:border-[#B8863C]/40 hover:-translate-y-1 transition-all duration-300"
              >
                <span className="font-[family-name:var(--font-mono)] text-sm text-[#B8863C]">
                  No. {item.step}
                </span>
                <h3 className="text-lg font-[family-name:var(--font-display)] text-[#1F1B16] mt-3 mb-2">
                  {item.title}
                </h3>
                <p className="text-[#4A443C] text-sm leading-relaxed">
                  {item.desc}
                </p>
              </div>
            ))}
          </div>
        </div>
      </section>

      <Footer />
    </main>
  );
}
