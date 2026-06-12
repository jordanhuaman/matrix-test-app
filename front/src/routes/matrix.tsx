import { Outlet, createFileRoute, useNavigate } from "@tanstack/react-router"
import { useEffect } from "react"
import { useAuth } from "../presentation/providers/auth-provider"
import { MatrixLayout } from "../presentation/layouts/matrix-layout"
import { HistorySidebar } from "../presentation/components/matrix/history-sidebar"

export const Route = createFileRoute("/matrix")({
  component: MatrixLayoutShell,
})

function MatrixLayoutShell() {
  const { isAuthenticated } = useAuth()
  const navigate = useNavigate()

  useEffect(() => {
    if (!isAuthenticated) {
      navigate({ to: "/login", replace: true })
    }
  }, [isAuthenticated, navigate])

  if (!isAuthenticated) return null

  return (
    <MatrixLayout sidebar={<HistorySidebar />}>
      <Outlet />
    </MatrixLayout>
  )
}
