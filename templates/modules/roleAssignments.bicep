param storageAccountName string
param containerGroupName string
param serviceBusNamespace string

resource containerGroup 'Microsoft.ContainerInstance/containerGroups@2021-09-01' existing = {
  name: containerGroupName
}

resource storageAccount 'Microsoft.Storage/storageAccounts@2021-09-01' existing = {
  name: storageAccountName
}

resource serviceBus 'Microsoft.ServiceBus/namespaces@2022-01-01-preview' existing = {
  name: serviceBusNamespace
}

var roles = {
  resourceGroupReader: 'acdd72a7-3385-48ef-bd42-f606fba81ae7'
  storageAccountContributor: '0c867c2a-1d8c-454a-a3db-ab2ea1bdc8bb'

  serviceBusDataOwner: '090c5cfd-751d-490a-894a-3ce6f1109419'
  
  eventGridContributor: '1e241071-0855-49ea-94dc-649edcd759de'
  eventGridDataSender: 'd5a91429-5739-47e2-a06b-3470a27159e7'
  eventGridEventSubscriptionContributor: '428e0ff0-5e57-4d9c-a221-2c70d0e0a443'
  eventGridEventSubscriptionReader: '2414bbcf-6497-4faf-8c65-045460748405'
}

resource resourceReaderAssignment 'Microsoft.Authorization/roleAssignments@2020-04-01-preview' = {
  scope: resourceGroup()
  name: guid(resourceGroup().id, containerGroup.name, roles.resourceGroupReader)
  properties: {
    roleDefinitionId: resourceId('Microsoft.Authorization/roleDefinitions', roles.resourceGroupReader)
    principalId: containerGroup.identity.principalId
    principalType: 'ServicePrincipal'
  }
}

resource storageAccountAssignment 'Microsoft.Authorization/roleAssignments@2020-04-01-preview' = {
  scope: storageAccount 
  name: guid(storageAccount.id, containerGroup.name, roles.storageAccountContributor)
  properties: {
    roleDefinitionId: resourceId('Microsoft.Authorization/roleDefinitions', roles.storageAccountContributor)
    principalId: containerGroup.identity.principalId
    principalType: 'ServicePrincipal'
  }
}

resource serviceBusDataOwnerAssignment 'Microsoft.Authorization/roleAssignments@2020-04-01-preview' = {
  scope: serviceBus
  name: guid(serviceBus.id, containerGroup.name, roles.serviceBusDataOwner)
  properties: {
    roleDefinitionId: resourceId('Microsoft.Authorization/roleDefinitions', roles.serviceBusDataOwner)
    principalId: containerGroup.identity.principalId
    principalType: 'ServicePrincipal'
  }
}

resource eventGridContributorAssignment 'Microsoft.Authorization/roleAssignments@2020-04-01-preview' = {
  scope: resourceGroup()
  name: guid(containerGroup.id, containerGroup.name, roles.eventGridContributor)
  properties: {
    roleDefinitionId: resourceId('Microsoft.Authorization/roleDefinitions', roles.eventGridContributor)
    principalId: containerGroup.identity.principalId
    principalType: 'ServicePrincipal'
  }
}

resource eventGridDataSenderAssignment 'Microsoft.Authorization/roleAssignments@2020-04-01-preview' = {
  scope: resourceGroup()
  name: guid(containerGroup.id, containerGroup.name, roles.eventGridDataSender)
  properties: {
    roleDefinitionId: resourceId('Microsoft.Authorization/roleDefinitions', roles.eventGridDataSender)
    principalId: containerGroup.identity.principalId
    principalType: 'ServicePrincipal'
  }
}

resource eventGridEventSubscriptionContributorAssignment 'Microsoft.Authorization/roleAssignments@2020-04-01-preview' = {
  scope: resourceGroup()
  name: guid(containerGroup.id, containerGroup.name, roles.eventGridEventSubscriptionContributor)
  properties: {
    roleDefinitionId: resourceId('Microsoft.Authorization/roleDefinitions', roles.eventGridEventSubscriptionContributor)
    principalId: containerGroup.identity.principalId
    principalType: 'ServicePrincipal'
  }
}

resource eventGridEventSubscriptionReaderAssignment 'Microsoft.Authorization/roleAssignments@2020-04-01-preview' = {
  scope: resourceGroup()
  name: guid(containerGroup.id, containerGroup.name, roles.eventGridEventSubscriptionReader)
  properties: {
    roleDefinitionId: resourceId('Microsoft.Authorization/roleDefinitions', roles.eventGridEventSubscriptionReader)
    principalId: containerGroup.identity.principalId
    principalType: 'ServicePrincipal'
  }
}

