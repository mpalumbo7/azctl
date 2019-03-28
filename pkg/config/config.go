package config

import (
	"fmt"

	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/spf13/viper"
)

// Config contains the data from the configuration file
type Config struct {
	CurrentContext string     `mapstructure:"current-context"`
	Contexts       []*Context `mapstructure:"contexts"`
	Users          []*User    `mapstructure:"users"`
}

// SaveToken persists a token for a given context
func (cxt *Context) SaveToken(t adal.Token) error {
	fmt.Printf("Saving context %s for resource: %s\n", cxt.Name, t.Resource)
	return nil
}

// User contains the data for a given identity
type User struct {
	Name            string `mapstructure:"name"`
	TenantID        string `mapstructure:"tenantId"`
	CertificateData string `mapstructure:"certificate-data,omitempty"`
	ClientID        string `mapstructure:"clientId,omitempty"`
	Endpoint        string `mapstructure:"endpoint,omitempty"`
}

// Context contains a complete Azure configuration spec
type Context struct {
	Name            string       `mapstructure:"name"`
	TenantID        string       `mapstructure:"tenantId"`
	Endpoint        string       `mapstructure:"endpoint,omitempty"`
	ClientID        string       `mapstructure:"clientId,omitempty"`
	CertificateData string       `mapstructure:"certificate-data,omitempty"`
	SubscriptionID  string       `mapstructure:"subscriptionId"`
	Resource        string       `mapstructure:"resource"`
	Tokens          []adal.Token `mapstructure:"tokens"`
}

// GetConfig reads the entire config file
func GetConfig() (*Config, error) {
	// Read the configuration file
	var cfg *Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// GetConfigAsMap will reads the config into a string map array
func GetConfigAsMap() (map[string]interface{}, error) {
	// Read the configuration file
	var cfg map[string]interface{}
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// GetContext retrieves a context by name
func GetContext(name string) (*Context, error) {
	// Read the configuration file
	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}

	// Select only what we need
	if name == "" {
		return GetCurrentContext()
	}
	return selectContext(name, cfg.Contexts)
}

// GetAllContexts retrieves all available contexts
func GetAllContexts() ([]*Context, error) {
	// Read the configuration file
	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}

	return cfg.Contexts, nil
}

// GetCurrentContext retrieves all available contexts
func GetCurrentContext() (*Context, error) {
	// Read the configuration file
	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}

	// Select only what we need
	return selectContext(cfg.CurrentContext, cfg.Contexts)
}

func selectContext(name string, contexts []*Context) (*Context, error) {
	for _, element := range contexts {
		if element.Name == name {
			return element, nil
		}
	}
	return nil, fmt.Errorf("No context with name %s could be found", name)
}
