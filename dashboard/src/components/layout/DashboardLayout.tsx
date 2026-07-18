import { useState } from 'react'
import { NavLink, Outlet, useNavigate } from 'react-router-dom'
import {
  LayoutDashboard,
  Images,
  Wallet,
  LogOut,
  Menu,
  X,
  Camera,
  Aperture,
  Users,
  TrendingUp,
  Settings,
  ChevronRight,
  Bell,
  Moon,
  Sun,
  Calendar,
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { useAuth } from '@/store/auth'
import { formatBytes } from '@/lib/utils'

const photographerNav = [
  { to: '/dashboard', label: 'Dashboard', icon: LayoutDashboard, exact: true },
  { to: '/dashboard/events', label: 'Event', icon: Calendar },
  { to: '/dashboard/galleries', label: 'Galeri', icon: Images },
  { to: '/dashboard/wallet', label: 'Dompet', icon: Wallet },
  { to: '/dashboard/settings', label: 'Pengaturan', icon: Settings },
]

const adminNav = [
  { to: '/dashboard', label: 'Dashboard', icon: LayoutDashboard, exact: true },
  { to: '/dashboard/users', label: 'Pengguna', icon: Users },
  { to: '/dashboard/revenue', label: 'Revenue', icon: TrendingUp },
  { to: '/dashboard/withdrawals', label: 'Penarikan', icon: Wallet },
]

export function DashboardLayout() {
  const { user, logout } = useAuth()
  const navigate = useNavigate()
  const [sidebarOpen, setSidebarOpen] = useState(true)
  const [darkMode, setDarkMode] = useState(() =>
    document.documentElement.classList.contains('dark')
  )

  const navItems = user?.role === 'super_admin' ? adminNav : photographerNav

  const toggleDark = () => {
    document.documentElement.classList.toggle('dark')
    setDarkMode(!darkMode)
  }

  const storageLimit = user?.subscription_tier === 'premium' ? 107374182400 : 3221225472
  const storagePercent = user ? (user.storage_used_bytes / storageLimit) * 100 : 0

  return (
    <div className="flex h-screen bg-background overflow-hidden">
      {/* Sidebar */}
      <aside
        className={cn(
          'flex flex-col border-r bg-card transition-all duration-300 ease-in-out shadow-sm',
          sidebarOpen ? 'w-64' : 'w-16'
        )}
      >
        {/* Logo */}
        <div className="flex items-center gap-3 px-4 py-5 border-b">
          <div className="flex-shrink-0 w-8 h-8 bg-foreground rounded-full flex items-center justify-center">
            <Aperture className="text-[#D8AE6E]" size={16} strokeWidth={1.75} />
          </div>
          {sidebarOpen && (
            <span
              className="font-medium text-lg tracking-tight"
              style={{ fontFamily: 'var(--font-display)' }}
            >
              Foto<span className="text-brand italic">IN</span>
            </span>
          )}
        </div>

        {/* Nav */}
        <nav className="flex-1 px-2 py-4 space-y-0.5">
          {navItems.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              end={item.exact}
              className={({ isActive }) =>
                cn(
                  'flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all duration-150 group',
                  isActive
                    ? 'bg-brand/10 text-brand dark:bg-brand/15 dark:text-brand-light'
                    : 'text-muted-foreground hover:bg-accent hover:text-foreground'
                )
              }
            >
              <item.icon size={17} className="flex-shrink-0" />
              {sidebarOpen && <span>{item.label}</span>}
              {sidebarOpen && (
                <ChevronRight
                  size={13}
                  className="ml-auto opacity-0 group-hover:opacity-60 transition-opacity"
                />
              )}
            </NavLink>
          ))}
        </nav>

        {/* Film-strip divider */}
        {sidebarOpen && <div className="film-strip mx-4 mb-0" />}

        {/* Storage indicator (photographer only) */}
        {sidebarOpen && user?.role === 'photographer' && (
          <div className="px-4 py-4 border-t">
            <div className="flex justify-between text-xs text-muted-foreground mb-1.5">
              <span>Storage</span>
              <span>{formatBytes(user.storage_used_bytes)}</span>
            </div>
            <div className="h-1.5 bg-muted rounded-full overflow-hidden">
              <div
                className={cn(
                  'h-full rounded-full transition-all',
                  storagePercent > 80 ? 'bg-destructive' : 'bg-brand'
                )}
                style={{ width: `${Math.min(storagePercent, 100)}%` }}
              />
            </div>
            <div className="mt-2.5">
              <Badge
                variant={user.subscription_tier === 'premium' ? 'brand' : 'secondary'}
                className="text-xs"
              >
                {user.subscription_tier === 'premium' ? '✦ Premium' : 'Basic'}
              </Badge>
            </div>
          </div>
        )}

        {/* User info */}
        {sidebarOpen && (
          <div className="px-4 py-4 border-t flex items-center gap-3">
            <div className="w-8 h-8 rounded-full bg-brand flex items-center justify-center text-white text-sm font-medium flex-shrink-0">
              {user?.name.charAt(0).toUpperCase()}
            </div>
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium truncate">{user?.name}</p>
              <p className="text-xs text-muted-foreground truncate">{user?.email}</p>
            </div>
            <Button variant="ghost" size="icon" onClick={logout} title="Logout">
              <LogOut size={15} className="text-muted-foreground" />
            </Button>
          </div>
        )}
      </aside>

      {/* Main */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Topbar */}
        <header className="flex items-center gap-3 px-6 h-14 border-b bg-card">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => setSidebarOpen(!sidebarOpen)}
          >
            {sidebarOpen ? <X size={17} /> : <Menu size={17} />}
          </Button>

          <div className="flex-1" />

          {/* Dark mode toggle */}
          <Button variant="ghost" size="icon" onClick={toggleDark} title="Toggle tema">
            {darkMode ? <Sun size={17} /> : <Moon size={17} />}
          </Button>
          <Button variant="ghost" size="icon" title="Notifikasi">
            <Bell size={17} />
          </Button>

          {/* User avatar chip */}
          {!sidebarOpen && user && (
            <div className="flex items-center gap-2 px-2.5 py-1.5 rounded-full border border-border">
              <div className="w-5 h-5 rounded-full bg-brand flex items-center justify-center text-white text-xs font-medium">
                {user.name.charAt(0).toUpperCase()}
              </div>
              <span className="text-sm font-medium text-foreground max-w-[100px] truncate">
                {user.name.split(' ')[0]}
              </span>
            </div>
          )}
        </header>

        {/* Page content */}
        <main className="flex-1 overflow-y-auto p-6 bg-background">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
