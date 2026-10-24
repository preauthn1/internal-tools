# -*- encoding=utf8 -*-
__author__ = "User"

from airtest.core.api import *

auto_setup(__file__)# -*- encoding=utf8 -*-
__author__ = "Frank"

from airtest.core.api import *


from poco.drivers.android.uiautomation import AndroidUiautomationPoco
poco = AndroidUiautomationPoco(use_airtest_input=True, screenshot_each_action=False)
from PIL import Image
from PIL import ImageGrab
import time
import os
import sys
import requests
auto_setup(__file__)

def callbot(Msg):
    r=requests.get("https://example.invalid/bot1361458875:AAG3Ekdf5LCwPXLE1zffjDOeiG32rOXQkOQ/sendMessage?chat_id=-487213300&text=%s"%Msg)
    return()

def screenshot():
    saveDir= "D:\\"     
    t=time.localtime()
    filename = time.strftime("%Y-%m-%d_%H-%M-%S", t) + '.jpg'
    img =ImageGrab.grab()
    img.save(os.path.join(saveDir, filename), 'jpeg')
    files = {'photo': open(os.path.join(saveDir, filename), 'rb')}
    url = "https://example.invalid/bot1361458875:AAG3Ekdf5LCwPXLE1zffjDOeiG32rOXQkOQ/sendPhoto"
    data = {'chat_id' : "-487213300"}
    #url = "https://example.invalid/bot871711435:AAGYzFP4tNBQ-5E2j45a8qFBiDZDbmIlY2M/sendPhoto"  #Lisa
    #data = {'chat_id' : "-216697951"}  #Lisa    
    r= requests.post(url, files=files, data=data)


