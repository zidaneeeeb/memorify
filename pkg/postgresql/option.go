package postgresql

import sqlx "github.com/jmoiron/sqlx"

// Option controls the behavior of ClientManager.
type Option func(cm *ClientManager) error

// WithClientConfig adds a new client configuration to
// ClientManager. The given name will be used as client name.
func WithClientConfig(name string, cfg ClientConfig) Option {
	return func(cm *ClientManager) error {
		cm.clientConfigs[name] = cfg
		return nil
	}
}

// WithClient adds a new client to ClientManager. The given
// name will be used as client name.
//
// Only for testing purpose. For real case, please use
// WithClientConfig() instead.
func WithClient(name string, db *sqlx.DB) Option {
	return func(cm *ClientManager) error {
		cm.clients[name] = client{
			db: db,
		}
		return nil
	}
}
