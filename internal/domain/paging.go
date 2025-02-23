package domain

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// PagingData struct contains the pagination data for a set of records
type PagingData struct {
	Cursor      string
	HasNextPage bool
	Limit       int
}

type Cursorable interface {
	Cursor() (string, error)
}

// Encodes data as a cusor
func EncodeCursor(data any) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(jsonData), nil
}

// Decodes a raw cursor string
func DecodeCursor(rawCursor string, v interface{}) error {
	jsonData, err := base64.StdEncoding.DecodeString(rawCursor)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonData, v)
}

func GetPagingData[T Cursorable](data []T, pageLimit int, prevCursor string) (*PagingData, error) {
	hasNextPage := len(data) == pageLimit

	nextCursor := prevCursor
	if len(data) > 0 {
		last := data[len(data)-1]

		var err error
		nextCursor, err = last.Cursor()
		if err != nil {
			return nil, fmt.Errorf("failed to get cursor: %w", err)
		}
	}

	paging := &PagingData{
		Cursor:      nextCursor,
		Limit:       pageLimit,
		HasNextPage: hasNextPage,
	}

	return paging, nil
}
