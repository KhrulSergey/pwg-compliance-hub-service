package router

import (
	"github.com/fasthttp/router"
	fastHttpRouter "github.com/fasthttp/router"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/config"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/docs"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/http/middleware"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/http/rest"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/logger"
)

// AppRouterHandler is responsible for routing components
type AppRouterHandler struct {
	complianceController *rest.ComplianceController
	logger               logger.Logger
	appConfig            config.AppConfig
}

//	@title			Swagger Entitlements-engine API
//	@version		1.0
//	@description	This is a microservice serving permission mission
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//
// Deprecated
// NewRouterHandler constructs a router
func NewRouterHandler(logger logger.Logger, appConfig config.AppConfig, complianceController *rest.ComplianceController) *AppRouterHandler {
	return &AppRouterHandler{
		complianceController: complianceController,
		logger:               logger,
		appConfig:            appConfig,
	}
}

// Deprecated
// InitRouter registers routing patterns for the auth service in the global router
func (h *AppRouterHandler) InitRouter() *router.Router {

	r := fastHttpRouter.New()
	r.GlobalOPTIONS = middleware.MiddlewareSetupResponse(func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusOK)
	})

	// Define router with Swagger
	docs.SwaggerInfo.Version = h.appConfig.AppVersion
	docs.SwaggerInfo.Host = h.appConfig.ServiceExternalUrl
	docs.SwaggerInfo.Schemes = []string{h.appConfig.ServiceProtocol}
	docs.SwaggerInfo.BasePath = ""

	r.GET("/docs/{filepath:*}", fasthttpadaptor.NewFastHTTPHandlerFunc(httpSwagger.WrapHandler))

	// Define api group
	apiGroup := r.Group("")

	apiGroup.POST("/checkCompliance", middleware.MiddlewareSetupResponse(h.complianceController.CheckCompliance))

	return r
}
