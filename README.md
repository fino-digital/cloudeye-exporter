# cloudeye-exporter

Prometheus cloudeye exporter for [Huaweicloud](https://www.huaweicloud.com/) compatible cloud providers

## Download

```sh
git clone https://github.com/huaweicloud/cloudeye-exporter
```

## (Option) Building The Discovery with Exact steps on clean Ubuntu 16.04

```sh
wget https://dl.google.com/go/go1.12.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.12.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin # You should put in your .profile or .bashrc
go version # to verify it runs and version #

go get github.com/huaweicloud/cloudeye-exporter
cd ~/go/src/github.com/huaweicloud/cloudeye-exporter
go build
```

## Usage

```sh
./cloudeye-exporter  -config=clouds.yml
```

The default port is 8087, default config file location is ./clouds.yml.

Visit metrics in <http://localhost:8087/metrics?services=SYS.ELB,SYS.RDS,SYS.DCS,SYS.VPC,SYS.EVS,SYS.ECS,SYS.AS,SYS.NAT>

## Help

```sh
Usage of ./cloudeye-exporter:
  -config string
        Path to the cloud configuration file (default "./clouds.yml")
  -debug
        If debug the code.
```

## Example of config file(clouds.yml)

The "URL" value can be get from [Identity and Access Management (IAM) endpoint list](https://developer.huaweicloud.com/en-us/endpoint).

```yaml
global:
  prefix: "huaweicloud"
  port: ":8087"
  metric_path: "/metrics"

auth:
  auth_url: "https://iam.xxx.yyy.com/v3"
  project_name: "{project_name}"
  access_key: "{access_key}"
  secret_key: "{secret_key}"
  region: "{region}"
```

or

```yaml
auth:
  auth_url: "https://iam.xxx.yyy.com/v3"
  project_name: "{project_name}"
  user_name: "{username}"
  password: "{password}"
  region: "{region}"
  domain_name: "{domain_name}"
```

## Prometheus Configuration

The huaweicloud exporter needs to be passed the address as a parameter, this can be done with relabelling.

Example config:

```yaml
global:
  scrape_interval: 1m # Set the scrape interval to every 1 minute seconds. Default is every 1 minute.
  scrape_timeout: 1m
scrape_configs:
  - job_name: 'huaweicloud'
    static_configs:
    - targets: ['10.0.0.10:8087']
    params:
      services: ['SYS.ELB,SYS.RDS,SYS.DCS,SYS.VPC,SYS.EVS,SYS.ECS,SYS.AS,SYS.NAT']
```
