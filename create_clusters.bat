@echo off

echo クラスタを作成します
gcloud container clusters create jmeter --preemptible --machine-type=g1-small --num-nodes=3 --disk-size=10 --zone=asia-northeast1-a --enable-basic-auth --issue-client-certificate --no-enable-ip-alias --metadata disable-legacy-endpoints=true
echo クラスタ作成完了

echo クラスタに認証を通します
gcloud container clusters get-credentials jmeter
echo 完了

echo namespaceを作成します
kubectl create namespace jmslave
kubectl create namespace jmmaster
echo 完了

echo slave用deployment、serviceを設定します
kubectl apply -f jmeter-slave.yaml
echo 設定完了

echo master用deployment、serviceを設定します
kubectl apply -f jmeter-master.yaml
echo 設定完了