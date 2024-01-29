/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
limitations under the License.
*/
package conjur

import (
	"fmt"

	// nolint
	. "github.com/onsi/ginkgo/v2"
	v1 "k8s.io/api/core/v1"

	"github.com/external-secrets/external-secrets-e2e/framework"
	"github.com/external-secrets/external-secrets-e2e/suites/provider/cases/common"
	esapi "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
)

const (
	withTokenAuth = "with apikey auth"
	withJWT       = "with jwt provider"
	withJWTK8s    = "with jwt k8s provider"
	// withK8s          = "with kubernetes provider"
)

var _ = Describe("[conjur]", Label("conjur"), func() {
	f := framework.New("eso-conjur")
	prov := newConjurProvider(f)

	DescribeTable("sync secrets",
		framework.TableFunc(f, prov),
		// uses token auth
		framework.Compose(withTokenAuth, f, common.SimpleDataSync, useApiKeyAuth),
		framework.Compose(withTokenAuth, f, common.SyncWithoutTargetName, useApiKeyAuth),
		// framework.Compose(withTokenAuth, f, common.FindByName, useApiKeyAuth),
		// framework.Compose(withTokenAuth, f, common.FindByNameAndRewrite, useApiKeyAuth),
		// framework.Compose(withTokenAuth, f, common.JSONDataFromSync, useApiKeyAuth),
		// framework.Compose(withTokenAuth, f, common.JSONDataFromRewrite, useApiKeyAuth),
		// framework.Compose(withTokenAuth, f, common.JSONDataWithProperty, useApiKeyAuth),
		// framework.Compose(withTokenAuth, f, common.JSONDataWithTemplate, useApiKeyAuth),
		// framework.Compose(withTokenAuth, f, common.DataPropertyDockerconfigJSON, useApiKeyAuth),
		// framework.Compose(withTokenAuth, f, common.JSONDataWithoutTargetName, useApiKeyAuth),
		// framework.Compose(withTokenAuth, f, common.SyncV1Alpha1, useApiKeyAuth),
		// framework.Compose(withTokenAuth, f, common.DecodingPolicySync, useApiKeyAuth),
		// framework.Compose(withTokenAuth, f, common.JSONDataWithTemplateFromLiteral, useApiKeyAuth),
		// framework.Compose(withTokenAuth, f, common.TemplateFromConfigmaps, useApiKeyAuth),
		// use jwt provider
		// framework.Compose(withJWT, f, common.FindByName, useJWTProvider),
		// framework.Compose(withJWT, f, common.FindByNameAndRewrite, useJWTProvider),
		// framework.Compose(withJWT, f, common.JSONDataFromSync, useJWTProvider),
		// framework.Compose(withJWT, f, common.JSONDataFromRewrite, useJWTProvider),
		// framework.Compose(withJWT, f, common.JSONDataWithProperty, useJWTProvider),
		// framework.Compose(withJWT, f, common.JSONDataWithTemplate, useJWTProvider),
		// framework.Compose(withJWT, f, common.DataPropertyDockerconfigJSON, useJWTProvider),
		// framework.Compose(withJWT, f, common.JSONDataWithoutTargetName, useJWTProvider),
		// // use jwt k8s provider
		// framework.Compose(withJWTK8s, f, common.JSONDataFromSync, useJWTK8sProvider),
		// framework.Compose(withJWTK8s, f, common.JSONDataFromRewrite, useJWTK8sProvider),
		// framework.Compose(withJWTK8s, f, common.JSONDataWithProperty, useJWTK8sProvider),
		// framework.Compose(withJWTK8s, f, common.JSONDataWithTemplate, useJWTK8sProvider),
		// framework.Compose(withJWTK8s, f, common.DataPropertyDockerconfigJSON, useJWTK8sProvider),
		// framework.Compose(withJWTK8s, f, common.JSONDataWithoutTargetName, useJWTK8sProvider),
		// use kubernetes provider
		// framework.Compose(withK8s, f, common.FindByName, useKubernetesProvider),
		// framework.Compose(withK8s, f, common.FindByNameAndRewrite, useKubernetesProvider),
		// framework.Compose(withK8s, f, common.JSONDataFromSync, useKubernetesProvider),
		// framework.Compose(withK8s, f, common.JSONDataFromRewrite, useKubernetesProvider),
		// framework.Compose(withK8s, f, common.JSONDataWithProperty, useKubernetesProvider),
		// framework.Compose(withK8s, f, common.JSONDataWithTemplate, useKubernetesProvider),
		// framework.Compose(withK8s, f, common.DataPropertyDockerconfigJSON, useKubernetesProvider),
		// framework.Compose(withK8s, f, common.JSONDataWithoutTargetName, useKubernetesProvider),
		// conjur-specific test cases
		Entry("secret value via data without property should return json-encoded string", Label("json"), testJSONWithoutProperty),
		// Entry("secret value via data with property should return json-encoded string", Label("json"), testJSONWithProperty),
		// Entry("dataFrom without property should extract key/value pairs", Label("json"), testDataFromJSONWithoutProperty),
		// Entry("dataFrom with property should extract key/value pairs", Label("json"), testDataFromJSONWithProperty),
	)
})

