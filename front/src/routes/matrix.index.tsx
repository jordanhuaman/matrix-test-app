import { createFileRoute } from "@tanstack/react-router"
import { useState, useEffect } from "react"
import { useNavigate, useLocation } from "@tanstack/react-router"
import { Button } from "#/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "#/components/ui/card"
import { MatrixGrid } from "../presentation/components/matrix/matrix-grid"
import { MatrixResultView } from "../presentation/components/matrix/matrix-result-view"
import { useSubmitMatrix, useMatrixResults } from "../presentation/hooks/use-matrix"
import { PendingMatrixStorage } from "../infrastructure/storage/pending-matrix.storage"

const pendingStorage = new PendingMatrixStorage()

export const Route = createFileRoute("/matrix/")({
  component: MatrixDashboard,
})

function MatrixDashboard() {
  const navigate = useNavigate()
  const location = useLocation()
  const submit = useSubmitMatrix()
  const { data: resultsData } = useMatrixResults(1, 20)
  const [data, setData] = useState<number[][]>([
    [1, 2, 3],
    [4, 5, 6],
  ])

  useEffect(() => {
    const pending = pendingStorage.get()
    const statePending = (location.state as { pendingMatrix?: number[][] } | null)?.pendingMatrix
    const matrix = pending ?? statePending ?? null
    if (matrix) {
      setData(matrix)
      pendingStorage.clear()
      submit.mutate(matrix, {
        onSuccess: (result) => {
          navigate({ to: "/matrix/$id", params: { id: result.id } })
        },
      })
    }
  }, [])

  const lastResult = submit.data ?? resultsData?.results?.[0]

  const handleSubmit = () => {
    submit.mutate(data, {
      onSuccess: (result) => {
        navigate({ to: "/matrix/$id", params: { id: result.id } })
      },
    })
  }

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>New Matrix</CardTitle>
          <CardDescription>Enter a matrix to factorize</CardDescription>
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

      {lastResult && (
        <div>
          <h2 className="text-lg font-semibold mb-4">Latest Result</h2>
          <MatrixResultView result={lastResult} />
        </div>
      )}
    </div>
  )
}
