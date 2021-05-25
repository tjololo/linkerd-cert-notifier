# linkerd-cert-notifier
Get notified before linkerd certificates expires.

Currently the application can send notifications to slack webhhoks.

Example message in slack:
![notification](slack-notification.png)

The application could be run as a cronjob in the same namespace as your linkerd install.

See the chart for example deployment.

## Local testing
To test the applicaiton locally run the script ```_localtests/start-test-cluster.sh```.
This will start a kind cluster with local registry and linkerd installed

Build and deploy the application with the command ```make local-deploy```

Delete the test cluster with: ```kind delete cluster --name kind```