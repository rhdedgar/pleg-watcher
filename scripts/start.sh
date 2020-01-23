#!/bin/bash -e

# This is useful so we can debug containers running inside of OpenShift that are
# failing to start properly.

if [ "$OO_PAUSE_ON_START" = "true" ] ; then
  echo
  echo "This container's startup has been paused indefinitely because OO_PAUSE_ON_START has been set."
  echo
  while true; do
    sleep 10    
  done
fi


echo "Creating mount points for crio and docker sockets and dependencies."
touch /var/run/docker.sock
touch /var/run/crio/crio.sock
touch /usr/bin/docker-current
touch /etc/sysconfig/docker

echo "Mounting crio and docker dependencies."
mount --bind -o ro /host/var/run/docker.sock /var/run/docker.sock
mount --bind -o ro /host/var/run/crio/crio.sock /var/run/crio/crio.sock
mount --bind -o ro /host/usr/bin/docker-current /usr/bin/docker-current
mount --bind -o ro /host/etc/sysconfig/docker /etc/sysconfig/docker

echo This container hosts the following applications:
echo
echo '/usr/bin/pleg-watcher'
echo
echo 'Always listen for PLEG events from sdjournal.'
echo '----------------'
/usr/bin/pleg-watcher
