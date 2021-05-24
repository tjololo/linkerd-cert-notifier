package linkerd

type LinkerdConfig struct {
	LinkerdVersion string `yaml:"linkerdVersion"`
	IdentityTrustAnchorsPEM string `yaml:"identityTrustAnchorsPEM"`
	Identity Identity `yaml:"identity"`
}

type Identity struct {
	Issuer Issuer `yaml:"issuer"`
}

type Issuer struct {
	Tls Tls `yaml:"tls"`
	Expiry string `yaml:"crtExpiry"`
}

type Tls struct {
	CrtPEM string `yaml:"crtPEM"`
}