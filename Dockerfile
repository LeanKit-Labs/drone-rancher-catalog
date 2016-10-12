FROM docker:1.11-dind

ADD drone-rancher-catalog /bin/
ENTRYPOINT ["/usr/local/bin/dockerd-entrypoint.sh", "/bin/drone-rancher-catalog"]