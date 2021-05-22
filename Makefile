local-deploy:
	KO_DOCKER_REPO=kind.local  ko apply -f _localtests/deployment.yaml

local-cleanup:
	kubectl delete -f _localtests/deployment.yaml