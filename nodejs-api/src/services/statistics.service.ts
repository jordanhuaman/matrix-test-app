import { Statistics } from "../interfaces/matrix.interface"

const EPSILON = 1e-10

function isDiagonal(m: number[][]): boolean {
  for (let i = 0; i < m.length; i++) {
    for (let j = 0; j < m[i].length; j++) {
      if (i !== j && Math.abs(m[i][j]) > EPSILON) {
        return false
      }
    }
  }
  return true
}

export class StatisticsService {
  calculate(q: number[][], r: number[][]): Statistics {
    let max = -Infinity
    let min = Infinity
    let sum = 0
    let count = 0

    const accumulate = (m: number[][]) => {
      for (const row of m) {
        for (const val of row) {
          if (val > max) max = val
          if (val < min) min = val
          sum += val
          count++
        }
      }
    }

    accumulate(q)
    accumulate(r)

    const average = count > 0 ? sum / count : 0

    return {
      max,
      min,
      average,
      sum,
      qIsDiagonal: isDiagonal(q),
      rIsDiagonal: isDiagonal(r),
    }
  }
}
