package linkerd

//Config struct used to unmarshal linkerd configuration from configmap in kubernetes
//Structs have fields for both linkerd v2.9.x and v2.10.x
type Config struct {
	LinkerdVersion          string   `yaml:"linkerdVersion"`
	IdentityTrustAnchorsPEM string   `yaml:"identityTrustAnchorsPEM"`
	Identity                Identity `yaml:"identity"`
	Global                  Global   `yaml:"global"`
}

//Identity partially defines identity part of linkerd configuration
type Identity struct {
	Issuer Issuer `yaml:"issuer"`
}

//Issuer partially defines issuer part of linkerd configuration
type Issuer struct {
	TLS    TLS    `yaml:"tls"`
	Expiry string `yaml:"crtExpiry"`
}

//TLS defines TLS part of linkerd configuration
type TLS struct {
	CrtPEM string `yaml:"crtPEM"`
}

//Global defines global part of linkerd 2.9.x configuration
type Global struct {
	IdentityTrustAnchorsPEM string `yaml:"identityTrustAnchorsPEM"`
	LinkerdVersion          string `yaml:"linkerdVersion"`
}
