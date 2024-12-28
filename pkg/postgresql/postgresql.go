package postgresql

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

// Followings are the known errors from go-lib/postgresql.
var (
	errClientNotInitialized          = errors.New("postgresql: client not initialized")
	errInvalidClientConnectionString = errors.New("postgresql: invalid client connection string")
)

// client is the PostgreSQL client.
type client struct {
	db *sqlx.DB
}

// ClientConfig contains the available configuration for a
// client.
type ClientConfig struct {
	// ConnectionString usually contains database name and
	// it's connection information for connecting to the
	// database.
	ConnectionString string

	// ConnectionTimeout is the time limit when trying to
	// connect to a database.
	ConnectionTimeout time.Duration
}

// ClientManager is the clients manager which creates and
// stores PostgreSQL client and it's configuration.
type ClientManager struct {
	clients       map[string]client
	clientConfigs map[string]ClientConfig
}

// NewClientManager creates a new ClientManager.
func NewClientManager(options ...Option) (*ClientManager, error) {
	cm := &ClientManager{
		clients:       make(map[string]client),
		clientConfigs: make(map[string]ClientConfig),
	}

	// apply options
	for _, opt := range options {
		if err := opt(cm); err != nil {
			return nil, err
		}
	}

	// create all clients
	err := cm.newClients()
	if err != nil {
		return nil, err
	}

	return cm, nil
}

// newClients creates clients based on the stored client
// configs The created clients are stored in ClientManager.
func (cm *ClientManager) newClients() error {
	// create client for each client configs
	for name, cfg := range cm.clientConfigs {
		cli, err := newClient(cfg)
		if err != nil {
			return err
		}
		cm.clients[name] = cli
	}

	return nil
}

// GetDatabase returns a PostgreSQL database based on the
// given client name.
func (cm *ClientManager) GetDatabase(clientName string) (*sqlx.DB, error) {
	cli, ok := cm.clients[clientName]
	if !ok {
		return nil, errClientNotInitialized
	}
	return cli.db, nil
}

// newClient creates client with the given client config.
func newClient(cfg ClientConfig) (client, error) {
	var cli client

	// validate connection string
	connStr := cfg.ConnectionString
	if connStr == "" {
		return cli, errInvalidClientConnectionString
	}

	// add connection timeout to connection string
	if !strings.Contains(strings.ToLower(connStr), "connect_timeout") {
		timeoutSec := int(cfg.ConnectionTimeout.Seconds())
		connStr = fmt.Sprintf("%s connect_timeout=%d", connStr, timeoutSec)
	}

	// connect to dabatabase
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return cli, err
	}

	cli = client{
		db: db,
	}
	return cli, nil
}
