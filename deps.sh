#!/bin/bash -e
                                                                                                                                        
tag=$1

if [ "x$tag" = "x" ]; then
  tag="master"
fi

for m in $(find ./ -name 'go.mod'); do
  d=$(dirname $m);
  pushd $d >/dev/null;
  echo $d;
  PKGS=$(grep github.com/micro/go-micro/v2 go.mod | tr -s '\n' ' ' | grep -v replace | sed 's|require||g' | awk '{print $1}')
  for PKG in $PKGS; do
    /bin/bash -c "go get $PKG@$tag && go mod tidy"
  done
  popd >/dev/null
done

