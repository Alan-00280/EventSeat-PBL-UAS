package helper

import (
	"errors"
	"reflect"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// validate dibuat SEKALI untuk seluruh aplikasi.
//
// validator.New() melakukan refleksi dan menyimpan hasilnya dalam cache
// internal. Membuatnya ulang pada setiap request berarti membuang cache
// tersebut berkali-kali — mahal dan tidak ada gunanya.
type AppValidator struct {
	validate          *validator.Validate
	passwordCommonSet *PasswordCommonSet
}

func NewValidator(passwordCommonSet *PasswordCommonSet) *AppValidator {
	v := validator.New()

	// Tanpa ini, pesan error menyebut nama field Go (misalnya "Username"),
	// padahal client mengirim dan membaca nama JSON (misalnya "username").
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "" || name == "-" {
			return field.Name
		}
		return name
	})

	// Aturan custom: tidak boleh mengandung spasi, tab, atau newline.
	_ = v.RegisterValidation("nospace", func(fl validator.FieldLevel) bool {
		return !strings.ContainsAny(fl.Field().String(), " \t\n\r")
	})

	// Aturan custom: username hanya boleh berisi huruf, angka, titik, dan underscore.
	_ = v.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		for _, r := range fl.Field().String() {
			if !unicode.IsLetter(r) &&
				!unicode.IsDigit(r) &&
				r != '.' &&
				r != '_' {
				return false
			}
		}
		return true
	})

	// Aturan custom: password harus memenuhi syarat kekuatan minimum.
	_ = v.RegisterValidation("strongpassword", func(fl validator.FieldLevel) bool {
		return checkPasswordStrength(fl.Field().String(), *passwordCommonSet) == ""
	})

	return &AppValidator{
		validate:          v,
		passwordCommonSet: passwordCommonSet,
	}
}

// ValidateStruct menjalankan seluruh aturan validasi pada tag struct dan
// mengembalikan map nama field ke pesan error berbahasa Indonesia.
// Mengembalikan nil berarti tidak ada pelanggaran.
func ValidateStruct(s any, appValidator AppValidator) map[string]string {
	validate := appValidator.validate

	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var invalid *validator.InvalidValidationError
	if errors.As(err, &invalid) {
		return map[string]string{"_": "objek yang divalidasi tidak sah"}
	}

	var fieldErrors validator.ValidationErrors
	if !errors.As(err, &fieldErrors) {
		return map[string]string{"_": "validasi gagal"}
	}

	result := make(map[string]string, len(fieldErrors))
	for _, fe := range fieldErrors {
		if _, exists := result[fe.Field()]; !exists {
			result[fe.Field()] = messageFor(fe, appValidator)
		}
	}

	return result
}

// messageFor menerjemahkan nama tag validasi menjadi kalimat yang dapat dibaca
// pengguna dalam bahasa Indonesia. Daftar ini terpusat: menambah satu tag baru
// cukup menambah satu case di sini, tidak menyebar ke banyak file.
func messageFor(fe validator.FieldError, appValidator AppValidator) string {
	switch fe.Tag() {
	case "required":
		return "wajib diisi"
	case "email":
		return "format email tidak valid"
	case "min":
		switch fe.Kind() {
		case reflect.String:
			return "minimal " + fe.Param() + " karakter"
		case reflect.Slice, reflect.Array, reflect.Map:
			return "minimal " + fe.Param() + " item"
		}
		return "nilai minimal " + fe.Param()
	case "max":
		switch fe.Kind() {
		case reflect.String:
			return "maksimal " + fe.Param() + " karakter"
		case reflect.Slice, reflect.Array, reflect.Map:
			return "maksimal " + fe.Param() + " item"
		}
		return "nilai maksimal " + fe.Param()
	case "alphanum":
		return "hanya boleh berisi huruf dan angka"
	case "nospace":
		return "tidak boleh mengandung spasi"
	case "username":
		return "hanya boleh huruf, angka, titik, dan garis bawah"
	case "strongpassword":
		if value, ok := fe.Value().(string); ok {
			return checkPasswordStrength(value, *appValidator.passwordCommonSet)
		}
		return "password tidak memenuhi syarat"
	case "oneof":
		return "harus salah satu dari: " +
			strings.ReplaceAll(fe.Param(), " ", ", ")
	default:
		// Jaring pengaman. Bila muncul di log, artinya ada tag yang
		// dipakai tetapi belum diterjemahkan di sini.
		return "tidak memenuhi aturan " + fe.Tag()
	}
}

// checkPasswordStrength memeriksa kekuatan password dan mengembalikan pesan error
// jika tidak memenuhi syarat. Mengembalikan string kosong berarti password valid.
func checkPasswordStrength(password string, pwCommon PasswordCommonSet) string {
	if len(password) < minPasswordLength {
		return "minimum 8 character of password"
	}

	var hasLetter, hasDigit bool = false, false
	for _, c := range password {
		switch {
		case unicode.IsLetter(c):
			hasLetter = true
		case unicode.IsDigit(c):
			hasDigit = true
		}
	}

	if !hasLetter || !hasDigit {
		return "passsword must contain mix of letter and numbers"
	}

	commonSet := pwCommon.PasswordSet
	if _, exists := commonSet[password]; exists {
		return "password's too common"
	}

	return ""

}
