package appconfig

// ContextConfig represents the configuration for managing a single context.
type ContextConfig struct {
	// Name holds the context name.
	Name string `yaml:"name"`
	// Cluster holds the cluster name.
	Cluster string `yaml:"cluster"`
	// User holds the username.
	User string `yaml:"user"`
}

// GetCluster returns the cluster configuration for the context.
func (c *ContextConfig) GetCluster(ctx AppConfig) *ClusterConfig {
	for _, cluster := range ctx.Clusters {
		if cluster.Name == c.Cluster {
			return &cluster
		}
	}
	return nil
}

// GetUser returns the user configuration for the context.
func (c *ContextConfig) GetUser(ctx AppConfig) *UserConfig {
	for _, user := range ctx.Users {
		if user.Name == c.User {
			return &user
		}
	}
	return nil
}
