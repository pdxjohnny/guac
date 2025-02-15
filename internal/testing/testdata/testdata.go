//
// Copyright 2022 The GUAC Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testdata

import (
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"reflect"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/guacsec/guac/internal/testing/keyutil"
	"github.com/guacsec/guac/pkg/assembler"
	model "github.com/guacsec/guac/pkg/assembler/clients/generated"
	asmhelpers "github.com/guacsec/guac/pkg/assembler/helpers"
	"github.com/guacsec/guac/pkg/certifier/components/root_package"
	"github.com/guacsec/guac/pkg/handler/processor"
	"github.com/guacsec/guac/pkg/ingestor/parser/common"
	"github.com/secure-systems-lab/go-securesystemslib/dsse"
)

var (
	// based off https://github.com/spdx/spdx-examples/blob/master/example7/spdx/example7-third-party-modules.spdx.json
	//go:embed exampledata/small-spdx.json
	SpdxExampleSmall []byte

	//go:embed exampledata/alpine-spdx.json
	SpdxExampleBig []byte

	//go:embed exampledata/alpine-small-spdx.json
	SpdxExampleAlpine []byte

	// Invalid types for field spdxVersion
	//go:embed exampledata/invalid-spdx.json
	SpdxInvalidExample []byte

	// Example scorecard
	//go:embed exampledata/kubernetes-scorecard.json
	ScorecardExample []byte

	// Invalid scorecard
	//go:embed exampledata/invalid-scorecard.json
	ScorecardInvalid []byte

	//go:embed exampledata/alpine-cyclonedx.json
	CycloneDXExampleAlpine []byte

	//go:embed exampledata/quarkus-deps-cyclonedx.json
	CycloneDXExampleQuarkusDeps []byte

	//go:embed exampledata/small-deps-cyclonedx.json
	CycloneDXExampleSmallDeps []byte

	//go:embed exampledata/invalid-cyclonedx.json
	CycloneDXInvalidExample []byte

	//go:embed exampledata/distroless-cyclonedx.json
	CycloneDXDistrolessExample []byte

	//go:embed exampledata/busybox-cyclonedx.json
	CycloneDXBusyboxExample []byte

	//go:embed exampledata/big-mongo-cyclonedx.json
	CycloneDXBigExample []byte

	//go:embed exampledata/npm-cyclonedx-dependencies-missing-depends-on.json
	CycloneDXDependenciesMissingDependsOn []byte

	//go:embed exampledata/laravel-cyclonedx.xml
	CycloneDXExampleLaravelXML []byte

	//go:embed exampledata/invalid-cyclonedx.xml
	CycloneDXInvalidExampleXML []byte

	//go:embed exampledata/no-dependent-components-cyclonedx.json
	CycloneDXExampleNoDependentComponents []byte

	//go:embed exampledata/crev-review.json
	ITE6CREVExample []byte

	//go:embed exampledata/github-review.json
	ITE6ReviewExample []byte

	//go:embed exampledata/certify-vuln.json
	ITE6VulnExample []byte

	//go:embed exampledata/oci-dsse-att.json
	OCIDsseAttExample []byte

	//go:embed exampledata/oci-spdx.json
	OCISPDXExample []byte

	//go:embed exampledata/go-spdx-multi-arch_1.json
	OCIGoSPDXMulti1 []byte

	//go:embed exampledata/go-spdx-multi-arch_2.json
	OCIGoSPDXMulti2 []byte

	//go:embed exampledata/go-spdx-multi-arch_3.json
	OCIGoSPDXMulti3 []byte

	// DSSE/SLSA Testdata

	// Taken from: https://slsa.dev/provenance/v0.1#example
	ite6SLSA = `
	{
		"_type": "https://in-toto.io/Statement/v0.1",
		"subject": [{"name": "helloworld", "digest": {"sha256": "3a2bd2c5cc4c978e8aefd8bd0ef335fb42ee31d1"}}],
		"predicateType": "https://slsa.dev/provenance/v0.2",
		"predicate": {
			"builder": { "id": "https://github.com/Attestations/GitHubHostedActions@v1" },
			"buildType": "https://github.com/Attestations/GitHubActionsWorkflow@v1",
			"invocation": {
			  "configSource": {
				"uri": "git+https://github.com/curl/curl-docker@master",
				"digest": { "sha1": "d6525c840a62b398424a78d792f457477135d0cf" },
				"entryPoint": "build.yaml:maketgz"
			  }
			},
			"metadata": {
			  "buildStartedOn": "2020-08-19T08:38:00Z",
			  "completeness": {
				  "environment": true
			  }
			},
			"materials": [
			  {
				"uri": "git+https://github.com/curl/curl-docker@master",
				"digest": { "sha1": "24279c5185ddc042896e3748f47fa89b48c1c14e" }
			  }, {
				"uri": "github_hosted_vm:ubuntu-18.04:20210123.1",
				"digest": { "sha1": "0bcaaa161e719bca41b6d33fc02547c0f97d5397" }
			  }
			]
		}
	}`

	b64ITE6SLSA    = base64.StdEncoding.EncodeToString([]byte(ite6SLSA))
	Ite6Payload, _ = json.Marshal(dsse.Envelope{
		PayloadType: "https://in-toto.io/Statement/v0.1",
		Payload:     b64ITE6SLSA,
		Signatures: []dsse.Signature{{
			KeyID: "id1",
			Sig:   "test",
		}},
	})
	Ite6DSSEDoc = processor.Document{
		Blob:   Ite6Payload,
		Type:   processor.DocumentDSSE,
		Format: processor.FormatJSON,
		SourceInformation: processor.SourceInformation{
			Collector: "TestCollector",
			Source:    "TestSource",
		},
	}
	Ite6SLSADoc = processor.Document{
		Blob:   []byte(ite6SLSA),
		Type:   processor.DocumentITE6SLSA,
		Format: processor.FormatJSON,
		SourceInformation: processor.SourceInformation{
			Collector: "TestCollector",
			Source:    "TestSource",
		},
	}

	art = model.ArtifactInputSpec{
		Algorithm: "sha256",
		Digest:    "3a2bd2c5cc4c978e8aefd8bd0ef335fb42ee31d1",
	}

	artPkg, _ = asmhelpers.PurlToPkg(asmhelpers.GuacGenericPurl("helloworld"))

	mat1 = model.ArtifactInputSpec{
		Algorithm: "sha1",
		Digest:    "24279c5185ddc042896e3748f47fa89b48c1c14e",
	}

	mat1Src, _ = asmhelpers.VcsToSrc("git+https://github.com/curl/curl-docker@master")

	mat2 = model.ArtifactInputSpec{
		Algorithm: "sha1",
		Digest:    "0bcaaa161e719bca41b6d33fc02547c0f97d5397",
	}

	mat2Pkg, _ = asmhelpers.PurlToPkg(asmhelpers.GuacGenericPurl("github_hosted_vm:ubuntu-18.04:20210123.1"))

	build = model.BuilderInputSpec{
		Uri: "https://github.com/Attestations/GitHubHostedActions@v1",
	}

	EcdsaPubKey, pemBytes, _ = keyutil.GetECDSAPubKey()
	// Not currently used due to skipping of DSSE and Trust information
	// keyHash, _               = dsse.SHA256KeyID(EcdsaPubKey)
	// Ident = assembler.IdentityNode{
	// 	ID:        "test",
	// 	Digest:    keyHash,
	// 	Key:       base64.StdEncoding.EncodeToString(pemBytes),
	// 	KeyType:   "ecdsa",
	// 	KeyScheme: "ecdsa",
	// 	NodeData: *assembler.NewObjectMetadata(
	// 		processor.SourceInformation{
	// 			Collector: "TestCollector",
	// 			Source:    "TestSource",
	// 		},
	// 	),
	// }

	slsaIsOccurrence = model.IsOccurrenceInputSpec{
		Justification: "from SLSA definition of checksums for subject/materials",
	}

	slsaStartTime, _ = time.Parse(time.RFC3339, "2020-08-19T08:38:00Z")
	SlsaPreds        = assembler.IngestPredicates{
		IsOccurence: []assembler.IsOccurenceIngest{
			{Pkg: artPkg, Artifact: &art, IsOccurence: &slsaIsOccurrence},
			{Src: mat1Src, Artifact: &mat1, IsOccurence: &slsaIsOccurrence},
			{Pkg: mat2Pkg, Artifact: &mat2, IsOccurence: &slsaIsOccurrence},
		},
		HasSlsa: []assembler.HasSlsaIngest{
			{
				HasSlsa: &model.SLSAInputSpec{
					BuildType:   "https://github.com/Attestations/GitHubActionsWorkflow@v1",
					SlsaVersion: "https://slsa.dev/provenance/v0.2",
					StartedOn:   slsaStartTime,
					SlsaPredicate: []model.SLSAPredicateInputSpec{
						{Key: "slsa.metadata.completeness.environment", Value: "true"},
						{Key: "slsa.metadata.buildStartedOn", Value: "2020-08-19T08:38:00Z"},
						{Key: "slsa.metadata.completeness.materials", Value: "false"},
						{Key: "slsa.buildType", Value: "https://github.com/Attestations/GitHubActionsWorkflow@v1"},
						{Key: "slsa.invocation.configSource.entryPoint", Value: "build.yaml:maketgz"},
						{Key: "slsa.invocation.configSource.uri", Value: "git+https://github.com/curl/curl-docker@master"},
						{Key: "slsa.metadata.reproducible", Value: "false"},
						{Key: "slsa.materials.0.uri", Value: "git+https://github.com/curl/curl-docker@master"},
						{Key: "slsa.builder.id", Value: "https://github.com/Attestations/GitHubHostedActions@v1"},
						{Key: "slsa.invocation.configSource.digest.sha1", Value: "d6525c840a62b398424a78d792f457477135d0cf"},
						{Key: "slsa.metadata.completeness.parameters", Value: "false"},
						{Key: "slsa.materials.0.digest.sha1", Value: "24279c5185ddc042896e3748f47fa89b48c1c14e"},
						{Key: "slsa.materials.1.uri", Value: "github_hosted_vm:ubuntu-18.04:20210123.1"},
						{Key: "slsa.materials.1.digest.sha1", Value: "0bcaaa161e719bca41b6d33fc02547c0f97d5397"},
					},
				},
				Artifact:  &art,
				Builder:   &build,
				Materials: []model.ArtifactInputSpec{mat1, mat2},
			},
		},
	}

	// TODO: needs to be resolved by https://github.com/guacsec/guac/issues/75
	Ident = []common.TrustInformation{}
	// Ident = assembler.IdentityNode{
	// 	ID:        "test",
	// 	Digest:    keyHash,
	// 	Key:       base64.StdEncoding.EncodeToString(pemBytes),
	// 	KeyType:   "ecdsa",
	// 	KeyScheme: "ecdsa",
	// 	NodeData: *assembler.NewObjectMetadata(
	// 		processor.SourceInformation{
	// 			Collector: "TestCollector",
	// 			Source:    "TestSource",
	// 		},
	// 	),
	// }

	DssePredicates = &assembler.IngestPredicates{}

	// SPDX Testdata

	topLevelPack, _       = asmhelpers.PurlToPkg("pkg:guac/oci/gcr.io/google-containers/alpine-latest")
	baselayoutPack, _     = asmhelpers.PurlToPkg("pkg:alpine/alpine-baselayout@3.2.0-r22?arch=x86_64&upstream=alpine-baselayout&distro=alpine-3.16.2")
	keysPack, _           = asmhelpers.PurlToPkg("pkg:alpine/alpine-keys@2.4-r1?arch=x86_64&upstream=alpine-keys&distro=alpine-3.16.2")
	baselayoutdataPack, _ = asmhelpers.PurlToPkg("pkg:alpine/alpine-baselayout-data@3.2.0-r22?arch=x86_64&upstream=alpine-baselayout&distro=alpine-3.16.2")

	worldFilePack, _  = asmhelpers.PurlToPkg(asmhelpers.GuacFilePurl("sha256", "713e3907167dce202d7c16034831af3d670191382a3e9026e0ac0a4023013201", strP("/etc/apk/world")))
	worldFileArtifact = &model.ArtifactInputSpec{
		Algorithm: "sha256",
		Digest:    "713e3907167dce202d7c16034831af3d670191382a3e9026e0ac0a4023013201",
	}

	rootFilePack, _  = asmhelpers.PurlToPkg(asmhelpers.GuacFilePurl("sha256", "575d810a9fae5f2f0671c9b2c0ce973e46c7207fbe5cb8d1b0d1836a6a0470e3", strP("/etc/crontabs/root")))
	rootFileArtifact = &model.ArtifactInputSpec{
		Algorithm: "sha256",
		Digest:    "575d810a9fae5f2f0671c9b2c0ce973e46c7207fbe5cb8d1b0d1836a6a0470e3",
	}

	triggersFilePack, _  = asmhelpers.PurlToPkg(asmhelpers.GuacFilePurl("sha256", "5415cfe5f88c0af38df3b7141a3f9bc6b8178e9cf72d700658091b8f5539c7b4", strP("/lib/apk/db/triggers")))
	triggersFileArtifact = &model.ArtifactInputSpec{
		Algorithm: "sha256",
		Digest:    "5415cfe5f88c0af38df3b7141a3f9bc6b8178e9cf72d700658091b8f5539c7b4",
	}

	rsaPubFilePack, _  = asmhelpers.PurlToPkg(asmhelpers.GuacFilePurl("sha256", "9a4cd858d9710963848e6d5f555325dc199d1c952b01cf6e64da2c15deedbd97", strP("/usr/share/apk/keys/alpine-devel@lists.alpinelinux.org-58cbb476.rsa.pub")))
	rsaPubFileArtifact = &model.ArtifactInputSpec{
		Algorithm: "sha256",
		Digest:    "9a4cd858d9710963848e6d5f555325dc199d1c952b01cf6e64da2c15deedbd97",
	}

	isDepJustifyTopPkg = &model.IsDependencyInputSpec{
		Justification: "top-level package GUAC heuristic connecting to each file/package",
	}
	isDepJustifyContains = &model.IsDependencyInputSpec{
		Justification: "Derived from SPDX CONTAINS relationship",
	}
	isDepJustifyDepends = &model.IsDependencyInputSpec{
		Justification: "Derived from SPDX DEPENDS_ON relationship",
	}
	isCDXDepJustifyDepends = &model.IsDependencyInputSpec{
		Justification: "CDX BOM Dependency",
	}

	isOccJustifyFile = &model.IsOccurrenceInputSpec{
		Justification: "spdx file with checksum",
	}

	isOccJustifyPkg = &model.IsOccurrenceInputSpec{
		Justification: "spdx package with checksum",
	}

	SpdxDeps = []assembler.IsDependencyIngest{
		{
			Pkg:          topLevelPack,
			DepPkg:       baselayoutPack,
			IsDependency: isDepJustifyTopPkg,
		},
		{
			Pkg:          topLevelPack,
			DepPkg:       baselayoutdataPack,
			IsDependency: isDepJustifyTopPkg,
		},
		{
			Pkg:          topLevelPack,
			DepPkg:       keysPack,
			IsDependency: isDepJustifyTopPkg,
		},
		{
			Pkg:          topLevelPack,
			DepPkg:       worldFilePack,
			IsDependency: isDepJustifyTopPkg,
		},
		{
			Pkg:          topLevelPack,
			DepPkg:       rootFilePack,
			IsDependency: isDepJustifyTopPkg,
		},
		{
			Pkg:          topLevelPack,
			DepPkg:       triggersFilePack,
			IsDependency: isDepJustifyTopPkg,
		},
		{
			Pkg:          topLevelPack,
			DepPkg:       rsaPubFilePack,
			IsDependency: isDepJustifyTopPkg,
		},
		{
			Pkg:          baselayoutPack,
			DepPkg:       keysPack,
			IsDependency: isDepJustifyDepends,
		},
		{
			Pkg:          rootFilePack,
			DepPkg:       rsaPubFilePack,
			IsDependency: isDepJustifyDepends,
		},
		{
			Pkg:          baselayoutPack,
			DepPkg:       rootFilePack,
			IsDependency: isDepJustifyContains,
		},
		{
			Pkg:          keysPack,
			DepPkg:       rsaPubFilePack,
			IsDependency: isDepJustifyContains,
		},
	}

	SpdxOccurences = []assembler.IsOccurenceIngest{
		{
			Pkg:         worldFilePack,
			Artifact:    worldFileArtifact,
			IsOccurence: isOccJustifyFile,
		},
		{
			Pkg:         rootFilePack,
			Artifact:    rootFileArtifact,
			IsOccurence: isOccJustifyFile,
		},
		{
			Pkg:         rsaPubFilePack,
			Artifact:    rsaPubFileArtifact,
			IsOccurence: isOccJustifyFile,
		},
		{
			Pkg:         triggersFilePack,
			Artifact:    triggersFileArtifact,
			IsOccurence: isOccJustifyFile,
		},
	}

	SpdxIngestionPredicates = assembler.IngestPredicates{
		IsDependency: SpdxDeps,
		IsOccurence:  SpdxOccurences,
	}

	// CycloneDX Testdata
	cdxTopLevelPack, _ = asmhelpers.PurlToPkg("pkg:oci/static@sha256:6ad5b696af3ca05a048bd29bf0f623040462638cb0b29c8d702cbb2805687388?repository_url=gcr.io/distroless/static&tag=nonroot")

	cdxTzdataPack, _ = asmhelpers.PurlToPkg("pkg:deb/debian/tzdata@2021a-1+deb11u6?arch=all&distro=debian-11")

	cdxNetbasePack, _ = asmhelpers.PurlToPkg("pkg:deb/debian/netbase@6.3?arch=all&distro=debian-11")

	cdxBasefilesPack, _ = asmhelpers.PurlToPkg("pkg:deb/debian/base-files@11.1+deb11u5?arch=amd64&distro=debian-11")

	CdxDeps = []assembler.IsDependencyIngest{
		{
			Pkg:          cdxTopLevelPack,
			DepPkg:       cdxBasefilesPack,
			IsDependency: isDepJustifyTopPkg,
		},
		{
			Pkg:          cdxTopLevelPack,
			DepPkg:       cdxNetbasePack,
			IsDependency: isDepJustifyTopPkg,
		},
		{
			Pkg:          cdxTopLevelPack,
			DepPkg:       cdxTzdataPack,
			IsDependency: isDepJustifyTopPkg,
		},
	}

	CdxIngestionPredicates = assembler.IngestPredicates{
		IsDependency: CdxDeps,
	}

	cdxTopQuarkusPack, _ = asmhelpers.PurlToPkg("pkg:maven/org.acme/getting-started@1.0.0-SNAPSHOT?type=jar")

	cdxResteasyPack, _ = asmhelpers.PurlToPkg("pkg:maven/io.quarkus/quarkus-resteasy-reactive@2.13.4.Final?type=jar")

	cdxReactiveCommonPack, _ = asmhelpers.PurlToPkg("pkg:maven/io.quarkus/quarkus-resteasy-reactive-common@2.13.4.Final?type=jar")

	CdxQuarkusDeps = []assembler.IsDependencyIngest{
		{
			Pkg:          cdxTopQuarkusPack,
			DepPkg:       cdxResteasyPack,
			IsDependency: isDepJustifyTopPkg,
		},
		{
			Pkg:          cdxTopQuarkusPack,
			DepPkg:       cdxReactiveCommonPack,
			IsDependency: isDepJustifyTopPkg,
		},
		{
			Pkg:          cdxResteasyPack,
			DepPkg:       cdxReactiveCommonPack,
			IsDependency: isCDXDepJustifyDepends,
		},
	}

	CdxQuarkusIngestionPredicates = assembler.IngestPredicates{
		IsDependency: CdxQuarkusDeps,
	}

	cdxWebAppPackage, _ = asmhelpers.PurlToPkg("pkg:npm/web-app@1.0.0")

	cdxBootstrapPackage, _ = asmhelpers.PurlToPkg("pkg:npm/bootstrap@4.0.0-beta.2")

	CdxNpmDeps = []assembler.IsDependencyIngest{
		{
			Pkg:          cdxWebAppPackage,
			DepPkg:       cdxBootstrapPackage,
			IsDependency: isDepJustifyTopPkg,
		},
	}

	CdxNpmIngestionPredicates = assembler.IngestPredicates{
		IsDependency: CdxNpmDeps,
	}

	CdxEmptyIngestionPredicates = assembler.IngestPredicates{}

	// ceritifer testdata

	Text4ShellVulAttestation = `{
		"_type":"https://in-toto.io/Statement/v0.1",
		"predicateType":"https://in-toto.io/attestation/vuln/v0.1",
		"subject":[
		   {
			  "name":"pkg:maven/org.apache.commons/commons-text@1.9",
			  "digest":null
		   }
		],
		"predicate":{
		   "invocation":{
			  "uri":"guac",
			  "producer_id":"guacsec/guac"
		   },
		   "scanner":{
			  "uri":"osv.dev",
			  "version":"0.0.14",
			  "db":{
			  },
			  "result":[
				 {
					"vulnerability_id":"GHSA-599f-7c49-w659"
				 }
			  ]
		   },
		   "metadata":{
			  "scannedOn":"2022-11-22T13:18:58.063182-05:00"
		   }
		}
	 }`
	SecondLevelVulAttestation = `{
		"_type":"https://in-toto.io/Statement/v0.1",
		"predicateType":"https://in-toto.io/attestation/vuln/v0.1",
		"subject":[
		   {
			  "name":"pkg:oci/vul-secondLevel-latest?repository_url=grc.io",
			  "digest":{"sha256":"fe608dbc4894fc0b9c82908ece9ddddb63bb79083e5b25f2c02f87773bde1aa1"}
		   }
		],
		"predicate":{
		   "invocation":{
			  "uri":"guac",
			  "producer_id":"guacsec/guac"
		   },
		   "scanner": {
			"uri": "osv.dev",
			"version": "0.0.14",
			"db": {}
		   },
		   "metadata":{
			  "scannedOn":"2022-11-22T13:19:18.825699-05:00"
		   }
		}
	 }`
	RootVulAttestation = `{
		"_type":"https://in-toto.io/Statement/v0.1",
		"predicateType":"https://in-toto.io/attestation/vuln/v0.1",
		"subject":[
		   {
			  "name":"pkg:oci/vul-image-latest?repository_url=grc.io",
			  "digest":null
		   }
		],
		"predicate":{
		   "invocation":{
			  "uri":"guac",
			  "producer_id":"guacsec/guac"
		   },
		   "scanner": {
			"uri": "osv.dev",
			"version": "0.0.14",
			"db": {}
		   },
		   "metadata":{
			  "scannedOn":"2022-11-22T13:19:18.825699-05:00"
		   }
		}
	 }`
	Log4JVulAttestation = `{
		"_type":"https://in-toto.io/Statement/v0.1",
		"predicateType":"https://in-toto.io/attestation/vuln/v0.1",
		"subject":[
		   {
			  "name":"pkg:maven/org.apache.logging.log4j/log4j-core@2.8.1",
			  "digest":null
		   }
		],
		"predicate":{
		   "invocation":{
			  "uri":"guac",
			  "producer_id":"guacsec/guac"
		   },
		   "scanner":{
			  "uri":"osv.dev",
			  "version":"0.0.14",
			  "db":{
			  },
			  "result":[
				 {
					"vulnerability_id":"GHSA-7rjr-3q55-vv33"
				 },
				 {
					"vulnerability_id":"GHSA-8489-44mv-ggj8"
				 },
				 {
					"vulnerability_id":"GHSA-fxph-q3j8-mv87"
				 },
				 {
					"vulnerability_id":"GHSA-jfh8-c2jp-5v3q"
				 },
				 {
					"vulnerability_id":"GHSA-p6xc-xr62-6r2g"
				 },
				 {
					"vulnerability_id":"GHSA-vwqq-5vrc-xw9h"
				 }
			  ]
		   },
		   "metadata":{
			  "scannedOn":"2022-11-22T13:18:31.607996-05:00"
		   }
		}
	 }`

	RootPackage = root_package.PackageNode{
		Purl: "pkg:oci/vul-image-latest?repository_url=grc.io",
	}

	SecondLevelPackage = root_package.PackageNode{
		Purl:      "pkg:oci/vul-secondLevel-latest?repository_url=grc.io",
		Algorithm: "sha256",
		Digest:    "fe608dbc4894fc0b9c82908ece9ddddb63bb79083e5b25f2c02f87773bde1aa1",
	}

	Log4JPackage = root_package.PackageNode{
		Purl: "pkg:maven/org.apache.logging.log4j/log4j-core@2.8.1",
	}

	Text4ShelPackage = root_package.PackageNode{
		Purl: "pkg:maven/org.apache.commons/commons-text@1.9",
	}

	VertxWebCommonAttestation = `{
		"_type": "https://in-toto.io/Statement/v0.1",
		"predicateType": "https://in-toto.io/attestation/vuln/v0.1",
		"subject": [
			{
				"name": "pkg:maven/io.vertx/vertx-web-common@4.3.7?type=jar",
				"digest": null
			}
		],
		"predicate": {
			"invocation": {
				"uri": "guac",
				"producer_id": "guacsec/guac"
			},
			"scanner": {
				"uri": "osv.dev",
				"version": "0.0.14",
				"db": {}
			},
			"metadata": {
				"scannedOn": "2023-02-15T11:10:08.986308-08:00"
			}
		}
	}`

	VertxAuthCommonAttestation = `{
		"_type": "https://in-toto.io/Statement/v0.1",
		"predicateType": "https://in-toto.io/attestation/vuln/v0.1",
		"subject": [
			{
				"name": "pkg:maven/io.vertx/vertx-auth-common@4.3.7?type=jar",
				"digest": null
			}
		],
		"predicate": {
			"invocation": {
				"uri": "guac",
				"producer_id": "guacsec/guac"
			},
			"scanner": {
				"uri": "osv.dev",
				"version": "0.0.14",
				"db": {}
			},
			"metadata": {
				"scannedOn": "2023-02-15T11:10:08.986401-08:00"
			}
		}
	}`

	VertxBridgeCommonAttestation = `{
		"_type": "https://in-toto.io/Statement/v0.1",
		"predicateType": "https://in-toto.io/attestation/vuln/v0.1",
		"subject": [
			{
				"name": "pkg:maven/io.vertx/vertx-bridge-common@4.3.7?type=jar",
				"digest": null
			}
		],
		"predicate": {
			"invocation": {
				"uri": "guac",
				"producer_id": "guacsec/guac"
			},
			"scanner": {
				"uri": "osv.dev",
				"version": "0.0.14",
				"db": {}
			},
			"metadata": {
				"scannedOn": "2023-02-15T11:10:08.98646-08:00"
			}
		}
	}`

	VertxCoreCommonAttestation = `{
		"_type": "https://in-toto.io/Statement/v0.1",
		"predicateType": "https://in-toto.io/attestation/vuln/v0.1",
		"subject": [
			{
				"name": "pkg:maven/io.vertx/vertx-core@4.3.7?type=jar",
				"digest": null
			}
		],
		"predicate": {
			"invocation": {
				"uri": "guac",
				"producer_id": "guacsec/guac"
			},
			"scanner": {
				"uri": "osv.dev",
				"version": "0.0.14",
				"db": {}
			},
			"metadata": {
				"scannedOn": "2023-02-15T11:10:08.986506-08:00"
			}
		}
	}`

	VertxWebAttestation = `{
		"_type": "https://in-toto.io/Statement/v0.1",
		"predicateType": "https://in-toto.io/attestation/vuln/v0.1",
		"subject": [
			{
				"name": "pkg:maven/io.vertx/vertx-web@4.3.7?type=jar",
				"digest": null
			}
		],
		"predicate": {
			"invocation": {
				"uri": "guac",
				"producer_id": "guacsec/guac"
			},
			"scanner": {
				"uri": "osv.dev",
				"version": "0.0.14",
				"db": {},
				"result": [
					{
						"vulnerability_id": "GHSA-53jx-vvf9-4x38"
					}
				]
			},
			"metadata": {
				"scannedOn": "2023-02-15T11:10:08.986592-08:00"
			}
		}
	}`

	VertxWebCommonPackage = root_package.PackageNode{
		Purl: "pkg:maven/io.vertx/vertx-web-common@4.3.7?type=jar",
	}

	VertxAuthCommonPackage = root_package.PackageNode{
		Purl: "pkg:maven/io.vertx/vertx-auth-common@4.3.7?type=jar",
	}

	VertxBridgeCommonPackage = root_package.PackageNode{
		Purl: "pkg:maven/io.vertx/vertx-bridge-common@4.3.7?type=jar",
	}

	VertxCoreCommonPackage = root_package.PackageNode{
		Purl: "pkg:maven/io.vertx/vertx-core@4.3.7?type=jar",
	}

	VertxWebPackage = root_package.PackageNode{
		Purl: "pkg:maven/io.vertx/vertx-web@4.3.7?type=jar",
	}

	// Deps.dev

	CollectedForeignTypes = `{
	"CurrentPackage":{
	   "type":"cargo",
	   "namespace":"",
	   "name":"foreign-types",
	   "version":"0.3.2",
	   "qualifiers":null,
	   "subpath":""
	},
	"Source":{
	   "type":"git",
	   "namespace":"github.com/sfackler",
	   "name":"foreign-types",
	   "tag":null,
	   "commit":null
	},
	"Vulnerabilities":null,
	"Scorecard":{
	   "checks":[
		  {
			 "check":"Code-Review",
			 "score":2
		  },
		  {
			 "check":"Maintained",
			 "score":0
		  },
		  {
			 "check":"CII-Best-Practices",
			 "score":0
		  },
		  {
			 "check":"Vulnerabilities",
			 "score":10
		  },
		  {
			 "check":"Signed-Releases",
			 "score":-1
		  },
		  {
			 "check":"Branch-Protection",
			 "score":0
		  },
		  {
			 "check":"License",
			 "score":10
		  },
		  {
			 "check":"Pinned-Dependencies",
			 "score":7
		  },
		  {
			 "check":"Binary-Artifacts",
			 "score":10
		  },
		  {
			 "check":"Token-Permissions",
			 "score":0
		  },
		  {
			 "check":"Dangerous-Workflow",
			 "score":10
		  },
		  {
			 "check":"SAST",
			 "score":0
		  },
		  {
			 "check":"Packaging",
			 "score":-1
		  },
		  {
			 "check":"Dependency-Update-Tool",
			 "score":0
		  },
		  {
			 "check":"Fuzzing",
			 "score":0
		  },
		  {
			 "check":"Security-Policy",
			 "score":0
		  }
	   ],
	   "aggregateScore":3.700000047683716,
	   "timeScanned":"2022-11-21T17:45:50.52Z",
	   "scorecardVersion":"v4.8.0-78-gfb07860",
	   "scorecardCommit":"fb07860d86065cdcbd2d0d5c6b998ff4542d53fe",
	   "origin":"deps.dev",
	   "collector":"deps.dev"
	},
	"DepPackages":[
	   {
		  "CurrentPackage":{
			 "type":"cargo",
			 "namespace":"",
			 "name":"foreign-types-shared",
			 "version":"0.1.1",
			 "qualifiers":[
				
			 ],
			 "subpath":""
		  },
		  "Source":{
			 "type":"git",
			 "namespace":"github.com/sfackler",
			 "name":"foreign-types",
			 "tag":null,
			 "commit":null
		  },
		  "Vulnerabilities":null,
		  "Scorecard":{
			 "checks":[
				{
				   "check":"Code-Review",
				   "score":2
				},
				{
				   "check":"Maintained",
				   "score":0
				},
				{
				   "check":"CII-Best-Practices",
				   "score":0
				},
				{
				   "check":"Vulnerabilities",
				   "score":10
				},
				{
				   "check":"Signed-Releases",
				   "score":-1
				},
				{
				   "check":"Branch-Protection",
				   "score":0
				},
				{
				   "check":"License",
				   "score":10
				},
				{
				   "check":"Pinned-Dependencies",
				   "score":7
				},
				{
				   "check":"Binary-Artifacts",
				   "score":10
				},
				{
				   "check":"Token-Permissions",
				   "score":0
				},
				{
				   "check":"Dangerous-Workflow",
				   "score":10
				},
				{
				   "check":"SAST",
				   "score":0
				},
				{
				   "check":"Packaging",
				   "score":-1
				},
				{
				   "check":"Dependency-Update-Tool",
				   "score":0
				},
				{
				   "check":"Fuzzing",
				   "score":0
				},
				{
				   "check":"Security-Policy",
				   "score":0
				}
			 ],
			 "aggregateScore":3.700000047683716,
			 "timeScanned":"2022-11-21T17:45:50.52Z",
			 "scorecardVersion":"v4.8.0-78-gfb07860",
			 "scorecardCommit":"fb07860d86065cdcbd2d0d5c6b998ff4542d53fe",
			 "origin":"deps.dev",
			 "collector":"deps.dev"
		  },
		  "DepPackages":null,
		  "UpdateTime":"2022-11-21T17:45:50.52Z"
	   }
	],
	"UpdateTime":"2022-11-21T17:45:50.52Z"
 }`
	CollectedYargsParser = `
 {
	 "CurrentPackage":{
		"type":"npm",
		"namespace":"",
		"name":"yargs-parser",
		"version":"4.2.1",
		"qualifiers":null,
		"subpath":""
	 },
	 "Source":{
		"type":"git",
		"namespace":"github.com/yargs",
		"name":"yargs-parser.git",
		"tag":null,
		"commit":null
	 },
	 "Vulnerabilities":[
		{
		   "osvId":"GHSA-p9pc-299p-vxgp"
		}
	 ],
	 "Scorecard":null,
	 "DepPackages":[
		{
		   "CurrentPackage":{
			  "type":"npm",
			  "namespace":"",
			  "name":"camelcase",
			  "version":"3.0.0",
			  "qualifiers":[
				 
			  ],
			  "subpath":""
		   },
		   "Source":{
			  "type":"git",
			  "namespace":"github.com/sindresorhus",
			  "name":"camelcase.git",
			  "tag":null,
			  "commit":null
		   },
		   "Vulnerabilities":null,
		   "Scorecard":null,
		   "DepPackages":null,
		   "UpdateTime":"2022-11-21T17:45:50.52Z"
		}
	 ],
	 "UpdateTime":"2022-11-21T17:45:50.52Z"
  }`
)

