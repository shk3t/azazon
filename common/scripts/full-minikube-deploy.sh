#!/usr/bin/env bash

scriptpath=$(dirname $0)
basepath=$(dirname $(dirname $scriptpath))
source $scriptpath/envs.sh

cat << END
Before installation:
1. Follow \"Deploy Strimzi using installation files\" guidelines for Kafka pre-installation: \"https://strimzi.io/quickstarts/\".
2. Rename \`example.env\`->\`.env\` and fill accordingly.
END

if [[ ! -e "${basepath}/.env" ]]; then
    echo "Can not find \`.env\` file"
    exit 1
fi

kubectl apply -f $basepath/common/deployments/namespace.yaml
kubectl apply -f $basepath/common/deployments/gatewayclass.yaml

bash $scriptpath/gen-proto.sh
bash $scriptpath/sync-helm-env.sh

helm uninstall -n $namespace $baseapp --ignore-not-found
helm install -n $namespace $baseapp $basepath/common/deployments/$baseapp || exit 1
kubectl wait -n $namespace kafka/$kafka_cluster_name --for=condition=Ready --timeout=300s || exit 1
bash $scriptpath/recreate-kafka-topics.sh --minikube || exit 1

bash $scriptpath/rebuild-app-images.sh || exit 1
bash $scriptpath/reinstall-database-charts.sh || exit 1
bash $scriptpath/reinstall-app-charts.sh || exit 1

echo
kubectl get pods -n $namespace
echo
cat << END
Use \`minikube tunnel\` to get access to the deployed cluster.
Pass gateway IP address from \`kubectl get services -n istio-ingress\` to \`.env\`.
Use \`common/scripts/read-all-logs.sh --minikube\` to read each service log file.
END
