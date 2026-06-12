import { createFileRoute, useParams, useNavigate } from "@tanstack/react-router"
import { useEffect } from "react"
import { useAuth } from "../presentation/providers/auth-provider"
import { MatrixResultPage } from "../presentation/pages/matrix-result-page"

export const Route = createFileRoute("/matrix/$id")({
  component: MatrixResultWrapper,
})

function MatrixResultWrapper() {
  const { id } = useParams({ from: "/matrix/$id" })
  const { isAuthenticated } = useAuth()
  const navigate = useNavigate()

  useEffect(() => {
    if (!isAuthenticated) {
      navigate({ to: "/login", replace: true })
    }
  }, [isAuthenticated, navigate])

  if (!isAuthenticated) return null

  return <MatrixResultPage id={id} />
}
