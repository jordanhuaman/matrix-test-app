import { createFileRoute } from "@tanstack/react-router"
import { AuthLayout } from "../presentation/layouts/auth-layout"
import { RegisterPage } from "../presentation/pages/register-page"

export const Route = createFileRoute("/register")({
  component: () => (
    <AuthLayout>
      <RegisterPage />
    </AuthLayout>
  ),
})
