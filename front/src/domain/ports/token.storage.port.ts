import type { AuthTokens } from "../models/auth.model"

export interface ITokenStorage {
  get(): AuthTokens | null
  set(tokens: AuthTokens): void
  clear(): void
}
