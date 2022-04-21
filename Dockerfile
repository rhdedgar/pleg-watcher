# Copyright 2019 Doug Edgar.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# begin build container definition
FROM registry.access.redhat.com/ubi8/ubi-minimal as build

# Install prerequisites 
RUN microdnf install -y golang \
                   gcc \
                   git \
                   systemd-libs \
                   systemd-devel

ENV GOBIN=/bin \
    GOPATH=/go

# install pleg-watcher
RUN /usr/bin/go install github.com/rhdedgar/pleg-watcher@master


# begin run container definition
FROM registry.access.redhat.com/ubi8/ubi-minimal as run

ADD scripts/ /usr/local/bin/

COPY --from=build /bin/pleg-watcher /usr/local/bin

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
          /etc/sysconfig/docker

USER 0

CMD /usr/local/bin/start.sh
