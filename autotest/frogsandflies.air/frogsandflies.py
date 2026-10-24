# -*- encoding=utf8 -*-
__author__ = "Frank"

from airtest.core.api import *

auto_setup(__file__)# -*- encoding=utf8 -*-
__author__ = "Frank"

from PIL import Image
from PIL import ImageGrab
import time
import os
import sys
import requests
from airtest.core.api import *
from poco.drivers.cocosjs import CocosJsPoco
poco = CocosJsPoco()

auto_setup(__file__)
#賠率表
symbol_s1 = [100, 250, 1000]
symbol_s2 = [75, 200, 750]
symbol_s3 = [50, 80, 250]
symbol_s4 = [10, 40, 125]
symbol_s5 = [5, 40, 100]
symbol_s6 = [5, 40, 75]
symbol_s7 = [5, 10, 50]
symbol_s8 = [5, 10, 50]
symbol_s9 = [5, 10, 50]
symbol_c1 = [0, 0, 40, 200, 5000]

def callbot(Msg):
    r=requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=%s"%Msg)
    #r=requests.get("https://example.invalid/bot871711435:AAGYzFP4tNBQ-5E2j45a8qFBiDZDbmIlY2M/sendMessage?chat_id=-216697951&text=%s"%Msg)   #Lisa
    return()
   
def screenshot():
    saveDir= "D:\\QA"     
    t=time.localtime()
    filename = time.strftime("%Y-%m-%d_%H-%M-%S", t) + '.jpg'
    img =ImageGrab.grab()
    img.save(os.path.join(saveDir, filename), 'jpeg')
    files = {'photo': open(os.path.join(saveDir, filename), 'rb')}
    url = "https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendPhoto"
    data = {'chat_id' : "-270103520"}
    #url = "https://example.invalid/bot871711435:AAGYzFP4tNBQ-5E2j45a8qFBiDZDbmIlY2M/sendPhoto"  #Lisa
    #data = {'chat_id' : "-216697951"}  #Lisa    
    r= requests.post(url, files=files, data=data)

#賠率表
def GetBaseReward(sname, num, BetAmount):
    if sname == "s1":
        return(symbol_s1[num])
    elif sname == "s2":
        return(symbol_s2[num])
    elif sname == "s3":
        return(symbol_s3[num])
    elif sname == "s4":
        return(symbol_s4[num])
    elif sname == "s5":
        return(symbol_s5[num])
    elif sname == "s6":
        return(symbol_s6[num])
    elif sname == "s7":
        return(symbol_s7[num])
    elif sname == "s8":
        return(symbol_s8[num])
    elif sname == "s9":
        return(symbol_s9[num])
    elif sname == "s10":
        return(symbol_s10[num])
    elif sname == "s11":
        return(symbol_s11[num])
    else:
        return(symbol_c1[num])
    
def SymbolResult(checkReel):  #抓取輪軸資料
    SList = []    
    g = poco('ReelMain').child('New Node').child(checkReel).children()
    for i in range(4, 8):
        symbol = str(g[i].attr('name')) 
        SList.append(symbol)  
    return(SList)
    
    
    
def ReelResult(BetAmount):   #抓取輪軸資料並計算連線金額
    RList1 = SymbolResult("Reel1")
    RList2 = SymbolResult("Reel2")
    RList3 = SymbolResult("Reel3")
    RList4 = SymbolResult("Reel4")
    RList5 = SymbolResult("Reel5")    
    TotalReward = CaculateReward(RList1, RList2, RList3, RList4, RList5, BetAmount) + scatter(RList1, RList2, RList3, RList4, RList5)
    gresult = poco('WinEffect').children()    
    g = gresult[0].attr('name')
    g = float(g.replace(',',''))
    if TotalReward != g:
        RewardDifferent(RList1, RList2, RList3, RList4, RList5, TotalReward, g)
    return(TotalReward)

