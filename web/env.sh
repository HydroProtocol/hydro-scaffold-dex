#!/bin/sh
if [ $CONFIG_VARS ]; then
  # clear
  echo -n > ${CONFIG_FILE_PATH}/config.js

  SPLIT=$(echo $CONFIG_VARS | tr "," "\n")
  echo "window._env = {" >> ${CONFIG_FILE_PATH}/config.js

  for VAR in ${SPLIT}; do
      VALUE=$(printenv ${VAR})
      echo "  ${VAR}: \"${VALUE}\"," >> ${CONFIG_FILE_PATH}/config.js
  done

  echo "}" >> ${CONFIG_FILE_PATH}/config.js
fi

# disable broswer cache
sed -i "s/config.js?v=[0-9]*/config.js?v=$(date +'%s')/g" /srv/http/index.html

# exec CMD
exec "$@"
