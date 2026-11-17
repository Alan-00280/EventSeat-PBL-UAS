package helper

import (
	"net/http"

	"github.com/Alan-00280/EventSeat-PBL-UAS.git/app/model"
	"github.com/gofiber/fiber/v2"
)

// Success mengirim respons JSON sukses dengan status code dan data yang diberikan.
func Success(c *fiber.Ctx, status int, message string, data any) error {
	return c.Status(status).JSON(model.WebResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessList mengirim respons JSON sukses untuk data berupa daftar dengan
// informasi paginasi berbasis offset (page/limit) di dalam field meta.
func SuccessList(c *fiber.Ctx, message string, data any, meta *model.Meta) error {
	return c.Status(http.StatusOK).JSON(model.WebResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// SuccessCursor mengirim respons JSON sukses untuk data berupa daftar dengan
// informasi paginasi berbasis cursor di dalam field meta.
// Metanya memakai CursorMeta (limit, next_cursor, has_more) — tanpa Total dan
// TotalPages, karena keduanya menuntut COUNT(*) atas seluruh tabel, persis biaya
// yang ingin dihindari keyset pagination.
func SuccessCursor(c *fiber.Ctx, message string, data any, meta *model.CursorMeta) error {
	return c.Status(http.StatusOK).JSON(model.WebResponseCursor{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// Created mengirim respons 201 Created sekaligus memasang header Location
// yang menunjuk ke resource yang baru dibuat.
func Created(c *fiber.Ctx, message string, data any, location string) error {
	c.Set("Location", location)
	return c.Status(http.StatusCreated).JSON(model.WebResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// NoContent mengirim respons 204 No Content tanpa body,
// dipakai untuk operasi DELETE atau aksi yang tidak menghasilkan data baru.
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}