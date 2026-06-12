import { Link } from "@tanstack/react-router"
import { Button } from "#/components/ui/button"
import { Skeleton } from "#/components/ui/skeleton"
import { ArrowLeft } from "lucide-react"
import { MatrixResultView } from "../components/matrix/matrix-result-view"
import { useMatrixResult } from "../hooks/use-matrix"

export function MatrixResultPage({ id }: { id: string }) {
  const { data: result, isLoading, isError, error } = useMatrixResult(id)

  return (
    <div className="max-w-4xl mx-auto space-y-6 p-6">
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" asChild>
          <Link to="/matrix">
            <ArrowLeft className="h-4 w-4" />
          </Link>
        </Button>
        <div>
          <h1 className="text-xl font-semibold">Matrix Result</h1>
          <p className="text-sm text-muted-foreground font-mono">{id}</p>
        </div>
      </div>

      {isLoading && (
        <div className="space-y-4">
          <Skeleton className="h-64 w-full rounded-lg" />
          <Skeleton className="h-32 w-full rounded-lg" />
        </div>
      )}

      {isError && (
        <div className="rounded-md bg-destructive/10 p-4 text-destructive text-sm">
          {error?.message ?? "Failed to load result"}
        </div>
      )}

      {result && <MatrixResultView result={result} />}
    </div>
  )
}