func GuacNodeSliceEqual(slice1, slice2 []assembler.GuacNode) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	result := true

	for _, node1 := range slice1 {
		e := false
		for _, node2 := range slice2 {
			if node1.Type() == "Package" && node2.Type() == "Package" {
				if node1.(assembler.PackageNode).Name == node2.(assembler.PackageNode).Name {
					if reflect.DeepEqual(node1, node2) {
						e = true
						break
					}
				}
			} else if node1.Type() == "Artifact" && node2.Type() == "Artifact" {
				if node1.(assembler.ArtifactNode).Name == node2.(assembler.ArtifactNode).Name {
					if reflect.DeepEqual(node1, node2) {
						e = true
						break
					}
				}
			} else if node1.Type() == "Attestation" && node2.Type() == "Attestation" {
				if node1.(assembler.AttestationNode).FilePath == node2.(assembler.AttestationNode).FilePath {
					if reflect.DeepEqual(node1, node2) {
						e = true
						break
					}
				}
			} else if node1.Type() == "Builder" && node2.Type() == "Builder" {
				if node1.(assembler.BuilderNode).BuilderId == node2.(assembler.BuilderNode).BuilderId {
					if reflect.DeepEqual(node1, node2) {
						e = true
						break
					}
				}
			} else if node1.Type() == "Identity" && node2.Type() == "Identity" {
				if node1.(assembler.IdentityNode).ID == node2.(assembler.IdentityNode).ID {
					if reflect.DeepEqual(node1, node2) {
						e = true
						break
					}
				}
			} else if node1.Type() == "Vulnerability" && node2.Type() == "Vulnerability" {
				if node1.(assembler.VulnerabilityNode).ID == node2.(assembler.VulnerabilityNode).ID {
					if reflect.DeepEqual(node1, node2) {
						e = true
						break
					}
				}
			}
		}
		if !e {
			result = false
			break
		}
	}
	return result
}

