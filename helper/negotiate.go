package helper

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const (
	MediaTypeJSON = "application/json"
	MediaTypeCSV  = "text/csv"
)

type acceptCandidate struct {
	mediaType string
	quality   float64
	order     int
}

func NegotiateMediaType(c *fiber.Ctx, supported ...string) (string, error) {
	if len(supported) == 0 {
		supported = []string{MediaTypeJSON}
	}

	acceptHeader := strings.TrimSpace(c.Get(fiber.HeaderAccept))
	if acceptHeader == "" || acceptHeader == "*/*" {
		return supported[0], nil
	}

	candidates := parseAcceptHeader(acceptHeader)
	if len(candidates) == 0 {
		return supported[0], nil
	}

	bestMatch := ""
	bestQuality := -1.0

	for _, cand := range candidates {
		if cand.quality <= 0 {
			continue
		}

		for _, sup := range supported {
			if matchesMediaType(cand.mediaType, sup) {
				if cand.quality > bestQuality {
					bestQuality = cand.quality
					bestMatch = sup
				}
				break
			}
		}
	}

	if bestMatch == "" {
		return "", ErrNotAcceptable(
			fmt.Sprintf("format Accept tidak didukung; gunakan salah satu dari: %s", strings.Join(supported, ", ")),
		)
	}

	return bestMatch, nil
}

func WriteCSVResponse(c *fiber.Ctx, filename string, headers []string, rows [][]string) error {
	cleanName := sanitizeCSVFilename(filename)

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	if len(headers) > 0 {
		if err := writer.Write(headers); err != nil {
			return ErrInternal(err)
		}
	}

	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return ErrInternal(err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return ErrInternal(err)
	}

	c.Set(fiber.HeaderContentType, "text/csv; charset=utf-8")
	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf(`attachment; filename="%s"`, cleanName))
	return c.Status(fiber.StatusOK).Send(buf.Bytes())
}

func parseAcceptHeader(header string) []acceptCandidate {
	parts := strings.Split(header, ",")
	result := make([]acceptCandidate, 0, len(parts))

	for idx, part := range parts {
		segments := strings.Split(strings.TrimSpace(part), ";")
		media := strings.ToLower(strings.TrimSpace(segments[0]))
		if media == "" {
			continue
		}

		quality := 1.0
		for _, param := range segments[1:] {
			kv := strings.SplitN(strings.TrimSpace(param), "=", 2)
			if len(kv) == 2 && strings.ToLower(strings.TrimSpace(kv[0])) == "q" {
				if parsed, err := strconv.ParseFloat(strings.TrimSpace(kv[1]), 64); err == nil {
					quality = parsed
				}
			}
		}

		result = append(result, acceptCandidate{
			mediaType: media,
			quality:   quality,
			order:     idx,
		})
	}

	return result
}

func matchesMediaType(pattern, target string) bool {
	p := strings.ToLower(strings.TrimSpace(pattern))
	t := strings.ToLower(strings.TrimSpace(target))

	if p == "*/*" || p == t {
		return true
	}

	if strings.HasSuffix(p, "/*") {
		prefix := strings.TrimSuffix(p, "*")
		return strings.HasPrefix(t, prefix)
	}

	return false
}

func sanitizeCSVFilename(filename string) string {
	base := filepath.Base(strings.TrimSpace(filename))
	base = strings.NewReplacer(`"`, "", "\r", "", "\n", "", "/", "_", `\`, "_").Replace(base)
	if base == "" || base == "." {
		return "export.csv"
	}
	if !strings.HasSuffix(strings.ToLower(base), ".csv") {
		base += ".csv"
	}
	return base
}