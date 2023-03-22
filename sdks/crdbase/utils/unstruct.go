package utils

import (
	"errors"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
)

var ErrInvalidUnstructuredContent = errors.New("invalid unstructured content or path")

// GetValueFormUnstructuredContent get value from unstructured content
func GetValueFormUnstructuredContent(in map[string]interface{}, path string) (interface{}, error) {
	// deep copy the content to avoid changing the original content
	uc := runtime.DeepCopyJSON(in)

	parts := strings.Split(path, ".")
	for _, part := range parts {
		// if unstructured content is nil, return nil and error
		if uc == nil {
			return nil, ErrInvalidUnstructuredContent
		}
		if v, ok := uc[part]; ok {
			// if value is map, set unstructured content to the map and do next loop
			if m, ok := v.(map[string]interface{}); ok {
				uc = m
				continue
			}
			return v, nil
		}
		// if value not found, return nil and error
		return nil, ErrInvalidUnstructuredContent
	}
	// if path empty, return the whole content
	return uc, nil
}
