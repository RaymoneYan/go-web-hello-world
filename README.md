Task 0: Install a ubuntu 16.04 server 64-bit

  1. install VirtualBox
  2. download ubuntu-16.04.6-server-amd64.iso
  3. create a new vm in VirtualBxo
  10.210.149.231

Task 1: Update system

  https://www.howtoforge.com/tutorial/how-to-upgrade-linux-kernel-in-ubuntu-1604-server/  
  1. sudo apt update
  2. sudo apt upgrade -y
  3. sudo reboot
  4. sudo apt list --upgradeable
  5. uname -msr  
     Linux 4.4.0-170-generic x86_64
  6. sudo mkdir -p ~/4.11.2
     cd ~/4.11.2
  7. wget http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.11.2/linux-headers-4.11.2-041102_4.11.2-041102.201705201036_all.deb
     wget http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.11.2/linux-headers-4.11.2-041102-generic_4.11.2-041102.201705201036_amd64.deb
     wget http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.11.2/linux-image-4.11.2-041102-generic_4.11.2-041102.201705201036_amd64.deb
  8. dpkg -i *.deb
  9. sudo update-grub
  10. sudo reboot
  11. uname -msr
      Linux 4.11.2-041102-generic x86_64
  Remove the old kernel
  12. sudo apt install byobu
  13. dpkg -l | grep linux-image
  14. sudo purge-old-kernels
  15. No kernels are eligible for removal
  16. purge-old-kernels --keep 1 -q
  17. sudo update-grub

Task 2: install gitlab-ce version in the host

  1. sudo apt-get install -y curl openssh-server ca-certificates
  2. curl -sS https://packages.gitlab.com/install/repositories/gitlab/gitlab-ce/script.deb.sh | sudo bash
  3. sudo EXTERNAL_URL="https://127.0.0.1" apt-get install gitlab-ce
  4. ³ö´í£¬Ð¶ÔØ¡£
  5. gitlab-ctl stop
  6. ps -ef | grep gitlab
  7. kill -9 xxx
  8. cd /
  9. find / -name *gitlab* | xargs rm -rf
  10. sudo EXTERNAL_URL="http://127.0.0.1" apt-get install gitlab-ce
  when you restart your VM, please run the command gitlab-ctl start to start the gitlab.

Task 3: create a demo group/project in gitlab

  1. Create a group called demo, then create a project called go-web-hello-world.
  2. git clone http://127.0.0.1:8080/demo/go-web-hello-world.git
  3. refer to rick's repo to create go file.
  4. git push the repo

Task 4: build the app and expose ($ go run) the service to 8081 port

  2. go build main.go
  3. ./main
  4. curl http://127.0.0.1:8081  

Task 5: install docker

  1. Reference: https://docs.docker.com/install/linux/docker-ce/ubuntu/
  2.  sudo apt install docker-ce=5:18.09.9~3-0~ubuntu-xenial docker-ce-cli=5:18.09.9~3-0~ubuntu-xenial 

Task 6: run the app in container

  1. docker build -t go-web-hello-world .¡¢
  2. docker images
  3. netstat -anp | grep 8082
  4. vi /etc/gitlab/gitlab.rb
  5. sidekiq['listen_port'] = 8787
  6. gitlab-ctl reconfigure
  7. netstat -anp | grep 8082
  8. docker run -p 8082:8081 -it --rm --name run-web-hello-world go-web-hello-world

Task 7: push image to dockerhub

  1. docker tag go-web-hello-world:latest raymoneyan126/go-web-hello-world:v0.1
  2. docker login 
  3. docker push  raymoneyan126/go-web-hello-world:v0.1
  4. Repo : https://hub.docker.com/repository/docker/raymoneyan126/go-web-hello-world

Task 8: document the procedure of step 0-7 in a MarkDown file

Done

Task 9: install a single node Kubernetes cluster using kubeadm

  1. do it in ce-tu node-10-210-149-231
  2. sudo apt-get update && sudo apt-get install -y apt-transport-https curl
     curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
     cat <<EOF | sudo tee /etc/apt/sources.list.d/kubernetes.list
     deb https://apt.kubernetes.io/ kubernetes-xenial main
     EOF
     sudo apt-get update
     sudo apt-get install -y kubelet kubeadm kubectl
  3. intall containerd https://kubernetes.io/docs/setup/production-environment/container-runtimes/
  4. ---intit cluster---
     kubeadm init
     mkdir -p $HOME/.kube
     sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
     sudo chown $(id -u):$(id -g) $HOME/.kube/config
  5. kubectl get node
     kubectl get po -A
  6. kubectl taint nodes --all node-role.kubernetes.io/master-

Task 10: deploy the hello world container

  1. kubectl run go-web-hello-world --image raymoneyan126/go-web-hello-world:v0.1  --replicas=1
  2. vi go-web-hello-world.yaml
      apiVersion: v1
      kind: Service
      metadata:
        creationTimestamp: "2020-01-02T08:03:30Z"
        labels:
          run: go-web-hello-world
        name: go-web-hello-world
        namespace: default
        resourceVersion: "7921"
        selfLink: /api/v1/namespaces/default/services/go-web-hello-world
        uid: 1f3726f2-521c-43b9-bbc6-a8b719db2d8d
      spec:
        clusterIP: 10.96.251.104
        externalTrafficPolicy: Cluster
        ports:
        - nodePort: 31222
          port: 8081
          protocol: TCP
          targetPort: 8081
        selector:
          run: go-web-hello-world
        sessionAffinity: None
        type: NodePort
      status:
        loadBalancer: {}
  3. kubectl apply -f go-web-hello-world.yaml
  4. http://10.210.149.231:31222/
Task 11: install kubernetes dashboard and expose the service to nodeport 31081
  
  1. kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0-beta8/aio/deploy/recommended.yaml
  2. kubectl get svc -n kubernetes-dashboard
  3. kubectl edit svc kubernetes-dashboard  -n kubernetes-dashboard    
    spec:
      clusterIP: 10.96.166.28
      ports:
      - port: 443
        protocol: TCP
        targetPort: 8443
        nodePort: 31081   
      selector:
        k8s-app: kubernetes-dashboard
      sessionAffinity: None
      type: NodePort
  4. mkdir certs
  5. openssl req -nodes -newkey rsa:2048 -keyout certs/dashboard.key -out certs/dashboard.csr -subj "/C=/ST=/L=/O=/OU=/CN=kubernetes-dashboard"
  6. openssl x509 -req -sha256 -days 365 -in certs/dashboard.csr -signkey certs/dashboard.key -out certs/dashboard.crt
  7. kubectl delete secret/kubernetes-dashboard-certs -n kubernetes-dashboard
  8. kubectl create secret generic kubernetes-dashboard-certs --from-file=certs -n kubernetes-dashboard
  9. kubectl get po -n kubernetes-dashboard
  10. kubectl delete po/kubernetes-dashboard-5996555fd8-r4zjx  -n kubernetes-dashboard
  11. kubectl get po -n kubernetes-dashboard

Task 12: generate token for dashboard login in task 11

  1. https://github.com/kubernetes/dashboard/blob/master/docs/user/access-control/creating-sample-user.md
  2. vi SA.yaml
  3. vi CRB.yaml
  4. kubectl -n kubernetes-dashboard describe secret $(kubectl -n kubernetes-dashboard get secret | grep admin-user | awk '{print $1}')
  5. copy the token to the dashboard.

  


     






















