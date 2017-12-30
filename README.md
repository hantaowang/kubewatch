# Kubewatch for Research

### Installing

Git clone this repo anywhere.

### Running

First have a cluster up and running with kops. Then in this repo run

    python server.py
    ./ngrok http 9000

Grab the http or https forwarding address. It should look like `http://6998c493.ngrok.io`. Open the file
`kubewatch-configmap.yaml` and replace the word `<url>` with the forwarding address. Then run

    kubectl create -f kubewatch-configmap.yaml
    kubectl create -f kubewatch.yaml

Once the pod is up and running (~30 sec), you should be able to visit the forwarding address and see events.

### TODO

The most important thing is to find a better way to run the receiving server. The current ngrok solution is a hack
and because it is free tier it is limited to 20 GET/POST requests per minute. Ideas include running on quilt or as
part of the kubernetes cluster. 