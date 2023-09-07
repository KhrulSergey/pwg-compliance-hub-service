package mapper

import (
	"database/sql"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/http/rest/dto"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/model"
)

func ToComplianceRequestEntity(requestDto dto.ComplianceCheckRequestDto) model.ComplianceCheckRequest {
	return model.ComplianceCheckRequest{
		PwgEntityGuid: requestDto.PwgEntityGuid,
		PwgEntityType: requestDto.PwgEntityType,
		Provider:      requestDto.Provider,
		CheckRules:    requestDto.CheckRules,
		RawRequest:    requestDto,
	}
}
func ToComplianceResponseDtoFromRequestEntity(checkRequestEntry model.ComplianceCheckRequest) dto.ComplianceCheckShortResponseDto {
	return dto.ComplianceCheckShortResponseDto{
		ExternalGuid: checkRequestEntry.RequestExternalGuid,
		Status:       checkRequestEntry.Status,
	}
}

func ApplyChangesToComplianceRequestEntity(checkRequestEntry *model.ComplianceCheckRequest, responseDto *dto.ComplianceCheckResponseDto) {
	checkRequestEntry.Status = responseDto.Status
	checkRequestEntry.UpdatedAt = responseDto.Timestamp
	if !responseDto.PassedAt.IsZero() {
		checkRequestEntry.FinishedAt = sql.NullTime{
			Time:  responseDto.PassedAt,
			Valid: true,
		}
	}
}

func ToComplianceCheckEntity(responseDto dto.ComplianceCheckResponseDto, request model.ComplianceCheckRequest) model.ComplianceCheck {
	return model.ComplianceCheck{
		ExternalGuid:             responseDto.ExternalGuid,
		ComplianceCheckRequest:   request,
		ComplianceCheckRequestID: request.Id,
		Provider:                 responseDto.Provider,
		Status:                   responseDto.Status,
		PassedAt:                 responseDto.PassedAt,
		ExpiredAt:                responseDto.ExpiredAt,
	}
}

func ApplyChangesToComplianceCheckEntity(existedCheck *model.ComplianceCheck, newCheck *model.ComplianceCheck) {
	existedCheck.Status = newCheck.Status
	existedCheck.PassedAt = newCheck.PassedAt
	existedCheck.ExpiredAt = newCheck.ExpiredAt
}
