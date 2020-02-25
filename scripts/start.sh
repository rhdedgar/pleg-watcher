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

touch /host/tmp/clamd.sock

if [ ! -S /host/tmp/clamd.sock ]; then
  n=0
  until mount -o bind /clam/clamd.sock /host/tmp/clamd.sock
  do
    n=$($n+30)
    echo "Failed to mount clam socket, trying again in $n seconds."
    sleep $n
  done
fi

echo This container hosts the following applications:
echo
echo '/usr/bin/pleg-watcher'
echo
echo 'Always listen for PLEG events from sdjournal.'
echo '----------------'
/usr/bin/pleg-watcher
