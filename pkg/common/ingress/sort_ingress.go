/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ingress

import (
	"crypto/x509"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// SSLCert describes a SSL certificate to be used in a server
type SSLCert struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Certificate       *x509.Certificate `json:"certificate,omitempty"`
	// CAFileName contains the path to the file with the root certificate
	CAFileName string `json:"caFileName"`
	// CRLFileName containst the path to the CA CRL file
	CRLFileName string `json:"crlFileName"`
	// PemFileName contains the path to the file with the certificate and key concatenated
	PemFileName string `json:"pemFileName"`
	// PemSHA contains the sha1 of the pem file.
	// This is used to detect changes in the secret that contains the certificates
	PemSHA string `json:"pemSha"`
	// CN contains all the common names defined in the SSL certificate
	CN []string `json:"cn"`
	// ExpiresTime contains the expiration of this SSL certificate in timestamp format
	ExpireTime time.Time `json:"expires"`
}

// GetObjectKind implements the ObjectKind interface as a noop
func (s SSLCert) GetObjectKind() schema.ObjectKind {
	return schema.EmptyObjectKind
}
