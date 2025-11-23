// Package server
package server

import (
	"chart-provider/src/interfaces/cleaner"
	"chart-provider/src/interfaces/content"
	"chart-provider/src/interfaces/global"
	"chart-provider/src/interfaces/server/dto"
	controllerImpl "chart-provider/src/server/controller"
	mid "chart-provider/src/server/middleware"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"time"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	slogecho "github.com/samber/slog-echo"
)

func StartServer(content *content.ApplicationContent) {
	config := content.ConfigManager().GetConfig()
	logger := content.Logger()

	logger.Info("Http server starting...")

	e := echo.New()
	e.Logger.SetOutput(io.Discard)
	e.Logger.SetLevel(log.OFF)

	httpConfig := config.ServerConfig.HttpServerConfig

	switch httpConfig.ProxyType {
	case 0:
		e.IPExtractor = echo.ExtractIPDirect()
	case 1:
		trustOperations := make([]echo.TrustOption, 0, len(httpConfig.TrustIps))
		for _, ip := range httpConfig.TrustIps {
			_, network, err := net.ParseCIDR(ip)
			if err != nil {
				logger.Warnf("%s is not a valid CIDR string, skipping it", ip)
				continue
			}
			trustOperations = append(trustOperations, echo.TrustIPRange(network))
		}
		e.IPExtractor = echo.ExtractIPFromXFFHeader(trustOperations...)
	case 2:
		e.IPExtractor = echo.ExtractIPFromRealIPHeader()
	default:
		logger.Warnf("Invalid proxy type %d, using default (direct)", httpConfig.ProxyType)
		e.IPExtractor = echo.ExtractIPDirect()
	}

	if httpConfig.SSLConfig.Enable && httpConfig.SSLConfig.ForceHttps {
		e.Use(middleware.HTTPSRedirect())
	}

	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: *global.RequestTimeout,
	}))

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: func(ctx echo.Context, err error, stack []byte) error {
			logger.Errorf("Recovered from a fatal error: %v, stack: %s", err, string(stack))
			return err
		},
	}))

	e.Use(slogecho.NewWithConfig(logger.LogHandler(), slogecho.Config{
		DefaultLevel:     slog.LevelInfo,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,
	}))

	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            httpConfig.SSLConfig.HSTSConfig.MaxAge,
		HSTSExcludeSubdomains: !httpConfig.SSLConfig.HSTSConfig.IncludeSubdomains,
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.OPTIONS},
	}))

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: *global.GzipLevel,
	}))

	if httpConfig.BodyLimit != "" {
		e.Use(middleware.BodyLimit(httpConfig.BodyLimit))
	} else {
		logger.Warn("No body limit set, be aware of possible DDOS attacks")
	}

	if httpConfig.RateLimit > 0 {
		ipPathLimiter := mid.NewSlidingWindowLimiter(time.Minute, httpConfig.RateLimit)
		ipPathLimiter.StartCleanup(2 * time.Minute)
		e.Use(mid.RateLimitMiddleware(ipPathLimiter, mid.CombinedKeyFunc))
		logger.Infof("Rate limit: %dQPS", httpConfig.RateLimit)
	} else {
		logger.Warn("No rate limit was set, be aware of possible DDOS attacks")
	}

	if httpConfig.JWTSecret != "" {
		jwtMiddleware := echojwt.WithConfig(echojwt.Config{
			SigningKey:    []byte(httpConfig.JWTSecret),
			TokenLookup:   "header:Authorization:Bearer ",
			SigningMethod: *global.SigningMethod,
			ErrorHandler: func(c echo.Context, err error) error {
				var data *dto.ApiResponse[any]
				switch {
				case errors.Is(err, echojwt.ErrJWTMissing):
					data = dto.NewApiResponse[any](dto.ErrMissingOrMalformedJwt, nil)
				case errors.Is(err, echojwt.ErrJWTInvalid):
					data = dto.NewApiResponse[any](dto.ErrInvalidOrExpiredJwt, nil)
				default:
					data = dto.NewApiResponse[any](dto.ErrUnknownJwtError, nil)
				}
				return data.Response(c)
			},
		})

		e.Use(jwtMiddleware)
	}

	logger.Info("Service initializing...")

	logger.Info("Controller initializing...")

	chartController := controllerImpl.NewChartController(logger)

	logger.Info("Applying router...")

	e.GET("/api/v1/charts/*", chartController.HandleProxy)

	e.Any("*", func(c echo.Context) error {
		return dto.ErrorResponse(c, dto.ErrNoMatchRoute)
	})

	content.Cleaner().Add("Http Server", func(e *echo.Echo) cleaner.ShutdownCallback {
		return func(ctx context.Context) error {
			return e.Shutdown(ctx)
		}
	}(e))

	protocol := "http"
	if httpConfig.SSLConfig.Enable {
		protocol = "https"
	}
	address := fmt.Sprintf("%s:%d", httpConfig.Host, httpConfig.Port)
	logger.Infof("Server started at %s://%s", protocol, address)

	var err error
	if httpConfig.SSLConfig.Enable {
		err = e.StartTLS(
			address,
			httpConfig.SSLConfig.Cert,
			httpConfig.SSLConfig.Key,
		)
	} else {
		err = e.Start(address)
	}

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatalf("Http server error: %v", err)
	}
}
