#############################################
#! /bin/bash
#############################################

array=( dfdc wlzb xgpd jjbx xmbx dhzw wlzbjp csdjp wlhc )
for i in "${array[@]}"
do
    tmux send-keys -t gs-$i.0 "./$i --proto $1" ENTER
    sleep 1s
    tmux send-keys -t gs-$i.0 "log level 2" ENTER
    sleep 1s
    tmux send-keys -t gs-$i.0 "log lite 1" ENTER
    sleep 1s
    tmux send-keys -t gs-$i.0 "log on" ENTER
    sleep 1s
    tmux send-keys -t gs-$i.0 "version" ENTER
    echo "$i Start ..."
    sleep 1s
done

# ### DFDC
# tmux send-keys -t gs-dfdc.0 "./dfdc --proto $1" ENTER
# sleep 1s
# tmux send-keys -t gs-dfdc.0 "log on" ENTER
# sleep 1s
# tmux send-keys -t gs-dfdc.0 "log level 2" ENTER
# sleep 1s
# tmux send-keys -t gs-dfdc.0 "log lite 1" ENTER
# sleep 1s
# tmux send-keys -t gs-dfdc.0 "version" ENTER
# echo "DFDC Start ..."

# ### WLZB
# tmux send-keys -t gs-wlzb.0 "./wlzb --proto $1" ENTER
# tmux send-keys -t gs-wlzb.0 "log on" ENTER
# tmux send-keys -t gs-wlzb.0 "log level 2" ENTER
# tmux send-keys -t gs-wlzb.0 "log lite 1" ENTER
# tmux send-keys -t gs-wlzb.0 "version" ENTER

# ### XGPD
# tmux send-keys -t gs-xgpd.0 "./xgpd --proto $1" ENTER
# tmux send-keys -t gs-xgpd.0 "log on" ENTER
# tmux send-keys -t gs-xgpd.0 "log level 2" ENTER
# tmux send-keys -t gs-xgpd.0 "log lite 1" ENTER
# tmux send-keys -t gs-xgpd.0 "version" ENTER

# ### JJBX
# tmux send-keys -t gs-jjbx.0 "./jjbx --proto $1" ENTER
# tmux send-keys -t gs-jjbx.0 "log on" ENTER
# tmux send-keys -t gs-jjbx.0 "log level 2" ENTER
# tmux send-keys -t gs-jjbx.0 "log lite 1" ENTER
# tmux send-keys -t gs-jjbx.0 "version" ENTER

# ### XMBX
# tmux send-keys -t gs-xmbx.0 "./xmbx --proto $1" ENTER
# tmux send-keys -t gs-xmbx.0 "log on" ENTER
# tmux send-keys -t gs-xmbx.0 "log level 2" ENTER
# tmux send-keys -t gs-xmbx.0 "log lite 1" ENTER
# tmux send-keys -t gs-xmbx.0 "version" ENTER

# ### DHZW
# tmux send-keys -t gs-dhzw.0 "./dhzw --proto $1" ENTER
# tmux send-keys -t gs-dhzw.0 "log on" ENTER
# tmux send-keys -t gs-dhzw.0 "log level 2" ENTER
# tmux send-keys -t gs-dhzw.0 "log lite 1" ENTER
# tmux send-keys -t gs-dhzw.0 "version" ENTER

# ### WLZBJP
# tmux send-keys -t gs-wlzbjp.0 "./wlzbjp --proto $1" ENTER
# tmux send-keys -t gs-wlzbjp.0 "log on" ENTER
# tmux send-keys -t gs-wlzbjp.0 "log level 2" ENTER
# tmux send-keys -t gs-wlzbjp.0 "log lite 1" ENTER
# tmux send-keys -t gs-wlzbjp.0 "version" ENTER

# ### CSDJP
# tmux send-keys -t gs-csdjp.0 "./csdjp --proto $1" ENTER
# tmux send-keys -t gs-csdjp.0 "log on" ENTER
# tmux send-keys -t gs-csdjp.0 "log level 2" ENTER
# tmux send-keys -t gs-csdjp.0 "log lite 1" ENTER
# tmux send-keys -t gs-csdjp.0 "version" ENTER

# ### WLHC
# tmux send-keys -t gs-wlhc.0 "./wlhc --proto $1" ENTER
# tmux send-keys -t gs-wlhc.0 "log on" ENTER
# tmux send-keys -t gs-wlhc.0 "log level 2" ENTER
# tmux send-keys -t gs-wlhc.0 "log lite 1" ENTER
# tmux send-keys -t gs-wlhc.0 "version" ENTER
