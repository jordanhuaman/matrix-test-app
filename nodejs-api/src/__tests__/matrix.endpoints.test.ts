import request from "supertest"
import app from "../index"

describe("POST /api/matrix/statistics", () => {
  it("should return 200 with statistics for valid matrices", async () => {
    const response = await request(app)
      .post("/api/matrix/statistics")
      .send({ q: [[1, 2], [3, 4]], r: [[5, 6], [7, 8]] })

    expect(response.status).toBe(200)
    expect(response.body).toEqual({
      max: 8,
      min: 1,
      average: 4.5,
      sum: 36,
      qIsDiagonal: false,
      rIsDiagonal: false,
    })
  })

  it("should detect diagonal matrices in response", async () => {
    const response = await request(app)
      .post("/api/matrix/statistics")
      .send({ q: [[1, 0], [0, 2]], r: [[3, 0], [0, 4]] })

    expect(response.status).toBe(200)
    expect(response.body.qIsDiagonal).toBe(true)
    expect(response.body.rIsDiagonal).toBe(true)
  })

  it("should return 400 when q is missing", async () => {
    const response = await request(app)
      .post("/api/matrix/statistics")
      .send({ r: [[1, 2], [3, 4]] })

    expect(response.status).toBe(400)
    expect(response.body).toEqual({
      status: "error",
      message: "q and r must be non-empty arrays",
      data: null,
    })
  })

  it("should return 400 when r is missing", async () => {
    const response = await request(app)
      .post("/api/matrix/statistics")
      .send({ q: [[1, 2], [3, 4]] })

    expect(response.status).toBe(400)
    expect(response.body).toEqual({
      status: "error",
      message: "q and r must be non-empty arrays",
      data: null,
    })
  })

  it("should return 400 when q is an empty array", async () => {
    const response = await request(app)
      .post("/api/matrix/statistics")
      .send({ q: [], r: [[1, 2], [3, 4]] })

    expect(response.status).toBe(400)
  })

  it("should return 400 when r is an empty array", async () => {
    const response = await request(app)
      .post("/api/matrix/statistics")
      .send({ q: [[1, 2], [3, 4]], r: [] })

    expect(response.status).toBe(400)
  })

  it("should return 400 when q is not an array", async () => {
    const response = await request(app)
      .post("/api/matrix/statistics")
      .send({ q: "not-an-array", r: [[1, 2], [3, 4]] })

    expect(response.status).toBe(400)
    expect(response.body.message).toMatch(/must be non-empty arrays/)
  })

  it("should return 400 when r is not an array", async () => {
    const response = await request(app)
      .post("/api/matrix/statistics")
      .send({ q: [[1, 2], [3, 4]], r: 42 })

    expect(response.status).toBe(400)
    expect(response.body.message).toMatch(/must be non-empty arrays/)
  })

  it("should return 400 when a row in q is not an array", async () => {
    const response = await request(app)
      .post("/api/matrix/statistics")
      .send({ q: [[1, 2], "bad-row"], r: [[5, 6], [7, 8]] })

    expect(response.status).toBe(400)
    expect(response.body.message).toContain("q[1] must be an array")
  })

  it("should return 400 when a row in r is not an array", async () => {
    const response = await request(app)
      .post("/api/matrix/statistics")
      .send({ q: [[1, 2], [3, 4]], r: [[5, 6], 42] })

    expect(response.status).toBe(400)
    expect(response.body.message).toContain("r[1] must be an array")
  })

  it("should return 400 when q has inconsistent column count", async () => {
    const response = await request(app)
      .post("/api/matrix/statistics")
      .send({ q: [[1, 2], [3, 4, 5]], r: [[5, 6], [7, 8]] })

    expect(response.status).toBe(400)
    expect(response.body.message).toContain("q[1] has inconsistent column count")
  })

  it("should return 400 when r has inconsistent column count", async () => {
    const response = await request(app)
      .post("/api/matrix/statistics")
      .send({ q: [[1, 2], [3, 4]], r: [[5, 6], [7]] })

    expect(response.status).toBe(400)
    expect(response.body.message).toContain("r[1] has inconsistent column count")
  })

  it("should return 400 when q has a row with 0 columns", async () => {
    const response = await request(app)
      .post("/api/matrix/statistics")
      .send({ q: [[], [1, 2]], r: [[5, 6], [7, 8]] })

    expect(response.status).toBe(400)
  })
})

describe("404 handling", () => {
  it("should return 404 for unknown routes", async () => {
    const response = await request(app).get("/api/unknown")

    expect(response.status).toBe(404)
    expect(response.body).toEqual({
      status: "error",
      message: "Route not found",
      data: null,
    })
  })
})
