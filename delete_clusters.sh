echo master用deployment、serviceを削除します
kubectl delete -f jmeter-master.yaml
echo 削除完了

echo slave用deployment、serviceを削除します
kubectl delete -f jmeter-slave.yaml
echo 削除完了

echo クラスタを作成します
gcloud container clusters delete jmeter -y
echo クラスタ
