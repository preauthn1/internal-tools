# -*- encoding=utf8 -*-
__author__ = "Frank"

from airtest.core.api import *


from poco.drivers.unity3d import UnityPoco
poco = UnityPoco()
from PIL import Image
from PIL import ImageGrab
import time
import os
import sys
import requests
auto_setup(__file__)


def callbot(Msg):
    r=requests.get("https://example.invalid/bot965536905:AAFSSpfGRYUzsm2548QHNroViU06Rz6YQNI/sendMessage?chat_id=-332128122&text=%s"%Msg)
    return()

def screenshot():
    saveDir= "D:\\QA\\QA_AutoTest\\temp"     
    t=time.localtime()
    filename = time.strftime("%Y-%m-%d_%H-%M-%S", t) + '.jpg'
    img =ImageGrab.grab()
    img.save(os.path.join(saveDir, filename), 'jpeg')
    files = {'photo': open(os.path.join(saveDir, filename), 'rb')}
    url = "https://example.invalid/bot965536905:AAFSSpfGRYUzsm2548QHNroViU06Rz6YQNI/sendPhoto"
    data = {'chat_id' : "-332128122"}   
    r= requests.post(url, files=files, data=data)
    
def pokercheck(pokerflowerlist,pokerlist): #判斷牌型
    fure = 0
    three = 0
    pair = 0
    same = 0
    samecount = 0
    sameflower = 0
    samestraight = 0
    samebigstraight = 0
    straight = 0
    single = 0
    pokerlist.sort() 
    for j  in range(0 , 4):
        s = str(j)
        sameflower = pokerflowerlist.count(s)
        if  sameflower == 5 :
            same = 1
    p1 =int(pokerlist[0])
    p2 =int(pokerlist[1])-1
    p3 =int(pokerlist[2])-2
    p4 =int(pokerlist[3])-3
    p5 =int(pokerlist[4])-4
    if p1 == p2:
        if p2 == p3:
            if p3 == p4:
                if p4 == p5:
                    straight = 1 
    if p1 == 1 and p2 == 9:
        if p2 == p3:
            if p3 == p4:
                if p4 == p5:                   
                    straight = 1
    if p1 == 1 and p2 == 9 and straight == 1 and same == 1:
        samebigstraight = 1
    for i  in range(1, 14):
        if i < 10:
            s =  '0' + str(i)
        else:
            s = str(i)
        samecount = pokerlist.count(s)
        if samecount == 2:
            pair += 1
        elif samecount == 3:
            three += 1 
        elif samecount == 4:
            fure += 1
    if p5 == 9 and p1 == 1:
        single = 1
    #牌型判斷，0:不成局，1:高牌，2:對子，3:兩對，4:三條，5:順子，6:同花，7:葫蘆，8:四條，9:同花順，10:皇家同花順
    if samebigstraight == 1 :
        pok = 10
    elif same == 1 and straight == 1:
        pok = 9
    elif fure == 1:
        pok = 8
    elif three == 1 and pair == 1:
        pok = 7
    elif same == 1:
        pok = 6
    elif straight == 1:
        pok = 5
    elif three == 1:
        pok = 4
    elif pair == 2:
        pok = 3
    elif pair == 1:
        pok = 2
    elif single == 1: 
        pok = 1
    else:
        pok = 0
    return(pok)


def cardpoint(Player):
    Plist = []
    pc = 0
    playerpoint = 0
    playercard = poco(Player).child("CardRoot").children()
    for i in range(0,5):
        pc = playercard[i].attr('texture')
        if len(pc) > 15:
            playerpoint = pc[14:16]
        else:
            playerpoint = "0" + pc[14]
        Plist.append(playerpoint)
    return(Plist)

def cardflower(Player):
    Flist = []
    pf = 0
    playerflower = 0
    playercard = poco(Player).child("CardRoot").children()
    for i in range(0,5):
        pf = playercard[i].attr('texture')
        playerflower = pf[12]
        Flist.append(playerflower)
    return(Flist)

def cardtype(Player):
    cdtype = 0
    playertype = poco(Player).child("PatternImage").attr('texture')[5:7]
    if playertype == '07' :
        cdtype = 2
    elif playertype == '14':
        cdtype = 10
    elif playertype == '13':
        cdtype = 9
    elif playertype == '10':
        cdtype = 8
    elif playertype == '22':
        cdtype = 7
    elif playertype == '12':
        cdtype = 6
    elif playertype == '11':
        cdtype = 5
    elif playertype == '09':
        cdtype = 4
    elif playertype == '08':
        cdtype = 3
    elif playertype == '06':
        cdtype = 1
    elif playertype == '15':
        cdtype = 0
    return(cdtype)



