---
title: "Command line"
linkTitle: "Command line"
weight: 2
description: >
  Static command-line configuration options.
---

The following command-line options are supported:

| Name                                                    | Type                       | Default                 | Since |
|---------------------------------------------------------|----------------------------|-------------------------|-------|
| [`--allow-cross-namespace`](#allow-cross-namespace)     | [true\|false]              | `false`                 |       |
| [`--annotation-prefix`](#annotation-prefix)             | prefix without `/`         | `ingress.kubernetes.io` | v0.8  |
| [`--default-backend-service`](#default-backend-service) | namespace/servicename      | (mandatory)             |       |
| [`--default-ssl-certificate`](#default-ssl-certificate) | namespace/secretname       | (mandatory)             |       |
| [`--ingress-class`](#ingress-class)                     | name                       | `haproxy`               |       |
| [`--kubeconfig`](#kubeconfig)                           | /path/to/kubeconfig        | in cluster config       |       |
| [`--max-old-config-files`](#max-old-config-files)       | num of files               | `0`                     |       |
| [`--publish-service`](#publish-service)                 | namespace/servicename      |                         |       |
| [`--rate-limit-update`](#rate-limit-update)             | uploads per second (float) | `0.5`                   |       |
| [`--reload-strategy`](#reload-strategy)                 | [native\|reusesocket]      | `reusesocket`           |       |
| [`--sort-backends`](#sort-backends)                     | [true\|false]              | `false`                 |       |
| [`--tcp-services-configmap`](#tcp-services-configmap)   | namespace/configmapname    | no tcp svc              |       |
| [`--verify-hostname`](#verify-hostname)                 | [true\|false]              | `true`                  |       |
| [`--wait-before-shutdown`](#wait-before-shutdown)       | seconds as integer         | `0`                     | v0.8  |
| [`--watch-namespace`](#watch-namespace)                 | namespace                  | all namespaces          |       |

---

## --allow-cross-namespace

`--allow-cross-namespace` argument, if added, will allow reading secrets from one namespace to an
ingress resource of another namespace. The default behavior is to deny such cross namespace reading.
This adds a breaking change from `v0.4` to `v0.5` on `ingress.kubernetes.io/auth-tls-secret`
annotation, where cross namespace reading were allowed without any configuration.

---

## --annotation-prefix

Changes the annotation prefix the controller should look for when parsing services and ingress
objects. The default value is `ingress.kubernetes.io` if not declared, which means SSL Redirect
should be configured with the annotation name `ingress.kubernetes.io/ssl-redirect`. Annotations
with other prefix are ignored. This allows using HAProxy Ingress with other ingress controllers
that shares ingress and service objects without conflicting each other.

---

## --default-backend-service

Defines the `namespace/servicename` that should be used if the incoming request doesn't match any
hostname, or the requested path doesn't match any location within the desired hostname.

This is a mandatory argument used in the [deployment](https://github.com/jcmoraisjr/haproxy-ingress/tree/master/examples/deployment) example page.

---

## --default-ssl-certificate

Defines the `namespace/secretname` of the default certificate that should be used if ingress
resources using TLS configuration doesn't provide it's own certificate.

This is a mandatory argument used in the [deployment](https://github.com/jcmoraisjr/haproxy-ingress/tree/master/examples/deployment) and
[TLS termination](https://github.com/jcmoraisjr/haproxy-ingress/tree/master/examples/tls-termination) example pages.

---

## --ingress-class

More than one ingress controller is supported per Kubernetes cluster. The `--ingress-class`
argument allow to override the class name of ingress resources that this instance of the
controller should listen to. Class names that match will be used in the HAProxy configuration.
Other classes will be ignored.

The ingress resource must use the `kubernetes.io/ingress.class` annotation to name it's
ingress class.

---

## --kubeconfig

Ingress controller will try to connect to the Kubernetes master using environment variables and a
service account. This behavior can be changed using `--kubeconfig` argument that reference a
kubeconfig file with master endpoint and credentials. This is a mandatory argument if the controller
is deployed outside of the Kubernetes cluster.

---

## --max-old-config-files

Everytime a configuration change need to update HAProxy, a configuration file is rewritten even if
dynamic update is used. By default the same file is recreated and the old configuration is lost.
Use `--max-old-config-files` to configure after how much files Ingress controller should start to
remove old configuration files. If `0`, the default value, a single `haproxy.cfg` is used.

---

## --publish-service

Some infrastructure tools like `external-DNS` relay in the ingress status to created access routes to the services exposed with ingress object.

```
apiVersion: extensions/v1beta1
kind: Ingress
...
status:
  loadBalancer:
    ingress:
    - hostname: <ingressControllerLoadbalancerFQDN>
```

Use `--publish-service=namespace/servicename` to indicate the services fronting the ingress controller. The controller mirrors the address of this service's endpoints to the load-balancer status of all Ingress objects it satisfies.

---

## --rate-limit-update

Use `--rate-limit-update` to change how much time to wait between HAProxy reloads. Note that the first
update is always immediate, the delay will only prevent two or more updates in the same time frame.
Moreover reloads will only occur if the cluster configuration has changed, otherwise no reload will
occur despite of the rate limit configuration.

This argument receives the allowed reloads per second. The default value is `0.5` which means no more
than one reload will occur within `2` seconds. The lower limit is `0.05` which means one reload within
`20` seconds. The highest one is `10` which will allow ingress controller to reload HAProxy up to 10
times per second.

---

## --reload-strategy

The `--reload-strategy` command-line argument is used to select which reload strategy
HAProxy should use. The following options are available:

* `native`: Uses native HAProxy reload option `-sf`.
* `reusesocket`: (starting on v0.6) Uses HAProxy `-x` command-line option to pass the listening sockets between old and new HAProxy process, allowing hitless reloads. This is the default option since v0.8.
* `multibinder`: (deprecated on v0.6) Uses GitHub's [multibinder](https://github.com/github/multibinder). This [link](https://githubengineering.com/glb-part-2-haproxy-zero-downtime-zero-delay-reloads-with-multibinder/)
describes how it works.

---

## --sort-backends

Ingress will randomly shuffle backends and server endpoints on each reload in order to avoid
requesting always the same backends just after reloads, depending on the balancing algorithm.
Use `--sort-backends` to avoid this behavior and always declare backends and upstream servers
in the same order.

---

## --tcp-services-configmap

Configure `--tcp-services-configmap` argument with `namespace/configmapname` resource with TCP
services and ports that HAProxy should listen to. Use the HAProxy's port number as the key of the
configmap.

The value of the configmap entry is a colon separated list of the following items:

1. `<namespace>/<service-name>`, mandatory, is the well known notation of the service that will receive incoming connections.
1. `<portnumber>`, mandatory, is the port number the upstream service is listening - this is not related to the listening port of HAProxy.
1. `<in-proxy>`, optional, should be defined as `PROXY` if HAProxy should expect requests using the [PROXY](http://www.haproxy.org/download/1.9/doc/proxy-protocol.txt) protocol. Leave empty to not use PROXY protocol. This is usually used only if there is another load balancer in front of HAProxy which supports the PROXY protocol. PROXY protocol v1 and v2 are supported.
1. `<out-proxy>`, optional, should be defined as `PROXY` or `PROXY-V2` if the upstream service expect connections using the PROXY protocol v2. Use `PROXY-V1` instead if the upstream service only support v1 protocol. Leave empty to connect without using the PROXY protocol.
1. `<namespace/secret-name>`, optional, used to configure SSL/TLS over the TCP connection. Secret should have `tls.crt` and `tls.key` pair used on TLS handshake. Leave empty to not use ssl-offload.

Optional fields should be skipped using two consecutive colons.

In the example below:

```
...
data:
  "5432": "default/pgsql:5432"
  "8000": "system-prod/http:8000::PROXY-V1"
  "9900": "system-prod/admin:9900:PROXY::system-prod/tcp-9900"
  "9990": "system-prod/admin:9999::PROXY-V2"
  "9999": "system-prod/admin:9999:PROXY:PROXY"
```

HAProxy will listen 5 new ports:

* `5432` will proxy to a `pgsql` service on `default` namespace.
* `8000` will proxy to `http` service, port `8000`, on the `system-prod` namespace. The upstream service will expect connections using the PROXY protocol but it only supports v1.
* `9900` will proxy to `admin` service, port `9900`, on the `system-prod` namespace. Clients should connect using the PROXY protocol v1 or v2. Upcoming connections should be encrypted, HAProxy will ssl-offload data using crt/key provided by `system-prod/tcp-9900` secret.
* `9990` and `9999` will proxy to the same `admin` service and `9999` port and the upstream service will expect connections using the PROXY protocol v2. The HAProxy frontend, however, will only expect PROXY protocol v1 or v2 on it's port `9999`.

---

## --verify-hostname

Ingress resources has `spec/tls[]/secretName` attribute to override the default X509 certificate.
As a default behavior the certificates are validated against the hostname in order to match the
SAN extension or CN (CN only up to `v0.4`). Invalid certificates, ie certificates which doesn't
match the hostname are discarded and a warning is logged into the ingress controller logging.

Use `--verify-hostname=false` argument to bypass this validation. If used, HAProxy will provide
the certificate declared in the `secretName` ignoring if the certificate is or is not valid.

---

## --wait-before-shutdown

If argument `--wait-before-shutdown` is defined, controller will wait defined time in seconds
before it starts shutting down components when SIGTERM was received. By default, it's 0, which means
the controller starts shutting down itself right after signal was sent.

---

## --watch-namespace

By default the proxy will be configured using all namespaces from the Kubernetes cluster. Use
`--watch-namespace` with the name of a namespace to watch and build the configuration of a
single namespace.