func useApiKeyAuth(tc *framework.TestCase) {
	tc.ExternalSecret.Spec.SecretStoreRef.Name = tc.Framework.Namespace.Name
}

// func useJWTProvider(tc *framework.TestCase) {
// 	tc.ExternalSecret.Spec.SecretStoreRef.Name = jwtProviderName
// }

// func useJWTK8sProvider(tc *framework.TestCase) {
// 	tc.ExternalSecret.Spec.SecretStoreRef.Name = jwtK8sProviderName
// }

// func useKubernetesProvider(tc *framework.TestCase) {
// 	tc.ExternalSecret.Spec.SecretStoreRef.Name = kubernetesProviderName
// }

const jsonVal = `{"foo":{"nested":{"bar":"mysecret","baz":"bang"}}}`

// when no property is set it should return the json-encoded at path.
func testJSONWithoutProperty(tc *framework.TestCase) {
	secretKey := fmt.Sprintf("%s-%s", tc.Framework.Namespace.Name, "json")
	tc.Secrets = map[string]framework.SecretEntry{
		secretKey: {Value: jsonVal},
	}
	tc.ExpectedSecret = &v1.Secret{
		Type: v1.SecretTypeOpaque,
		Data: map[string][]byte{
			secretKey: []byte(jsonVal),
		},
	}
	tc.ExternalSecret.Spec.Data = []esapi.ExternalSecretData{
		{
			SecretKey: secretKey,
			RemoteRef: esapi.ExternalSecretDataRemoteRef{
				Key: secretKey,
			},
		},
	}
}

// when property is set it should return the json-encoded at path.
func testJSONWithProperty(tc *framework.TestCase) {
	secretKey := fmt.Sprintf("%s-%s", tc.Framework.Namespace.Name, "json")
	expectedVal := `{"bar":"mysecret","baz":"bang"}`
	tc.Secrets = map[string]framework.SecretEntry{
		secretKey: {Value: jsonVal},
	}
	tc.ExpectedSecret = &v1.Secret{
		Type: v1.SecretTypeOpaque,
		Data: map[string][]byte{
			secretKey: []byte(expectedVal),
		},
	}
	tc.ExternalSecret.Spec.Data = []esapi.ExternalSecretData{
		{
			SecretKey: secretKey,
			RemoteRef: esapi.ExternalSecretDataRemoteRef{
				Key:      secretKey,
				Property: "foo.nested",
			},
		},
	}
}

// when no property is set it should extract the key/value pairs at the given path
// note: it should json-encode if a value contains nested data
func testDataFromJSONWithoutProperty(tc *framework.TestCase) {
	secretKey := fmt.Sprintf("%s-%s", tc.Framework.Namespace.Name, "json")
	tc.Secrets = map[string]framework.SecretEntry{
		secretKey: {Value: jsonVal},
	}
	tc.ExpectedSecret = &v1.Secret{
		Type: v1.SecretTypeOpaque,
		Data: map[string][]byte{
			"foo": []byte(`{"nested":{"bar":"mysecret","baz":"bang"}}`),
		},
	}
	tc.ExternalSecret.Spec.DataFrom = []esapi.ExternalSecretDataFromRemoteRef{
		{
			Extract: &esapi.ExternalSecretDataRemoteRef{
				Key: secretKey,
			},
		},
	}
}

// when property is set it should extract values with dataFrom at the given path.
func testDataFromJSONWithProperty(tc *framework.TestCase) {
	secretKey := fmt.Sprintf("%s-%s", tc.Framework.Namespace.Name, "json")
	tc.Secrets = map[string]framework.SecretEntry{
		secretKey: {Value: jsonVal},
	}
	tc.ExpectedSecret = &v1.Secret{
		Type: v1.SecretTypeOpaque,
		Data: map[string][]byte{
			"bar": []byte(`mysecret`),
			"baz": []byte(`bang`),
		},
	}
	tc.ExternalSecret.Spec.DataFrom = []esapi.ExternalSecretDataFromRemoteRef{
		{
			Extract: &esapi.ExternalSecretDataRemoteRef{
				Key:      secretKey,
				Property: "foo.nested",
			},
		},
	}
}
