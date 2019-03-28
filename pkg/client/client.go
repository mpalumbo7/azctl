package client

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"golang.org/x/crypto/pkcs12"

	configpkg "github.com/mpalumbo7/azctl/pkg/config"
)

// NewClient returns a bootstrapped ADAL client
func NewClient(cxt *configpkg.Context) (*autorest.Client, error) {
	endpoint := cxt.Endpoint
	if endpoint == "" {
		endpoint = azure.PublicCloud.ActiveDirectoryEndpoint
	}
	oauthConfig, err := adal.NewOAuthConfig(endpoint, cxt.TenantID)
	if err != nil {
		return nil, err
	}

	// Generate token
	var spt *adal.ServicePrincipalToken
	if cxt.CertificateData != "" {
		spt, err = getSptFromCertificate(*oauthConfig, cxt.ClientID, cxt.Resource, cxt.CertificateData, cxt.SaveToken)
	} else {
		spt, err = getSptFromDeviceFlow(*oauthConfig, cxt.ClientID, cxt.Resource, cxt.SaveToken)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve token.\n %s", err)
	}

	client := &autorest.Client{}
	client.Authorizer = autorest.NewBearerAuthorizer(spt)
	return client, nil
}

// GetSptFromDeviceFlow retrieves a Service Principal Token from the Azure device login
func getSptFromDeviceFlow(oauthConfig adal.OAuthConfig, clientID string, resource string, callbacks ...adal.TokenRefreshCallback) (*adal.ServicePrincipalToken, error) {
	oauthClient := &autorest.Client{}
	deviceCode, err := adal.InitiateDeviceAuth(oauthClient, oauthConfig, clientID, resource)
	if err != nil {
		return nil, fmt.Errorf("failed to start device auth flow: %s", err)
	}
	// Required for user to know what to do
	fmt.Println(*deviceCode.Message)

	token, err := adal.WaitForUserCompletion(oauthClient, deviceCode)
	if err != nil {
		return nil, fmt.Errorf("failed to finish device auth flow: %s", err)
	}

	spt, err := adal.NewServicePrincipalTokenFromManualToken(
		oauthConfig,
		clientID,
		resource,
		*token,
		callbacks...)
	if err != nil {
		return nil, fmt.Errorf("failed to get oauth token from device flow: %v", err)
	}

	return spt, nil
}

// GetSptFromCertificate retrieves a Service Principal Token from an Azure certificate
func getSptFromCertificate(oauthConfig adal.OAuthConfig, clientID, resource, certificatePath string, callbacks ...adal.TokenRefreshCallback) (*adal.ServicePrincipalToken, error) {
	certData, err := ioutil.ReadFile(certificatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read the certificate file (%s): %v", certificatePath, err)
	}

	certificate, rsaPrivateKey, err := decodePkcs12(certData, "")
	if err != nil {
		return nil, fmt.Errorf("failed to decode pkcs12 certificate while creating spt: %v", err)
	}

	spt, _ := adal.NewServicePrincipalTokenFromCertificate(
		oauthConfig,
		clientID,
		certificate,
		rsaPrivateKey,
		resource,
		callbacks...)

	return spt, nil
}

func decodePkcs12(pkcs []byte, password string) (*x509.Certificate, *rsa.PrivateKey, error) {
	privateKey, certificate, err := pkcs12.Decode(pkcs, password)
	if err != nil {
		return nil, nil, err
	}

	rsaPrivateKey, isRsaKey := privateKey.(*rsa.PrivateKey)
	if !isRsaKey {
		return nil, nil, fmt.Errorf("PKCS#12 certificate must contain an RSA private key")
	}

	return certificate, rsaPrivateKey, nil
}
