package configuration

import (
	"testing"

	"github.com/containous/traefik/tls"
	"github.com/containous/traefik/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseEntryPointsConfiguration(t *testing.T) {
	testCases := []struct {
		name           string
		value          string
		expectedResult map[string]string
	}{
		{
			name: "all parameters",
			value: "Name:foo " +
				"Address::8000 " +
				"TLS:goo,gii " +
				"TLS " +
				"CA:car " +
				"CA.Optional:true " +
				"Redirect.EntryPoint:https " +
				"Redirect.Regex:http://localhost/(.*) " +
				"Redirect.Replacement:http://mydomain/$1 " +
				"Redirect.Permanent:true " +
				"Compress:true " +
				"WhiteListSourceRange:10.42.0.0/16,152.89.1.33/32,afed:be44::/16 " +
				"ProxyProtocol.TrustedIPs:192.168.0.1 " +
				"ForwardedHeaders.TrustedIPs:10.0.0.3/24,20.0.0.3/24 " +
				"Auth.Basic.Users:test:$apr1$H6uskkkW$IgXLP6ewTrSuBkTrqE8wj/,test2:$apr1$d9hr9HBB$4HxwgUir3HP4EsggP/QNo0 " +
				"Auth.Digest.Users:test:traefik:a2688e031edb4be6a3797f3882655c05,test2:traefik:518845800f9e2bfb1f1f740ec24f074e " +
				"Auth.HeaderField:X-WebAuth-User " +
				"Auth.Forward.Address:https://authserver.com/auth " +
				"Auth.Forward.TrustForwardHeader:true " +
				"Auth.Forward.TLS.CA:path/to/local.crt " +
				"Auth.Forward.TLS.CAOptional:true " +
				"Auth.Forward.TLS.Cert:path/to/foo.cert " +
				"Auth.Forward.TLS.Key:path/to/foo.key " +
				"Auth.Forward.TLS.InsecureSkipVerify:true ",
			expectedResult: map[string]string{
				"address":                             ":8000",
				"auth_basic_users":                    "test:$apr1$H6uskkkW$IgXLP6ewTrSuBkTrqE8wj/,test2:$apr1$d9hr9HBB$4HxwgUir3HP4EsggP/QNo0",
				"auth_digest_users":                   "test:traefik:a2688e031edb4be6a3797f3882655c05,test2:traefik:518845800f9e2bfb1f1f740ec24f074e",
				"auth_forward_address":                "https://authserver.com/auth",
				"auth_forward_tls_ca":                 "path/to/local.crt",
				"auth_forward_tls_caoptional":         "true",
				"auth_forward_tls_cert":               "path/to/foo.cert",
				"auth_forward_tls_insecureskipverify": "true",
				"auth_forward_tls_key":                "path/to/foo.key",
				"auth_forward_trustforwardheader":     "true",
				"auth_headerfield":                    "X-WebAuth-User",
				"ca":                                  "car",
				"ca_optional":                         "true",
				"compress":                            "true",
				"forwardedheaders_trustedips":         "10.0.0.3/24,20.0.0.3/24",
				"name": "foo",
				"proxyprotocol_trustedips": "192.168.0.1",
				"redirect_entrypoint":      "https",
				"redirect_permanent":       "true",
				"redirect_regex":           "http://localhost/(.*)",
				"redirect_replacement":     "http://mydomain/$1",
				"tls":                  "goo,gii",
				"tls_acme":             "TLS",
				"whitelistsourcerange": "10.42.0.0/16,152.89.1.33/32,afed:be44::/16",
			},
		},
		{
			name:  "compress on",
			value: "name:foo Compress:on",
			expectedResult: map[string]string{
				"name":     "foo",
				"compress": "on",
			},
		},
		{
			name:  "TLS",
			value: "Name:foo TLS:goo TLS",
			expectedResult: map[string]string{
				"name":     "foo",
				"tls":      "goo",
				"tls_acme": "TLS",
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			conf := parseEntryPointsConfiguration(test.value)

			assert.Len(t, conf, len(test.expectedResult))
			assert.Equal(t, test.expectedResult, conf)
		})
	}
}

