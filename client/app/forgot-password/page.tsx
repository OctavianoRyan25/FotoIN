"use client";

import { useState, FormEvent } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Eye, EyeOff, Mail, Lock, AlertCircle, KeyRound, ArrowRight, CheckCircle2 } from "lucide-react";
import { toast } from "sonner";
import { FONT_CLASSES, Logo, Footer } from "../components/site-ui";

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api/v1";

export default function ForgotPasswordPage() {
  const router = useRouter();
  
  // Steps: 1 = Enter Email, 2 = Enter OTP and New Password, 3 = Success
  const [step, setStep] = useState<1 | 2 | 3>(1);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  // Form states
  const [email, setEmail] = useState("");
  const [otp, setOtp] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [showPass, setShowPass] = useState(false);

  const handleSendOTP = async (e: FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      const res = await fetch(`${API_URL}/auth/forgot-password`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email }),
      });
      const data = await res.json();

      if (!res.ok || !data.success) {
        setError(data.message ?? "Gagal mengirim OTP. Pastikan email terdaftar.");
        return;
      }

      toast.success(data.message ?? "Kode OTP telah dikirim ke email Anda");
      setStep(2);
    } catch {
      setError("Tidak dapat terhubung ke server. Coba lagi.");
    } finally {
      setLoading(false);
    }
  };

  const handleResetPassword = async (e: FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      const res = await fetch(`${API_URL}/auth/reset-password`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, otp, new_password: newPassword }),
      });
      const data = await res.json();

      if (!res.ok || !data.success) {
        setError(data.message ?? "Gagal mereset password. Pastikan kode OTP benar.");
        return;
      }

      toast.success("Password berhasil diubah!");
      setStep(3);
    } catch {
      setError("Tidak dapat terhubung ke server. Coba lagi.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <main className={`${FONT_CLASSES} min-h-screen flex flex-col bg-[#F6F3EC]`}>
      {/* Header */}
      <header className="sticky top-0 z-40 bg-[#F6F3EC]/95 backdrop-blur-sm border-b border-[#1F1B16]/8 px-4 py-3">
        <div className="max-w-6xl mx-auto flex items-center justify-between gap-4">
          <Logo small />
          <Link
            href="/login"
            className="text-sm font-medium text-[#4A443C]/70 hover:text-[#1F1B16] transition-colors"
          >
            ← Kembali ke Login
          </Link>
        </div>
      </header>

      {/* Main Content */}
      <div className="flex-1 flex items-center justify-center px-4 py-16">
        <div className="w-full max-w-md">
          {/* Card */}
          <div className="bg-white rounded-3xl border border-[#1F1B16]/8 shadow-xl shadow-[#1F1B16]/5 overflow-hidden">
            
            {/* Success State */}
            {step === 3 ? (
              <div className="px-8 py-12 text-center">
                <div className="w-16 h-16 bg-emerald-50 rounded-full flex items-center justify-center mx-auto mb-5">
                  <CheckCircle2 size={32} className="text-emerald-500" />
                </div>
                <h2 className="text-2xl font-[family-name:var(--font-display)] text-[#1F1B16] font-medium mb-2">
                  Password Berhasil Diubah!
                </h2>
                <p className="text-sm text-[#4A443C]/70 mb-8">
                  Anda sekarang dapat masuk menggunakan password baru.
                </p>
                <Link
                  href="/login"
                  className="inline-flex items-center justify-center w-full py-3.5 rounded-2xl bg-[#1F1B16] text-[#F6F3EC] font-medium text-sm hover:bg-[#B8863C] transition-all duration-300"
                >
                  Masuk Sekarang
                </Link>
              </div>
            ) : (
              <>
                {/* Card Header */}
                <div className="relative bg-[#1F1B16] px-8 pt-10 pb-8 text-center overflow-hidden">
                  <div aria-hidden className="absolute -right-8 -top-8 w-32 h-32 rounded-full border border-[#B8863C]/30" />
                  <div aria-hidden className="absolute -left-6 bottom-0 w-24 h-24 rounded-full border border-[#B8863C]/20" />
                  <div className="relative">
                    <div className="w-14 h-14 bg-[#B8863C]/15 border border-[#B8863C]/40 rounded-2xl flex items-center justify-center mx-auto mb-4">
                      <KeyRound size={24} className="text-[#D8AE6E]" />
                    </div>
                    <h1 className="text-2xl font-[family-name:var(--font-display)] text-[#F6F3EC] font-medium">
                      Reset Password
                    </h1>
                    <p className="text-[#F6F3EC]/50 text-sm mt-1.5 font-[family-name:var(--font-mono)]">
                      {step === 1 ? "Masukkan email Anda untuk OTP" : "Masukkan OTP dan password baru"}
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

                  {/* Step 1 Form */}
                  {step === 1 && (
                    <form onSubmit={handleSendOTP} className="space-y-4">
                      <div>
                        <label className="block text-xs font-medium text-[#4A443C] mb-1.5 font-[family-name:var(--font-mono)]">
                          Email Terdaftar
                        </label>
                        <div className="relative">
                          <Mail size={15} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-[#4A443C]/40" />
                          <input
                            type="email"
                            required
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            placeholder="nama@email.com"
                            className="w-full pl-10 pr-4 py-3 rounded-xl border border-[#1F1B16]/12 text-sm text-[#1F1B16] placeholder:text-[#4A443C]/40 focus:outline-none focus:ring-2 focus:ring-[#B8863C]/30 focus:border-[#B8863C] transition-all bg-white"
                          />
                        </div>
                      </div>
                      <button
                        type="submit"
                        disabled={loading || !email}
                        className="w-full py-3.5 rounded-2xl bg-[#1F1B16] text-[#F6F3EC] font-medium text-sm hover:bg-[#B8863C] disabled:opacity-60 disabled:cursor-not-allowed transition-all duration-300 mt-2 flex items-center justify-center gap-2"
                      >
                        {loading ? (
                          <div className="w-4 h-4 border-2 border-[#F6F3EC]/30 border-t-[#F6F3EC] rounded-full animate-spin" />
                        ) : (
                          <>Kirim Kode OTP <ArrowRight size={16} /></>
                        )}
                      </button>
                    </form>
                  )}

                  {/* Step 2 Form */}
                  {step === 2 && (
                    <form onSubmit={handleResetPassword} className="space-y-4">
                      {/* Readonly Email */}
                      <div>
                        <label className="block text-xs font-medium text-[#4A443C] mb-1.5 font-[family-name:var(--font-mono)]">
                          Email
                        </label>
                        <div className="relative">
                          <Mail size={15} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-[#4A443C]/40" />
                          <input
                            type="email"
                            disabled
                            value={email}
                            className="w-full pl-10 pr-4 py-3 rounded-xl border border-[#1F1B16]/12 text-sm text-[#4A443C] bg-[#F6F3EC]"
                          />
                        </div>
                      </div>

                      {/* OTP Input */}
                      <div>
                        <label className="block text-xs font-medium text-[#4A443C] mb-1.5 font-[family-name:var(--font-mono)]">
                          Kode OTP
                        </label>
                        <input
                          type="text"
                          required
                          value={otp}
                          onChange={(e) => setOtp(e.target.value.toUpperCase())}
                          placeholder="Masukkan OTP dari email"
                          className="w-full px-4 py-3 rounded-xl border border-[#1F1B16]/12 text-sm text-[#1F1B16] font-[family-name:var(--font-mono)] tracking-widest focus:outline-none focus:ring-2 focus:ring-[#B8863C]/30 focus:border-[#B8863C] transition-all bg-white text-center uppercase"
                        />
                      </div>

                      {/* New Password */}
                      <div>
                        <label className="block text-xs font-medium text-[#4A443C] mb-1.5 font-[family-name:var(--font-mono)]">
                          Password Baru
                        </label>
                        <div className="relative">
                          <Lock size={15} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-[#4A443C]/40" />
                          <input
                            type={showPass ? "text" : "password"}
                            required
                            minLength={8}
                            value={newPassword}
                            onChange={(e) => setNewPassword(e.target.value)}
                            placeholder="Min. 8 karakter"
                            className="w-full pl-10 pr-11 py-3 rounded-xl border border-[#1F1B16]/12 text-sm text-[#1F1B16] placeholder:text-[#4A443C]/40 focus:outline-none focus:ring-2 focus:ring-[#B8863C]/30 focus:border-[#B8863C] transition-all bg-white"
                          />
                          <button
                            type="button"
                            onClick={() => setShowPass((p) => !p)}
                            className="absolute right-3.5 top-1/2 -translate-y-1/2 text-[#4A443C]/40 hover:text-[#4A443C] transition-colors"
                          >
                            {showPass ? <EyeOff size={15} /> : <Eye size={15} />}
                          </button>
                        </div>
                      </div>

                      <button
                        type="submit"
                        disabled={loading || !otp || !newPassword}
                        className="w-full py-3.5 rounded-2xl bg-[#1F1B16] text-[#F6F3EC] font-medium text-sm hover:bg-[#B8863C] disabled:opacity-60 disabled:cursor-not-allowed transition-all duration-300 mt-2 flex items-center justify-center gap-2"
                      >
                        {loading ? (
                          <div className="w-4 h-4 border-2 border-[#F6F3EC]/30 border-t-[#F6F3EC] rounded-full animate-spin" />
                        ) : (
                          "Simpan Password Baru"
                        )}
                      </button>

                      <div className="text-center mt-4">
                        <button
                          type="button"
                          onClick={() => setStep(1)}
                          className="text-xs font-medium text-[#4A443C]/60 hover:text-[#B8863C] transition-colors"
                        >
                          Ganti Email
                        </button>
                      </div>
                    </form>
                  )}
                </div>
              </>
            )}
          </div>
        </div>
      </div>

      <Footer />
    </main>
  );
}
