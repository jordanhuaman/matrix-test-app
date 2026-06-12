import { createFileRoute } from "@tanstack/react-router"
import { AuthLayout } from "../presentation/layouts/auth-layout"
import { LoginPage } from "../presentation/pages/login-page"

export const Route = createFileRoute("/login")({
  component: () => (
    <AuthLayout>
      <LoginPage />
    </AuthLayout>
  ),
})
