---
title: Run ML Experiments
weight: 2
---
Here we use three repositories
- Dataset for training
- A base model for transfer learning of fine-tuning
- Experiment output. Versioned by the timestamp of an experiment.
In this use case, we use `get` and `put` commands to simplify the commands for programmatic use cases.

1. Clone the training code
   ```
   git clone https://github.com/my-org/my-ml-project.git
   cd my-ml-project
   ```
1. Download the dataset
   ```
   art get -o dataset/ s3://mybucket/datasets/flowers-classification@v0.1.0
   ```
1. Download the base model   
   ```
   art get -o base/ s3://mybucket/models/my-base-model@v0.3.0
   ```   
1. Train and output your training result (trained model, experiment log, hyperparams, etc) to `artifacts/` folder
   ```
   python ./train.py
   ``` 
1. Upload the artifacts 
   ```
   art put artifacts/ s3://mybucket/experiments/project1@202220303-100504
   ```


## Case 3: Expose the dataset by HTTP
1. [Make your S3 bucket public](https://aws.amazon.com/premiumsupport/knowledge-center/read-access-objects-s3-bucket/?nc1=h_ls)
1. Copy the public URL of your repository. For example
   ```
   https://mybucket.s3.ap-northeast-1.amazonaws.com/datasets/flowers-classification
   ```
1. Download the data
   ```
   art get -o /tmp/output https://mybucket.s3.ap-northeast-1.amazonaws.com/datasets/flowers-classification
   ```   