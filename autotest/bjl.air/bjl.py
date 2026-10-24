# -*- encoding=utf8 -*-
__author__ = "user16"

from airtest.core.api import *

from poco.drivers.unity3d import UnityPoco
poco = UnityPoco()

auto_setup(__file__)
import time
import os

count = 0
while (int(count) < 10): #控制次數
    betbanker = 0
    betplayer = 0
    bethan = 0
    betbankerpair = 0
    betplayerpair = 0
    player = 0
    banker = 0
    han =0
    bankerpair = 0
    playerpair = 0
    end = 0
    money = 0
    poco('time').wait(60).exists()
    count = int(count) + 1
    usermoney = poco('coin').attr('text')
    poco('ButtonBet1').click()  #押莊
    betbanker = 1                  #押莊flag
    #poco('ButtonBet2').click() #押閒
    #betplayer = 1                  #押閒flag
    #poco('ButtonBet3').click() #押和
    #bethan = 1                 #押和flag
    #poco('ButtonBet6').click() #押莊對
    #betbankerpair = 1                  #押莊對flag
    #poco('ButtonBet7').click() #押閒對
    #betplayerpair = 1                 #押閒對flag
    while (int(end) == 0): #判斷哪一方贏
        if poco('win1').exists():
            player = 1
            end = 1
        elif poco('win2').exists():
            banker = 1
            end = 1
        elif poco('win3').exists():
            han = 1            
            end = 1
        else:
            end = 0
    bankcard_1 = poco('Canvas').offspring('banker').offspring('card1').attr('texture')
    bankcard_2 = poco('Canvas').offspring('banker').offspring('card2').attr('texture')
    playercard_1 = poco('Canvas').offspring('user').offspring('card1').attr('texture')
    playercard_2 = poco('Canvas').offspring('user').offspring('card2').attr('texture')
    if bankcard_1[11]+bankcard_1[12] == bankcard_2[11]+bankcard_2[12]: #判斷莊家是否有對子
        bankerpair = 1
    if playercard_1[11]+playercard_1[12] == playercard_2[11]+playercard_2[12]: #判斷閒家是否有對子
        playerpair = 1
    if betbanker == 1:
        if banker == 1:
            money = int(money) + 1           
        elif player == 1:
            money = int(money) - 1  
    if betplayer == 1:
        if player == 1:
            money = int(money) + 1 
        elif banker == 1:
            money = int(money) - 1 
    if bethan == 1:
        if han == 1:
            money = int(money) + 8
        else:
            money = int(money) - 1 
    if betbankerpair == 1:      
        if bankerpair == 1:
            money = int(money) + 11
        else:
            money = int(money) - 1
    if betplayerpair == 1:      
        if playerpair == 1:
            money = int(money) + 11
        else:
            money = int(money) - 1     
    if int(money) > 0 :
        usermoney = round(((float(usermoney)*100)+float(money)*0.95*100)/100,2)
    else:
        usermoney = round(float(usermoney) + float(money),2)
    sleep(1)
    if float(poco('coin').attr('text')) == float(usermoney): #判斷上局玩家金額是否正確        
        print('正確')
        break
    else:
        print('錯誤')
        print (poco('coin').attr('text'))
        print (usermoney)
        print (time.asctime(time.localtime(time.time())))
        break   