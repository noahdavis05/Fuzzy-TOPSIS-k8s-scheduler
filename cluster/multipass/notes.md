# Multipass Cluster

- This setup uses multipass and microk8s to set up a multi node Kubernetes cluster on one machine.
- Multipass is used to create multiple Ubuntu VMs, and each node is set up with Microk8s and joined into the same cluster.
- This is done in the script `setup.py`. You can delete everything and start again in `purge.py`.
- Once the setup has been complete install the proemtheus grafana stack using `microk8s enable prometheus`.

**TO-DO** - add the commands to port forward from the multipass VM to browser to access grafana.