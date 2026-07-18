"use client";

import { useEffect, useState, FormEvent } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import {
  User,
  Lock,
  KeyRound,
  Trash2,
  LogOut,
  ChevronRight,
  AlertCircle,
  CheckCircle2,
  Eye,
  EyeOff,
  Loader2,
  Mail,
  ShieldAlert,
  Aperture,
} from "lucide-react";
import { toast } from "sonner";
import { FONT_CLASSES, Logo, Footer } from "../components/site-ui";

const API_URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api/v1";

// ── Types ─────────────────────────────────────────────────────────────────────

interface ClientUser {
  id: string;
  name: string;
  email: string;
  role: string;
  subscription_tier: string;
  is_verified: boolean;
  has_password?: boolean; // true if the user has set a password
}

// ── Small reusable atoms ──────────────────────────────────────────────────────

function SectionCard({
  title,
  description,
  icon: Icon,
  children,
}: {
  title: string;
  description: string;
  icon: React.ElementType;
  children: React.ReactNode;
}) {
  return (
    <div className="bg-white rounded-2xl border border-[#1F1B16]/8 shadow-sm overflow-hidden">
      <div className="px-6 py-5 border-b border-[#1F1B16]/8 flex items-start gap-3">
        <div className="w-9 h-9 rounded-xl bg-[#B8863C]/10 border border-[#B8863C]/20 flex items-center justify-center flex-shrink-0 mt-0.5">
          <Icon size={17} className="text-[#B8863C]" />
        </div>
        <div>
          <h2 className="font-semibold text-[#1F1B16] text-sm">{title}</h2>
          <p className="text-xs text-[#4A443C]/60 mt-0.5">{description}</p>
        </div>
      </div>
      <div className="px-6 py-5">{children}</div>
    </div>
  );
}

function InputField({
  id,
  label,
  type = "text",
  value,
  onChange,
  placeholder,
  disabled,
  showToggle,
  onToggle,
}: {
  id: string;
  label: string;
  type?: string;
  value: string;
  onChange: (v: string) => void;
  placeholder?: string;
  disabled?: boolean;
  showToggle?: boolean;
  onToggle?: () => void;
}) {
  return (
    <div>
      <label
        htmlFor={id}
        className="block text-xs font-medium text-[#4A443C] mb-1.5 font-[family-name:var(--font-mono)]"
      >
        {label}
      </label>
      <div className="relative">
        <input
          id={id}
          type={type}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder={placeholder}
          disabled={disabled}
          className="w-full px-3.5 py-2.5 rounded-xl border border-[#1F1B16]/12 text-sm text-[#1F1B16] placeholder:text-[#4A443C]/40 focus:outline-none focus:ring-2 focus:ring-[#B8863C]/30 focus:border-[#B8863C] disabled:bg-[#F6F3EC] disabled:cursor-not-allowed transition-all bg-white pr-10"
        />
        {showToggle && onToggle && (
          <button
            type="button"
            onClick={onToggle}
            className="absolute right-3 top-1/2 -translate-y-1/2 text-[#4A443C]/40 hover:text-[#4A443C] transition-colors"
          >
            {type === "text" ? <EyeOff size={15} /> : <Eye size={15} />}
          </button>
        )}
      </div>
    </div>
  );
}

function SubmitButton({
  loading,
  label,
  id,
}: {
  loading: boolean;
  label: string;
  id: string;
}) {
  return (
    <button
      id={id}
      type="submit"
      disabled={loading}
      className="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-[#1F1B16] text-[#F6F3EC] text-sm font-medium hover:bg-[#B8863C] disabled:opacity-60 disabled:cursor-not-allowed transition-all duration-200"
    >
      {loading && <Loader2 size={14} className="animate-spin" />}
      {label}
    </button>
  );
}

// ── Hooks ─────────────────────────────────────────────────────────────────────

function useClientAuth() {
  const [user, setUser] = useState<ClientUser | null>(null);
  const [token, setToken] = useState<string | null>(null);

  useEffect(() => {
    const t = localStorage.getItem("client_token");
    const u = localStorage.getItem("client_user");
    if (t && u) {
      setToken(t);
      try {
        setUser(JSON.parse(u));
      } catch {}
    }
  }, []);

  const refresh = async (t: string) => {
    try {
      const res = await fetch(`${API_URL}/me`, {
        headers: { Authorization: `Bearer ${t}` },
      });
      if (res.ok) {
        const data = await res.json();
        const updated = data.data as ClientUser;
        setUser(updated);
        localStorage.setItem("client_user", JSON.stringify(updated));
      }
    } catch {}
  };

  return { user, token, setUser, refresh };
}

