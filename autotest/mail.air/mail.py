# -*- encoding=utf8 -*-
import smtplib
from email.mime.text import MIMEText

gmail_user = 'user@example.invalid'
gmail_password = [REDACTED]' # your gmail password

msg = MIMEText('content')
msg['Subject'] = 'Test'
msg['From'] = gmail_user
msg['To'] = 'user@example.invalid'

server = smtplib.SMTP_SSL('example.invalid', 465)
server.ehlo()
server.login(gmail_user, gmail_password)
server.send_message(msg)    
server.quit()

print('Email sent!')