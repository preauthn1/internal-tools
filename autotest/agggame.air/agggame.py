# -*- encoding=utf8 -*-
__author__ = "frank"

from airtest.core.api import *

auto_setup(__file__)
# -*- encoding=utf8 -*-
__author__ = "Frank"

from airtest.core.api import *
from PIL import Image
from PIL import ImageGrab
import time
import os
import sys
import requests
from selenium import webdriver
from selenium.webdriver.common.keys import Keys
from airtest_selenium.proxy import WebChrome
auto_setup(__file__)

def callbot(Msg):
    r=requests.get("https://example.invalid/bot1343006432:AAHEmnjGyaOnOgDxr8luTgZ1oItQqGBdVAU/sendMessage?chat_id=-401714840&text=%s"%Msg)
    return()

def time():
    sleep(10)
    return()


def play():
    sleep(3)
    driver.airtest_touch(Template(r"tpl1624346773505.png", record_pos=(-0.167, 0.222), resolution=(1920, 1080)))
    time()
    driver.back()
    #driver.airtest_touch(Template(r"tpl1624346782056.png", record_pos=(-0.369, -0.203), resolution=(1920, 1080)))
    return()


def agg(): 
    driver.find_element_by_class_name("game_290008").click() #西部牛仔
    play()
    driver.find_element_by_class_name("game_290009").click() #成吉思汗
    play()
    driver.find_element_by_class_name("game_290013").click() #地產大亨
    play()
    driver.find_element_by_class_name("game_290051").click() #太極財神
    play()
    driver.find_element_by_class_name("game_290041").click() #發發發
    play()
    driver.find_element_by_class_name("game_290050").click() #太極童子
    play()
    driver.find_element_by_class_name("game_290038").click() #黃金彌勒佛
    play()
    driver.find_element_by_class_name("game_290039").click() #武則天
    play()
    driver.find_element_by_class_name("game_290046").click() #幸運轉輪
    play()
    driver.find_element_by_class_name("game_290048").click() #巨龍寶藏
    play()
    driver.find_element_by_class_name("game_290030").click() #星月魔女
    play()
    driver.find_element_by_class_name("game_290034").click() #金礦哥布林
    play()
    driver.find_element_by_class_name("game_290031").click() #焱狼
    play()
    driver.find_element_by_class_name("game_290042").click() #亞瑟王
    play()
    driver.find_element_by_class_name("game_290033").click() #波斯帝國
    play()
    driver.find_element_by_class_name("game_290043").click() #糖果屋
    play()
    driver.find_element_by_class_name("game_290049").click() #埃及豔后
    play()
    driver.find_element_by_class_name("game_290025").click() #加勒比寶藏
    play()
    driver.find_element_by_class_name("game_290036").click() #海克力斯
    play()
    driver.find_element_by_class_name("game_290037").click() #波賽頓
    play()
    driver.find_element_by_class_name("game_290044").click() #太極舞獅
    play()
    driver.find_element_by_class_name("game_290047").click() #法老王
    play()
    driver.find_element_by_class_name("game_290040").click() #快樂農場
    play()
    driver.find_element_by_class_name("game_290045").click() #極速賽馬
    play()
    driver.find_element_by_class_name("game_290059").click() #梁祝
    play()
    driver.find_element_by_class_name("game_290052").click() #福爾摩沙 
    play()
    driver.find_element_by_class_name("game_290058").click() #四大美人
    play()
    driver.find_element_by_class_name("game_290057").click() #甜蜜花園
    play()
    driver.find_element_by_class_name("game_290054").click() #幸運轉輪 Deluxe
    play()
    driver.find_element_by_class_name("game_290060").click() #綠野仙蹤 
    play()
    driver.find_element_by_class_name("game_290053").click() #本能寺之變
    play()
    driver.find_element_by_class_name("game_290055").click() #大魔術師
    play()
    driver.find_element_by_class_name("game_290056").click() #怪物星球
    play()
    driver.find_element_by_class_name("game_290061").click() #瑪雅之謎
    play()
    driver.find_element_by_class_name("game_290062").click() #黃金福娃
    play()
    driver.find_element_by_class_name("game_290063").click() #阿拉丁
    play()
    driver.find_element_by_class_name("game_290064").click() #超級鑽
    play()
    driver.find_element_by_class_name("game_290065").click() #法海傳
    play()
    
    return
    
    


driver = WebChrome()
driver.implicitly_wait(20)
driver.get("http://0.0.0.0:8080/?playerName=frankvnd001&token=[REDACTED]")
agg()
callbot("1倍比例vn正式AGG遊戲一次投注完成")