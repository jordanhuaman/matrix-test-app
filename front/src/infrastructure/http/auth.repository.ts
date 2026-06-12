interface GoLoginData {
  access_token: string
  refresh_token: string
}

interface GoRefreshData {
  token: string
  refresh_token: string
}

import type { IAuthRepository } from "../../domain/ports/auth.repository.port"
import type { LoginRequest, LoginResponse } from "../../domain/models/user.model"
import type { AuthTokens } from "../../domain/models/auth.model"
import { httpClient } from "./http-client"

export class AuthRepository implements IAuthRepository {
  async login(credentials: LoginRequest): Promise<LoginResponse> {
    const result = await httpClient.post<GoLoginData>("/auth/login", credentials)
    return {
      accessToken: result.access_token,
      refreshToken: result.refresh_token,
    }
  }

  async register(email: string, username: string, password: string): Promise<void> {
    await httpClient.post("/auth/register", { email, username, password })
  }

  async refreshToken(token: string): Promise<AuthTokens> {
    const result = await httpClient.post<GoRefreshData>("/auth/refresh-token", { refresh_token: token })
    return {
      accessToken: result.token,
      refreshToken: result.refresh_token,
    }
  }
}
