import { useNavigate } from "@tanstack/react-router"
import { ScrollArea } from "#/components/ui/scroll-area"
import { Card, CardContent } from "#/components/ui/card"
import { Skeleton } from "#/components/ui/skeleton"
import { useMatrixResults } from "../../hooks/use-matrix"
import { History, ChevronRight, Clock } from "lucide-react"

export function HistorySidebar() {
  const { data, isLoading } = useMatrixResults(1, 20)
  const navigate = useNavigate()

  const results = data?.results ?? []

  return (
    <div className="flex flex-col h-full">
      <div className="p-4 border-b border-border">
        <div className="flex items-center gap-2 text-sm font-semibold">
          <History className="h-4 w-4" />
          History
        </div>
      </div>
      <ScrollArea className="flex-1">
        <div className="p-2 space-y-1">
          {isLoading &&
            Array.from({ length: 5 }).map((_, i) => (
              <Skeleton key={i} className="h-16 w-full rounded-lg" />
            ))}
          {results.map((r) => (
            <Card
              key={r.id}
              className="cursor-pointer hover:bg-accent/50 transition-colors border-0 shadow-none"
              onClick={() => navigate({ to: "/matrix/$id", params: { id: r.id } })}
            >
              <CardContent className="p-3 flex items-center justify-between">
                <div className="space-y-1 min-w-0">
                  <div className="flex items-center gap-2 text-xs text-muted-foreground">
                    <Clock className="h-3 w-3 shrink-0" />
                    <span className="truncate">
                      {r.matrixInput
                        ? `${r.matrixInput.rows}×${r.matrixInput.columns}`
                        : "—"}
                    </span>
                  </div>
                  <p className="text-xs text-muted-foreground truncate">
                    {r.createdAt
                      ? new Date(r.createdAt).toLocaleDateString("es-PE", {
                          day: "2-digit",
                          month: "short",
                          hour: "2-digit",
                          minute: "2-digit",
                        })
                      : "—"}
                  </p>
                </div>
                <ChevronRight className="h-4 w-4 text-muted-foreground shrink-0" />
              </CardContent>
            </Card>
          ))}
          {results.length === 0 && !isLoading && (
            <p className="text-xs text-muted-foreground text-center py-8">
              No results yet
            </p>
          )}
        </div>
      </ScrollArea>
      {results.length > 0 && (
        <div className="p-3 border-t border-border">
          <p className="text-xs text-muted-foreground text-center">
            {data?.total ?? 0} total
          </p>
        </div>
      )}
    </div>
  )
}
