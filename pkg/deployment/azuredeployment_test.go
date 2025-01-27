package deployment_test

import (
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/stretchr/testify/assert"
)

func TestStartDeployment(t *testing.T) {
	fullPath := filepath.Join("../../test/testdata/nameviolation/nestedfailure", "mainTemplate.json")
	template, err := utils.ReadJson(fullPath)
	assert.NoError(t, err)
	assert.NotNil(t, template)

	resources := deployment.FindResourcesByType(template, "Microsoft.Resources/deployments")
	assert.Greater(t, len(resources), 0)
}

func TestGetDeploymentParamsNested(t *testing.T) {
	templatePath := filepath.Join("../../test/testdata/missingparam", "mainTemplate.json")
	template, err := utils.ReadJson(templatePath)
	assert.NoError(t, err)
	assert.NotNil(t, template)

	paramsPath := filepath.Join("../../test/testdata/missingparam", "parameters.json")
	params, err := utils.ReadJson(paramsPath)
	assert.NoError(t, err)
	assert.NotNil(t, params)

	azureDeployment := &deployment.AzureDeployment{
		SubscriptionId: uuid.NewString(),
		Location:       "eastus",
		ResourceGroupName: "TestResourceGroup",
		DeploymentName: "TestDeployment",
		DeploymentType: deployment.AzureResourceManager,
		Template:       template,
		Params:         params,
	}

	result := azureDeployment.GetParams()
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result))
	assert.NotNil(t, result["aapName"])
	assert.NotNil(t, result["testName"])
}

func TestGetDeploymentParamsUnNested(t *testing.T) {
	templatePath := filepath.Join("../../test/testdata/missingparam", "mainTemplate.json")
	template, err := utils.ReadJson(templatePath)
	assert.NoError(t, err)
	assert.NotNil(t, template)

	paramsMap := make(map[string]interface{})
	
	aapNameMap := make(map[string]interface{})
	aapNameMap["value"] = "test"

	testNameMap := make(map[string]interface{})
	testNameMap["value"] = "test2"
	
	paramsMap["aapName"] = aapNameMap
	paramsMap["testName"] = testNameMap
	
	azureDeployment := &deployment.AzureDeployment{
		SubscriptionId: uuid.NewString(),
		Location:       "eastus",
		ResourceGroupName: "TestResourceGroup",
		DeploymentName: "TestDeployment",
		DeploymentType: deployment.AzureResourceManager,
		Template:       template,
		Params:         paramsMap,
	}

	result := azureDeployment.GetParams()
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result))
	assert.NotNil(t, result["aapName"])
	assert.NotNil(t, result["testName"])

	aapNameValue := result["aapName"].(map[string]interface{})["value"]
	assert.Equal(t, "test", aapNameValue)

	testNameValue := result["testName"].(map[string]interface{})["value"]
	assert.Equal(t, "test2", testNameValue)
}

func TestGetDeploymentTemplateParamsNested(t *testing.T) {
	templatePath := filepath.Join("../../test/testdata/missingparam", "mainTemplate.json")
	template, err := utils.ReadJson(templatePath)
	assert.NoError(t, err)
	assert.NotNil(t, template)

	paramsPath := filepath.Join("../../test/testdata/missingparam", "parameters.json")
	params, err := utils.ReadJson(paramsPath)
	assert.NoError(t, err)
	assert.NotNil(t, params)

	azureDeployment := &deployment.AzureDeployment{
		SubscriptionId: uuid.NewString(),
		Location:       "eastus",
		ResourceGroupName: "TestResourceGroup",
		DeploymentName: "TestDeployment",
		DeploymentType: deployment.AzureResourceManager,
		Template:       template,
		Params:         params,
	}

	result := azureDeployment.GetTemplateParams()
	assert.NotNil(t, result)
	assert.Equal(t, 3, len(result))
	assert.NotNil(t, result["aapName"])
	assert.NotNil(t, result["testName"])
	assert.NotNil(t, result["testName2"])
}