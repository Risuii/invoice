package main

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

type User struct {
	Name       string `validate:"required,alphanumeric50"`
	Email      string `validate:"omitempty,email"`
	Age        int    `validate:"omitempty,numeric,gte=0,lte=99"`
	IsActive   bool   `validate:"omitempty"`
	Iat        int64  `json:"iat" validate:"required,unixtime"`
	AuthMethod string `json:"authMethod" validate:"required,oneof=00 01 02 03 04 05 06"`
	AuthCode   string `validate:"required_unless=authMethod 00 02 03 04 05"`
}

func alphanumericValidator(a int) validator.Func {
	return func(fl validator.FieldLevel) bool {
		re := regexp.MustCompile(fmt.Sprintf(`^[a-zA-Z0-9]{1,%d}$`, a))
		return re.MatchString(fl.Field().String())
	}
}

func unixtime(fl validator.FieldLevel) bool {
	iat := fl.Field().Int()
	if iat <= 0 {
		return false
	}
	return true
}

func main() {
	validate := validator.New()
	validate.RegisterValidation("alphanumeric50", alphanumericValidator(50))
	validate.RegisterValidation("unixtime", unixtime)

	user := &User{
		// Name: "1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50",
		Name:       "12345678910",
		Email:      "john.doe@example.com",
		Age:        99,
		IsActive:   false,
		Iat:        1645345600,
		AuthMethod: "01",
		AuthCode:   "00",
	}

	err := validate.Struct(user)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Validation passed!")
}
