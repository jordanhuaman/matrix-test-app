import { useState } from "react"
import { useNavigate } from "@tanstack/react-router"
import { Button } from "#/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "#/components/ui/card"
import { MatrixGrid } from "../components/matrix/matrix-grid"
import { useAuth } from "../providers/auth-provider"
import { PendingMatrixStorage } from "../../infrastructure/storage/pending-matrix.storage"
import { useSubmitMatrix } from "../hooks/use-matrix"

const pendingStorage = new PendingMatrixStorage()

export function HomePage() {
  const navigate = useNavigate()
  const { isAuthenticated } = useAuth()
  const submit = useSubmitMatrix()
  const [data, setData] = useState<number[][]>([
    [1, 2, 3],
    [4, 5, 6],
  ])

  const handleSubmit = () => {
    if (!isAuthenticated) {
      pendingStorage.save(data)
      navigate({ to: "/register" })
      return
    }
    submit.mutate(data, {
      onSuccess: () => navigate({ to: "/matrix" }),
    })
  }

  return (
    <div className="page-wrap py-12 space-y-8">
      <div className="text-center space-y-4">
        <h1 className="display-title text-4xl sm:text-5xl font-bold">
          Matrix QR Decomposition
        </h1>
        <p className="text-lg text-muted-foreground max-w-xl mx-auto">
          Enter a rectangular matrix and get its QR factorization using Gram-Schmidt.
          View Q (orthogonal) and R (upper triangular) matrices with statistics.
        </p>
      </div>

      <Card className="max-w-2xl mx-auto">
        <CardHeader>
          <CardTitle>Input Matrix</CardTitle>
          <CardDescription>
            Add rows and columns, then click Process to factorize
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <MatrixGrid data={data} onChange={setData} />
          <Button
            className="w-full"
            size="lg"
            onClick={handleSubmit}
            disabled={submit.isPending}
          >
            {submit.isPending ? "Processing..." : "Process Matrix"}
          </Button>
          {submit.isError && (
            <p className="text-sm text-destructive text-center">
              {submit.error?.message ?? "An error occurred"}
            </p>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
