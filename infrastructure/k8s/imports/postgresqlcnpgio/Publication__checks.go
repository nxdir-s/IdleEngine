//go:build !no_runtime_type_checking

package postgresqlcnpgio

import (
	"fmt"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/constructs-go/constructs/v10"
)

func validatePublication_IsApiObjectParameters(o interface{}) error {
	if o == nil {
		return fmt.Errorf("parameter o is required, but nil was provided")
	}

	return nil
}

func validatePublication_IsConstructParameters(x interface{}) error {
	if x == nil {
		return fmt.Errorf("parameter x is required, but nil was provided")
	}

	return nil
}

func validatePublication_ManifestParameters(props *PublicationProps) error {
	if props == nil {
		return fmt.Errorf("parameter props is required, but nil was provided")
	}
	if err := _jsii_.ValidateStruct(props, func() string { return "parameter props" }); err != nil {
		return err
	}

	return nil
}

func validatePublication_OfParameters(c constructs.IConstruct) error {
	if c == nil {
		return fmt.Errorf("parameter c is required, but nil was provided")
	}

	return nil
}

func validateNewPublicationParameters(scope constructs.Construct, id *string, props *PublicationProps) error {
	if scope == nil {
		return fmt.Errorf("parameter scope is required, but nil was provided")
	}

	if id == nil {
		return fmt.Errorf("parameter id is required, but nil was provided")
	}

	if props == nil {
		return fmt.Errorf("parameter props is required, but nil was provided")
	}
	if err := _jsii_.ValidateStruct(props, func() string { return "parameter props" }); err != nil {
		return err
	}

	return nil
}

