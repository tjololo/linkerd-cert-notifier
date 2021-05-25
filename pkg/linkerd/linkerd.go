package linkerd

import (
	"context"
	"fmt"

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
