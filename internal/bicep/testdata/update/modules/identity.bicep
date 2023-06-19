/// Parameters ///

@minLength(3)
@maxLength(128)
@description('Name of the user managed identity')
param name string

@description('Location of the user managed identity')
param location string

/// Resources ///

resource identity 'Microsoft.ManagedIdentity/userAssignedIdentities@2022-01-31-preview' = {
  name: name
  location: location
}

/// Outputs ///

output resource_id string = identity.id
output client_id string = identity.properties.clientId
output principal_id string = identity.properties.principalId
output tenant_id string = identity.properties.tenantId
