package service

import (
	"github.com/google/uuid"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/http/rest/dto"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/model"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/model/enum"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/logger"
	"time"
)

// ExternalComplianceService specifies the function required for the auth complianceService
type ExternalComplianceService interface {
	SendComplianceRequestToProvider(complianceRequest *dto.ComplianceCheckRequestDto) (*dto.ComplianceCheckShortResponseDto, error)
	GetComplianceStatusFromProvider(complianceRequest model.ComplianceCheckRequest) (*dto.ComplianceCheckResponseDto, error)
}

// externalComplianceService wrapper for the interface implementation
type externalComplianceService struct {
	logger logger.Logger
}

// InitExternalComplianceService constructs an instance of the Auth ComplianceService
func InitExternalComplianceService(logger logger.Logger) ExternalComplianceService {
	return &externalComplianceService{
		logger: logger,
	}
}

func (ecs externalComplianceService) SendComplianceRequestToProvider(complianceRequest *dto.ComplianceCheckRequestDto) (*dto.ComplianceCheckShortResponseDto, error) {
	//TODO implement me
	return &dto.ComplianceCheckShortResponseDto{
		ExternalGuid: uuid.New().String(), //mock provider response
	}, nil
}
func (ecs externalComplianceService) GetComplianceStatusFromProvider(complianceRequest model.ComplianceCheckRequest) (*dto.ComplianceCheckResponseDto, error) {
	//TODO implement me
	return &dto.ComplianceCheckResponseDto{
		ExternalGuid:     complianceRequest.RequestExternalGuid,
		PwgEntityGuid:    complianceRequest.PwgEntityGuid,
		PwgEntityType:    complianceRequest.PwgEntityType,
		ComplianceEntity: complianceRequest.RawRequest.ComplianceEntity,
		Provider:         complianceRequest.Provider,
		PassedAt:         time.Now(),
		Status:           enum.ComplianceStatusPassed,
	}, nil
}
