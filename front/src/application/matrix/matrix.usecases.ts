import type { IMatrixRepository } from "../../domain/ports/matrix.repository.port"
import type { MatrixResult, PaginatedResponse, QRRequest } from "../../domain/models/matrix.model"

export class MatrixUseCases {
  constructor(private readonly matrixRepo: IMatrixRepository) {}

  async submit(data: number[][]): Promise<MatrixResult> {
    return await this.matrixRepo.submit({ data } as QRRequest)
  }

  async list(page = 1, limit = 20): Promise<PaginatedResponse<MatrixResult>> {
    return await this.matrixRepo.list(page, limit)
  }

  async getById(id: string): Promise<MatrixResult> {
    return await this.matrixRepo.getById(id)
  }
}
