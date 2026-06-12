import { Request, Response, NextFunction } from "express"
import { StatisticsService } from "../services/statistics.service"
import { QRRequest } from "../interfaces/matrix.interface"

const statisticsService = new StatisticsService()

export function calculate(req: Request, res: Response, next: NextFunction): void {
  const { q, r } = req.body as QRRequest

  if (!Array.isArray(q) || !Array.isArray(r) || q.length === 0 || r.length === 0) {
    res.status(400).json({
      status: "error",
      message: "q and r must be non-empty arrays",
      data: null,
    })
    return
  }

  const validateMatrix = (m: number[][], name: string): string | null => {
    if (!Array.isArray(m)) return `${name} must be an array`
    if (m.length === 0) return `${name} cannot be empty`
    const cols = m[0].length
    for (let i = 0; i < m.length; i++) {
      if (!Array.isArray(m[i])) return `${name}[${i}] must be an array`
      if (m[i].length !== cols) return `${name}[${i}] has inconsistent column count`
    }
    return null
  }

  const qErr = validateMatrix(q, "q")
  if (qErr) {
    res.status(400).json({ status: "error", message: qErr, data: null })
    return
  }

  const rErr = validateMatrix(r, "r")
  if (rErr) {
    res.status(400).json({ status: "error", message: rErr, data: null })
    return
  }

  try {
    const statistics = statisticsService.calculate(q, r)
    res.json(statistics)
  } catch (err) {
    next(err)
  }
}