// ── Sections ──────────────────────────────────────────────────────────────────

function ProfileSection({
  user,
  token,
  onUpdated,
}: {
  user: ClientUser;
  token: string;
  onUpdated: () => void;
}) {
  const [name, setName] = useState(user.name);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!name.trim() || name.trim().length < 2) {
      toast.error("Nama minimal 2 karakter");
      return;
    }
    setLoading(true);
    try {
      const res = await fetch(`${API_URL}/me/profile`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ name: name.trim() }),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.message ?? "Gagal memperbarui profil");
      // Update localStorage
      const stored = JSON.parse(localStorage.getItem("client_user") ?? "{}");
      localStorage.setItem(
        "client_user",
        JSON.stringify({ ...stored, name: name.trim() }),
      );
      toast.success("Profil berhasil diperbarui");
      onUpdated();
    } catch (err: unknown) {
      toast.error(err instanceof Error ? err.message : "Gagal memperbarui profil");
    } finally {
      setLoading(false);
    }
  };

  return (
    <SectionCard title="Profil Saya" description="Ubah nama yang tampil di platform" icon={User}>
      <form onSubmit={handleSubmit} className="space-y-4">
        {/* Email — read only */}
        <div>
          <label className="block text-xs font-medium text-[#4A443C] mb-1.5 font-[family-name:var(--font-mono)]">
            Email
          </label>
          <div className="flex items-center gap-2.5 px-3.5 py-2.5 rounded-xl border border-[#1F1B16]/8 bg-[#F6F3EC] text-sm text-[#4A443C]">
            <Mail size={14} className="flex-shrink-0 text-[#4A443C]/50" />
            {user.email}
            {user.is_verified && (
              <CheckCircle2 size={13} className="text-emerald-500 ml-auto flex-shrink-0" />
            )}
          </div>
        </div>

        <InputField
          id="profile-name"
          label="Nama Lengkap"
          value={name}
          onChange={setName}
          placeholder="Nama Anda"
        />

        <div className="flex items-center gap-3 pt-1">
          <SubmitButton loading={loading} label="Simpan Perubahan" id="btn-save-profile" />
        </div>
      </form>
    </SectionCard>
  );
}

