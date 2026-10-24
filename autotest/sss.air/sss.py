# -*- encoding=utf8 -*-
from airtest.core.api import *


from poco.drivers.unity3d import UnityPoco
poco = UnityPoco()
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

#牌點在前面,花色在後面
#花色:方塊:0 梅花:1 紅心:2 黑桃:3
def pokercheck(node,conut):
    Flist = []
    Plist = []
    pc = 0
    playerpoint = 0
    playerflower = 0
    playercard1 = poco(node).children()
    for i in range(0,conut):
        pc = playercard1[i].attr('texture')
        if len(pc) > 6:
            playerpoint = pc[5:7]
        else:
            playerpoint = pc[5]
        playerflower = int(playerpoint) // 13
        if int(playerpoint) % 13 == 0:
            playerflower += 1 
        playerpoint = int(playerpoint) % 13
        if playerpoint == 0:
            playerpoint = 13
        Flist.append(playerflower)
        Plist.append(playerpoint)
    pclist = Plist + Flist
    return(pclist)


def firstcheck(pclist):
    plist = pclist[0:3]
    cardtype = "散牌"
    for i  in range(1, 14):
        samecount = plist.count(i)
        if samecount == 3:
            cardtype = "冲三"
        elif samecount == 2:
            cardtype = "一对"
    return(cardtype)

def lastcheck(pclist):
    plist = pclist[0:5]
    flist = pclist[5:10]
    fure = 0
    three = 0
    pair = 0
    same = 0
    samecount = 0
    sameflower = 0
    samestraight = 0
    samebigstraight = 0
    straight = 0
    single = 0
    for j  in range(0 , 4):
        sameflower = flist.count(j)
        if  sameflower == 5 :
            same = 1
    p1 =int(pclist[0])
    p2 =int(pclist[1])-1
    p3 =int(pclist[2])-2
    p4 =int(pclist[3])-3
    p5 =int(pclist[4])-4
    if p1 == p2:
        if p2 == p3:
            if p3 == p4:
                if p4 == p5:
                    straight = 1 
    if p1 == 1 and p2 == 9:
        if p2 == p3:
            if p3 == p4:
                if p4 == p5:                   
                    straight = 1
    if p1 == 1 and p2 == 9 and straight == 1 and same == 1:
        samebigstraight = 1
    for i  in range(1, 14):
        samecount = plist.count(i)
        if samecount == 2:
            pair += 1
        elif samecount == 3:
            three += 1 
        elif samecount == 4:
            fure += 1
    #牌型判斷，0:散牌，1:對子，2:兩對，3:三條，4:順子，5:同花，6:葫蘆，7:炸彈，8:同花順            
    if same == 1 and straight == 1:
        pok = "同花顺"
    elif fure == 1:
        pok = "炸弹"
    elif three == 1 and pair == 1:
        pok = "葫芦"
    elif same == 1:
        pok = "同花"
    elif straight == 1:
        pok = "顺子"
    elif three == 1:
        pok = "三条"
    elif pair == 2:
        pok = "两对"
    elif pair == 1:
        pok = "一对"
    else:
        pok = "散牌"
    return(pok)

def checkall():
    if len(poco("PokerNode4").children()) == 0:
        p1card = pokercheck("PokerNode1",3)
        firsttype = firstcheck(p1card)
        if firsttype != poco('PokerName1').attr('text'):
            msg1 += "判斷牌型 :" + str(firsttype)+ "\n"
            msg1 += "顯示牌型 :" + poco('PokerName1').attr('text') + "\n"
            callbot(msg1)
            screenshot()
        p2card = pokercheck("PokerNode2",5)
        secondtype = lastcheck(p2card)
        if secondtype != poco('PokerName2').attr('text'):
            msg1 += "判斷牌型 :" + str(secondtype)+ "\n"
            msg1 += "顯示牌型 :" + poco('PokerName2').attr('text') + "\n"
            callbot(msg1)
            screenshot()
        p3card = pokercheck("PokerNode3",5)
        lasttype = lastcheck(p3card)
        if lasttype != poco('PokerName3').attr('text'):
            msg1 += "判斷牌型 :" + str(secondtype)+ "\n"
            msg1 += "顯示牌型 :" + poco('PokerName2').attr('text') + "\n"
            callbot(msg1)
            screenshot()
    else:
        sleep(1)

msg = "十三水測試開始" +"\n"
callbot(msg)        
try:    
    while(True):        
        poco("FinishButton").wait(60).exists()
        poco("FinishButton").click()
        poco("ReadyButton").wait(60).exists()
        poco("Player1").child("TipsButton").click()
        checkall()
        poco("Player2").child("TipsButton").click()
        poco("Player2").child("TipsButton").click()
        checkall()
        poco("Player3").child("TipsButton").click()
        poco("Player3").child("TipsButton").click()
        checkall()
        poco("Player4").child("TipsButton").click()
        poco("Player4").child("TipsButton").click()
        checkall()
        poco("ReadyButton").click()
        poco("ReadyButton").click()
except:
    Msg1 = "十三水異常結束,原因為: " +  str(sys.exc_info()[0]) + " ,結束畫面如下:"
    callbot(Msg1)
    screenshot()    
else:
    callbot("十三水測試完畢")