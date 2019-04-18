#!/bin/sh
if [ $CONFIG_VARS ]; then
  SPLIT=$(echo $CONFIG_VARS | tr "," "\n")
  echo "window._env = {" >> ${CONFIG_FILE_PATH}/config.js

  for VAR in ${SPLIT}; do
      VALUE=$(printenv ${VAR})
      echo "  ${VAR}: \"${VALUE}\"," >> ${CONFIG_FILE_PATH}/config.js
  done

  echo "}" >> ${CONFIG_FILE_PATH}/config.js
fi

exec "$@"
