package postgresqlcnpgio


// dataSourceRef specifies the object from which to populate the volume with data, if a non-empty volume is desired.
//
// This may be any object from a non-empty API group (non
// core object) or a PersistentVolumeClaim object.
// When this field is specified, volume binding will only succeed if the type of
// the specified object matches some installed volume populator or dynamic
// provisioner.
// This field will replace the functionality of the dataSource field and as such
// if both fields are non-empty, they must have the same value. For backwards
// compatibility, when namespace isn't specified in dataSourceRef,
// both fields (dataSource and dataSourceRef) will be set to the same
// value automatically if one of them is empty and the other is non-empty.
// When namespace is specified in dataSourceRef,
// dataSource isn't set to the same value and must be empty.
// There are three important differences between dataSource and dataSourceRef:
// * While dataSource only allows two specific types of objects, dataSourceRef
// allows any non-core object, as well as PersistentVolumeClaim objects.
// * While dataSource ignores disallowed values (dropping them), dataSourceRef
// preserves all values, and generates an error if a disallowed value is
// specified.
// * While dataSource only allows local objects, dataSourceRef allows objects
// in any namespaces.
// (Beta) Using this field requires the AnyVolumeDataSource feature gate to be enabled.
// (Alpha) Using the namespace field of dataSourceRef requires the CrossNamespaceVolumeDataSource feature gate to be enabled.
type ClusterSpecTablespacesStoragePvcTemplateDataSourceRef struct {
	// Kind is the type of resource being referenced.
	Kind *string `field:"required" json:"kind" yaml:"kind"`
	// Name is the name of resource being referenced.
	Name *string `field:"required" json:"name" yaml:"name"`
	// APIGroup is the group for the resource being referenced.
	//
	// If APIGroup is not specified, the specified Kind must be in the core API group.
	// For any other third-party types, APIGroup is required.
	ApiGroup *string `field:"optional" json:"apiGroup" yaml:"apiGroup"`
	// Namespace is the namespace of resource being referenced Note that when a namespace is specified, a gateway.networking.k8s.io/ReferenceGrant object is required in the referent namespace to allow that namespace's owner to accept the reference. See the ReferenceGrant documentation for details. (Alpha) This field requires the CrossNamespaceVolumeDataSource feature gate to be enabled.
	Namespace *string `field:"optional" json:"namespace" yaml:"namespace"`
}

