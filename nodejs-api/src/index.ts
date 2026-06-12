import express from "express"
import matrixRoutes from "./routes/matrix.routes"
import { errorHandler, notFoundHandler } from "./middleware/error.middleware"

const app = express()
const PORT = process.env.PORT || "3001"

app.use(express.json())
app.use(matrixRoutes)
app.use(notFoundHandler)
app.use(errorHandler)

app.listen(PORT, () => {
  console.log(`[Node.js API] Running on http://localhost:${PORT}/`)
})
