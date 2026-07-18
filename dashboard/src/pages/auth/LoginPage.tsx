import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import {
  Aperture,
  Eye,
  EyeOff,
  Loader2,
  Camera,
  Images,
  Wallet,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useAuth } from "@/store/auth";
import { toast } from "sonner";

export function LoginPage() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPass, setShowPass] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      await login(email, password);
      navigate("/dashboard");
    } catch (err: unknown) {
      const message =
        (err as { response?: { data?: { message?: string } } })?.response?.data
          ?.message ?? "Email atau password salah";
      toast.error(message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      className="min-h-screen flex"
      style={{ fontFamily: "var(--font-body)" }}
    >
      {/* ── Left brand panel ── */}
      <div
        className="hidden lg:flex flex-col justify-between w-1/2 p-12 text-white relative overflow-hidden"
        style={{
          background:
            "linear-gradient(135deg, #1F1B16 0%, #23474A 60%, #1F1B16 100%)",
        }}
      >
        {/* Decorative film-strip dots */}
        <div
          className="absolute inset-0 opacity-10"
          style={{
            backgroundImage:
              "radial-gradient(circle, #F6F3EC 1px, transparent 1px)",
            backgroundSize: "24px 24px",
          }}
        />
        {/* Gold accent circle */}
        <div className="absolute -right-16 -top-16 w-64 h-64 rounded-full border border-[#B8863C]/20" />
        <div className="absolute -left-12 bottom-24 w-48 h-48 rounded-full border border-[#B8863C]/15" />

        <div className="relative flex items-center gap-3">
          <div className="w-10 h-10 bg-[#1F1B16] rounded-full flex items-center justify-center border border-[#B8863C]/40">
            <Aperture size={20} className="text-[#D8AE6E]" strokeWidth={1.75} />
          </div>
          <span
            className="text-2xl font-medium"
            style={{ fontFamily: "var(--font-display)" }}
          >
            Foto<span className="text-[#B8863C] italic">IN</span>
          </span>
        </div>

        <div className="relative">
          <p
            className="text-[#F6F3EC]/50 text-xs tracking-[0.25em] uppercase mb-4"
            style={{ fontFamily: "var(--font-mono)" }}
          >
            Platform Fotografer
          </p>
          <h1
            className="text-4xl font-medium leading-tight mb-5"
            style={{ fontFamily: "var(--font-display)" }}
          >
            Galeri digital
            <br />
            untuk fotografer
            <br />
            <em className="text-[#B8863C]">profesional</em>
          </h1>
          <p className="text-[#F6F3EC]/60 leading-relaxed">
            Unggah foto event, buat galeri terproteksi, dan jual karya Anda
            langsung ke peserta.
          </p>
        </div>

        <div className="relative grid grid-cols-3 gap-4">
          {[
            { icon: Camera, label: "Upload ke S3", sub: "Presigned URL" },
            {
              icon: Images,
              label: "Watermark otomatis",
              sub: "Perlindungan karya",
            },
            { icon: Wallet, label: "Pembayaran aman", sub: "QRIS · Transfer" },
          ].map((f) => (
            <div
              key={f.label}
              className="bg-white/5 border border-white/10 rounded-xl p-4"
            >
              <f.icon
                size={20}
                className="text-[#D8AE6E] mb-2"
                strokeWidth={1.5}
              />
              <p className="text-sm font-medium text-white">{f.label}</p>
              <p
                className="text-xs text-white/50 mt-0.5"
                style={{ fontFamily: "var(--font-mono)" }}
              >
                {f.sub}
              </p>
            </div>
          ))}
        </div>
      </div>

      {/* ── Right form panel ── */}
      <div className="flex-1 flex items-center justify-center p-8 bg-background">
        <div className="w-full max-w-md">
          {/* Mobile logo */}
          <div className="lg:hidden flex items-center gap-2 mb-8">
            <div className="w-8 h-8 bg-foreground rounded-full flex items-center justify-center">
              <Aperture
                size={16}
                className="text-[#D8AE6E]"
                strokeWidth={1.75}
              />
            </div>
            <span
              className="font-medium text-lg"
              style={{ fontFamily: "var(--font-display)" }}
            >
              Foto<span className="text-brand italic">IN</span>
            </span>
          </div>

          <div className="mb-8">
            <h2 className="text-2xl font-semibold mb-2">Masuk ke Dashboard</h2>
            <p className="text-muted-foreground text-sm">
              Belum punya akun?{" "}
              <Link
                to="/register"
                className="text-brand font-medium hover:underline"
              >
                Daftar gratis
              </Link>
            </p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label
                htmlFor="login-email"
                className="text-sm font-medium mb-1.5 block"
              >
                Email
              </label>
              <Input
                id="login-email"
                type="email"
                placeholder="foto@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>
            <div>
              <label
                htmlFor="login-password"
                className="text-sm font-medium mb-1.5 block"
              >
                Password
              </label>
              <div className="relative">
                <Input
                  id="login-password"
                  type={showPass ? "text" : "password"}
                  placeholder="Min. 8 karakter"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                  className="pr-10"
                />
                <button
                  type="button"
                  onClick={() => setShowPass(!showPass)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                >
                  {showPass ? <EyeOff size={16} /> : <Eye size={16} />}
                </button>
              </div>
            </div>

            <Button
              id="login-submit"
              type="submit"
              className="w-full"
              disabled={loading}
            >
              {loading && <Loader2 size={16} className="animate-spin" />}
              Masuk
            </Button>
          </form>

          <p
            className="text-center text-xs text-muted-foreground mt-8"
            style={{ fontFamily: "var(--font-mono)" }}
          >
            © {new Date().getFullYear()} FotoIN · Platform Fotografi Event
            Indonesia
          </p>
        </div>
      </div>
    </div>
  );
}
