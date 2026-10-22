"use client";
import Link from "next/link";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Aperture, MapPin, Calendar, LogIn, Scan, LogOut, User } from "lucide-react";
import { Fraunces, Inter, IBM_Plex_Mono } from "next/font/google";

// ---------------------------------------------------------------------------
// Fonts — Fraunces (display/serif), Inter (body), IBM Plex Mono (metadata).
// Move this block into app/layout.tsx and apply FONT_CLASSES to <body> once,
// instead of importing it per-page, to avoid re-loading fonts on every route.
// ---------------------------------------------------------------------------
export const fraunces = Fraunces({
  subsets: ["latin"],
  weight: ["400", "500", "600"],
  style: ["normal", "italic"],
  variable: "--font-display",
});
export const inter = Inter({ subsets: ["latin"], variable: "--font-body" });
export const plexMono = IBM_Plex_Mono({
  subsets: ["latin"],
  weight: ["400", "500"],
  variable: "--font-mono",
});
export const FONT_CLASSES = `${fraunces.variable} ${inter.variable} ${plexMono.variable} font-[family-name:var(--font-body)]`;

// ---------------------------------------------------------------------------
// Palette (for reference, used via Tailwind arbitrary values throughout):
// paper #F6F3EC · paper-2 #EFE9DA · ink #1F1B16 · gold #B8863C
// gold-light #D8AE6E · teal #23474A · rose #A8433F
// ---------------------------------------------------------------------------

export interface Photo {
  id: string;
  filename: string;
  preview_url: string;
  size_bytes: number;
  width: number;
  height: number;
}

export interface Gallery {
  id: string;
  title: string;
  slug: string;
  event_date: string;
  location: string;
  description: string;
  price_per_photo: number;
  photo_count: number;
  cover_photo_url?: string;
  photographer: { name: string };
  event?: {
    id: string;
    slug: string;
    title: string;
    location: string;
    event_date: string;
  };
}

export interface Event {
  id: string;
  photographer_id: string;
  title: string;
  slug: string;
  event_date?: string | null;
  location?: string;
  description?: string;
  cover_photo_url?: string;
  is_active: boolean;
  expires_at?: string | null;
  created_at: string;
  updated_at: string;
  photographer?: { name: string };
  galleries?: Gallery[];
}

export const formatRupiah = (amount: number) =>
  new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    minimumFractionDigits: 0,
  }).format(amount);

export const formatDate = (d: string) =>
  new Intl.DateTimeFormat("id-ID", {
    day: "numeric",
    month: "long",
    year: "numeric",
  }).format(new Date(d));

