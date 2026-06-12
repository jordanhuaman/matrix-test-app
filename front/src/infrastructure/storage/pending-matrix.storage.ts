const STORAGE_KEY = "pending_matrix"

export class PendingMatrixStorage {
  save(data: number[][]): void {
    sessionStorage.setItem(STORAGE_KEY, JSON.stringify(data))
  }

  get(): number[][] | null {
    try {
      const raw = sessionStorage.getItem(STORAGE_KEY)
      return raw ? (JSON.parse(raw) as number[][]) : null
    } catch {
      return null
    }
  }

  clear(): void {
    sessionStorage.removeItem(STORAGE_KEY)
  }
}
