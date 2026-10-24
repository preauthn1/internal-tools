# -*- encoding=utf8 -*-
__author__ = "user16"

from airtest.core.api import *
from airtest.core.api import using
from poco.drivers.unity3d import UnityPoco
poco = UnityPoco()

# start your script here
import time
import decimal
import random


bolog = 0
pairlog = 0
lucklog = 0
count = 0
betpaircount = 0
betluckycount = 0
betpairwinmoney = 0
betluckwinmoney = 0
bowshencount = 0
bowshenwinmoney = 0
betzhuangboom = 0
betzhuangboomcount = 0
betzhuangboomwinmoney = 0
userbet = 0
betpairmoney = 0
fo = open("pairtest.log", "w")
while (count < 3):
    
    count = int(count) + 1
    sleep(1)
    if poco('btn').exists():
        usermoney = round(float(poco('textMoney').attr('text'))+1,2)
    else:
        usermoney = poco('textMoney').attr('text')
    poco('Canvas').offspring('CastPool').offspring('bet').wait(60).exists() #判斷是否已下注
    bowshenbet = 0
    betpair = 1
    usersub = 0
    betlucky = 0
    betzhuangboom = 0
    luckyflag = 0   
    if int(betpairmoney) >= 1: #判斷是否要下對子
        betpairmoney = int(betpairmoney)*2
        betmoney = betpairmoney
        if int(betmoney) > 1000:
            betmoney = 1000
        if int(betmoney) >= 500:
            n=int(betmoney)//500
            poco('Canvas').offspring('Coin').offspring("PfChip(Clone)")[5].child('btn').click()
            betmoney = int(betmoney)-500*int(n)
            while(int(n) > 0):               
                poco('BetArea_1').click()
                n = int(n)- 1
            
        if int(betmoney) >= 100:
            n=int(betmoney)//100
            poco('Canvas').offspring('Coin').offspring("PfChip(Clone)")[4].child('btn').click()
            betmoney = int(betmoney)-100*int(n)
            while(int(n) > 0):               
                poco('BetArea_1').click()
                n = int(n)- 1
                
        if int(betmoney) >= 50:
            n=int(betmoney)//50
            poco('Canvas').offspring('Coin').offspring("PfChip(Clone)")[3].child('btn').click()
            betmoney = int(betmoney)-50*int(n)
            while(int(n) > 0):               
                poco('BetArea_1').click()
                n = int(n)- 1
                
        if int(betmoney) >= 10:
            n=int(betmoney)//10
            poco('Canvas').offspring('Coin').offspring("PfChip(Clone)")[2].child('btn').click()
            betmoney = int(betmoney)-10*int(n)
            while(int(n) > 0):               
                poco('BetArea_1').click()
                n = int(n)- 1
                
        if int(betmoney) >= 5:
            n=int(betmoney)//5
            poco('Canvas').offspring('Coin').offspring("PfChip(Clone)")[1].child('btn').click()
            betmoney = int(betmoney)-5*int(n)
            while(int(n) > 0):               
                poco('BetArea_1').click()
                n = int(n)- 1 
            sleep(1)
            
        if int(betmoney) >= 1:
            poco('Canvas').offspring('Coin').offspring("PfChip(Clone)")[0].child('btn').click()
            n=int(betmoney)//1    
            betmoney = int(betmoney)-1*int(n)
            while(int(n) > 0):               
                poco('BetArea_1').click()
                n = int(n)- 1  
                
    else:
        poco("PfChip(Clone)")[0].child('btn').click()
        poco('BetArea_1').click()
        betpairmoney = 1
        
        betpaircount = int(betpaircount) + 1  
    poco('Canvas').offspring('player1').offspring('Text').wait(60).exists()   
    userpoint = poco('Canvas').offspring('player1').offspring('Text').attr('text')  #取得玩家手牌
    userpoker = poco('Canvas').offspring('player1').offspring('poka(Clone)')
    userpoker_1 = userpoker[0].attr('texture')
    userpoker_2 = userpoker[1].attr('texture')  
    if len(userpoker[0].attr('texture')) > 16:
        userpokerpoint_1 = userpoker_1[15]+userpoker_1[16]
    else:           
        userpokerpoint_1 = userpoker_1[15]
    if len(userpoker[1].attr('texture')) > 16:
        userpokerpoint_2 = userpoker_2[15]+userpoker_2[16]
    else:           
        userpokerpoint_2 = userpoker_2[15]
    if betpair == 1:  #對子押注點勝負判斷
        if userpokerpoint_1 == userpokerpoint_2:              
            betpairwinmoney = int(betpairwinmoney) + round(float(betpairmoney)*0.95,2)   
            betpairmoney = 0
        else:
            betpairwinmoney = int(betpairwinmoney) - int(betpairmoney)
    Tmpstr = "第" + str(count) + "局" + "對子總輸贏:" + str(betpairwinmoney) + "\n"
    fo.write(Tmpstr)
    
    while userpoint != "BJ":   #判斷牌型不是bj且小於x進行要牌       
        poco('operbtn_2').wait(60).exists()  
        #poco('operbtn_2').click()
        sleep(1)       
        if poco('Canvas').offspring('BetSub_1').exists():
            usersub = 1
            userpoint = int(poco('Canvas').offspring('player1').offspring('Point').offspring('Text').attr('text'))
        while int(userpoint) < 16:        
            poco('operbtn_4').click()
            sleep(1)
            userpoint = int(poco('Canvas').offspring('player1').offspring('Point').offspring('Text').attr('text'))
        if  int(userpoint) >= 15 and int(userpoint) < 21: 
            poco('operbtn_1').click()
        if usersub == 1:
            usersubpoint = int(poco('Canvas').offspring('player1').offspring('PointSub').offspring('Text').attr('text'))
            while int(usersubpoint) < 16:        
                poco('operbtn_4').click()
                sleep(1)
                usersubpoint = int(poco('Canvas').offspring('player1').offspring('PointSub').offspring('Text').attr('text'))
            if  int(usersubpoint) >= 15 and int(usersubpoint) < 21: 
                    poco('operbtn_1').click()              
        poco('Canvas').offspring('Banker').offspring('Text').wait(60).exists()
    poco('Canvas').offspring('ready').wait(180)
    poco('Canvas').offspring('ready').click()
    sleep(1)
fo.close()