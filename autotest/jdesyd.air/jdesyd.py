# -*- encoding=utf8 -*-
__author__ = "Frank"

from airtest.core.api import *
from poco.drivers.unity3d import UnityPoco
from PIL import Image
from PIL import ImageGrab
import time
import os
import sys
import requests
import decimal
import random
poco = UnityPoco()
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

def checkpoint(player):
    ppoint = 0
    one = 0
    playercard = poco(player).offspring('CardPrefab(Clone)')
    lpc= len(playercard)
    for i in range(0,lpc):
        pc = playercard[i].attr('texture')
        if len(pc) > 15:
            ppoint += 10
        else:
            if int(pc[14]) == 1 :
                if ppoint < 11:
                    one += 1
                    ppoint += 11
                else:
                    ppoint += 1
            else:    
                ppoint += int(pc[14])
        if one > 0 and ppoint > 21:
            ppoint -= (one*10)
            one -= 1
    return(ppoint)

def playerpt():
    if poco('Player1').offspring('PatternImage').exists() == True:
        playerpoint = 23
    if len(poco("Player1").offspring('CardGroup1').offspring('Point').child('Text').attr('text')) == 4:
        playerpoint = int(poco("Player1").offspring('CardGroup1').offspring('Point').child('Text').attr('text')[2:4])
    elif len(poco("Player1").offspring('CardGroup1').offspring('Point').child('Text').attr('text')) == 5:
        playerpoint = int(poco("Player1").offspring('CardGroup1').offspring('Point').child('Text').attr('text')[3:5])
    else:
        playerpoint = int(poco("Player1").offspring('CardGroup1').offspring('Point').child('Text').attr('text'))
    return(playerpoint)

def dealerpt():
    if poco('UIDealer').offspring('PatternImage').exists() == True:
        dealerpoint = 22
    else:
        if len(poco('UIDealer').offspring('Text').attr('text')) == 4:
            dealerpoint = int(poco('UIDealer').offspring('Text').attr('text')[2:4])
        elif len(poco('UIDealer').offspring('Text').attr('text')) == 5:
            dealerpoint = int(poco('UIDealer').offspring('Text').attr('text')[3:5])
        else:    
            dealerpoint = int(poco('UIDealer').offspring('Text').attr('text'))
    return(dealerpoint)

joinusermoney = round(float(poco("MoneyText").attr("text")), 2)  
bjcount = 0
ttocount = 0
doublewincount = 0
dealwincount = 0
deal = 0
double = 0
bowshanflag = 0
wincount = 0
losecount = 0
duescount = 0
gamecount = 0
msg = "經典21點測試開始" +"\n"
msg += "攜帶金幣:" + poco('MoneyText').attr('text')
callbot(msg)
try: 
    while(True):
        gamecount += 1
        poco("Player1").offspring('CardGroup1').offspring('CardPrefab(Clone)').wait(60).exists()
        sleep(1)
        playerpoint = checkpoint("Player1")
        usermoney = round(float(poco("MoneyText").attr("text")), 2) + 1 
        dealerpoker = poco('UIDealer').offspring('CardPrefab(Clone)')[0].attr('texture')  
        if playerpoint == 21 :
            bjcount += 1           
            dealerpoker2 = poco('UIDealer').offspring('CardPrefab(Clone)')[1].attr('texture')
            while(dealerpoker2[12] == "4" ) :
                sleep(0.5)
                dealerpoker2 = poco('UIDealer').offspring('CardPrefab(Clone)')[1].attr('texture')
            sleep(1)
            if len(poco('UIDealer').offspring('CardPrefab(Clone)')) == 2 and poco('UIDealer').offspring('PatternImage').exists() == True:
                duescount += 1
            else:
                usermoney += 1.43
                wincount += 1                         
        else:
            
            #        bowshan = random.randint(1,100)
            #        if int(bowshan) <= 50:
            #            poco('InsuranceButton').click()
            #            usermoney -= 1
            #        else:
            #            poco('NoInsuranceButton').click()
            playerpoint = playerpt()      
            while(playerpoint < 17) :
                poco('HitButton').wait(20).exists()
                poco('HitButton').click()
                if poco('Player1').offspring('PatternImage').exists() == True:
                    playerpoint = 23
                else:
                    playerpoint = playerpt()
            if playerpoint < 22:
                if playerpoint == 21:
                    ttocount += 1
                else:    
                    poco('StandButton').wait(20).exists()
                    poco('StandButton').click()
                dealerpoker2 = poco('UIDealer').offspring('CardPrefab(Clone)')[1].attr('texture')
                while(dealerpoker2[12] == "4" ) :
                    sleep(1)
                    dealerpoker2 = poco('UIDealer').offspring('CardPrefab(Clone)')[1].attr('texture')
                if poco('UIDealer').offspring('PatternImage').exists() == False:                 
                    dealerpoint = checkpoint('UIDealer')
                    while(dealerpoint < 17):                     
                        if poco('UIDealer').offspring('PatternImage').exists() == True:
                            dealerpoint = 22 
                        else:
                            dealerpoint = checkpoint('UIDealer')
                    #qaplayerpoint = checkpoint("Player1")
                    #if playerpoint != qaplayerpoint:       
                    #    msg1 = "牌點計算有誤"+ "\n"
                    #    msg1 += "計算牌點:" + str(qaplayerpoint) +"\n"
                    #    msg1 += "實際牌點:" + str(playerpoint) +"\n"
                    #    callbot(msg1)
                    #    screenshot()
                    #    break
                    if playerpoint > dealerpoint and playerpoint < 22 or dealerpoint > 21:
                        if deal == 1 and double == 1 :
                            dealwincount += 1
                        elif deal == 0 and double == 1:
                            doublewincount += 1
                        else:
                            wincount += 1
                        usermoney += 0.95
                    elif playerpoint == dealerpoint and playerpoint < 22:
                        duescount += 1
                    else:
                        losecount += 1
                        usermoney -= 1
                else:
                    losecount += 1
                    usermoney -= 1
            else:
                losecount += 1
                usermoney -= 1
        poco('SquareCountdown').wait(60).exists()
        if round(usermoney,2) != float(poco('MoneyText').attr('text')):
            msg1 = "金額計算有誤"+ "\n"
            msg1 += "計算金額:" + str(round(usermoney,2)) + "\n"
            msg1 += "畫面金額:" + poco('MoneyText').attr('text')
            msg1 += "莊家點數" + str(dealerpoint)
            msg1 += "玩家點數" + str(playerpoint)
            screenshot()
            callbot(msg1)
            break
        tt = gamecount % 50
        if tt == 0:
            msg2 = "目前進行場次 :" + str(gamecount)+ "\n"
            msg2 += "勝場數(不加倍,不分牌):" + str(wincount) + "\n"
            msg2 += "勝場數(分牌):" + str(dealwincount) + "\n"
            msg2 += "勝場數(加倍):" + str(doublewincount) + "\n"
            msg2 += "玩家拿到BJ的場次:" + str(bjcount) + "\n"
            msg2 += "玩家拿到21點的場次" + str(ttocount) + "\n"
            msg2 += "敗場數:" + str(losecount) + "\n"
            msg2 += "平手數:" + str(duescount) + "\n"
            msg2 += "輸贏金額:" + str(round(usermoney-joinusermoney,2)) + "\n"
            callbot(msg2)
except:
    Msg1 = "經典21點異常結束,原因為: " +  str(sys.exc_info()[0]) + " ,結束畫面如下:"
    callbot(Msg1)
    screenshot()    
else:
    callbot("經典21點測試完畢")




   