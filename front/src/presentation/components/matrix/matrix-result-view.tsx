import { Badge } from "#/components/ui/badge"
import { Card, CardContent, CardHeader, CardTitle } from "#/components/ui/card"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "#/components/ui/table"
import { Check, X } from "lucide-react"
import type { MatrixResult } from "../../../domain/models/matrix.model"

interface MatrixResultViewProps {
  result: MatrixResult
}

function formatValue(v: number): string {
  if (Number.isInteger(v)) return v.toString()
  return v.toFixed(4)
}

function MatrixTable({ data, label }: { data: number[][]; label: string }) {
  return (
    <Card className="flex-1">
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-semibold">{label}</CardTitle>
      </CardHeader>
      <CardContent className="overflow-x-auto">
        <Table>
          <TableHeader>
            <TableRow>
              {data[0]?.map((_, c) => (
                <TableHead key={c} className="text-right text-xs">{c + 1}</TableHead>
              ))}
            </TableRow>
          </TableHeader>
          <TableBody>
            {data.map((row, r) => (
              <TableRow key={r}>
                {row.map((val, c) => (
                  <TableCell key={c} className="text-right font-mono text-xs tabular-nums">
                    {formatValue(val)}
                  </TableCell>
                ))}
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  )
}

export function MatrixResultView({ result }: MatrixResultViewProps) {
  const { qrResult, statistics } = result

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row gap-4">
        <MatrixTable data={qrResult.q} label="Q Matrix (Orthogonal)" />
        <MatrixTable data={qrResult.r} label="R Matrix (Upper Triangular)" />
      </div>

      <Card>
        <CardHeader className="pb-2">
          <CardTitle className="text-sm font-semibold">Statistics</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
            <div className="space-y-1">
              <p className="text-xs text-muted-foreground">Max</p>
              <p className="text-lg font-bold">{formatValue(statistics.max)}</p>
            </div>
            <div className="space-y-1">
              <p className="text-xs text-muted-foreground">Min</p>
              <p className="text-lg font-bold">{formatValue(statistics.min)}</p>
            </div>
            <div className="space-y-1">
              <p className="text-xs text-muted-foreground">Average</p>
              <p className="text-lg font-bold">{formatValue(statistics.average)}</p>
            </div>
            <div className="space-y-1">
              <p className="text-xs text-muted-foreground">Sum</p>
              <p className="text-lg font-bold">{formatValue(statistics.sum)}</p>
            </div>
          </div>
          <div className="flex gap-4 mt-4 pt-4 border-t border-border">
            <div className="flex items-center gap-2">
              <span className="text-xs text-muted-foreground">Q Diagonal:</span>
              <Badge variant={statistics.qIsDiagonal ? "default" : "secondary"}>
                {statistics.qIsDiagonal ? <Check className="h-3 w-3" /> : <X className="h-3 w-3" />}
                <span className="ml-1">{statistics.qIsDiagonal ? "Yes" : "No"}</span>
              </Badge>
            </div>
            <div className="flex items-center gap-2">
              <span className="text-xs text-muted-foreground">R Diagonal:</span>
              <Badge variant={statistics.rIsDiagonal ? "default" : "secondary"}>
                {statistics.rIsDiagonal ? <Check className="h-3 w-3" /> : <X className="h-3 w-3" />}
                <span className="ml-1">{statistics.rIsDiagonal ? "Yes" : "No"}</span>
              </Badge>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
