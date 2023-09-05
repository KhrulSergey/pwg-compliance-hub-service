package model

import (
	"database/sql"
	"github.com/google/uuid"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/http/rest/dto"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/model/enum"
	"gorm.io/gorm"
	"time"
)

// KYC CHECK    Example
// guid    fee681e7-0fed-408f-bbb2-3154fbc5221e
// kyc_external_guid    140634866
// kyc_check_request_id    56
// kyc_provider    FinClusive
// kyc_status    pending/passed
// passed_at    12.02.2023 19-42
// expired_at    null

// ComplianceCheck represent model of processed kyc/kyb check on external provider
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

// ComplianceCheckRequest represent model of submitted request for kyc/kyb check TO external provider
type ComplianceCheckRequest struct {
	Id                  uuid.UUID                     `gorm:"id,type:uuid;primary_key"`
	PwgEntityGuid       string                        `gorm:"pwg_entity_guid;not null;"`
	PwgEntityType       enum.PwgEntityType            `gorm:"pwg_entity_type;not null;default:'users_account'"`
	RequestExternalGuid string                        `gorm:"request_external_guid;not null;"`
	Provider            enum.KYCProviderType          `gorm:"provider;not null;"`
	CheckRules          []dto.ComplianceCheckRule     `gorm:"type:jsonb;serializer:json"`
	RawRequest          dto.ComplianceCheckRequestDto `gorm:"type:jsonb;serializer:json"`
	Status              enum.ComplianceStatusType     `gorm:"status;not null;"`
	RequestedAt         time.Time                     `gorm:"requested_at;not null;"`
	FinishedAt          sql.NullTime                  `gorm:"finished_at;type:TIMESTAMP NULL"`
	CreatedAt           time.Time                     `gorm:"created_at;not null;"`
	UpdatedAt           time.Time                     `gorm:"updated_at;not null;"`
	//RawRequest               pq.StringArray       `gorm:"type:text[]"`
}

// todo delete!
type Compliance struct {
	gorm.Model
	UserID                 string    `gorm:"user_id;unique;not null;" json:"transuser_id"`
	ServiceProvider        string    `gorm:"compliance_service_provider;not null;" json:"compliance_service_provider"`
	ServiceProviderUserId  string    `gorm:"service_provider_user_id;not null;" json:"service_provider_user_id"`
	ComplianceCreatedAt    time.Time `gorm:"compliance_created_at;not null;" json:"compliance_created_at"`
	ComplianceUpdatedAt    time.Time `gorm:"compliance_updated_at;not null;" json:"compliance_updated_at"`
	ComplianceExpirationAt time.Time `gorm:"compliance_expiration_at;not null;" json:"compliance_expiration_at"`
	ComplianceStatus       string    `gorm:"compliance_status;not null;" json:"compliance_status"`
}

type ComplianceInput struct {
	ServiceProvider        string `json:"compliance_service_provider"`
	ServiceProviderUserId  string `json:"service_provider_user_id"`
	ComplianceCreatedAt    string `json:"compliance_created_at"`
	ComplianceUpdatedAt    string `json:"compliance_updated_at"`
	ComplianceExpirationAt string `json:"compliance_expiration_at"`
	ComplianceStatus       string `json:"compliance_status"`
}

type ComplianceStatus struct {
	Status string `gorm:"compliance_status"`
}
