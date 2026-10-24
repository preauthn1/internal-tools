# -*- encoding=utf8 -*-
__author__ = "user16"

from airtest.core.api import *
from poco.drivers.unity3d import UnityPoco
poco = UnityPoco()

auto_setup(__file__)
import time

count = 0
while (count < 11):
    count = int(count) + 1
    betlon = 0
    bethu = 0
    bethan = 0    
    poco('timeDown').wait(60).exists()
    usermoney = poco('coin').attr('text')
    poco('ButtonBet1').click()  #押龍
    betlon = 1                  #押龍flag
    #poco('ButtonBet3').click() #押虎
    #bethu = 1                  #押虎flag
    #poco('ButtonBet2').click() #押和
    #bethan = 1                 #押和flag
    while (int(poco('timeDown').attr('text')) > 2):
           sleep(1)    
    poco('timeDown').wait(60).exists()
    winner = poco('Canvas').offspring('content').offspring('_icon(Clone)')
    if betlon == 1: #判斷龍贏
        if winner[(len(winner)-1)].offspring('point1').exists():
            usermoney = round(((float(usermoney)*100)+(0.95*100))/100,2)
        elif winner[(len(winner)-1)].offspring('point3').exists():
            usermoney = round(float(usermoney) - 1,2)
    if bethu == 1: #判斷虎贏
        if winner[(len(winner)-1)].offspring('point3').exists():
            usermoney = round(((float(usermoney)*100)+(0.95*100))/100,2)
        elif winner[(len(winner)-1)].offspring('point1').exists():
            usermoney = round(float(usermoney) - 1,2)
    if bethan == 1:
        if winner[(len(winner)-1)].offspring('point2').exists():
            usermoney = round(((float(usermoney)*100)+(7.6*100))/100,2)
        else:
            usermoney = round(float(usermoney) - 1,2)
    if float(poco('coin').attr('text')) == float(usermoney): #判斷上局玩家金額是否正確
        print('正確')
    else:
        print('錯誤')
        print (poco('textMoney').attr('text'))
        print (usermoney)
        print (time.asctime(time.localtime(time.time())))
        break   