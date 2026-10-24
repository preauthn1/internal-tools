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
while (count < 1000):
    count = int(count) + 1
    sleep(1)
    if poco('btn').exists():
        usermoney = round(float(poco('textMoney').attr('text'))+1,2)
    else:
        usermoney = poco('textMoney').attr('text')
    poco('Canvas').offspring('CastPool').offspring('bet').wait(60).exists() #判斷是否已下注
    bowshenbet = 0
    usersub = 0
    betpair = 0
    betlucky = 0
    betzhuangboom = 0
    luckyflag = 0
    betpairProb=random.randint(1,100)
    betluckyProb =random.randint(1,100)
    betzhuangboomProb =random.randint(1,100)
    #if int(betzhuangboomProb) < 50:   
    #    poco('BetArea_2').click()
    #    betzhuangboom = 1
    #    betzhuangboomcount = int(betzhuangboomcount) + 1
    #if int(betluckyProb) < 50:   #判斷是否要下幸運20
    #    poco('Canvas').offspring('BetArea_3').click()
    #    betlucky = 1
    #    betluckycount = int(betluckycount) + 1
    #if int(betpairProb) < 50:    #判斷是否要下對子
    #    poco('BetArea_1').click()
    #    betpair = 1
    #    betpaircount = int(betpaircount) + 1
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
            usermoney = round(((float(usermoney)*100)+(10.45*100))/100,2)
            betpairwinmoney = round(float(betpairwinmoney)+10.45,2)            
        else:
            usermoney = round(float(usermoney) - 1,2)
    pairlog = usermoney        
    while userpoint != "BJ":   #判斷牌型不是bj且小於x進行要牌       
        #if poco('msgbg').wait(1).exists():   #判斷是否需要保險
        #    bowshen = random.randint(1,100)
        #    if int(bowshen) < 50:
        #        poco('savebtn_3').click()
        #    else:
        #        poco('savebtn_1').click()
        #        bowshenbet = 1
        #        bowshencount = int(bowshencount) + 1
        if betlucky == 1: #幸運押注點勝負判斷
            if int(userpokerpoint_1) == 14:
                userpokerpoint_1 = 1
            elif int(userpokerpoint_1) > 10:
                userpokerpoint_1 = 10
            if int(userpokerpoint_1) == 1:
                userpokerpoint_1 = 11
            if int(userpokerpoint_2) == 14:
                userpokerpoint_2= 1
            elif int(userpokerpoint_2) > 10:
                userpokerpoint_2 = 10
            if int(userpokerpoint_2) == 1:
                userpokerpoint_2 = 11
            if int(userpokerpoint_1)+int(userpokerpoint_2)==20:
                luckyflag = 4.75                
                if int(userpokerpoint_1) > 9 and int(userpokerpoint_2) > 9:                  
                    if int(userpoker_1[15]+userpoker_1[16]) == int(userpoker_2[15]+userpoker_2[16]) :        
                        luckyflag = 9.5
                        if userpoker_1[13] == userpoker_2[13] and int(userpoker_1[15]+userpoker_1[16])!=10:
                            luckyflag= 23.75
            if float(luckyflag) != 0:
                usermoney = round((float(usermoney)*100+(float(luckyflag)*100))/100,2)
                betluckwinmoney=round((float(betluckwinmoney)+float(luckyflag)*100)/100,2)
            else:
                usermoney = round(float(usermoney) - 1,2)                
            betlucky = 0
            lucklog = usermoney
        poco('operbtn_2').wait(60).exists()  
        #poco('operbtn_2').click()
        #sleep(1)       
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
        bankpoint = int(poco('Canvas').offspring('Banker').offspring('Text').attr('text'))
        bankpoker = poco('Canvas').offspring('Banker').offspring('poka(Clone)')
        if bowshenbet == 1: #有保險後進入保險金額判斷
            if bankpoint == 21 and len(bankpoker) == 2:
                usermoney = round(((float(usermoney)*100)+(1*100))/100,2)
                bowshenwinmoney = float(bowshenwinmoney) + 1
            else:
                usermoney = round(float(usermoney) - 0.5,2)
        bolog = usermoney
        if int(userpoint) < 22:
            while int(bankpoint) < 17:
                sleep(1)
                bankpoint = int(poco('Canvas').offspring('Banker').offspring('Text').attr('text'))
        if betzhuangboom == 1:  #莊爆判斷    
            if int(bankpoint) > 21:
                bankpoker = poco('Canvas').offspring('Banker').offspring('poka(Clone)')
                boompokercount = len(bankpoker)
                if boompokercount == 5:
                    usermoney = round(((float(usermoney)*100)+(2.85*100))/100,2)
                    betzhuangboomwinmoney =round((float(betzhuangboomwinmoney)+2.85),2)
                elif boompokercount == 6:
                    usermoney = round(((float(usermoney)*100)+(14.25*100))/100,2)
                    betzhuangboomwinmoney =round((float(betzhuangboomwinmoney)+14.25),2)
                elif boompokercount == 7:
                    usermoney = round(((float(usermoney)*100)+(47.5*100))/100,2)
                    betzhuangboomwinmoney =round((float(betzhuangboomwinmoney)+47.5),2)
                elif boompokercount == 8:
                    usermoney = round(((float(usermoney)*100)+(237.5*100))/100,2)
                    betzhuangboomwinmoney =round((float(betzhuangboomwinmoney)+237.5),2)
                elif boompokercount == 9:
                    usermoney = round(((float(usermoney)*100)+(475*100))/100,2)
                    betzhuangboomwinmoney =round((float(betzhuangboomwinmoney)+475),2)
                else:
                    usermoney = round(float(usermoney) - 1,2)
            else:
                usermoney = round(float(usermoney) - 1,2)
        if int(userpoint) > 21 or (int(bankpoint) < 22 and int(bankpoint) > int(userpoint)):
            usermoney = round(float(usermoney) - 1,2)
        elif (int(userpoint) < 22 and int(userpoint) > int(bankpoint)) or (int(bankpoint) > 21 and int(userpoint) < 22):
            usermoney = round(((float(usermoney)*100)+(0.95*100))/100,2)
        if usersub == 1:
            if int(usersubpoint) > 21 or (int(bankpoint) < 22 and int(bankpoint) > int(usersubpoint)):
                usermoney = round(float(usermoney) - 1,2)
            elif (int(usersubpoint) < 22 and int(usersubpoint) > int(bankpoint)) or (int(bankpoint) > 21 and int(usersubpoint) < 22):
                usermoney = round(((float(usermoney)*100)+(0.95*100))/100,2)
            if usersubpoint == 'BJ':
                usermoney = round(float(usermoney) + 1.43,2)
        break                  
    if userpoint == 'BJ': #BJ直接結算
        if betlucky == 1:
            usermoney = round(float(usermoney) - 1,2)
        if betzhuangboom == 1:
            usermoney = round(float(usermoney) - 1,2)
        usermoney = round(float(usermoney) + 1.43,2)
    poco('Canvas').offspring('ready').wait(180)
    poco('Canvas').offspring('ready').click()  
    sleep(1)
  #  if float(poco('textMoney').attr('text')) == float(usermoney): #判斷上局玩家金額是否正確
  #      print('正確')               
  #  elif poco('btn').exists():
  #      if round(float(poco('textMoney').attr('text'))+1,2) == float(usermoney): #判斷上局玩家金額是否正確
  #          print('正確')
  #      else:
  #          print('錯誤')
  #          print('判斷已開始下注')
  #          print (poco('textMoney').attr('text'))
  #          print (usermoney)
  #          print (time.asctime(time.localtime(time.time())))
  #          break    
  #  else:
  #      print('錯誤')
  #      print (poco('textMoney').attr('text'))
  #      print (usermoney)
  #      print (time.asctime(time.localtime(time.time())))
  #      break   
    