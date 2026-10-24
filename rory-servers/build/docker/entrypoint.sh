#! /bin/sh

tmux new -d -s ${APP_NAME} -c "/app" "./${APP_NAME}"
echo "Started ${APP_NAME} tmux session"

# loop to monitor for the taisun session
while true
do
  if [[ `tmux ls | grep ${APP_NAME} | awk '{print $1}' | cut -c 1-${#APP_NAME}` != "${APP_NAME}" ]]; then
    tmux new -d -s ${APP_NAME} -c "/app" "./${APP_NAME}"
    echo "Retarted ${APP_NAME} tmux session"
  fi
  sleep 5
done