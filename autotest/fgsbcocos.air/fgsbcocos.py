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
    saveDir= "D:\\QA"     
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
    playercard = poco(Player).child("CardRoot").children()
    for i in range(0,3):
        pc = playercard[i].attr('name')
        if len(pc) > 8:
            playerpoint = pc[7:9]
        else:
            playerpoint = "0" + pc[7]
        Plist.append(playerpoint)
    return(Plist)

def cardflower(Player):
    Flist = []
    pf = 0
    playerflower = 0
    playercard = poco(Player).child("CardRoot").children()
    for i in range(0,3):
        pf = playercard[i].attr('name')
        playerflower = pf[5]
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
    if s1 == s2 and s2 == s3:
        same = 1
    if (p1 == (p2-1) and p2 == (p3-1)):
        straight = 1
    if p1 == p2 or p2 == p3:
        pair = 1
        if p1 == p2 and p2 == p3:
            three = 1
    if same == 1 and straight == 1 :
        royalstraight = 1
         
    if royalstraight == 1:
        pokertype = 6
    elif three == 1:
        pokertype = 5
    elif straight == 1:
        pokertype = 4
    elif same == 1:
        pokertype = 3
    elif pair == 1:
        pokertype = 2
    elif p3 > 11:
        pokertype = 1
    else:
        pokertype = 0
    return(pokertype)

straightwin = 0
straightnon = 0
samedownwin = 0
samedownnon = 0
samedownlose = 0
playerwin = 0
dues = 0
dealerwin = 0
gamecount = 0
tt = 0
msg = "富貴三寶測試開始" +"\n"
msg += "攜帶金幣:" + poco('UserMoney').child('Text').attr('text')
callbot(msg)
firstusermoney = float(poco('UserMoney').child('Text').attr('text'))
try:
    while(True):
        pwin = 0      
        dwin = 0
        gamecount += 1
        usermoney = float(poco('UserMoney').child('Text').attr('text'))
        poco('UIChip').child('Chips').offspring('ChipPrefab(1)').wait(60).exists()
        poco("Background2").child("bet1")[1].click()
        poco("Background2").child("bet2")[1].click()
        poco("ActionButton1").wait(30).exists()
        poco("ActionButton1").click()
        pcard = cardpoint("User1")
        pcard.sort()
        pcardf = cardflower("User1")
        ptype = int(pokercheck(pcardf,pcard))
        if ptype == 0:
            ptype = 1
        sleep(2)
        dcard = cardpoint("UIDealer")
        dcard.sort()
        dcardf = cardflower("UIDealer")
        dtype = int(pokercheck(dcardf,dcard))    
        if ptype > dtype :
            playerwin += 1
            if ptype == 4 :
                if dtype == 0:
                    straightnon += 1
                else:
                    straightwin += 1
            if ptype < 4 :
                if dtype == 0:
                    samedownnon += 1
                else:
                    samedownwin += 1
            if dtype == 0:
                usermoney += 0.95
            else:
                usermoney += 1.9
        elif ptype < dtype:
            dealerwin += 1
            usermoney -= 2
            if ptype < 4 :
                samedownlose += 1
            
        elif ptype == dtype:
            if ptype == 1 or ptype == 3:
                if int(pcard[2]) == int(dcard[2]):
                    pcard[2] = pcard[1]
                    dcard[2] = dcard[1]
                    if int(pcard[2]) == int(dcard[2]):
                        pcard[2] = pcard[0]
                        dcard[2] = dcard[0]
            if ptype == 2:
                for i  in range(2, 15):
                    if i < 10:
                        s =  '0' + str(i)
                    else:
                        s = str(i)
                    samecount = pcard.count(s)
                    samecount2 = dcard.count(s)
                    if samecount == 2:                      
                        ppair = i
                    if samecount2 == 2:
                        dpair = i
                if ppair == dpair:
                    psingle = pcard
                    psingle.remove(str(ppair))
                    psingle.remove(str(ppair))
                    dsingle = dcard
                    dsingle.remove(str(dpair))
                    dsingle.remove(str(dpair))
                    if psingle > dsingle:                  
                        playerwin += 1
                        samedownwin += 1
                        usermoney += 1.9
                    elif psingle < dsingle:
                        dealerwin += 1
                        samedownlose += 1
                        usermoney -= 2
                    else :
                        dues += 1
                        
                else:
                    pcard[2] = ppair
                    dcard[2] = dpair
                    if int(pcard[2]) > int(dcard[2]):
                        playerwin += 1
                        usermoney += 1.9              
                    elif int(pcard[2]) < int(dcard[2]):
                        dealerwin += 1
                        usermoney -= 2
                    else:               
                        dues += 1
                        
            elif int(pcard[2]) > int(dcard[2]):
                playerwin += 1
                usermoney += 1.9
                if ptype == 4 :
                    straightwin += 1
                if ptype < 4 :
                    samedownwin += 1
            elif int(pcard[2]) < int(dcard[2]):
                dealerwin += 1
                usermoney -= 2
                if ptype < 4 :
                    samedownlose += 1
            else: 
                dues += 1                
        if ptype == 2:
            usermoney += 0.95
        elif ptype == 3:
            usermoney += 3.8
        elif ptype == 4:
            usermoney += 6.65 #5.7+0.95
        elif ptype == 5:
            usermoney += 32.3 #28.5+38
        elif ptype == 6:
            usermoney += 42.75 #38+4.75
        else:
            usermoney -= 1    
        
        tt = gamecount % 100
        if tt == 0:
            Msg2 = "目前進行場次 :" + str(gamecount)+ "\n"
            Msg2 += "玩家順子贏"+ str(straightwin)+ "\n"
            Msg2 += "玩家順子,但莊家不成局:" + str(straightnon)+ "\n"
            Msg2 += "玩家同花或以下的牌勝利:" + str(samedownwin)+ "\n"
            Msg2 += "玩家同花或以下的牌,但莊家不成局:" + str(samedownnon)+ "\n"
            Msg2 += "玩家同花或以下的牌輸:" + str(samedownlose)+ "\n"
            Msg2 += "玩家勝:" + str(playerwin)+ "\n"
            Msg2 += "平手:" + str(dues) + "\n"
            Msg2 += "莊家勝:" + str(dealerwin)+ "\n" 
            Msg2 += "玩家輸贏金額"+ str(round(usermoney-firstusermoney,2))+ "\n"
            callbot(Msg2)
            
        poco('UIChip').child('Chips').offspring('ChipPrefab(1)').wait(60).exists()
        if round(usermoney,2) != round(float(poco('UserMoney').child('Text').attr('text')),2):
            msg1 = "金額計算有誤"+ "\n"
            msg1 += "計算金額:" + str(round(usermoney,2)) + "\n"
            msg1 += "畫面金額:" + poco('UserMoney').child('Text').attr('text')
            callbot(msg1)
            poco("ExitGameButton").click()
            break

        

except:
    Msg1 = "富貴三寶異常結束,原因為: " +  str(sys.exc_info()[0]) + " ,結束畫面如下:"
    callbot(Msg1)
    #screenshot()    
else:
    callbot("富貴三寶測試完畢")