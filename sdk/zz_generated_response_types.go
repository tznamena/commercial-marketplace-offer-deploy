//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package sdk

// DeploymentManagementClientCreateDeploymentResponse contains the response from method DeploymentManagementClient.CreateDeployment.
type DeploymentManagementClientCreateDeploymentResponse struct {
	Deployment
}

// DeploymentManagementClientCreateEvenHookResponse contains the response from method DeploymentManagementClient.CreateEvenHook.
type DeploymentManagementClientCreateEvenHookResponse struct {
	CreateEventHookResponse
}

// DeploymentManagementClientDeleteEventHookResponse contains the response from method DeploymentManagementClient.DeleteEventHook.
type DeploymentManagementClientDeleteEventHookResponse struct {
	// placeholder for future response values
}

// DeploymentManagementClientGetDeploymentResponse contains the response from method DeploymentManagementClient.GetDeployment.
type DeploymentManagementClientGetDeploymentResponse struct {
	Deployment
}

// DeploymentManagementClientGetEventHookResponse contains the response from method DeploymentManagementClient.GetEventHook.
type DeploymentManagementClientGetEventHookResponse struct {
	EventHookResponse
}

// DeploymentManagementClientGetEventTypesResponse contains the response from method DeploymentManagementClient.GetEventTypes.
type DeploymentManagementClientGetEventTypesResponse struct {
	// Array of EventType
	EventTypeArray []*EventType
}

// DeploymentManagementClientGetInvokedDeploymentOperationResponse contains the response from method DeploymentManagementClient.GetInvokedDeploymentOperation.
type DeploymentManagementClientGetInvokedDeploymentOperationResponse struct {
	GetInvokedOperationResponse
}

// DeploymentManagementClientInvokeDeploymentOperationResponse contains the response from method DeploymentManagementClient.InvokeDeploymentOperation.
type DeploymentManagementClientInvokeDeploymentOperationResponse struct {
	InvokedDeploymentOperationResponse
}

// DeploymentManagementClientListDeploymentsResponse contains the response from method DeploymentManagementClient.ListDeployments.
type DeploymentManagementClientListDeploymentsResponse struct {
	// Array of Deployment
	DeploymentArray []*Deployment
}

// DeploymentManagementClientListEventHooksResponse contains the response from method DeploymentManagementClient.ListEventHooks.
type DeploymentManagementClientListEventHooksResponse struct {
	// Array of EventHookResponse
	EventHookResponseArray []*EventHookResponse
}

// DeploymentManagementClientListInvokedOperationsResponse contains the response from method DeploymentManagementClient.ListInvokedOperations.
type DeploymentManagementClientListInvokedOperationsResponse struct {
	ListInvokedOperationResponse
}

// DeploymentManagementClientUpdateDeploymentResponse contains the response from method DeploymentManagementClient.UpdateDeployment.
type DeploymentManagementClientUpdateDeploymentResponse struct {
	// placeholder for future response values
}

