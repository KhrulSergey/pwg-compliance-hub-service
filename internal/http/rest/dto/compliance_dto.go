package dto

import (
	"database/sql"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/model/enum"
	"time"
)

type ComplianceCheckRequestDto struct {
	PwgEntityGuid    string                   `json:"pwgEntityGuid"`
	PwgEntityType    enum.PwgEntityType       `json:"pwgEntityType"`
	ComplianceEntity ComplianceEntityDto      `json:"complianceEntity"`
	Provider         enum.KYCProviderType     `json:"complianceProvider"`
	CheckRules       []ComplianceCheckRuleDto `json:"checkRules"`
}

type ComplianceCheckResponseDto struct {
	ExternalGuid     string                    `json:"complianceExternalGuid"`
	PwgEntityGuid    string                    `json:"pwgEntityGuid"`
	PwgEntityType    enum.PwgEntityType        `json:"pwgEntityType"`
	ComplianceEntity ComplianceEntityDto       `json:"complianceEntity"`
	Provider         enum.KYCProviderType      `json:"complianceProvider"`
	CheckRules       []ComplianceCheckRuleDto  `json:"checkRules"`
	Status           enum.ComplianceStatusType `json:"complianceStatus"`
	PassedAt         time.Time                 `json:"passedAt"`
	ExpiredAt        sql.NullTime              `json:"expiredAt"`
	Timestamp        time.Time                 `json:"timestamp"`
}

type ComplianceCheckShortResponseDto struct {
	ExternalGuid string                    `json:"complianceExternalGuid"`
	Status       enum.ComplianceStatusType `json:"complianceStatus"`
}

type ComplianceCheckRuleDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Result      string `json:"result"`
	Details     string `json:"details"`
}

// ComplianceEntityDto represent object and its personal data that will be sent to kyc/kyb check to provider
type ComplianceEntityDto struct {
	Institution      *InstitutionDto      `json:"institution,omitempty"`
	IndividualPerson *IndividualPersonDto `json:"individualPerson,omitempty"`
}

type InstitutionDto struct {
	InstitutionGuid string                `json:"institutionPwgGuid"`
	LegalName       string                `json:"legalName"`
	CompanyAddress  AddressDto            `json:"companyAddress"`
	ControlPersons  []IndividualPersonDto `json:"controlPersons"`
}

type IndividualPersonDto struct {
	UserGuid  string     `json:"userPwgGuid"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Address   AddressDto `json:"address"`
}

type AddressDto struct {
	FullAddress    string `json:"full_address"`
	City           string `json:"city"`
	State          string `json:"state"`
	PostalCode     string `json:"postalCode"`
	ISOCountryCode string `json:"isoCountryCode"`
}
