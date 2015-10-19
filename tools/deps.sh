#!/bin/sh

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <dynamic>"
  exit 1
fi

dep() {
  for i in `ldd $1 | sed -e 's:^\(.*\)* \(.*\) (.*)$:\2:g' -e 's:^\s*\(.*\) (.*)$:\1:g'`; do
    if ! echo $i | grep -q '/'; then
      continue
    fi

    echo $i
    $0 $i
  done
}

dep $1 | sort | uniq
