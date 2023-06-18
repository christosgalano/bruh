targetScope = 'subscription'

/// Parameters ///

@description('Azure region used for the deployment of all resources')
param location string

@description('Abbreviation fo the location')
param location_abbreviation string

@description('Name of the workload that will be deployed')
param workload string

@description('Name of the workloads environment')
param environment string

@description('Tags to be applied on the resource group')
param rg_tags object = {}

/// Variables ///

var tags = union({
    workload: workload
    environment: environment
  }, rg_tags)

var rg_name = 'rg-${suffix}'
var webapp_name = 'app-${suffix}'
var suffix = '${workload}-${environment}-${location_abbreviation}'

/// Resources ///

resource rg 'Microsoft.Resources/resourceGroups@2022-09-01' = {
  name: rg_name
  location: location
  tags: tags
}

module compute 'modules/compute.bicep' = {
  scope: resourceGroup(rg_name)
  name: 'compute-deployment'
  params: {
    location: location

    plan_name: 'plan-${suffix}'
    plan_kind: 'linux'
    plan_reserved: true
    plan_sku_name: 'P1V3'
    plan_sku_tier: 'PremiumV3'
    plan_sku_capacity: 2

    webapp_name: webapp_name
    webapp_always_on: true
    webapp_min_tls_version: '1.2'
    webapp_health_check_path: '/health'
    webapp_linux_fx_version: 'DOCKER|mcr.microsoft.com/appsvc/staticsite:latest'
  }
  dependsOn: [
    rg
  ]
}

/// Outputs ///

output webapp_name string = webapp_name
output resource_groups array = [ rg.name ]
