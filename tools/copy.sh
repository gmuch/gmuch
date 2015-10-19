#!/bin/bash

dest=$1
shift

for i in $@; do
  p="${i}"
  dp="${dest}/${p}"
  if echo $i | grep -q ":"; then
    p="`echo ${i} | cut -d: -f1`"
    dp="${dest}/`echo $i | cut -d: -f2`"
  fi

  mkdir -p `dirname $dp`
  cp $p $dp
done
