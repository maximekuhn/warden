package handlerutils

import (
	"net/http"

	"github.com/maximekuhn/warden/internal/domain/valueobjects"
	"github.com/maximekuhn/warden/internal/ui/components/errors"
)

func ToEmailOrReturnErrorBox(w http.ResponseWriter, emailStr string) (valueobjects.Email, error) {
	email, err := valueobjects.NewEmail(emailStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeErrorBox(w, err.Error())
		return valueobjects.Email{}, err
	}
	return email, nil
}

func ToPasswordOrReturnErrorBox(w http.ResponseWriter, passwordStr string) (valueobjects.Password, error) {
	password, err := valueobjects.NewPassword(passwordStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeErrorBox(w, err.Error())
		return valueobjects.Password{}, err
	}
	return password, nil
}

func writeErrorBox(w http.ResponseWriter, errMsg string) {
	if err := errors.BoxError(errMsg); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
