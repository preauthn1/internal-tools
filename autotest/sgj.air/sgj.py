# -*- encoding=utf8 -*-
__author__ = "user16"

from airtest.core.api import *
from poco.drivers.unity3d import UnityPoco
poco = UnityPoco()

auto_setup(__file__)
import time
count = 0

while(int(count) < 10):#控制次數
    count = int(count) + 1
    betbtn_1 = 0
    betbtn_2 = 0
    betbtn_3 = 0
    betbtn_4 = 0
    betbtn_5 = 0
    betbtn_6 = 0
    betbtn_7 = 0
    betbtn_8 = 0
    times = 0
    totaltimes = 0    
    usermoney = poco('_glopUserCoin').attr('text')
    print('初始金額:',usermoney)
    
    poco('Canvas').offspring('center').child('1').child('_glfrtet').click()
    betbtn_1 = 1
    poco('Canvas').offspring('center').child('2').child('_glfrtet').click()
    betbtn_2 = 1
    poco('Canvas').offspring('center').child('3').child('_glfrtet').click()
    betbtn_3 = 1
    poco('Canvas').offspring('center').child('4').child('_glfrtet').click()
    betbtn_4 = 1
    poco('Canvas').offspring('center').child('5').child('_glfrtet').click()
    betbtn_5 = 1
    poco('Canvas').offspring('center').child('6').child('_glfrtet').click()
    betbtn_6 = 1
    poco('Canvas').offspring('center').child('7').child('_glfrtet').click()
    betbtn_7 = 1
    poco('Canvas').offspring('center').child('8').child('_glfrtet').click()
    betbtn_8 = 1    
    
    poco('Canvas').offspring('zhongjing').offspring('icon').wait(60).exists()
    sleep(2)
    winicon = poco("lottery(Clone)")[0].child('icon').attr('texture')
    if poco("lottery(Clone)")[0].child('bet').exists():
        times = poco("lottery(Clone)")[0].child('bet').attr('text')
        lentimes = len(times)
        if int(lentimes) < 3:
            times = times[1]
        else:
            times = times[1]+times[2]
    if betbtn_1 == 1:
        if winicon == poco('Canvas').offspring('center').child('1').child('icon').attr('texture'):
            totaltimes = int(totaltimes)+int(times)-1
        else:
            totaltimes = int(totaltimes) - 1
    if betbtn_2 == 1:
        if winicon == poco('Canvas').offspring('center').child('2').child('icon').attr('texture'):
            totaltimes = int(totaltimes)+int(times)-1
        else:
            totaltimes = int(totaltimes) - 1
    if betbtn_3 == 1:
        if winicon == poco('Canvas').offspring('center').child('3').child('icon').attr('texture'):
            totaltimes = int(totaltimes)+int(times)-1
        else:
            totaltimes = int(totaltimes) - 1
    if betbtn_4 == 1:
        if winicon == poco('Canvas').offspring('center').child('4').child('icon').attr('texture'):
            totaltimes = int(totaltimes)+int(times)-1
        else:
            totaltimes = int(totaltimes) - 1
    if betbtn_5 == 1:
        if winicon == poco('Canvas').offspring('center').child('5').child('icon').attr('texture'):
            totaltimes = int(totaltimes)+int(times)-1
        else:
            totaltimes = int(totaltimes) - 1
    if betbtn_6 == 1:
        if winicon == poco('Canvas').offspring('center').child('6').child('icon').attr('texture'):
            totaltimes = int(totaltimes)+int(times)-1
        else:
            totaltimes = int(totaltimes) - 1
    if betbtn_7 == 1:
        if winicon == poco('Canvas').offspring('center').child('7').child('icon').attr('texture'):
            totaltimes = int(totaltimes)+int(times)-1
        else:
            totaltimes = int(totaltimes) - 1
    if betbtn_8 == 1:
        if winicon[0:3] == 'BAR':
            if len(winicon) < 10:
                totaltimes = int(totaltimes)+49
            else:
                totaltimes = int(totaltimes)+99
    if int(totaltimes) > 0 :
        usermoney = round(((float(usermoney))+float(totaltimes)*0.95),2)
    else:
        usermoney = round(float(usermoney) + float(totaltimes),2)
    poco('_glbankTime').wait(60).exists()
    if float(poco('_glopUserCoin').attr('text')) == float(usermoney): #判斷上局玩家金額是否正確
        print('正確')
    else:
        print ('錯誤')
        print (poco('_glopUserCoin').attr('text'))
        print (usermoney)
        print (time.asctime(time.localtime(time.time())))
        break