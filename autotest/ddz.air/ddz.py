# -*- encoding=utf8 -*-
__author__ = "user16"

from airtest.core.api import *

auto_setup(__file__)

from poco.drivers.unity3d import UnityPoco
poco = UnityPoco()
count = 0
while(int(count) < 10):  #控制次數
    count = int(count) + 1
    usermoney = float(poco('Canvas').offspring('Node_Down').offspring('Text_Money').attr('text'))   
    poco('BtnTuoGuang').wait(120).exists()
    poco('BtnTuoGuang').click()  
    poco('Button_JXYX').wait(180).exists()
    if poco('Button_JXYX').exists():
        odd = poco('Text_bs1').attr('text')
        if poco('youxi_shibai').exists(): 
            usermoney = round((float(usermoney)-(1*int(odd))),2)
        else:
            usermoney = round((float(usermoney)+(0.95*int(odd))),2)
        poco('Button_JXYX').click()       
    if float(usermoney) ==  float(poco('Canvas').offspring('Node_Down').offspring('Text_Money').attr('text')):
        print('正確')
    else:
        print('錯誤')
        
    
    
