// Copyright 2023 The Perses Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package secret

import (
	"crypto/tls"

	promConfig "github.com/prometheus/common/config"
)

// PublicTLSConfig is the public struct of TLSConfig.
// It's used when the API returns a response to a request
type PublicTLSConfig struct {
	CA                 Hidden `yaml:"ca,omitempty" json:"ca,omitempty"`
	Cert               Hidden `yaml:"cert,omitempty" json:"cert,omitempty"`
	Key                Hidden `yaml:"key,omitempty" json:"key,omitempty"`
	CAFile             string `yaml:"caFile,omitempty" json:"caFile,omitempty"`
	CertFile           string `yaml:"certFile,omitempty" json:"certFile,omitempty"`
	KeyFile            string `yaml:"keyFile,omitempty" json:"keyFile,omitempty"`
	ServerName         string `yaml:"serverName,omitempty" json:"serverName,omitempty"`
	InsecureSkipVerify bool   `yaml:"insecureSkipVerify,omitempty" json:"insecureSkipVerify,omitempty"`
	MinVersion         string `yaml:"minVersion,omitempty" json:"minVersion,omitempty"`
	MaxVersion         string `yaml:"maxVersion,omitempty" json:"maxVersion,omitempty"`
}

func NewPublicTLSConfig(t *TLSConfig) *PublicTLSConfig {
	if t == nil {
		return nil
	}
	return &PublicTLSConfig{
		CA:                 Hidden(t.CA),
		Cert:               Hidden(t.Cert),
		Key:                Hidden(t.Key),
		CAFile:             t.CAFile,
		CertFile:           t.CertFile,
		KeyFile:            t.KeyFile,
		ServerName:         t.ServerName,
		InsecureSkipVerify: t.InsecureSkipVerify,
		MinVersion:         t.MinVersion,
		MaxVersion:         t.MaxVersion,
	}
}

func BuildTLSConfig(cfg *TLSConfig) (*tls.Config, error) {
	if cfg == nil {
		return &tls.Config{MinVersion: tls.VersionTLS12, MaxVersion: tls.VersionTLS13}, nil
	}
	minVersion := promConfig.TLSVersions["TLS12"]
	maxVersion := promConfig.TLSVersions["TLS13"]
	if len(cfg.MinVersion) == 0 {
		minVersion = promConfig.TLSVersions[cfg.MinVersion]
	}
	if len(cfg.MaxVersion) == 0 {
		maxVersion = promConfig.TLSVersions[cfg.MaxVersion]
	}
	preConfig := &promConfig.TLSConfig{
		CA:                 cfg.CA,
		Cert:               cfg.Cert,
		Key:                promConfig.Secret(cfg.Key),
		CAFile:             cfg.CAFile,
		CertFile:           cfg.CertFile,
		KeyFile:            cfg.KeyFile,
		ServerName:         cfg.ServerName,
		InsecureSkipVerify: cfg.InsecureSkipVerify,
		MinVersion:         minVersion,
		MaxVersion:         maxVersion,
	}
	return promConfig.NewTLSConfig(preConfig)
}

type TLSConfig struct {
	// Text of the CA cert to use for the targets.
	CA string `yaml:"ca,omitempty" json:"ca,omitempty"`
	// Text of the client cert file for the targets.
	Cert string `yaml:"cert,omitempty" json:"cert,omitempty"`
	// Text of the client key file for the targets.
	Key string `yaml:"key,omitempty" json:"key,omitempty"`
	// The CA cert to use for the targets.
	CAFile string `yaml:"caFile,omitempty" json:"caFile,omitempty"`
	// The client cert file for the targets.
	CertFile string `yaml:"certFile,omitempty" json:"certFile,omitempty"`
	// The client key file for the targets.
	KeyFile string `yaml:"keyFile,omitempty" json:"keyFile,omitempty"`
	// Used to verify the hostname for the targets.
	ServerName string `yaml:"serverName,omitempty" json:"serverName,omitempty"`
	// Disable target certificate validation.
	InsecureSkipVerify bool `yaml:"insecureSkipVerify,omitempty" json:"insecureSkipVerify,omitempty"`
	// Minimum acceptable TLS version. Accepted values: TLS10 (TLS 1.0), TLS11 (TLS 1.1), TLS12 (TLS 1.2), TLS13 (TLS 1.3).
	// If unset, Perses will use Go default minimum version, which is TLS 1.2.
	// See MinVersion in https://pkg.go.dev/crypto/tls#Config.
	MinVersion string `yaml:"minVersion,omitempty" json:"minVersion,omitempty"`
	// Maximum acceptable TLS version. Accepted values: TLS10 (TLS 1.0), TLS11 (TLS 1.1), TLS12 (TLS 1.2), TLS13 (TLS 1.3).
	// If unset, Perses will use Go default maximum version, which is TLS 1.3.
	// See MaxVersion in https://pkg.go.dev/crypto/tls#Config.
	MaxVersion string `yaml:"maxVersion,omitempty" json:"maxVersion,omitempty"`
}

// Verify checks if the TLSConfig is valid.
// It also set the default value if needed
// This method is called when Perses is loading the configuration.
func (t *TLSConfig) Verify() error {
	if t == nil {
		return nil
	}
	if len(t.MinVersion) == 0 {
		t.MinVersion = "TLS12"
	}
	if len(t.MaxVersion) == 0 {
		t.MaxVersion = "TLS13"
	}
	return nil
}
