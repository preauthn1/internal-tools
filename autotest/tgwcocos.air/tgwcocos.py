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


symbol_p1 = [0.3, 2.7, 18, 60]
symbol_p2 = [0, 2.1, 12, 45]
symbol_p3 = [0, 1.5, 6, 30]
symbol_p4 = [0, 0.9, 4.5, 15]
symbol_p5 = [0, 0.6, 2.4, 6]
symbol_p6 = [0, 0.6, 0.9, 3]
symbol_p7 = [0, 0.3, 0.9, 3]
symbol_p8 = [0, 0.3, 0.9, 3]
symbol_p9 = [0, 0.3, 0.9, 3]
symbol_p10 = [5, 10, 50]
symbol_p11 = [5, 10, 50]
symbol_p12 = [100, 200, 1000]
symbol_p13 = [50, 100, 500]
symbol_p14 = [40, 80, 400]
symbol_p15 = [0, 0, 0, 0]
symbol_w1 = [0.6, 3, 30, 150]
SymbolList = [2, 3, 4, 1, 2, 3, 5, 1, 2, 4, 5, 1, 3, 4, 5]

def callbot(Msg):
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
    
def SymbolResult(checkReel):
    #print('SymbolResult')
    SList = []
    g = poco("UISlot").child(checkReel).children()
    for i in range(0, 3):
        symbol = g[i].attr('name')
        s = poco("UISlot").child(checkReel).child(symbol).children()    
        s1 = len(s)
        s2 = s[(s1-2)].attr('name')
        if len(s2) == 23:
            s2 = s2[20:22]
        else:
            s2 = s2[20:23]
        SList.append(s2)  
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
    elif sname == "w1":
        return(symbol_w1[num])
    else:
        return(symbol_p15[num])
    
    
def CaculateReward(RList1, RList2, RList3, RList4, RList5, BetAmount):
    #print("CaculateReward")
    RewardMoney = 0
    linknum = 0
    K = 0
    for i in RList1:
        symbolname = i
        K += 1
        #print(K)
        if K == 1:
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, symbolname, RewardMoney, 0, 0, 0, 0, 0) #線條2
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, symbolname, RewardMoney, 0, 1, 2, 1, 0) #線條4
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, symbolname, RewardMoney, 0, 0, 1, 2, 2) #線條8
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, symbolname, RewardMoney, 0, 1, 1, 1, 0) #線條12
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, symbolname, RewardMoney, 0, 1, 0, 1, 0) #線條14
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, symbolname, RewardMoney, 0, 0, 2, 0, 0) #線條18
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, symbolname, RewardMoney, 0, 2, 2, 2, 0) #線條20
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, symbolname, RewardMoney, 0, 0, 1, 0, 0) #線條22
        if K == 2:
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, symbolname, RewardMoney, 1, 1, 1, 1, 1) #線條1
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, symbolname, RewardMoney, 1, 0, 0, 0, 1) #線條6
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, symbolname, RewardMoney, 1, 2, 2, 2, 1) #線條7
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, symbolname, RewardMoney, 1, 0, 1, 2, 1) #線條10
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, symbolname, RewardMoney, 1, 2, 1, 0, 1) #線條11
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, symbolname, RewardMoney, 1, 1, 0, 1, 1) #線條16
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, symbolname, RewardMoney, 1, 1, 2, 1, 1) #線條17
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, symbolname, RewardMoney, 1, 0, 1, 0, 1) #線條24
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, symbolname, RewardMoney, 1, 2, 1, 2, 1) #線條25
                     
        if K == 3:
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, symbolname, RewardMoney, 2, 2, 2, 2, 2) #線條3
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, symbolname, RewardMoney, 2, 1, 0, 1, 2) #線條5
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, symbolname, RewardMoney, 2, 2, 1, 0, 0) #線條9
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, symbolname, RewardMoney, 2, 1, 1, 1, 2) #線條13
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, symbolname, RewardMoney, 2, 1, 2, 1, 2) #線條15
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, symbolname, RewardMoney, 2, 2, 0, 2, 2) #線條19
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, symbolname, RewardMoney, 2, 0, 0, 0, 2) #線條21
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, symbolname, RewardMoney, 2, 2, 1, 2, 2) #線條23
                       
    return(RewardMoney)     


