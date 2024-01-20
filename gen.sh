#!/bin/bash

# github repos
REPO_E=apache/echarts
REPO_X=apache/echarts-examples

SRC=$(realpath $(cd -P "$(dirname "${BASH_SOURCE[0]}")" && pwd))

set -e

WORKDIR=$HOME/src/charts

mkdir -p $WORKDIR

git_latest_tag() {
  git -C "$1" describe --tags
}

git_checkout_reset() {
  local dir="$WORKDIR/$1" name="$1" repo="$2" branch="$3"
  if [ ! -d "$dir" ]; then
    (set -x;
      git clone "$repo" "$dir"
    )
  fi
  (set -x;
    git -C "$dir" fetch origin
  )
  (set -x;
    git -C "$dir" reset --hard
    git -C "$dir" clean -f -x -d -e node_modules
    git -C "$dir" checkout origin/$branch &> /dev/null
  )
  local ver=$(git_latest_tag "$dir")
  echo "$name $ver"
  if [ "$name" = "echarts" ]; then
    echo "$ver" > "$SRC/version.txt"
  fi
}

git_checkout_reset echarts          "https://github.com/${REPO_E}.git" "master"
git_checkout_reset echarts-examples "https://github.com/${REPO_X}.git" "gh-pages"

(set -x;
  cp $WORKDIR/echarts/dist/echarts.min.js $SRC
)

pushd $WORKDIR/echarts-examples &> /dev/null
(set -x;
  pnpm install
  pnpm add -D typescript
)
popd &> /dev/null

pushd $WORKDIR/echarts-examples/public/examples &> /dev/null
(set -x;
  rm -rf js out
  npx tsc -p tsconfig.json || :
  rsync \
    -rvP \
    --exclude gl \
    --exclude archive \
    --delete \
    ./js/. $SRC/testdata/.
)
popd &> /dev/null