def pokerbattle(playerlist,playerflowerlist,banklist,bankflowerlist,pok):
    p1 = int(playerlist[0])
    p3 = int(playerlist[2])
    p4 = int(playerlist[3])
    p5 = int(playerlist[4])
    b1 = int(banklist[0])
    b3 = int(banklist[2])
    b4 = int(banklist[3])
    b5 = int(banklist[4])
    pair = 0
    win = 0
    samecount = 0
    samecount2 = 0
                
    if pok == 1: #判斷輸贏 0=玩家贏 1=莊家贏 2=平手
        if p1 == 1:
            p5 = 14
        if b1 == 1:
            b5 = 14
        if p5 > b5:
            win = 0
        elif p5 < b5:
            win = 1
        else:
            if p5 == b5:
                if p4 > b4:
                    win =0
                elif p4 < b4:
                    win =1
            
    elif pok == 2:
        for i  in range(1, 14):
            if i < 10:
                s =  '0' + str(i)
            else:
                s = str(i)
            samecount = playerlist.count(s)
            samecount2 = banklist.count(s)
            if samecount == 2:
                if i == 1:
                    i = 14
                p5 = i
            if samecount2 == 2:
                if i == 1:
                    i = 14
                b5 = i
        if p5 > b5:
             win = 0
        elif p5 < b5:
            win = 1
        else:
            win = 2
    elif pok == 3:
        if p1 == 1:
            p4 = 14
        if b1 == 1:
            b4 = 14
        if p4 > b4:
            win = 0
        elif p4 < b4:
            win = 1
        else:
            if p2 > b2:
                win = 0
            elif p2 < b2:
                win = 1
            else:
                win = 2
    elif pok == 9 or pok == 5:
        if p1 == 1:
            p5 = 14
        if b1 == 1:
            b5 = 14
        if p5 > b5:
            win = 0
        elif p5 < b5:
            win = 1
        else:
            win = 2
    else:
        if p3 == 1:
            p3 = 14
        if b3 == 1:
            b3 = 14
        if p3 > b3:
            win = 0
        elif p3 < b3:
            win = 1
        else:
            win = 2
    return(win)        

