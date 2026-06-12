import { Outlet, Link, useLocation } from "@tanstack/react-router"
import { useAuth } from "../providers/auth-provider"

export function RootLayout() {
  const { isAuthenticated, logout } = useAuth()
  const location = useLocation()

  const isActive = (path: string) => location.pathname === path

  return (
    <div className="min-h-screen flex flex-col">
      <header className="border-b border-border bg-background/80 backdrop-blur-sm sticky top-0 z-50">
        <div className="page-wrap flex items-center justify-between h-16">
          <div className="flex items-center gap-8">
            <Link to="/" className="text-lg font-bold text-foreground no-underline">
              Matrix QR
            </Link>
            <nav className="flex items-center gap-6">
              <Link
                to="/home"
                className={`nav-link text-sm ${isActive("/home") ? "is-active" : ""}`}
              >
                Home
              </Link>
              {isAuthenticated && (
                <Link
                  to="/matrix"
                  className={`nav-link text-sm ${isActive("/matrix") || location.pathname.startsWith("/matrix/") ? "is-active" : ""}`}
                >
                  Dashboard
                </Link>
              )}
            </nav>
          </div>
          <div className="flex items-center gap-4">
            {isAuthenticated ? (
              <button onClick={logout} className="text-sm text-muted-foreground hover:text-foreground cursor-pointer">
                Logout
              </button>
            ) : (
              <div className="flex items-center gap-3">
                <Link to="/login" className="text-sm text-muted-foreground hover:text-foreground">
                  Login
                </Link>
                <Link
                  to="/register"
                  className="inline-flex items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 no-underline"
                >
                  Register
                </Link>
              </div>
            )}
          </div>
        </div>
      </header>
      <main className="flex-1">
        <Outlet />
      </main>
      <footer className="site-footer py-6 text-center text-sm text-muted-foreground">
        Matrix QR Decomposition &mdash; Built with Go + Node.js + TanStack
      </footer>
    </div>
  )
}
