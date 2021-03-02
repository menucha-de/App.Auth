# Util.Auth
## Manual upgrade
    NAME=auth
    ctr -n system i pull -k ghcr.io/peramic/$NAME:latest
    systemctl stop $NAME
    ctr -n system c rm $NAME
    ctr -n system c create --label NAME=$NAME --label IS_ACTIVE=true --env LOGHOST=$NAME --with-ns network:/var/run/netns/$NAME ghcr.io/peramic/$NAME:latest --mount type=bind,src=/etc/hosts,dst=/etc/hosts,options=rbind:ro --mount type=bind,src=/etc/resolv.conf,dst=/etc/resolv.conf,options=rbind:ro $NAME
    systemctl start $NAME
