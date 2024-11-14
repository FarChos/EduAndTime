package graph

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

// MarshalLong convierte un int64 en un tipo Marshalable
func MarshalLong(i int64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.FormatInt(i, 10)))
	})
}

// UnmarshalLong convierte un tipo GraphQL en un int64
func UnmarshalLong(v interface{}) (int64, error) {
	switch v := v.(type) {
	case int:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("Long must be an integer")
	}
}
