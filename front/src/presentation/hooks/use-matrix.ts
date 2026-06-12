import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { MatrixUseCases } from "../../application/matrix/matrix.usecases"
import { MatrixRepository } from "../../infrastructure/http/matrix.repository"
import type { PaginatedResponse, MatrixResult } from "../../domain/models/matrix.model"

const matrixRepo = new MatrixRepository()
const matrixUseCases = new MatrixUseCases(matrixRepo)

export function useSubmitMatrix() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: number[][]) => matrixUseCases.submit(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["matrix-results"] })
    },
  })
}

export function useMatrixResults(page = 1, limit = 20) {
  return useQuery<PaginatedResponse<MatrixResult>>({
    queryKey: ["matrix-results", page, limit],
    queryFn: () => matrixUseCases.list(page, limit),
  })
}

export function useMatrixResult(id: string) {
  return useQuery<MatrixResult>({
    queryKey: ["matrix-result", id],
    queryFn: () => matrixUseCases.getById(id),
    enabled: !!id,
  })
}
