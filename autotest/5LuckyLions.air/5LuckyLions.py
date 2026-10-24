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

symbol_s1 = [50, 100, 500, 1000]
symbol_s2 = [40, 80, 450, 900]
symbol_s3 = [40, 60, 400, 800]
symbol_s4 = [30, 55, 350, 700]
symbol_s5 = [30, 50, 300, 600]
symbol_s6 = [20, 45, 250, 500]
symbol_s7 = [20, 40, 250, 500]
symbol_s8 = [10, 35, 150, 300]
symbol_s9 = [10, 30, 100, 200]
symbol_s10 = [5, 25, 75, 4100]
symbol_s11 = [5, 20, 50, 100]
symbol_w1 = [50, 100, 500, 1000]
symbol_c1 = [0, 0, 264, 1320, 4400, 8800]  


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
    elif sname == "w1":
        return(symbol_w1[num])
    elif sname == "c1":
        return(symbol_c1[num])
    
    
def SymbolResult(checkReel):
    SList = []    
    g = poco('ReelMain').child('New Node').child(checkReel).children()
    for i in range(4, 8):
        symbol = str(g[i].attr('name')) 
        SList.append(symbol)  
    return(SList)
    
    
    
def ReelResult(BetAmount):  
    #print("ReelResult")
    gresult = poco('WinEffect').children()    
    RList1 = SymbolResult("Reel1")
    RList2 = SymbolResult("Reel2")
    RList3 = SymbolResult("Reel3")
    RList4 = SymbolResult("Reel4")
    RList5 = SymbolResult("Reel5")
    RList6 = SymbolResult("Reel6")
    TotalReward = CaculateReward(RList1, RList2, RList3, RList4, RList5, RList6, BetAmount) 
    gresult = poco('WinEffect').children()    
    g = gresult[0].attr('name')
    g = float(g.replace(',',''))
    if TotalReward != g:
        RewardDifferent(RList1, RList2, RList3, RList4, RList5, RList6, TotalReward, g)
    return(TotalReward)

