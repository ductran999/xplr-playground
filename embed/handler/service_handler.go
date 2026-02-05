package handler

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

type ServiceParams struct {
	AppName       string `json:"app_name"`
	Namespace     string `json:"namespace"`
	ServiceType   string `json:"service_type"`
	ServicePort   int    `json:"service_port"`
	ContainerPort int    `json:"container_port"`
}

func (hdl *manifestHandler) GetServiceManifest(c *gin.Context) {
	req := ServiceParams{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := hdl.bindServiceValue(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"data":    result,
	})
}

func (hdl *manifestHandler) bindServiceValue(req *ServiceParams) (*unstructured.Unstructured, error) {
	var outputBuffer bytes.Buffer

	if err := hdl.serviceTmpl.Execute(&outputBuffer, req); err != nil {
		return nil, fmt.Errorf("error executing template: %w", err)
	}

	obj := &unstructured.Unstructured{}
	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	_, _, err := dec.Decode(outputBuffer.Bytes(), nil, obj)
	if err != nil {
		return nil, fmt.Errorf("failed to decode YAML: %w", err)
	}

	return obj, nil
}
