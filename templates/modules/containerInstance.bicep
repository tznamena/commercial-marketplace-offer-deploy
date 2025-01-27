
@description('Application version in this format: v1.0.0')
param appVersion string

@description('container image')
param containerImage string

@description('Location for all resources.')
param location string = resourceGroup().location

@description('Port to open on the container')
param port int = 8080

@description('The number of CPU cores to allocate to the container.')
param cpuCores int = 1

@description('The amount of memory to allocate to the container in gigabytes.')
param memoryInGb int = 2

@description('The Azure Tenant Id')
param tenantId string

@description('The Azure Subscription Id')
param subscriptionId string

@description('The Azure Resource Group')
param resourceGroupName string

@description('The Azure service bus name')
param serviceBusNamespace string

@description('The email address used for the acme account')
param acmeEmail string

@description('The public http port')
param publicHttpPort int = 80

@description('The public https port')
param publicHttpsPort int = 8443

@description('The behavior of Azure runtime if container has stopped.')
@allowed([
  'Always'
  'Never'
  'OnFailure'
])
param restartPolicy string = 'Always'

var versionSuffix = replace(appVersion, '.', '')

resource storageAccount 'Microsoft.Storage/storageAccounts@2021-09-01' = {
  name: 'modmstor0${versionSuffix}'
  location: location
  kind: 'StorageV2'
  sku: {
    name: 'Standard_LRS'
  }
}

var fileShareName = '${storageAccount.name}/default/share'
resource fileStore 'Microsoft.Storage/storageAccounts/fileServices/shares@2021-09-01' = {
  name: fileShareName
  properties: {
    shareQuota: 1
    enabledProtocols: 'SMB'
  }
}

var sharedVolumeName = 'filestore'
var fileShareMountPath = '/opt/modm'
var containerName = 'modm-${versionSuffix}'
var readinessFilePath = '${fileShareMountPath}/ready'

resource containerGroup 'Microsoft.ContainerInstance/containerGroups@2022-10-01-preview' = {
  name: 'modm-group-${versionSuffix}'
  location: location
  identity: {
    type: 'SystemAssigned'
  }
  properties: {
    imageRegistryCredentials: [
      {
        password: '<placeholderRegistryPassword>'
        server: 'ghcr.io'
        username: 'gpsuscodewith'
      }
    ] 
    volumes: [
      {
        name: sharedVolumeName
        azureFile: {
          readOnly: false
          shareName: 'share'
          storageAccountName: storageAccount.name
          storageAccountKey: storageAccount.listKeys().keys[0].value
        }
      }
    ]
    containers: [
      {
        name: 'proxy-server'
        properties: {
          image: 'caddy:2-alpine'
          ports: [
            {
              port: 80
              protocol: 'TCP'
            }
            {
              port: 443
              protocol: 'TCP'
            }
          ]
          resources: {
            requests: {
              cpu: cpuCores
              memoryInGB: memoryInGb
            }
          }
          volumeMounts: [
            {
              name: sharedVolumeName
              mountPath: '/etc/caddy'
              readOnly: false
            }
          ]
          environmentVariables: [
            {
              name: 'ACME_ACCOUNT_EMAIL'
              value: acmeEmail
            }
            {
              name: 'MODM_SITE_ADDRESS'
              value: '${containerName}.${location}.azurecontainer.io'
            }
            // the hostname and port will be used to configure caddy and point to modm
            {
              name: 'MODM_LOCAL_HOSTNAME'
              value: 'localhost'
            }
            {
              name: 'MODM_LOCAL_PORT'
              value: '${port}'
            }
          ]
        }
      }
      {
        name: containerName
        properties: {
          image: containerImage
          ports: [
            {
              port: port
              protocol: 'TCP'
            }
          ]
          resources: {
            requests: {
              cpu: cpuCores
              memoryInGB: memoryInGb
            }
          }
          volumeMounts: [
            {
              name: sharedVolumeName
              mountPath: fileShareMountPath
              readOnly: false
            }
          ]
          environmentVariables: [
            {
              name: 'MODM_DB_PATH'
              value: fileShareMountPath
            }
            {
              name: 'MODM_AZURE_TENANT_ID'
              value: tenantId
            }
            {
              name: 'MODM_AZURE_SUBSCRIPTION_ID'
              value: subscriptionId
            }
            {
              name: 'MODM_AZURE_RESOURCE_GROUP'
              value: resourceGroupName
            }
            {
              name: 'MODM_AZURE_LOCATION'
              value: location
            }
            {
              name: 'MODM_AZURE_SERVICEBUS_NAMESPACE'
              value: serviceBusNamespace
            }
            {
              name: 'ACME_ACCOUNT_EMAIL'
              value: acmeEmail
            }
            {
              name: 'MODM_PUBLIC_BASE_URL'
              value: 'https://${containerName}.${location}.azurecontainer.io/'
            }
            {
              name: 'MODM_LOG_FILE_PATH'
              value: fileShareMountPath
            }
            {
              name: 'MODM_READINESS_FILE_PATH'
              value: readinessFilePath
            }
          ]
          readinessProbe: {
            exec: {
              command: [
                'cat'
                readinessFilePath
              ]
            }
            initialDelaySeconds: 120
            periodSeconds: 30
          }
        }
      }
    ]
    osType: 'Linux'
    restartPolicy: restartPolicy
    ipAddress: {
      type: 'Public'
      ports: [
        {
          port: 80
          protocol: 'TCP'
        }
        {
          port: 443
          protocol: 'TCP'
        }
      ]
      dnsNameLabel: containerName
    }
  }
  dependsOn: [
    caddyFileDeploymentScript
  ]
}

var filename = 'caddyFile'
@description('this deployment script supports sucking in the caddyFile for caddy')
resource caddyFileDeploymentScript 'Microsoft.Resources/deploymentScripts@2020-10-01' = {
  name: 'caddyFile'
  location: location
  kind: 'AzureCLI'
  properties: {
    azCliVersion: '2.26.1'
    timeout: 'PT5M'
    retentionInterval: 'PT1H'
    environmentVariables: [
      {
        name: 'AZURE_STORAGE_ACCOUNT'
        value: storageAccount.name
      }
      {
        name: 'AZURE_STORAGE_KEY'
        secureValue: storageAccount.listKeys().keys[0].value
      }
      {
        name: 'CONFIG_FILE_CONTENT'
        value: loadTextContent('../../deployments/caddy/caddyFile')
      }
    ]
    scriptContent: 'echo "$CONFIG_FILE_CONTENT" > ${filename} && az storage file upload --source ${filename} -s share'
  }
}

output storageAccountName string = storageAccount.name
output containerGroupName string = containerGroup.name
output containerIPv4Address string = containerGroup.properties.ipAddress.ip
