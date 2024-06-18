package main

import (
	"cognix.ch/api/v2/api/handler"
	"cognix.ch/api/v2/core/ai"
	"cognix.ch/api/v2/core/bll"
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/oauth"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/security"
	"cognix.ch/api/v2/core/server"
	"cognix.ch/api/v2/core/storage"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/fx"
	"net/http"
)

var Module = fx.Options(
	repository.DatabaseModule,
	repository.RepositoriesModule,
	bll.BLLModule,
	storage.MinioModule,
	messaging.NatsModule,
	ai.EmbeddingModule,
	storage.MilvusModule,
	fx.Provide(ReadConfig,
		NewRouter,
		newGoogleOauthProvider,
		newJWTService,
		//newStorage,
		ai.NewBuilder,
		server.NewAuthMiddleware,
		handler.NewAuthHandler,
		handler.NewCollectorHandler,
		handler.NewSwaggerHandler,
		handler.NewPersonaHandler,
		handler.NewChatHandler,
		handler.NewEmbeddingModelHandler,
		handler.NewTenantHandler,
		handler.NewDocumentHandler,
		newOauthHandler,
	),
	fx.Invoke(
		MountRoute,
		RunServer,
	),
)

func MountRoute(param MountParams) error {
	param.AutHandler.Mount(param.Router, param.AuthMiddleware.RequireAuth)
	param.SwaggerHandler.Mount(param.Router)
	param.ConnectorHandler.Mount(param.Router, param.AuthMiddleware.RequireAuth)
	param.ChatHandler.Mount(param.Router, param.AuthMiddleware.RequireAuth)
	param.PersonaHandler.Mount(param.Router, param.AuthMiddleware.RequireAuth)
	param.EmbeddingModelHandler.Mount(param.Router, param.AuthMiddleware.RequireAuth)
	param.TenantHandler.Mount(param.Router, param.AuthMiddleware.RequireAuth)
	param.DocumentHandler.Mount(param.Router, param.AuthMiddleware.RequireAuth)
	param.OAuthHandler.Mount(param.Router)
	return nil
}

func newGoogleOauthProvider(cfg *Config) oauth.Proxy {
	return oauth.NewGoogleProvider(cfg.OAuth.Google, cfg.RedirectURL)
}
func newJWTService(cfg *Config) security.JWTService {
	return security.NewJWTService(cfg.JWTSecret, cfg.JWTExpiredTime)
}

//	func newStorage(cfg *Config) (storage.Storage, error) {
//		return storage.NewNutsDbStorage(cfg.StoragePath)
//	}
func newOauthHandler(cfg *Config) *handler.OAuthHandler {
	return handler.NewOAuthHandler(cfg.OAuth)
}
func NewRouter() *gin.Engine {
	router := gin.Default()
	router.Use(otelgin.Middleware("service-name"))
	corsConfig := cors.DefaultConfig()

	corsConfig.CustomSchemas = cors.DefaultSchemas
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowWildcard = true
	router.Use(cors.New(corsConfig))
	return router
}

func RunServer(cfg *Config, router *gin.Engine) {
	srv := http.Server{}
	srv.Addr = fmt.Sprintf("0.0.0.0:%d", cfg.Port)
	srv.Handler = router
	otelzap.S().Infof("Start server %s ", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		otelzap.S().Errorf("HTTP server: %s", err.Error())
	}
}