func Test_toBool(t *testing.T) {
	testCases := []struct {
		name         string
		value        string
		key          string
		expectedBool bool
	}{
		{
			name:         "on",
			value:        "on",
			key:          "foo",
			expectedBool: true,
		},
		{
			name:         "true",
			value:        "true",
			key:          "foo",
			expectedBool: true,
		},
		{
			name:         "enable",
			value:        "enable",
			key:          "foo",
			expectedBool: true,
		},
		{
			name:         "arbitrary string",
			value:        "bar",
			key:          "foo",
			expectedBool: false,
		},
		{
			name:         "no existing entry",
			value:        "bar",
			key:          "fii",
			expectedBool: false,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			conf := map[string]string{
				"foo": test.value,
			}

			result := toBool(conf, test.key)

			assert.Equal(t, test.expectedBool, result)
		})
	}
}

func TestEntryPoints_Set(t *testing.T) {
	testCases := []struct {
		name                   string
		expression             string
		expectedEntryPointName string
		expectedEntryPoint     *EntryPoint
	}{
		{
			name: "all parameters camelcase",
			expression: "Name:foo " +
				"Address::8000 " +
				"TLS:goo,gii " +
				"TLS " +
				"CA:car " +
				"CA.Optional:true " +
				"Redirect.EntryPoint:https " +
				"Redirect.Regex:http://localhost/(.*) " +
				"Redirect.Replacement:http://mydomain/$1 " +
				"Redirect.Permanent:true " +
				"Compress:true " +
				"WhiteListSourceRange:10.42.0.0/16,152.89.1.33/32,afed:be44::/16 " +
				"ProxyProtocol.TrustedIPs:192.168.0.1 " +
				"ForwardedHeaders.TrustedIPs:10.0.0.3/24,20.0.0.3/24 " +
				"Auth.Basic.Users:test:$apr1$H6uskkkW$IgXLP6ewTrSuBkTrqE8wj/,test2:$apr1$d9hr9HBB$4HxwgUir3HP4EsggP/QNo0 " +
				"Auth.Digest.Users:test:traefik:a2688e031edb4be6a3797f3882655c05,test2:traefik:518845800f9e2bfb1f1f740ec24f074e " +
				"Auth.HeaderField:X-WebAuth-User " +
				"Auth.Forward.Address:https://authserver.com/auth " +
				"Auth.Forward.TrustForwardHeader:true " +
				"Auth.Forward.TLS.CA:path/to/local.crt " +
				"Auth.Forward.TLS.CAOptional:true " +
				"Auth.Forward.TLS.Cert:path/to/foo.cert " +
				"Auth.Forward.TLS.Key:path/to/foo.key " +
				"Auth.Forward.TLS.InsecureSkipVerify:true ",
			expectedEntryPointName: "foo",
			expectedEntryPoint: &EntryPoint{
				Address: ":8000",
				TLS: &tls.TLS{
					Certificates: tls.Certificates{
						{
							CertFile: tls.FileOrContent("goo"),
							KeyFile:  tls.FileOrContent("gii"),
						},
					},
					ClientCA: tls.ClientCA{
						Files:    []string{"car"},
						Optional: true,
					},
				},
				Redirect: &types.Redirect{
					EntryPoint:  "https",
					Regex:       "http://localhost/(.*)",
					Replacement: "http://mydomain/$1",
					Permanent:   true,
				},
				Auth: &types.Auth{
					Basic: &types.Basic{
						Users: types.Users{
							"test:$apr1$H6uskkkW$IgXLP6ewTrSuBkTrqE8wj/",
							"test2:$apr1$d9hr9HBB$4HxwgUir3HP4EsggP/QNo0",
						},
					},
					Digest: &types.Digest{
						Users: types.Users{
							"test:traefik:a2688e031edb4be6a3797f3882655c05",
							"test2:traefik:518845800f9e2bfb1f1f740ec24f074e",
						},
					},
					Forward: &types.Forward{
						Address: "https://authserver.com/auth",
						TLS: &types.ClientTLS{
							CA:                 "path/to/local.crt",
							CAOptional:         true,
							Cert:               "path/to/foo.cert",
							Key:                "path/to/foo.key",
							InsecureSkipVerify: true,
						},
						TrustForwardHeader: true,
					},
					HeaderField: "X-WebAuth-User",
				},
				WhitelistSourceRange: []string{
					"10.42.0.0/16",
					"152.89.1.33/32",
					"afed:be44::/16",
				},
				Compress: true,
				ProxyProtocol: &ProxyProtocol{
					Insecure:   false,
					TrustedIPs: []string{"192.168.0.1"},
				},
				ForwardedHeaders: &ForwardedHeaders{
					Insecure: false,
					TrustedIPs: []string{
						"10.0.0.3/24",
						"20.0.0.3/24",
					},
				},
			},
		},
		{
			name: "all parameters lowercase",
			expression: "Name:foo " +
				"address::8000 " +
				"tls:goo,gii " +
				"tls " +
				"ca:car " +
				"ca.Optional:true " +
				"redirect.entryPoint:https " +
				"redirect.regex:http://localhost/(.*) " +
				"redirect.replacement:http://mydomain/$1 " +
				"redirect.permanent:true " +
				"compress:true " +
				"whiteListSourceRange:10.42.0.0/16,152.89.1.33/32,afed:be44::/16 " +
				"proxyProtocol.TrustedIPs:192.168.0.1 " +
				"forwardedHeaders.TrustedIPs:10.0.0.3/24,20.0.0.3/24 " +
				"auth.basic.users:test:$apr1$H6uskkkW$IgXLP6ewTrSuBkTrqE8wj/,test2:$apr1$d9hr9HBB$4HxwgUir3HP4EsggP/QNo0 " +
				"auth.digest.users:test:traefik:a2688e031edb4be6a3797f3882655c05,test2:traefik:518845800f9e2bfb1f1f740ec24f074e " +
				"auth.headerField:X-WebAuth-User " +
				"auth.forward.address:https://authserver.com/auth " +
				"auth.forward.trustForwardHeader:true " +
				"auth.forward.tls.ca:path/to/local.crt " +
				"auth.forward.tls.caOptional:true " +
				"auth.forward.tls.cert:path/to/foo.cert " +
				"auth.forward.tls.key:path/to/foo.key " +
				"auth.forward.tls.insecureSkipVerify:true ",
			expectedEntryPointName: "foo",
			expectedEntryPoint: &EntryPoint{
				Address: ":8000",
				TLS: &tls.TLS{
					Certificates: tls.Certificates{
						{
							CertFile: tls.FileOrContent("goo"),
							KeyFile:  tls.FileOrContent("gii"),
						},
					},
					ClientCA: tls.ClientCA{
						Files:    []string{"car"},
						Optional: true,
					},
				},
				Redirect: &types.Redirect{
					EntryPoint:  "https",
					Regex:       "http://localhost/(.*)",
					Replacement: "http://mydomain/$1",
					Permanent:   true,
				},
				Auth: &types.Auth{
					Basic: &types.Basic{
						Users: types.Users{
							"test:$apr1$H6uskkkW$IgXLP6ewTrSuBkTrqE8wj/",
							"test2:$apr1$d9hr9HBB$4HxwgUir3HP4EsggP/QNo0",
						},
					},
					Digest: &types.Digest{
						Users: types.Users{
							"test:traefik:a2688e031edb4be6a3797f3882655c05",
							"test2:traefik:518845800f9e2bfb1f1f740ec24f074e",
						},
					},
					Forward: &types.Forward{
						Address: "https://authserver.com/auth",
						TLS: &types.ClientTLS{
							CA:                 "path/to/local.crt",
							CAOptional:         true,
							Cert:               "path/to/foo.cert",
							Key:                "path/to/foo.key",
							InsecureSkipVerify: true,
						},
						TrustForwardHeader: true,
					},
					HeaderField: "X-WebAuth-User",
				},
				WhitelistSourceRange: []string{
					"10.42.0.0/16",
					"152.89.1.33/32",
					"afed:be44::/16",
				},
				Compress: true,
				ProxyProtocol: &ProxyProtocol{
					Insecure:   false,
					TrustedIPs: []string{"192.168.0.1"},
				},
				ForwardedHeaders: &ForwardedHeaders{
					Insecure: false,
					TrustedIPs: []string{
						"10.0.0.3/24",
						"20.0.0.3/24",
					},
				},
			},
		},
		{
			name:                   "default",
			expression:             "Name:foo",
			expectedEntryPointName: "foo",
			expectedEntryPoint: &EntryPoint{
				ForwardedHeaders: &ForwardedHeaders{Insecure: true},
			},
		},
		{
			name:                   "ForwardedHeaders insecure true",
			expression:             "Name:foo ForwardedHeaders.Insecure:true",
			expectedEntryPointName: "foo",
			expectedEntryPoint: &EntryPoint{
				ForwardedHeaders: &ForwardedHeaders{Insecure: true},
			},
		},
		{
			name:                   "ForwardedHeaders insecure false",
			expression:             "Name:foo ForwardedHeaders.Insecure:false",
			expectedEntryPointName: "foo",
			expectedEntryPoint: &EntryPoint{
				ForwardedHeaders: &ForwardedHeaders{Insecure: false},
			},
		},
		{
			name:                   "ForwardedHeaders TrustedIPs",
			expression:             "Name:foo ForwardedHeaders.TrustedIPs:10.0.0.3/24,20.0.0.3/24",
			expectedEntryPointName: "foo",
			expectedEntryPoint: &EntryPoint{
				ForwardedHeaders: &ForwardedHeaders{
					TrustedIPs: []string{"10.0.0.3/24", "20.0.0.3/24"},
				},
			},
		},
		{
			name:                   "ProxyProtocol insecure true",
			expression:             "Name:foo ProxyProtocol.Insecure:true",
			expectedEntryPointName: "foo",
			expectedEntryPoint: &EntryPoint{
				ForwardedHeaders: &ForwardedHeaders{Insecure: true},
				ProxyProtocol:    &ProxyProtocol{Insecure: true},
			},
		},
		{
			name:                   "ProxyProtocol insecure false",
			expression:             "Name:foo ProxyProtocol.Insecure:false",
			expectedEntryPointName: "foo",
			expectedEntryPoint: &EntryPoint{
				ForwardedHeaders: &ForwardedHeaders{Insecure: true},
				ProxyProtocol:    &ProxyProtocol{},
			},
		},
		{
			name:                   "ProxyProtocol TrustedIPs",
			expression:             "Name:foo ProxyProtocol.TrustedIPs:10.0.0.3/24,20.0.0.3/24",
			expectedEntryPointName: "foo",
			expectedEntryPoint: &EntryPoint{
				ForwardedHeaders: &ForwardedHeaders{Insecure: true},
				ProxyProtocol: &ProxyProtocol{
					TrustedIPs: []string{"10.0.0.3/24", "20.0.0.3/24"},
				},
			},
		},
		{
			name:                   "compress on",
			expression:             "Name:foo Compress:on",
			expectedEntryPointName: "foo",
			expectedEntryPoint: &EntryPoint{
				Compress:         true,
				ForwardedHeaders: &ForwardedHeaders{Insecure: true},
			},
		},
		{
			name:                   "compress true",
			expression:             "Name:foo Compress:true",
			expectedEntryPointName: "foo",
			expectedEntryPoint: &EntryPoint{
				Compress:         true,
				ForwardedHeaders: &ForwardedHeaders{Insecure: true},
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			eps := EntryPoints{}
			err := eps.Set(test.expression)
			require.NoError(t, err)

			ep := eps[test.expectedEntryPointName]
			assert.EqualValues(t, test.expectedEntryPoint, ep)
		})
	}
}
