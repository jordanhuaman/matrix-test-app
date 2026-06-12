import type { IMatrixRepository } from "../../domain/ports/matrix.repository.port"
import type { MatrixResult, PaginatedResponse, QRRequest } from "../../domain/models/matrix.model"
import { httpClient } from "./http-client"

export class MatrixRepository implements IMatrixRepository {
  async submit(data: QRRequest): Promise<MatrixResult> {
    return await httpClient.post<MatrixResult>("/matrix/qr", data)
  }

  async list(page = 1, limit = 20): Promise<PaginatedResponse<MatrixResult>> {
    return await httpClient.get<PaginatedResponse<MatrixResult>>(`/matrix/qr?page=${page}&limit=${limit}`)
  }

  async getById(id: string): Promise<MatrixResult> {
    return await httpClient.get<MatrixResult>(`/matrix/qr/${id}`)
  }
}
