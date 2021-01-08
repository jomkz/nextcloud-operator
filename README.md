# NextCloud Operator

A Kubernetes operator for managing Nextcloud clusters.

## Development

Clone this repo somewhere and navigate to the directory. 

Run the following command to build the operator and run tests.

```
make test
```

Build and upload the container image for the new build of the operator.

```
make img-build img-push
```

Deploy the operator to the kubernetes cluster configured in `~/.kube/config`.

```
make deploy
```

Run the following command to remove the operator when done.

```
make undeploy
```
