package search

import (
	"fmt"
	"time"
)

// SearchError struct для обработки ошибок при поиске
type SearchError struct {
	err  error
	date time.Time
}

func (e *SearchError) Error() string {
	return fmt.Sprintf("%s search error: %s", e.date, e.err)
}

func (e *SearchError) Unwrap() error {
	return e.err
}

// WrapSearchError Создание своей обёртки для ошибки
func WrapSearchError(err error) error {
	return &SearchError{
		err:  err,
		date: time.Now(),
	}
}
