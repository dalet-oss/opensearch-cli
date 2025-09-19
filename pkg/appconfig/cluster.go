package appconfig

// ClusterConfig represents the configuration for managing a single cluster.
type ClusterConfig struct {
	// Name holds the cluster name.
	Name string `yaml:"name"`
	// Params holds the cluster parameters.
	Params ClusterParams `yaml:"params"`
}

// ClusterParams represents the parameters for managing a single cluster.
type ClusterParams struct {
	// Server holds the cluster server address.
	Server string `yaml:"server"`
	// Tls indicates whether TLS encryption is enabled for the cluster connection.
	Tls bool `yaml:"tls"`
}
