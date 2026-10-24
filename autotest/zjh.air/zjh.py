# -*- encoding=utf8 -*-
__author__ = "user16"

from airtest.core.api import *

from poco.drivers.unity3d import UnityPoco
poco = UnityPoco()

auto_setup(__file__)

count = 0

while(int(count) < 10000): #控制次數
    coin = 0
    usermoney = float(poco('Canvas').offspring('Player3').child('MoneyNumber').attr('text'))
    poco('AnimCard1').wait(60)

    while int(poco('RoundText').attr('text')) < 3:
        sleep(2)
        if poco('AddNumButton1').exists():
            poco('AddNumButton1').click()
        poco('FollowButton').wait(60).exists()
        poco('FollowButton').click()
        if poco('ChangeButton').exists():
            coin = 1
    while int(coin) == 0:
        if poco('Canvas').offspring('SelectCompare').offspring('PokerStart').exists():
            compare = poco('Canvas').offspring('SelectCompare').offspring('PokerStart')
            compare[0].click() 
        else:
            sleep(1)
        if poco('WinFont').exists():
            a =poco('WinFont').attr('pos')
            print(a[0][3])
            if int(poco('WinFont').attr('pos')) == int(poco('Canvas').offspring('Player3').child('HeroDB').attr('pos')):
                betmoney = float(poco('CalcMoney').attr('text'))
                totalmoney = round(usermoney + betmoney,2)                                            
                coin = 1
        if poco('ChangeButton').exists():
            betmoney = float(poco('Canvas').offspring('Player3').child('BetMoneyNumber').attr('text'))
            totalmoney = round(usermoney - betmoney,2)            
            coin = 1
    #poco('WaiveButton').wait(5).exists()
    #poco('WaiveButton').click()

    poco('ChangeButton').wait(5).exists()
    poco('ChangeButton').click()
    if float(totalmoney) == float(poco('Canvas').offspring('Player3').child('MoneyNumber').attr('text')):
        print('正確')
    else:
        print('錯誤')
