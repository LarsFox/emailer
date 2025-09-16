package entities

import "fmt"

// TGError — ошибка Телеграма.
type TGError struct {
	Code        int64
	Description string
}

func (e *TGError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Description)
}
