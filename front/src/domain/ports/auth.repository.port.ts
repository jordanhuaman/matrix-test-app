import type { LoginRequest, LoginResponse } from "../models/user.model"
import type { AuthTokens } from "../models/auth.model"

export interface IAuthRepository {
  login(credentials: LoginRequest): Promise<LoginResponse>
  register(email: string, username: string, password: string): Promise<void>
  refreshToken(token: string): Promise<AuthTokens>
}
