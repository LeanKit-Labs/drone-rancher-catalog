FROM rancher/docker:v1.10.2

ADD drone-rancher-catalog /go/bin/
VOLUME /var/lib/docker
ENTRYPOINT ["/usr/bin/dockerlaunch", "/go/bin/drone-rancher-catalog"]