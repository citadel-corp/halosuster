package user

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var itNIPValidationRule = validation.NewStringRule(func(s string) bool {
	if s[0:3] != "615" {
		return false
	}
	if !strings.Contains("12", s[3:4]) {
		return false
	}
	currentYear := time.Now().Year()
	nipYear, _ := strconv.Atoi(s[4:8])
	if nipYear < 2000 || nipYear > currentYear {
		return false
	}
	nipMonth, _ := strconv.Atoi(s[8:10])
	if nipMonth < 1 || nipMonth > 12 {
		return false
	}
	return true
}, "NIP must be valid")

var nurseNIPValidationRule = validation.NewStringRule(func(s string) bool {
	if s[0:3] != "303" {
		return false
	}
	if !strings.Contains("12", s[3:4]) {
		return false
	}
	currentYear := time.Now().Year()
	nipYear, _ := strconv.Atoi(s[4:8])
	if nipYear < 2000 || nipYear > currentYear {
		return false
	}
	nipMonth, _ := strconv.Atoi(s[8:10])
	if nipMonth < 1 || nipMonth > 12 {
		return false
	}
	return true
}, "NIP must be valid")

var imgUrlValidationRule = validation.NewStringRule(func(s string) bool {
	match, _ := regexp.MatchString(`^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/|\/|\/\/)?[A-z0-9_-]*?[:]?[A-z0-9_-]*?[@]?[A-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/{1}[A-z0-9_\-\:x\=\(\)]+)*(\.(jpg|jpeg|png))?$`, s)
	return match
}, "image url is not valid")

type CreateITUserPayload struct {
	NIP      int    `json:"nip"`
	Name     string `json:"name"`
	Password string `json:"password"`

	nipStr string
}

func (p CreateITUserPayload) Validate() error {
	p.nipStr = strconv.Itoa(p.NIP)
	return validation.ValidateStruct(&p,
		validation.Field(&p.nipStr, validation.Required, validation.Length(13, 15), itNIPValidationRule),
		validation.Field(&p.Name, validation.Required, validation.Length(5, 50)),
		validation.Field(&p.Password, validation.Required, validation.Length(5, 33)),
	)
}

type CreateNurseUserPayload struct {
	NIP                 int    `json:"nip"`
	Name                string `json:"name"`
	IdentityCardScanImg string `json:"identityCardScanImg"`

	nipStr string
}

func (p CreateNurseUserPayload) Validate() error {
	p.nipStr = strconv.Itoa(p.NIP)
	return validation.ValidateStruct(&p,
		validation.Field(&p.nipStr, validation.Required, validation.Length(13, 15), nurseNIPValidationRule),
		validation.Field(&p.Name, validation.Required, validation.Length(5, 50)),
		validation.Field(&p.IdentityCardScanImg, validation.Required, imgUrlValidationRule),
	)
}

type ITUserLoginPayload struct {
	NIP      int    `json:"nip"`
	Password string `json:"password"`

	nipStr string
}

func (p ITUserLoginPayload) Validate() error {
	p.nipStr = strconv.Itoa(p.NIP)

	return validation.ValidateStruct(&p,
		validation.Field(&p.nipStr, validation.Required, validation.Length(13, 15), itNIPValidationRule),
		validation.Field(&p.Password, validation.Required, validation.Length(5, 33)),
	)
}

type NurseUserLoginPayload struct {
	NIP      int    `json:"nip"`
	Password string `json:"password"`

	nipStr string
}

func (p NurseUserLoginPayload) Validate() error {
	p.nipStr = strconv.Itoa(p.NIP)

	return validation.ValidateStruct(&p,
		validation.Field(&p.nipStr, validation.Required, validation.Length(13, 15), nurseNIPValidationRule),
		validation.Field(&p.Password, validation.Required, validation.Length(5, 33)),
	)
}

type ListUserPayload struct {
	UserID    string `schema:"userId" binding:"omitempty"`
	Limit     int    `schema:"limit" binding:"omitempty"`
	Offset    int    `schema:"offset" binding:"omitempty"`
	Name      string `schema:"name" binding:"omitempty"`
	NIP       int    `schema:"nip" binding:"omitempty"`
	Role      string `schema:"role" binding:"omitempty"`
	CreatedAt string `schema:"createdAt" binding:"omitempty"`

	nipStr        string
	CreatedAtType CreatedAtType
	RoleType      RoleType
}

type CreatedAtType int
type RoleType int

const (
	Ascending CreatedAtType = iota
	Descending
	IgnoreCreatedAt

	ITType RoleType = iota
	NurseType
	IgnoreRole
)

type UpdateNursePayload struct {
	NIP  int    `json:"nip"`
	Name string `json:"name"`

	nipStr string
}

func (p UpdateNursePayload) Validate() error {
	p.nipStr = strconv.Itoa(p.NIP)

	return validation.ValidateStruct(&p,
		validation.Field(&p.nipStr, validation.Required, validation.Length(13, 15), nurseNIPValidationRule),
		validation.Field(&p.Name, validation.Required, validation.Length(5, 50)),
	)
}

type GrantNurseAccessPayload struct {
	Password string `json:"password"`
}

func (p GrantNurseAccessPayload) Validate() error {

	return validation.ValidateStruct(&p,
		validation.Field(&p.Password, validation.Required, validation.Length(5, 33)),
	)
}
