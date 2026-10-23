'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { CheckCircle, AlertCircle } from 'lucide-react'

type State = 'loading' | 'success' | 'error'

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080/api/v1'

// Pesan error dipetakan ke teks yang aman, bukan menampilkan string mentah dari URL
const ERROR_MESSAGES: Record<string, string> = {
  access_denied: 'Kamu membatalkan login Google.',
  invalid_token: 'Sesi login tidak valid. Silakan coba lagi.',
}
const DEFAULT_ERROR = 'Terjadi kesalahan saat login. Silakan coba lagi.'

export default function AuthCallbackPage() {
  const router = useRouter()
  const [state, setState] = useState<State>('loading')
  const [errorMsg, setErrorMsg] = useState('')

  useEffect(() => {
    const controller = new AbortController()
    let redirectTimer: ReturnType<typeof setTimeout> | undefined

    const fail = (msg: string) => {
      localStorage.removeItem('client_token')
      localStorage.removeItem('client_user')
      setErrorMsg(msg)
      setState('error')
    }

    const hashParams = new URLSearchParams(window.location.hash.slice(1))
    const searchParams = new URLSearchParams(window.location.search)

    const token = hashParams.get('token')
    const error = hashParams.get('error') ?? searchParams.get('error')

    // Hapus token dari address bar & history secepatnya
    window.history.replaceState(null, '', window.location.pathname)

    if (error) {
      fail(ERROR_MESSAGES[error] ?? DEFAULT_ERROR)
      return
    }

    if (!token) {
      fail('Token tidak ditemukan. Silakan login lagi.')
      return
    }

    const timeout = setTimeout(() => controller.abort(), 10_000)

    ;(async () => {
      try {
        const res = await fetch(`${API_URL}/me`, {
          headers: { Authorization: `Bearer ${token}` },
          signal: controller.signal,
        })
        if (!res.ok) throw new Error(`HTTP ${res.status}`)

        const data = await res.json()
        if (!data?.success || !data?.data) throw new Error('Respons tidak valid')

        localStorage.setItem('client_token', token)
        localStorage.setItem('client_user', JSON.stringify(data.data))

        setState('success')
        redirectTimer = setTimeout(() => router.replace('/orders'), 1200)
      } catch (err) {
        if (controller.signal.aborted && !timeout) return
        const aborted = err instanceof DOMException && err.name === 'AbortError'
        fail(
          aborted
            ? 'Server terlalu lama merespons. Silakan coba lagi.'
            : 'Gagal mengambil data akun. Silakan coba lagi.'
        )
      } finally {
        clearTimeout(timeout)
      }
    })()

    return () => {
      controller.abort()
      clearTimeout(timeout)
      if (redirectTimer) clearTimeout(redirectTimer)
    }
  }, [router])

  return (
    <main className="min-h-screen flex items-center justify-center bg-gray-50 px-4">
      <div
        className="bg-white rounded-2xl shadow-lg p-8 max-w-sm w-full text-center"
        role="status"
        aria-live="polite"
      >
        {state === 'loading' && (
          <>
            <div
              aria-hidden="true"
              className="w-12 h-12 border-4 border-violet-200 border-t-violet-600 rounded-full animate-spin mx-auto mb-4"
            />
            <p className="text-gray-600">Memproses login Google…</p>
          </>
        )}

        {state === 'success' && (
          <>
            <CheckCircle size={48} aria-hidden="true" className="text-green-500 mx-auto mb-4" />
            <h1 className="text-xl font-bold text-gray-900 mb-2">Login Berhasil!</h1>
            <p className="text-gray-600 text-sm mb-4">Mengalihkan ke halaman pesanan…</p>
            <button
              onClick={() => router.replace('/orders')}
              className="text-sm text-violet-600 hover:underline"
            >
              Lanjut sekarang
            </button>
          </>
        )}

        {state === 'error' && (
          <>
            <AlertCircle size={48} aria-hidden="true" className="text-red-500 mx-auto mb-4" />
            <h1 className="text-xl font-bold text-gray-900 mb-2">Login Gagal</h1>
            <p className="text-gray-600 text-sm mb-4">{errorMsg}</p>
            <button
              onClick={() => router.push('/')}
              className="px-6 py-2 bg-violet-600 text-white rounded-xl font-medium hover:bg-violet-700 focus:outline-none focus:ring-2 focus:ring-violet-400"
            >
              Kembali ke Beranda
            </button>
          </>
        )}
      </div>
    </main>
  )
}