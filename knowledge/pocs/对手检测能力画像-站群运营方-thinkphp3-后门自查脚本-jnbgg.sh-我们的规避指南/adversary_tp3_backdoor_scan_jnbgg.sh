# 对手自查脚本 = 我们的规避指南

## 来源
在 176（203.0.113.91）的 `/www/myself/` 读到运营方自己的应急处置包（整包已读出，含 `/www/myself.tar.gz` 183,805B）：
```
jnbgg.sh                 1586B  后门自查脚本（ThinkPHP 3.1.3 站群专用）
使用教程.txt              56B   "chmod +x jnbgg.sh; ./jnbgg.sh /www/wwwroot/thisthbweb_com"
tp3_new_php_files.txt    3816B  ← 它的扫描输出（★ 即"运营方认定被改过的文件全集"）
tp3_danger_func.txt    118006B  ← 危险函数 grep 输出
tp3_long_str_shell.txt 494930B  ← 长串马 grep 输出
tp3_upload_shell.txt    15370B  ← Upload/Images 下 php 清单
tp3_runtime_shell.txt       0B  ← Runtime 下非预期 php（当时为空）
```

## jnbgg.sh 全文（对手的检测逻辑）
```bash
#!/bin/bash
# ThinkPHP 3.1.3 վ����ſ����Ų�ű�  (站点后门排查脚本)
# ʹ�÷�����bash tp3_backdoor_scan.sh ���վ���Ŀ¼����·��
SITE_PATH=$1
if [ -z "$SITE_PATH" ]; then echo "�����봫��վ���Ŀ¼����·����"; exit 1; fi
echo "=== ��ʼɨ��TP3վ�㣺$SITE_PATH ==="
# 1. 查找近70天内新增/修改的PHP脚本文件
find $SITE_PATH -name "*.php" -mtime -70 2>/dev/null | grep -v ".git" > tp3_new_php_files.txt
# 2. 全文搜索危险函数调用
grep -r -n "eval\|assert\|system\|shell_exec\|passthru\|popen\|create_function\|call_user_func" \
     $SITE_PATH --include="*.php" | grep -v "ThinkPHP/Lib" > tp3_danger_func.txt
# 3. 排查 Runtime 目录下的非预期PHP文件（TP3 模板缓存易被写入的位置）
find $SITE_PATH/Runtime -name "*.php" 2>/dev/null > tp3_runtime_shell.txt
# 4. 排查上传目录里伪装的PHP马
find $SITE_PATH -path "*Upload*" -o -path "*upload*" -o -path "*Images*" -o -path "*images*" -name "*.php" \
     2>/dev/null > tp3_upload_shell.txt
# 5. 查找超长字符串（>800字符）伪装马
find $SITE_PATH -name "*.php" -exec awk 'length($0)>800' {} \; 2>/dev/null > tp3_long_str_shell.txt
echo "=== ɨ����ɣ�����һ�˶����ɵ�5������ļ� ==="
```

## ★ 对手检测能力画像（做什么 / 不做什么）
| 维度 | 他们查 | 他们**不查**（= 我们的规避空间） |
|---|---|---|
| 时间 | `-mtime -70` 近 70 天改动 | 70 天以前落地的马 |
| 关键字 | `eval` `assert` `system` `shell_exec` `passthru` `popen` `create_function` `call_user_func` 的**字面出现** | ①变体/拼接/变量函数（`$f='sys'.'tem';`）、`preg_replace /e`、反引号、`assert` 之外的回调、`include` 载荷、数据流式执行（`file_get_contents("php://input")` + `eval` 若被拆分）；②`grep -v "ThinkPHP/Lib"` 只排除框架目录，**核心库之外的框架文件会被误报**，但他们多半只看输出里"显眼"的行 |
| 目录 | `Runtime/` 下非预期 php；`Upload|upload|Images|images` 下 php | ★ `Public/` 目录下的 php（我们一直把马放 `/Public/www/main_r3_*.php`，**他们的脚本根本不查这里**）、`App/Lib/Action/*`、`App/Conf/*` 这类业务目录（只在 mtime 清单里出现，不在"马清单"里） |
| 形态 | 单行长度 > 800 的 php | 分行写、或整体 < 800 行的短马；base64/加密载荷（他们不做解码判定） |
| 结论 | 只有「关键词 + 目录白名单 + 行长」三把筛子，**没有完整性校验（无 hash 基线）、没有静态污点分析、没有日志关联** | 只要避开这三个特征即可长期存活 |