try:
    while(True):
        poco(text="game_icon_filter_slots").click()
        poco("android.widget.FrameLayout").offspring("com.android.browser:id/webview_wrapper").child("android.webkit.WebView").offspring("app ").child("C711Mobile").child("android.view.View").child("C711Mobile").child("android.view.View").child("android.view.View").child("android.view.View")[0].child("android.widget.Image")[4].swipe([-0.0314, -0.1344])

        poco(text="159765181430").click()
        sleep(30)
        msg1 = str("88forture")
        assert_exists(Template(r"tpl1598606433535.png", record_pos=(0.001, 0.618), resolution=(720, 1280)), "Please fill in the test point.")

        touch(Template(r"tpl1600166656772.png", record_pos=(-0.411, -0.665), resolution=(720, 1280)))
        sleep(5)
        poco(text="game_icon_filter_slots").click()                
        poco("android.widget.FrameLayout").offspring("com.android.browser:id/webview_wrapper").child("android.webkit.WebView").offspring("app ").child("C711Mobile").child("android.view.View").child("C711Mobile").child("android.view.View").child("android.view.View").child("android.view.View")[0].child("android.widget.Image")[4].swipe([-0.0314, -0.1344])
        poco(text="159765180096").click()
        sleep(30)
        msg1 = str('jokerjewels')
        
        assert_exists(Template(r"tpl1598607098860.png", record_pos=(-0.004, 0.611), resolution=(720, 1280)), "Please fill in the test point.")
        touch(Template(r"tpl1600166656772.png", record_pos=(-0.411, -0.665), resolution=(720, 1280)))
        sleep(10)
        poco(text="game_icon_filter_slots").click()
        poco(text="159765153566").click()
        sleep(30)
        msg1 = str('luckybat')
        assert_exists(Template(r"tpl1598607435655.png", record_pos=(0.004, 0.613), resolution=(720, 1280)), "Please fill in the test point.")
        touch(Template(r"tpl1600166656772.png", record_pos=(-0.411, -0.665), resolution=(720, 1280)))
        sleep(10)
        poco(text="game_icon_filter_slots").click()
        poco(text="159765150220").click()
        sleep(30)
        msg1 = str('frogsinflies')
        assert_exists(Template(r"tpl1598607480790.png", record_pos=(0.001, 0.61), resolution=(720, 1280)), "Please fill in the test point.") 
        touch(Template(r"tpl1600166656772.png", record_pos=(-0.411, -0.665), resolution=(720, 1280)))
        sleep(10)
        poco(text="game_icon_filter_slots").click()
        poco(text="159765147848").click()
        sleep(30)
        msg1 = str('queen')
        assert_exists(Template(r"tpl1598607601078.png", record_pos=(0.003, 0.61), resolution=(720, 1280)), "Please fill in the test point.")
        touch(Template(r"tpl1600166656772.png", record_pos=(-0.411, -0.665), resolution=(720, 1280)))
        sleep(10)
        poco(text="game_icon_filter_slots").click()
        poco(text="159765139163").click()
        sleep(30)
        msg1 = str('buffalo')
        assert_exists(Template(r"tpl1598607669560.png", record_pos=(0.003, 0.617), resolution=(720, 1280)), "Please fill in the test point.")
        touch(Template(r"tpl1600166656772.png", record_pos=(-0.411, -0.665), resolution=(720, 1280)))
        sleep(10)
        poco(text="game_icon_filter_slots").click()
        poco(text="159765145320").click()
        sleep(30)
        msg1 = str('888')
        assert_exists(Template(r"tpl1598607720554.png", record_pos=(0.003, 0.607), resolution=(720, 1280)), "Please fill in the test point.")
        touch(Template(r"tpl1600166656772.png", record_pos=(-0.411, -0.665), resolution=(720, 1280)))
        sleep(10)
        poco(text="game_icon_filter_slots").click()
        poco(text="159980400946").click()
        sleep(30)
        msg1 = str('5luckylion')
        assert_exists(Template(r"tpl1598607776073.png", record_pos=(0.001, 0.608), resolution=(720, 1280)), "Please fill in the test point.")
        touch(Template(r"tpl1600166656772.png", record_pos=(-0.411, -0.665), resolution=(720, 1280)))
        sleep(10)
        poco(text="159765158151").click()
        sleep(30)
        msg1 = str('多福多财')
        assert_exists(Template(r"tpl1598607855976.png", record_pos=(-0.239, 0.689), resolution=(720, 1280)), "Please fill in the test point.")
        touch(Template(r"tpl1600166656772.png", record_pos=(-0.411, -0.665), resolution=(720, 1280)))
        sleep(10)
        poco(text="159765153566").click()
        sleep(30)
        msg1 = str('糖果屋')
        assert_exists(Template(r"tpl1598608035716.png", record_pos=(-0.368, 0.576), resolution=(720, 1280)), "Please fill in the test point.")

        touch(Template(r"tpl1600166656772.png", record_pos=(-0.411, -0.665), resolution=(720, 1280)))
        sleep(10)
        poco(text="159765150220").click()
        sleep(30)
        msg1 = str('聚宝盆')
        assert_exists(Template(r"tpl1598606803563.png", record_pos=(-0.201, 0.175), resolution=(720, 1280)), "Please fill in the test point.")
        touch(Template(r"tpl1600166656772.png", record_pos=(-0.411, -0.665), resolution=(720, 1280)))
        sleep(10)
        poco(text="159765147848").click()
        sleep(30)
        msg1 = str('龙虎斗')
        assert_exists(Template(r"tpl1598426163694.png", record_pos=(11.655, 7.145), resolution=(100, 100)), "判斷有按鈕")
        touch(Template(r"tpl1600166656772.png", record_pos=(-0.411, -0.665), resolution=(720, 1280)))
        sleep(10)
        poco(text="159765139163").click()
        sleep(30)
        msg1 = str('百家乐')
        assert_exists(Template(r"tpl1598608499602.png", record_pos=(0.228, -0.215), resolution=(720, 1280)), "Please fill in the test point.")
        touch(Template(r"tpl1600166656772.png", record_pos=(-0.411, -0.665), resolution=(720, 1280)))
        sleep(10)
        poco(text="159765145320").click()
        sleep(30)
        msg1 = str('红黑大战')
        assert_exists(Template(r"tpl1598608601728.png", record_pos=(-0.285, 0.535), resolution=(720, 1280)), "Please fill in the test point.")
        touch(Template(r"tpl1600166656772.png", record_pos=(-0.411, -0.665), resolution=(720, 1280)))
        sleep(10)
        #sleep(600)
except:
    msg1 += "移工站大廳異常結束,原因為: " +  str(sys.exc_info()[0]) + " ,結束畫面如下:"
    callbot(msg1) 
    screenshot()
else:
    callbot("移工站大廳測試完畢")driver.find_element_by_xpath("//*[@id=\"appRoot\"]/div/div").click()
driver.find_element_by_xpath("//input[@maxlength='30']").click()
driver.find_element_by_xpath("//input[@maxlength='30']").click()
driver.find_element_by_xpath("//input[@maxlength='30']").click()
