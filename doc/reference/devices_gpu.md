(devices-gpu)=
# Type: `gpu`

```{youtube} https://www.youtube.com/watch?v=T0aV2LsMpoA
```

GPU devices make the specified GPU device or devices appear in the instance.

```{note}
For containers, a `gpu` device may match multiple GPUs at once.
For VMs, each device can match only a single GPU.
```

The following types of GPUs can be added using the `gputype` device option:

- [`physical`](gpu-physical) (container and VM): Passes an entire GPU through into the instance.
  This value is the default if `gputype` is unspecified.
- [`mdev`](gpu-mdev) (VM only): Creates and passes a virtual GPU through into the instance.
- [`mig`](gpu-mig) (container only): Creates and passes a MIG (Multi-Instance GPU) through into the instance.
- [`sriov`](gpu-sriov) (VM only): Passes a virtual function of an SR-IOV-enabled GPU into the instance.

The available device options depend on the GPU type and are listed in the tables in the following sections.

(gpu-physical)=
## `gputype`: `physical`

```{note}
The `physical` GPU type is supported for both containers and VMs.
It supports hotplugging only for containers, not for VMs.
```

A `physical` GPU device passes an entire GPU through into the instance.

### Device options

GPU devices of type `physical` have the following device options:

% Include content from [../metadata.txt](../metadata.txt)
```{include} ../metadata.txt
    :start-after: <!-- config group device-gpu-physical-device-conf start -->
    :end-before: <!-- config group device-gpu-physical-device-conf end -->
```

### Configuration examples

Add all GPUs from the host system as a `physical` GPU device to an instance:

    lxc config device add <instance_name> <device_name> gpu gputype=physical

Add a specific GPU from the host system as a `physical` GPU device to an instance by specifying its PCI address:

    lxc config device add <instance_name> <device_name> gpu gputype=physical pci=<pci_address>

