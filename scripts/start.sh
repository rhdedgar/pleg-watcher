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

# Wait for the clam socket to become available before launching.
if [ ! -S /clam/clamd.sock ]; then
  until [ -S /clam/clamd.sock ]
  do
    echo "Failed to find clam socket, trying again in 30 seconds."
    sleep 30
  done
fi

echo This container hosts the following applications:
echo
echo '/usr/bin/pleg-watcher'
echo
echo 'Always listen for PLEG events from sdjournal.'
echo '----------------'
/usr/bin/pleg-watcher
