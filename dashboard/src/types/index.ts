export interface User {
  id: string
  name: string
  email: string
  role: 'photographer' | 'super_admin' | 'client'
  subscription_tier: 'basic' | 'premium'
  subscription_expires_at?: string
  is_verified: boolean
  is_suspended: boolean
  storage_used_bytes: number
  wallet_balance: number
  created_at: string
  updated_at: string
}

export interface Gallery {
  id: string
  photographer_id: string
  title: string
  slug: string
  event_date: string
  location: string
  description: string
  price_per_photo: number
  cover_photo_url?: string
  is_active: boolean
  expires_at?: string
  photo_count: number
  total_sales?: number
  created_at: string
  updated_at: string
  photographer?: User
}

export interface Photo {
  id: string
  gallery_id: string
  filename: string
  size_bytes: number
  width: number
  height: number
  preview_url?: string
  created_at: string
  updated_at: string
}

export interface Order {
  id: string
  buyer_email: string
  gallery_id: string
  total_amount: number
  platform_fee: number
  net_amount: number
  payment_status: 'pending' | 'paid' | 'failed' | 'expired'
  snap_url?: string
  download_expires_at?: string
  created_at: string
  gallery?: Gallery
  items?: OrderItem[]
}

export interface OrderItem {
  id: string
  order_id: string
  photo_id: string
  price: number
  download_url?: string
  photo?: Photo
}

export interface WalletTransaction {
  id: string
  photographer_id: string
  order_id?: string
  type: 'credit' | 'debit'
  amount: number
  description: string
  balance_after: number
  created_at: string
}

export interface WithdrawalRequest {
  id: string
  photographer_id: string
  amount: number
  bank_name: string
  account_number: string
  account_holder: string
  status: 'pending' | 'approved' | 'rejected'
  processed_by?: string
  processed_at?: string
  notes?: string
  created_at: string
  photographer?: User
}

export interface WalletSummary {
  balance: number
  transactions: WalletTransaction[]
  total: number
  total_pages: number
}

export interface ApiResponse<T> {
  success: boolean
  message: string
  data?: T
  error?: string
}

export interface PaginatedResponse<T> {
  success: boolean
  message: string
  data: T[]
  meta: {
    page: number
    limit: number
    total: number
    total_pages: number
  }
}

export interface AuthResponse {
  access_token: string
  user: User
}

export interface UploadURLResponse {
  upload_url: string
  s3_key: string
  s3_key_preview: string
}

export interface Subscription {
  id: string
  user_id: string
  tier: 'basic' | 'premium'
  amount: number
  duration_months: number
  status: 'pending' | 'paid' | 'expired'
  qr_code_url?: string
  qr_code_string?: string
  qr_expires_at?: string
  paid_at?: string
  expires_at?: string
  created_at: string
  user?: User
}

export interface SubscriptionStatusResponse {
  user: User
  current_tier: 'basic' | 'premium'
  expires_at?: string
  storage_limit_bytes: number
  storage_used_bytes: number
  commission_rate: number
  latest_subscription?: Subscription
}

export interface RevenueOverview {
  total_revenue: number
  order_commission: number
  basic_tier_revenue: number
  premium_tier_revenue: number
  subscription_revenue: number
  total_orders: number
  total_photographers: number
  active_subscribers: number
}
