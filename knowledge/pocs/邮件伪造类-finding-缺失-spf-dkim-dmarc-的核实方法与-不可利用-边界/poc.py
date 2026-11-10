#!/usr/bin/env python3
# 只读核实：5 域 SPF/DKIM/DMARC/MX 与常见邮件前缀
# 依赖 dnspython: pip install dnspython
import dns.resolver, dns.exception

R = dns.resolver.Resolver(configure=False)
R.nameservers = ['8.8.8.8', '1.1.1.1']; R.timeout = 5; R.lifetime = 8
doms = ['target506.example.com', 'target510.example.com', 'target511.example.com', 'target512.example.com', 'target513.example.com']
selectors = ['default','selector1','selector2','google','k1','mail','dkim','s1','s2','smtp','protonmail','pm','proton']
prefixes = ['mail','smtp','webmail','imap','pop','mx','mx1','email','newsletter','em','edm','notify','support']

def q(name, rt):
    try:
        return 'OK ' + ' | '.join(str(x) for x in R.resolve(name, rt))
    except dns.resolver.NXDOMAIN: return 'NXDOMAIN'
    except dns.resolver.NoAnswer: return 'NODATA'
    except dns.exception.Timeout: return 'TIMEOUT'
    except Exception as e: return type(e).__name__ + ':' + str(e)[:60]

for d in doms:
    print('===== ' + d + ' =====')
    for rt in ('TXT','MX','NS'):
        print(rt, q(d, rt))
    print('DMARC', q('_dmarc.' + d, 'TXT'))
    print('DKIM', '; '.join(s+'='+q(s+'._domainkey.'+d,'TXT') for s in selectors))
    for p in prefixes:
        fq = p + '.' + d
        a, mx = q(fq,'A'), q(fq,'MX')
        if a not in ('NXDOMAIN','NODATA') or mx not in ('NXDOMAIN','NODATA'):
            print('MAILHOST', fq, a, mx)

# 判定：
# 全部 _dmarc=NXDOMAIN/NODATA、无 v=spf1、DKIM 无、无 MX/无邮件前缀主机
#  => 伪造技术前提成立，但无服务器侧可打面；唯一利用是向第三方投递伪造邮件（通常 RoE 排除）。
# 若出现 v=spf1 且 -all / _dmarc p=quarantine|reject => 该 finding 应被证否。
