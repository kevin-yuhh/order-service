# Soter Order Service

> BTFS soter project - order service

The Soter order service is a order service in the btfs gateway. The controller receives the user request, calls the order interface to calculate the fee and completes the deduction.

## Table of Contents

- [Security](#security)
- [Background](#background)
- [Install](#install)
- [Deploy](#deploy)
- [API](#api)

## Security

The order server and the controller are in the same network segment and the external network is invisible, so the order server is secure.

## Background

BTFS stores pictures that consume node memory and network bandwidth, which comes at a price. Therefore, a billing service is required to request a fee from the user.

## Install

```
cd $PROJECT_HOME
GOOS=linux GOARCH=amd64 go build
```

## Deploy

### Prerequisites

1. Set max open files

   ```shell
   # for once
   ulimit -n 65535
   
   # permanent
   echo '*  -  nofile  65535' >> /etc/security/limits.conf
   shutdown -r now
   ```

2. Adjust linux kernel parameters

   ```shell
   # for once
   net.ipv4.tcp_tw_reuse = 1
   net.ipv4.tcp_tw_recycle = 1
   
   # permanent
   echo 'net.ipv4.tcp_tw_reuse = 1' >> /etc/sysctl.conf
   echo 'net.ipv4.tcp_tw_recycle = 1' >> /etc/sysctl.conf
   /sbin/sysctl -p
   ```



### Compile & Deploy

```
scp soter-order-service target
scp config.yml target
nohup soter-order-service -d (config path) -n (config name) &
```

## API

* QueryBalance
  * Request Params
    * address
  * Response Params
    * balance
* CreateOrder
  * Request Params
    * address
    * request_id
    * file_name
    * file_size
  * Response Params
    * order_id
* SubmitOrder
  * Request Params
    * order_id
    * file_hash
  * None Response Params
* CloseOrder
  * Request Params
    * order_id
  * None Response Params