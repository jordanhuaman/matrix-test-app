import type { ITokenStorage } from "../../domain/ports/token.storage.port"
import type { AuthTokens } from "../../domain/models/auth.model"

const STORAGE_KEY = "auth_tokens"

export class TokenStorage implements ITokenStorage {
  get(): AuthTokens | null {
    try {
      const raw = localStorage.getItem(STORAGE_KEY)
      return raw ? (JSON.parse(raw) as AuthTokens) : null
    } catch {
      return null
    }
  }

  set(tokens: AuthTokens): void {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(tokens))
  }

  clear(): void {
    localStorage.removeItem(STORAGE_KEY)
  }
}
