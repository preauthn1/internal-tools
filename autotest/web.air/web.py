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
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
auto_setup(__file__)


def callbot(Msg):
    r=requests.get("https://example.invalid/bot1343006432:AAHEmnjGyaOnOgDxr8luTgZ1oItQqGBdVAU/sendMessage?chat_id=-401714840&text=%s"%Msg)
    return()

def time():
    sleep(20)    
    return()

def screenshot():
    saveDir= "D:\\"     
    t=time.localtime()
    filename = time.strftime("%Y-%m-%d_%H-%M-%S", t) + '.jpg'
    img =ImageGrab.grab()
    img.save(os.path.join(saveDir, filename), 'jpeg')
    files = {'photo': open(os.path.join(saveDir, filename), 'rb')}
    url = "https://example.invalid/bot1343006432:AAHEmnjGyaOnOgDxr8luTgZ1oItQqGBdVAU/sendPhoto"
    data = {'chat_id' : "-401714840"}
    #url = "https://example.invalid/bot871711435:AAGYzFP4tNBQ-5E2j45a8qFBiDZDbmIlY2M/sendPhoto"  #Lisa
    #data = {'chat_id' : "-216697951"}  #Lisa    
    r= requests.post(url, files=files, data=data)



#用網頁模擬手機型號
# mobileEmulation = {'deviceName': 'test'}
mobileEmulation = {'deviceMetrics': {'width': 360, 'height': 640, 'pixelRatio': 3}, 'userAgent': 'Mozilla/5.0 (Linux; Android 6.0.1; Moto G (4)) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Mobile Safari/537.36'}
options = webdriver.ChromeOptions()
options.add_experimental_option('mobileEmulation', mobileEmulation)
gamecount = 0
#callbot("阿伯出發!")   

try:
    while(True):  
        #開網頁
        gamecount += 1
        driver = webdriver.Chrome(executable_path='chromedriver.exe', chrome_options=options)
        driver.get('https://example.invalid/Lobby')
        sleep(10)
        
        #登入帳密
        userid = driver.find_element_by_id('InputUser')
        userid.send_keys("frank05")
        password = [REDACTED]'InputPassword')
        password.send_keys("123456")
        button = driver.find_element_by_xpath('//*[@id="C711"]/div[1]/div[3]/div[2]/form/div[3]/button')
        button.click()
        time()
        if exists(Template(r"tpl1623307033844.png", threshold=0.5999999999999999, record_pos=(-0.471, 0.109), resolution=(1920, 1080))):
    
        #if exists(Template(r"tpl1612427265698.png", record_pos=(-0.316, -0.189), resolution=(1920, 1080))):
        #if exists(Template(r"tpl1621618228893.png", record_pos=(-0.315, -0.189), resolution=(1920, 1080))):
            try:
                result = driver.find_element_by_xpath('//*[@id="C711"]/div/div/div[1]/div[3]/div[1]/div/div[2]/div[2]/div/div').is_displayed()
            except:    
                driver.close()
            else:
                msg1 = "vworld token錯誤了 阿伯"
                callbot(msg1)
                driver.close()
        else:
            #screenshot()
            msg1 = "vworld出事了 阿伯"
            callbot(msg1)
            driver.close()
            continue
            
            
        driver = webdriver.Chrome(executable_path='chromedriver.exe', chrome_options=options)    
        driver.get('https://example.invalid/login#')
        sleep(10)
        
        #登入帳密
        userid = driver.find_element_by_name('login')
        userid.send_keys("vwmfrank")
        password = [REDACTED]'password')
        password.send_keys("a123123")
        button = driver.find_element_by_id('loginBtn')
        button.click()
        time()
        
        if exists(Template(r"tpl1622432098292.png", threshold=0.8999999999999999, rgb=True, record_pos=(-0.332, -0.186), resolution=(1920, 1080))):
            driver.close()
        else:
            msg1 = "cs出事了 阿伯"
            callbot(msg1)
            driver.close()
            continue
        
        #driver.get('https://example.invalid/Setting')
        #time()
        #driver.find_element_by_xpath("//button[@type='submit']").click()
        #time()
except:
    msg1 = "vworld登入異常結束,原因為: " +  str(sys.exc_info()[0]) + " ,結束畫面如下:"
    callbot(msg1) 
    #screenshot()
else:
    callbot("vworld登入測試完畢")    
    
    

