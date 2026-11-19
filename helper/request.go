package helper

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Alan-00280/EventSeat-PBL-UAS.git/app/model"
	"github.com/gofiber/fiber/v2"
)

// RequestContext memberi timeout 5 detik untuk setiap operasi database,
// diturunkan dari context request Fiber agar propagasi cancel berjalan benar.
func RequestContext(c *fiber.Ctx) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.UserContext(), 5*time.Second)
}

// ParamUUID membaca parameter rute berdasarkan nama key dan memvalidasi
// bahwa nilainya tidak kosong dan memiliki panjang UUID yang benar (36 karakter).
// Mengembalikan string UUID dan true jika valid, string kosong dan false jika tidak.
func ParamUUID(c *fiber.Ctx, key string) (string, bool) {
	id := c.Params(key)
	if len(id) != 36 {
		return "", false
	}
	return id, true
}

// allowedSort adalah daftar putih kolom yang boleh dipakai sebagai sort key.
// Ini mencegah SQL injection melalui query string ?sort=... pada layer repository.
var allowedSort = map[string]bool{
	"id":         true,
	"created_at": true,
	"updated_at": true,
	"username":   true,
	"email":      true,
}

// ParseListQuery membaca query string dan menghasilkan ListQuery dengan nilai
// bawaan yang aman. Nilai di luar rentang wajar dikoreksi diam-diam karena
// ini preferensi tampilan — client tidak perlu tahu bahwa limit-nya dipotong.
func ParseListQuery(c *fiber.Ctx) model.ListQuery {
	q := model.ListQuery{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
		Search: strings.TrimSpace(c.Query("search")),
		Sort:   c.Query("sort", "created_at"),
		Order:  strings.ToLower(c.Query("order", "desc")),
	}

	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 {
		q.Limit = 10
	}
	if q.Limit > 100 {
		q.Limit = 100
	}
	if !allowedSort[q.Sort] {
		q.Sort = "created_at"
	}
	if q.Order != "asc" {
		q.Order = "desc"
	}
	if raw := c.Query("is_active"); raw != "" {
		if v, err := strconv.ParseBool(raw); err == nil {
			q.IsActive = &v
		}
	}

	return q
}

// ParseCursorQuery membaca query string untuk paginasi berbasis cursor.
// Cursor rusak atau is_active yang bukan boolean dikembalikan sebagai error 400,
// bukan dikoreksi diam-diam — client berhak tahu bahwa request-nya tidak valid.
func ParseCursorQuery(c *fiber.Ctx) (model.CursorQuery, error) {
	q := model.CursorQuery{
		Limit:  c.QueryInt("limit", 10),
		Search: strings.TrimSpace(c.Query("search")),
	}

	if q.Limit < 1 {
		q.Limit = 10
	}
	if q.Limit > 100 {
		q.Limit = 100
	}

	if raw := c.Query("is_active"); raw != "" {
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return model.CursorQuery{}, NewAppError(
				http.StatusBadRequest, "BAD_REQUEST", "is_active harus berupa boolean", nil,
			)
		}
		q.IsActive = &v
	}

	if raw := strings.TrimSpace(c.Query("cursor")); raw != "" {
		cursor, err := DecodeCursor(raw)
		if err != nil {
			return model.CursorQuery{}, NewAppError(
				http.StatusBadRequest, "BAD_REQUEST", "cursor tidak valid atau sudah kedaluwarsa", nil,
			)
		}
		q.After = &cursor
	}

	return q, nil
}
