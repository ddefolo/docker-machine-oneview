#!/bin/bash

# build docker environment options based on current env
if [ -z "$1" ];then
  envfile=$(mktemp)
else
  envfile="$1"
fi

env_opts="--env-file $envfile "
for i in $(env | awk -F'=' '{print $1}'); do
  [ "$i" = "_" ] && continue
  # check for newlines greater than 1 and pass those with -e options
  if [ "$(eval 'echo -n "$'$i'"' | grep -c '$')" -gt 1 ]; then
      env_opts=$env_opts' -e '$i'='"'"$(eval 'echo "$'$i'"')"'";
#  elif echo $(eval 'echo -n $'$i) | grep -e '\s' > /dev/null 2<&1; then
#      eval 'echo '$i'=\"$'$i'\"' >> $envfile;
  else
      eval 'echo "'$i'="$'$i >> $envfile;
  fi
done

echo "$env_opts"
