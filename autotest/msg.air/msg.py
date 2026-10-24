from airtest.core.api import *
using("hall.air")
def add():
    import requests
    from hall import msg1 
    # GET Method
    string = str(msg1)
    r = requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=%s" % string)