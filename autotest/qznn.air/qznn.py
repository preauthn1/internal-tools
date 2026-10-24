# -*- encoding=utf8 -*-
__author__ = "user16"

from airtest.core.api import *

auto_setup(__file__)

from poco.drivers.unity3d import UnityPoco
poco = UnityPoco()

while(True):
    poco("sanpai").wait(60).exists()
    poco("sanpai").click()
    poco("changedesk").wait(60).exists()
    poco("changedesk").click()
