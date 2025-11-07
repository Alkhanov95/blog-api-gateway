package config

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

// Проверяет, что указан существующий уровень логов
func validateLogLevel(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	var lvl slog.Level
	return lvl.UnmarshalText([]byte(val)) == nil
}

// Проверяет, что список адресов NATS валидный
func validateNatsUrls(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	re := regexp.MustCompile(`^((nats|tls)://([a-zA-Z0-9.-]+):([0-9]{1,5}),?)*$`)
	return re.MatchString(val)
}

// Проверяет валидность порта (0-65535)
func validatePort(fl validator.FieldLevel) bool {
	port, err := strconv.Atoi(fl.Field().String())
	if err != nil {
		return false
	}
	return port >= 0 && port <= 65535
}

func Validate(cfg interface{}) error {
	v := validator.New()

	if err := v.RegisterValidation("log_level", validateLogLevel); err != nil {
		return errors.Wrap(err, "register custom validation: log_level")
	}
	if err := v.RegisterValidation("nats_urls", validateNatsUrls); err != nil {
		return errors.Wrap(err, "register custom validation: nats_urls")
	}
	if err := v.RegisterValidation("port", validatePort); err != nil {
		return errors.Wrap(err, "register custom validation: port")
	}

	if err := v.Struct(cfg); err != nil {
		var messages []string
		for _, e := range err.(validator.ValidationErrors) {
			messages = append(messages, fmt.Sprintf("%s(%s)", e.StructNamespace(), e.Tag()))
		}
		return errors.New("invalid configuration: " + strings.Join(messages, ", "))
	}

	return nil
}
