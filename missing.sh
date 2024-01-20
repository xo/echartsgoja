#/bin/bash

SRC=$(realpath $(cd -P "$(dirname "${BASH_SOURCE[0]}")" && pwd))

pushd $SRC/testdata &> /dev/null

for i in $(find -type f -iname \*.js|sort -h); do
  if [ ! -f $i.svg ]; then
    echo "TestRender/$(sed -e 's/\.js$//' <<< "$i"|sed -e 's/^\.\///')"
  fi
done

popd &> /dev/null