// Demo events — used whenever the API isn't reachable, so the site stays
// browsable end-to-end (home → events → event detail → gallery) without a backend.
export const DEMO_EVENTS: Event[] = [
  {
    id: "evt-1",
    photographer_id: "p1",
    title: "Sunmori Semarang",
    slug: "sunmori-semarang",
    event_date: "2026-06-14",
    location: "Semarang, Jawa Tengah",
    description:
      "Touring pagi komunitas motor klasik mengelilingi Kota Lama dan Simpang Lima.",
    cover_photo_url: "https://picsum.photos/seed/sunmori-semarang/700/500",
    is_active: true,
    created_at: "",
    updated_at: "",
    photographer: { name: "Rangga Photo" },
  },
  {
    id: "evt-2",
    photographer_id: "p2",
    title: "Lari 10K Borobudur Sunrise",
    slug: "lari-10k-borobudur",
    event_date: "2026-05-10",
    location: "Magelang, Jawa Tengah",
    description:
      "Lomba lari pagi dengan latar Candi Borobudur saat matahari terbit.",
    cover_photo_url: "https://picsum.photos/seed/lari-10k-borobudur/700/500",
    is_active: true,
    created_at: "",
    updated_at: "",
    photographer: { name: "Studio Kalem" },
  },
  {
    id: "evt-3",
    photographer_id: "p3",
    title: "Motocross Kejurnas Solo",
    slug: "motocross-kejurnas-solo",
    event_date: "2026-04-22",
    location: "Solo, Jawa Tengah",
    description:
      "Seri kejuaraan nasional motocross kelas pemula hingga profesional.",
    cover_photo_url:
      "https://picsum.photos/seed/motocross-kejurnas-solo/700/500",
    is_active: true,
    created_at: "",
    updated_at: "",
    photographer: { name: "Ardi Lens" },
  },
  {
    id: "evt-4",
    photographer_id: "p1",
    title: "Fun Bike Kota Lama",
    slug: "fun-bike-kota-lama",
    event_date: "2026-03-30",
    location: "Semarang, Jawa Tengah",
    description: "Gowes santai keliling kawasan heritage Kota Lama Semarang.",
    cover_photo_url: "https://picsum.photos/seed/fun-bike-kota-lama/700/500",
    is_active: true,
    created_at: "",
    updated_at: "",
    photographer: { name: "Rangga Photo" },
  },
  {
    id: "evt-5",
    photographer_id: "p4",
    title: "Triathlon Rawa Pening",
    slug: "triathlon-rawa-pening",
    event_date: "2026-07-05",
    location: "Ambarawa, Jawa Tengah",
    description:
      "Kombinasi renang, sepeda, dan lari mengelilingi Danau Rawa Pening.",
    cover_photo_url: "https://picsum.photos/seed/triathlon-rawa-pening/700/500",
    is_active: true,
    created_at: "",
    updated_at: "",
    photographer: { name: "Studio Kalem" },
  },
  {
    id: "evt-6",
    photographer_id: "p3",
    title: "Touring Dieng Sunrise",
    slug: "touring-dieng-sunrise",
    event_date: "2026-02-14",
    location: "Dieng, Wonosobo",
    description: "Konvoi motor menembus kabut dini hari menuju Bukit Sikunir.",
    cover_photo_url: "https://picsum.photos/seed/touring-dieng-sunrise/700/500",
    is_active: true,
    created_at: "",
    updated_at: "",
    photographer: { name: "Ardi Lens" },
  },
  {
    id: "evt-7",
    photographer_id: "p2",
    title: "Marathon Kota Lama",
    slug: "marathon-kota-lama-semarang",
    event_date: "2026-08-17",
    location: "Semarang, Jawa Tengah",
    description:
      "Lari 21K dan 42K memperingati Kemerdekaan, start dari Tugu Muda.",
    cover_photo_url: "https://picsum.photos/seed/marathon-kota-lama/700/500",
    is_active: true,
    created_at: "",
    updated_at: "",
    photographer: { name: "Rangga Photo" },
  },
  {
    id: "evt-8",
    photographer_id: "p4",
    title: "Drift Battle Bandengan",
    slug: "drift-battle-bandengan",
    event_date: "2026-01-25",
    location: "Jepara, Jawa Tengah",
    description: "Ajang unjuk skill drifting komunitas otomotif pesisir utara.",
    cover_photo_url:
      "https://picsum.photos/seed/drift-battle-bandengan/700/500",
    is_active: true,
    created_at: "",
    updated_at: "",
    photographer: { name: "Ardi Lens" },
  },
];

// Signature motif: a film-strip sprocket divider, echoed on every page.
export function FilmStrip({ tone = "ink" }: { tone?: "ink" | "paper" }) {
  const dot = tone === "ink" ? "#1F1B16" : "#F6F3EC";
  return (
    <div
      aria-hidden
      className="h-3 w-full shrink-0"
      style={{
        backgroundImage: `radial-gradient(circle, ${dot} 1.6px, transparent 1.6px)`,
        backgroundSize: "16px 100%",
        backgroundPosition: "center",
        opacity: tone === "ink" ? 0.18 : 0.35,
      }}
    />
  );
}

