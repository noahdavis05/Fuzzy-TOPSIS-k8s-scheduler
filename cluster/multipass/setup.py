## This is a script which sets up a kubernetes cluster with kubeadm on multipass

import subprocess
import time 

CONTROL_PLANE_NAME = "cp"

def run(cmd):
    print("Running:", cmd)
    subprocess.run(cmd, shell=True, check=True)

def get_join_command():
    output = subprocess.check_output(
        f"multipass exec {CONTROL_PLANE_NAME} -- microk8s add-node",
        shell=True,
    ).decode()

    for line in output.splitlines():
        if "microk8s join" in line:
            return line.strip()

    raise Exception("Could not retrieve join command")

# dictionary with VM details
vms = {
    CONTROL_PLANE_NAME : {"cpus": "2", "memory":  "4G", "disk": "15G"},
    "w1": {"cpus": "2", "memory":  "4G", "disk": "15G"},
    "w2": {"cpus": "2", "memory":  "4G", "disk": "15G"}
}

# create these vms
for name, specs in vms.items():
    command = f"multipass launch 24.04 --name {name} --cpus {specs['cpus']} --memory {specs['memory']} --disk {specs['disk']}"
    run(command)


# Install MicroK8s on all nodes
INSTALL_COMMANDS = [
    "sudo snap install microk8s --classic --channel=1.30/stable",
    "sudo usermod -a -G microk8s ubuntu",
    "sudo iptables -P FORWARD ACCEPT"
]

for name in vms.keys():
    for cmd in INSTALL_COMMANDS:
        run(f'multipass exec {name} -- bash -c "{cmd}"')

# pause
time.sleep(30)

# enable required addons on control plane
ADDONS = "dns storage ingress metrics-server"

run(
    f'multipass exec {CONTROL_PLANE_NAME} -- '
    f'bash -c "microk8s status --wait-ready && microk8s enable {ADDONS}"'
)

# join worker nodes to cluster
for name in vms.keys():
    if name != CONTROL_PLANE_NAME:
        join_cmd = get_join_command()
        print(f"{name} joining with: {join_cmd}")
        run(f'multipass exec {name} -- bash -c "sudo {join_cmd}"')

print("Cluster setup complete")

