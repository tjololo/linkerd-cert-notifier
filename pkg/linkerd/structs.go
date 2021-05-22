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
}

type Tls struct {
	crtPEM string `yaml:"crtPEM"`
}