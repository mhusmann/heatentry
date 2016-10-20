#!/bin/bash

# heatentry als Shell Programm, ich will das nicht mehr als GUI

# Mo 17. Okt 14:19:17 CEST 2016
# ich will heatentry nur dann aufrufen, wenn ich es heute noch
# nicht gemacht habe. Sonst passiert gar nichts

today=$(date +'%Y-%m-%d')
# echo $today

st=$(stat --printf='%y' /home/mhusmann/Documents/src/pyt/heizung/heatpump.db)
# ich will nur die ersten 10 Zeichen von st
st=${st:0:10}
# echo $st

if [[ "$st" != "$today" ]]; then
      # echo "starte heatentry"
        DISPLAY=:0 LANG=de_DE.UTF-8 /usr/bin/sakura -x /home/mhusmann/usr/src/gocode/bin/heatentry
        # DISPLAY=:0 LANG=de_DE.UTF-8 /home/mhusmann/usr/src/gocode/bin/heatentry
fi
