/*
Copyright 2019 Doug Edgar.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/


# /usr/local/bin/start.sh will start the service

FROM registry.access.redhat.com/ubi8/ubi-minimal

# Pause indefinitely if asked to do so.
ARG OO_PAUSE_ON_BUILD
RUN test "$OO_PAUSE_ON_BUILD" = "true" && while sleep 10; do true; done || :

ADD scripts/ /usr/local/bin/

RUN microdnf install -y golang \
                   gcc \
                   git \
                   systemd-libs \
                   systemd-devel && \
    microdnf clean all

ENV GOBIN=/bin \
    GOPATH=/go

# Creating mount points for crio and docker sockets and dependencies.
RUN mkdir -p /host/usr/bin \
             /logs \
             /var/log/journal \
             /var/run/crio \
             /usr/bin \
             /etc/sysconfig && \
    touch /var/run/docker.sock \
          /var/run/crio/crio.sock \
          /usr/bin/docker-current \ 
          /etc/sysconfig/docker && \
    /usr/bin/go get github.com/rhdedgar/pleg-watcher && \
    cd /go/src/github.com/rhdedgar/pleg-watcher && \
    /usr/bin/go install && \
    cd && \
    rm -rf /go

USER 0

# Start processes
CMD /usr/local/bin/start.sh
