package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateTargetNamespace(t *testing.T) {

	vtn := []struct {
		targetNamespace   string
		excludeNamespaces string
		valid             bool
	}{
		{
			targetNamespace:   "default",
			excludeNamespaces: "kube-system",
			valid:             true,
		},
		{
			targetNamespace:   "",
			excludeNamespaces: "kube-system",
			valid:             false,
		},
		{
			targetNamespace:   "",
			excludeNamespaces: "",
			valid:             false,
		},
		{
			targetNamespace:   "default",
			excludeNamespaces: "",
			valid:             true,
		},
		{
			targetNamespace:   "kube-system",
			excludeNamespaces: "kube-system",
			valid:             false,
		},
		{
			targetNamespace:   "default",
			excludeNamespaces: "kube-system,default",
			valid:             false,
		},
		{
			targetNamespace:   "default",
			excludeNamespaces: "kube-system,workloads",
			valid:             true,
		},
	}

	for _, element := range vtn {
		valid := validateTargetNamespace(element.targetNamespace, element.excludeNamespaces)
		assert.Equal(t, valid, element.valid)
	}
}
