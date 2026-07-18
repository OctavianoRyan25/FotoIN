import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { Toaster } from 'sonner'
import { AuthProvider, useAuth } from '@/store/auth'
import { DashboardLayout } from '@/components/layout/DashboardLayout'
import { LoginPage } from '@/pages/auth/LoginPage'
import { RegisterPage } from '@/pages/auth/RegisterPage'
import { DashboardHomePage } from '@/pages/dashboard/DashboardHomePage'
import { GalleryListPage } from '@/pages/gallery/GalleryListPage'
import { GalleryDetailPage } from '@/pages/gallery/GalleryDetailPage'
import { WalletPage } from '@/pages/wallet/WalletPage'
import EventListPage from '@/pages/event/EventListPage'
import { Loader2 } from 'lucide-react'

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { user, isLoading } = useAuth()
  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Loader2 className="animate-spin text-violet-600" size={32} />
      </div>
    )
  }
  if (!user) return <Navigate to="/login" replace />
  return <>{children}</>
}

function PublicRoute({ children }: { children: React.ReactNode }) {
  const { user, isLoading } = useAuth()
  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Loader2 className="animate-spin text-violet-600" size={32} />
      </div>
    )
  }
  if (user) return <Navigate to="/dashboard" replace />
  return <>{children}</>
}

function AppRoutes() {
  return (
    <Routes>
      {/* Public */}
      <Route path="/login" element={<PublicRoute><LoginPage /></PublicRoute>} />
      <Route path="/register" element={<PublicRoute><RegisterPage /></PublicRoute>} />

      {/* Protected Dashboard */}
      <Route path="/dashboard" element={<ProtectedRoute><DashboardLayout /></ProtectedRoute>}>
        <Route index element={<DashboardHomePage />} />
        <Route path="events" element={<EventListPage />} />
        <Route path="galleries" element={<GalleryListPage />} />
        <Route path="galleries/:id" element={<GalleryDetailPage />} />
        <Route path="wallet" element={<WalletPage />} />
      </Route>

      {/* Default redirect */}
      <Route path="/" element={<Navigate to="/dashboard" replace />} />
      <Route path="*" element={<Navigate to="/dashboard" replace />} />
    </Routes>
  )
}

function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <AppRoutes />
        <Toaster position="top-right" richColors />
      </BrowserRouter>
    </AuthProvider>
  )
}

export default App
