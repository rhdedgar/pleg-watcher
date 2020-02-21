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

n=0
until [ $n -ge 5 ]
do
   mount -o bind /clam/clamd.sock /host/tmp/clamd.sock

   e=$?
   if [[ $e == 0 ]] ; then
     break
   elif [ n == 5 ] ; then
     echo "Failed to mount clam socket, scans will be unavailable."
   fi

   n=$[$n+1]
   t=$[$n*30]
   sleep t
done

echo This container hosts the following applications:
echo
echo '/usr/bin/pleg-watcher'
echo
echo 'Always listen for PLEG events from sdjournal.'
echo '----------------'
/usr/bin/pleg-watcher
