import { useCallback, type ChangeEvent } from "react"
import { Button } from "#/components/ui/button"
import { Input } from "#/components/ui/input"

interface MatrixGridProps {
  data: number[][]
  onChange: (data: number[][]) => void
  readOnly?: boolean
}

export function MatrixGrid({ data, onChange, readOnly = false }: MatrixGridProps) {
  const rows = data.length
  const cols = data[0]?.length ?? 0

  const handleCellChange = useCallback(
    (r: number, c: number, value: string) => {
      if (readOnly) return
      const next = data.map((row) => [...row])
      next[r][c] = value === "" ? 0 : Number.parseFloat(value)
      onChange(next)
    },
    [data, onChange, readOnly],
  )

  const addRow = useCallback(() => {
    const next = [...data.map((row) => [...row]), Array(cols).fill(0)]
    onChange(next)
  }, [data, cols, onChange])

  const addCol = useCallback(() => {
    const next = data.map((row) => [...row, 0])
    onChange(next)
  }, [data, onChange])

  const removeRow = useCallback(() => {
    if (rows <= 1) return
    onChange(data.slice(0, -1))
  }, [data, rows, onChange])

  const removeCol = useCallback(() => {
    if (cols <= 1) return
    onChange(data.map((row) => row.slice(0, -1)))
  }, [data, cols, onChange])

  return (
    <div className="space-y-3">
      <div className="overflow-x-auto">
        <table className="border-collapse">
          <thead>
            <tr>
              <th className="w-8" />
              {Array.from({ length: cols }).map((_, c) => (
                <th key={c} className="text-xs text-muted-foreground font-medium px-1 pb-1 text-center">
                  Col {c + 1}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {data.map((row, r) => (
              <tr key={r}>
                <td className="text-xs text-muted-foreground font-medium pr-2 text-right align-middle">
                  {r + 1}
                </td>
                {row.map((val, c) => (
                  <td key={c} className="p-0.5">
                    <Input
                      type="number"
                      step="any"
                      value={Number.isFinite(val) ? val : 0}
                      onChange={(e: ChangeEvent<HTMLInputElement>) => handleCellChange(r, c, e.target.value)}
                      className="w-20 h-9 text-center text-sm [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
                      readOnly={readOnly}
                    />
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {!readOnly && (
        <div className="flex items-center gap-2 flex-wrap">
          <Button type="button" variant="outline" size="sm" onClick={addRow}>
            + Row
          </Button>
          <Button type="button" variant="outline" size="sm" onClick={addCol}>
            + Col
          </Button>
          {rows > 1 && (
            <Button type="button" variant="ghost" size="sm" onClick={removeRow}>
              - Row
            </Button>
          )}
          {cols > 1 && (
            <Button type="button" variant="ghost" size="sm" onClick={removeCol}>
              - Col
            </Button>
          )}
          <div className="ml-auto text-xs text-muted-foreground">
            {rows} × {cols}
          </div>
        </div>
      )}
    </div>
  )
}
