package linkerd

import (
	"context"
	"fmt"

	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type LinkerdReader struct {
	Client *kubernetes.Clientset
}

func (l *LinkerdReader) FetchTrustAnchor(namespace string, ctx context.Context) (anchorPem []byte, err error) {
	configMapName := "linkerd-config"
	configmap, err := l.Client.CoreV1().ConfigMaps(namespace).Get(ctx, configMapName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil, fmt.Errorf("configmap %s not found in namespace %s", configMapName, namespace)
	} else if err != nil {
		return nil, fmt.Errorf("failed to retrive configmap %s from namespace %s error: %s", configMapName, namespace, err)
	}
	values := configmap.Data["values"]
	config := LinkerdConfig{}
	err = yaml.Unmarshal([]byte(values), &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal values from configmap. %s", err)
	}
	return []byte(config.IdentityTrustAnchorsPEM), nil
}

func (l *LinkerdReader) FetchIssuerCert(namespace string, ctx context.Context) (issuerPem []byte, err error) {
	configMapName := "linkerd-config"
	configmap, err := l.Client.CoreV1().ConfigMaps(namespace).Get(ctx, configMapName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil, fmt.Errorf("configmap %s not found in namespace %s", configMapName, namespace)
	} else if err != nil {
		return nil, fmt.Errorf("failed to retrive configmap %s from namespace %s error: %s", configMapName, namespace, err)
	}
	values := configmap.Data["values"]
	config := LinkerdConfig{}
	err = yaml.Unmarshal([]byte(values), &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal values from configmap. %s", err)
	}
	return []byte(config.Identity.Issuer.Tls.CrtPEM), nil
}