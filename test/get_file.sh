#!/bin/bash

for file in $(cat urls.txt);do
  if [ ! -f "${file##*/}" ];then
    echo "getting ${file##*/}"
    wget $file
  else
     echo "${file##*/} exists"
  fi
done
echo done