func GuacEdgeSliceEqual(slice1, slice2 []assembler.GuacEdge) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	result := true
	for _, edge1 := range slice1 {
		e := false
		for _, edge2 := range slice2 {
			if edge1.Type() == "DependsOn" && edge2.Type() == "DependsOn" {
				if reflect.DeepEqual(edge1, edge2) {
					e = true
					break
				}
			} else if edge1.Type() == "Contains" && edge2.Type() == "Contains" {
				if reflect.DeepEqual(edge1, edge2) {
					e = true
					break
				}
			} else if edge1.Type() == "Attestation" && edge2.Type() == "Attestation" {
				if reflect.DeepEqual(edge1, edge2) {
					e = true
					break
				}
			} else if edge1.Type() == "Identity" && edge2.Type() == "Identity" {
				if reflect.DeepEqual(edge1, edge2) {
					e = true
					break
				}
			} else if edge1.Type() == "BuiltBy" && edge2.Type() == "BuiltBy" {
				if reflect.DeepEqual(edge1, edge2) {
					e = true
					break
				}
			} else if edge1.Type() == "Vulnerable" && edge2.Type() == "Vulnerable" {
				if reflect.DeepEqual(edge1, edge2) {
					e = true
					break
				}
			}
		}
		if !e {
			result = false
			break
		}
	}
	return result
}

