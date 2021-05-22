package linkerd

type LinkerdConfig struct {
	LinkerdVersion string
	IdentityTrustAnchorsPEM string
	Identity Identity
}

type Identity struct {
	Issuer Issuer
}

type Issuer struct {
	Tls Tls
}

type Tls struct {
	crtPEM string
}