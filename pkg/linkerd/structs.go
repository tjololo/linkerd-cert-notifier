package linkerd

//Config struct used to unmarshal linkerd configuration from configmap in kubernetes
type Config struct {
	LinkerdVersion          string   `yaml:"linkerdVersion"`
	IdentityTrustAnchorsPEM string   `yaml:"identityTrustAnchorsPEM"`
	Identity                Identity `yaml:"identity"`
}

//Identity partialy defines identity part of linkerd configuration
type Identity struct {
	Issuer Issuer `yaml:"issuer"`
}

//Issuer partialy defines issuer part of linkerd configuration
type Issuer struct {
	TLS    TLS    `yaml:"tls"`
	Expiry string `yaml:"crtExpiry"`
}

//TLS defines TLS part of linkerd configuration
type TLS struct {
	CrtPEM string `yaml:"crtPEM"`
}