zerocountlog = 0
onecountlog = 0
paircountlog = 0
twopaircountlog = 0
threecountlog = 0
straightcountlog = 0
samecountlog = 0
fullhousecountlog = 0
fourcountlog = 0
samestraightcountlog = 0
royalsamestraightcountlog = 0
win = 0
lose = 0
dues = 0
none = 0
msg = "海盜船測試開始" +"\n"
msg += "攜帶金幣:" + poco('MoneyText').attr('text')
callbot(msg)
gamecount = 0
side = 0 
winjackpot = 0
try: 
    while(True):
        battleresult = 0
        gamecount += 1
        poco('UIChip').offspring('ChipPrefab(Clone)').wait(60).exists()
        poco('BetPoolButton1').click()
        poco('BetPoolButton2').click()
        side = 1
        poco('DoubleButton').wait(60).exists()
        jackpotmoney = float(poco('JackpotMoney').child('Text').attr('text'))
        usermoney = float((poco('User1').offspring('MoneyText').attr('text')))
        usermoney += 3
        poco('DoubleButton').click()
        playerlist = cardpoint('User1')
        playerflowerlist = cardflower('User1')
        usertype = pokercheck(playerflowerlist,playerlist)
        poco('User1').child("PatternImage").wait(60).exists()
        playertype = int(cardtype('User1'))
        if usertype == 0:
            usertype = 1
        if usertype != playertype:
            screenshot()
            Msg1 = "玩家牌型判斷有誤 "
            Msg1 += " 牌型: " + str(playerlist) 
            Msg1 += " 圖片顯示: " + str(playerlist)
            Msg1 += " 計算牌型: " + str(usertype)
            callbot(Msg1)
            break
        bankertype = int(cardtype('UIDealer'))
        bankerlist = cardpoint('UIDealer')
        bankerflowerlist = cardflower('UIDealer')
        btype = pokercheck(bankerflowerlist,bankerlist)             
        if playertype > bankertype:
            battleresult = 13
        elif playertype < bankertype :
            battleresult = 14
        else:
            if playertype == bankertype:
                battlecheck = pokerbattle(playerlist,playerflowerlist,bankerlist,bankerflowerlist,usertype)
                if battlecheck == 0:
                    battleresult = 13
                elif battlecheck == 1:
                    battleresult = 14
                else:
                    battleresult = 15
        if bankertype == 0:
            none += 1
        if side == 1:
            if usertype < 6:
                usermoney -= 1
            else:
                if usertype == 6:
                    usermoney += 49
                    winjackpot += 49
                elif usertype == 7:
                    usermoney += 74
                    winjackpot += 74
                elif usertype == 8:
                    usermoney += 99
                    winjackpot += 99
                elif usertype == 9:
                    usermoney += round((jackpotmoney * 0.1) - 1,2)
                    winjackpot += round((jackpotmoney * 0.1) - 1,2)
                else: 
                    usermoney += round(jackpotmoney - 1,2)
                    winjackpot += round(jackpotmoney - 1,2)
            if bankertype < 6:
                usermoney -= 1
            else:
                if bankertype == 6:
                    usermoney += 49
                    winjackpot += 49
                elif bankertype == 7:
                    usermoney += 74
                    winjackpot += 74
                elif bankertype == 8:
                    usermoney += 99
                    winjackpot += 99
                elif bankertype == 9:
                    usermoney += round((jackpotmoney * 0.1) - 1,2)
                    winjackpot += round((jackpotmoney * 0.1) - 1,2)
                else: 
                    usermoney += round(jackpotmoney - 1,2)
                    winjackpot += round(jackpotmoney - 1,2)
            if playertype >= 6 or bankertype >= 6:
                sleep(5)
        if battleresult == 13:
            win += 1
            if bankertype == 0:
                usermoney += 0.95
            else:
                if usertype == 3:
                    usermoney += 4.75
                elif usertype == 4:
                    usermoney += 6.65
                elif usertype == 5:
                    usermoney += 8.55
                elif usertype == 6:
                    usermoney += 14.25
                elif usertype == 7:
                    usermoney += 18.05
                elif usertype == 8:
                    usermoney += 38.95
                elif usertype == 9:
                    usermoney += 57.95
                elif usertype == 10:
                    usermoney += 95.95
                else:
                    usermoney += 2.85
        elif battleresult == 14:
            lose += 1
            usermoney -= 3
        else:
            dues += 1 
        if bankertype != btype:
            screenshot()
            Msg1 = "莊家牌型判斷有誤 "
            Msg1 += " 牌型: " + str(bankerlist) 
            Msg1 += " 圖片顯示: " + str(bankertype)
            Msg1 += " 計算牌型: " + str(btype)
            callbot(Msg1)
            break
        
        #牌型判斷，0:不成局，1:高牌，2:對子，3:兩對，4:三條，5:順子，6:同花，7:葫蘆，8:四條，9:同花順，10:皇家同花順
        if playertype == 1:
            onecountlog += 1
        elif playertype == 2:
            paircountlog += 1
        elif playertype == 3:
            twopaircountlog += 1
        elif playertype == 4:
            threecountlog += 1
        elif playertype == 5:
            straightcountlog += 1
        elif playertype == 6:
            samecountlog += 1
        elif playertype == 7:
            fullhousecountlog += 1
        elif playertype == 8:
            fourcountlog += 1
        elif playertype == 9:
            samestraightcountlog += 1
        elif playertype == 10:
            royalsamestraightcountlog += 1
        tt = gamecount % 50
        if tt == 0:
            Msg2 = "目前進行場次 :" + str(gamecount)+ "\n"
            Msg2 += "玩家獲得牌型場數"+ "\n"
            Msg2 += "高牌:" + str(onecountlog)+ "\n"
            Msg2 += "對子:" + str(paircountlog)+ "\n"
            Msg2 += "兩對:" + str(twopaircountlog)+ "\n"
            Msg2 += "三條:" + str(threecountlog)+ "\n"
            Msg2 += "順子:" + str(straightcountlog)+ "\n"
            Msg2 += "同花:" + str(samecountlog)+ "\n"
            Msg2 += "葫蘆:" + str(fullhousecountlog)+ "\n"
            Msg2 += "四條:" + str(fourcountlog)+ "\n"
            Msg2 += "同花順:" + str(samestraightcountlog)+ "\n"
            Msg2 += "皇家同花順:" + str(royalsamestraightcountlog) + "\n"
            Msg2 += "勝場數:" + str(win) + "\n"
            Msg2 += "敗場數:" + str(lose) + "\n"
            Msg2 += "和場數:" + str(dues) + "\n"
            Msg2 += "不成局勝場數:" + str(none) + "\n"
            Msg2 += "贏得JP金額:" + str(winjackpot)      
            callbot(Msg2)
        sleep(5)    
        if round(usermoney,2) != round(float(poco('User1').offspring('MoneyText').attr('text')),2):
            msg1 = "金額計算有誤"+ "\n"
            msg1 += "計算金額:" + str(round(usermoney,2)) + "\n"
            msg1 += "畫面金額:" + poco('User1').offspring('MoneyText').attr('text')
            callbot(msg1)
            break
except:
    Msg1 = "海盜船異常結束,原因為: " +  str(sys.exc_info()[0]) + " ,結束畫面如下:"
    callbot(Msg1)
    #screenshot()    
else:
    callbot("海盜船測試完畢")
