# clusterawsadm

Utility to manage IAM objects for Kubernetes Cluster API Provider AWS.  
The real `clusterawsadm` uses CloudFormation and lives [here](https://github.com/kubernetes-sigs/cluster-api-provider-aws/tree/master/cmd/clusterawsadm).  
This is a proof of concept implementation using [h0tbird/terrago](https://github.com/h0tbird/terrago).

#### Development
Run this to use the latest `terrago` code:
```
go get -u github.com/h0tbird/terrago@master
```