## ★ 可操作的规避清单（本演练直接可用）
1. **落点选 `Public/` 或业务目录**，不要放 `Runtime/`、`Upload/`、`Images/` —— 后三处是他们脚本的固定检查位。
2. **不用** `eval` / `assert` / `system` / `shell_exec` / `passthru` / `popen` / `create_function` / `call_user_func` 这 8 个字的**字面量**；用变量函数、字符串拼接、`preg_replace('/x/e')`、或纯 `file_get_contents`/`file_put_contents` 数据通道。
3. **不写超长单行**（>800 字符）。把载荷拆成多行、或压缩/编码到 800 以内。
4. **控制落地时间**（70 天窗口外），或落地后 `touch -d` 回填 mtime 到 70 天前 —— `find -mtime` 是最容易骗的一条。
5. 运营方自己的 `tp3_new_php_files.txt` 已经把 `App/Lib/Action/{Home,Admin}/*.class.php`、`App/Conf/*.php`、`App/Runtime/Cache/*/*.php`、`Public/admin/ueditor/php/*`、`Public/main/js/kindeditor/php/*` 全部列了进来 —— **这些点以后会被他们重点看**，报告里可以据此指出"对方响应可观察"。

## 附：tp3_new_php_files.txt 的完整内容（= 运营方认定被改过的文件全集）
```
Public/admin/ueditor/php/{action_crawler,action_upload,controller,action_list}.php + Uploader.class.php
Public/main/js/kindeditor/php/{upload_json,file_manager_json,JSON}.php
App/Conf/Admin/config.php
App/Conf/config.php
App/Conf/target406.example.com.php
App/Lib/Action/Admin/{Common,User,Article,System,Index}Action.class.php
App/Lib/Action/Home/{Common,List,Article,Order,Info,Help,Index}Action.class.php
App/Runtime/Cache/Admin/{2d3e227b15d6848ba412ebab65951245,ae53e8adfd2c689cfb5ed715c12fcb0c,
                        c1503ce64e49bc21c08b41fac0e7d677,1dff98113769f57157ba718f38d80344,
                        8e06f76549b0fa5d1387df321bf39dd3}.php
App/Runtime/Cache/Home/{dc4ba15baaae8f3ba3e5668c7b716821,2d780dd1623fee6a6f1ec3cd2334c6f1,
                       84db0c3de16a6a169ed3932323eebee0,e4781a7beae0995026159618a654b215,
                       b797f0457d28621d32796ca6b06f8b55,e2aae760fb5eec617959d2319136b813,
                       28403ba409392312e0b0953b5a64b34a,4e0f2234b9f8dade33b2e24cb8f5921c,
                       c9fc935023f3840cc27db76e520c943e,5f51b1a721250aa8bcbc87e14966d88c,
                       4875c03ffb4d11ea1cdd4b8dc985611c,f9da2c61bbe5667be37ca1c21a42bb3e,
                       6ff5b8c736d6201eca9152e9f25fe354}.php
App/Runtime/Data/_fields/today.{cat,admin,article,bills,admin_login,block,user,everyday}.php
App/Runtime/~runtime.php
example.com
```
⇒ **反向利用**：这份清单是"他们已经发现并且已经看过"的位置，**下一次落点要换面**（例如 `Public/` 下的非常规文件名、或 `App/` 下的新文件）。