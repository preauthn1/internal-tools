# -*- encoding=utf8 -*-
__author__ = "user16"

from airtest.core.api import *
from poco.drivers.unity3d import UnityPoco
poco = UnityPoco()

auto_setup(__file__)
keyevent("POWER")
swipe(Template(r"tpl1551251594987.png", record_pos=(0.008, 0.679), resolution=(1080, 1920)), vector=[-0.0083, -0.7754])

#touch(Template(r"tpl1551250092473.png", record_pos=(-0.363, -0.182), resolution=(1080, 1920)))
#sleep(15)
#wait(Template(r"tpl1551250214800.png", record_pos=(0.374, -0.234), resolution=(1920, 1080)),20)
#touch(Template(r"tpl1551250223707.png", record_pos=(0.373, -0.233), resolution=(1920, 1080)))
#touch(Template(r"tpl1551250289035.png", record_pos=(-0.009, -0.404), resolution=(1080, 1920)))
#keyevent("HOME")
#touch(Template(r"tpl1551250672252.png", record_pos=(-0.369, 0.1), resolution=(1080, 1920)))
#touch(Template(r"tpl1551250726083.png", record_pos=(0.297, -0.143), resolution=(1920, 1080)))
