#!/bin/sh

cmdArgs="$*"
if [ -n "$cmdArgs" ]; then
  /opt/aliddns $cmdArgs
  exit 0
fi

AccessKeyID=${AccessKeyID:-LTAIlzrfC9p85z8a}
AccessKeySecret=${AccessKeySecret:-wYDg4epPF4dGhCvjoIREehFYRAR0ll}
Domain=${Domain:-www.wpeak.win}
Redo=${Redo:-0}

cat > /opt/supervisord.conf <<EOF
[supervisord]
nodaemon=true

[program:aliddns]
command=/opt/aliddns --id ${AccessKeyID} --secret ${AccessKeySecret} auto-update --domain ${Domain} --redo ${Redo}
autorestart=true

EOF

/usr/bin/supervisord -c /opt/supervisord.conf
