# -*- encoding=utf8 -*-
__author__ = "user16"

from airtest.core.api import *
import requests
auto_setup(__file__)
import time
from poco.drivers.unity3d import UnityPoco
poco = UnityPoco()

localtime = time.asctime(time.localtime(time.time()))
while (True):
    #from poco.drivers.android.uiautomation import AndroidUiautomationPoco
    #poco = AndroidUiautomationPoco(use_airtest_input=True, screenshot_each_action=False)
    #poco("超級大贏家-伺服器").click()
    #sleep(10)

    #poco("ButtonEnter").click()
    #if poco('_HeadIcon').wait(30).exists():
    #    sleep(5)
    #else:
    #    string = str('大廳目前無法登入')
    #    r = requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=%s" % string)
    poco('GameLogo_1').click()  #點擊gameid"1"的遊戲
    poco('2').click()
    if poco("_show").wait(60).exists():
        poco("_show").click()
        poco("_exit").click()
        poco("_RoomlistBack").click()
    else:
        string = str('鬥地主無法登入')
        r = requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=%s" % string)

    poco('GameLogo_2').click()
    poco('2').click()
    if poco("_show").wait(60).exists():
        poco("_show").click()
        poco("_exit").click()
        poco("_RoomlistBack").click()
    else:
        string = str('搶莊牛牛無法登入')
        r = requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=%s" % string)

    poco('GameLogo_3').click()
    poco('2').click()
    if poco("_show").wait(60).exists():
        poco("_show").click()
        poco("_exit").click()
        poco("_RoomlistBack").click()
    else:
        string = str('炸金花無法登入')
        r = requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=%s" % string)
        
    if poco('GameLogo_9').wait(60).exists():
        poco('GameLogo_9').click()
    else:
        string = str('水果機未開啟')
        r = requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=%s" % string)
    if poco("_show").wait(60).exists():
        sleep(10)
        poco("_show").click()
        poco("_exit").click()
    else:
        string = str('水果機無法登入')
        r = requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=%s" % string)   

    if poco('GameLogo_13').wait(60).exists():
        sleep(3)
        poco('GameLogo_13').click()
    else:
        string = str('百家樂未開啟')
        r = requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=%s" % string)
    if poco("_show").wait(60).exists():
        sleep(10)
        poco("_show").click()
        poco("_exit").click()
    else:
        string = str('百家樂無法登入')
        r = requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=%s" % string)

    if poco('GameLogo_11').wait(60).exists():
        sleep(3)
        poco('GameLogo_11').click()
    else:
        string = str('龍虎鬥未開啟')
        r = requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=%s" % string)
    if poco("_show").wait(60).exists():
        sleep(10)
        poco("_show").click()
        poco("_exit").click()
    else:
        string = str('龍虎鬥無法登入')
        r = requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=%s" % string)

    #poco('GameLogo_15').click()
    #poco('1').click()
    #if poco("_show").wait(60).exists():
    #    poco("_show").click()
    #    poco("_exit").click()
    #    poco("_RoomlistBack").click()
    #else:
    #    string = str('四川麻將無法登入')
    #    r = requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=%s" % string)




    poco('_RightPage').click()
    if poco('GameLogo_14').wait(60).exists():
        poco('GameLogo_14').click()
    else:
        string = str('捕魚未開啟')
        r = requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=%s" % string)
    poco('1').click()
    if poco("_ShowMenu").wait(60).exists():
        sleep(15)
        poco("_ShowMenu").click()
        poco("_Exit").click()
        poco("_RoomlistBack").click()
    else:
        string = str('捕魚無法登入')
        r = requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=%s" % string)

    sleep(5)
    if poco('GameLogo_20').wait(60).exists():
        sleep(3)
        poco('GameLogo_20').click()
    else:
        string = str('紅黑大戰未開啟')
        r = requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=%s" % string)           
    if poco("_show").wait(60).exists():
        sleep(10)
        poco("_show").click()
        poco("_exit").click()
    else:
        string = str('紅黑大戰無法登入')
        r = requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=%s" % string)

    if poco('GameLogo_21').wait(60).exists():
        sleep(3)
        poco('GameLogo_21').click()
    else:
        string = str('紅黑牛牛未開啟')
        r = requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=%s" % string)       
    if poco("_show").wait(60).exists():
        sleep(10)
        poco("_show").click()
        poco("_exit").click()
    else:
        string = str('紅黑牛牛無法登入')
        r = requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=%s" % string)    
    if poco('GameLogo_19').wait(60).exists():
        sleep(3)
        poco('GameLogo_19').click()
    else:
        string = str('百人牛牛未開啟')
        r = requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=%s" % string)       
    if poco("_show").wait(60).exists():
        sleep(10)
        poco("_show").click()
        poco("_exit").click()
    else:
        string = str('百人牛牛無法登入')
        r = requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=%s" % string)
    poco('_LeftPage').click()  
