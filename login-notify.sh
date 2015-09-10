#!/usr/bin/env bash
#
#Author: HieuHT (HieuHT@vnoss.org)
#
if [ -n "$SSH_CLIENT" ];then 
    LOGIN_NOTIFY_HOST="192.168.10.100:8080"
    REMOTE_IP=`echo $SSH_CLIENT|awk '{print $1}'`
    HOSTNAME=`hostname -s`
    LOGIN_NOTIFY_API="http://${LOGIN_NOTIFY_HOST}/notify/login?user=${USER}&remoteip=$REMOTE_IP&servername=$HOSTNAME($SERVERIP)"
    /usr/bin/curl "$LOGIN_NOTIFY_API" &> /dev/null
fi
