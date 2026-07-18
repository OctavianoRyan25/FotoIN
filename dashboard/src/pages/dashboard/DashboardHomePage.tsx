import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import {
  TrendingUp,
  Images,
  HardDrive,
  Wallet,
  Plus,
  ArrowRight,
  Camera,
  Star,
  Sparkles,
  BarChart3,
} from "lucide-react";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Cell,
} from "recharts";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { useAuth } from "@/store/auth";
import { formatRupiah, formatBytes, formatDateTime } from "@/lib/utils";
import api from "@/lib/api";
import type { Gallery, WalletSummary, User } from "@/types";

interface DashboardStats {
  galleries: Gallery[];
  wallet: WalletSummary | null;
  user: User | null;
}

const CHART_COLOR = "#7c3aed"; // brand violet, matches accent used elsewhere on the page

export function DashboardHomePage() {
  const { user } = useAuth();
  const [stats, setStats] = useState<DashboardStats>({
    galleries: [],
    wallet: null,
    user: user,
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [galRes, walRes, userRes] = await Promise.all([
          api.get("/galleries?limit=5"),
          api.get("/wallet"),
          api.get("/me"),
        ]);
        setStats({
          galleries: galRes.data.data ?? [],
          wallet: walRes.data.data ?? null,
          user: userRes.data.data ?? null,
        });
      } catch {
        // Silently fail — user will see empty state
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  const storageLimit =
    stats.user?.subscription_tier === "premium" ? 107374182400 : 3221225472;
  const storagePercent = stats.user
    ? (stats.user.storage_used_bytes / storageLimit) * 100
    : 0;
  const metricCards = [
    {
      id: "total-balance",
      label: "Saldo Dompet",
      value: stats.wallet ? formatRupiah(stats.wallet.balance) : "Rp 0",
      icon: Wallet,
      color: "text-emerald-600",
      bg: "bg-emerald-50 dark:bg-emerald-900/20",
    },
    {
      id: "total-galleries",
      label: "Total Galeri",
      value: String(stats.galleries.length),
      icon: Images,
      color: "text-blue-600",
      bg: "bg-blue-50 dark:bg-blue-900/20",
    },
    {
      id: "storage-used",
      label: "Storage Digunakan",
      value: stats.user ? formatBytes(stats.user.storage_used_bytes) : "0 B",
      icon: HardDrive,
      color: storagePercent > 80 ? "text-destructive" : "text-brand",
      bg: storagePercent > 80 ? "bg-destructive/10" : "bg-brand/10",
    },
    {
      id: "recent-earnings",
      label: "Pendapatan Bulan Ini",
      value: "Rp 0", // Will be computed from transactions
      icon: TrendingUp,
      color: "text-amber-600",
      bg: "bg-amber-50 dark:bg-amber-900/20",
    },
  ];

  // Chart data derived from the galleries we already fetched — top 5 by photo count
  const chartData = [...stats.galleries]
    .sort((a, b) => b.photo_count - a.photo_count)
    .slice(0, 5)
    .map((g) => ({
      name: g.title.length > 14 ? `${g.title.slice(0, 14)}…` : g.title,
      photos: g.photo_count,
    }));

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">
            Selamat datang, {stats.user?.name?.split(" ")[0]} 👋
          </h1>
          <p className="text-muted-foreground mt-1">
            Ini ringkasan aktivitas galeri Anda
          </p>
        </div>
        <div className="flex items-center gap-3">
          <Badge
            variant={
              stats.user?.subscription_tier === "premium"
                ? "brand"
                : "secondary"
            }
          >
            {stats.user?.subscription_tier === "premium" ? (
              <>
                <Star size={10} className="mr-1" /> Premium
              </>
            ) : (
              "Basic Plan"
            )}
          </Badge>
          {stats.user?.subscription_tier === "basic" && (
            <Button
              size="sm"
              variant="outline"
              className="text-brand border-brand/20 hover:bg-brand/5 gap-1.5"
            >
              <Sparkles size={14} /> Upgrade Premium
            </Button>
          )}
        </div>
      </div>

      {/* Metric Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-4">
        {metricCards.map((card) => (
          <Card
            key={card.id}
            id={card.id}
            className="hover:shadow-md hover:-translate-y-0.5 transition-all"
          >
            <CardContent className="p-5">
              <div className="flex items-start justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">{card.label}</p>
                  <p className="text-2xl font-bold mt-1">
                    {loading ? "—" : card.value}
                  </p>
                </div>
                <div className={`p-2.5 rounded-lg ${card.bg}`}>
                  <card.icon size={20} className={card.color} />
                </div>
              </div>
              {card.id === "storage-used" && (
                <div className="mt-3">
                  <div className="h-1.5 bg-muted rounded-full overflow-hidden">
                    <div
                      className={`h-full rounded-full transition-all ${storagePercent > 80 ? "bg-destructive" : "bg-brand"}`}
                      style={{ width: `${Math.min(storagePercent, 100)}%` }}
                    />
                  </div>
                  <p className="text-xs text-muted-foreground mt-1">
                    {storagePercent.toFixed(1)}% dari{" "}
                    {user?.subscription_tier === "premium" ? "100 GB" : "3 GB"}
                  </p>
                </div>
              )}
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Recent Galleries + Quick Actions */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Recent Galleries */}
        <Card className="lg:col-span-2">
          <CardHeader className="pb-3">
            <div className="flex items-center justify-between">
              <CardTitle className="text-base">Galeri Terbaru</CardTitle>
              <Link to="/dashboard/galleries">
                <Button
                  variant="ghost"
                  size="sm"
                  className="gap-1 text-violet-600"
                >
                  Lihat semua <ArrowRight size={14} />
                </Button>
              </Link>
            </div>
          </CardHeader>
          <CardContent>
            {loading ? (
              <div className="space-y-3">
                {[1, 2, 3].map((i) => (
                  <div
                    key={i}
                    className="h-14 bg-muted rounded-lg animate-pulse"
                  />
                ))}
              </div>
            ) : stats.galleries.length === 0 ? (
              <div className="text-center py-12">
                <Camera
                  size={40}
                  className="mx-auto text-muted-foreground mb-3"
                />
                <p className="text-muted-foreground text-sm mb-4">
                  Belum ada galeri
                </p>
                <Link to="/dashboard/galleries">
                  <Button
                    size="sm"
                    id="create-first-gallery"
                    className="bg-violet-600 hover:bg-violet-700 text-white gap-1.5"
                  >
                    <Plus size={14} /> Buat Galeri Pertama
                  </Button>
                </Link>
              </div>
            ) : (
              <div className="space-y-2">
                {stats.galleries.map((gallery) => (
                  <Link
                    key={gallery.id}
                    to={`/dashboard/galleries/${gallery.id}`}
                    className="flex items-center gap-4 p-3 rounded-lg hover:bg-muted/50 transition-colors group"
                  >
                    <div className="w-10 h-10 rounded-lg bg-brand/10 flex items-center justify-center flex-shrink-0">
                      <Images size={16} className="text-brand" />
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium truncate">
                        {gallery.title}
                      </p>
                      <p className="text-xs text-muted-foreground">
                        {gallery.photo_count} foto ·{" "}
                        {formatRupiah(gallery.price_per_photo)}/foto
                      </p>
                    </div>
                    <div className="flex items-center gap-2">
                      <Badge
                        variant={gallery.is_active ? "success" : "secondary"}
                        className="text-xs"
                      >
                        {gallery.is_active ? "Aktif" : "Nonaktif"}
                      </Badge>
                      <ArrowRight
                        size={14}
                        className="text-muted-foreground opacity-0 group-hover:opacity-100 transition-opacity"
                      />
                    </div>
                  </Link>
                ))}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Quick Actions — styled like a contact-sheet strip: numbered frames, sprocket edge */}
        <Card className="relative overflow-hidden">
          <div
            className="absolute inset-y-0 left-0 w-4 border-r border-dashed border-muted-foreground/15"
            style={{
              backgroundImage:
                "radial-gradient(circle, hsl(var(--muted-foreground) / 0.25) 1.4px, transparent 1.4px)",
              backgroundSize: "100% 16px",
              backgroundPosition: "center 8px",
            }}
          />
          <CardHeader className="pl-8">
            <p className="text-[10px] font-mono uppercase tracking-[0.15em] text-muted-foreground/70 mb-0.5">
              Roll 01 · Aksi
            </p>
            <CardTitle className="text-base">Aksi Cepat</CardTitle>
          </CardHeader>
          <CardContent className="p-0 pb-2">
            {[
              {
                id: "quick-create-gallery",
                frame: "01",
                href: "/dashboard/galleries",
                icon: Plus,
                label: "Buat Galeri Baru",
                caption: "Mulai jual foto dari event baru",
                accent: "text-violet-600",
                accentBg: "bg-violet-50 dark:bg-violet-900/20",
                accentBorder: "border-violet-200 dark:border-violet-800/60",
                hoverBg: "hover:bg-violet-50/60 dark:hover:bg-violet-900/10",
              },
              {
                id: "quick-wallet",
                frame: "02",
                href: "/dashboard/wallet",
                icon: Wallet,
                label: "Dompet & Penarikan",
                caption: "Cek saldo dan tarik dana",
                accent: "text-emerald-600",
                accentBg: "bg-emerald-50 dark:bg-emerald-900/20",
                accentBorder: "border-emerald-200 dark:border-emerald-800/60",
                hoverBg: "hover:bg-emerald-50/60 dark:hover:bg-emerald-900/10",
              },
            ].map((action, i, arr) => (
              <Link
                key={action.id}
                to={action.href}
                id={action.id}
                className={`group flex items-center gap-3 pl-8 pr-4 py-3.5 transition-colors ${action.hoverBg} ${
                  i < arr.length - 1
                    ? "border-b border-dashed border-muted-foreground/15"
                    : ""
                }`}
              >
                <span
                  className={`font-mono text-[10px] leading-none w-6 h-6 flex-shrink-0 rounded-sm border flex items-center justify-center ${action.accentBorder} ${action.accent}`}
                >
                  {action.frame}
                </span>
                <span
                  className={`w-9 h-9 rounded-lg flex items-center justify-center flex-shrink-0 ${action.accentBg}`}
                >
                  <action.icon size={16} className={action.accent} />
                </span>
                <span className="flex-1 min-w-0">
                  <span className="block text-sm font-medium">
                    {action.label}
                  </span>
                  <span className="block text-xs text-muted-foreground truncate">
                    {action.caption}
                  </span>
                </span>
                <ArrowRight
                  size={14}
                  className="text-muted-foreground/40 group-hover:text-foreground group-hover:translate-x-0.5 transition-all flex-shrink-0"
                />
              </Link>
            ))}
          </CardContent>
        </Card>
      </div>

      {/* Photo Distribution Chart */}
      <Card>
        <CardHeader className="pb-3">
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="text-base flex items-center gap-2">
                <BarChart3 size={16} className="text-brand" />
                Distribusi Foto per Galeri
              </CardTitle>
              <CardDescription>
                5 galeri teratas berdasarkan jumlah foto
              </CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {loading ? (
            <div className="h-52 bg-muted rounded-lg animate-pulse" />
          ) : chartData.length === 0 ? (
            <div className="text-center py-10">
              <BarChart3
                size={32}
                className="mx-auto text-muted-foreground/40 mb-2"
              />
              <p className="text-muted-foreground text-sm">
                Belum ada data untuk ditampilkan
              </p>
            </div>
          ) : (
            <div className="h-52">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart
                  data={chartData}
                  margin={{ top: 4, right: 8, left: -16, bottom: 0 }}
                >
                  <CartesianGrid
                    strokeDasharray="3 3"
                    vertical={false}
                    className="stroke-muted"
                  />
                  <XAxis
                    dataKey="name"
                    tick={{
                      fontSize: 12,
                      fill: "hsl(var(--muted-foreground))",
                    }}
                    axisLine={false}
                    tickLine={false}
                  />
                  <YAxis
                    allowDecimals={false}
                    tick={{
                      fontSize: 12,
                      fill: "hsl(var(--muted-foreground))",
                    }}
                    axisLine={false}
                    tickLine={false}
                    width={32}
                  />
                  <Tooltip
                    cursor={{ fill: "hsl(var(--muted))", opacity: 0.4 }}
                    contentStyle={{
                      borderRadius: 8,
                      border: "1px solid hsl(var(--border))",
                      fontSize: 12,
                    }}
                    formatter={(value: number) => [`${value} foto`, "Jumlah"]}
                  />
                  <Bar dataKey="photos" radius={[6, 6, 0, 0]} maxBarSize={48}>
                    {chartData.map((_, i) => (
                      <Cell
                        key={i}
                        fill={CHART_COLOR}
                        fillOpacity={1 - i * 0.12}
                      />
                    ))}
                  </Bar>
                </BarChart>
              </ResponsiveContainer>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
