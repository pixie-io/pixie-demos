# Detect Monero Miners with bpftrace
This demo accompanies the [Detect Monero Miners with bpftrace blogpost](https://blog.px.dev/detect-monero-miners).

## Prereqs
1. (Optional) Deploy k0s
**Cloud Providers have strict policies around cryptomining. We strongly recommend not deploying to a cloud provider or you will risk account deactivation.**
One option is to [deploy k0s](https://docs.k0sproject.io/v1.23.3+k0s.0/install/), exposing the cri to your local docker instance
```sh
# Download k0s
curl -sSLf https://get.k0s.sh | sudo sh
# Install the controll
sudo k0s install controller --single --enable-worker --cri-socket docker:unix:///var/run/docker.sock
# Start k0s
sudo k0s start
#$ Copy the kube config over 
sudo cp /var/lib/k0s/pki/admin.conf admin.conf


# Before you run kubectl commands
export KUBECONFIG=admin.conf
kubectl apply -f mydeployment.yaml
```
2. [Deploy Pixie](https://docs.px.dev/installing-pixie/install-guides/)

## Deploying xmrig to your cluster
**Cloud Providers have strict policies around cryptomining. We strongly recommend not deploying to a cloud provider or you will risk account deactivation.**
For this demo, we deployed the popular open source Monero miner, [XMRig](https://github.com/xmrig/xmrig).
I built the docker image locally. I couldn't find a reliable looking public image.
1. `git clone` this repo and cd into this direcory.
```sh
cd ./detect-monero-demo
```

2. Download xmrig and verify the sha256sum.
```sh
# Download the xmrig binary and verify the sha256sum from
# https://github.com/xmrig/xmrig/releases
# Instructions added for convenience, but please double check shas and download paths.
curl -LO https://github.com/xmrig/xmrig/releases/download/v6.16.4/xmrig-6.16.4-linux-static-x64.tar.gz
# Make sure the grep matches. Double check the sum with the release page.
sha256sum xmrig-6.16.4-linux-static-x64.tar.gz  | grep bf1e10f389d119fe4f72950a6a59bc6a74ba99faa48e5c959edabcdc234ac457
```
3. Unpack the tar file and move the xmrig binary out to this directory.
```sh
tar -xzvf xmrig-6.16.4-linux-static-x64.tar.gz
# Move the binary out of the directory
mv xmrig-6.16.4/xmrig .
```
4. Create a config file using https://xmrig.com/wizard and paste it in `config.json`
5. Build the docker image and apply the kubernetes yamls
```sh
# You might have to change your docker-env to push to your local environment 
docker build . -t xmrig
kubectl apply -f xmrig_deployment.yaml
```

## Running the bpftrace script
### Pixie CLI
```sh
px run -f detectrandomx.pxl
```
### Pixie UI 
Copy and paste the contents of `detectrandomx.pxl` into the [scratchpad](https://docs.px.dev/using-pixie/using-live-ui#write-your-own-pxl-scripts-use-the-scratch-pad).
### bpftrace CLI
[bpftrace install guide](https://github.com/iovisor/bpftrace/blob/master/INSTALL.md)
```sh
sudo bpftrace detectrandomx.bt
```

## Caveats
1. This script only works for x86 processors. There is probably a similar detection opportunity on ARM processors.
2. The script was tested on Linux Kernel version 5.13. You'll have to update this for Linux kernel >=5.16 changed the structure.
3. Minikube virtualizes the CPU so this script won't work inside Pixie running on Minikube. I used [k0s](https://k0sproject.io/). 

