package sql

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func (s *sqlReader) QueryAsTOON(ctx context.Context, section string, query string, args ...any) (string, error) {
	rows, err := s.c.QueryContext(ctx, query, args...)
	if err != nil {
		return "", fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return "", fmt.Errorf("columns: %w", err)
	}

	var b strings.Builder

	b.WriteString(section)
	b.WriteString(" {")
	b.WriteString(strings.Join(cols, ", "))
	b.WriteString("}:")

	values := make([]any, len(cols))
	raw := make([]any, len(cols))
	for i := range values {
		values[i] = &raw[i]
	}

	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			return "", fmt.Errorf("scan: %w", err)
		}

		b.WriteRune('\n')

		for i, v := range raw {
			if i > 0 {
				b.WriteRune(' ')
			}
			val := formatToonValue(v)

			b.WriteRune('"')
			b.WriteString(val)
			b.WriteRune('"')
		}
	}

	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("iterate rows: %w", err)
	}

	return b.String(), nil
}

func formatToonValue(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case []byte:
		return sanitizeToonValue(string(x))
	case time.Time:
		return x.UTC().Format(time.RFC3339)
	default:
		return sanitizeToonValue(fmt.Sprint(x))
	}
}

func sanitizeToonValue(s string) string {
	return strings.NewReplacer(
		"\n", " ",
		"\r", " ",
		`"`, `'`,
	).Replace(s)
}
