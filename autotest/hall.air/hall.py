# -*- encoding=utf8 -*-
__author__ = "user16"

from airtest.core.api import *
import requests
from poco.drivers.unity3d import UnityPoco
poco = UnityPoco()

auto_setup(__file__)
lhdgamecount = 0
bjlgamecount = 0 
hhnngamecount = 0
hhdzgamecount = 0
import time

poco('GameLogo_11').click()
while(int(lhdgamecount) < 3):
    betlon = 0
    bethu = 0
    bethan = 0    
    poco('Canvas').offspring('timeBet').child('timeBet').wait(60).exists()
    usermoney = poco('coin').attr('text')
    print(usermoney)
    poco('ButtonBet1').click()  #押龍
    betlon = 1                  #押龍flag
    #poco('ButtonBet3').click() #押虎
    #bethu = 1                  #押虎flag
    #poco('ButtonBet2').click() #押和
    #bethan = 1                 #押和flag
    poco('Canvas').offspring('timeRst').wait(60).exists()
    poco('Canvas').offspring('timeBet').child('timeBet').wait(60).exists()
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
        string =str('龍虎鬥正常運行')
        r = requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=%s" % string)
    else:
        print('錯誤')        
        print (poco('textMoney').attr('text'))
        print (usermoney)
        print (time.asctime(time.localtime(time.time())))
        string =str('龍虎鬥金錢計算錯誤')
        r = requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=%s" % string)
        break   
    lhdgamecount = int(lhdgamecount + 1)
poco("_show").click()
poco("_exit").click()

poco('GameLogo_13').click()
while(int(bjlgamecount) < 3):
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
    usermoney = poco('coin').attr('text')
    print('初始金額:',usermoney)
    sleep(2)
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
    poco('time_open').wait(60).exists()
    sleep(3)
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
        print(bankcard_1[11],bankcard_1[12])
    if playercard_1[11]+playercard_1[12] == playercard_2[11]+playercard_2[12]: #判斷閒家是否有對子
        playerpair = 1
        print(playercard_1[11],playercard_1[12])
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
    else:
        string =str('百家樂金錢計算錯誤')
        r = requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=%s" % string)
        print('錯誤')
        print (poco('coin').attr('text'))
        print (usermoney)
        print (time.asctime(time.localtime(time.time())))
        print('開閒',player)
        print('開莊',banker)
        print('開和',han)
        print('開莊對',bankerpair)
        print('開閒對',playerpair)
        break   
    bjlgamecount = int(bjlgamecount) + 1
poco("_show").click()
poco("_exit").click()

poco('_RightPage').click()
poco('GameLogo_21').click()
while( int(hhnngamecount) < 3):
    poco('TextBetTimer').wait(60).exists()
    usermoney = poco('TextCoin').attr('text')
    print('初始金額:',usermoney)
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
    print("甚麼牛(0:五小,1:四炸,2:五花,3:牛牛,4:牛?,5:沒牛):",cardtype[3])
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
    print("賠率:",times)
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
    print("好運連連賠率:",times_2)
    if betblack == 1:
        if red == 1:
            times = 0-int(times)
    if betred == 1:
        if black == 1:
            times = 0-int(times)
    totaltimes =int(times)+int(times_2)
    print("總賠率:",totaltimes)
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
    hhnngamecount = int(hhnngamecount) + 1
poco("_show").click()
poco("_exit").click()

poco('GameLogo_20').click()
while( int(hhdzgamecount) < 3):
    poco('timeDown').wait(60).exists()
    usermoney = poco('coin').attr('text')
    print('初始金額:',usermoney)
    times = 0
    times_2 = 0
    black = 0
    red = 0
    lucky = 0
    betblack = 0
    betred = 0
    betlucky = 0
    totaltimes = 0
    end = 0
    
    poco('ButtonBet1').click()  #押紅
    betred = 1               #押紅flag
    #poco('ButtonBet2').click() #押黑
    #betblack = 1              #黑flag
    poco('ButtonBet3').click() #押致命一擊flag
    betlucky = 1                 #押致命一擊flag
    

    poco('timeRst').wait(60).exists()
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
    print("牌型種類:單:1,對:2,9-A對:3,順:4,金:5,順金:6,豹子:7",cardtype[16])
    if betred == 1:
        if red == 1:
            times = 1
        else:
            times = -1
    if betblack == 1:
        if black == 1:
            times = 1
        else:
            times = -1        
    if betlucky == 1:
        if cardtype[16] == '3':
            times_2 = 1
        elif cardtype[16] == '4':
            times_2 = 2
        elif cardtype[16] == '5':
            times_2 = 3
        elif cardtype[16] == '6':
            times_2 = 5
        elif cardtype[16] == '7':
            times_2 = 10
        else:
            times_2 = -1
    print("致命一擊賠率:",times_2)
    totaltimes = int(times)+int(times_2)
    print("總輸贏:",totaltimes)
    if int(totaltimes) > 0 :
        usermoney = round(((float(usermoney))+float(totaltimes)*0.95),2)
    else:
        usermoney = round(float(usermoney) + float(totaltimes),2)
    poco('timeDown').wait(60).exists()
    if float(poco('coin').attr('text')) == float(usermoney): #判斷上局玩家金額是否正確
        print('正確')
    else:
        print('錯誤')
        print (poco('coin').attr('text'))
        print (usermoney)
        print (time.asctime(time.localtime(time.time())))
        break
    hhdzgamecount = int(hhdzgamecount) + 1
poco("_show").click()
poco("_exit").click()

    
