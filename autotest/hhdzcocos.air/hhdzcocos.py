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
    Plist = []
    pc = 0
    playerpoint = 0
    playercard = poco(Player).children()
    for i in range(0,3):
        pc = playercard[i].attr('name')
        if len(pc) == 9:
            playerpoint = pc[5:7]
        else:
            playerpoint = "0" + pc[5]
        if playerpoint == "01":
            playerpoint = 14
        Plist.append(playerpoint)
    return(Plist)

def cardflower(Player):
    Flist = []
    pf = 0
    playerflower = 0
    playercard = poco(Player).children()
    for i in range(0,3):
        pf = playercard[i].attr('name')
        if len(pf) == 9:
            playerflower = pf[8]
        else:
            playerflower = pf[7]
        Flist.append(playerflower)
    return(Flist)
def pokercheck(pokerflowerlist,pokerlist):
    p1 = int(pokerlist[0])
    p2 = int(pokerlist[1])
    p3 = int(pokerlist[2])
    s1 = int(pokerflowerlist[0])
    s2 = int(pokerflowerlist[1])
    s3 = int(pokerflowerlist[2])
    pokertype = 0
    pair = 0
    same = 0
    straight = 0
    three = 0
    royalstraight = 0
    pairbig = 0
    if s1 == s2 and s2 == s3:
        same = 1
    if p3 == 14:
        if p1 == (p2-1) and p2 == (p3-1):
            straight = 1
        if p1 == (p2-1) and p2 == (p3-11):
            straight = 1
    else:
        if (p1 == (p2-1) and p2 == (p3-1)):
            straight = 1
    if p1 == p2 or p2 == p3:
        pair = 1
        if p1 == p2:
            if p1 > 8:
                pairbig = 1
        elif p2 == p3:
            if p2 > 8:
                pairbig = 1
        if p1 == p2 and p2 == p3:
            three = 1
    if same == 1 and straight == 1 :
        royalstraight = 1        
    if three == 1:
        pokertype = 6
    elif royalstraight == 1:
        pokertype = 5
    elif same == 1:
        pokertype = 4
    elif straight == 1:
        pokertype = 3
    elif pair == 1 and pairbig == 1:
        pokertype = 2
    elif pair == 1:
        pokertype = 1
    else:
        pokertype = 0
    return(pokertype)

msg = "cocos紅黑大戰測試開始" +"\n"
callbot(msg)   
gamecount = 0
redwin = 0
blackwin = 0
lucky = 0
try:
    while(True):
        winmoney = 0
        end = 0
        luckypoint = 0
        gamecount += 1       
        usermoney = float(poco('coin').attr('text'))
        poco("ButtonBet1").click()
        poco("ButtonBet3").click()
        while(end < 1): #判斷秒數讀完
            x = len(poco('countdown').children())
            if x == 0 :
                end = 1
        sleep(6)
        blackcard = cardpoint("HeiCard")
        blackf = cardflower("HeiCard")
        blacktype = int(pokercheck(blackf,blackcard))
        redcard = cardpoint("HungCard")
        redf = cardflower("HungCard")
        redtype = int(pokercheck(redf,redcard)) 
        if redtype > blacktype:
            redwin += 1 
            luckypoint = redtype
            winmoney += 1
        elif redtype < blacktype:
            blackwin += 1
            luckypoint = blacktype
            winmoney -= 1
        else:
            luckypoint = redtype
            if redtype == 1 :
                if redcard[0] == redcard[1]:
                    redcard[2] = redcard[0]
                    redcard[0] = redcard[1]
                if blackcard[0] == blackcard[1]:
                    blackcard[2] = blackcard[0]
                    blackcard[0] = blackcard[1]                
            redbig = int(redcard[2])
            blackbig = int(blackcard[2])
            if redbig > blackbig :
                redwin += 1 
                winmoney += 1
            elif redbig < blackbig:
                blackwin += 1
                winmoney -= 1
            else:
                if redbig == blackbig:
                    redbig = int(redcard[1])
                    blackbig = int(blackcard[1])
                if redbig > blackbig :
                    redwin += 1 
                    winmoney += 1
                elif redbig < blackbig:
                    blackwin += 1
                    winmoney -= 1
                else:
                    if redbig == blackbig:
                        redbig = int(redcard[0])
                        blackbig = int(blackcard[0])
                    if redbig > blackbig :
                        redwin += 1 
                        winmoney += 1
                    elif redbig < blackbig:
                        blackwin += 1
                        winmoney -= 1
                    else:
                        if int(redf[2]) > int(blackf[2]):
                            redwin += 1 
                            winmoney += 1
                        else:
                            blackwin += 1
                            winmoney -= 1
        if luckypoint == 6:
            lucky += 1
            winmoney += 10
        elif luckypoint == 5:
            lucky += 1
            winmoney += 5
        elif luckypoint == 4:
            lucky += 1
            winmoney += 3
        elif luckypoint == 3:
            lucky += 1
            winmoney += 2
        elif luckypoint == 2:
            lucky += 1
            winmoney += 1
        else:
            winmoney -= 1
        if winmoney > 0:
            usermoney += (winmoney*0.95)
        else:
            usermoney += winmoney    
        sleep(2)
        if round(float(usermoney),2) != round(float(poco('coin').attr('text')),2):
            msg1 = "金額計算有誤"+ "\n"
            msg1 += "計算金額:" + str(round(usermoney,2)) + "\n"
            msg1 += "畫面金額:" + poco('coin').attr('text')           
            callbot(msg1)
            #screenshot()
            break   
        tt = gamecount % 200
        if tt == 0:
            Msg2 = "目前進行場次 :" + str(gamecount)+ "\n"
            Msg2 += "紅贏:" + str(redwin)+ "\n"
            Msg2 += "黑贏:" + str(blackwin)+ "\n"
            Msg2 += "幸運一擊:" + str(lucky)+ "\n"
            #screenshot()
            break                               
        sleep(6)
        
        
except:
    Msg1 = "cocos紅黑大戰異常結束,原因為: " +  str(sys.exc_info()[0]) + " ,結束畫面如下:"
    callbot(Msg1) 
else:
    callbot("cocos紅黑大戰測試完畢")