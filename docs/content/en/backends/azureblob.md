---
title: Azure Blob Storage
weight: 13
---

Use [Azure Blob Storage](https://azure.microsoft.com/services/storage/blobs/) as the repository backend.

## Configuration

Before using the backend, you have to setup the credential. There are two methods to configure.

- **Use Azure CLI to login:** Suitable for development environment.
- **Use environment variables:** Suitable for production or CI environment


{{< hint warning >}}
**Assign the Permission**\
The logged-in account requires **Storage Blob Data Contributor** role to the storage account. Assign it in the **Azure Portal**

*Storage Accounts* > *my account* > *Access Control (IAM)* > *Role assignments*

For more information, please see https://docs.microsoft.com/azure/storage/blobs/assign-azure-role-data-access
{{< /hint >}}

### Use Azure CLI to login

This backend suppport to use [Azure CLI](https://docs.microsoft.com/cli/azure/install-azure-cli) to configure the login account. It will open the browser and start the login process. 

```
az login
```

It also supports other login options provided by az login, such as

```
az login --service-principal -u <client-id> -p <client-password> -t <tenant-id>
```

### Use Environment Variables

- Service principal with a secret

    | Name | Description
    | --- | --- |
    AZURE_TENANT_ID	| ID of the application's Azure AD tenant
    AZURE_CLIENT_ID	| Application ID of an Azure service principal
    AZURE_CLIENT_SECRET	| Password of the Azure service principal

- Service principal with certificate

    | Name | Description
    | --- | --- |
    AZURE_TENANT_ID	| ID of the application's Azure AD tenant
    AZURE_CLIENT_ID	| ID of an Azure AD application
    AZURE_CLIENT_CERTIFICATE_PATH	| Path to a certificate file including private key (without password protection)

- Username and password

    | Name | Description
    | --- | --- |
    AZURE_CLIENT_ID	| ID of an Azure AD application
    AZURE_USERNAME	| A username (usually an email address)
    AZURE_PASSWORD	| That user's password

- Managed identity

    [Managed identities](https://docs.microsoft.com/azure/active-directory/managed-identities-azure-resources/overview) eliminate the need for developers to manage credentials. By connecting to resources that support Azure AD authentication, applications can use Azure AD tokens instead of credentials.

    | Name | Description
    | --- | --- |
    AZURE_CLIENT_ID	| User assigned managed identity client id

- Storage account key

    | Name | Description
    | --- | --- |
    AZURE_STORAGE_ACCOUNT_KEY | The access key of the storage account

## Usage

Init a workspace
```shell
avc init https://mystorageaccount.blob.core.windows.net/mycontainer/path/to/mydataset
```

Clone a repository
```shell
avc clone https://mystorageaccount.blob.core.windows.net/mycontainer/path/to/mydataset
cd mydataset/
```
