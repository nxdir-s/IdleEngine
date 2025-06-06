package postgresqlcnpgio

import (
	_init_ "example.com/charts/imports/postgresqlcnpgio/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

// Name or number of the port to access on the container.
//
// Number must be in the range 1 to 65535.
// Name must be an IANA_SVC_NAME.
type PoolerSpecTemplateSpecInitContainersLifecyclePostStartHttpGetPort interface {
	Value() interface{}
}

// The jsii proxy struct for PoolerSpecTemplateSpecInitContainersLifecyclePostStartHttpGetPort
type jsiiProxy_PoolerSpecTemplateSpecInitContainersLifecyclePostStartHttpGetPort struct {
	_ byte // padding
}

func (j *jsiiProxy_PoolerSpecTemplateSpecInitContainersLifecyclePostStartHttpGetPort) Value() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"value",
		&returns,
	)
	return returns
}


func PoolerSpecTemplateSpecInitContainersLifecyclePostStartHttpGetPort_FromNumber(value *float64) PoolerSpecTemplateSpecInitContainersLifecyclePostStartHttpGetPort {
	_init_.Initialize()

	if err := validatePoolerSpecTemplateSpecInitContainersLifecyclePostStartHttpGetPort_FromNumberParameters(value); err != nil {
		panic(err)
	}
	var returns PoolerSpecTemplateSpecInitContainersLifecyclePostStartHttpGetPort

	_jsii_.StaticInvoke(
		"postgresqlcnpgio.PoolerSpecTemplateSpecInitContainersLifecyclePostStartHttpGetPort",
		"fromNumber",
		[]interface{}{value},
		&returns,
	)

	return returns
}

func PoolerSpecTemplateSpecInitContainersLifecyclePostStartHttpGetPort_FromString(value *string) PoolerSpecTemplateSpecInitContainersLifecyclePostStartHttpGetPort {
	_init_.Initialize()

	if err := validatePoolerSpecTemplateSpecInitContainersLifecyclePostStartHttpGetPort_FromStringParameters(value); err != nil {
		panic(err)
	}
	var returns PoolerSpecTemplateSpecInitContainersLifecyclePostStartHttpGetPort

	_jsii_.StaticInvoke(
		"postgresqlcnpgio.PoolerSpecTemplateSpecInitContainersLifecyclePostStartHttpGetPort",
		"fromString",
		[]interface{}{value},
		&returns,
	)

	return returns
}

