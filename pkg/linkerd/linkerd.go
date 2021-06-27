package linkerd

import (
	"context"
	"fmt"
	"regexp"

	semver "github.com/Masterminds/semver/v3"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//Reader struct used for reading linkerd config from kubernetes
type Reader struct {
	Client *kubernetes.Clientset
}

//FetchTrustAnchor reads trustAnchorPEM from linkerd-config in kubernetes
func (l *Reader) FetchTrustAnchor(ctx context.Context, namespace string) (anchorPem []byte, err error) {
	configMapName := "linkerd-config"
	configmap, err := l.Client.CoreV1().ConfigMaps(namespace).Get(ctx, configMapName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil, fmt.Errorf("configmap %s not found in namespace %s", configMapName, namespace)
	} else if err != nil {
		return nil, fmt.Errorf("failed to retrive configmap %s from namespace %s error: %s", configMapName, namespace, err)
	}
	values := configmap.Data["values"]
	config := Config{}
	err = yaml.Unmarshal([]byte(values), &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal values from configmap. %s", err)
	}
	version, err := getLinkerdSemver(config.Global.LinkerdVersion, config.LinkerdVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch linkerd version. %s", err)
	}
	if version.LessThan(semver.MustParse("2.10.0")) {
		return []byte(config.Global.IdentityTrustAnchorsPEM), nil
	}
	return []byte(config.IdentityTrustAnchorsPEM), nil
}

//FetchIssuerCert reads Issuer crtPEM from linkerd-config in kubernetes
func (l *Reader) FetchIssuerCert(ctx context.Context, namespace string) (issuerPem []byte, err error) {
	configMapName := "linkerd-config"
	configmap, err := l.Client.CoreV1().ConfigMaps(namespace).Get(ctx, configMapName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil, fmt.Errorf("configmap %s not found in namespace %s", configMapName, namespace)
	} else if err != nil {
		return nil, fmt.Errorf("failed to retrive configmap %s from namespace %s error: %s", configMapName, namespace, err)
	}
	values := configmap.Data["values"]
	config := Config{}
	err = yaml.Unmarshal([]byte(values), &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal values from configmap. %s", err)
	}
	return []byte(config.Identity.Issuer.TLS.CrtPEM), nil
}

func getLinkerdSemver(globalversionString string, versionString string) (*semver.Version, error) {
	r := regexp.MustCompile(`^stable-(?P<Version>\d+\.\d+\.\d+)$`)
	if r.MatchString(versionString) {
		n := r.FindAllStringSubmatch(versionString, 1)
		semverString := n[0][1]
		version, err := semver.NewVersion(semverString)
		return version, err
	}
	if r.MatchString(globalversionString) {
		n := r.FindAllStringSubmatch(globalversionString, 1)
		semverString := n[0][1]
		version, err := semver.NewVersion(semverString)
		return version, err
	}
	return nil, fmt.Errorf("linkerd version did not match expected stable format stable-X.Y.Z, got %s", versionString)
}
