package deployment

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

type ExecutionStatus string

const (
	Started           ExecutionStatus = "Started"
	Failed            ExecutionStatus = "Failed"
	PermanentlyFailed ExecutionStatus = "PermanentlyFailed"
	Succeeded         ExecutionStatus = "Succeeded"
	Restart           ExecutionStatus = "Restart"
	Restarted         ExecutionStatus = "Restarted"
	RestartTimedOut   ExecutionStatus = "RestartTimedOut"
	Canceled          ExecutionStatus = "Canceled"
)

type AzureDeployment struct {
	subscriptionId string
	location string
	resourceGroupName string
	deploymentName string
	deploymentType DeploymentType
	template Template
	params TemplateParams
	resumeToken string
}

type AzureDeploymentResult struct {
	ID                string                 `json:"id"`
	CorrelationID     string                 `json:"correlationId"`
	Duration          string                 `json:"duration"`
	Timestamp         time.Time              `json:"timestamp"`
	ProvisioningState string                 `json:"provisioningState"`
	Outputs           map[string]interface{} `json:"outputs"`
	Status            ExecutionStatus
}

func (ad *AzureDeployment) GetDeploymentType() DeploymentType {
	return ad.deploymentType
}

func (ad *AzureDeployment) GetTemplate() map[string]interface{} {
	return ad.template
}

func (ad *AzureDeployment) GetTemplateParams() map[string]interface{} {
	return ad.params
}

type Deployer interface {
	Deploy(d *AzureDeployment) (*AzureDeploymentResult, error)
}

type ArmTemplateDeployer struct {
	deployerType DeploymentType
}

func (armDeployer *ArmTemplateDeployer) Deploy(ad *AzureDeployment) (*AzureDeploymentResult, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Print(err)
	}
	ctx := context.Background()
	deploymentsClient, err := armresources.NewDeploymentsClient(ad.subscriptionId, cred, nil)
	if err != nil {
		return nil, err
	}

	log.Printf("About to Create a deployment")

	deploymentPollerResp, err := deploymentsClient.BeginCreateOrUpdate(
		ctx,
		ad.resourceGroupName,
		ad.deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Template:   ad.template,
				Parameters: ad.params,
				Mode:       to.Ptr(armresources.DeploymentModeIncremental),
			},
		},
		nil)

	if err != nil {
		return nil, fmt.Errorf("cannot create deployment: %v", err)
	}

	//todo: capture the state of the started deployment
	resp, err := deploymentPollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot get the create deployment future respone: %v", err)
	}

	mappedResult, err := armDeployer.mapDeploymentResult(resp)
	if err != nil {
		return nil, err
	}

	return mappedResult, nil
}

func (armDeployer *ArmTemplateDeployer) mapDeploymentResult(resp armresources.DeploymentsClientCreateOrUpdateResponse)	(*AzureDeploymentResult, error) {
	var status ExecutionStatus
	deploymentExtended := resp.DeploymentExtended
	provisioningState := *deploymentExtended.Properties.ProvisioningState
	switch provisioningState {
	case armresources.ProvisioningStateSucceeded:
		status = Succeeded
	case armresources.ProvisioningStateCanceled:
		status = Canceled
	default:
		status = Failed
	}
	// make sure response outputs are always there, even if empty
	var responseOutputs map[string]interface{}
	if deploymentExtended.Properties.Outputs != nil {
		responseOutputs = deploymentExtended.Properties.Outputs.(map[string]interface{})
	} else {
		responseOutputs = make(map[string]interface{})
	}
	res := AzureDeploymentResult{}
	if &provisioningState != nil {
		res.ProvisioningState = string(*deploymentExtended.Properties.ProvisioningState)
	}
	if deploymentExtended.ID != nil {
		res.ID = *deploymentExtended.ID
	}
	if deploymentExtended.Properties.CorrelationID != nil {
		res.CorrelationID = *deploymentExtended.Properties.CorrelationID
	}
	if deploymentExtended.Properties.Duration != nil {
		res.Duration = *deploymentExtended.Properties.Duration
	}
	if deploymentExtended.Properties.Timestamp != nil {
		res.Timestamp = *deploymentExtended.Properties.Timestamp
	}
	res.Status = status
	res.Outputs = responseOutputs
	return &res, nil
}

func CreateNewDeployer(deployment AzureDeployment) Deployer {
	return &ArmTemplateDeployer{
		deployerType: deployment.deploymentType,
	}
}

func (azureDeployment *AzureDeployment) validate() (error) {
	if len(azureDeployment.subscriptionId) == 0 {
		return errors.New("subscriptionId is not set on azureDeployment input struct")
	}
	if len(azureDeployment.location) == 0 {
		return errors.New("location is not set on azureDeployment input struct")
	}
	if len(azureDeployment.resourceGroupName) == 0 {
		return errors.New("resourceGroupName is not set on azureDeployment input struct")
	}
	if len(azureDeployment.resourceGroupName) == 0 {
		return errors.New("resourceGroupName is not set on azureDeployment input struct")
	}
	if azureDeployment.template == nil {
		return errors.New("template is not set on deployment azureDeployment struct")
	}
	// allow params to be empty to support all default params
	return nil
}

// ErrorAdditionalInfo - The resource management error additional info.
type ErrorAdditionalInfo struct {
	// READ-ONLY; The additional info.
	Info interface{} `json:"info,omitempty" azure:"ro"`

	// READ-ONLY; The additional info type.
	Type *string `json:"type,omitempty" azure:"ro"`
}