def CaculateReward(RList1, RList2, RList3, RList4, RList5, RList6, BetAmount):
    #print("CaculateReward")
    RewardMoney = 0
    linknum = 0
    K = 0
    for i in RList1:
        symbolname = i
        K += 1
        #print(K)
        if K == 1:
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 0, 0, 0, 0, 0, 0) #線條1
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 0, 3, 2, 2, 2, 0) #線條7
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 0, 0, 0, 2, 3, 2) #線條13
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 0, 1, 0, 0, 1, 1) #線條17
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 0, 3, 2, 1, 3, 2) #線條19
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 0, 3, 1, 0, 2, 3) #線條23
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 0, 3, 2, 0, 2, 2) #線條30
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 0, 2, 2, 0, 3, 1) #線條33
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 0, 0, 2, 0, 3, 0) #線條36
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 0, 1, 2, 3, 3, 0) #線條37
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 0, 1, 0, 2, 3, 0) #線條39
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 0, 3, 3, 2, 2, 0) #線條43
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 0, 0, 1, 1, 0, 0) #線條51
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 0, 3, 2, 3, 0, 1) #線條53
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 0, 1, 0, 3, 0, 1) #線條56
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 0, 1, 3, 3, 1, 0) #線條57
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 0, 2, 3, 1, 2, 0) #線條61
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 0, 0, 1, 0, 1, 3) #線條65
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 0, 3, 1, 1, 2, 3) #線條78
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 0, 1, 1, 0, 0, 2) #線條81
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 0, 1, 0, 2, 0, 1) #線條84
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 0, 1, 0, 1, 1, 2) #線條86
                        
        if K == 2:   
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 1, 1, 1, 1, 1, 1) #線條2
                                    
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 1, 0, 1, 1, 1, 0) #線條5    
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 1, 2, 3, 2, 1, 0) #線條10
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 1, 0, 3, 0, 1, 3) #線條12
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 1, 3, 0, 1, 3, 1) #線條16
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 1, 0, 2, 2, 0, 0) #線條22
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 1, 3, 0, 3, 0, 2) #線條25
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 1, 1, 0, 0, 0, 2) #線條28
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 1, 2, 2, 0, 3, 3) #線條31
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 1, 1, 2, 0, 2, 1) #線條41
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 1, 1, 1, 2, 2, 1) #線條46
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 1, 3, 3, 2, 3, 1) #線條48
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 1, 0, 1, 0, 2, 3) #線條50
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 1, 3, 0, 1, 3, 2) #線條60
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 1, 0, 0, 3, 1, 1) #線條63
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 1, 1, 0, 3, 1, 1) #線條67
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 1, 2, 1, 2, 3, 3) #線條69
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 1, 1, 2, 2, 3, 1) #線條72
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 1, 1, 0, 0, 1, 3) #線條74
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 1, 3, 0, 3, 2, 3) #線條75
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 1, 1, 2, 1, 2, 2) #線條79
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 1, 2, 3, 2, 0, 3) #線條87
                                             
        if K == 3:    
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 2, 2, 2, 2, 2, 2) #線條3
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 2, 3, 2, 2, 2, 3) #線條6
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 2, 1, 0, 1, 2, 3) #線條9
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 2, 3, 0, 3, 2, 0) #線條11
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 2, 0, 3, 2, 0, 2) #線條15
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 2, 3, 1, 1, 3, 3) #線條21
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 2, 0, 3, 0, 3, 1) #線條26
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 2, 2, 3, 3, 3, 1) #線條27
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 2, 1, 1, 3, 0, 0) #線條32
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 2, 2, 1, 3, 1, 2) #線條42
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 2, 2, 2, 1, 1, 2) #線條45
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 2, 0, 0, 1, 0, 2) #線條47
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 2, 3, 2, 3, 1, 0) #線條49
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 2, 0, 3, 2, 0, 1) #線條59
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 2, 3, 3, 0, 2, 2) #線條64
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 2, 2, 3, 0, 2, 2) #線條68
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 2, 1, 2, 1, 0, 0) #線條70
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 2, 2, 1, 1, 0, 2) #線條71
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 2, 2, 3, 3, 2, 0) #線條73
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 2, 0, 3, 0, 1, 0) #線條76
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 2, 2, 1, 2, 1, 1) #線條80
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 2, 1, 0, 1, 3, 0) #線條88
                  
        if K == 4: 
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 3, 3, 3, 3, 3, 3) #線條4
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 3, 0, 1, 1, 1, 3) #線條8
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 3, 3, 3, 1, 0, 1) #線條14
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 3, 2, 3, 3, 2, 2) #線條18
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 3, 0, 1, 2, 0, 1) #線條20
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 3, 0, 2, 3, 1, 0) #線條24
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 3, 0, 1, 3, 1, 1) #線條29
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 3, 1, 1, 3, 0, 2) #線條34
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 3, 3, 1, 3, 0, 3) #線條35
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 3, 2, 1, 0, 0, 3) #線條38
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 3, 2, 3, 1, 0, 3) #線條40
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 3, 0, 0, 1, 1, 3) #線條44
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 3, 3, 2, 2, 3, 3) #線條52
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 3, 0, 1, 0, 3, 2) #線條54
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 3, 2, 3, 0, 3, 2) #線條55
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 3, 2, 0, 0, 2, 3) #線條58
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 3, 1, 0, 2, 1, 3) #線條62
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 3, 3, 2, 3, 2, 0) #線條66
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 3, 0, 2, 2, 1, 0) #線條77
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 3, 2, 2, 3, 3, 1) #線條82
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 3, 2, 3, 1, 3, 2) #線條83
            
            RewardMoney = checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, 3, 2, 3, 2, 2, 1) #線條85
            
                       
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
        if RList6[i] == "c1":
            linknum += 1
    if linknum > 0 :
        linknum -= 1
    RewardMoney = RewardMoney + GetBaseReward("c1", linknum, BetAmount) 
    return(RewardMoney)

