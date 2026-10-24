# -*- encoding=utf8 -*-
__author__ = "Frank"

from airtest.core.api import *
from poco.drivers.cocosjs import CocosJsPoco
poco = CocosJsPoco()
from PIL import Image
from PIL import ImageGrab
import time
import os
import sys
import requests
auto_setup(__file__)

def callbot(Msg):
    r=requests.get("https://example.invalid/bot965536905:AAFSSpfGRYUzsm2548QHNroViU06Rz6YQNI/sendMessage?chat_id=-332128122&text=%s"%Msg)
    return()

def screenshot():
    saveDir= "D:\\"     
    t=time.localtime()
    filename = time.strftime("%Y-%m-%d_%H-%M-%S", t) + '.jpg'
    img =ImageGrab.grab()
    img.save(os.path.join(saveDir, filename), 'jpeg')
    files = {'photo': open(os.path.join(saveDir, filename), 'rb')}
    url = "https://example.invalid/bot965536905:AAFSSpfGRYUzsm2548QHNroViU06Rz6YQNI/sendPhoto"
    data = {'chat_id' : "-332128122"}   
    r= requests.post(url, files=files, data=data)

def cardcheck(pokername):
    Plist = []
    pc = 0
    playerpoint = 0
    playercard = poco(pokername).child('card').children()
    g =len(playercard)
    for i in range(0,g):
        pc = playercard[i].attr('name')
        if len(pc) == 10:
            playerpoint = pc[4:6]
        else:
            playerpoint = "0" + pc[4]
        Plist.append(playerpoint)
    return(Plist)
    
def paircheck(cardlist):
    a = cardlist[0]
    b = cardlist[1]
    if a == b:
        pairpoint = 1
    else:
        pairpoint = 0
    return(pairpoint)

def duel(player,banker):
    x = len(player)
    y = len(banker)
    totalppoint = 0
    totalbpoint = 0
    for i in range(0 , x):
        ppoint = int(player[i])
        if ppoint > 10:
            ppoint = 10
        totalppoint += ppoint
    for i in range(0 , y):
        bpoint = int(banker[i])
        if bpoint > 10:
            bpoint = 10
        totalbpoint += bpoint
    totalppoint = totalppoint % 10
    totalbpoint = totalbpoint % 10
    if totalppoint > totalbpoint :
        reward = 0
    elif totalppoint < totalbpoint:
        reward = 1
    else:
        reward = 2
    return(reward)
       
    
msg = "cocos百家樂測試開始" +"\n"
callbot(msg) 
gamecount = 0
playerwin = 0
playerpair = 0
bankerpair = 0
bankerwin = 0
dues = 0
try:
    while(True):
        end = 0
        gamecount += 1
        usermoney = float(poco('coin').attr('text'))
        winmoney = 0
        poco("ButtonBet2").click()
        poco("ButtonBet6").click()
        poco("ButtonBet7").click()
        while(end < 1): #判斷秒數讀完
            x = len(poco('countdown').children())
            if x == 0 :
                end = 1
            sleep(1)
        sleep(6)
        pcard = cardcheck('user')
        bcard = cardcheck('banker')
        ppair = paircheck(pcard)
        bpair = paircheck(bcard)
        if ppair == 1:
            winmoney += 11
            playerpair += 1
        else:
            winmoney -= 1
        if bpair == 1:
            winmoney += 11
            bankerpair += 1
        else:
            winmoney -= 1
        reward = duel(pcard,bcard)
        if reward == 0 :
            playerwin = 0
            winmoney += 1
        elif reward == 1:
            bankerwin += 1
            winmoney -= 1
        else:
            dues += 1
            winmoney += 0
        if winmoney > 0:
            usermoney += (winmoney*0.95)
        else:
            usermoney += winmoney

        sleep(4)    
        if round(float(usermoney),2) != round(float(poco('coin').attr('text')),2):
            msg1 = "金額計算有誤"+ "\n"
            msg1 += "計算金額:" + str(round(usermoney,2)) + "\n"
            msg1 += "畫面金額:" + poco('coin').attr('text')           
            callbot(msg1)
            #screenshot()
            break
        tt = gamecount % 100
        if tt == 0:
            Msg2 = "目前進行場次 :" + str(gamecount)+ "\n"
            Msg2 += "閒贏:" + str(playerwin)+ "\n"
            Msg2 += "莊贏:" + str(bankerwin)+ "\n"
            Msg2 += "和:" + str(dues)+ "\n"
            Msg2 += "閒對:" + str(playerpair)+ "\n"
            Msg2 += "莊對:" + str(bankerpair)+ "\n"
            #screenshot()
            break  
        sleep(3)

except:
    Msg1 = "cocos百家樂異常結束,原因為: " +  str(sys.exc_info()[0]) + " ,結束畫面如下:"
    screenshot()
    callbot(Msg1) 
else:
    callbot("cocos百家樂測試完畢")
    
    
