# -*- encoding=utf8 -*-
__author__ = "user16"

from airtest.core.api import *

auto_setup(__file__)

from poco.drivers.unity3d import UnityPoco
poco = UnityPoco()

count = 0
while (count < 31):
    poco("TextBetTimer").wait(60).exists()
    sleep(1)
    poco("ButtonBet1").click()
    poco("ButtonBet2").click()
    poco("ButtonBet3").click()
    poco("ButtonBet4").click()
    sleep(11)
    count = int(count) + 1