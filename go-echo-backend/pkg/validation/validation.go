package validation

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	english "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
)

// Validator instance
type Validator struct {
	validator *validator.Validate
	trans     ut.Translator
}

// Validate validate
func (v *Validator) Validate(i interface{}) error {
	err := v.validator.Struct(i)
	if err == nil {
		return nil
	}

	if ers, ok := err.(validator.ValidationErrors); ok {
		var e = ValidateErrors{
			ers,
		}

		var err = e.TransformError().WithTranslate(v.trans)
		var transformError = getErrorForTag(err)
		if transformError != nil {
			return transformError
		}

		return errs.New(http.StatusUnprocessableEntity, err.Error())

	}

	return nil
}

// RegisterValidation register
func (v *Validator) RegisterValidation(tag string, fc validator.Func, msg string) error {
	err := v.validator.RegisterValidation(tag, fc)
	if err != nil {
		return fmt.Errorf("RegisterValidation error for %s", tag)
	}

	err = v.validator.RegisterTranslation(tag, v.trans, registrationFunc(tag, msg), translateFunc)
	if err != nil {
		return fmt.Errorf("RegisterTranslation error for %s", tag)
	}

	return nil
}

// RegisterValidation override echo's validator
func RegisterValidation() *Validator {
	v := validator.New()
	eng := english.New()
	uni := ut.New(eng, eng)
	trans, _ := uni.GetTranslator("en")

	customValidator := &Validator{validator: v, trans: trans}

	err := en.RegisterDefaultTranslations(customValidator.validator, trans)
	if err != nil {
		panic(err)
	}

	err = customValidator.validator.RegisterTranslation("required", customValidator.trans, registrationFunc("required", "{0} is required"), translateFunc)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	err = customValidator.validator.RegisterTranslation("required_unless", customValidator.trans, registrationFunc("required_unless", "{0} is required"), translateFunc)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	err = customValidator.RegisterValidation("isURL", isURL, fmt.Sprintf("{0} is invalid URL"))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	err = customValidator.RegisterValidation("isBirthday", isBirthday, fmt.Sprintf("{0} is invalid birthday"))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	err = customValidator.RegisterValidation("isTimezone", isTimezone, fmt.Sprintf("{0} is invalid timezone"))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	err = customValidator.RegisterValidation("isPhone", isPhone, fmt.Sprintf("{0} is invalid phone"))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return customValidator
}

func registrationFunc(tag string, translation string) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) error {
		return ut.Add(tag, translation, true)
	}

}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field())

	if err != nil {
		return fe.(error).Error()
	}
	return t
}

func getErrorForTag(err *ValidateError) error {
	if err.ValidateType != "required" {
		return nil
	}

	switch err.FieldName {
	case "Email":
		if err.ValidateType == "required_without" {
			return errs.ErrPhoneOrEmailRequired
		}
		return errs.ErrEmailInvalid
	case "AppVersion":
		return errs.ErrAppVersionInvalid
	case "Phone":
		if err.ValidateType == "required_without" {
			return errs.ErrPhoneOrEmailRequired
		}
		return errs.ErrPhoneInvalid
	case "Avatar":
		return errs.ErrAvatarURLInvalid
	case "Birthday":
		return errs.ErrBirthdayInvalid
	case "Password":
		if err.ValidateType == "min" {
			return errs.ErrPasswordInvalid.WithMessage(fmt.Sprintf("Password must have at least %v but got %v", err.ExpectedResult, err.GotResult))
		}
		return errs.ErrPasswordIncorrect
	case "Lat":
		return errs.ErrLatInvalid

	case "Lng":
		return errs.ErrLngInvalid
	case "Timezone":
		return errs.ErrTimezoneInvalid

	default:

	}

	return nil
}

func hasValue(fl validator.FieldLevel) bool {
	return requireCheckFieldKind(fl, "")
}

func requireCheckFieldKind(fl validator.FieldLevel, param string) bool {
	field := fl.Field()
	if len(param) > 0 {
		if fl.Parent().Kind() == reflect.Ptr {
			field = fl.Parent().Elem().FieldByName(param)
		} else {
			field = fl.Parent().FieldByName(param)
		}
	}
	switch field.Kind() {
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return !field.IsNil()
	default:
		_, _, nullable := fl.ExtractType(field)
		if nullable && field.Interface() != nil {
			return true
		}
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}

func isEq(field reflect.Value, value string) bool {
	switch field.Kind() {

	case reflect.String:
		return field.String() == value

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(value)

		return int64(field.Len()) == p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(value)

		return field.Int() == p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(value)

		return field.Uint() == p

	case reflect.Float32, reflect.Float64:
		p := asFloat(value)

		return field.Float() == p
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

func asInt(param string) int64 {

	i, err := strconv.ParseInt(param, 0, 64)
	panicIf(err)

	return i
}

func asUint(param string) uint64 {

	i, err := strconv.ParseUint(param, 0, 64)
	panicIf(err)

	return i
}

func asFloat(param string) float64 {

	i, err := strconv.ParseFloat(param, 64)
	panicIf(err)

	return i
}

func panicIf(err error) {
	if err != nil {
		panic(err.Error())
	}
}
