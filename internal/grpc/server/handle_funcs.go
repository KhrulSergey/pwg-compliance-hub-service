package grpc

import (
	"context"
	"fmt"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/http/rest/dto"

	grpc "gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/pkg/grpc/compliance_hub_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *grpcServer) CheckCompliance(ctx context.Context, req *grpc.ComplianceCheckRequestDto) (*grpc.ComplianceCheckShortResponseDto, error) {
	checkRequestDto := convertCheckRequestRpcToDto(req)
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

func convertCheckShortResponseDtoToRpc(resp dto.ComplianceCheckShortResponseDto) *grpc.ComplianceCheckShortResponseDto {

	return &grpc.ComplianceCheckShortResponseDto{
		ComplianceExternalGuid: resp.ExternalGuid,
		ComplianceStatus:       grpc.ComplianceStatusType_ACTION_REQUIRED, //todo convert all fields
		//PwgEntityType    enum.PwgEntityType    `json:"pwgEntityType"`
		//ComplianceEntity ComplianceEntity      `json:"complianceEntity"`
		//Provider         enum.KYCProviderType  `json:"complianceProvider"`
		//CheckRules       []ComplianceCheckRule `json:"checkRules"`
	}
}

func convertCheckRequestRpcToDto(req *grpc.ComplianceCheckRequestDto) *dto.ComplianceCheckRequestDto {

	return &dto.ComplianceCheckRequestDto{
		PwgEntityGuid: req.PwgEntityGuid, //todo convert all fields
		//PwgEntityType    enum.PwgEntityType    `json:"pwgEntityType"`
		//ComplianceEntity ComplianceEntity      `json:"complianceEntity"`
		//Provider         enum.KYCProviderType  `json:"complianceProvider"`
		//CheckRules       []ComplianceCheckRule `json:"checkRules"`
	}
}
