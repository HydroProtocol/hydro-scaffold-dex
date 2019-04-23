#!/bin/sh
if [ $CONFIG_VARS ]; then
  SPLIT=$(echo $CONFIG_VARS | tr "," "\n")
  echo "{" >> ${CONFIG_FILE_PATH}/tmp

  for VAR in ${SPLIT}; do
      VALUE=$(printenv ${VAR})
      echo "  \"${VAR}\": \"${VALUE}\"," >> ${CONFIG_FILE_PATH}/tmp
  done

  echo $(sed '$ s/.$//' ${CONFIG_FILE_PATH}/tmp) >> ${CONFIG_FILE_PATH}/env.json
  rm ${CONFIG_FILE_PATH}/tmp
  echo "}" >> ${CONFIG_FILE_PATH}/env.json
fi

exec "$@"
