package helper

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"
)

type TimeIDCursor struct {
	CreatedAt time.Time `json:"c"`
	ID        string    `json:"i"`
}


func EncodeTimeIDCursor(createdAt time.Time, id string) (string, error) {
	cleanID := strings.TrimSpace(id)
	if createdAt.IsZero() || cleanID == "" {
		return "", ErrBadRequest("INVALID_CURSOR", "data cursor tidak valid")
	}

	payload := TimeIDCursor{
		CreatedAt: createdAt.UTC(),
		ID:        cleanID,
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return "", ErrInternal(err)
	}

	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func DecodeTimeIDCursor(encoded string) (TimeIDCursor, error) {
	trimmed := strings.TrimSpace(encoded)
	if trimmed == "" {
		return TimeIDCursor{}, ErrBadRequest("INVALID_CURSOR", "cursor tidak boleh kosong")
	}

	raw, err := base64.RawURLEncoding.DecodeString(trimmed)
	if err != nil {
		return TimeIDCursor{}, ErrBadRequest("INVALID_CURSOR", "format cursor tidak valid")
	}

	var cursor TimeIDCursor
	if err := json.Unmarshal(raw, &cursor); err != nil {
		return TimeIDCursor{}, ErrBadRequest("INVALID_CURSOR", "struktur payload cursor tidak valid")
	}

	cursor.ID = strings.TrimSpace(cursor.ID)
	if cursor.CreatedAt.IsZero() || cursor.ID == "" {
		return TimeIDCursor{}, ErrBadRequest("INVALID_CURSOR", "isi cursor tidak lengkap")
	}

	cursor.CreatedAt = cursor.CreatedAt.UTC()
	return cursor, nil
}

type ObjectIDCursor struct {
	ID string `json:"id"`
}

func EncodeObjectIDCursor(objectIDHex string) (string, error) {
	cleanID := strings.ToLower(strings.TrimSpace(objectIDHex))
	if !isValidHexObjectID(cleanID) {
		return "", ErrBadRequest("INVALID_CURSOR", "object id cursor tidak valid")
	}

	raw, err := json.Marshal(ObjectIDCursor{ID: cleanID})
	if err != nil {
		return "", ErrInternal(err)
	}

	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func DecodeObjectIDCursor(encoded string) (string, error) {
	trimmed := strings.TrimSpace(encoded)
	if trimmed == "" {
		return "", ErrBadRequest("INVALID_CURSOR", "cursor tidak boleh kosong")
	}

	raw, err := base64.RawURLEncoding.DecodeString(trimmed)
	if err != nil {
		return "", ErrBadRequest("INVALID_CURSOR", "format cursor tidak valid")
	}

	var cursor ObjectIDCursor
	if err := json.Unmarshal(raw, &cursor); err != nil {
		return ObjectIDCursor{}, ErrBadRequest("INVALID_CURSOR", "struktur payload cursor tidak valid")
	}

	cleanID := strings.ToLower(strings.TrimSpace(cursor.ID))
	if !isValidHexObjectID(cleanID) {
		return "", ErrBadRequest("INVALID_CURSOR", "nilai object id pada cursor tidak valid")
	}

	return cleanID, nil
}

func NormalizePageLimit(limit, defaultLimit, maxLimit int) int {
	if defaultLimit <= 0 {
		defaultLimit = 10
	}
	if maxLimit < defaultLimit {
		maxLimit = 50
	}
	if limit <= 0 {
		return defaultLimit
	}
	if limit > maxLimit {
		return maxLimit
	}
	return limit
}

func isValidHexObjectID(s string) bool {
	if len(s) != 24 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}