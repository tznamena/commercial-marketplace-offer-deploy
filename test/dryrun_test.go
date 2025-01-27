//go:build integration
// +build integration
package test

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armpolicy"

	"path/filepath"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type dryRunSuite struct {
	suite.Suite
	subscriptionId    string
	resourceGroupName string
	location          string
	endpoint          string
	deploymentId      int
}

func TestDryRunSuite(t *testing.T) {
	suite.Run(t, &dryRunSuite{})
}

func (s *dryRunSuite) SetupSuite() {
	s.subscriptionId = "31e9f9a0-9fd2-4294-a0a3-0101246d9700"
	s.resourceGroupName = "aMODMTestb"
	s.location = "eastus"
	s.endpoint = "http://localhost:8080"
}

func (s *dryRunSuite) SetupTest() {
	s.SetupResourceGroup()
	s.DeployPolicyDefintion()
}

func (s *dryRunSuite) runDeploymentTest(path string, errorExpected bool, template map[string]interface{}, params map[string]interface{}) *sdk.DryRunResult {
	azureDeployment := &deployment.AzureDeployment{
		SubscriptionId:    s.subscriptionId,
		Location:          s.location,
		ResourceGroupName: s.resourceGroupName,
		DeploymentName:    "DryRunDeploymentTest",
		DeploymentType:    deployment.AzureResourceManager,
		Template:          template,
		Params:            params,
	}

	s.Assert().NotNil(azureDeployment)

	ctx := context.TODO()
	resp, err := deployment.DryRun(ctx, azureDeployment)
	if errorExpected {
		s.Assert().NoError(err)
	}

	s.Assert().NotNil(resp)
	return resp
}

func (s *dryRunSuite) getJsonAsMap(path string) map[string]interface{} {
	jsonMap, err := utils.ReadJson(path)
	if err != nil {
		log.Println(err)
	}
	return jsonMap
}

func (s *dryRunSuite) TestNamePolicyFailure() {
	nameViolationPath := "./testdata/nameviolation/failure"
	result := s.runDeploymentTest(nameViolationPath, true, s.getTemplate(nameViolationPath), s.getParameters(nameViolationPath))
	log.Print("TestNamePolicyFailure Results:\n %s" + *s.prettify(result))
}

func (s *dryRunSuite) TestExistingStorageFailure() {
	nameViolationPath := "./testdata/existingstorage"
	result := s.runDeploymentTest(nameViolationPath, true, s.getTemplate(nameViolationPath), s.getParameters(nameViolationPath))
	log.Print("TestNamePolicyFailure Results:\n %s" + *s.prettify(result))
}

func (s *dryRunSuite) TestTaggedDeployment() {
	taggedDeploymentPath := "./testdata/taggeddeployment"
	result := s.runDeploymentTest(taggedDeploymentPath, true, s.getTemplate(taggedDeploymentPath), s.getParameters(taggedDeploymentPath))
	log.Print("TestNamePolicyFailure Results:\n %s" + *s.prettify(result))
}

func (s *dryRunSuite) TestMissingParameter() {
	taggedDeploymentPath := "./testdata/missingparam"
	result := s.runDeploymentTest(taggedDeploymentPath, true, s.getTemplate(taggedDeploymentPath), s.getParameters(taggedDeploymentPath))
	log.Print("TestMissingParameter Results:\n %s" + *s.prettify(result))
}

func (s *dryRunSuite) TestNestedPolicyFailure() {
	nameViolationPath := "./testdata/nameviolation/nestedfailure"
	result := s.runDeploymentTest(nameViolationPath, true, s.getTemplate(nameViolationPath), s.getParameters(nameViolationPath))
	log.Print("TestNamePolicyFailure Results:\n %s" + *s.prettify(result))
}

func (s *dryRunSuite) TestDirectParamsMap() {
	paramsMapPath := "./testdata/directparamsmap"
	valueMap := map[string]interface{}{
		"value"	: "bobjacbicep2",
	}
	paramsMap := map[string]interface{}{
		"testName": valueMap,
	}
	result := s.runDeploymentTest(paramsMapPath, true, s.getTemplate(paramsMapPath), paramsMap)
	s.Assert().NotNil(result)
}

func (s *dryRunSuite) TestDirectParamsMapStringValue() {
	paramsMapPath := "./testdata/directparamsmap"
	paramsMap := map[string]interface{}{
		"testName": "bobjacbicep2",
	}
	result := s.runDeploymentTest(paramsMapPath, true, s.getTemplate(paramsMapPath), paramsMap)
	s.Assert().NotNil(result)
}

func (s *dryRunSuite) TestQuotaViolation() {
	quotaViolationPath := "./testdata/quotaviolation"
	result := s.runDeploymentTest(quotaViolationPath, true, s.getTemplate(quotaViolationPath), s.getParameters(quotaViolationPath))
	require.NotNil(s.T(), result)
	log.Print("TestQuotaViolation Results:\n %s" + *s.prettify(result))
}

