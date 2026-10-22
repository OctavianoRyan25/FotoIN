import axios from 'axios'

const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor — inject JWT token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('fotoin_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Response interceptor — handle 401
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('fotoin_token')
      localStorage.removeItem('fotoin_user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default api
