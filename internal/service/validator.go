package service

import (
	"fmt"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/model"
	"regexp"
	"time"
)

const (
	StatusPassed     = "KYC_verification_passed"
	StatusFailed     = "KYC_verification_failed"
	StatusInProgress = "KYC_verification_is_in_progress"
)

type ParsedAssertComplianceTime struct {
	ComplianceCreatedAt    time.Time
	ComplianceUpdatedAt    time.Time
	ComplianceExpirationAt time.Time
}

func validateUserID(id string) error {
	match, _ := regexp.MatchString("^[Uu][a-zA-Z\\d]{17}$", id)
	if !match {
		return fmt.Errorf("userId didn't match pattern ''")
	}
	// TBD check if user exist
	return nil
}

func validateAssertCompliance(input *model.ComplianceInput) (*ParsedAssertComplianceTime, error) {
	if input.ServiceProvider == "" {
		return nil, fmt.Errorf("'compliance_service_provider' should not be empty")
	}
	if input.ServiceProviderUserId == "" {
		return nil, fmt.Errorf("'service_provider_user_id' should not be empty")
	}
	complianceCreatedAt, err := time.Parse(time.RFC3339, input.ComplianceCreatedAt)
	if err != nil {
		return nil, err
	}
	complianceUpdatedAt, err := time.Parse(time.RFC3339, input.ComplianceUpdatedAt)
	if err != nil {
		return nil, err
	}
	complianceExpirationAt, err := time.Parse(time.RFC3339, input.ComplianceExpirationAt)
	if err != nil {
		return nil, err
	}

	statusMap := map[string]bool{
		StatusPassed:     true,
		StatusFailed:     true,
		StatusInProgress: true,
	}
	if _, ok := statusMap[input.ComplianceStatus]; !ok {
		return nil, fmt.Errorf("wrong ComplianceStatus")
	}
	return &ParsedAssertComplianceTime{
		ComplianceCreatedAt:    complianceCreatedAt,
		ComplianceUpdatedAt:    complianceUpdatedAt,
		ComplianceExpirationAt: complianceExpirationAt,
	}, nil
}
