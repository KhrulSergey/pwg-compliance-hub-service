package grpc

import (
	"context"
	"errors"
	"fmt"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/http/rest/dto"
	grpc "gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/pkg/grpc/compliance_hub_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrRequestMissing            = errors.New("request is not provided")
	ErrComplianceRequestNotFound = errors.New("compliance request not found")
)

func (s *grpcServer) CheckCompliance(ctx context.Context, request *grpc.ComplianceCheckRequestRpc) (*grpc.ComplianceCheckShortResponseRpc, error) {
	//validate request
	if request.PwgEntityGuid == "" || request.GetPwgEntityType() == grpc.PwgEntityType_NULL_ENTITY ||
		request.ComplianceProvider == grpc.KYCProviderType_NULL_PROVIDER || request.ComplianceEntity == nil {
		return nil, ErrRequestMissing
	}

	checkRequestDto := convertCheckRequestRpcToDto(request)
	shortResponse, err := s.complianceCheckService.CreateRequestAndCheckCompliance(checkRequestDto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Can't get compliance status. Error: %v", err.Error()))
	}
	return convertCheckShortResponseDtoToRpc(shortResponse), nil
}

//todo use for convert enum or delete
//func getEnumValueOptions(enumName string, valueName string) string {
//	enumDesc, err := protoregistry.GlobalTypes.FindEnumByName(protoreflect.FullName(enumName))
//	if err != nil {
//		fmt.Printf("Failed to find enum %s: %v\n", enumName, err)
//		return ""
//	}
//	valueDesc := enumDesc.Values().ByName(protoreflect.Name(valueName))
//	if valueDesc == nil {
//		fmt.Printf("Failed to find value %s in enum %s\n", valueName, enumName)
//		return ""
//	}
//	options := proto.GetExtension(valueDesc.Options().Interface(), protoiface.E_EnumValueOptions).(*protoiface.EnumValueOptions)
//	if options != nil {
//		return options.GetDisplayValue()
//	}
//	return ""
//}

func convertCheckShortResponseDtoToRpc(resp dto.ComplianceCheckShortResponseDto) *grpc.ComplianceCheckShortResponseRpc {

	return &grpc.ComplianceCheckShortResponseRpc{
		ComplianceExternalGuid: resp.ExternalGuid,
		ComplianceStatus:       grpc.ComplianceStatusType_ACTION_REQUIRED, //todo use enum converter
	}
}

func convertCheckRequestRpcToDto(req *grpc.ComplianceCheckRequestRpc) *dto.ComplianceCheckRequestDto {

	complianceEntity := *convertComplianceEntityRpcToDto(req.ComplianceEntity)
	checkRules := make([]dto.ComplianceCheckRuleDto, 0, len(req.GetCheckRules()))
	for i := range req.GetCheckRules() {
		checkRules = append(checkRules, *convertCheckRulesRpcToDto(req.GetCheckRules()[i]))
	}
	return &dto.ComplianceCheckRequestDto{
		PwgEntityGuid: req.PwgEntityGuid,
		//PwgEntityType: req.GetPwgEntityType(), todo use enum converter
		ComplianceEntity: complianceEntity,
		CheckRules:       checkRules,
		//Provider:  req.GetComplianceProvider(), todo use enum converter
	}
}

func convertCheckRulesRpcToDto(req *grpc.ComplianceCheckRuleRpc) *dto.ComplianceCheckRuleDto {
	if req == nil {
		return nil
	}
	return &dto.ComplianceCheckRuleDto{
		Description: req.Description,
		Result:      req.Result,
		Details:     req.Details,
		Name:        req.Name,
	}
}

func convertComplianceEntityRpcToDto(req *grpc.ComplianceEntityRpc) *dto.ComplianceEntityDto {
	individualPerson := convertIndividualPersonRpcToDto(req.GetIndividualPerson())
	institution := convertInstitutionRpcToDto(req.GetInstitution())
	return &dto.ComplianceEntityDto{
		IndividualPerson: individualPerson,
		Institution:      institution,
	}
}

func convertInstitutionRpcToDto(req *grpc.InstitutionRpc) *dto.InstitutionDto {
	if req == nil {
		return nil
	}
	controlPersonsDto := make([]dto.IndividualPersonDto, 0, len(req.GetControlPersons()))
	for i := range req.GetControlPersons() {
		controlPersonsDto = append(controlPersonsDto, *convertIndividualPersonRpcToDto(req.GetControlPersons()[i]))
	}
	return &dto.InstitutionDto{
		CompanyAddress:  *convertAddressRpcToDto(req.GetCompanyAddress()),
		InstitutionGuid: req.InstitutionPwgGuid,
		LegalName:       req.LegalName,
		ControlPersons:  controlPersonsDto,
	}
}

func convertIndividualPersonRpcToDto(req *grpc.IndividualPersonRpc) *dto.IndividualPersonDto {
	if req == nil {
		return nil
	}
	return &dto.IndividualPersonDto{
		Address:   *convertAddressRpcToDto(req.GetAddress()),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		UserGuid:  req.UserPwgGuid,
	}
}

func convertAddressRpcToDto(req *grpc.AddressRpc) *dto.AddressDto {
	if req == nil {
		return nil
	}
	return &dto.AddressDto{
		FullAddress:    req.FullAddress,
		City:           req.City,
		ISOCountryCode: req.IsoCountryCode,
		PostalCode:     req.PostalCode,
		State:          req.State,
	}
}
