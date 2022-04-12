---
title: Google Cloud Storage
weight: 12
---

{{< toc >}}

Use [Google Cloud Storage (GCS)](https://cloud.google.com/storage) as the repository backend.

Note that Google Cloud Storage is not [Google Drive](https://www.google.com.tw/drive/). They are different google product.

## Configuration

Before using the backend, you have to configure the service account credential. There are three method to configure it.

1. Use application default credentials. It is recommended way to use in your development environment.

    ```
    gcloud auth application-default login  
    ```

    It will open the browser and start the login process.

1. Use service account credentials. It is recommended way to use in CI, job, or production environment.

   ```
   export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account-credentials.json
   ```

   to get this json file, please see the [Passing credential manually](https://cloud.google.com/docs/authentication/production#manually) document

1. Use the service account in the GCP resources (e.g. GCE, GKE). It is recommended way if the `ArtiVC` is run in the GCP environment. Please see [default service accounts](https://cloud.google.com/iam/docs/service-accounts#default) document


The GCS backend finds credentials by a default procedure defined by [Google Cloud](https://cloud.google.com/docs/authentication/production)



## Usage

Init a workspace
```shell
avc init gs://mybucket/path/to/mydataset
```

Clone a repository
```shell
avc clone gs://mybucket/path/to/mydataset
cd mydataset/
```


## Environment Variables

| Name | Description | Default value |
| --- | --- | --- |
| `GOOGLE_APPLICATION_CREDENTIALS` | The location of service account keys in JSON |  |