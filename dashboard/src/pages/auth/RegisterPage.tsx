import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { Aperture, Eye, EyeOff, Loader2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useAuth } from '@/store/auth'
import { toast } from 'sonner'

export function RegisterPage() {
  const { register } = useAuth()
  const navigate = useNavigate()
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [showPass, setShowPass] = useState(false)
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (password.length < 8) {
      toast.error('Password minimal 8 karakter')
      return
    }
    setLoading(true)
    try {
      await register(name, email, password)
      toast.success('Akun fotografer berhasil dibuat!')
      navigate('/dashboard')
    } catch (err: unknown) {
      const message =
        (err as { response?: { data?: { message?: string } } })?.response?.data?.message ??
        'Registrasi gagal'
      toast.error(message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex" style={{ fontFamily: 'var(--font-body)' }}>
      {/* ── Left brand panel ── */}
      <div
        className="hidden lg:flex flex-col justify-between w-1/2 p-12 text-white relative overflow-hidden"
        style={{ background: 'linear-gradient(135deg, #23474A 0%, #1F1B16 100%)' }}
      >
        <div
          className="absolute inset-0 opacity-10"
          style={{
            backgroundImage: 'radial-gradient(circle, #F6F3EC 1px, transparent 1px)',
            backgroundSize: '24px 24px',
          }}
        />
        <div className="absolute -right-16 -top-16 w-64 h-64 rounded-full border border-[#B8863C]/20" />

        <div className="relative flex items-center gap-3">
          <div className="w-10 h-10 bg-[#1F1B16] rounded-full flex items-center justify-center border border-[#B8863C]/40">
            <Aperture size={20} className="text-[#D8AE6E]" strokeWidth={1.75} />
          </div>
          <span className="text-2xl font-medium" style={{ fontFamily: 'var(--font-display)' }}>
            Foto<span className="text-[#B8863C] italic">IN</span>
          </span>
        </div>

        <div className="relative">
          <p className="text-[#F6F3EC]/50 text-xs tracking-[0.25em] uppercase mb-4"
             style={{ fontFamily: 'var(--font-mono)' }}>
            Daftar Fotografer
          </p>
          <h1 className="text-4xl font-medium leading-tight mb-5"
              style={{ fontFamily: 'var(--font-display)' }}>
            Mulai jual foto Anda<br />
            <em className="text-[#B8863C]">hari ini</em>
          </h1>
          <p className="text-[#F6F3EC]/60 leading-relaxed mb-8">
            Daftar gratis dan dapatkan 3 GB storage untuk mulai monetisasi karya foto event Anda.
          </p>

          <div className="grid grid-cols-2 gap-3">
            {[
              { label: 'Basic Plan',      value: 'Gratis' },
              { label: 'Storage',         value: '3 GB' },
              { label: 'Komisi Platform', value: '20%' },
              { label: 'Upgrade',         value: 'Premium 4%' },
            ].map((item) => (
              <div key={item.label} className="bg-white/5 border border-white/10 rounded-xl p-3">
                <p className="text-[#F6F3EC]/50 text-xs" style={{ fontFamily: 'var(--font-mono)' }}>
                  {item.label}
                </p>
                <p className="text-white font-semibold mt-0.5">{item.value}</p>
              </div>
            ))}
          </div>
        </div>

        <p className="relative text-[#F6F3EC]/40 text-xs" style={{ fontFamily: 'var(--font-mono)' }}>
          Dengan mendaftar, Anda menyetujui syarat layanan FotoIN.
        </p>
      </div>

      {/* ── Right form panel ── */}
      <div className="flex-1 flex items-center justify-center p-8 bg-background">
        <div className="w-full max-w-md">
          {/* Mobile logo */}
          <div className="lg:hidden flex items-center gap-2 mb-8">
            <div className="w-8 h-8 bg-foreground rounded-full flex items-center justify-center">
              <Aperture size={16} className="text-[#D8AE6E]" strokeWidth={1.75} />
            </div>
            <span className="font-medium text-lg" style={{ fontFamily: 'var(--font-display)' }}>
              Foto<span className="text-brand italic">IN</span>
            </span>
          </div>

          <div className="mb-8">
            <h2 className="text-2xl font-semibold mb-2">Buat Akun Fotografer</h2>
            <p className="text-muted-foreground text-sm">
              Sudah punya akun?{' '}
              <Link to="/login" className="text-brand font-medium hover:underline">
                Masuk
              </Link>
            </p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label htmlFor="register-name" className="text-sm font-medium mb-1.5 block">
                Nama Lengkap
              </label>
              <Input
                id="register-name"
                placeholder="Budi Santoso"
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
              />
            </div>
            <div>
              <label htmlFor="register-email" className="text-sm font-medium mb-1.5 block">
                Email
              </label>
              <Input
                id="register-email"
                type="email"
                placeholder="foto@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>
            <div>
              <label htmlFor="register-password" className="text-sm font-medium mb-1.5 block">
                Password
              </label>
              <div className="relative">
                <Input
                  id="register-password"
                  type={showPass ? 'text' : 'password'}
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
              id="register-submit"
              type="submit"
              className="w-full"
              disabled={loading}
            >
              {loading && <Loader2 size={16} className="animate-spin" />}
              Daftar Gratis
            </Button>
          </form>

          <p className="text-center text-xs text-muted-foreground mt-8"
             style={{ fontFamily: 'var(--font-mono)' }}>
            © {new Date().getFullYear()} FotoIN · Platform Fotografi Event Indonesia
          </p>
        </div>
      </div>
    </div>
  )
}
