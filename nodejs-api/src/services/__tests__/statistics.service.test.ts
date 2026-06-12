import { StatisticsService } from "../statistics.service"

const EPSILON = 1e-10

describe("StatisticsService", () => {
  let service: StatisticsService

  beforeEach(() => {
    service = new StatisticsService()
  })

  describe("calculate", () => {
    it("should compute max, min, sum, and average across both matrices", () => {
      const q = [[1, 2], [3, 4]]
      const r = [[5, 6], [7, 8]]

      const result = service.calculate(q, r)

      expect(result.max).toBe(8)
      expect(result.min).toBe(1)
      expect(result.sum).toBe(36)
      expect(result.average).toBe(4.5)
    })

    it("should detect diagonal matrices", () => {
      const q = [[1, 0], [0, 2]]
      const r = [[5, 0, 0], [0, 6, 0], [0, 0, 7]]

      const result = service.calculate(q, r)

      expect(result.qIsDiagonal).toBe(true)
      expect(result.rIsDiagonal).toBe(true)
    })

    it("should detect non-diagonal matrices", () => {
      const q = [[1, 2], [3, 4]]
      const r = [[5, 0], [1, 6]]

      const result = service.calculate(q, r)

      expect(result.qIsDiagonal).toBe(false)
      expect(result.rIsDiagonal).toBe(false)
    })

    it("should treat near-diagonal matrices (within EPSILON) as diagonal", () => {
      const q = [[1, EPSILON / 2], [0, 2]]
      const r = [[5, 0, 0], [EPSILON / 3, 6, 0], [0, 0, 7]]

      const result = service.calculate(q, r)

      expect(result.qIsDiagonal).toBe(true)
      expect(result.rIsDiagonal).toBe(true)
    })

    it("should treat off-diagonal values above EPSILON as non-diagonal", () => {
      const q = [[1, EPSILON + 1e-15], [0, 2]]

      const result = service.calculate(q, [[1]])

      expect(result.qIsDiagonal).toBe(false)
    })

    it("should handle single-element matrices", () => {
      const q = [[42]]
      const r = [[-7]]

      const result = service.calculate(q, r)

      expect(result.max).toBe(42)
      expect(result.min).toBe(-7)
      expect(result.sum).toBe(35)
      expect(result.average).toBe(17.5)
      expect(result.qIsDiagonal).toBe(true)
      expect(result.rIsDiagonal).toBe(true)
    })

    it("should handle matrices with negative values", () => {
      const q = [[-10, -20], [-30, -40]]
      const r = [[5, 15], [25, 35]]

      const result = service.calculate(q, r)

      expect(result.max).toBe(35)
      expect(result.min).toBe(-40)
      expect(result.sum).toBe(-20)
      expect(result.average).toBe(-2.5)
    })

    it("should handle matrices of different sizes", () => {
      const q = [[1, 2, 3]]
      const r = [[4], [5], [6]]

      const result = service.calculate(q, r)

      expect(result.max).toBe(6)
      expect(result.min).toBe(1)
      expect(result.sum).toBe(21)
      expect(result.average).toBe(3.5)
    })

    it("should return average = 0 when there are no elements", () => {
      const q: number[][] = []
      const r: number[][] = []

      const result = service.calculate(q, r)

      expect(result.max).toBe(-Infinity)
      expect(result.min).toBe(Infinity)
      expect(result.sum).toBe(0)
      expect(result.average).toBe(0)
    })

    it("should handle a matrix with one row and one column being diagonal", () => {
      const q = [[1]]
      const r = [[2]]

      const result = service.calculate(q, r)

      expect(result.qIsDiagonal).toBe(true)
      expect(result.rIsDiagonal).toBe(true)
      expect(result.max).toBe(2)
      expect(result.min).toBe(1)
      expect(result.sum).toBe(3)
      expect(result.average).toBe(1.5)
    })

    it("should handle matrices with duplicate values", () => {
      const q = [[5, 5], [5, 5]]
      const r = [[5, 5], [5, 5]]

      const result = service.calculate(q, r)

      expect(result.max).toBe(5)
      expect(result.min).toBe(5)
      expect(result.sum).toBe(40)
      expect(result.average).toBe(5)
    })

    it("should treat an empty matrix (no rows) as diagonal", () => {
      const q: number[][] = []

      const result = service.calculate(q, [[1]])

      expect(result.qIsDiagonal).toBe(true)
    })
  })
})
