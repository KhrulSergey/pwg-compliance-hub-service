package model

import (
	"database/sql"
	"github.com/google/uuid"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/http/rest/dto"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/model/enum"
	"time"
)

// ComplianceCheck represent model of processed kyc/kyb personal data check on external provider
type ComplianceCheck struct {
	Id                       uuid.UUID                 `gorm:"id,type:uuid;primary_key"`
	ExternalGuid             string                    `gorm:"external_guid;unique;not null;"`
	ComplianceCheckRequestID uuid.UUID                 `gorm:"foreignKey:ComplianceCheckRequestForeignKey"` //compliance_check_request_id;type:uuid;not null;
	ComplianceCheckRequest   ComplianceCheckRequest    `json:"compliance_check_request"`
	Provider                 enum.KYCProviderType      `gorm:"provider;not null;"`
	Status                   enum.ComplianceStatusType `gorm:"status;not null;"`
	PassedAt                 time.Time                 `gorm:"passed_at"`
	ExpiredAt                sql.NullTime              `gorm:"expired_at;type:TIMESTAMP NULL"`
	CreatedAt                time.Time                 `gorm:"created_at;not null;"`
	UpdatedAt                time.Time                 `gorm:"updated_at;not null;"`
}

// ComplianceCheckRequest represent model of submitted request for kyc/kyb personal data check TO external provider
type ComplianceCheckRequest struct {
	Id                  uuid.UUID                     `gorm:"id,type:uuid;primary_key"`
	PwgEntityGuid       string                        `gorm:"pwg_entity_guid;not null;"`
	PwgEntityType       enum.PwgEntityType            `gorm:"pwg_entity_type;not null;default:'users_account'"`
	RequestExternalGuid string                        `gorm:"request_external_guid;not null;"`
	Provider            enum.KYCProviderType          `gorm:"provider;not null;"`
	CheckRules          []dto.ComplianceCheckRuleDto  `gorm:"type:jsonb;serializer:json"`
	RawRequest          dto.ComplianceCheckRequestDto `gorm:"type:jsonb;serializer:json"`
	Status              enum.ComplianceStatusType     `gorm:"status;not null;"`
	RequestedAt         time.Time                     `gorm:"requested_at;not null;"`
	FinishedAt          sql.NullTime                  `gorm:"finished_at;type:TIMESTAMP NULL"`
	CreatedAt           time.Time                     `gorm:"created_at;not null;"`
	UpdatedAt           time.Time                     `gorm:"updated_at;not null;"`
}
