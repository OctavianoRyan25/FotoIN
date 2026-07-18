"use client";

import { Suspense, useEffect, useMemo, useState } from "react";
import { useSearchParams } from "next/navigation";
import { Search, SearchX } from "lucide-react";
import {
  FONT_CLASSES,
  Header,
  Footer,
  EventCard,
  Event,
} from "../components/site-ui";

const API_URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api/v1";

function OrdersPageContent() {
  const searchParams = useSearchParams();
  const [query, setQuery] = useState(searchParams.get("q") ?? "");
  const [events, setEvents] = useState<Event[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchEvents = async () => {
      setLoading(true);
      try {
        const res = await fetch(`${API_URL}/public/events?limit=100`);
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

  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase();
    if (!q) return events;
    return events.filter(
      (e) =>
        e.title.toLowerCase().includes(q) ||
        (e.location ?? "").toLowerCase().includes(q) ||
        (e.photographer?.name ?? "").toLowerCase().includes(q),
    );
  }, [events, query]);

  return (
    <main className={`${FONT_CLASSES} min-h-screen flex flex-col bg-[#F6F3EC]`}>
      <Header active="events" />

      {/* Page intro + search */}
      <section className="px-4 pt-14 pb-10 text-center border-b border-[#1F1B16]/8">
        <p className="uppercase tracking-[0.25em] text-xs font-[family-name:var(--font-mono)] text-[#23474A] mb-3">
          Jelajahi
        </p>
        <h1 className="text-3xl sm:text-5xl font-[family-name:var(--font-display)] text-[#1F1B16] mb-8">
          Semua Event
        </h1>
        <div className="relative max-w-xl mx-auto">
          <Search
            size={17}
            className="absolute left-4 top-1/2 -translate-y-1/2 text-[#4A443C]/50"
          />
          <input
            id="events-search-input"
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Cari event, lokasi, atau fotografer..."
            className="w-full bg-white border border-[#1F1B16]/12 rounded-full pl-11 pr-4 py-3.5 text-sm text-[#1F1B16] placeholder:text-[#4A443C]/40 focus:outline-none focus:ring-2 focus:ring-[#B8863C] shadow-sm"
          />
        </div>
      </section>

      {/* Results */}
      <section className="flex-1 py-14 px-4">
        <div className="max-w-6xl mx-auto">
          {loading ? (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
              {Array.from({ length: 6 }).map((_, i) => (
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
          ) : filtered.length === 0 ? (
            <div className="flex flex-col items-center justify-center text-center py-20">
              <div className="w-14 h-14 rounded-full bg-[#EFE9DA] flex items-center justify-center mb-5">
                <SearchX size={22} className="text-[#B8863C]" />
              </div>
              <h2 className="text-xl font-[family-name:var(--font-display)] text-[#1F1B16] mb-2">
                Tidak ada event ditemukan
              </h2>
              <p className="text-[#4A443C] text-sm max-w-sm">
                Belum ada data event dari backend. Pastikan API berjalan dan
                event sudah tersedia.
              </p>
            </div>
          ) : (
            <>
              <p className="text-sm text-[#4A443C]/60 font-[family-name:var(--font-mono)] mb-8">
                {filtered.length} event ditemukan
              </p>
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
                {filtered.map((event) => (
                  <EventCard key={event.id} event={event} />
                ))}
              </div>
            </>
          )}
        </div>
      </section>

      <Footer />
    </main>
  );
}

export default function OrdersPage() {
  return (
    <Suspense
      fallback={
        <main
          className={`${FONT_CLASSES} min-h-screen flex items-center justify-center bg-[#F6F3EC]`}
        >
          <p className="text-[#23474A]">Memuat halaman pesanan...</p>
        </main>
      }
    >
      <OrdersPageContent />
    </Suspense>
  );
}
