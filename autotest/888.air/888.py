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

symbol_s1 = [100, 500, 2000]
symbol_s2 = [50, 100, 500]
symbol_s3 = [15, 50, 200]
symbol_s4 = [15, 40, 150]
symbol_s5 = [10, 25, 100]
symbol_s6 = [10, 20, 75]
symbol_s7 = [10, 15, 50]
symbol_s8 = [20, 30, 50]  #混三八
symbol_s9 = [0, 0, 10]  #混bar

def callbot(Msg):
    r=requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=%s"%Msg)
    #r=requests.get("https://example.invalid/bot871711435:AAGYzFP4tNBQ-5E2j45a8qFBiDZDbmIlY2M/sendMessage?chat_id=-216697951&text=%s"%Msg)   #Lisa
    return()
   
def screenshot():
    saveDir= "D:\\QA\\QA_AutoTest\\temp"     
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


def GetBaseReward(sname, num, BetAmount): #賠率表
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
    else: 
        return(symbol_s9[num])
    
    
def SymbolResult(checkReel):  #抓輪軸資料
    SList = []    
    g = poco('ReelMain').child('New Node').child(checkReel).children()
    for i in range(3, 6):
        symbol = str(g[i].attr('name')) 
        SList.append(symbol)  
    return(SList)
    
    
    
def ReelResult(BetAmount):  #抓輪軸資料並計算金額
    RList1 = SymbolResult("Reel1")
    RList2 = SymbolResult("Reel2")
    RList3 = SymbolResult("Reel3")
    RList4 = SymbolResult("Reel4")
    RList5 = SymbolResult("Reel5")
    TotalReward = CaculateReward(RList1, RList2, RList3, RList4, RList5, BetAmount)
    gresult = poco('WinEffect').children()    #抓贏得金額的資料
    g = gresult[0].attr('name')
    g = float(g.replace(',',''))
    if TotalReward != g:
        RewardDifferent(RList1, RList2, RList3, RList4, RList5, TotalReward, g)
    return(TotalReward)

def CaculateReward(RList1, RList2, RList3, RList4, RList5, BetAmount): #計算連線
    RewardMoney = 0
    K = 0
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
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5,  symbolname, RewardMoney, 0, 0, 0, 0, 0) #線條2
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5,  symbolname, RewardMoney, 0, 1, 2, 1, 0) #線條4
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5,  symbolname, RewardMoney, 0, 0, 1, 2, 2) #線條6
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5,  symbolname, RewardMoney, 0, 1, 1, 1, 2) #線條10
            
        if K == 2:   
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5,  symbolname, RewardMoney, 1, 1, 1, 1, 1) #線條1
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5,  symbolname, RewardMoney, 1, 0, 1, 2, 1) #線條8
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5,  symbolname, RewardMoney, 1, 2, 1, 0, 1) #線條9
            
        if K == 3:    
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5,  symbolname, RewardMoney, 2, 2, 2, 2, 2) #線條3
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5,  symbolname, RewardMoney, 2, 1, 0, 1, 2) #線條5
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5,  symbolname, RewardMoney, 2, 2, 1, 0, 0) #線條7            
            
    return(RewardMoney)

