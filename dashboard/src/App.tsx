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
import { SubscriptionPage } from '@/pages/subscription/SubscriptionPage'
import { SettingsPage } from '@/pages/settings/SettingsPage'
import { UserManagementPage } from '@/pages/admin/UserManagementPage'
import { RevenuePage } from '@/pages/admin/RevenuePage'
import { WithdrawalApprovalPage } from '@/pages/admin/WithdrawalApprovalPage'
import { SubscriptionManagementPage } from '@/pages/admin/SubscriptionManagementPage'
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
        {/* Photographer & Common routes */}
        <Route path="events" element={<EventListPage />} />
        <Route path="galleries" element={<GalleryListPage />} />
        <Route path="galleries/:id" element={<GalleryDetailPage />} />
        <Route path="wallet" element={<WalletPage />} />
        <Route path="subscription" element={<SubscriptionPage />} />
        <Route path="settings" element={<SettingsPage />} />

        {/* Super Admin routes */}
        <Route path="users" element={<UserManagementPage />} />
        <Route path="revenue" element={<RevenuePage />} />
        <Route path="withdrawals" element={<WithdrawalApprovalPage />} />
        <Route path="subscriptions" element={<SubscriptionManagementPage />} />
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
