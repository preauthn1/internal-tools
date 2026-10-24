# -*- encoding=utf8 -*-
__author__ = "Frank"

from PIL import Image
from PIL import ImageGrab
import time
import os
import sys
import requests
from airtest.core.api import *
from poco.drivers.unity3d import UnityPoco
poco = UnityPoco()
auto_setup(__file__)

symbol_p1 = [1.5, 3, 5, 10]
symbol_p2 = [0, 3, 5, 10]
symbol_p3 = [0, 3, 5, 7]
symbol_p4 = [0, 3, 5, 7]
symbol_p5 = [0, 3, 5, 7]
symbol_p6 = [0, 1, 3, 5]
symbol_p7 = [0, 1, 3, 5]
symbol_p8 = [0, 1, 2, 3]
symbol_p9 = [0, 1, 2, 3]
symbol_p10 = [0, 0.5, 1, 2]
symbol_p11 = [0, 0.5, 1, 2]
symbol_p12 = [100, 200, 1000]
symbol_p13 = [50, 100, 500]
symbol_p14 = [40, 80, 400]
symbol_p15 = [0, 0, 0, 0]
SymbolList = [2, 3, 4, 1, 2, 3, 5, 1, 2, 4, 5, 1, 3, 4, 5]
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

def timediff():
    runtime = time.time()    
    return(int((runtime - starttime) // 60))  
    
def SymbolResult(checkReel):
    SList = []
    y = poco("UISlot").child(checkReel).children()
    for i in range(0, 4):
        symbol = str(y[i].attr('name'))  
        s = poco("UISlot").child(checkReel).child(symbol).children().attr("name")        
        if len(s) > 30 :
            s1 = "symbol_"+ s[20:23]
        else:
            s1 = "symbol_"+ s[20:22]
        SList.append(s1)  
    return(SList)
    
def GetBaseReward(sname, num, BetAmount):
    if sname == "symbol_p1":
        return(symbol_p1[num])
    elif sname == "symbol_p2":
        return(symbol_p2[num])
    elif sname == "symbol_p3":
        return(symbol_p3[num])
    elif sname == "symbol_p4":
        return(symbol_p4[num])
    elif sname == "symbol_p5":
        return(symbol_p5[num])
    elif sname == "symbol_p6":
        return(symbol_p6[num])
    elif sname == "symbol_p7":
        return(symbol_p7[num])
    elif sname == "symbol_p8":
        return(symbol_p8[num])
    elif sname == "symbol_p9":
        return(symbol_p9[num])
    elif sname == "symbol_p10":
        return(symbol_p10[num])
    elif sname == "symbol_p11":
        return(symbol_p11[num])
    elif sname == "symbol_s1":
        return(symbol_s1[num])

def CaculateReward(RList1, RList2, RList3, RList4, RList5, BetAmount):
    print("CaculateReward")
    RewardMoney = 0
    for i in RList1:
        R2 = 0
        R3 = 0
        R4 = 0
        R5 = 0
        linknum = 0 
        if(RList2.count(i) != 0) or (RList2.count("symbol_w1") != 0):
            R2 = R2 + RList2.count(i) + RList2.count("symbol_w1")
            if (RList3.count(i) != 0) or (RList3.count("symbol_w1") != 0):
                linknum = 1
                R3 = R3 + RList3.count(i) + RList3.count("symbol_w1")
                if (RList4.count(i) != 0) or (RList4.count("symbol_w1") != 0):
                    linknum = 2
                    R4 = R4 + RList4.count(i) + RList4.count("symbol_w1")
                    if(RList5.count(i) != 0) or (RList5.count("symbol_w1") != 0):
                        linknum = 3
                        R5 = R5 + RList5.count(i) + RList5.count("symbol_w1")
        if R5 > 0:            
            RewardMoney = RewardMoney + GetBaseReward(i, linknum, BetAmount) * R2 * R3 * R4 * R5 
        elif R4 > 0:            
            RewardMoney = RewardMoney + GetBaseReward(i, linknum, BetAmount) * R2 * R3 * R4
        elif R3 > 0:            
            RewardMoney = RewardMoney + GetBaseReward(i, linknum, BetAmount) * R2 * R3   
        elif R2 > 0:            
            RewardMoney = RewardMoney + GetBaseReward(i, linknum, BetAmount) * R2
    
    return(RewardMoney)     

def RewardDifferent(RList1, RList2, RList3, RList4, RList5, TotalReward, gresult, FAmount):
    print("RewardDifferent")
    screenshot()
    Msg1 = RList1[0] + ", " + RList2[0] + ", "  + RList3[0] + ", " + RList4[0] + ", " + RList5[0] + "\n"
    Msg1 += RList1[1] + ", " + RList2[1] + ", " + RList3[1] + ", " + RList4[1] + ", " + RList5[1] + "\n"
    Msg1 += RList1[2] + ", " + RList2[2] + ", " + RList3[2] + ", " + RList4[2] + ", " + RList5[2] + "\n"
    Msg1 += RList1[3] + ", " + RList2[3] + ", " + RList3[3] + ", " + RList4[3] + ", " + RList5[3] + "\n"
    if FAmount > 0:
        Msg1 += "QA: " + str(TotalReward + FAmount) + "  User: " + str(gresult)
        Msg1 += "上盤金額: " + str(FAmount)
    else:
        Msg1 += "QA: " + str(TotalReward) + "  User: " + str(gresult)
    callbot(Msg1)        
    
def ReelResult(FAmount, BetAmount):  
    print("ReelResult")
    pay1 = ""
    RList1 = SymbolResult("Reel1")
    RList2 = SymbolResult("Reel2")
    RList3 = SymbolResult("Reel3")
    RList4 = SymbolResult("Reel4")
    RList5 = SymbolResult("Reel5")    
    TotalReward = CaculateReward(RList1, RList2, RList3, RList4, RList5, BetAmount) + float(FAmount)
    poco("Button1").click()
    poco("Button1").click()
    poco("Button1").click()
    gresult = poco("TextWinMoney").attr("text")
    pay1 = gresult
    if (gresult == "") or (pay1 == ""):
        gresult = 0
        pay1 = 0
    if (TotalReward != float(gresult)) and (TotalReward != float(pay1)):
        RewardDifferent(RList1, RList2, RList3, RList4, RList5, TotalReward, gresult, FAmount)
    return(TotalReward)

def GetJPAmount():    
    print("GetJPAmount")    
    JPList = []
    for a in range(1, 5):
        flag = "Jackpot" + str(a)
        s = poco("UIJackpotPool").child(flag).child("TextSystem").attr("text")
        s1 = ""
        for j in s:
            if j != ",":
                s1 += j
        JPList.append(s1)
    return(JPList)

def CollectBonusPic(o, p, BonusPicList):
    print("CollectBonusPic")    
    Group = "Group" + str(o)
    ButtonOption = "ButtonOption" + str(p)
    for i in range(1, 5):
        image = "Image" + str(i)
        if poco("UIBonusGame").child("ImageBg").child(Group).child(ButtonOption).child(image).exists():           
            BonusPicList.append(image)
            break    
    return(BonusPicList)

def CompareAmount(UserAmount, TotalBet, spincount, SpinTotalAmount):
    print("CompareAmount")    
    UserAmount1 = round((UserAmount - float(TotalBet) + SpinTotalAmount), 2) 
    PAmount = poco("TextCarryMoney").attr("text")
    PAmount = float(UserAmount.replace(',',''))
    if UserAmount1 != PAmount:
        #screenshot()
        Msg1 = "玩家金額不正確;玩家初始金額: " + str(UserAmount) 
        Msg1 += " 總贏金: " + str(SpinTotalAmount) 
        Msg1 += " 總下注: " + str(TotalBet) + " (遊戲次數: " + str(spincount) + ")"
        Msg1 += " 預期金額: " + str(UserAmount1) + " 目前結算金額: " + str(PAmount)
        callbot(Msg1) 
        return(PAmount)
    else:
        return(PAmount)


try: 
    isfirstspin = True       
    for i in range(0, 2):                          
        UserAmount = poco("TextCarryMoney").attr("text")
        UserAmount = float(UserAmount.replace(',',''))
        starttime = time.time()
        isneedclick = True 
        isFreespin = False
        isBonus = False      
        SpinTotalAmount = 0
        spincount = 0
        TotalBet = 0
        Totalspin = 50
        doublenum = 2
        if i == 0:            
            BetAmount = 2.8                       
            Msg1 = "東海之王測試開始,游戲底住2.8, 每局押注金額" + str(BetAmount) + "元,玩家起始金額" + str(UserAmount)            
            callbot(Msg1)                       
        else:            
            BetAmount = 2.8                       
            Msg1 = "東海之王測試開始,游戲底住2.8, 每局押注金額" + str(BetAmount) + "元,玩家起始金額" + str(UserAmount)            
            callbot(Msg1)                        
        while (timediff() < 480):
            tt = spincount % Totalspin
            if (spincount > 0) and (tt == 0):
                Totalspin =  Totalspin * doublenum                
                Msg1 = "遊戲次數: " + str(spincount) + " 目前spin運作正常"
                callbot(Msg1) 
            if poco("TextState").attr("text") == "start":
                poco("ButtonStart").click()
                spincount += 1                
                TotalBet += BetAmount                
            while poco("TextState").attr("text") == "spin":
                    time.sleep(1)            
            Amount = ReelResult(0, BetAmount) 
            SpinTotalAmount += Amount
            print("SpinTotalAmount = %f"%SpinTotalAmount)
except:            
    Msg1 = "東海之王異常結束,原因為: " +  str(sys.exc_info()[0]) + " ,結束畫面如下:"
    callbot(Msg1)
    screenshot()
else:
    callbot("東海之王測試完畢")