# -*- encoding=utf8 -*-
__author__ = "user16"

from airtest.core.api import *
from poco.drivers.unity3d import UnityPoco
poco = UnityPoco()

auto_setup(__file__)
import time

count = 0
redcount = 0
blackcount = 0
luckcount = 0

while(count < 11):
    count = int(count) + 1
    poco('TextBetTimer').wait(60).exists()
    usermoney = poco('TextCoin').attr('text')
    times = 0
    times_2 = 0
    totaltimes = 0
    black = 0
    red = 0
    lucky = 0
    betblack = 0
    betred = 0
    betlucky = 0
    end = 0
    
    poco('ButtonBet1').click()  #押紅
    betred = 1               #押紅flag
    #poco('ButtonBet2').click() #押黑
    #betblack = 1              #黑flag
    poco('ButtonBet3').click() #押好運連連
    betlucky = 1                 #押好運連連flag
    
        
    poco('Canvas').offspring('PokerGrop1').child('Result').wait(60).exists()
    while (int(end) == 0): #判斷哪一方贏
        if poco('ButtonBet1').child('blink').exists():
            red = 1      
            end = 1
        elif poco('ButtonBet2').child('blink').exists():
            black = 1
            end = 1
        else:
            end = 0
    if poco('ButtonBet3').child('blink').exists():
        lucky = 1
    cardtypeIcon = poco('Canvas').offspring('CardType').offspring('Icon')
    cardlist = len(cardtypeIcon)
    cardtype = cardtypeIcon[int(cardlist)-1].attr('texture')
    if cardtype[3] == '0':
        times = 10
    elif cardtype[3] == '1':
        times = 10
    elif cardtype[3] == '2':
        times = 10
    elif cardtype[3] == '3':
        times = 10
    elif cardtype[3] == '4':
        times = cardtype[5]
    else:
        times = 1
    if betlucky == 1:
        if cardtype[3] == '0':
            times_2 = 20
        elif cardtype[3] == '1':
            times_2 = 15
        elif cardtype[3] == '2':
            times_2 = 10
        elif cardtype[3] == '3' and lucky == 1:
            times_2 = 2
        else:
            times_2 = -1
    if betblack == 1:
        if red == 1:
            times = 0-int(times)
    if betred == 1:
        if black == 1:
            times = 0-int(times)
    totaltimes =int(times)+int(times_2)
    if betblack == 1 or betred == 1:
        if int(totaltimes) > 0 :
            usermoney = round(((float(usermoney))+float(totaltimes)*0.95),2)
        else:
            usermoney = round(float(usermoney) + float(totaltimes),2)
    poco('TextBetTimer').wait(60).exists()        
    if float(poco('TextCoin').attr('text')) == float(usermoney): #判斷上局玩家金額是否正確
        print('正確')
    else:
        print('錯誤')
        print (poco('TextCoin').attr('text'))
        print (usermoney)
        print (time.asctime(time.localtime(time.time())))
        break