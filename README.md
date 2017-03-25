# GIBS

The **G**o **I**CMP **B**ind **S**hell is a remote bindshell tunneled over ICMP. Although, since this is ICMP, we use raw sockets. It technically doesn't "bind" to any port, and requires root privileges to run.

### Usage

Set up the bindshell on a target host:

```
sudo gibs -bind-shell
```

Execute commands remotely:

```
sudo gibs -host 127.0.0.1 -cmd 'uname -a'
```