func (s *dryRunSuite) prettify(obj any) *string {
	bytes, _ := json.MarshalIndent(obj, "", "  ")
	result := string(bytes)
	return &result
}

func (s *dryRunSuite) getParameters(path string) map[string]interface{} {
	paramsPath := filepath.Join(path, "parameters.json")
	parameters, err := utils.ReadJson(paramsPath)
	require.NoError(s.T(), err)
	return parameters
}

func (s *dryRunSuite) createDeployment(ctx context.Context, client *sdk.Client, templatePath string) *sdk.Deployment {
	name := "DryRunDeploymentTest"
	template := s.getTemplate(templatePath)

	deployment, err := client.Create(ctx, sdk.CreateDeployment{
		Name:           &name,
		SubscriptionID: &s.subscriptionId,
		ResourceGroup:  &s.resourceGroupName,
		Location:       &s.location,
		Template:       template,
	})
	require.NoError(s.T(), err)

	return deployment
}

func (s *dryRunSuite) getTemplate(path string) map[string]interface{} {
	fullPath := filepath.Join(path, "mainTemplate.json")
	template, err := utils.ReadJson(fullPath)
	require.NoError(s.T(), err)
	return template
}

func (s *dryRunSuite) ResourceGroupExists() bool {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(s.T(), err)

	ctx := context.Background()

	resourceGroupClient, err := armresources.NewResourceGroupsClient(s.subscriptionId, cred, nil)
	require.NoError(s.T(), err)

	resp, err := resourceGroupClient.CheckExistence(ctx, s.resourceGroupName, nil)
	require.NoError(s.T(), err)

	return resp.Success
}

func (s *dryRunSuite) SetupResourceGroup() {
	if exists := s.ResourceGroupExists(); exists {
		return
	}

	_, err := s.CreateResourceGroup()
	require.NoError(s.T(), err)
}

func (s *dryRunSuite) CreateResourceGroup() (*armresources.ResourceGroup, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Print(err)
	}
	ctx := context.Background()

	resourceGroupClient, err := armresources.NewResourceGroupsClient(s.subscriptionId, cred, nil)
	if err != nil {
		return nil, err
	}

	resourceGroupResp, err := resourceGroupClient.CreateOrUpdate(
		ctx,
		s.resourceGroupName,
		armresources.ResourceGroup{
			Location: to.Ptr(s.location),
		},
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &resourceGroupResp.ResourceGroup, nil
}

func (s *dryRunSuite) DeployPolicyDefintion() {
	log.Printf("Inside deployPolicyDefinition()")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(s.T(), err)

	ctx := context.Background()
	client, err := armpolicy.NewDefinitionsClient(s.subscriptionId, cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.CreateOrUpdate(ctx,
		"ResourceNaming",
		armpolicy.Definition{
			Properties: &armpolicy.DefinitionProperties{
				Description: to.Ptr("Force resource names to begin with given 'prefix' and/or end with given 'suffix'"),
				DisplayName: to.Ptr("Enforce resource naming convention"),
				Metadata: map[string]interface{}{
					"category": "Naming",
				},
				Mode: to.Ptr("All"),
				PolicyRule: map[string]interface{}{
					"if": map[string]interface{}{
						"not": map[string]interface{}{
							"field": "name",
							"like":  "a*b",
						},
					},
					"then": map[string]interface{}{
						"effect": "deny",
					},
				},
			},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

func (s *dryRunSuite) DeployPolicy() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(s.T(), err)

	ctx := context.Background()
	client, err := armpolicy.NewAssignmentsClient(s.subscriptionId, cred, nil)
	require.NoError(s.T(), err)

	scope := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", s.subscriptionId, s.resourceGroupName)
	policyDefinitionId := fmt.Sprintf("/subscriptions/%s/providers/Microsoft.Authorization/policyDefinitions/ResourceNaming", s.subscriptionId)

	_, err = client.Create(ctx,
		scope,
		"aResourceNameb",
		armpolicy.Assignment{
			Properties: &armpolicy.AssignmentProperties{
				Description: to.Ptr("Enforce resource naming conventions"),
				DisplayName: to.Ptr("Enforce Resource Names"),
				Scope:       &scope,
				Metadata: map[string]interface{}{
					"assignedBy": "John Doe",
				},
				NonComplianceMessages: []*armpolicy.NonComplianceMessage{
					{
						Message: to.Ptr("A resource name was non-complaint.  It must be in the format 'a*b'."),
					}},
				PolicyDefinitionID: to.Ptr(policyDefinitionId),
			},
		},
		nil)
	require.NoError(s.T(), err)
}
