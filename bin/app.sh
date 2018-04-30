#!/bin/bash
# &> /dev/null
# app.sh

APP_DIR=app/bin
APP_NAME=oceansf

usage() {
    echo "usage : ${0} {start|stop|status|restart}"
}

log() {
    echo [`date '+%Y-%m-%d %H:%M:%S'`] $*
}

# 파라미터 파싱
while [ "$1" != "" ]; do
    case $1 in
        start|stop|status|restart|check_restart)
            COMMAND=$1
            ;;
        *)
        usage
        exit 1
    esac
    shift
done


psgrep() {
    /bin/ps -ef | /bin/grep "${APP_NAME}" | /bin/grep -v grep
}

start() {
    cd $(dirname "$0")
    /usr/bin/nohup  /${APP_DIR}/${APP_NAME} >/dev/null 2>&1 &
}

stop() {
    /bin/kill -9 `psgrep | /bin/awk '{print $2}'` &>/dev/null
}

restart() {
    /bin/kill -USR2 `psgrep | /bin/awk '{print $2}'` &>/dev/null
}

status() {
    retval=`psgrep | wc -l`
    if [ ${retval} -eq 1 ]; then
        log ${APP_NAME} - ok
    else
        log ${APP_NAME} - failure
    fi
}

# 명령어 실행
${COMMAND}