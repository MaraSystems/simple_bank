package api

import (
	"github.com/MaraSystems/simple_bank/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		// check currency support
		return util.IsSupportedCurrency(currency)
	}
	return false
}
