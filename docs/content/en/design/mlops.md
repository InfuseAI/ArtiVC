---
title: MLOps    
weight: 5
---

Machine learning application is a data-driven approach to solving problems. The MLOps is a methodology to has a feedback loop from data, model training, evaluation, model publish, deploy, monitoring. There are three core components in an ML application.

- Code
- Datasets
- Models

The code can be training code, application code. Mostly, it is versioned by git and we have been familiar with the way to version it. And we also use git as the single source of truth to drive the whole DevOps lifecycle.

However, for datasets and models, there is still no defacto solution to version them. Usually, these data are stored in cloud object storage, on-premise object storage like MinIO, or NFS. There is still a gap between data storage and version metadata storage. Here is why we would like to build the *ArtiVC*.

In addition, we are thinking about how to drive the automation when an artifact store event is triggered. In git, we can trigger a job whenever a git event. In the artifact store, we lack the fundamentals to trigger this event. *ArtiVC* reference the *git* design and provide the *commits* and *references* primitives to make it possible to define a commit or a version that is created. It makes it possible to listen to the object storage or the file system event to trigger an automation job accordingly.