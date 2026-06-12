export interface Statistics {
  max: number
  min: number
  average: number
  sum: number
  qIsDiagonal: boolean
  rIsDiagonal: boolean
}

export interface QRRequest {
  q: number[][]
  r: number[][]
}
