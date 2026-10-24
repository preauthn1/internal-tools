# -*- encoding=utf8 -*-
__author__ = "user16"

from airtest.core.api import *

auto_setup(__file__)

from poco.drivers.unity3d import UnityPoco
poco = UnityPoco()
while(True):
    poco("background").wait(1800).exists()
    poco("Node_Left").offspring("ico").click()
    poco("ReadyGame").wait(180).exists()
    poco("ReadyGame").click()
