package services

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/jordanhuaman/go-api/src/clients"
	"github.com/jordanhuaman/go-api/src/models"
)

type MatrixService struct {
	userRepo         *models.UserRepository
	matrixInputRepo  *models.MatrixInputRepository
	matrixResultRepo *models.MatrixResultRepository
	nodeClient       *clients.NodeClient
}

func NewMatrixService(
	userRepo *models.UserRepository,
	matrixInputRepo *models.MatrixInputRepository,
	matrixResultRepo *models.MatrixResultRepository,
	nodeClient *clients.NodeClient,
) *MatrixService {
	return &MatrixService{
		userRepo:         userRepo,
		matrixInputRepo:  matrixInputRepo,
		matrixResultRepo: matrixResultRepo,
		nodeClient:       nodeClient,
	}
}

func (s *MatrixService) ProcessMatrix(ctx context.Context, userID uuid.UUID, data [][]float64) (*models.MatrixResult, error) {
	if len(data) == 0 || len(data[0]) == 0 {
		return nil, errors.New("matrix cannot be empty")
	}

	rows := len(data)
	cols := len(data[0])

	for i, row := range data {
		if len(row) != cols {
			return nil, fmt.Errorf("matrix is not rectangular: row %d has %d columns, expected %d", i, len(row), cols)
		}
	}

	if rows < cols {
		return nil, errors.New("QR decomposition requires at least as many rows as columns")
	}

	input := &models.MatrixInput{
		Data:    data,
		Rows:    rows,
		Columns: cols,
	}

	q, r := gramSchmidt(data)

	var stats *models.Statistics
	var calcErr error
	if s.nodeClient != nil {
		stats, calcErr = s.nodeClient.CalculateStatistics(ctx, q, r)
		if calcErr != nil {
			return nil, fmt.Errorf("node service: %w", calcErr)
		}
	}
	if stats == nil {
		s := calculateStatistics(q, r)
		stats = &s
	}

	result := &models.MatrixResult{
		UserID:     userID,
		QRResult:   models.QRMatrices{Q: q, R: r},
		Statistics: *stats,
		Status:     "completed",
	}

	if err := s.matrixInputRepo.Create(input); err != nil {
		return nil, fmt.Errorf("save input: %w", err)
	}

	result.MatrixInputID = input.ID
	if err := s.matrixResultRepo.Create(result); err != nil {
		return nil, fmt.Errorf("save result: %w", err)
	}

	return result, nil
}

func gramSchmidt(a [][]float64) (q, r [][]float64) {
	m := len(a)
	n := len(a[0])

	q = make([][]float64, m)
	for i := range q {
		q[i] = make([]float64, n)
	}
	r = make([][]float64, n)
	for i := range r {
		r[i] = make([]float64, n)
	}

	for k := 0; k < n; k++ {
		v := make([]float64, m)
		for i := 0; i < m; i++ {
			v[i] = a[i][k]
		}

		for i := 0; i < k; i++ {
			qCol := make([]float64, m)
			for j := 0; j < m; j++ {
				qCol[j] = q[j][i]
			}
			r[i][k] = dot(qCol, v)
			for j := 0; j < m; j++ {
				v[j] -= r[i][k] * qCol[j]
			}
		}

		r[k][k] = norm(v)
		if r[k][k] != 0 {
			for i := 0; i < m; i++ {
				q[i][k] = v[i] / r[k][k]
			}
		}
	}

	return q, r
}

func dot(a, b []float64) float64 {
	var sum float64
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}

func norm(v []float64) float64 {
	return math.Sqrt(dot(v, v))
}

func calculateStatistics(q, r [][]float64) models.Statistics {
	var max, min float64
	var sum float64
	var count int
	first := true

	collect := func(m [][]float64) {
		for _, row := range m {
			for _, val := range row {
				if first {
					max = val
					min = val
					first = false
				}
				if val > max {
					max = val
				}
				if val < min {
					min = val
				}
				sum += val
				count++
			}
		}
	}

	collect(q)
	collect(r)

	var avg float64
	if count > 0 {
		avg = sum / float64(count)
	}

	return models.Statistics{
		Max:       max,
		Min:       min,
		Average:   avg,
		Sum:       sum,
		QDiagonal: isDiagonal(q),
		RDiagonal: isDiagonal(r),
	}
}

func isDiagonal(m [][]float64) bool {
	epsilon := 1e-10
	for i, row := range m {
		for j, val := range row {
			if i != j && math.Abs(val) > epsilon {
				return false
			}
		}
	}
	return true
}
