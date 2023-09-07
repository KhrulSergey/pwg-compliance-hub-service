package service

import (
	"errors"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/http/rest/dto"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/logger"
)

// AccountOperatorServiceClient specifies the function required for the auth complianceService
type AccountOperatorServiceClient interface {
	SendComplianceResponse(complianceResponse dto.ComplianceCheckResponseDto) error
}

// accountOperatorServiceClient wrapper for the interface implementation
type accountOperatorServiceClient struct {
	logger logger.Logger
}

// InitAccountOperatorService constructs an instance of the Auth ComplianceService
func InitAccountOperatorServiceClient(logger logger.Logger) AccountOperatorServiceClient {
	return &accountOperatorServiceClient{
		logger: logger,
	}
}

func (aos accountOperatorServiceClient) SendComplianceResponse(complianceResponse dto.ComplianceCheckResponseDto) error {
	aos.logger.Debugf("Try to send compliance check response to account operator service, data: %s", complianceResponse)
	if complianceResponse.PwgEntityGuid == "" || complianceResponse.PwgEntityType == "" {
		return errors.New("response data is not valid. Check params")
	}
	return nil
}
