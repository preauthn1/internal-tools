import requests
from bjl import usermoney
# GET Method

string = str(usermoney)
r = requests.get("https://example.invalid/bot746757886:AAFZdatm5kvME63ryezZ3Vp-d7A1bAQbbR0/sendMessage?chat_id=-270103520&text=%s" % string)
print(r.text)  
