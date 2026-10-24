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
    saveDir= "D:\\QA\\QA_AutoTest\\temp"     
    t=time.localtime()
    filename = time.strftime("%Y-%m-%d_%H-%M-%S", t) + '.jpg'
    img =ImageGrab.grab()
    img.save(os.path.join(saveDir, filename), 'jpeg')
    files = {'photo': open(os.path.join(saveDir, filename), 'rb')}
    url = "https://example.invalid/bot965536905:AAFSSpfGRYUzsm2548QHNroViU06Rz6YQNI/sendPhoto"
    data = {'chat_id' : "-332128122"}   
    r= requests.post(url, files=files, data=data)
    
def pokerpoint(username,i):
    card = poco(username).child('CardRoot').children()
    if len(card[i].attr('name')) == 8:
        point = card[i].attr('name')[7]
    else:
        point = card[i].attr('name')[7:9]
    if point == "1":
        point = 14
    return(point)

gamecount = 0
wincount = 0
losecount = 0
warwincount = 0
warlosecount = 0
wartie = 0
msg = "賭場戰爭測試開始" +"\n"
msg += "攜帶金幣:" + poco('Player1').child('PlayerInfo').child('MoneyText').attr('text')
callbot(msg)
startmoney = float(poco('Player1').child('PlayerInfo').child('MoneyText').attr('text'))
try:
    while(True):
        gamecount += 1
        poco('BetPoolButton1').click()
        usermoney = float(poco('Player1').child('PlayerInfo').child('MoneyText').attr('text'))
        usermoney += 2
        poco('Player1').child("CardRoot").children().wait(20).exists()
        sleep(1)
        playerpoint = int(pokerpoint('Player1',0))
        bankerpoint = int(pokerpoint('UIDealer',0))
        if playerpoint > bankerpoint:
            usermoney -= 0.05
            wincount += 1
        elif playerpoint < bankerpoint:
            usermoney -= 2
            losecount += 1
        else:
            usermoney += 9.5
            sleep(2)
            poco('WarButton').click()
            sleep(3)
            playerpoint = int(pokerpoint('Player1',4))
            bankerpoint = int(pokerpoint('UIDealer',4))
            if playerpoint > bankerpoint:
                usermoney += 0.95 
                warwincount += 1
            elif playerpoint < bankerpoint:
                usermoney -= 2
                warlosecount += 1
            else:
                usermoney += 1.9
                wartie += 1        
        poco('Chips').offspring('ChipPrefab').wait(30).exists()
        tt = gamecount % 50
        if tt == 0:
            Msg2 = "目前進行場次:" + str(gamecount) + "\n"
            Msg2 += "勝場數:" + str(wincount) + "\n"
            Msg2 += "敗場數:" + str(losecount) + "\n"
            Msg2 += "戰爭勝場數:" + str(warwincount) + "\n"
            Msg2 += "戰爭敗場數:" + str(warlosecount) + "\n"
            Msg2 += "戰爭平手數:" + str(wartie) + "\n"
            Msg2 += "輸贏金額:" + str(round((usermoney - startmoney),2))
            callbot(Msg2)
        if round(usermoney,2) != (float(poco('Player1').child('PlayerInfo').child('MoneyText').attr('text'))+1):
            msg1 = "金額計算有誤"+ "\n"
            msg1 += "計算金額:" + str(round(usermoney,2)) + "\n"
            msg1 += "畫面金額:" + poco('Player1').child('PlayerInfo').child('MoneyText').attr('text')  + "\n"
            msg1 += "莊家點數:" + str(bankerpoint)  + "\n"
            msg1 += "玩家點數:" + str(playerpoint)  + "\n"
            callbot(msg1)
            break
except:
    Msg1 = "異常結束,原因為: " +  str(sys.exc_info()[0]) + " ,結束畫面如下:"
    callbot(Msg1)
    #screenshot()    
else:
    callbot("賭場戰爭測試完畢") 