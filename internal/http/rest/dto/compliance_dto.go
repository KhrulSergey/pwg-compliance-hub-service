package dto

import (
	"database/sql"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/model/enum"
	"time"
)

type ComplianceCheckRequestDto struct {
	PwgEntityGuid    string                `json:"pwgEntityGuid"`
	PwgEntityType    enum.PwgEntityType    `json:"pwgEntityType"`
	ComplianceEntity ComplianceEntity      `json:"complianceEntity"`
	Provider         enum.KYCProviderType  `json:"complianceProvider"`
	CheckRules       []ComplianceCheckRule `json:"checkRules"`
}

type ComplianceCheckResponseDto struct {
	ExternalGuid     string                    `json:"complianceExternalGuid"`
	PwgEntityGuid    string                    `json:"pwgEntityGuid"`
	PwgEntityType    enum.PwgEntityType        `json:"pwgEntityType"`
	ComplianceEntity ComplianceEntity          `json:"complianceEntity"`
	Provider         enum.KYCProviderType      `json:"complianceProvider"`
	CheckRules       []ComplianceCheckRule     `json:"checkRules"`
	Status           enum.ComplianceStatusType `json:"complianceStatus"`
	PassedAt         time.Time                 `json:"passedAt"`
	ExpiredAt        sql.NullTime              `json:"expiredAt"`
	Timestamp        time.Time                 `json:"timestamp"`
}

type ComplianceCheckShortResponseDto struct {
	ExternalGuid string                    `json:"complianceExternalGuid"`
	Status       enum.ComplianceStatusType `json:"complianceStatus"`
}

type ComplianceCheckRule struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Result      string `json:"result"`
	Details     string `json:"details"`
}

// ComplianceEntity represent model that will be sent to kyc/kyb check to provider
type ComplianceEntity struct {
	Institution      *Institution      `json:"institution,omitempty"`
	IndividualPerson *IndividualPerson `json:"individualPerson,omitempty"`
}

type Institution struct {
	InstitutionGuid string             `json:"institutionPwgGuid"`
	LegalName       string             `json:"legalName"`
	CompanyAddress  Address            `json:"companyAddress"`
	ControlPersons  []IndividualPerson `json:"controlPersons"`
}

type IndividualPerson struct {
	UserGuid  string  `json:"userPwgGuid"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Address   Address `json:"address"`
}

type Address struct {
	FullAddress    string `json:"full_address"`
	City           string `json:"city"`
	State          string `json:"state"`
	PostalCode     string `json:"postalCode"`
	ISOCountryCode string `json:"isoCountryCode"`
}

// Deprecated: FunctionName is deprecated.
type ComplianceAssertOutput struct {
	UserID           string `json:"user_id"`
	ComplianceStatus string `json:"compliance_status"`
}

// Deprecated: FunctionName is deprecated.
type ComplianceStatusOutput struct {
	Status string `json:"status"`
}
