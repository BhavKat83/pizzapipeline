
# Lambda code for Pizza Meme Generator

## Introduction

This repository has the code for a  AWS Lamda based meme generator in golang.

The repository is integrated with AWS CodePipeline that can compile and deploy it to the designated function.

## Contents

* gomeme.go - main file with Lambda handlers
* demo.html - frontend HTML hosted on S3 that triggers the Lambda function through AWS API Gateway
* meme.png - image used in the frontend also packaged with the hosting on s3

