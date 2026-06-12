import { Router } from "express"
import { calculate } from "../controllers/matrix.controller"

const router = Router()

router.post("/api/matrix/statistics", calculate)

export default router
