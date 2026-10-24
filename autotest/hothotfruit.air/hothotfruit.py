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
symbol_s1 = [0, 0, 10, 40, 100, 150, 300, 400, 500, 600, 750, 800, 1200, 4000, 25000]
symbol_s2 = [0, 0, 5, 30, 75, 175, 200, 250, 300, 900]
symbol_s3 = [0, 0, 5, 25, 60, 120, 150, 200, 250, 750]
symbol_s4 = [0, 0, 5, 25, 50, 100, 125, 175, 200, 600]
symbol_s5 = [0, 0, 5, 25, 40, 75, 100, 150, 175, 500]


def callbot(Msg):
    #print(Msg)
    r=requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=%s"%Msg)
    #r=requests.get("https://example.invalid/bot871711435:AAGYzFP4tNBQ-5E2j45a8qFBiDZDbmIlY2M/sendMessage?chat_id=-216697951&text=%s"%Msg)   #Lisa
    return()
   
def screenshot():
    saveDir= "D:\\"     
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


def GetBaseReward(sname, num, BetAmount):  #賠率表
    if sname == "s1":
        return(symbol_s1[num])
    elif sname == "s2":
        return(symbol_s2[num])
    elif sname == "s3":
        return(symbol_s3[num])
    elif sname == "s4":
        return(symbol_s4[num])
    else: 
        return(symbol_s5[num])
    
    
def SymbolResult(checkReel):  #抓取輪軸資料
    SList = []    
    g = poco('ReelMain').child('New Node').child(checkReel).children()
    for i in range(3, 6):
        symbol = str(g[i].attr('name')) 
        SList.append(symbol)
    return(SList)
    
def symbolchange(symbollist):  #抓取輪軸資料
    SList = [] 
    linknum = 0
    if symbollist == "w2":
        symbollist = "w1"
        linknum = 1
    elif symbollist == "ss1":
        symbollist = "s1"
        linknum = 2
    elif symbollist == "ss2":
        symbollist = "s2"
        linknum = 1
    elif symbollist == "ss3":
        symbollist = "s3"
        linknum = 1
    elif symbollist == "ss4":
        symbollist = "s4"
        linknum = 1
    elif symbollist == "ss5":
        symbollist = "s5"
        linknum = 1   
    SList.append(symbollist)
    SList.append(linknum)
    return(SList)    
    
def ReelResult(BetAmount):  #抓取輪軸資料並計算連線金額
    RList1 = SymbolResult("Reel1")
    RList2 = SymbolResult("Reel2")
    RList3 = SymbolResult("Reel3")
    RList4 = SymbolResult("Reel4")
    RList5 = SymbolResult("Reel5")
    TotalReward = CaculateReward(RList1, RList2, RList3, RList4, RList5, BetAmount)
    gresult = poco('WinEffect').children()    
    g = gresult[0].attr('name')
    g = float(g.replace(',',''))
    if TotalReward != g:
        RewardDifferent(RList1, RList2, RList3, RList4, RList5, TotalReward, g)
    return(TotalReward)

def CaculateReward(RList1, RList2, RList3, RList4, RList5, BetAmount):
    RewardMoney = 0
    K = 0
    line1 = 0
    line2 = 0
    line3 = 0
    line4 = 0
    line5 = 0
    for i in RList1:
        R1 = 0
        R2 = 0
        R3 = 0
        R4 = 0
        R5 = 0        
        linknum = 0
        symbolname = i
        K += 1       
        if K == 1:
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5,  symbolname, RewardMoney, 0, 0, 0, 0, 0) #線條1
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5,  symbolname, RewardMoney, 0, 0, 2, 0, 0) #線條4

            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5,  symbolname, RewardMoney, 0, 1, 2, 1, 2) #線條11
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5,  symbolname, RewardMoney, 0, 0, 1, 2, 2) #線條12
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5,  symbolname, RewardMoney, 0, 1, 0, 1, 0) #線條15
            
        if K == 2:   
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5,  symbolname, RewardMoney, 1, 1, 1, 1, 1) #線條2
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5,  symbolname, RewardMoney, 1, 2, 0, 2, 1) #線條6
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5,  symbolname, RewardMoney, 1, 0, 2, 0, 1) #線條7
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5,  symbolname, RewardMoney, 1, 0, 1, 2, 1) #線條8
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5,  symbolname, RewardMoney, 1, 2, 1, 0, 1) #線條9
            
        if K == 3:    
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5,  symbolname, RewardMoney, 2, 2, 2, 2, 2) #線條3
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5,  symbolname, RewardMoney, 2, 2, 0, 2, 2) #線條5
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5,  symbolname, RewardMoney, 2, 1, 0, 1, 0) #線條10
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5,  symbolname, RewardMoney, 2, 2, 1, 0, 0) #線條13
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5,  symbolname, RewardMoney, 2, 1, 2, 1, 2) #線條14
            
    return(RewardMoney)

