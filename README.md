# Kubewatch for Research

### Installing

Git clone this repo anywhere.

### Running

Do these steps regardless of in or out of cluster. First have a cluster up and running with kops. Then in this repo run

    python server.py
    ./ngrok http 9000

Grab the http or https forwarding address. It should look like `http://6998c493.ngrok.io`. 

#### In Cluster

Open the file `kubewatch-configmap.yaml` and replace the word `<url>` with the forwarding address. Then run

    kubectl create -f kubewatch-configmap.yaml
    kubectl create -f kubewatch.yaml

Once the pod is up and running (~30 sec), you should be able to visit the forwarding address and see events.

#### Out of Cluster

    export KW_URL='http://6998c493.ngrok.io'
    ./kubewatch

### Building

To modify kubewatch, first clone the repo and then edit the go code. 

#### In Cluster
To build for in cluster, create a kubewatch repo on dockerhub. Then run

    make binary-image
    docker tag kubewatch yourdockername/kubewatch
    docker push yourdockername/kubewatch
    
Then edit the image in `kubewatch.yaml` to point at `yourdockername/kubewatch`

#### Out of Cluster
Move the kubewatch repo to `$GOPATH/src/github.com/skippbox/kubewatch`. Then run `make build` and it will create the binary image.

### TODO

The most important thing is to find a better way to run the receiving server. The current ngrok solution is a hack
and because it is free tier it is limited to 20 GET/POST requests per minute. Ideas include running on quilt or as
part of the kubernetes cluster. 

-----------------------------------------------------

# Kubewatch
[![Build Status](https://travis-ci.org/skippbox/kubewatch.svg?branch=master)](https://travis-ci.org/skippbox/kubewatch) [![Join us on Slack](https://s3.eu-central-1.amazonaws.com/ngtuna/join-us-on-slack.png)](https://skippbox.herokuapp.com)

`kubewatch` is a Kubernetes watcher that currently publishes notification to Slack. Run it in your k8s cluster, and you will get event notifications in a slack channel.

## Run kubewatch in a Kubernetes cluster

In order to run kubewatch in a Kubernetes cluster quickly, the easiest way is for you to create a [ConfigMap](https://github.com/skippbox/kubewatch/blob/master/kubewatch-configmap.yaml) to hold kubewatch configuration. It contains a SLACK API token, channel.

An example is provided at [`kubewatch-configmap.yaml`](https://github.com/skippbox/kubewatch/blob/master/kubewatch-configmap.yaml), do not forget to update your own slack channel and token parameters. Alternatively, you could use secrets.

Create k8s configmap:

```console
$ kubectl create -f kubewatch-configmap.yaml
```
Create the [Pod](https://github.com/skippbox/kubewatch/blob/master/kubewatch.yaml) directly, or create your own deployment:

```console
$ kubectl create -f kubewatch.yaml
```

A `kubewatch` container will be created along with `kubectl` sidecar container in order to reach the API server.

Once the Pod is running, you will start seeing Kubernetes events in your configured Slack channel. Here is a screenshot:

![slack](./docs/slack.png)

To modify what notifications you get, update the `kubewatch` ConfigMap and turn on and off (true/false) resources:

```
resource:
      deployment: false
      replicationcontroller: false
      replicaset: false
      daemonset: false
      services: true
      pod: true
```

## Building

### Building with go

- you need go v1.5 or later.
- if your working copy is not in your `GOPATH`, you need to set it accordingly.

```console
$ go build -o kubewatch main.go
```

You can also use the Makefile directly:

```console
$ make build
```

### Building with Docker

Buiding builder image:

```console
$ make builder-image
```

Using the `kubewatch-builder` image to build `kubewatch` binary:

```console
$ make binary-image
$ docker images
REPOSITORY          TAG                 IMAGE ID            CREATED              SIZE
kubewatch           latest              f1ade726c6e2        31 seconds ago       33.08 MB
kubewatch-builder   latest              6b2d325a3b88        About a minute ago   514.2 MB
```

## Download kubewatch package

```console
$ go get -u github.com/skippbox/kubewatch
```

## Configuration
Kubewatch supports `config` command for configuration. Config file will be saved at $HOME/.kubewatch.yaml

### Configure slack

```console
$ kubewatch config slack --channel <slack_channel> --token <slack_token>
```

### Configure resources to be watched

```console
// rc, po and svc will be watched
$ kubewatch config resource --rc --po --svc

// only svc will be watched
$ kubewatch config resource --svc
```

### Environment variables
You have an altenative choice to set your SLACK token, channel via environment variables:

```console
$ export KW_SLACK_TOKEN='XXXXXXXXXXXXXXXX'
$ export KW_SLACK_CHANNEL='#channel_name'
```

### Run kubewatch locally

```console
$ kubewatch
```
