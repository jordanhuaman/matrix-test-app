import { Request, Response, NextFunction } from "express"

export function errorHandler(err: Error, _req: Request, res: Response, _next: NextFunction): void {
  console.error("[Error]", err.message)
  res.status(500).json({
    status: "error",
    message: "Internal server error",
    data: null,
  })
}

export function notFoundHandler(_req: Request, res: Response): void {
  res.status(404).json({
    status: "error",
    message: "Route not found",
    data: null,
  })
}
