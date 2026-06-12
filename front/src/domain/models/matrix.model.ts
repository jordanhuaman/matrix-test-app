export type Matrix2D = number[][]

export interface QRRequest {
  data: Matrix2D
}

export interface QRMatrices {
  q: Matrix2D
  r: Matrix2D
}

export interface Statistics {
  max: number
  min: number
  average: number
  sum: number
  qIsDiagonal: boolean
  rIsDiagonal: boolean
}

export interface MatrixResult {
  id: string
  userId: string
  matrixInputId: string
  matrixInput?: MatrixInput
  qrResult: QRMatrices
  statistics: Statistics
  status: string
  errorMsg?: string
  createdAt: string
}

export interface MatrixInput {
  id: string
  data: Matrix2D
  rows: number
  columns: number
  createdAt: string
}

export interface PaginatedResponse<T> {
  results: T[]
  total: number
  page: number
  limit: number
}