export function Logo({ small }: { small?: boolean }) {
  return (
    <Link href="/" className="flex items-center gap-2 shrink-0">
      <div
        className={`${small ? "w-8 h-8" : "w-9 h-9"} bg-[#1F1B16] rounded-full flex items-center justify-center`}
      >
        <Aperture
          size={small ? 16 : 18}
          className="text-[#D8AE6E]"
          strokeWidth={1.75}
        />
      </div>
      <span className="font-[family-name:var(--font-display)] text-[#1F1B16] font-medium text-lg">
        Foto<span className="text-[#B8863C] italic">IN</span>
      </span>
    </Link>
  );
}

export function Header({
  active,
  right,
}: {
  active?: "home" | "events" | "gallery" | "face-scan";
  right?: React.ReactNode;
}) {
  const router = useRouter();
  const [user, setUser] = useState<{ name: string; email: string } | null>(null);
  const [menuOpen, setMenuOpen] = useState(false);

  useEffect(() => {
    const raw = localStorage.getItem("client_user");
    if (raw) {
      try { setUser(JSON.parse(raw)); } catch { /* ignore */ }
    }
  }, []);

  const handleLogout = () => {
    localStorage.removeItem("client_token");
    localStorage.removeItem("client_user");
    setUser(null);
    setMenuOpen(false);
    router.push("/");
  };

  const linkClass = (key: string) =>
    `text-sm font-medium transition-colors ${
      active === key
        ? "text-[#1F1B16]"
        : "text-[#4A443C]/70 hover:text-[#1F1B16]"
    }`;

  return (
    <header className="sticky top-0 z-40 bg-[#F6F3EC]/95 backdrop-blur-sm border-b border-[#1F1B16]/8 px-4 py-3">
      <div className="max-w-6xl mx-auto flex items-center justify-between gap-4">
        <Logo small />
        <nav className="hidden sm:flex items-center gap-7">
          <Link href="/" className={linkClass("home")}>
            Beranda
          </Link>
          <Link href="/events" className={linkClass("events")}>
            Event
          </Link>
          <Link href="/gallery" className={linkClass("gallery")}>
            Galeri
          </Link>
          <Link href="/face-scan" className={`flex items-center gap-1.5 ${linkClass("face-scan")}`}>
            <Scan size={13} />
            Scan Wajah
          </Link>
          <a
            href="http://localhost:5173"
            className="text-sm font-medium text-[#4A443C]/70 hover:text-[#1F1B16] transition-colors"
          >
            Saya Fotografer
          </a>
        </nav>

        <div className="flex items-center gap-3">
          {right}
          {user ? (
            <div className="relative">
              <button
                id="header-user-menu"
                onClick={() => setMenuOpen(o => !o)}
                className="flex items-center gap-2 px-3 py-2 rounded-full border border-[#1F1B16]/12 hover:border-[#B8863C]/40 transition-all"
              >
                <div className="w-6 h-6 rounded-full bg-[#B8863C] flex items-center justify-center">
                  <User size={13} className="text-white" />
                </div>
                <span className="text-sm font-medium text-[#1F1B16] max-w-[100px] truncate">
                  {user.name.split(" ")[0]}
                </span>
              </button>

              {menuOpen && (
                <div className="absolute right-0 top-full mt-2 w-48 bg-white rounded-2xl border border-[#1F1B16]/8 shadow-lg py-2 z-50">
                  <div className="px-4 py-2 border-b border-[#1F1B16]/8">
                    <p className="text-xs font-medium text-[#1F1B16] truncate">{user.name}</p>
                    <p className="text-xs text-[#4A443C]/60 truncate">{user.email}</p>
                  </div>
                  <Link
                    href="/orders"
                    onClick={() => setMenuOpen(false)}
                    className="flex items-center gap-2 px-4 py-2.5 text-sm text-[#4A443C] hover:text-[#1F1B16] hover:bg-[#F6F3EC] transition-colors"
                  >
                    Pesanan Saya
                  </Link>
                  <Link
                    href="/profile"
                    onClick={() => setMenuOpen(false)}
                    className="flex items-center gap-2 px-4 py-2.5 text-sm text-[#4A443C] hover:text-[#1F1B16] hover:bg-[#F6F3EC] transition-colors"
                  >
                    <User size={13} />
                    Profil & Akun
                  </Link>
                  <Link
                    href="/face-scan"
                    onClick={() => setMenuOpen(false)}
                    className="flex items-center gap-2 px-4 py-2.5 text-sm text-[#4A443C] hover:text-[#1F1B16] hover:bg-[#F6F3EC] transition-colors"
                  >
                    <Scan size={13} />
                    Scan Wajah
                  </Link>
                  <button
                    id="header-logout-btn"
                    onClick={handleLogout}
                    className="w-full flex items-center gap-2 px-4 py-2.5 text-sm text-red-600 hover:bg-red-50 transition-colors"
                  >
                    <LogOut size={13} />
                    Keluar
                  </button>
                </div>
              )}
            </div>
          ) : (
            <Link
              id="header-login-btn"
              href="/login"
              className="flex items-center gap-1.5 px-4 py-2 rounded-full bg-[#1F1B16] text-[#F6F3EC] text-sm font-medium hover:bg-[#B8863C] transition-all duration-200"
            >
              <LogIn size={13} />
              Masuk
            </Link>
          )}
        </div>
      </div>
    </header>
  );
}

