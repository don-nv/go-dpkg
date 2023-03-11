package dhttp

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
)

// MustNewTLSConfig - see NewTLSConfig.
func MustNewTLSConfig(options ...TLSConfigOption) *tls.Config {
	config, err := NewTLSConfig(options...)
	if err != nil {
		panic(err)
	}

	return config
}

/*
NewTLSConfig

Defaults:
  - tls.ClientAuthType is tls.RequireAndVerifyClientCert. See OptionTLSConfigServerWithAuth() option;
  - tls min version is tls.VersionTLS12. See OptionTLSConfigWithMinVersion() option;
*/
func NewTLSConfig(options ...TLSConfigOption) (*tls.Config, error) {
	var config = tls.Config{} //nolint:gosec

	err := OptionTLSConfigServerWithAuth(tls.RequireAndVerifyClientCert)(&config)
	if err != nil {
		return nil, fmt.Errorf("applying tls server with auth option: %w", err)
	}
	err = OptionTLSConfigWithMinVersion(tls.VersionTLS12)(&config)
	if err != nil {
		return nil, fmt.Errorf("applying tls with min version option: %w", err)
	}

	for _, option := range options {
		oerr := option(&config)
		if oerr != nil {
			err = errors.Join(err, fmt.Errorf("optioning: %w", oerr))
		}
	}

	return &config, err
}

type TLSConfigOption func(config *tls.Config) error

func OptionTLSConfigClientWithMTLSPEM(cert, key, ca []byte) TLSConfigOption {
	return func(config *tls.Config) error {
		err := OptionTLSConfigWithCertPEM(cert, key)(config)
		if err != nil {
			return fmt.Errorf("OptionTLSConfigWithCertPEM: %w", err)
		}

		err = OptionTLSConfigClientWithCAPEM(ca)(config)
		if err != nil {
			return fmt.Errorf("OptionTLSConfigClientWithCAPEM: %w", err)
		}

		err = OptionTLSConfigClientWithSecure()(config)
		if err != nil {
			return fmt.Errorf("OptionTLSConfigClientWithSecure: %w", err)
		}

		return nil
	}
}

func OptionTLSConfigServerWithMTLSPEM(cert, key, ca []byte) TLSConfigOption {
	return func(config *tls.Config) error {
		err := OptionTLSConfigWithCertPEM(cert, key)(config)
		if err != nil {
			return fmt.Errorf("OptionTLSConfigWithCertPEM: %w", err)
		}

		err = OptionTLSConfigServerWithCAPEM(ca)(config)
		if err != nil {
			return fmt.Errorf("OptionTLSConfigServerWithCAPEM: %w", err)
		}

		err = OptionTLSConfigServerWithAuth(tls.RequireAndVerifyClientCert)(config)
		if err != nil {
			return fmt.Errorf("OptionTLSConfigServerWithAuth: %w", err)
		}

		err = OptionTLSConfigWithMinVersion(tls.VersionTLS12)(config)
		if err != nil {
			return fmt.Errorf("OptionTLSConfigWithMinVersion: %w", err)
		}

		return nil
	}
}

func OptionTLSConfigClientWithSecure() TLSConfigOption {
	return OptionTLSCConfigClientWithSecurity(true)
}

func OptionTLSConfigClientWithInsecure() TLSConfigOption {
	return OptionTLSCConfigClientWithSecurity(false)
}

func OptionTLSCConfigClientWithSecurity(verifies bool) TLSConfigOption {
	return func(config *tls.Config) error {
		config.InsecureSkipVerify = !verifies

		return nil
	}
}

func OptionTLSConfigServerWithAuth(authType tls.ClientAuthType) TLSConfigOption {
	return func(config *tls.Config) error {
		config.ClientAuth = authType

		if authType < tls.NoClientCert {
			config.ClientAuth = tls.NoClientCert
		}

		return nil
	}
}

func OptionTLSConfigWithCertPEM(cert, key []byte) TLSConfigOption {
	return func(config *tls.Config) error {
		crt, err := tls.X509KeyPair(cert, key)
		if err != nil {
			return fmt.Errorf("parsing x509 pem key pair: %w", err)
		}

		config.Certificates = append(config.Certificates, crt)

		return nil
	}
}

func OptionTLSConfigWithMinVersion(v uint16) TLSConfigOption {
	return func(config *tls.Config) error {
		config.MinVersion = v

		return nil
	}
}

/*
OptionTLSConfigClientWithCAPEM - adds CA to the set of root certificate authorities that clients use when verifying
server certificates.
*/
func OptionTLSConfigClientWithCAPEM(ca []byte) TLSConfigOption {
	return func(config *tls.Config) error {
		if config.RootCAs == nil {
			config.RootCAs = x509.NewCertPool()
		}

		if !config.RootCAs.AppendCertsFromPEM(ca) {
			return errors.New("appending servers pem ca: not appended")
		}

		return nil
	}
}

/*
OptionTLSConfigServerWithCAPEM - adds CA to the set of root certificate authorities that servers use if required
to verify a client certificate by the policy in tls.ClientAuthType.
*/
func OptionTLSConfigServerWithCAPEM(ca []byte) TLSConfigOption {
	return func(config *tls.Config) error {
		if config.ClientCAs == nil {
			config.ClientCAs = x509.NewCertPool()
		}

		if !config.ClientCAs.AppendCertsFromPEM(ca) {
			return errors.New("appending clients pem ca: not appended")
		}

		return nil
	}
}
