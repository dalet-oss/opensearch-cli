package appconfig

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/fp"
	"fmt"
	"slices"
	"strings"
	"time"
)

// this is package is managing context for the application

const ApiVersionV1 = "v1"
const DefaultServerTimeoutSeconds = 60

// AppConfig represents the configuration for managing multiple clusters, users, and contexts settings.
// Clusters defines a slice of cluster configurations, including name and parameters.
// Users defines a slice of user configurations, including name and password.
// Contexts defines a slice of context configurations mapping clusters, users, and contexts.
// Current specifies the currently active context name.
type AppConfig struct {
	CliParams  *CliParams `yaml:"params,omitempty"`
	ApiVersion string     `yaml:"apiVersion"`
	// Clusters holds the configuration for multiple cluster settings, including name and related parameters.
	Clusters []ClusterConfig `yaml:"clusters"`
	// Users holds the configuration for multiple user settings, including name and password.
	Users []UserConfig `yaml:"users"`
	// Contexts holds the configuration for multiple context settings, including name, cluster and user.
	Contexts []ContextConfig `yaml:"contexts"`
	// Current holds the current context name.
	Current string `yaml:"current"`
}

type CliParams struct {
	ServerTimeoutSeconds *int `yaml:"serverTimeoutSeconds"`
}

func (p *CliParams) GetServerTimeoutSeconds() int {
	if p.ServerTimeoutSeconds == nil {
		return DefaultServerTimeoutSeconds
	}
	return *p.ServerTimeoutSeconds
}

func (c *AppConfig) ShowContextInfo(name string) string {
	info := []string{}
	if ctx := c.GetContext(name); ctx != nil {
		info = append(info, fmt.Sprintf("✅Context: %s", name))
		if cluster := c.GetCluster(ctx.Cluster); cluster != nil {
			info = append(info, fmt.Sprintf("	✅Cluster: %s", cluster.Name))
		} else {
			info = append(info, fmt.Sprintf("	❌Cluster: %s", ctx.Cluster))
		}
		if user := c.GetUser(ctx.User); user != nil {
			info = append(info, fmt.Sprintf("	✅User: %s", user.Name))
		} else {
			info = append(info, fmt.Sprintf("	❌User: %s", ctx.User))
		}
	} else {
		info = append(info, fmt.Sprintf("❌Context: %s", name))
	}
	return strings.Join(info, "\n")
}

func (c *AppConfig) ShowContextInfoExtended(name string) string {
	info := []string{}
	if ctx := c.GetContext(name); ctx != nil {
		info = append(info, fmt.Sprintf("Context:%s", name))
		if cluster := c.GetCluster(ctx.Cluster); cluster != nil {
			info = append(info, fmt.Sprintf("	Cluster: %s", cluster.Name))
			info = append(info, fmt.Sprintf("		Server: %s", cluster.Params.Server))
			info = append(info, fmt.Sprintf("		TLS: %t", cluster.Params.Tls))
		} else {
			info = append(info, "	❌cluster is not found")
		}
		if user := c.GetUser(ctx.User); user != nil {
			info = append(info, fmt.Sprintf("	User: %s", user.Name))
			info = append(info, fmt.Sprintf("		Token: %s", user.User.Token))
		} else {
			info = append(info, fmt.Sprintf("	❌User is not found: %s", ctx.User))
		}
	} else {
		info = append(info, fmt.Sprintf("❌Context is not found: %s", name))
	}
	return strings.Join(info, "\n")
}

func (c *AppConfig) GetCluster(name string) *ClusterConfig {
	if idx := slices.IndexFunc(c.Clusters, func(p ClusterConfig) bool {
		return p.Name == name
	}); idx == -1 {
		return nil
	} else {
		return &c.Clusters[idx]
	}
}
func (c *AppConfig) GetContext(name string) *ContextConfig {
	if idx := slices.IndexFunc(c.Contexts, func(p ContextConfig) bool {
		return p.Name == name
	}); idx == -1 {
		return nil
	} else {
		return &c.Contexts[idx]
	}
}
func (c *AppConfig) GetUser(name string) *UserConfig {
	if idx := slices.IndexFunc(c.Users, func(p UserConfig) bool {
		return p.Name == name
	}); idx == -1 {
		return nil
	} else {
		return &c.Users[idx]
	}
}
func (c *AppConfig) HasCluster(candidate ClusterConfig) bool {
	return slices.ContainsFunc(c.Clusters, func(p ClusterConfig) bool {
		return p.Name == candidate.Name
	})
}

func (c *AppConfig) HasUser(candidate UserConfig) bool {
	return slices.ContainsFunc(c.Users, func(p UserConfig) bool {
		return p.Name == candidate.Name
	})
}
func (c *AppConfig) HasContext(candidate ContextConfig) bool {
	return slices.ContainsFunc(c.Contexts, func(p ContextConfig) bool {
		return p.Name == candidate.Name
	})
}
func (c *AppConfig) GetActiveContext() *ContextConfig {
	if c.Current == "" {
		return nil
	}
	for _, ctx := range c.Contexts {
		if ctx.Name == c.Current {
			return &ctx
		}
	}
	return nil
}

func (c *AppConfig) Push(cluster ClusterConfig, user UserConfig, ctx ContextConfig) {
	c.pushCluster(cluster)
	c.pushUser(user)
	c.pushContext(ctx)
}

func (c *AppConfig) pushCluster(e ClusterConfig) {
	c.Clusters = append(c.Clusters, e)
}

func (c *AppConfig) pushUser(e UserConfig) {
	c.Users = append(c.Users, e)
}

func (c *AppConfig) pushContext(e ContextConfig) {
	c.Contexts = append(c.Contexts, e)
}

func (c *AppConfig) ListContexts() {
	for _, ctx := range c.Contexts {
		fmt.Println(c.ShowContextInfo(ctx.Name))
	}
	fmt.Println("Current context: ", c.Current)
}

func (c *AppConfig) GetContextList() []string {
	ctxs := []string{}
	for _, ctx := range c.Contexts {
		ctxs = append(ctxs, ctx.Name)
	}
	return ctxs
}

// Example returns an example AppConfig.
func Example() AppConfig {
	return AppConfig{
		ApiVersion: ApiVersionV1,
		CliParams:  &CliParams{ServerTimeoutSeconds: fp.AsPointer(DefaultServerTimeoutSeconds)},
		Clusters: []ClusterConfig{
			{
				Name: "example-cluster",
				Params: ClusterParams{
					Server: "http://localhost:9200",
					Tls:    false,
				},
			},
		},
		Users: []UserConfig{
			{
				Name: "example-user",
				User: User{
					Token: "UUID | This is a not real thing please use cli to push actual values as it's using your local OS secret storage",
				},
			},
		},
		Contexts: []ContextConfig{
			{
				Name:    "example-context",
				Cluster: "example-cluster",
				User:    "example-user",
			},
		},
		Current: "example-context",
	}
}

func (c *AppConfig) ServerCallTimeout() time.Duration {
	v := DefaultServerTimeoutSeconds
	if c.CliParams != nil && c.CliParams.ServerTimeoutSeconds != nil {
		if *c.CliParams.ServerTimeoutSeconds <= 0 {
			v = DefaultServerTimeoutSeconds
		} else {
			v = *c.CliParams.ServerTimeoutSeconds
		}
	}
	return time.Duration(v) * time.Second
}