def checksymbolline(RList1, RList2, RList3, RList4, RList5, RList6, symbolname, RewardMoney, R1list, R2list, R3list, R4list, R5list, R6list): #連線的規則  
    R6 = 0
    R5 = 0
    R4 = 0
    R3 = 0
    R2 = 0
    r4sym = RList4[R4list]
    r5sym = RList5[R5list]
    linknum = 0
    i = symbolname
    if len(r4sym) == 3 :
        symfour = int(r4sym[1:3])
    else:
        symfour = int(r4sym[1])
    if len(r5sym) == 3 :
        symfive = int(r5sym[1:3])
    else:
        symfive = int(r5sym[1])
    if symbolname == "w1" and RList2[R2list] != "c1":
        symbolname = RList2[R2list]
        if symbolname == "w1" and RList3[R3list] != "c1":
            symbolname = RList3[R3list]
            if symbolname == "w1" and RList4[R4list] != "c1" and (symfour < 6  or RList5[R5list] == "w1" or RList4[R4list] == RList5[R5list]) :
                symbolname = RList4[R4list]
                if symbolname == "w1" and RList5[R5list] != "c1" and (symfive < 10   or RList6[R6list] == "w1" or RList5[R5list] == RList6[R6list]):
                    symbolname = RList5[R5list]
                    if symbolname == "w1" and RList6[R6list] != "c1":
                        symbolname = RList6[R6list]
    #print(symbolname)       
    if symbolname != "c1" and symbolname != "w1":
        if (RList2[R2list] == symbolname) or (RList2[R2list].count("w1") != 0):
            R2 = R2 + RList2[R2list].count(symbolname) + RList2[R2list].count("w1")
            if (RList3[R3list] == symbolname) or (RList3[R3list].count("w1") != 0) :
                R3 = R3 + RList3[R3list].count(symbolname) + RList3[R3list].count("w1")
                if (RList4[R4list] == symbolname) or (RList4[R4list].count("w1") != 0):
                    linknum = 1
                    R4 = R4 + RList4[R4list].count(symbolname) + RList4[R4list].count("w1")
                    if(RList5[R5list] == symbolname) or (RList5[R5list].count("w1") != 0):
                        linknum = 2
                        R5 = R5 + RList5[R5list].count(symbolname) + RList5[R5list].count("w1")
                        if(RList6[R6list] == symbolname) or (RList6[R6list].count("w1") != 0):
                            linknum = 3
                            R6 = R6 + RList6[R6list].count(symbolname) + RList6[R6list].count("w1")
    elif symbolname != "c1" and symbolname == "w1":
        if (RList2[R2list] == symbolname) :
            R2 = R2 + RList2[R2list].count(symbolname) 
            if (RList3[R3list] == symbolname)  :
                R3 = R3 + RList3[R3list].count(symbolname) 
                if (RList4[R4list] == symbolname) :
                    linknum = 1
                    R4 = R4 + RList4[R4list].count(symbolname) 
                    if(RList5[R5list] == symbolname) :
                        linknum = 2
                        R5 = R5 + RList5[R5list].count(symbolname)
                        if(RList6[R6list] == symbolname) :
                            linknum = 3
                            R6 = R6 + RList6[R6list].count(symbolname) 
    
    if R6 > 0:            
        RewardMoney = RewardMoney + GetBaseReward(symbolname, linknum, BetAmount) * R2 * R3 * R4 * R5 * R6
    elif R5 > 0:            
        RewardMoney = RewardMoney + GetBaseReward(symbolname, linknum, BetAmount) * R2 * R3 * R4 * R5
    elif R4 > 0:            
        RewardMoney = RewardMoney + GetBaseReward(symbolname, linknum, BetAmount) * R2 * R3 * R4
    elif R3 > 0:    
        RewardMoney = RewardMoney + GetBaseReward(symbolname, linknum, BetAmount) * R2 * R3
    symbolname = i    
    #print(RewardMoney)
    return(RewardMoney)


def RewardDifferent(RList1, RList2, RList3, RList4, RList5, RList6, TotalReward, gresult):
    #print("RewardDifferent")
    Msg1 = RList1[0] + ", " + RList2[0] + ", "  + RList3[0] + ", " + RList4[0] + ", " + RList5[0] + ", " + RList6[0] +"\n"
    Msg1 += RList1[1] + ", " + RList2[1] + ", " + RList3[1] + ", " + RList4[1] + ", " + RList5[1] + ", " + RList6[1] +"\n"
    Msg1 += RList1[2] + ", " + RList2[2] + ", " + RList3[2] + ", " + RList4[2] + ", " + RList5[2] + ", " + RList6[2] +"\n"
    Msg1 += RList1[3] + ", " + RList2[3] + ", " + RList3[3] + ", " + RList4[3] + ", " + RList5[3] + ", " + RList6[3] +"\n"
    Msg1 += "QA: " + str(TotalReward) + "  User: " + str(gresult)
    screenshot()
    callbot(Msg1)
    
    
    
msg = "5Luckylions測試開始" +"\n"
callbot(msg) 
usermoney = float(poco('CreditL').attr('text'))
SpinTotalAmount = 0    
TotalBet = 0    
spincount = 0
BetAmount = 88
totalbet = poco('TotalBetL').attr('text')
while totalbet != "88":
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
            #screenshot()
            break
        tt = spincount % 100
        if (spincount > 0) and (tt == 0):
            Msg1 = "5Luckylions遊戲次數: " + str(spincount)
            callbot(Msg1)  
except:
    Msg1 = "5Luckylions異常結束,原因為: " +  str(sys.exc_info()[0]) + " ,結束畫面如下:"
    screenshot()
    callbot(Msg1) 
else:
    callbot("5Luckylions測試完畢")
