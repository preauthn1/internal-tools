# -*- encoding=utf8 -*-

__author__ = "user16"



from poco.drivers.unity3d import UnityPoco
poco = UnityPoco()

coinstart=float(poco('coin').attr('text'))
wait(Template(r"tpl1550199502192.png", threshold=0.85, record_pos=(-0.001, -0.014), resolution=(1280, 720)),30)
i=2
while(i<100):
 poco('ButtonBet1').click()
 wait(Template(r"tpl1550221995339.png", threshold=0.9, target_pos=5, rgb=False, record_pos=(-0.077, -0.013), resolution=(1280, 720)),30)
 coinend=poco('coin').attr('text')
 if exists(Template(r"tpl1550200159023.png", threshold=0.9000000000000001, record_pos=(0.291, 0.169), resolution=(1280, 720))):
     print ("龍贏") 
     if (float(poco('coin').attr('text')) == (coinstart+0.95)):
        coinstart=coinstart+0.95
        print('正確')
     else:
        print(coinstart)
        print(float(poco('coin').attr('text')))
        print('錯誤') 
        break
 elif exists(Template(r"tpl1550200429812.png", threshold=0.9000000000000001, record_pos=(0.287, 0.17), resolution=(1280, 720))):
    print ("虎贏")
    if (float(poco('coin').attr('text')) == (coinstart-1)):
        coinstart=coinstart-1
        print('正確')
    else:
        print(coinstart)
        print(poco('coin').attr('text'))
        print('錯誤') 
        break
 else:
    print ("和贏")
    if (float(poco('coin').attr('text')) == coinstart):
        print('正確')
    else:
        print('錯誤') 
        break
 i=i+1

