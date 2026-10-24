#############################################
#! /bin/bash
#############################################

array=( dfdc wlzb xgpd jjbx xmbx dhzw wlzbjp csdjp wlhc )
for i in "${array[@]}"
do
    tmux send-keys -t "gs-$i.0" "roomserv stop" ENTER
    echo "$i Stoped ..."
done

# ### DFDC
# tmux send-keys -t gs-dfdc.0 "roomserv stop" ENTER

# ### WLZB
# tmux send-keys -t gs-wlzb.0 "roomserv stop" ENTER

# ### XGPD
# tmux send-keys -t gs-xgpd.0 "roomserv stop" ENTER

# ### JJBX
# tmux send-keys -t gs-jjbx.0 "roomserv stop" ENTER

# ### XMBX
# tmux send-keys -t gs-xmbx.0 "roomserv stop" ENTER

# ### DHZW
# tmux send-keys -t gs-dhzw.0 "roomserv stop" ENTER

# ### WLZBJP
# tmux send-keys -t gs-wlzbjp.0 "roomserv stop" ENTER

# ### CSDJP
# tmux send-keys -t gs-csdjp.0 "roomserv stop" ENTER

# ### WLHC
# tmux send-keys -t gs-wlhc.0 "roomserv stop" ENTER