function PasswordSection({
  user,
  token,
}: {
  user: ClientUser;
  token: string;
}) {
  // Has password = update flow; no password (OAuth user) = create flow
  const hasPassword = Boolean(user.has_password);

  const [oldPass, setOldPass] = useState("");
  const [newPass, setNewPass] = useState("");
  const [showOld, setShowOld] = useState(false);
  const [showNew, setShowNew] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (newPass.length < 8) {
      toast.error("Password baru minimal 8 karakter");
      return;
    }
    setLoading(true);
    try {
      const body = hasPassword
        ? { old_password: oldPass, new_password: newPass }
        : { new_password: newPass };

      const res = await fetch(`${API_URL}/me/password`, {
        method: hasPassword ? "PUT" : "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(body),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.message ?? "Gagal memperbarui password");
      toast.success(hasPassword ? "Password berhasil diubah" : "Password berhasil dibuat");
      setOldPass("");
      setNewPass("");
    } catch (err: unknown) {
      toast.error(err instanceof Error ? err.message : "Gagal memperbarui password");
    } finally {
      setLoading(false);
    }
  };

  return (
    <SectionCard
      title={hasPassword ? "Ubah Password" : "Buat Password"}
      description={
        hasPassword
          ? "Gunakan password yang kuat dan unik"
          : "Akun Google Anda belum memiliki password — buat sekarang untuk bisa login via email"
      }
      icon={hasPassword ? Lock : KeyRound}
    >
      {!hasPassword && (
        <div className="flex items-start gap-2.5 p-3.5 rounded-xl bg-amber-50 border border-amber-200 mb-4 text-sm text-amber-800">
          <AlertCircle size={15} className="flex-shrink-0 mt-0.5" />
          <p>
            Anda masuk via Google. Membuat password memungkinkan login dengan email &amp; password.
          </p>
        </div>
      )}
      <form onSubmit={handleSubmit} className="space-y-4">
        {hasPassword && (
          <InputField
            id="old-password"
            label="Password Lama"
            type={showOld ? "text" : "password"}
            value={oldPass}
            onChange={setOldPass}
            placeholder="••••••••"
            showToggle
            onToggle={() => setShowOld((p) => !p)}
          />
        )}
        <InputField
          id="new-password"
          label="Password Baru"
          type={showNew ? "text" : "password"}
          value={newPass}
          onChange={setNewPass}
          placeholder="Min. 8 karakter"
          showToggle
          onToggle={() => setShowNew((p) => !p)}
        />
        <div className="pt-1">
          <SubmitButton
            loading={loading}
            label={hasPassword ? "Ubah Password" : "Buat Password"}
            id="btn-save-password"
          />
        </div>
      </form>
    </SectionCard>
  );
}

function DangerSection({
  token,
  onLogout,
}: {
  token: string;
  onLogout: () => void;
}) {
  const router = useRouter();
  const [confirm, setConfirm] = useState("");
  const [loading, setLoading] = useState(false);
  const CONFIRM_WORD = "HAPUS";

  const handleDelete = async (e: FormEvent) => {
    e.preventDefault();
    if (confirm !== CONFIRM_WORD) {
      toast.error(`Ketik "${CONFIRM_WORD}" untuk konfirmasi`);
      return;
    }
    setLoading(true);
    try {
      const res = await fetch(`${API_URL}/me`, {
        method: "DELETE",
        headers: { Authorization: `Bearer ${token}` },
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.message ?? "Gagal menghapus akun");
      localStorage.removeItem("client_token");
      localStorage.removeItem("client_user");
      toast.success("Akun berhasil dihapus");
      router.replace("/");
    } catch (err: unknown) {
      toast.error(err instanceof Error ? err.message : "Gagal menghapus akun");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-white rounded-2xl border border-red-200 shadow-sm overflow-hidden">
      <div className="px-6 py-5 border-b border-red-100 flex items-start gap-3">
        <div className="w-9 h-9 rounded-xl bg-red-50 border border-red-200 flex items-center justify-center flex-shrink-0 mt-0.5">
          <ShieldAlert size={17} className="text-red-500" />
        </div>
        <div>
          <h2 className="font-semibold text-red-700 text-sm">Zona Berbahaya</h2>
          <p className="text-xs text-red-500/70 mt-0.5">Tindakan ini tidak bisa dibatalkan</p>
        </div>
      </div>
      <div className="px-6 py-5 space-y-5">
        {/* Logout */}
        <div className="flex items-center justify-between py-3 border-b border-[#1F1B16]/6">
          <div>
            <p className="text-sm font-medium text-[#1F1B16]">Keluar dari akun</p>
            <p className="text-xs text-[#4A443C]/60 mt-0.5">Hapus sesi aktif di perangkat ini</p>
          </div>
          <button
            id="btn-logout"
            onClick={onLogout}
            className="inline-flex items-center gap-2 px-4 py-2 rounded-xl border border-[#1F1B16]/15 text-sm font-medium text-[#4A443C] hover:border-red-300 hover:text-red-600 transition-all"
          >
            <LogOut size={14} /> Logout
          </button>
        </div>

        {/* Delete account */}
        <div>
          <p className="text-sm font-medium text-red-700 mb-1">Hapus Akun Permanen</p>
          <p className="text-xs text-[#4A443C]/60 mb-4">
            Semua data Anda (pesanan, riwayat) akan dihapus selamanya.
          </p>
          <form onSubmit={handleDelete} className="space-y-3">
            <div>
              <label
                htmlFor="delete-confirm"
                className="block text-xs font-medium text-[#4A443C] mb-1.5 font-[family-name:var(--font-mono)]"
              >
                Ketik &ldquo;<span className="text-red-600">{CONFIRM_WORD}</span>&rdquo; untuk konfirmasi
              </label>
              <input
                id="delete-confirm"
                type="text"
                value={confirm}
                onChange={(e) => setConfirm(e.target.value)}
                placeholder={CONFIRM_WORD}
                className="w-full max-w-xs px-3.5 py-2.5 rounded-xl border border-red-200 text-sm text-[#1F1B16] placeholder:text-[#4A443C]/30 focus:outline-none focus:ring-2 focus:ring-red-300 focus:border-red-400 transition-all bg-white"
              />
            </div>
            <button
              id="btn-delete-account"
              type="submit"
              disabled={loading || confirm !== CONFIRM_WORD}
              className="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-red-600 text-white text-sm font-medium hover:bg-red-700 disabled:opacity-40 disabled:cursor-not-allowed transition-all"
            >
              {loading ? (
                <Loader2 size={14} className="animate-spin" />
              ) : (
                <Trash2 size={14} />
              )}
              Hapus Akun Saya
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}

// ── Main Page ─────────────────────────────────────────────────────────────────

export default function ProfilePage() {
  const router = useRouter();
  const { user, token, refresh } = useClientAuth();

  // Redirect to login if not authenticated
  useEffect(() => {
    if (user === null && typeof window !== "undefined") {
      const t = localStorage.getItem("client_token");
      if (!t) router.replace("/login?redirect=/profile");
    }
  }, [user, router]);

  const handleLogout = () => {
    localStorage.removeItem("client_token");
    localStorage.removeItem("client_user");
    toast.success("Berhasil logout");
    router.replace("/");
  };

  if (!user || !token) {
    return (
      <div className={`${FONT_CLASSES} min-h-screen flex items-center justify-center bg-[#F6F3EC]`}>
        <Loader2 size={28} className="animate-spin text-[#B8863C]" />
      </div>
    );
  }

  return (
    <main className={`${FONT_CLASSES} min-h-screen flex flex-col bg-[#F6F3EC]`}>
      {/* Header */}
      <header className="sticky top-0 z-40 bg-[#F6F3EC]/95 backdrop-blur-sm border-b border-[#1F1B16]/8 px-4 py-3">
        <div className="max-w-4xl mx-auto flex items-center justify-between gap-4">
          <Logo small />
          <nav className="flex items-center gap-1 text-sm">
            <Link href="/events" className="px-3 py-1.5 text-[#4A443C]/70 hover:text-[#1F1B16] transition-colors rounded-lg hover:bg-[#1F1B16]/5">
              Event
            </Link>
            <Link href="/orders" className="px-3 py-1.5 text-[#4A443C]/70 hover:text-[#1F1B16] transition-colors rounded-lg hover:bg-[#1F1B16]/5">
              Pesanan
            </Link>
            <Link href="/profile" className="px-3 py-1.5 font-medium text-[#B8863C] rounded-lg bg-[#B8863C]/8">
              Profil
            </Link>
          </nav>
        </div>
      </header>

      <div className="flex-1 py-10 px-4">
        <div className="max-w-2xl mx-auto">
          {/* Page title */}
          <div className="mb-8">
            {/* Breadcrumb */}
            <p className="text-xs font-[family-name:var(--font-mono)] text-[#4A443C]/50 uppercase tracking-[0.2em] mb-3">
              <Link href="/" className="hover:text-[#B8863C] transition-colors">Beranda</Link>
              <span className="mx-2">/</span>
              Profil
            </p>
            {/* User identity card */}
            <div className="flex items-center gap-4 bg-white rounded-2xl border border-[#1F1B16]/8 shadow-sm px-5 py-4">
              <div className="w-14 h-14 rounded-full bg-[#1F1B16] flex items-center justify-center flex-shrink-0">
                <span className="text-xl font-semibold text-[#D8AE6E] font-[family-name:var(--font-display)]">
                  {user.name.charAt(0).toUpperCase()}
                </span>
              </div>
              <div className="flex-1 min-w-0">
                <p className="font-semibold text-[#1F1B16] truncate">{user.name}</p>
                <p className="text-sm text-[#4A443C]/60 truncate">{user.email}</p>
              </div>
              <div className="text-right flex-shrink-0">
                <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium bg-[#B8863C]/10 text-[#B8863C] border border-[#B8863C]/20 font-[family-name:var(--font-mono)]">
                  <Aperture size={11} />
                  {user.subscription_tier === "premium" ? "Premium" : "Basic"}
                </span>
              </div>
            </div>
          </div>

          {/* Sections */}
          <div className="space-y-5">
            <ProfileSection
              user={user}
              token={token}
              onUpdated={() => refresh(token)}
            />
            <PasswordSection user={user} token={token} />
            <DangerSection token={token} onLogout={handleLogout} />
          </div>
        </div>
      </div>

      <Footer />
    </main>
  );
}
