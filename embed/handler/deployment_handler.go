package handler

import (
	"bytes"
	"embed"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

//go:embed manifests
var manifestFS embed.FS

type DeployParams struct {
	AppName   string `json:"app_name"`
	Namespace string `json:"namespace"`

	// Deployment specific
	Replicas      int    `json:"replicas"`
	Image         string `json:"image"`
	ContainerPort int    `json:"container_port"`
}

type manifestHandler struct {
	deploymentTmpl *template.Template
	serviceTmpl    *template.Template
}

func NewManifestHandler() (*manifestHandler, error) {
	deploymentYaml, err := manifestFS.ReadFile("manifests/deployment.yaml")
	if err != nil {
		return nil, err
	}

	deploymentTmpl, err := template.New("deployment").Parse(string(deploymentYaml))
	if err != nil {
		return nil, err
	}

	serviceYaml, err := manifestFS.ReadFile("manifests/service.yaml")
	if err != nil {
		return nil, err
	}

	serviceTempl, err := template.New("service").Parse(string(serviceYaml))
	if err != nil {
		return nil, err
	}

	return &manifestHandler{
		deploymentTmpl: deploymentTmpl,
		serviceTmpl:    serviceTempl,
	}, nil
}

func (hdl *manifestHandler) GetDeploymentManifest(c *gin.Context) {
	req := DeployParams{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := hdl.bindDeploymentValue(&req)
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

func (hdl *manifestHandler) bindDeploymentValue(req *DeployParams) (*unstructured.Unstructured, error) {
	var buf bytes.Buffer

	if err := hdl.deploymentTmpl.Execute(&buf, req); err != nil {
		return nil, err
	}

	obj := &unstructured.Unstructured{}
	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

	_, _, err := dec.Decode(buf.Bytes(), nil, obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}
