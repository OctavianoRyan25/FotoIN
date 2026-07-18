import type { Metadata } from 'next'
import { Plus_Jakarta_Sans } from 'next/font/google'
import './globals.css'
import { Toaster } from 'sonner'

const jakarta = Plus_Jakarta_Sans({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'FotoIN — Temukan Foto Event Anda',
  description: 'Platform fotografi event Indonesia. Temukan dan beli foto berkualitas tinggi dari event olahraga dan Sunmori favoritmu.',
  openGraph: {
    title: 'FotoIN — Platform Fotografi Event Indonesia',
    description: 'Temukan foto Anda dari event olahraga dan Sunmori. Kualitas tinggi, harga terjangkau.',
    type: 'website',
  },
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="id">
      <body className={jakarta.className}>
        {children}
        <Toaster position="bottom-center" richColors />
      </body>
    </html>
  )
}
