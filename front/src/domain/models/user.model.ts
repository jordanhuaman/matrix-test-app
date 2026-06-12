export interface User {
  id: string
  email: string
  username: string
  isActive: boolean
  createdAt: string
}

export interface RegisterRequest {
  email: string
  username: string
  password: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  accessToken: string
  refreshToken: string
}