def checksymbolline(RList1, RList2, RList3, RList4, RList5, symbolname, RewardMoney, R1list, R2list, R3list, R4list, R5list): #連線的規則  
    R5 = 0
    R4 = 0
    R3 = 0
    R2 = 0
    linknum = 0
    Reel1symbol = symbolchange(RList1[R1list])
    Reel2symbol = symbolchange(RList2[R2list])
    Reel3symbol = symbolchange(RList3[R3list])
    Reel4symbol = symbolchange(RList4[R4list])
    Reel5symbol = symbolchange(RList5[R5list]) 
    symbolname = Reel1symbol[0]
    if symbolname == "w1":
        symbolname = Reel2symbol[0]
        if symbolname == "w1":
            symbolname = Reel3symbol[0]
    #print(symbolname)
    if Reel2symbol[0] == str(symbolname)  or Reel2symbol.count("w1") :
        R2 = R2 + Reel2symbol.count(symbolname) + Reel2symbol.count("w1")
        linknum = 1
        if Reel3symbol[0] == str(symbolname):
            linknum = 2
            R3 = R3 + Reel3symbol.count(symbolname)
            if Reel4symbol[0] == str(symbolname) or (Reel4symbol.count("w1") != 0) :
                linknum = 3
                R4 = R4 + Reel4symbol.count(symbolname) + Reel4symbol.count("w1")
                if Reel5symbol[0] == str(symbolname) or (Reel5symbol.count("w1") != 0) :
                    linknum = 4
                    R5 = R5 + Reel5symbol.count(symbolname) + Reel5symbol.count("w1")
    
    
    if R5 > 0:   
        linknum = linknum + int(Reel1symbol[1]) + int(Reel2symbol[1]) + int(Reel3symbol[1]) + int(Reel4symbol[1]) + int(Reel5symbol[1])
        RewardMoney = RewardMoney + GetBaseReward(symbolname, linknum, BetAmount) * R2 * R3 * R4 * R5 
    elif R4 > 0:
        linknum = linknum + int(Reel1symbol[1]) + int(Reel2symbol[1]) + int(Reel3symbol[1]) + int(Reel4symbol[1]) 
        RewardMoney = RewardMoney + GetBaseReward(symbolname, linknum, BetAmount) * R2 * R3 * R4
    elif R3 > 0:  
        linknum = linknum + int(Reel1symbol[1]) + int(Reel2symbol[1]) + int(Reel3symbol[1]) 
        RewardMoney = RewardMoney + GetBaseReward(symbolname, linknum, BetAmount) * R2 * R3
    elif R2 > 0:
        linknum = linknum + int(Reel1symbol[1]) + int(Reel2symbol[1])
        RewardMoney = RewardMoney + GetBaseReward(symbolname, linknum, BetAmount) * R2
    else :
        linknum = linknum + int(Reel1symbol[1])
        RewardMoney = RewardMoney + GetBaseReward(symbolname, linknum, BetAmount) 
    #print(RewardMoney)     
    return(RewardMoney)


def RewardDifferent(RList1, RList2, RList3, RList4, RList5, TotalReward, gresult):
    #print("RewardDifferent")
    Msg1 = RList1[0] + ", " + RList2[0] + ", "  + RList3[0] + ", " + RList4[0] + ", " + RList5[0] + "\n"
    Msg1 += RList1[1] + ", " + RList2[1] + ", " + RList3[1] + ", " + RList4[1] + ", " + RList5[1] + "\n"
    Msg1 += RList1[2] + ", " + RList2[2] + ", " + RList3[2] + ", " + RList4[2] + ", " + RList5[2] + "\n"
    Msg1 += "QA: " + str(TotalReward) + "  User: " + str(gresult)
    #screenshot()
    callbot(Msg1)
    
    
    
msg = "hothotfruit測試開始" +"\n"
callbot(msg) 
usermoney = float(poco('CreditL').attr('text'))
SpinTotalAmount = 0    
TotalBet = 0    
spincount = 0
BetAmount = 15
totalbet = poco('TotalBetL').attr('text')
while totalbet != "15":
    poco("BetPlus").click()
    totalbet = poco('TotalBetL').attr('text')
try:
    while(True):       
        poco('SpinButton').click()
        spincount += 1                
        TotalBet += BetAmount  
        #poco("Bar").click()
        sleep(4)        
        Amount = ReelResult(BetAmount)
        usermoney = usermoney + Amount - BetAmount
        SpinTotalAmount += Amount
        if round(usermoney,2) != round(float(poco('CreditL').attr('text')),2):
            msg1 = "金額計算有誤"+ "\n"
            msg1 += "計算金額:" + str(round(usermoney,2)) + "\n"
            msg1 += "畫面金額:" + poco('CreditL').attr('text')           
            callbot(msg1)
            #screenshot()
            break
        tt = spincount % 100
        if (spincount > 0) and (tt == 0):
            Msg1 = "hothotfruit遊戲次數: " + str(spincount)
            callbot(Msg1)  
except:
    Msg1 = "hothotfruit異常結束,原因為: " +  str(sys.exc_info()[0]) + " ,結束畫面如下:"
    #screenshot()
    callbot(Msg1) 
else:
    callbot("hothotfruit測試完畢")
