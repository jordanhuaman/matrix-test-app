export interface AuthTokens {
  accessToken: string
  refreshToken: string
}

export interface AuthState {
  user: { id: string; email: string; username: string } | null
  tokens: AuthTokens | null
  isAuthenticated: boolean
}

export interface PendingMatrix {
  data: number[][]
}
