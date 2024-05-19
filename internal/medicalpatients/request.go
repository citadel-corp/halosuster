package medicalpatients

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var phoneNumberValidationRule = validation.NewStringRule(func(s string) bool {
	return strings.HasPrefix(s, "+62")
}, "phone number must start with +62")

var imgUrlValidationRule = validation.NewStringRule(func(s string) bool {
	match, _ := regexp.MatchString(`^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/|\/|\/\/)?[A-z0-9_-]*?[:]?[A-z0-9_-]*?[@]?[A-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/{1}[A-z0-9_\-\:x\=\(\)]+)*(\.(jpg|jpeg|png))?$`, s)
	return match
}, "image url is not valid")

type PostMedicalPatients struct {
	IdentityNumber      int64     `json:"identityNumber"`
	PhoneNumber         string    `json:"phoneNumber"`
	Name                string    `json:"name"`
	Birthdate           time.Time `json:"birthDate"`
	Gender              Gender    `json:"gender"`
	IdentityCardScanImg string    `json:"identityCardScanImg"`
}

func (p PostMedicalPatients) Validate() error {
	idNumber := strconv.Itoa(int(p.IdentityNumber))
	if len(idNumber) != 16 {
		return fmt.Errorf("%s: %s", "identityNumber", "must be 16 characters")
	}

	return validation.ValidateStruct(&p,
		validation.Field(&p.IdentityNumber, validation.Required),
		validation.Field(&p.PhoneNumber, validation.Required, phoneNumberValidationRule, validation.Length(10, 15)),
		validation.Field(&p.Name, validation.Required, validation.Length(3, 30)),
		validation.Field(&p.Birthdate, validation.Required),
		validation.Field(&p.Gender, validation.Required, validation.In(Genders...)),
		validation.Field(&p.IdentityCardScanImg, validation.Required, imgUrlValidationRule),
	)
}
