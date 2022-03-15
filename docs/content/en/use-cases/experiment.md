---
title: ML Experiments
weight: 2
---
Here we use three repositories
- Dataset for training
- A base model for transfer learning of fine-tuning
- Experiment output. Versioned by the timestamp of an experiment.


In this use case, we use `get` and `put` commands to simplify the commands for programmatic use cases.

1. Clone the training code
   ```shell
   git clone https://github.com/my-org/my-ml-project.git
   cd my-ml-project
   ```
1. Download the dataset and the base model
   ```shell
   art get -o dataset/ s3://mybucket/datasets/flowers-classification@v0.1.0
   art get -o base/ s3://mybucket/models/my-base-model@v0.3.0
   ```   
1. Train and output your training result (trained model, experiment log, hyperparams, etc) to `artifacts/` folder
   ```shell
   python ./train.py
   ``` 
1. Upload the artifacts 
   ```shell
   art put artifacts/ s3://mybucket/experiments/project1@202220303-100504
   ```
