package postgresqlcnpgio


// Projection that may be projected along with other supported volume types.
//
// Exactly one of these fields must be set.
type ClusterSpecProjectedVolumeTemplateSources struct {
	// ClusterTrustBundle allows a pod to access the `.spec.trustBundle` field of ClusterTrustBundle objects in an auto-updating file.
	//
	// Alpha, gated by the ClusterTrustBundleProjection feature gate.
	//
	// ClusterTrustBundle objects can either be selected by name, or by the
	// combination of signer name and a label selector.
	//
	// Kubelet performs aggressive normalization of the PEM contents written
	// into the pod filesystem.  Esoteric PEM features such as inter-block
	// comments and block headers are stripped.  Certificates are deduplicated.
	// The ordering of certificates within the file is arbitrary, and Kubelet
	// may change the order over time.
	ClusterTrustBundle *ClusterSpecProjectedVolumeTemplateSourcesClusterTrustBundle `field:"optional" json:"clusterTrustBundle" yaml:"clusterTrustBundle"`
	// configMap information about the configMap data to project.
	ConfigMap *ClusterSpecProjectedVolumeTemplateSourcesConfigMap `field:"optional" json:"configMap" yaml:"configMap"`
	// downwardAPI information about the downwardAPI data to project.
	DownwardApi *ClusterSpecProjectedVolumeTemplateSourcesDownwardApi `field:"optional" json:"downwardApi" yaml:"downwardApi"`
	// secret information about the secret data to project.
	Secret *ClusterSpecProjectedVolumeTemplateSourcesSecret `field:"optional" json:"secret" yaml:"secret"`
	// serviceAccountToken is information about the serviceAccountToken data to project.
	ServiceAccountToken *ClusterSpecProjectedVolumeTemplateSourcesServiceAccountToken `field:"optional" json:"serviceAccountToken" yaml:"serviceAccountToken"`
}

