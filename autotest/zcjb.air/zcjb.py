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

def timediff():
    runtime = time.time()    
    return(int((runtime - starttime) // 60))  
    
def SymbolResult(checkReel):
    SList = []
    y = poco("UISlot").child(checkReel).child("Symbol_1").child("RootPos").attr("pos") 
    y1 = int(y[1] * 10) 
    y2 = int(y1 // 2)    
    for i in range(0, 3):
        k = y2 * 3 + i
        symbol = "Symbol_" + str(SymbolList[k])  
        s = poco("UISlot").child(checkReel).child(symbol).children().attr("name")   
        if len(s) > 30 :
            s1 = "symbol_"+ s[20:23]
        else:
            s1 = "symbol_"+ s[20:22]
        SList.append(s1)  
    return(SList)
    #y = poco("UISlot").child(checkReel).children()   
    #for i in range(0, 5):
    #    symbol = str(y[i].attr('name'))
    #    if poco("UISlot").child(checkReel).child(symbol).child("ImageIcon").exists() == False:
    #        s = poco("UISlot").child(checkReel).child(symbol).children().attr("name")
    #        s1 = "symbol_"+ s[20:22]
    #        print(s1)
    #        SList.append(s1)
    #    else:
    #        print('noting')       
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
    elif sname == "symbol_p12":
        return(symbol_p12[num])
    elif sname == "symbol_p13":
        return(symbol_p13[num])
    elif sname == "symbol_p14":
        return(symbol_p14[num])
    elif sname == "symbol_s1":
        if BetAmount == 8:
            return(symbol_s1[num])
        else:
            return(symbol_s1[num] * 11)    
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
        if (RList2.count(i) != 0) or (RList2.count("symbol_w1") != 0):
            R2 = R2 + RList2.count(i) + RList2.count("symbol_w1")
            if (RList3.count(i) != 0) or (RList3.count("symbol_w1") != 0):
                linknum = 0
                R3 = R3 + RList3.count(i) + RList3.count("symbol_w1")
                if (RList4.count(i) != 0) or (RList4.count("symbol_w1") != 0):
                    linknum = 1
                    R4 = R4 + RList4.count(i) + RList4.count("symbol_w1")
                    if(RList5.count(i) != 0):
                        linknum = 2
                        R5 = R5 + RList5.count(i)        
        if R5 > 0:            
            RewardMoney = RewardMoney + GetBaseReward(i, linknum, BetAmount) * R2 * R3 * R4 * R5 
        elif R4 > 0:            
            RewardMoney = RewardMoney + GetBaseReward(i, linknum, BetAmount) * R2 * R3 * R4
        elif R3 > 0:            
            RewardMoney = RewardMoney + GetBaseReward(i, linknum, BetAmount) * R2 * R3   
        if (i == "symbol_s1") and (R3 > 0):
            Msg1 = str(linknum + 3) + "個銅鑼連線,中奬金額" + str(RewardMoney)
            callbot(Msg1)
    
    return(RewardMoney)     

def RewardDifferent(RList1, RList2, RList3, RList4, RList5, TotalReward, gresult, FAmount):
    print("RewardDifferent")
    screenshot()
    Msg1 = RList1[0] + ", " + RList2[0] + ", "  + RList3[0] + ", " + RList4[0] + ", " + RList5[0] + "\n"
    Msg1 += RList1[1] + ", " + RList2[1] + ", " + RList3[1] + ", " + RList4[1] + ", " + RList5[1] + "\n"
    Msg1 += RList1[2] + ", " + RList2[2] + ", " + RList3[2] + ", " + RList4[2] + ", " + RList5[2] + "\n"
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
    #TPay = poco("Result").child("ImageNormal").child("TextPays").attr("text")
    #if TPay != "":        
    #    k = TPay.index(" ", 0, len(TPay)) + 1;        
    #    pay1 = TPay[k:-1]
    gresult = poco("User").child("TextStepCoin").attr("text")
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
    poco.click([0.7578, 0.9522])
    PAmount = round(float(poco("User").child("TextCoin").attr("text")), 2)
    if UserAmount1 != PAmount:
        screenshot()
        Msg1 = "玩家金額不正確;玩家初始金額: " + str(UserAmount) 
        Msg1 += " 總贏金: " + str(SpinTotalAmount) 
        Msg1 += " 總下注: " + str(TotalBet) + " (遊戲次數: " + str(spincount) + ")"
        Msg1 += " 預期金額: " + str(UserAmount1) + " 目前結算金額: " + str(PAmount)
        callbot(Msg1) 
        return(-1)
    else:
        return(PAmount)

def ProcessJP():
    print("ProcessJP")
    JPAmount = 0    
    JPList = GetJPAmount() 
    BonusPicList = []
    for i in range(0, 4):
        for k in range(0, 3):
            poco.click([0.23 + 0.17 * i, 0.24 + 0.25 * k])
            BonusPicList = CollectBonusPic(i + 1 , k + 1, BonusPicList)
            if poco("UIBonusGame").child("ImageMoney").exists():
                break
        if poco("UIBonusGame").child("ImageMoney").exists():
            break
    for i in range(1,5):
        pic = "Image" + str(i)       
        if BonusPicList.count(pic) == 3:
            index = i - 1 
            JPAmount = round(float(JPList[index]), 2)
            break
    return(JPAmount)

try:
    betstr=poco("ButtonChip").child("Text").attr("text")
    while betstr != "X1":
        poco("ButtonChipUp").click()
        betstr=poco("ButtonChip").child("Text").attr("text")    
    isfirstspin = True       
    for i in range(0, 2):                          
        UserAmount = round(float(poco("User").child("TextCoin").attr("text")), 2) 
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
            BetAmount = 8
            poco("UIGoldSymbol").offspring("Toggle1").child("Checkmark").click()                       
            Msg1 = "招財進寶測試開始,游戲倍數1倍, 每局押注金額" + str(BetAmount) + "元,玩家起始金額" + str(UserAmount)            
            callbot(Msg1)                       
        else:            
            BetAmount = 88
            poco("UIGoldSymbol").offspring("Toggle5").child("Checkmark").click()                       
            Msg1 = "招財進寶測試開始,游戲倍數5倍, 每局押注金額" + str(BetAmount) + "元,玩家起始金額" + str(UserAmount)            
            callbot(Msg1)                        
        while (timediff() < 300):
            tt = spincount % Totalspin
            if (spincount > 0) and (tt == 0):
                Totalspin =  Totalspin * doublenum
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
                    time.sleep(1)            
            Amount = ReelResult(0, BetAmount) 
            SpinTotalAmount += Amount
            print("SpinTotalAmount = %f"%SpinTotalAmount)
            if poco("UIBonusGame").exists():                    
                    JPAmount = ProcessJP()    
                    SpinTotalAmount += JPAmount
                    SpinTotalAmount = CompareAmount(UserAmount, TotalBet, spincount, SpinTotalAmount)            
            while poco("UIBigWin").exists():
                poco("TextStepCoin").click()
                poco("TextCoin").click()
                time.sleep(0.1) 
                poco("TextCoin").click()
            while poco("UIFreeSpin").exists():
                if isneedclick:
                    poco("UIFreeSpin").child("ButtonStart").click() 
                    isneedclick = False
                    isFreespin = True   
                if poco("UIBonusGame").exists():
                    JPAmount = ProcessJP(JPAmount)  
                    SpinTotalAmount += JPAmount
                while poco("UIBigWin").exists():
                    poco("TextCoin").click()
                    time.sleep(0.1) 
                    poco("TextCoin").click()
            if isFreespin:
                while poco("UIBigWin").exists():
                    poco("TextStepCoin").click()
                    poco("TextCoin").click()
                    time.sleep(0.1)
                    poco("TextCoin").click()
                isFreespin = False
                isneedclick = True
                while poco("TextState").attr("text") == "spin":
                    time.sleep(0.1)
                FAmount = poco("User").child("TextStepCoin").attr("text")
                SpinTotalAmount += round(float(FAmount), 2)                
                FAmount = 0        
except:
    Msg1 = "招財進寶異常結束,原因為: " +  str(sys.exc_info()[0]) + " ,結束畫面如下:"
    callbot(Msg1)
    screenshot()
else:
    callbot("招財進寶測試完畢")