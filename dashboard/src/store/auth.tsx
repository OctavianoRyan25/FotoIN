import React, { createContext, useContext, useState, useEffect } from 'react'
import type { User, AuthResponse } from '@/types'
import api from '@/lib/api'

interface AuthContextType {
  user: User | null
  token: string | null
  login: (email: string, password: string) => Promise<void>
  register: (name: string, email: string, password: string) => Promise<void>
  logout: () => void
  isLoading: boolean
}

const AuthContext = createContext<AuthContextType | null>(null)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const storedToken = localStorage.getItem('fotoin_token')
    const storedUser = localStorage.getItem('fotoin_user')
    if (storedToken && storedUser) {
      setToken(storedToken)
      setUser(JSON.parse(storedUser))
    }
    setIsLoading(false)
  }, [])

  const login = async (email: string, password: string) => {
    const res = await api.post<{ data: AuthResponse }>('/auth/login', { email, password })
    const { access_token, user } = res.data.data
    localStorage.setItem('fotoin_token', access_token)
    localStorage.setItem('fotoin_user', JSON.stringify(user))
    setToken(access_token)
    setUser(user)
  }

  const register = async (name: string, email: string, password: string) => {
    const res = await api.post<{ data: AuthResponse }>('/auth/register', { name, email, password })
    const { access_token, user } = res.data.data
    localStorage.setItem('fotoin_token', access_token)
    localStorage.setItem('fotoin_user', JSON.stringify(user))
    setToken(access_token)
    setUser(user)
  }

  const logout = () => {
    localStorage.removeItem('fotoin_token')
    localStorage.removeItem('fotoin_user')
    setToken(null)
    setUser(null)
    window.location.href = '/login'
  }

  return (
    <AuthContext.Provider value={{ user, token, login, register, logout, isLoading }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
