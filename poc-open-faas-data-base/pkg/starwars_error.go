package pkg

import (
	"fmt"
	"net/http"
)

type GenericException struct {
	Msj        string
	StatusCode int
}

func (b GenericException) Error() string {
	return fmt.Sprintf("%d - %s", b.StatusCode, b.Msj)
}

func GetErrorDetail(error error) (int, string) {
	switch errorType := error.(type) {
	case GenericException:
		return errorType.StatusCode, errorType.Msj
	default:
		return http.StatusInternalServerError, errorType.Error()
	}
}
