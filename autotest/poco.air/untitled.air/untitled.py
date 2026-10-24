# -*- encoding=utf8 -*-
__author__ = "user16"

from airtest.core.api import *

auto_setup(__file__)

from poco.drivers.unity3d import UnityPoco

poco = UnityPoco()
# for unity editor on windows
# poco = UnityPoco(('localhost', 5001), unity_editor=True)

ui = poco('...')
ui.click()

items = poco('Canvas').child('Gamelayer').offspring('BetGroup')
print(items[0].child('meterial_name').get_text())
                                                    