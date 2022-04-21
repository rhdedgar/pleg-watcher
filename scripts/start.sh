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

if [ "$SCHEDULED_SCAN" = "true" ] ; then
  echo "Evironment set to scheduled scan."
  
  # Stagger the start time of the scheduled scan by a random number of hours & minutes.
  # Times are calculated in seconds.
  hour_stagger="$((RANDOM % 24 * 60 * 60))"
  minute_stagger="$((RANDOM % 24 * 60))"

  while true; do
    current_time=$(date +%s)
    selection_time=$(date -d "this $SCHEDULED_SCAN_DAY" '+%s')

    sleep_seconds=$(( selection_time - current_time + hour_stagger + minute_stagger ))

    echo "Sleeping for $sleep_seconds seconds."
    sleep $sleep_seconds

    /usr/local/bin/pleg-watcher
    echo "Sleeping 1d before reassessing when the next scan will take place."
    sleep 1d
  done
else
  echo "Evironment set to active scan."
  echo This container hosts the following applications:
  echo
  echo '/usr/bin/pleg-watcher'
  echo
  echo 'Always listen for PLEG events from sdjournal.'
  echo '----------------'
  /usr/local/bin/pleg-watcher
fi

echo "Both scheduled and active scanning blocks have exited. This shouldn't happen."
echo "Staying active for troubleshooting."
while true; do
  sleep 10
done
