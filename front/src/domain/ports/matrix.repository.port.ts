import type { MatrixResult, PaginatedResponse, QRRequest } from "../models/matrix.model"

export interface IMatrixRepository {
  submit(data: QRRequest): Promise<MatrixResult>
  list(page?: number, limit?: number): Promise<PaginatedResponse<MatrixResult>>
  getById(id: string): Promise<MatrixResult>
}
