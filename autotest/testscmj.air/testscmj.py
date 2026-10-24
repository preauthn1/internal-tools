# -*- encoding=utf8 -*-
__author__ = "user16"

from poco.drivers.unity3d import UnityPoco
poco = UnityPoco()
import requests
from airtest.core.api import *

auto_setup(__file__)

while (True):
    poco("1").click()
    sleep(20)
    if poco("_show").exists():
        poco("_show").click()
        poco("_exit").click()
    else:
        break
r = requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=被踢出桌子")
       
