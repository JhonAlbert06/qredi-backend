package utils

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func IsDominicanIDValid(cedula string) bool {

	// Eliminar los guiones en caso de que los haya
	cedula = strings.ReplaceAll(cedula, "-", "")

	suma := 0

	peso := []int{1, 2, 1, 2, 1, 2, 1, 2, 1, 2}

	if len(cedula) != 11 {
		return false
	}

	_, err := strconv.Atoi(cedula)
	if err != nil {
		return false
	}

	for i := 0; i < 10; i++ {
		a, _ := strconv.Atoi(string(cedula[i]))
		b := peso[i]

		mult := strconv.Itoa(a * b)

		if len(mult) > 1 {
			a, _ = strconv.Atoi(string(mult[0]))
			b, _ = strconv.Atoi(string(mult[1]))
		} else {
			a = 0
			b, _ = strconv.Atoi(mult)
		}

		suma = suma + a + b
	}

	digito, _ := strconv.Atoi(string(cedula[10]))

	digitoVerificador := (10 - (suma % 10)) % 10

	return digito == digitoVerificador
}

func IsPasswordValid(password string) bool {
	return len(password) > 1 && len(password) < 40
}

func IsPhoneNumberValid(phoneNumber string) bool {

	// Verificar si el número de teléfono tiene exactamente 10 caracteres
	if len(phoneNumber) != 10 {
		return false
	}

	// Verificar si todos los caracteres son dígitos numéricos
	for _, char := range phoneNumber {
		if !unicode.IsDigit(char) {
			return false
		}
	}

	return true
}

func IsStringEmpty(s string) bool {
	return len(s) == 0
}

func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func FormatAndConvertFloat32(value float32) float32 {
	formattedValue := fmt.Sprintf("%.2f", value)

	result, err := strconv.ParseFloat(formattedValue, 32)
	if err != nil {
		//return 0.0, err
	}

	return float32(result)
}