export function Footer() {
  return (
    <footer className="bg-[#1F1B16] py-10 px-4">
      <div className="max-w-6xl mx-auto flex flex-col sm:flex-row items-center justify-between gap-4">
        <div className="flex items-center gap-2 text-[#F6F3EC]">
          <Aperture size={16} className="text-[#B8863C]" strokeWidth={1.75} />
          <span className="font-[family-name:var(--font-display)]">
            Foto<span className="text-[#B8863C] italic">IN</span>
          </span>
        </div>
        <p className="text-[#F6F3EC]/50 text-sm font-[family-name:var(--font-mono)]">
          © {new Date().getFullYear()} · Platform Fotografi Event Indonesia
        </p>
      </div>
    </footer>
  );
}

export function EventCard({ event }: { event: Event }) {
  const galleryCount = event.galleries?.length;
  return (
    <Link
      href={`/events/${event.slug}`}
      className="group block rounded-2xl overflow-hidden border border-[#1F1B16]/8 bg-white hover:border-[#B8863C]/40 hover:-translate-y-1 transition-all duration-300"
    >
      <div className="relative aspect-[4/3] overflow-hidden bg-[#EFE9DA]">
        <img
          src={
            event.cover_photo_url ||
            `https://picsum.photos/seed/${event.slug}/600/450`
          }
          alt={event.title}
          className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
          loading="lazy"
        />
        <div className="absolute inset-0 bg-gradient-to-t from-[#1F1B16]/75 via-[#1F1B16]/5 to-transparent" />
        {typeof galleryCount === "number" && galleryCount > 0 && (
          <span className="absolute top-3 right-3 px-2 py-1 rounded-full bg-[#1F1B16]/70 backdrop-blur-sm text-[10px] font-[family-name:var(--font-mono)] text-[#F6F3EC]">
            {galleryCount} galeri
          </span>
        )}
        <div className="absolute bottom-0 left-0 right-0 p-4">
          <h3 className="text-white font-[family-name:var(--font-display)] text-lg leading-snug">
            {event.title}
          </h3>
        </div>
      </div>
      <div className="p-4 flex flex-col gap-1.5">
        <span className="flex items-center gap-1.5 text-xs text-[#4A443C] font-[family-name:var(--font-mono)]">
          <MapPin size={12} className="text-[#B8863C] shrink-0" />
          {event.location || "Lokasi belum diisi"}
        </span>
        {event.event_date && (
          <span className="flex items-center gap-1.5 text-xs text-[#4A443C] font-[family-name:var(--font-mono)]">
            <Calendar size={12} className="text-[#B8863C] shrink-0" />
            {formatDate(event.event_date)}
          </span>
        )}
      </div>
    </Link>
  );
}
