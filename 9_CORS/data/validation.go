package data

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	validator.FieldError
}
type ValidationErrors []ValidationError

type Validation struct {
	validate *validator.Validate
}

func (v ValidationError) Error() string {
	return fmt.Sprintf(
		"Key: '%s' Error: Field Validation for '%s' failed on the '%s tag",
		v.Namespace(),
		v.Field(),
		v.Tag(),
	)
}

func (v ValidationErrors) Errors() []string {
	errs := []string{}
	for _, err := range v {
		errs = append(errs, err.Error())
	}
	return errs
}

func NewValidation() *Validation {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return &Validation{validate}
}

// Validate the item
// for more detail the returned error can be cast into a
// validator.ValidationErrors collection
//
//	if ve, ok := err.(validator.ValidationErrors); ok {
//				fmt.Println(ve.Namespace())
//				fmt.Println(ve.Field())
//				fmt.Println(ve.StructNamespace())
//				fmt.Println(ve.StructField())
//				fmt.Println(ve.Tag())
//				fmt.Println(ve.ActualTag())
//				fmt.Println(ve.Kind())
//				fmt.Println(ve.Type())
//				fmt.Println(ve.Value())
//				fmt.Println(ve.Param())
//				fmt.Println()
//		}
func (v *Validation) Validate(i interface{}) ValidationErrors {
	err := v.validate.Struct(i)
	if err == nil {
		return nil
	}
	// Not needed anymore if errs == nil is checked
	// if len(errs) == 0 {
	// 	return nil
	// }
	errs := err.(validator.ValidationErrors)
	var returnErrs []ValidationError

	for _, fieldErr := range errs {
		// cast the FieldError into our ValidationError and append to the slice
		ve := ValidationError{fieldErr.(validator.FieldError)}
		returnErrs = append(returnErrs, ve)
	}

	return returnErrs
}

func validateSKU(f1 validator.FieldLevel) bool {
	// sku is of format abc-def-gec
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	sku := re.FindAllString(f1.Field().String(), -1)

	return len(sku) == 1
}
