/*
Copyright 2020 The Tekton Authors
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

package artifacts

import (
	"fmt"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/go-containerregistry/pkg/name"
	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
	"github.com/tektoncd/chains/pkg/chains/objects"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logtesting "knative.dev/pkg/logging/testing"
)

const (
	digest1                 = "sha256:05f95b26ed10668b7183c1e2da98610e91372fa9f510046d4ce5812addad86b5"
	digest2                 = "sha256:05f95b26ed10668b7183c1e2da98610e91372fa9f510046d4ce5812addad86b6"
	digest3                 = "sha256:05f95b26ed10668b7183c1e2da98610e91372fa9f510046d4ce5812addad86b7"
	digest4                 = "sha256:05f95b26ed10668b7183c1e2da98610e91372fa9f510046d4ce5812addad86b8"
	digest_sha384           = "sha384:05f95b26ed10668b7183c1e2da98610e91372fa9f510046d4ce5812addad86b893c56eeba9ec70f74c9bfd297d951664"
	digest_sha512           = "sha512:05f95b26ed10668b7183c1e2da98610e91372fa9f510046d4ce5812addad86b805f95b26ed10668b7183c1e2da98610e91372fa9f510046d4ce5812addad86b8"
	digest_sha1             = "sha1:93c56eeba9ec70f74c9bfd297d9516642d366cb5"
	digest_incorrect_sha1   = "sha1:93c56eeba9ec70f74c9bfd297d9516642d366c5"
	digest_incorrect_sha512 = "sha512:05f95b26ed1066b7183c1e2da98610e91372fa9f510046d4ce5812addad86b805f95b26ed10668b7183c1e2da98610e91372fa9f510046d4ce5812addad86b8"
	digest_incorrect_sha384 = "sha384:0595b26ed10668b7183c1e2da98610e91372fa9f510046d4ce5812addad86b893c56eeba9ec70f74c9bfd297d951664"
)

var ignore = []cmp.Option{cmpopts.IgnoreUnexported(name.Registry{}, name.Repository{}, name.Digest{})}

func TestOCIArtifact_ExtractObjects(t *testing.T) {

	tests := []struct {
		name string
		obj  objects.TektonObject
		want []interface{}
	}{
		{
			name: "one image",
			obj: objects.NewTaskRunObject(&v1beta1.TaskRun{
				TypeMeta: metav1.TypeMeta{
					Kind: "TaskRun",
				},
				Status: v1beta1.TaskRunStatus{
					TaskRunStatusFields: v1beta1.TaskRunStatusFields{
						ResourcesResult: []v1beta1.PipelineResourceResult{
							{
								ResourceName: "my-image",
								Key:          "url",
								Value:        "gcr.io/foo/bar",
							},
							{
								ResourceName: "my-image",
								Key:          "digest",
								Value:        digest1,
							},
						},
						TaskSpec: &v1beta1.TaskSpec{
							Resources: &v1beta1.TaskResources{
								Outputs: []v1beta1.TaskResource{
									{
										ResourceDeclaration: v1beta1.ResourceDeclaration{
											Name: "my-image",
											Type: "image",
										},
									},
								},
							},
						},
					},
				},
			}),
			want: []interface{}{createDigest(t, "gcr.io/foo/bar@sha256:05f95b26ed10668b7183c1e2da98610e91372fa9f510046d4ce5812addad86b5")},
		},
		{
			name: "two images",
			obj: objects.NewTaskRunObject(&v1beta1.TaskRun{
				TypeMeta: metav1.TypeMeta{
					Kind: "TaskRun",
				},
				Status: v1beta1.TaskRunStatus{
					TaskRunStatusFields: v1beta1.TaskRunStatusFields{
						ResourcesResult: []v1beta1.PipelineResourceResult{
							{
								ResourceName: "my-image1",
								Key:          "url",
								Value:        "gcr.io/foo/bar",
							},
							{
								ResourceName: "my-image1",
								Key:          "digest",
								Value:        digest1,
							},
							{
								ResourceName: "my-image2",
								Key:          "url",
								Value:        "gcr.io/foo/baz",
							},
							{
								ResourceName: "my-image2",
								Key:          "digest",
								Value:        digest2,
							},
						},
						TaskSpec: &v1beta1.TaskSpec{
							Resources: &v1beta1.TaskResources{
								Outputs: []v1beta1.TaskResource{
									{
										ResourceDeclaration: v1beta1.ResourceDeclaration{
											Name: "my-image1",
											Type: "image",
										},
									},
									{
										ResourceDeclaration: v1beta1.ResourceDeclaration{
											Name: "my-image2",
											Type: "image",
										},
									},
								},
							},
						},
					},
				},
			}),
			want: []interface{}{
				createDigest(t, "gcr.io/foo/bar@sha256:05f95b26ed10668b7183c1e2da98610e91372fa9f510046d4ce5812addad86b5"),
				createDigest(t, "gcr.io/foo/baz@sha256:05f95b26ed10668b7183c1e2da98610e91372fa9f510046d4ce5812addad86b6"),
			},
		},
		{
			name: "resource and result",
			obj: objects.NewTaskRunObject(&v1beta1.TaskRun{
				TypeMeta: metav1.TypeMeta{
					Kind: "TaskRun",
				},
				Status: v1beta1.TaskRunStatus{
					TaskRunStatusFields: v1beta1.TaskRunStatusFields{
						ResourcesResult: []v1beta1.PipelineResourceResult{
							{
								ResourceName: "my-image",
								Key:          "url",
								Value:        "gcr.io/foo/bar",
							},
							{
								ResourceName: "my-image",
								Key:          "digest",
								Value:        digest1,
							},
						},
						TaskRunResults: []v1beta1.TaskRunResult{
							{
								Name:  "IMAGE_URL",
								Value: *v1beta1.NewArrayOrString("gcr.io/foo/bat"),
							},
							{
								Name:  "IMAGE_DIGEST",
								Value: *v1beta1.NewArrayOrString("sha256:05f95b26ed10668b7183c1e2da98610e91372fa9f510046d4ce5812addad86b4"),
							},
						},
						TaskSpec: &v1beta1.TaskSpec{
							Results: []v1beta1.TaskResult{
								{
									Name: "IMAGE_URL",
								},
								{
									Name: "IMAGE_DIGEST",
								},
							},
							Resources: &v1beta1.TaskResources{
								Outputs: []v1beta1.TaskResource{
									{
										ResourceDeclaration: v1beta1.ResourceDeclaration{
											Name: "my-image",
											Type: "image",
										},
									},
								},
							},
						},
					},
				},
			}),
			want: []interface{}{
				createDigest(t, "gcr.io/foo/bat@sha256:05f95b26ed10668b7183c1e2da98610e91372fa9f510046d4ce5812addad86b4"),
				createDigest(t, "gcr.io/foo/bar@sha256:05f95b26ed10668b7183c1e2da98610e91372fa9f510046d4ce5812addad86b5")},
		},
		{
			name: "extra",
			obj: objects.NewTaskRunObject(&v1beta1.TaskRun{
				TypeMeta: metav1.TypeMeta{
					Kind: "TaskRun",
				},
				Status: v1beta1.TaskRunStatus{
					TaskRunStatusFields: v1beta1.TaskRunStatusFields{
						TaskRunResults: []v1beta1.TaskRunResult{
							{
								Name:  "IMAGE_URL",
								Value: *v1beta1.NewArrayOrString("foo"),
							},
							{
								Name:  "gibberish",
								Value: *v1beta1.NewArrayOrString("baz"),
							},
						},
						ResourcesResult: []v1beta1.PipelineResourceResult{
							{
								ResourceName: "my-image",
								Key:          "url",
								Value:        "gcr.io/foo/bar",
							},
							{
								ResourceName: "my-image",
								Key:          "digest",
								Value:        digest1,
							},
							{
								ResourceName: "gibberish",
								Key:          "url",
								Value:        "gcr.io/foo/bar",
							},
							{
								ResourceName: "gobble-dygook",
								Key:          "digest",
								Value:        digest1,
							},
						},
						TaskSpec: &v1beta1.TaskSpec{
							Resources: &v1beta1.TaskResources{
								Outputs: []v1beta1.TaskResource{
									{
										ResourceDeclaration: v1beta1.ResourceDeclaration{
											Name: "my-image",
											Type: "image",
										},
									},
								},
							},
						},
					},
				},
			}),
			want: []interface{}{createDigest(t, "gcr.io/foo/bar@sha256:05f95b26ed10668b7183c1e2da98610e91372fa9f510046d4ce5812addad86b5")},
		}, {
			name: "images",
			obj: objects.NewTaskRunObject(&v1beta1.TaskRun{
				Status: v1beta1.TaskRunStatus{
					TaskRunStatusFields: v1beta1.TaskRunStatusFields{
						TaskRunResults: []v1beta1.TaskRunResult{
							{
								Name:  "IMAGES",
								Value: *v1beta1.NewArrayOrString(fmt.Sprintf("  \n \tgcr.io/foo/bar@%s\n,gcr.io/baz/bar@%s", digest1, digest2)),
							},
						},
					},
				},
			}),
			want: []interface{}{
				createDigest(t, "gcr.io/foo/bar@sha256:05f95b26ed10668b7183c1e2da98610e91372fa9f510046d4ce5812addad86b5"),
				createDigest(t, "gcr.io/baz/bar@sha256:05f95b26ed10668b7183c1e2da98610e91372fa9f510046d4ce5812addad86b6"),
			},
		}, {
			name: "images-newline",
			obj: objects.NewTaskRunObject(&v1beta1.TaskRun{
				Status: v1beta1.TaskRunStatus{
					TaskRunStatusFields: v1beta1.TaskRunStatusFields{
						TaskRunResults: []v1beta1.TaskRunResult{
							{
								Name:  "IMAGES",
								Value: *v1beta1.NewArrayOrString(fmt.Sprintf("gcr.io/foo/bar@%s\ngcr.io/baz/bar@%s\n\n", digest1, digest2)),
							},
						},
					},
				},
			}),
			want: []interface{}{
				createDigest(t, "gcr.io/foo/bar@sha256:05f95b26ed10668b7183c1e2da98610e91372fa9f510046d4ce5812addad86b5"),
				createDigest(t, "gcr.io/baz/bar@sha256:05f95b26ed10668b7183c1e2da98610e91372fa9f510046d4ce5812addad86b6"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := logtesting.TestLogger(t)
			oa := &OCIArtifact{
				Logger: logger,
			}
			got := oa.ExtractObjects(tt.obj)
			sort.Slice(got, func(i, j int) bool {
				a := got[i].(name.Digest)
				b := got[j].(name.Digest)
				return a.DigestStr() < b.DigestStr()
			})
			if !cmp.Equal(got, tt.want, ignore...) {
				t.Errorf("OCIArtifact.ExtractObjects() = %s", cmp.Diff(got, tt.want, ignore...))
			}
		})
	}
}

func TestExtractOCIImagesFromResults(t *testing.T) {
	tr := &v1beta1.TaskRun{
		Status: v1beta1.TaskRunStatus{
			TaskRunStatusFields: v1beta1.TaskRunStatusFields{
				TaskRunResults: []v1beta1.TaskRunResult{
					{Name: "img1_IMAGE_URL", Value: *v1beta1.NewArrayOrString("img1")},
					{Name: "img1_IMAGE_DIGEST", Value: *v1beta1.NewArrayOrString(digest1)},
					{Name: "img2_IMAGE_URL", Value: *v1beta1.NewArrayOrString("img2")},
					{Name: "img2_IMAGE_DIGEST", Value: *v1beta1.NewArrayOrString(digest2)},
					{Name: "IMAGE_URL", Value: *v1beta1.NewArrayOrString("img3")},
					{Name: "IMAGE_DIGEST", Value: *v1beta1.NewArrayOrString(digest1)},
					{Name: "img4_IMAGE_URL", Value: *v1beta1.NewArrayOrString("img4")},
					{Name: "img5_IMAGE_DIGEST", Value: *v1beta1.NewArrayOrString("sha123:abc")},
					{Name: "empty_str_IMAGE_DIGEST", Value: *v1beta1.NewArrayOrString("")},
					{Name: "empty_str_IMAGE_URL", Value: *v1beta1.NewArrayOrString("")},
				},
			},
		},
	}
	obj := objects.NewTaskRunObject(tr)
	want := []interface{}{
		createDigest(t, fmt.Sprintf("img1@%s", digest1)),
		createDigest(t, fmt.Sprintf("img2@%s", digest2)),
		createDigest(t, fmt.Sprintf("img3@%s", digest1)),
	}
	got := ExtractOCIImagesFromResults(obj, logtesting.TestLogger(t))
	sort.Slice(got, func(i, j int) bool {
		a := got[i].(name.Digest)
		b := got[j].(name.Digest)
		return a.String() < b.String()
	})
	if !cmp.Equal(got, want, ignore...) {
		t.Fatalf("not the same %s", cmp.Diff(want, got, ignore...))
	}
}

func TestExtractSignableTargetFromResults(t *testing.T) {
	tr := &v1beta1.TaskRun{
		Status: v1beta1.TaskRunStatus{
			TaskRunStatusFields: v1beta1.TaskRunStatusFields{
				TaskRunResults: []v1beta1.TaskRunResult{
					{Name: "mvn1_ARTIFACT_URI", Value: *v1beta1.NewArrayOrString("projects/test-project/locations/us-west4/repositories/test-repo/mavenArtifacts/com.google.guava:guava:31.0-jre")},
					{Name: "mvn1_ARTIFACT_DIGEST", Value: *v1beta1.NewArrayOrString(digest1)},
					{Name: "mvn1_pom_ARTIFACT_URI", Value: *v1beta1.NewArrayOrString("com.google.guava:guava:31.0-jre.pom")},
					{Name: "mvn1_pom_ARTIFACT_DIGEST", Value: *v1beta1.NewArrayOrString(digest2)},
					{Name: "mvn1_src_ARTIFACT_URI", Value: *v1beta1.NewArrayOrString("com.google.guava:guava:31.0-jre-sources.jar")},
					{Name: "mvn1_src_ARTIFACT_DIGEST", Value: *v1beta1.NewArrayOrString(digest3)},
					{Name: "mvn2_ARTIFACT_URI", Value: *v1beta1.NewArrayOrString("projects/test-project/locations/us-west4/repositories/test-repo/mavenArtifacts/a.b.c:d:1.0-jre")},
					{Name: "mvn2_ARTIFACT_DIGEST", Value: *v1beta1.NewArrayOrString(digest4)},
					{Name: "ARTIFACT_URI", Value: *v1beta1.NewArrayOrString("projects/test-project/locations/us-west4/repositories/test-repo/mavenArtifacts/empty_prefix")},
					{Name: "ARTIFACT_DIGEST", Value: *v1beta1.NewArrayOrString(digest1)},
					{Name: "miss_target_name_ARTIFACT_DIGEST", Value: *v1beta1.NewArrayOrString(digest1)},
					{Name: "wrong_digest_format_ARTIFACT_URI", Value: *v1beta1.NewArrayOrString("projects/test-project/locations/us-west4/repositories/test-repo/mavenArtifacts/wrong_digest_format")},
					{Name: "wrong_digest_format_ARTIFACT_DIGEST", Value: *v1beta1.NewArrayOrString("abc")},
				},
			},
		},
	}
	want := []*StructuredSignable{
		{URI: "projects/test-project/locations/us-west4/repositories/test-repo/mavenArtifacts/com.google.guava:guava:31.0-jre", Digest: digest1},
		{URI: "com.google.guava:guava:31.0-jre.pom", Digest: digest2},
		{URI: "com.google.guava:guava:31.0-jre-sources.jar", Digest: digest3},
		{URI: "projects/test-project/locations/us-west4/repositories/test-repo/mavenArtifacts/a.b.c:d:1.0-jre", Digest: digest4},
		{URI: "projects/test-project/locations/us-west4/repositories/test-repo/mavenArtifacts/empty_prefix", Digest: digest1},
	}
	got := ExtractSignableTargetFromResults(objects.NewTaskRunObject(tr), logtesting.TestLogger(t))
	sort.Slice(got, func(i, j int) bool {
		return got[i].URI < got[j].URI
	})
	sort.Slice(want, func(i, j int) bool {
		return want[i].URI < want[j].URI
	})
	if !cmp.Equal(got, want, ignore...) {
		t.Fatalf("not the same %s", cmp.Diff(want, got, ignore...))
	}
}

func TestExtractStructuredTargetFromResults(t *testing.T) {
	tr := &v1beta1.TaskRun{
		Status: v1beta1.TaskRunStatus{
			TaskRunStatusFields: v1beta1.TaskRunStatusFields{
				TaskRunResults: []v1beta1.TaskRunResult{
					{
						Name: "mvn1_pkg" + "_" + ArtifactsOutputsResultName,
						Value: *v1beta1.NewObject(map[string]string{
							"uri":           "projects/test-project/locations/us-west4/repositories/test-repo/mavenArtifacts/com.google.guava:guava:31.0-jre",
							"digest":        digest1,
							"signable_type": "",
						}),
					},
					{
						Name: "mvn1_pom_sha512" + "_" + ArtifactsOutputsResultName,
						Value: *v1beta1.NewObject(map[string]string{
							"uri":           "com.google.guava:guava:31.0-jre.pom",
							"digest":        digest2,
							"signable_type": "",
						}),
					},
					{
						Name: "img1_input" + "_" + ArtifactsInputsResultName,
						Value: *v1beta1.NewObject(map[string]string{
							"uri":    "gcr.io/foo/bar",
							"digest": digest3,
						}),
					},
					{
						Name: "img2_input_sha1" + "_" + ArtifactsInputsResultName,
						Value: *v1beta1.NewObject(map[string]string{
							"uri":    "gcr.io/foo/bar",
							"digest": digest_sha1,
						}),
					},
					{
						Name: "img2_input_incorrect_sha1" + "_" + ArtifactsInputsResultName,
						Value: *v1beta1.NewObject(map[string]string{
							"uri":    "gcr.io/foo/bar",
							"digest": digest_incorrect_sha1,
						}),
					},
					{
						Name: "img3_input_sha384" + "_" + ArtifactsInputsResultName,
						Value: *v1beta1.NewObject(map[string]string{
							"uri":    "gcr.io/foo/bar",
							"digest": digest_sha384,
						}),
					},
					{
						Name: "img3_input_incorrect_sha384" + "_" + ArtifactsInputsResultName,
						Value: *v1beta1.NewObject(map[string]string{
							"uri":    "gcr.io/foo/bar",
							"digest": digest_incorrect_sha384,
						}),
					},
					{
						Name: "img4_input_sha512" + "_" + ArtifactsInputsResultName,
						Value: *v1beta1.NewObject(map[string]string{
							"uri":    "gcr.io/foo/bar",
							"digest": digest_sha512,
						}),
					},
					{
						Name: "img4_input_incorrect_sha512" + "_" + ArtifactsInputsResultName,
						Value: *v1beta1.NewObject(map[string]string{
							"uri":    "gcr.io/foo/bar",
							"digest": digest_incorrect_sha512,
						}),
					},
					{
						Name: "img2_input_no_digest" + "_" + ArtifactsInputsResultName,
						Value: *v1beta1.NewObject(map[string]string{
							"uri":    "gcr.io/foo/foo",
							"digest": "",
						}),
					},
				},
			},
		},
	}

	wantInputs := []*StructuredSignable{
		{URI: "gcr.io/foo/bar", Digest: digest3},
		{URI: "gcr.io/foo/bar", Digest: digest_sha1},
		{URI: "gcr.io/foo/bar", Digest: digest_sha384},
		{URI: "gcr.io/foo/bar", Digest: digest_sha512},
	}
	gotInputs := ExtractStructuredTargetFromResults(objects.NewTaskRunObject(tr), ArtifactsInputsResultName, logtesting.TestLogger(t))
	if diff := cmp.Diff(gotInputs, wantInputs, cmpopts.SortSlices(func(x, y *StructuredSignable) bool { return x.Digest < y.Digest })); diff != "" {
		t.Errorf("Inputs are not as expected: %v", diff)
	}

	wantOutputs := []*StructuredSignable{
		{URI: "projects/test-project/locations/us-west4/repositories/test-repo/mavenArtifacts/com.google.guava:guava:31.0-jre", Digest: digest1},
		{URI: "com.google.guava:guava:31.0-jre.pom", Digest: digest2},
	}
	gotOutputs := ExtractStructuredTargetFromResults(objects.NewTaskRunObject(tr), ArtifactsOutputsResultName, logtesting.TestLogger(t))
	opts := append(ignore, cmpopts.SortSlices(func(x, y *StructuredSignable) bool { return x.Digest < y.Digest }))
	if diff := cmp.Diff(gotOutputs, wantOutputs, opts...); diff != "" {
		t.Error(diff)
	}
}

func TestRetrieveMaterialsFromStructuredResults(t *testing.T) {
	tr := &v1beta1.TaskRun{
		Status: v1beta1.TaskRunStatus{
			TaskRunStatusFields: v1beta1.TaskRunStatusFields{
				TaskRunResults: []v1beta1.TaskRunResult{
					{
						Name: "img1_input" + "_" + ArtifactsInputsResultName,
						Value: *v1beta1.NewObject(map[string]string{
							"uri":    "gcr.io/foo/bar",
							"digest": "sha256:05f95b26ed10668b7183c1e2da98610e91372fa9f510046d4ce5812addad86b7",
						}),
					},
					{
						Name: "img2_input_no_digest" + "_" + ArtifactsInputsResultName,
						Value: *v1beta1.NewObject(map[string]string{
							"uri":    "gcr.io/foo/foo",
							"digest": "",
						}),
					},
					{
						Name: "img2_input_invalid_digest" + "_" + ArtifactsInputsResultName,
						Value: *v1beta1.NewObject(map[string]string{
							"uri":    "gcr.io/foo/foo",
							"digest": "sha:123",
						}),
					},
				},
			},
		},
	}
	wantMaterials := []slsa.ProvenanceMaterial{
		{
			URI:    "gcr.io/foo/bar",
			Digest: map[string]string{"sha256": "05f95b26ed10668b7183c1e2da98610e91372fa9f510046d4ce5812addad86b7"},
		},
	}

	gotMaterials := RetrieveMaterialsFromStructuredResults(objects.NewTaskRunObject(tr), ArtifactsInputsResultName, logtesting.TestLogger(t))

	if diff := cmp.Diff(gotMaterials, wantMaterials, ignore...); diff != "" {
		t.Fatalf("Materials not the same %s", diff)
	}
}

func TestValidateResults(t *testing.T) {
	tests := []struct {
		name           string
		obj            objects.Result
		categoryMarker string
		wantResult     bool
		wantErr        error
	}{
		{
			name:           "valid result",
			categoryMarker: ArtifactsOutputsResultName,
			obj: objects.Result{
				Name: "valid_result-ARTIFACT_OUTPUTS",
				Value: v1beta1.ParamValue{
					ObjectVal: map[string]string{
						"uri":    "gcr.io/foo/bar",
						"digest": digest3,
					},
				},
			},
			wantResult: true,
			wantErr:    nil,
		},
		{
			name:           "invalid result without digest field",
			categoryMarker: ArtifactsOutputsResultName,
			obj: objects.Result{
				Name: "missing_digest-ARTIFACT_OUTPUTS",
				Value: v1beta1.ParamValue{
					ObjectVal: map[string]string{
						"uri": "gcr.io/foo/bar",
					},
				},
			},
			wantResult: false,
			wantErr:    fmt.Errorf("missing_digest-ARTIFACT_OUTPUTS should have digest field: map[uri:gcr.io/foo/bar]"),
		},
		{
			name:           "invalid result without uri field",
			categoryMarker: ArtifactsOutputsResultName,
			obj: objects.Result{
				Name: "missing_digest-ARTIFACT_OUTPUTS",
				Value: v1beta1.ParamValue{
					ObjectVal: map[string]string{
						"digest": digest3,
					},
				},
			},
			wantResult: false,
			wantErr:    fmt.Errorf("missing_digest-ARTIFACT_OUTPUTS should have uri field: map[digest:sha256:05f95b26ed10668b7183c1e2da98610e91372fa9f510046d4ce5812addad86b7]"),
		},
		{
			name:           "invalid result wrong digest format",
			categoryMarker: ArtifactsOutputsResultName,
			obj: objects.Result{
				Name: "missing_digest-ARTIFACT_OUTPUTS",
				Value: v1beta1.ParamValue{
					ObjectVal: map[string]string{
						"uri":    "gcr.io/foo/bar",
						"digest": "",
					},
				},
			},
			wantResult: false,
			wantErr:    fmt.Errorf("missing_digest-ARTIFACT_OUTPUTS should have digest field: map[digest: uri:gcr.io/foo/bar]"),
		},
		{
			name:           "invalid result wrong type hinting",
			categoryMarker: ArtifactsOutputsResultName,
			obj: objects.Result{
				Name: "missing_digest-ARTIFACTs_OUTPUTS",
				Value: v1beta1.ParamValue{
					ObjectVal: map[string]string{
						"uri":    "gcr.io/foo/bar",
						"digest": digest3,
					},
				},
			},
			wantResult: false,
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isStructuredResult(tt.obj, tt.categoryMarker)
			if got != tt.wantResult {
				t.Errorf("Validation result is not as the expected: got %v and wanted %v", got, tt.wantResult)
			}
			if !tt.wantResult && tt.wantErr != nil {
				if diff := cmp.Diff(err.Error(), tt.wantErr.Error()); diff != "" {
					t.Errorf("Validation error is not as the expected: %s", diff)
				}
			}

		})
	}
}

func createDigest(t *testing.T, dgst string) name.Digest {
	result, err := name.NewDigest(dgst)
	if err != nil {
		t.Fatal(err)
	}
	return result

}
