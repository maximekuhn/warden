package handlerutils

import (
	"net/http"

	"github.com/maximekuhn/warden/internal/apps/web/ui/components/errors"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
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

func ToMinecraftServerNameOrReturnErrorBox(w http.ResponseWriter, name string) (valueobjects.MinecraftServerName, error) {
	serverName, err := valueobjects.NewMinecraftServerName(name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeErrorBox(w, err.Error())
		return valueobjects.MinecraftServerName{}, err
	}
	return serverName, nil
}

func writeErrorBox(w http.ResponseWriter, errMsg string) {
	if err := errors.BoxError(errMsg); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