var IngestPredicatesCmpOpts = []cmp.Option{
	cmpopts.EquateEmpty(),
	cmpopts.SortSlices(certifyScorecardLess),
	cmpopts.SortSlices(isDependencyLess),
	cmpopts.SortSlices(isOccurenceLess),
	cmpopts.SortSlices(packageQualifierInputSpecLess),
	cmpopts.SortSlices(psaInputSpecLess),
	cmpopts.SortSlices(slsaPredicateInputSpecLess),
}

func certifyScorecardLess(e1, e2 assembler.CertifyScorecardIngest) bool {
	return gLess(e1, e2)
}

func isDependencyLess(e1, e2 assembler.IsDependencyIngest) bool {
	return gLess(e1, e2)
}

func isOccurenceLess(e1, e2 assembler.IsOccurenceIngest) bool {
	return gLess(e1, e2)
}

func packageQualifierInputSpecLess(e1, e2 model.PackageQualifierInputSpec) bool {
	return gLess(e1, e2)
}

func psaInputSpecLess(e1, e2 model.ArtifactInputSpec) bool {
	return gLess(e1, e2)
}

func slsaPredicateInputSpecLess(e1, e2 model.SLSAPredicateInputSpec) bool {
	return gLess(e1, e2)
}

func gLess(e1, e2 any) bool {
	s1, _ := json.Marshal(e1)
	s2, _ := json.Marshal(e2)
	return string(s1) < string(s2)
}

func strP(s string) *string {
	return &s
}