Add a specific GPU from the host system as a `physical` GPU device to an instance using the [Container Device Interface](https://github.com/cncf-tags/container-device-interface) (CDI) notation.

    lxc config device add <instance_name> <device_name> gpu gputype=physical id=<fq_CDI_name>

See {ref}`instances-configure-devices` for more information.

#### Passing an NVIDIA iGPU to a container

If you have NVIDIA runtime libraries and configuration installed on your host and you want to pass these files to your container, you should add a device with the CDI notation.
Consider the following scenario:

Your host is an NVIDIA single board computer that has a Tegra SoC with an iGPU, and you have an SDK for handling AI workloads installed on the host. You want to create a LXD container and run an inference job inside the container using the iGPU as a backend. This job should run inside a Docker container (or another OCI-compliant runtime).
To achieve this, complete the following steps:

1. Initialize a LXD container:

       lxc init ubuntu:24.04 t1 --config security.nesting=true --config security.privileged=true

1. Add an iGPU device to your container:

       lxc config device add t1 igpu0 gpu gputype=physical id=nvidia.com/igpu=0

1. Create a `cloud-init` script that installs the Docker runtime, the [NVIDIA Container Toolkit](https://github.com/NVIDIA/nvidia-container-toolkit), and a script to run a test [TensorRT](https://github.com/NVIDIA/TensorRT) workload:

       ```yaml
         #cloud-config
         package_update: true
         packages:
           - docker.io
         write_files:
           - path: /etc/docker/daemon.json
             permissions: '0644'
             owner: root:root
             content: |
               {
                 "max-concurrent-downloads": 12,
                 "max-concurrent-uploads": 12,
                 "runtimes": {
                 "nvidia": {
                     "args": [],
                     "path": "nvidia-container-runtime"
                   }
                 }
               }
           - path: /root/run_tensorrt.sh
             permissions: '0755'
             owner: root:root
             content: |
               #!/bin/bash
               echo "OS release,Kernel version"
               (. /etc/os-release; echo "${PRETTY_NAME}"; uname -r) | paste -s -d,
               echo
               nvidia-smi -q
               echo
               exec bash -o pipefail -c "
               cd /workspace/tensorrt/samples
               make -j4
               cd /workspace/tensorrt/bin
               ./sample_onnx_mnist
               retstatus=\${PIPESTATUS[0]}
               echo \"Test exited with status code: \${retstatus}\" >&2
               exit \${retstatus}
               "
         runcmd:
           - systemctl start docker
           - systemctl enable docker
           - usermod -aG docker root
           - curl -fsSL https://nvidia.github.io/libnvidia-container/gpgkey | gpg --dearmor -o /usr/share/keyrings/nvidia-container-toolkit-keyring.gpg
           - curl -s -L https://nvidia.github.io/libnvidia-container/stable/deb/nvidia-container-toolkit.list | sed 's#deb https://#deb [signed-by=/usr/share/keyrings/nvidia-container-toolkit-keyring.gpg] https:// #g' | tee /etc/apt/sources.list.d/nvidia-container-toolkit.list
           - sed -i -e '/experimental/ s/^#//g' /etc/apt/sources.list.d/nvidia-container-toolkit.list
           - apt-get update
           - DEBIAN_FRONTEND=noninteractive apt-get install -y nvidia-container-toolkit
           - nvidia-ctk runtime configure
           - systemctl restart docker
       ```

1. Apply this `cloud-init` setup to your instance:

       lxc config set t1 cloud-init.user-data - < cloud-init.yml

1. Start the instance:

       lxc start t1

1. Wait for the `cloud-init` process to finish:

       lxc exec t1 -- cloud-init status --wait

You can now run your inference job inside the LXD container.

```{note}
The NVIDIA runtime must be enabled with CSV mode.
Therefore, you must modify the `mode` of the NVIDIA Container Runtime inside the LXD container to the value `csv` and not `auto` when running Docker.
You can do this by editing the `/etc/nvidia-container-runtime/config.toml` configuration file:

       lxc shell t1
       root@t1 # docker run --gpus all --runtime nvidia --rm -v $(pwd):/sh_input nvcr.io/nvidia/tensorrt:24.02-py3-igpu bash /sh_input/run_tensorrt.sh
```

(gpu-mdev)=
## `gputype`: `mdev`

```{note}
The `mdev` GPU type is supported only for VMs.
It does not support hotplugging.
```

An `mdev` GPU device creates and passes a virtual GPU through into the instance.
You can check the list of available `mdev` profiles by running [`lxc info --resources`](lxc_info.md).

### Device options

GPU devices of type `mdev` have the following device options:

% Include content from [../metadata.txt](../metadata.txt)
```{include} ../metadata.txt
    :start-after: <!-- config group device-gpu-mdev-device-conf start -->
    :end-before: <!-- config group device-gpu-mdev-device-conf end -->
```

### Configuration examples

Add an `mdev` GPU device to an instance by specifying its `mdev` profile and the PCI address of the GPU:

    lxc config device add <instance_name> <device_name> gpu gputype=mdev mdev=<mdev_profile> pci=<pci_address>

See {ref}`instances-configure-devices` for more information.

(gpu-mig)=
## `gputype`: `mig`

```{note}
The `mig` GPU type is supported only for containers.
It does not support hotplugging.
```

A `mig` GPU device creates and passes a MIG compute instance through into the instance.
Currently, this requires NVIDIA MIG instances to be pre-created.

### Device options

GPU devices of type `mig` have the following device options:

% Include content from [../metadata.txt](../metadata.txt)
```{include} ../metadata.txt
    :start-after: <!-- config group device-gpu-mig-device-conf start -->
    :end-before: <!-- config group device-gpu-mig-device-conf end -->
```

You must set either {config:option}`device-gpu-mig-device-conf:mig.uuid` (NVIDIA drivers 470+) or both {config:option}`device-gpu-mig-device-conf:mig.ci` and {config:option}`device-gpu-mig-device-conf:mig.gi` (old NVIDIA drivers).

### Configuration examples

Add a `mig` GPU device to an instance by specifying its UUID and the PCI address of the GPU:

    lxc config device add <instance_name> <device_name> gpu gputype=mig mig.uuid=<mig_uuid> pci=<pci_address>

See {ref}`instances-configure-devices` for more information.

(gpu-sriov)=
## `gputype`: `sriov`

```{note}
The `sriov` GPU type is supported only for VMs.
It does not support hotplugging.
```

An `sriov` GPU device passes a virtual function of an SR-IOV-enabled GPU into the instance.

### Device options

GPU devices of type `sriov` have the following device options:

% Include content from [../metadata.txt](../metadata.txt)
```{include} ../metadata.txt
    :start-after: <!-- config group device-gpu-sriov-device-conf start -->
    :end-before: <!-- config group device-gpu-sriov-device-conf end -->
```

### Configuration examples

Add a `sriov` GPU device to an instance by specifying the PCI address of the parent GPU:

    lxc config device add <instance_name> <device_name> gpu gputype=sriov pci=<pci_address>

See {ref}`instances-configure-devices` for more information.