def CaculateReward(RList1, RList2, RList3, RList4, RList5, BetAmount): #計算連線金額
    RewardMoney = 0
    for i in RList1:
        R2 = 0
        R3 = 0
        R4 = 0
        R5 = 0
        linknum = 0        
        if (RList2.count(i) != 0) or (RList2.count("w1") != 0):
            R2 = R2 + RList2.count(i) + RList2.count("w1")
            if (RList3.count(i) != 0) or (RList3.count("w1") != 0):
                linknum = 0
                R3 = R3 + RList3.count(i) + RList3.count("w1")
                if (RList4.count(i) != 0) or (RList4.count("w1") != 0):
                    linknum = 1
                    R4 = R4 + RList4.count(i) + RList4.count("w1")
                    if(RList5.count(i) != 0):
                        linknum = 2
                        R5 = R5 + RList5.count(i)        
        if R5 > 0:            
            RewardMoney = RewardMoney + GetBaseReward(i, linknum, BetAmount) * R2 * R3 * R4 * R5 
        elif R4 > 0:            
            RewardMoney = RewardMoney + GetBaseReward(i, linknum, BetAmount) * R2 * R3 * R4
        elif R3 > 0:            
            RewardMoney = RewardMoney + GetBaseReward(i, linknum, BetAmount) * R2 * R3   
    return(RewardMoney)

def scatter(RList1, RList2, RList3, RList4, RList5):  #scatter存在及算獎金
    RewardMoney = 0
    linknum = 0        
    for i in  range(0, 4):
        if RList1[i] == "c1":
            linknum += 1
        if RList2[i] == "c1":
            linknum += 1
        if RList3[i] == "c1":
            linknum += 1
        if RList4[i] == "c1":
            linknum += 1
        if RList5[i] == "c1":
            linknum += 1
    if linknum > 0 :
        linknum -= 1
    RewardMoney = GetBaseReward("c1", linknum, BetAmount)   
    return(RewardMoney)


def RewardDifferent(RList1, RList2, RList3, RList4, RList5, TotalReward, gresult):
    print("RewardDifferent")
    Msg1 = RList1[0] + ", " + RList2[0] + ", "  + RList3[0] + ", " + RList4[0] + ", " + RList5[0] + "\n"
    Msg1 += RList1[1] + ", " + RList2[1] + ", " + RList3[1] + ", " + RList4[1] + ", " + RList5[1] + "\n"
    Msg1 += RList1[2] + ", " + RList2[2] + ", " + RList3[2] + ", " + RList4[2] + ", " + RList5[2] + "\n"
    Msg1 += "QA: " + str(TotalReward) + "  User: " + str(gresult)
    sult = poco('WinEffect').children()
    Msg1 += sult[0].attr('name')
    screenshot()
    callbot(Msg1)

msg = "frogsandflies測試開始" +"\n"
usermoney = float(poco('CreditL').attr('text'))
callbot(msg) 
SpinTotalAmount = 0    
TotalBet = 0    
spincount = 0
BetAmount = 40 
totalbet = poco('TotalBetL').attr('text')
while totalbet != "40":
    poco("BetPlus").click()
    totalbet = poco('TotalBetL').attr('text')
try:
    while(True):       
        poco('SpinButton').click()
        spincount += 1                
        TotalBet += BetAmount  
        sleep(10)        
        Amount = ReelResult(BetAmount) 
        usermoney = usermoney + Amount - BetAmount
        SpinTotalAmount += Amount
        if round(usermoney,2) != round(float(poco('CreditL').attr('text')),2):
            msg1 = "金額計算有誤"+ "\n"
            msg1 += "計算金額:" + str(round(usermoney,2)) + "\n"
            msg1 += "畫面金額:" + poco('CreditL').attr('text')           
            callbot(msg1)
            screenshot()
            break
        tt = spincount % 100
        if (spincount > 0) and (tt == 0):
            Msg1 = "frogsandflies遊戲次數: " + str(spincount)
            callbot(Msg1)  
except:
    Msg1 = "frogsandflies異常結束,原因為: " +  str(sys.exc_info()[0]) + " ,結束畫面如下:"
    screenshot()
    callbot(Msg1) 
else:
    callbot("frogsandflies測試完畢")
