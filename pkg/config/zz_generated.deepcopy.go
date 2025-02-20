//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package config

import (
	sets "k8s.io/apimachinery/pkg/util/sets"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Artifact) DeepCopyInto(out *Artifact) {
	*out = *in
	if in.StorageBackend != nil {
		in, out := &in.StorageBackend, &out.StorageBackend
		*out = make(sets.Set[string], len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Artifact.
func (in *Artifact) DeepCopy() *Artifact {
	if in == nil {
		return nil
	}
	out := new(Artifact)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArtifactConfigs) DeepCopyInto(out *ArtifactConfigs) {
	*out = *in
	in.TaskRuns.DeepCopyInto(&out.TaskRuns)
	in.OCI.DeepCopyInto(&out.OCI)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArtifactConfigs.
func (in *ArtifactConfigs) DeepCopy() *ArtifactConfigs {
	if in == nil {
		return nil
	}
	out := new(ArtifactConfigs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BuilderConfig) DeepCopyInto(out *BuilderConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BuilderConfig.
func (in *BuilderConfig) DeepCopy() *BuilderConfig {
	if in == nil {
		return nil
	}
	out := new(BuilderConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Config) DeepCopyInto(out *Config) {
	*out = *in
	in.Artifacts.DeepCopyInto(&out.Artifacts)
	out.Storage = in.Storage
	out.Signers = in.Signers
	out.Builder = in.Builder
	out.Transparency = in.Transparency
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Config.
func (in *Config) DeepCopy() *Config {
	if in == nil {
		return nil
	}
	out := new(Config)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DocDBStorageConfig) DeepCopyInto(out *DocDBStorageConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DocDBStorageConfig.
func (in *DocDBStorageConfig) DeepCopy() *DocDBStorageConfig {
	if in == nil {
		return nil
	}
	out := new(DocDBStorageConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GCSStorageConfig) DeepCopyInto(out *GCSStorageConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GCSStorageConfig.
func (in *GCSStorageConfig) DeepCopy() *GCSStorageConfig {
	if in == nil {
		return nil
	}
	out := new(GCSStorageConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KMSSigner) DeepCopyInto(out *KMSSigner) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KMSSigner.
func (in *KMSSigner) DeepCopy() *KMSSigner {
	if in == nil {
		return nil
	}
	out := new(KMSSigner)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OCIStorageConfig) DeepCopyInto(out *OCIStorageConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OCIStorageConfig.
func (in *OCIStorageConfig) DeepCopy() *OCIStorageConfig {
	if in == nil {
		return nil
	}
	out := new(OCIStorageConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SignerConfigs) DeepCopyInto(out *SignerConfigs) {
	*out = *in
	out.X509 = in.X509
	out.KMS = in.KMS
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SignerConfigs.
func (in *SignerConfigs) DeepCopy() *SignerConfigs {
	if in == nil {
		return nil
	}
	out := new(SignerConfigs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StorageConfigs) DeepCopyInto(out *StorageConfigs) {
	*out = *in
	out.GCS = in.GCS
	out.OCI = in.OCI
	out.Tekton = in.Tekton
	out.DocDB = in.DocDB
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StorageConfigs.
func (in *StorageConfigs) DeepCopy() *StorageConfigs {
	if in == nil {
		return nil
	}
	out := new(StorageConfigs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StorageOpts) DeepCopyInto(out *StorageOpts) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StorageOpts.
func (in *StorageOpts) DeepCopy() *StorageOpts {
	if in == nil {
		return nil
	}
	out := new(StorageOpts)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TektonStorageConfig) DeepCopyInto(out *TektonStorageConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TektonStorageConfig.
func (in *TektonStorageConfig) DeepCopy() *TektonStorageConfig {
	if in == nil {
		return nil
	}
	out := new(TektonStorageConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TransparencyConfig) DeepCopyInto(out *TransparencyConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TransparencyConfig.
func (in *TransparencyConfig) DeepCopy() *TransparencyConfig {
	if in == nil {
		return nil
	}
	out := new(TransparencyConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *X509Signer) DeepCopyInto(out *X509Signer) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new X509Signer.
func (in *X509Signer) DeepCopy() *X509Signer {
	if in == nil {
		return nil
	}
	out := new(X509Signer)
	in.DeepCopyInto(out)
	return out
}
