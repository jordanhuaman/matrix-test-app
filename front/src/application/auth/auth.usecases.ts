import type { IAuthRepository } from "../../domain/ports/auth.repository.port"
import type { ITokenStorage } from "../../domain/ports/token.storage.port"
import type { LoginRequest } from "../../domain/models/user.model"
import type { AuthTokens } from "../../domain/models/auth.model"

export class AuthUseCases {
  constructor(
    private readonly authRepo: IAuthRepository,
    private readonly tokenStorage: ITokenStorage,
  ) {}

  async login(credentials: LoginRequest): Promise<AuthTokens> {
    const response = await this.authRepo.login(credentials)
    const tokens: AuthTokens = {
      accessToken: response.accessToken,
      refreshToken: response.refreshToken,
    }
    this.tokenStorage.set(tokens)
    return tokens
  }

  async register(email: string, username: string, password: string): Promise<void> {
    await this.authRepo.register(email, username, password)
  }

  async refreshToken(): Promise<AuthTokens | null> {
    const current = this.tokenStorage.get()
    if (!current?.refreshToken) return null

    try {
      const tokens = await this.authRepo.refreshToken(current.refreshToken)
      this.tokenStorage.set(tokens)
      return tokens
    } catch {
      this.tokenStorage.clear()
      return null
    }
  }

  logout(): void {
    this.tokenStorage.clear()
  }

  isAuthenticated(): boolean {
    return this.tokenStorage.get() !== null
  }
}
