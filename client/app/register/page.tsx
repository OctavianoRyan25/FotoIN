"use client";

import { useState, FormEvent } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Eye, EyeOff, UserPlus, Mail, Lock, User, AlertCircle } from "lucide-react";
import { toast } from "sonner";
import { FONT_CLASSES, Logo, Footer } from "../components/site-ui";

const API_URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api/v1";

export default function RegisterPage() {
  const router = useRouter();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPass, setShowPass] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleRegister = async (e: FormEvent) => {
    e.preventDefault();
    setError("");
    if (password.length < 8) {
      setError("Password minimal 8 karakter.");
      return;
    }
    setLoading(true);
    try {
      const res = await fetch(`${API_URL}/auth/register`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name, email, password, role: "client" }),
      });
      const data = await res.json();
      if (!res.ok || !data.success) {
        setError(data.message ?? "Pendaftaran gagal. Coba lagi.");
        return;
      }
      localStorage.setItem("client_token", data.data.access_token);
      localStorage.setItem("client_user", JSON.stringify(data.data.user));
      toast.success("Akun berhasil dibuat! Selamat datang di FotoIN.");
      router.replace("/orders");
    } catch {
      setError("Tidak dapat terhubung ke server. Coba lagi.");
    } finally {
      setLoading(false);
    }
  };

  const handleGoogleLogin = () => {
    window.location.href = `${API_URL}/auth/google`;
  };

  return (
    <main className={`${FONT_CLASSES} min-h-screen flex flex-col bg-[#F6F3EC]`}>
      <header className="sticky top-0 z-40 bg-[#F6F3EC]/95 backdrop-blur-sm border-b border-[#1F1B16]/8 px-4 py-3">
        <div className="max-w-6xl mx-auto flex items-center justify-between gap-4">
          <Logo small />
          <Link href="/" className="text-sm font-medium text-[#4A443C]/70 hover:text-[#1F1B16] transition-colors">
            ← Kembali ke Beranda
          </Link>
        </div>
      </header>

      <div className="flex-1 flex items-center justify-center px-4 py-16">
        <div className="w-full max-w-md">
          <div className="bg-white rounded-3xl border border-[#1F1B16]/8 shadow-xl shadow-[#1F1B16]/5 overflow-hidden">
            <div className="relative bg-[#23474A] px-8 pt-10 pb-8 text-center overflow-hidden">
              <div aria-hidden className="absolute -right-8 -top-8 w-32 h-32 rounded-full border border-[#B8863C]/30" />
              <div aria-hidden className="absolute -left-6 bottom-0 w-24 h-24 rounded-full border border-[#B8863C]/20" />
              <div className="relative">
                <div className="w-14 h-14 bg-[#B8863C]/15 border border-[#B8863C]/40 rounded-2xl flex items-center justify-center mx-auto mb-4">
                  <UserPlus size={24} className="text-[#D8AE6E]" />
                </div>
                <h1 className="text-2xl font-[family-name:var(--font-display)] text-[#F6F3EC] font-medium">
                  Buat Akun FotoIN
                </h1>
                <p className="text-[#F6F3EC]/50 text-sm mt-1.5 font-[family-name:var(--font-mono)]">
                  Daftar gratis dan temukan foto Anda
                </p>
              </div>
            </div>

            <div className="px-8 py-8">
              {error && (
                <div className="flex items-start gap-3 mb-5 p-3.5 rounded-xl bg-red-50 border border-red-200">
                  <AlertCircle size={16} className="text-red-500 shrink-0 mt-0.5" />
                  <p className="text-red-700 text-sm">{error}</p>
                </div>
              )}

              <button
                id="btn-google-register"
                type="button"
                onClick={handleGoogleLogin}
                className="w-full flex items-center justify-center gap-3 py-3.5 px-4 rounded-2xl border border-[#1F1B16]/12 bg-white hover:bg-[#F6F3EC] hover:border-[#B8863C]/40 transition-all duration-200 group mb-5"
              >
                <svg width="18" height="18" viewBox="0 0 18 18" fill="none" xmlns="http://www.w3.org/2000/svg" className="shrink-0">
                  <path d="M17.64 9.205c0-.639-.057-1.252-.164-1.841H9v3.481h4.844a4.14 4.14 0 0 1-1.796 2.716v2.259h2.908c1.702-1.567 2.684-3.875 2.684-6.615Z" fill="#4285F4" />
                  <path d="M9 18c2.43 0 4.467-.806 5.956-2.18l-2.908-2.259c-.806.54-1.837.86-3.048.86-2.344 0-4.328-1.584-5.036-3.711H.957v2.332A8.997 8.997 0 0 0 9 18Z" fill="#34A853" />
                  <path d="M3.964 10.71A5.41 5.41 0 0 1 3.682 9c0-.593.102-1.17.282-1.71V4.958H.957A8.996 8.996 0 0 0 0 9c0 1.452.348 2.827.957 4.042l3.007-2.332Z" fill="#FBBC05" />
                  <path d="M9 3.58c1.321 0 2.508.454 3.44 1.345l2.582-2.58C13.463.891 11.426 0 9 0A8.997 8.997 0 0 0 .957 4.958L3.964 7.29C4.672 5.163 6.656 3.58 9 3.58Z" fill="#EA4335" />
                </svg>
                <span className="text-sm font-medium text-[#1F1B16] group-hover:text-[#B8863C] transition-colors">
                  Daftar dengan Google
                </span>
              </button>

              <div className="flex items-center gap-3 mb-5">
                <div className="flex-1 h-px bg-[#1F1B16]/10" />
                <span className="text-xs text-[#4A443C]/50 font-[family-name:var(--font-mono)]">atau daftar dengan email</span>
                <div className="flex-1 h-px bg-[#1F1B16]/10" />
              </div>

              <form onSubmit={handleRegister} className="space-y-4">
                <div>
                  <label htmlFor="reg-name" className="block text-xs font-medium text-[#4A443C] mb-1.5 font-[family-name:var(--font-mono)]">Nama Lengkap</label>
                  <div className="relative">
                    <User size={15} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-[#4A443C]/40" />
                    <input id="reg-name" type="text" autoComplete="name" required minLength={2} value={name} onChange={e => setName(e.target.value)} placeholder="Nama kamu" className="w-full pl-10 pr-4 py-3 rounded-xl border border-[#1F1B16]/12 text-sm text-[#1F1B16] placeholder:text-[#4A443C]/40 focus:outline-none focus:ring-2 focus:ring-[#B8863C]/30 focus:border-[#B8863C] transition-all bg-white" />
                  </div>
                </div>
                <div>
                  <label htmlFor="reg-email" className="block text-xs font-medium text-[#4A443C] mb-1.5 font-[family-name:var(--font-mono)]">Email</label>
                  <div className="relative">
                    <Mail size={15} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-[#4A443C]/40" />
                    <input id="reg-email" type="email" autoComplete="email" required value={email} onChange={e => setEmail(e.target.value)} placeholder="nama@email.com" className="w-full pl-10 pr-4 py-3 rounded-xl border border-[#1F1B16]/12 text-sm text-[#1F1B16] placeholder:text-[#4A443C]/40 focus:outline-none focus:ring-2 focus:ring-[#B8863C]/30 focus:border-[#B8863C] transition-all bg-white" />
                  </div>
                </div>
                <div>
                  <label htmlFor="reg-password" className="block text-xs font-medium text-[#4A443C] mb-1.5 font-[family-name:var(--font-mono)]">Password</label>
                  <div className="relative">
                    <Lock size={15} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-[#4A443C]/40" />
                    <input id="reg-password" type={showPass ? "text" : "password"} autoComplete="new-password" required minLength={8} value={password} onChange={e => setPassword(e.target.value)} placeholder="Min. 8 karakter" className="w-full pl-10 pr-11 py-3 rounded-xl border border-[#1F1B16]/12 text-sm text-[#1F1B16] placeholder:text-[#4A443C]/40 focus:outline-none focus:ring-2 focus:ring-[#B8863C]/30 focus:border-[#B8863C] transition-all bg-white" />
                    <button type="button" onClick={() => setShowPass(p => !p)} className="absolute right-3.5 top-1/2 -translate-y-1/2 text-[#4A443C]/40 hover:text-[#4A443C] transition-colors">
                      {showPass ? <EyeOff size={15} /> : <Eye size={15} />}
                    </button>
                  </div>
                </div>
                <button
                  id="btn-register-submit"
                  type="submit"
                  disabled={loading}
                  className="w-full py-3.5 rounded-2xl bg-[#23474A] text-[#F6F3EC] font-medium text-sm hover:bg-[#B8863C] disabled:opacity-60 disabled:cursor-not-allowed transition-all duration-300 mt-2 flex items-center justify-center gap-2"
                >
                  {loading ? (
                    <><div className="w-4 h-4 border-2 border-[#F6F3EC]/30 border-t-[#F6F3EC] rounded-full animate-spin" />Mendaftarkan…</>
                  ) : "Buat Akun"}
                </button>
              </form>

              <p className="text-center text-sm text-[#4A443C]/70 mt-6">
                Sudah punya akun?{" "}
                <Link href="/login" className="font-medium text-[#B8863C] hover:underline">
                  Masuk sekarang
                </Link>
              </p>
            </div>
          </div>

          <p className="text-center text-xs text-[#4A443C]/50 mt-6 font-[family-name:var(--font-mono)] leading-relaxed">
            Dengan mendaftar, Anda menyetujui <span className="text-[#B8863C]">Syarat & Ketentuan</span> serta <span className="text-[#B8863C]">Kebijakan Privasi</span> FotoIN.
          </p>
        </div>
      </div>

      <Footer />
    </main>
  );
}
