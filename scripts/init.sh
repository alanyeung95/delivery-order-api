    #!/bin/sh
     
    set -eu
     
    echo "Checking DB connection ..."
     
    i=0
    until [ $i -ge 1000 ]
    do
      nc -z mysql.network 3306 && break
     
      i=$(( i + 1 ))
     
      echo "$i: Waiting for DB 5 second ..."
      sleep 5
    done
     
    if [ $i -eq 1000 ]
    then
      echo "DB connection refused, terminating ..."
      exit 1
    fi
     
    echo "DB is up ..."
