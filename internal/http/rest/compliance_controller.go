package rest

import (
	"encoding/json"
	"fmt"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/http/rest/dto"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/service"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/logger"

	"github.com/valyala/fasthttp"
)

// ComplianceController is responsible for processing request and applying the business logic
type ComplianceController struct {
	logger            logger.Logger
	complianceService service.ComplianceService
}

// Deprecated
// InitComplianceController constructs a controller
func InitComplianceController(logger logger.Logger, complianceService service.ComplianceService) *ComplianceController {
	return &ComplianceController{
		logger:            logger,
		complianceService: complianceService,
	}
}

// Deprecated
// CheckCompliance godoc
//
//	@Summary		Send request for compliance check
//	@Description	Send request for compliance check
//	@Tags			compliance
//	@Accept			json
//	@Param			RequestBody	body	dto.ComplianceCheckRequestDto	true	"object to update"
//	@Produce		json
//	@Failure		400	{object}	service.ErrorOutput
//	@Failure		500	{object}	service.ErrorOutput
//	@Router			/checkCompliance [post]
//	@Security		JWT
func (l *ComplianceController) CheckCompliance(ctx *fasthttp.RequestCtx) {
	var inputRequestDto dto.ComplianceCheckRequestDto
	err := json.Unmarshal(ctx.PostBody(), &inputRequestDto)
	if err != nil {
		errBody := service.ErrorOutput{
			Code:    fasthttp.StatusBadRequest,
			Message: fmt.Sprintf("Bad Request [%v]", err),
		}
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		if err = json.NewEncoder(ctx).Encode(errBody); err != nil {
			ctx.SetBody([]byte(err.Error()))
			return
		}
		return
	}

	resp, err := l.complianceService.CreateRequestAndCheckCompliance(&inputRequestDto)
	if err != nil {

	}
	if err = json.NewEncoder(ctx).Encode(resp); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody([]byte(err.Error()))
		return
	}
	ctx.SetStatusCode(fasthttp.StatusCreated)
}
