package repository

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// Setinel Errs
var (
	ErrNotFound         = errors.New("data not found")
	ErrDuplicate        = errors.New("data already exists")
	ErrNotEnrolled      = errors.New("tidak terdaftar di kelas ini")
	ErrCourseFull       = errors.New("kuota kelas sudah penuh")
	ErrVersionConflict  = errors.New("data sudah diubah pihak lain")
	ErrQuotaBelowFilled = errors.New("kuota lebih kecil dari jumlah peserta")
)

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
