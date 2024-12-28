package server

import (
	"context"
	"fmt"
	"hbdtoyou/cmd/hbdtoyou-api-http/config"
	"hbdtoyou/internal/auth"
	authhttphandler "hbdtoyou/internal/auth/handler/http"
	authservice "hbdtoyou/internal/auth/service"
	authpgstore "hbdtoyou/internal/auth/store/postgresql"
	"hbdtoyou/internal/content"
	contenthttphandler "hbdtoyou/internal/content/handler/http"
	contentservice "hbdtoyou/internal/content/service"
	contentpgstore "hbdtoyou/internal/content/store/postgresql"
	"hbdtoyou/internal/payment"
	paymenthttphandler "hbdtoyou/internal/payment/handler/http"
	paymentservice "hbdtoyou/internal/payment/service"
	paymentpgstore "hbdtoyou/internal/payment/store/postgresql"
	"hbdtoyou/internal/template"
	templatehttphandler "hbdtoyou/internal/template/handler/http"
	templateservice "hbdtoyou/internal/template/service"
	templatepgstore "hbdtoyou/internal/template/store/postgresql"
	configlib "hbdtoyou/pkg/config"
	"hbdtoyou/pkg/graceful"
	pglib "hbdtoyou/pkg/postgresql"
	secretlocalfile "hbdtoyou/pkg/secret/client/localfile"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

// Following constants are the possible exit code returned
// when running a server.
const (
	CodeSuccess = iota
	CodeBadConfig
	CodeFailServeHTTP
)

// Option contains available options to run the server.
type Option struct {
	SecretPath string
	ConfigTest bool
}

// Run creates a server with the given Option and starts the
// server.
//
// Run returns a status code suitable for os.Exit() argument.
func Run(opt Option) int {
	s, err := new(opt)
	if err != nil {
		return CodeBadConfig
	}

	// do not start server for config testing
	if opt.ConfigTest {
		return CodeSuccess
	}

	return s.start()
}

// server is the long-runnning application.
type server struct {
	srv      *http.Server
	handlers []handler
	config   config.Config
}

// handler provides mechanism to start HTTP handler. All HTTP
// handlers must implements this interface.
type handler interface {
	Start(multiplexer *mux.Router) error
}

// new creates and returns a new server.
func new(opt Option) (*server, error) {
	s := &server{
		srv: &http.Server{
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}

	// initialize secrets
	var secrets map[string]string
	{
		client := secretlocalfile.New()
		bytes, err := client.Fetch(context.Background(), opt.SecretPath)
		if err != nil {
			log.Printf("[memorify-api-http] failed to fetch secrets: %s\n", err.Error())
			return nil, fmt.Errorf("failed to fetch secrets: %s", err.Error())
		}

		err = yaml.Unmarshal(bytes, &secrets)
		if err != nil {
			log.Printf("[memorify-api-http] failed to unmarshal secrets: %s\n", err.Error())
			return nil, fmt.Errorf("failed to unmarshal secrets: %s", err.Error())
		}
	}

	// initialize config
	{
		configPath, err := getConfigFilePath()
		if err != nil {
			log.Printf("[memorify-api-http] failed to get config filepath: %s\n", err.Error())
			return nil, fmt.Errorf("failed to get config filepath: %s", err.Error())
		}

		var cfg config.Config
		err = configlib.ReadFile(configPath, &cfg, configlib.WithStrictParsing(), configlib.WithSecrets(secrets))
		if err != nil {
			log.Printf("[memorify-api-http] failed to read config: %s\n", err.Error())
			return nil, fmt.Errorf("failed to read config: %s", err.Error())
		}
		s.config = cfg
	}

	// end of config testing
	if opt.ConfigTest {
		return s, nil
	}

	// initilize postgresql client manager
	var pgClientManager *pglib.ClientManager
	{
		clientNames := []string{
			config.PostgreSQLTenant,
		}

		var options []pglib.Option
		for _, clientName := range clientNames {
			cfg, ok := s.config.PostgreSQL[clientName]
			if !ok {
				log.Printf("[memorify-api-http] postgresql config not found for client name: %s\n", clientName)
				return nil, fmt.Errorf("postgresql config not found")
			}

			options = append(options, pglib.WithClientConfig(clientName, pglib.ClientConfig{
				ConnectionString:  cfg.ConnectionString,
				ConnectionTimeout: time.Duration(cfg.ConnectionTimeout),
			}))
		}

		var err error
		pgClientManager, err = pglib.NewClientManager(options...)
		if err != nil {
			log.Printf("[memorify-api-http] failed to initialize postgresql client manager: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize postgresql client manager: %s", err.Error())
		}
	}

	// initialize auth service
	var authSvc auth.Service
	{
		pgDb, err := pgClientManager.GetDatabase(config.PostgreSQLTenant)
		if err != nil {
			log.Printf("[auth-api-http] failed to get postgresql database: %s\n", err.Error())
			return nil, fmt.Errorf("failed to get postgresql database: %s", err.Error())
		}

		pgStore, err := authpgstore.New(pgDb)
		if err != nil {
			log.Printf("[auth-api-http] failed to initialize auth postgresql store: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize auth postgresql store: %s", err.Error())
		}

		svcOptions := []authservice.Option{}
		svcOptions = append(svcOptions, authservice.WithConfig(authservice.Config{
			PasswordSalt:    s.config.User.PasswordSalt,
			TokenExpiration: time.Duration(s.config.User.TokenExpiration),
			TokenSecretKey:  s.config.User.TokenSecretKey,
			ClientID:        s.config.User.ClientID,
		}))

		authSvc, err = authservice.New(pgStore, svcOptions...)
		if err != nil {
			log.Printf("[auth-api-http] failed to initialize auth service: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize auth service: %s", err.Error())
		}
	}

	// initialize template service
	var templateSvc template.Service
	{
		pgDb, err := pgClientManager.GetDatabase(config.PostgreSQLTenant)
		if err != nil {
			log.Printf("[template-api-http] failed to get postgresql database: %s\n", err.Error())
			return nil, fmt.Errorf("failed to get postgresql database: %s", err.Error())
		}

		pgStore, err := templatepgstore.New(pgDb)
		if err != nil {
			log.Printf("[template-api-http] failed to initialize template postgresql store: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize template postgresql store: %s", err.Error())
		}

		templateSvc, err = templateservice.New(pgStore)
		if err != nil {
			log.Printf("[template-api-http] failed to initialize template service: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize template service: %s", err.Error())
		}
	}

	// initialize content service
	var contentSvc content.Service
	{
		pgDb, err := pgClientManager.GetDatabase(config.PostgreSQLTenant)
		if err != nil {
			log.Printf("[content-api-http] failed to get postgresql database: %s\n", err.Error())
			return nil, fmt.Errorf("failed to get postgresql database: %s", err.Error())
		}

		pgStore, err := contentpgstore.New(pgDb)
		if err != nil {
			log.Printf("[content-api-http] failed to initialize content postgresql store: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize content postgresql store: %s", err.Error())
		}

		contentSvc, err = contentservice.New(pgStore, authSvc, templateSvc)
		if err != nil {
			log.Printf("[content-api-http] failed to initialize content service: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize content service: %s", err.Error())
		}
	}

	// initialize payment service
	var paymentSvc payment.Service
	{
		pgDb, err := pgClientManager.GetDatabase(config.PostgreSQLTenant)
		if err != nil {
			log.Printf("[payment-api-http] failed to get postgresql database: %s\n", err.Error())
			return nil, fmt.Errorf("failed to get postgresql database: %s", err.Error())
		}

		pgStore, err := paymentpgstore.New(pgDb)
		if err != nil {
			log.Printf("[payment-api-http] failed to initialize payment postgresql store: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize payment postgresql store: %s", err.Error())
		}

		paymentSvc, err = paymentservice.New(pgStore, authSvc, contentSvc)
		if err != nil {
			log.Printf("[payment-api-http] failed to initialize payment service: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize payment service: %s", err.Error())
		}
	}

	// initialize auth HTTP handler
	{
		var options []authhttphandler.Option
		for scopeName, cfg := range s.config.User.HTTP {
			options = append(options, authhttphandler.WithScopeSetting(scopeName, authhttphandler.ScopeSetting{
				Timeout: time.Duration(cfg.Timeout),
			}))
		}

		identities := []authhttphandler.HandlerIdentity{
			authhttphandler.HandlerLoginSocial,
			authhttphandler.HandlerUser,
		}

		for _, identity := range identities {
			options = append(options, authhttphandler.WithHandler(identity))
		}

		authHTTP, err := authhttphandler.New(authSvc, options...)
		if err != nil {
			log.Printf("[auth-api-http] failed to initialize auth http handlers: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize auth http handlers: %s", err.Error())
		}

		s.handlers = append(s.handlers, authHTTP)
	}

	// initialize content HTTP handler
	{
		var options []contenthttphandler.Option
		for scopeName, cfg := range s.config.Content.HTTP {
			options = append(options, contenthttphandler.WithScopeSetting(scopeName, contenthttphandler.ScopeSetting{
				Timeout: time.Duration(cfg.Timeout),
			}))
		}

		identities := []contenthttphandler.HandlerIdentity{
			contenthttphandler.HandlerContent,
			contenthttphandler.HandlerContents,
		}

		for _, identity := range identities {
			options = append(options, contenthttphandler.WithHandler(identity))
		}

		contentHTTP, err := contenthttphandler.New(contentSvc, authSvc, options...)
		if err != nil {
			log.Printf("[content-api-http] failed to initialize task http handlers: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize task http handlers: %s", err.Error())
		}

		s.handlers = append(s.handlers, contentHTTP)
	}

	// initialize template HTTP handler
	{
		var options []templatehttphandler.Option
		for scopeName, cfg := range s.config.Template.HTTP {
			options = append(options, templatehttphandler.WithScopeSetting(scopeName, templatehttphandler.ScopeSetting{
				Timeout: time.Duration(cfg.Timeout),
			}))
		}

		identities := []templatehttphandler.HandlerIdentity{
			templatehttphandler.HandlerTemplate,
			templatehttphandler.HandlerTemplates,
		}

		for _, identity := range identities {
			options = append(options, templatehttphandler.WithHandler(identity))
		}

		templateHTTP, err := templatehttphandler.New(templateSvc, authSvc, options...)
		if err != nil {
			log.Printf("[template-api-http] failed to initialize task http handlers: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize task http handlers: %s", err.Error())
		}

		s.handlers = append(s.handlers, templateHTTP)
	}

	// initialize payment HTTP handler
	{
		var options []paymenthttphandler.Option
		for scopeName, cfg := range s.config.Payment.HTTP {
			options = append(options, paymenthttphandler.WithScopeSetting(scopeName, paymenthttphandler.ScopeSetting{
				Timeout: time.Duration(cfg.Timeout),
			}))
		}

		identities := []paymenthttphandler.HandlerIdentity{
			paymenthttphandler.HandlerPayment,
			paymenthttphandler.HandlerPayments,
		}

		for _, identity := range identities {
			options = append(options, paymenthttphandler.WithHandler(identity))
		}

		paymentHTTP, err := paymenthttphandler.New(paymentSvc, authSvc, options...)
		if err != nil {
			log.Printf("[payment-api-http] failed to initialize task http handlers: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize task http handlers: %s", err.Error())
		}

		s.handlers = append(s.handlers, paymentHTTP)
	}

	return s, nil
}

// start starts the given server.
func (s *server) start() int {
	log.Println("[memorify-api-http] starting server...")

	// create multiplexer object
	rootMux := mux.NewRouter()
	appMux := rootMux.PathPrefix("/tenant").Subrouter()

	// use middlewares to app mux only
	// appMux.Use(prometheuslib.GetHTTPHandlerMiddleware("memorify-api-http"))

	// starts handlers
	for _, h := range s.handlers {
		if err := h.Start(appMux); err != nil {
			log.Printf("[memorify-api-http] failed to start handler: %s\n", err.Error())
			return CodeFailServeHTTP
		}
	}

	// handle prometheus pull endpoint
	// rootMux.Handle("/metrics", promhttp.Handler())

	// assign multiplexer as server handler
	s.srv.Handler = rootMux

	// serve using graceful mechanism
	address := fmt.Sprintf(":%d", s.config.Server.Port)
	err := graceful.ServeHTTP(s.srv, address, 0)
	if err != nil {
		log.Printf("[memorify-api-http] failed to start server: %s\n", err.Error())
		return CodeFailServeHTTP
	}

	return CodeSuccess
}
