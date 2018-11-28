@echo off

echo �N���X�^���쐬���܂�
gcloud container clusters create jmeter --preemptible --machine-type=g1-small --num-nodes=3 --disk-size=10 --zone=asia-northeast1-a --enable-basic-auth --issue-client-certificate --no-enable-ip-alias --metadata disable-legacy-endpoints=true
echo �N���X�^�쐬����

echo �N���X�^�ɔF�؂�ʂ��܂�
gcloud container clusters get-credentials jmeter
echo ����

echo namespace���쐬���܂�
kubectl create namespace jmslave
kubectl create namespace jmmaster
echo ����

echo slave�pdeployment�Aservice��ݒ肵�܂�
kubectl apply -f jmeter-slave.yaml
echo �ݒ芮��

echo master�pdeployment�Aservice��ݒ肵�܂�
kubectl apply -f jmeter-master.yaml
echo �ݒ芮��