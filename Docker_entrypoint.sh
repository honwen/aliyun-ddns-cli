#!/bin/sh

cmdArgs="$*"
if [ -n "$cmdArgs" ]; then
  /opt/aliddns $cmdArgs
  exit 0
fi

AccessKeyID=${AccessKeyID:-1234567890}
AccessKeySecret=${AccessKeySecret:-abcdefghijklmn}
Domain=${Domain:-ddns.example.win}
Redo=${Redo:-0}

cat > /opt/supervisord.conf <<EOF
[supervisord]
nodaemon=true

[program:aliddns]
command=/opt/aliddns --id ${AccessKeyID} --secret ${AccessKeySecret} auto-update --domain ${Domain} --redo ${Redo}
autorestart=true
redirect_stderr=true
stdout_logfile=/dev/stdout
stdout_logfile_maxbytes=0

EOF

/usr/bin/supervisord -c /opt/supervisord.conf
