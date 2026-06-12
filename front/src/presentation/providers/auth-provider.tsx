import { createContext, useContext, useCallback, useEffect, useState, type ReactNode } from "react"
import { AuthUseCases } from "../../application/auth/auth.usecases"
import { AuthRepository } from "../../infrastructure/http/auth.repository"
import { TokenStorage } from "../../infrastructure/storage/token.storage"
import type { AuthTokens } from "../../domain/models/auth.model"
import type { LoginRequest } from "../../domain/models/user.model"

const authRepo = new AuthRepository()
const tokenStorage = new TokenStorage()
const authUseCases = new AuthUseCases(authRepo, tokenStorage)

interface AuthContextValue {
  isAuthenticated: boolean
  login: (credentials: LoginRequest) => Promise<AuthTokens>
  register: (email: string, username: string, password: string) => Promise<void>
  logout: () => void
  refreshToken: () => Promise<AuthTokens | null>
}

const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState(() => authUseCases.isAuthenticated())

  const login = useCallback(async (credentials: LoginRequest) => {
    const tokens = await authUseCases.login(credentials)
    setIsAuthenticated(true)
    return tokens
  }, [])

  const register = useCallback(async (email: string, username: string, password: string) => {
    await authUseCases.register(email, username, password)
  }, [])

  const logout = useCallback(() => {
    authUseCases.logout()
    setIsAuthenticated(false)
  }, [])

  const refreshToken = useCallback(async () => {
    const tokens = await authUseCases.refreshToken()
    setIsAuthenticated(tokens !== null)
    return tokens
  }, [])

  useEffect(() => {
    const check = () => setIsAuthenticated(authUseCases.isAuthenticated())
    window.addEventListener("storage", check)
    return () => window.removeEventListener("storage", check)
  }, [])

  return (
    <AuthContext.Provider value={{ isAuthenticated, login, register, logout, refreshToken }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error("useAuth must be used within AuthProvider")
  return ctx
}
