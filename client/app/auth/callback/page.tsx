'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { CheckCircle, AlertCircle } from 'lucide-react'

/**
 * Google OAuth Callback Page
 *
 * The backend redirects here after Google authentication:
 *   http://localhost:3000/auth/callback#token=<JWT>
 *
 * The token is in the URL fragment (# hash) so it's NOT sent to any server.
 * We extract it, store in localStorage, then redirect to the client portal.
 */
export default function AuthCallbackPage() {
  const router = useRouter()
  const [state, setState] = useState<'loading' | 'success' | 'error'>('loading')
  const [errorMsg, setErrorMsg] = useState('')

  useEffect(() => {
    // Read fragment from URL (not accessible server-side)
    const hash = window.location.hash.slice(1) // remove leading '#'
    const params = new URLSearchParams(hash)

    const token = params.get('token')
    const error = params.get('error') ?? new URLSearchParams(window.location.search).get('error')

    if (error) {
      setState('error')
      setErrorMsg(decodeURIComponent(error))
      return
    }

    if (!token) {
      setState('error')
      setErrorMsg('Token tidak ditemukan. Coba login lagi.')
      return
    }

    // Persist JWT for the client portal
    localStorage.setItem('client_token', token)

    // Fetch user profile
    const API_URL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080/api/v1'
    fetch(`${API_URL}/me`, {
      headers: { Authorization: `Bearer ${token}` }
    })
      .then(res => res.json())
      .then(data => {
        if (data.success && data.data) {
          localStorage.setItem('client_user', JSON.stringify(data.data))
          setState('success')
          setTimeout(() => router.replace('/orders'), 1500)
        } else {
          throw new Error('Gagal memuat profil')
        }
      })
      .catch((err) => {
        setState('error')
        setErrorMsg('Gagal mengambil data user: ' + err.message)
      })
  }, [router])

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 px-4">
      <div className="bg-white rounded-2xl shadow-lg p-8 max-w-sm w-full text-center">
        {state === 'loading' && (
          <>
            <div className="w-12 h-12 border-4 border-violet-200 border-t-violet-600 rounded-full animate-spin mx-auto mb-4" />
            <p className="text-gray-600">Memproses login Google…</p>
          </>
        )}

        {state === 'success' && (
          <>
            <CheckCircle size={48} className="text-green-500 mx-auto mb-4" />
            <h1 className="text-xl font-bold text-gray-900 mb-2">Login Berhasil!</h1>
            <p className="text-gray-600 text-sm">Mengalihkan ke halaman pesanan…</p>
          </>
        )}

        {state === 'error' && (
          <>
            <AlertCircle size={48} className="text-red-500 mx-auto mb-4" />
            <h1 className="text-xl font-bold text-gray-900 mb-2">Login Gagal</h1>
            <p className="text-gray-600 text-sm mb-4">{errorMsg}</p>
            <button
              onClick={() => router.push('/')}
              className="px-6 py-2 bg-violet-600 text-white rounded-xl font-medium hover:bg-violet-700"
            >
              Kembali ke Beranda
            </button>
          </>
        )}
      </div>
    </div>
  )
}
