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

def cardpoint(Player):
    playercard = poco(Player).child('value').children()
    pc = playercard[0].attr('name')
    if len(pc) == 9:
        playerpoint = pc[5:7]
    else:
        playerpoint = "0" + pc[5]
    return(playerpoint)    

msg = "cocos百家樂測試開始" +"\n"
callbot(msg)
gamecount = 0
try:
    while(True):
        end = 0
        winmoney = 0
        gamecount += 1
        usermoney = float(poco('coin').attr('text'))
        poco("ButtonBet1").click()
        poco("ButtonBet2").click()
        while(end < 1): #判斷秒數讀完
            x = len(poco('countdown').children())
            if x == 0 :
                end = 1
            sleep(1)
        sleep(1)
        longpoint = int(cardpoint('longCard'))
        sleep(1)
        hupoint = int(cardpoint('huCard'))
        if longpoint > hupoint :
            winmoney = 0
        elif longpoint < hupoint :
            winmoney -= 2
        elif longpoint == hupoint :
            winmoney += 8
        if winmoney > 0:
            usermoney += (winmoney*0.95)
        else:
            usermoney += winmoney
        if round(float(usermoney),2) != round(float(poco('coin').attr('text')),2):
            msg1 = "金額計算有誤"+ "\n"
            msg1 += "計算金額:" + str(round(usermoney,2)) + "\n"
            msg1 += "畫面金額:" + poco('coin').attr('text')           
            callbot(msg1)
            #screenshot()
        sleep(3)    
except:
    Msg1 = "cocos百家樂異常結束,原因為: " +  str(sys.exc_info()[0]) 
    callbot(Msg1) 
else:
    callbot("cocos百家樂測試完畢")    

