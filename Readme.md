[![license](https://img.shields.io/github/license/RedisTimeSeries/netdata-redistimeseries-relay.svg)](https://github.com/RedisTimeSeries/netdata-redistimeseries-relay)

# Netdata metrics long term archiving with RedisTimeSeries
[![Forum](https://img.shields.io/badge/Forum-RedisTimeSeries-blue)](https://forum.redislabs.com/c/modules/redistimeseries)
[![Gitter](https://badges.gitter.im/RedisLabs/RedisTimeSeries.svg)](https://gitter.im/RedisLabs/RedisTimeSeries?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

Netdata supports  interoperability with other monitoring and visualization solutions for archiving the metrics, or providing long term dashboards, using Grafana or other tools, like it.  To this end, you can use Netdata Agent's exporting engine to send metrics to multiple external databases/services in parallel.

The following document explains the installation process of Netdata agent to collect thousands of metrics per server per second, and how to use netdata-redistimeseries-relay to store them into RedisTimeSeries. We will first dive into installing netdata locally - if you already have server(s) with netdata agent(s) running you can jump to section “Exporting quickstart”.

## Install netdata on your system

```
bash <(curl -Ss https://my-netdata.io/kickstart.sh)
```
Open your favorite browser and navigate to http://localhost:19999 to find Netdata’s real-time dashboard with hundreds of pre-configured charts and alarms.

## Exporting quickstart

This section covers the process of enabling an netdata exporting connector, using the JSON ( RedisTimeSeries ) connector. 

Open the exporting.conf file with edit-config.

```
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config exporting.conf
```

## Enable the exporting engine

Enable the exporting engine by setting enabled to yes:
```
[exporting:global]
    enabled = yes
```
Change how often the exporting engine sends metrics

By default, the exporting engine only sends metrics to external databases every 10 seconds to avoid congesting the destination with thousands of per-second metrics. Given that even one standalone RedisTimeSeries server can ingest hundreds of thousands of metrics per second we will change the frequency to every 1 second using the update every setting within the [exporting:global]section. In the end you should have an global section like the following:
```
[exporting:global]
    enabled = yes
    # send configured labels = yes
    # send automatic labels = no
    update every = 1 
```
## Enable the JSON (RedisTimeSeries) connector

We will now add a new RedisTimeSeries connector following the [<type>:<name>] format for defining connector instances. 


Given that Netdata does not support directly RedisTimeSeries we will add a json connector and use netdata-redistimeseries-relay to take JSON streams from Netdata client and write them to a RedisTimeseries DB ( this is the same approach as TimeScale with its netdata-timescale-relay ). You'll run this program in parallel with Netdata, and after a short configuration process, your metrics should start populating RedisTimeSeries, but first lets properly configure it. 


Add a new section to exporting.confnamed [json:redistimeseries]with the following configurations:

```
[json:redistimeseries]
    enabled = yes
    type = json
    destination = http://localhost:8080
    data source = as collected
    prefix = netdata
    hostname = my-hostname
    send charts matching = *
    send hosts matching = localhost *
    send names instead of ids = yes
```

## From Netdata JSON to RedisTimeSeries with netdata-redistimeseries-relay 

The easiest way to get and install netdata-redistimeseries-relay is to use go get and then issuing make:
```
# Fetch netdata-redistimeseries-relay and its dependencies
go get github.com/filipecosta90/netdata-redistimeseries-relay
cd $GOPATH/src/github.com/filipecosta90/netdata-redistimeseries-relay

# Install netdata-redistimeseries-relay binary:
make
```

## Running netdata-redistimeseries-relay

```
$ netdata-redistimeseries-relay --redistimeseries-host localhost:6379
Starting netdata-redistimeseries-relay...
Listening at 127.0.0.1:8080 for JSON inputs, and pushing RedisTimeSeries datapoints to localhost:6379...
```