def checksymbolline(RList1, RList2, RList3, RList4, RList5, symbolname, RewardMoney, R1list, R2list, R3list, R4list, R5list): #連線的規則  
    R5 = 0
    R4 = 0
    R3 = 0
    R2 = 0
    linknum = 0
    Reel2symbol = RList2[R2list]
    Reel3symbol = RList3[R3list]
    Reel4symbol = RList4[R4list]
    Reel5symbol = RList5[R5list]
    if (RList2[R2list].count(symbolname) != 0) :
        R2 = R2 + RList2[R2list].count(symbolname) 
        if (RList3[R3list].count(symbolname) != 0):
            R3 = R3 + RList3[R3list].count(symbolname)
            if (RList4[R4list].count(symbolname) != 0):
                linknum = 1
                R4 = R4 + RList4[R4list].count(symbolname)
                if(RList5[R5list].count(symbolname) != 0):
                    linknum = 2
                    R5 = R5 + RList5[R5list].count(symbolname)   
    
    if R3 == 0:
        if symbolname == "s1" or symbolname == "s2":
            symbolname = "s8"
        if symbolname == "s5" or symbolname == "s6" or symbolname == "s7":
            symbolname = "s9"
        if Reel2symbol == "s1" or Reel2symbol == "s2":
            Reel2symbol = "s8"
        if Reel2symbol == "s5" or Reel2symbol == "s6" or Reel2symbol == "s7":
            Reel2symbol = "s9"    
        if Reel3symbol == "s1" or Reel3symbol == "s2":
            Reel3symbol = "s8"
        if Reel3symbol == "s5" or Reel3symbol == "s6" or Reel3symbol == "s7":
            Reel3symbol = "s9"    
        if Reel4symbol == "s1" or Reel4symbol == "s2":
            Reel4symbol = "s8"
        if Reel4symbol == "s5" or Reel4symbol == "s6" or Reel4symbol == "s7":
            Reel4symbol = "s9"
        if Reel5symbol == "s1" or Reel5symbol == "s2":
            Reel5symbol = "s8"
        if Reel5symbol == "s5" or Reel5symbol == "s6" or Reel5symbol == "s7":
            Reel5symbol = "s9"    
        if (Reel2symbol.count(symbolname) != 0) :
            R2 = R2 + Reel2symbol.count(symbolname) 
            if (Reel3symbol.count(symbolname) != 0):
                R3 = R3 + Reel3symbol.count(symbolname)
                if (Reel4symbol.count(symbolname) != 0):
                    linknum = 1
                    R4 = R4 + Reel4symbol.count(symbolname)
                    if(Reel5symbol.count(symbolname) != 0):
                        linknum = 2
                        R5 = R5 + Reel5symbol.count(symbolname)
    
    if R5 > 0:            
        RewardMoney = RewardMoney + GetBaseReward(symbolname, linknum, BetAmount) * R3 * R4 * R5 
    elif R4 > 0:            
        RewardMoney = RewardMoney + GetBaseReward(symbolname, linknum, BetAmount) * R3 * R4
    elif R3 > 0:            
        RewardMoney = RewardMoney + GetBaseReward(symbolname, linknum, BetAmount) * R3 
        
    return(RewardMoney)


def RewardDifferent(RList1, RList2, RList3, RList4, RList5, TotalReward, gresult):
    print("RewardDifferent")
    Msg1 = RList1[0] + ", " + RList2[0] + ", "  + RList3[0] + ", " + RList4[0] + ", " + RList5[0] + "\n"
    Msg1 += RList1[1] + ", " + RList2[1] + ", " + RList3[1] + ", " + RList4[1] + ", " + RList5[1] + "\n"
    Msg1 += RList1[2] + ", " + RList2[2] + ", " + RList3[2] + ", " + RList4[2] + ", " + RList5[2] + "\n"
    Msg1 += "QA: " + str(TotalReward) + "  User: " + str(gresult)
    screenshot()
    callbot(Msg1)
    
    
    
msg = "888測試開始" +"\n"
callbot(msg) 
usermoney = float(poco('CreditL').attr('text'))
SpinTotalAmount = 0    
TotalBet = 0    
spincount = 0
BetAmount = 10
totalbet = poco('TotalBetL').attr('text')
while totalbet != "10":
    poco("BetPlus").click()
    totalbet = poco('TotalBetL').attr('text') 
try:
    while(True):       
        poco('SpinButton').click()
        spincount += 1                
        TotalBet += BetAmount  
        time.sleep(10)        
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
            Msg1 = "888遊戲次數: " + str(spincount)
            callbot(Msg1)  
except:
    Msg1 = "888異常結束,原因為: " +  str(sys.exc_info()[0]) + " ,結束畫面如下:"
    screenshot()
    callbot(Msg1) 
else:
    callbot("888測試完畢")
