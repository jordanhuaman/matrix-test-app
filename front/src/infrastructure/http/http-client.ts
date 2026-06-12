const API_BASE = import.meta.env.VITE_API_URL ?? "http://localhost:3000/api"

interface RequestConfig {
  method: string
  body?: unknown
  headers?: Record<string, string>
}

function getToken(): string | null {
  try {
    const raw = localStorage.getItem("auth_tokens")
    if (!raw) return null
    const parsed = JSON.parse(raw)
    return parsed?.accessToken ?? null
  } catch {
    return null
  }
}

export class HttpError extends Error {
  constructor(
    public status: number,
    message: string,
  ) {
    super(message)
    this.name = "HttpError"
  }
}

async function request<T>(path: string, config: RequestConfig = { method: "GET" }): Promise<T> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...config.headers,
  }

  const token = getToken()
  if (token) {
    headers["Authorization"] = `Bearer ${token}`
  }

  const response = await fetch(`${API_BASE}${path}`, {
    method: config.method,
    headers,
    body: config.body ? JSON.stringify(config.body) : undefined,
  })

  if (response.status === 401) {
    localStorage.removeItem("auth_tokens")
    window.location.href = "/login"
    throw new HttpError(401, "Unauthorized")
  }

  const json = await response.json()

  if (!response.ok) {
    throw new HttpError(response.status, json.message ?? "Request failed")
  }

  return json.data as T
}

export const httpClient = {
  get: <T>(path: string) => request<T>(path, { method: "GET" }),
  post: <T>(path: string, body?: unknown) => request<T>(path, { method: "POST", body }),
  patch: <T>(path: string, body?: unknown) => request<T>(path, { method: "PATCH", body }),
  delete: <T>(path: string) => request<T>(path, { method: "DELETE" }),
}
