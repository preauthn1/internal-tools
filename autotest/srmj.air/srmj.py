# -*- encoding=utf8 -*-
__author__ = "user16"

from airtest.core.api import *

auto_setup(__file__)

from poco.drivers.unity3d import UnityPoco
poco = UnityPoco()
while(True):
    poco("BtnBuJiaBei").wait(60).exists()
    poco("BtnBuJiaBei").click()
#    poco("_show").click()
#    poco("_set").click()
#    poco("Trustee").click()
#    poco("close").click()
#    poco("_show").click()
    poco("background").wait(1800).exists()
#    poco(texture="sprite_15").click()
    poco("Click_LSEnter").click()
    poco("ReadyGame").wait(180).exists()
    poco("ReadyGame").click()

