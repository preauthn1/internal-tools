# -*- encoding=utf8 -*-
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


symbol_p1 = [10, 20, 100]
symbol_p2 = [5, 10, 50]
symbol_p3 = [5, 10, 50]
symbol_p4 = [5, 10, 50]
symbol_p5 = [5, 10, 50]
symbol_p6 = [5, 10, 50]
symbol_p7 = [5, 10, 50]
symbol_p8 = [5, 10, 50]
symbol_p9 = [5, 10, 50]
symbol_p10 = [5, 10, 50]
symbol_p11 = [5, 10, 50]
symbol_p12 = [100, 200, 1000]
symbol_p13 = [50, 100, 500]
symbol_p14 = [40, 80, 400]
symbol_p15 = [25, 50, 250]
symbol_s1 = [40, 80, 400]
SymbolList = [2, 3, 4, 1, 2, 3, 5, 1, 2, 4, 5, 1, 3, 4, 5]

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
    
def SymbolResult(checkReel):
    SList = []
    g = poco("UISlot").child(checkReel).child("Symbol_1").children()
    y =g[0].attr('pos')
    y1 = int(y[1] * 10) 
    if y1 == 12:
        y1 = 1
    y2 = int(y1 // 2)  
    for i in range(0, 3):
        k = y2 * 3 + i
        symbol = "Symbol_" + str(SymbolList[k])  
        s = poco("UISlot").child(checkReel).child(symbol).children()    
        s1 = s[0].attr('name')
        SList.append(s1)  
    return(SList)

def GetBaseReward(sname, num, BetAmount):
    if sname == "p1":
        return(symbol_p1[num])
    elif sname == "p2":
        return(symbol_p2[num])
    elif sname == "p3":
        return(symbol_p3[num])
    elif sname == "p4":
        return(symbol_p4[num])
    elif sname == "p5":
        return(symbol_p5[num])
    elif sname == "p6":
        return(symbol_p6[num])
    elif sname == "p7":
        return(symbol_p7[num])
    elif sname == "p8":
        return(symbol_p8[num])
    elif sname == "p9":
        return(symbol_p9[num])
    elif sname == "p10":
        return(symbol_p10[num])
    elif sname == "p11":
        return(symbol_p11[num])
    elif sname == "p12":
        return(symbol_p12[num])
    elif sname == "p13":
        return(symbol_p13[num])
    elif sname == "p14":
        return(symbol_p14[num])
    elif sname == "s1":
        return(symbol_s1[num])
    else:
        return(symbol_p15[num])
    
    
def CaculateReward(RList1, RList2, RList3, RList4, RList5, BetAmount):
    print("CaculateReward")
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

def RewardDifferent(RList1, RList2, RList3, RList4, RList5, TotalReward, gresult, FAmount):
    print("RewardDifferent")
    Msg1 = RList1[0] + ", " + RList2[0] + ", "  + RList3[0] + ", " + RList4[0] + ", " + RList5[0] + "\n"
    Msg1 += RList1[1] + ", " + RList2[1] + ", " + RList3[1] + ", " + RList4[1] + ", " + RList5[1] + "\n"
    Msg1 += RList1[2] + ", " + RList2[2] + ", " + RList3[2] + ", " + RList4[2] + ", " + RList5[2] + "\n"
    if FAmount > 0:
        Msg1 += "QA: " + str(TotalReward + FAmount) + "  User: " + str(gresult)
        Msg1 += "上盤金額: " + str(FAmount)
    else:
        Msg1 += "QA: " + str(TotalReward) + "  User: " + str(gresult)
    callbot(Msg1)
    screenshot()
    
def ReelResult(FAmount, BetAmount):  
    print("ReelResult")
    pay1 = ""
    RList1 = SymbolResult("Reel1")
    RList2 = SymbolResult("Reel2")
    RList3 = SymbolResult("Reel3")
    RList4 = SymbolResult("Reel4")
    RList5 = SymbolResult("Reel5")    
    TotalReward = CaculateReward(RList1, RList2, RList3, RList4, RList5, BetAmount) + float(FAmount)
    poco("ImageTotalBet").click()
    gresult = poco("UIUser").child("TextWinMoney").attr("text")
    if (gresult == ""):
        gresult = 0
    if TotalReward != float(gresult):
        RewardDifferent(RList1, RList2, RList3, RList4, RList5, TotalReward, gresult, FAmount)
    return(TotalReward)


def CompareAmount(UserAmount, TotalBet, spincount, SpinTotalAmount):
    print("CompareAmount")    
    UserAmount1 = round((UserAmount - float(TotalBet) + SpinTotalAmount), 2) 
    PAmount =  poco("UIUser").child("TextCarryMoney").attr("text")
    PAmount =  round(float(PAmount.replace(',','')),2)
    if UserAmount1 != PAmount:
        Msg1 = "玩家金額不正確;玩家初始金額: " + str(UserAmount) 
        Msg1 += " 總贏金: " + str(SpinTotalAmount) 
        Msg1 += " 總下注: " + str(TotalBet) + " (遊戲次數: " + str(spincount) + ")"
        Msg1 += " 預期金額: " + str(UserAmount1) + " 目前結算金額: " + str(PAmount)
        callbot(Msg1) 
        return(-1)
    else:
        return(PAmount)

  
isfirstspin = True                                
UserAmount =  poco("UIUser").child("TextCarryMoney").attr("text")
UserAmount =  float(UserAmount.replace(',',''))

isneedclick = True 
SpinTotalAmount = 0
spincount = 0
TotalBet = 0
Totalspin = 50
doublenum = 2
BetAmount = 8
betstr=poco("ButtonChip").child("Text").attr("text")
while betstr != "X1.00":
    poco("ButtonChipUp").click()
    betstr=poco("ButtonChip").child("Text").attr("text")    
while (True):
    tt = spincount % 100
    if (spincount > 0) and (tt == 0):
        PAmount = CompareAmount(UserAmount, TotalBet, spincount, SpinTotalAmount)
        if PAmount == -1:
            break                
        p = int(PAmount) / int(UserAmount) * 100
        p = round(p, 2)
        Msg1 = "遊戲次數: " + str(spincount) + " 目前payout ratio = " + str(p) + "%"
        callbot(Msg1)
    if poco("TextState").attr("text") == "start":
        poco("ButtonStart").click()
        spincount += 1                
        TotalBet += BetAmount                
    while poco("TextState").attr("text") == "spin":
        time.sleep(0.1)            
    Amount = ReelResult(0, BetAmount) 
    SpinTotalAmount += Amount    