def checksymbolline(RList1, RList2, RList3, RList4, RList5,  symbolname, RewardMoney, R1list, R2list, R3list, R4list, R5list): #連線的規則  
    R5 = 0
    R4 = 0
    R3 = 0
    R2 = 0
    r3sym = RList3[R3list]
    r4sym = RList4[R4list]
    linknum = 0
    i = symbolname
    if len(r4sym) == 3 :
        symfour = int(r4sym[1:3])
    else:
        symfour = int(r4sym[1])
    if len(r3sym) == 3 :
        symthree = int(r3sym[1:3])
    else:
        symthree = int(r3sym[1])
    if symbolname == "w1" and RList2[R2list] != "b1":
        symbolname = RList2[R2list]
        if symbolname == "w1" and RList3[R3list] != "b1"and (symthree < 6  or RList4[R4list] == "w1" or RList3[R3list] == RList4[R4list]):
            symbolname = RList3[R3list]
            if symbolname == "w1" and RList4[R4list] != "b1" and (symfour < 5  or RList5[R5list] == "w1" or RList4[R4list] == RList5[R5list]) :
                symbolname = RList4[R4list]
                if symbolname == "w1" and RList5[R5list] != "b1" :
                    symbolname = RList5[R5list]
                    
    #print(symbolname)       
    if symbolname != "b1" and symbolname != "w1":
        if (RList2[R2list] == symbolname) or (RList2[R2list].count("w1") != 0):
            R2 = R2 + RList2[R2list].count(symbolname) + RList2[R2list].count("w1")
            if (RList3[R3list] == symbolname) or (RList3[R3list].count("w1") != 0) :
                linknum = 1
                R3 = R3 + RList3[R3list].count(symbolname) + RList3[R3list].count("w1")
                if (RList4[R4list] == symbolname) or (RList4[R4list].count("w1") != 0):
                    linknum = 2
                    R4 = R4 + RList4[R4list].count(symbolname) + RList4[R4list].count("w1")
                    if(RList5[R5list] == symbolname) or (RList5[R5list].count("w1") != 0):
                        linknum = 3
                        R5 = R5 + RList5[R5list].count(symbolname) + RList5[R5list].count("w1")
                        
    elif symbolname != "b1" and symbolname == "w1":
        if (RList2[R2list] == symbolname) :
            R2 = R2 + RList2[R2list].count(symbolname) 
            if (RList3[R3list] == symbolname)  :
                R3 = R3 + RList3[R3list].count(symbolname)
                linknum = 1
                if (RList4[R4list] == symbolname) :
                    linknum = 2
                    R4 = R4 + RList4[R4list].count(symbolname) 
                    if(RList5[R5list] == symbolname) :
                        linknum = 3
                        R5 = R5 + RList5[R5list].count(symbolname)
    if symbolname == "w1" and R5 > 0:
        jackpot = float(poco("Jackpot1").child('TextSystem').attr('text'))
        symbol_w1 = [0.6, 3, 30, jackpot]
        
    if R5 > 0:            
        RewardMoney = round(float(RewardMoney + GetBaseReward(symbolname, linknum, BetAmount) * R2 * R3 * R4 * R5),2) 
    elif R4 > 0:            
        RewardMoney = round(float(RewardMoney + GetBaseReward(symbolname, linknum, BetAmount) * R2 * R3 * R4),2)
    elif R3 > 0:    
        RewardMoney = round(float(RewardMoney + GetBaseReward(symbolname, linknum, BetAmount) * R2 * R3),2)
    elif R2 > 0:    
        RewardMoney = round(float(RewardMoney + GetBaseReward(symbolname, linknum, BetAmount) * R2),2) 
    symbolname = i    
    #print(RewardMoney)
    return(RewardMoney)

def RewardDifferent(RList1, RList2, RList3, RList4, RList5, TotalReward, gresult, FAmount):
    #print("RewardDifferent")
    Msg1 = RList1[0] + ", " + RList2[0] + ", "  + RList3[0] + ", " + RList4[0] + ", " + RList5[0] + "\n"
    Msg1 += RList1[1] + ", " + RList2[1] + ", " + RList3[1] + ", " + RList4[1] + ", " + RList5[1] + "\n"
    Msg1 += RList1[2] + ", " + RList2[2] + ", " + RList3[2] + ", " + RList4[2] + ", " + RList5[2] + "\n"
    if FAmount > 0:
        Msg1 += "QA: " + str(TotalReward + FAmount) + "  User: " + str(gresult)
        Msg1 += "上盤金額: " + str(FAmount)
    else:
        Msg1 += "QA: " + str(TotalReward) + "  User: " + str(gresult)
    callbot(Msg1)
    #screenshot()
    RList1 = SymbolResult("Reel6")
    
def ReelResult(FAmount, BetAmount):  
    #print("ReelResult")
    pay1 = ""
    RList1 = SymbolResult("Reel1")
    RList2 = SymbolResult("Reel2")
    RList3 = SymbolResult("Reel3")
    RList4 = SymbolResult("Reel4")
    RList5 = SymbolResult("Reel5")    
    TotalReward = CaculateReward(RList1, RList2, RList3, RList4, RList5, BetAmount) + float(FAmount)
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
msg = "糖果屋測試開始" +"\n"
callbot(msg) 
isneedclick = True 
SpinTotalAmount = 0
spincount = 0
TotalBet = 0
Totalspin = 50
doublenum = 2
BetAmount = 1.5 
try:
    while (True):
        tt = spincount % 100
        if (spincount > 0) and (tt == 0):
            Msg1 = "遊戲次數: " + str(spincount) 
            callbot(Msg1)
        if poco("TextState").attr("text") == "start":
            poco("ButtonStart").click()
            spincount += 1                
            TotalBet += BetAmount                
        while poco("TextState").attr("text") == "spin":
            time.sleep(0.1)  
        Amount = ReelResult(0, BetAmount) 
        SpinTotalAmount += Amount    
except:
    Msg1 = "糖果屋異常結束,原因為: " +  str(sys.exc_info()[0]) + " ,結束畫面如下:"
    #screenshot()
    callbot(Msg1) 
else:
    callbot("糖果屋測試完畢")