# Detect Monero Miners with bpftrace
This demo accompanies the [Detect Monero Miners with bpftrace blogpost]().

## Deploying xmrig to your cluster
**Cloud Providers have strict policies around cryptomining. We strongly recommend not deploying
to a cloud provider or you will risk account deactivation.**

I built my xmrig docker image locally. I couldn't find a public one that looked reliable.

```bash
# cd into this directory
cd ./detect-monero-demo
# Download the xmrig binary and verify the sha256sum from
# https://github.com/xmrig/xmrig/releases
# Instructions added for convenience, but please double check shas and download paths.
curl -LO https://github.com/xmrig/xmrig/releases/download/v6.16.4/xmrig-6.16.4-linux-static-x64.tar.gz
# Make sure the grep matches. Double check the sum with the release page.
sha256sum xmrig-6.16.4-linux-static-x64.tar.gz  | grep bf1e10f389d119fe4f72950a6a59bc6a74ba99faa48e5c959edabcdc234ac457
tar -xzvf xmrig-6.16.4-linux-static-x64.tar.gz
# Move the binary out of the directory
mv xmrig-6.16.4/xmrig .

# Create a config file using https://xmrig.com/wizard and paste it here.
vim config.json

# You might have to change your docker-env to push to your local environment 
docker build . -t xmrig
kubectl apply -f k8s/
```

## Running the bpftrace script
Using the bpftrace CLI

```bash
sudo bpftrace detectrandomx.bt
```

## Running the pxl script
**bpftrace in Pixie is still an alpha feature. APIs might change and features might break.**
### Pixie CLI
```bash
px run -f detectrandomx.pxl
```
### Pixie UI 
1. Copy over the script
2. Open up web UI
3. Navigate to the scratchpad script
4. Paste script
5. Run script

## Caveats
1. This script only works for x86 processors. There is probably a similar detection opportunity
on ARM processors.
2. The script was tested on Linux Kernel version 5.13. You'll have to update this for Linux kernel 
>=5.16 changed the structure.
3. Minikube virtualizes the CPU so this script won't work inside Pixie running on Minikube. I
used [k0s](https://k0sproject.io/). 
