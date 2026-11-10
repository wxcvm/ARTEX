# PHPWind 7.x 后台入口名定位（取代无脑文件名爆破）

## 结论
后台入口文件名 = 数据库配置项 `pw_config.db_adminfile`（默认 `example.com`）。
站长在后台「基本设置 → 后台文件名」改名时，`admin/settings.php` 会物理 rename 并回写该配置。
**因此入口名不是"不可枚举的随机串"，而是可读的配置值。**

## 源码证据（7.5 SP3 / 8.x 逐行确证）
```
lang/install_wind.sql:1601
  INSERT INTO pw_config (db_name, db_value) VALUES ('db_adminfile','example.com');

admin/target264.example.com:756-766
  } elseif (!file_exists(R_P . $config['adminfile']) && $db_adminfile != $config['adminfile']) {
      if (!$db_adminfile || !file_exists(R_P . $db_adminfile)) { $db_adminfile = 'example.com'; }
      if (!rename(R_P . $db_adminfile, R_P . $config['adminfile'])) { $config['adminfile'] = $db_adminfile; }

template/admin/target265.example.com:833
  <td class="td2"><input name="config[adminfile]" value="$db_adminfile" class="input input_wa" /></td>
```

## 前台渲染点（拿到这些页面的 HTML 就等于拿到入口名）
| 文件 | 行 | 输出 |
|---|---|---|
| `template/wind/header_login.htm` | 323 | `<a href="$db_adminfile" target="_blank">系统设置</a>` |
| `template/wind/read.htm` | 719-721 | `{$db_adminfile}?adminjob=superdel&admintype=article&action=deltpc&fid=-1&ifkeep=1&authorid=X&lines=30` |
| `template/wind/read_addfloor.htm` | 211-213 | 同上 + `adminjob=thread&action=replylist` / `adminjob=usermanage` |
| `template/wind/ajax.htm` | 314 | `<li><a href="$db_adminfile" target="_blank">系统设置</a></li>` |
| `mode/o/template/header.htm` | 287 | 同上 |
| `mode/area/template/m_header.htm` | 333 | 同上 |
| `u/themes/default/space_model.htm` | 248-252 | `{$db_adminfile}?adminjob=usermanage&...` |
| `template/wind/searcher_headbar.htm` | 47 | `| <a href="$db_adminfile" target="_blank">系统设置</a>` |

## 触发条件
`CkInArray($windid,$manager) || $SYSTEM['allowadmincp']`
- `$manager` = 版块版主列表（`pw_forums.manager`）
- `$SYSTEM['allowadmincp']` = `pw_permission` 表 `rkey='allowadmincp'` 且 `gid=<当前用户组>`
- 默认带该权限的组：`install_wind.sql:3166/3266/3364` → gid=3/4/5

## 操作步骤
1. 拿一个 allowadmincp=1 的用户组账号（或任一版块版主账号）登录 `/bbs/login.php`
2. `GET /bbs/example.com`、`GET /bbs/read.php?tid=<任一帖>`、`GET /bbs/u.php?uid=<自己>`
3. 正则 `href="([^"]+)\?adminjob=` / `>系统设置<` → 捕获组即后台入口文件名
4. 直接访问 `https://host/bbs/<入口名>` 进后台（前提：根目录无 install.php）

## 旁路（无管理权限时）
- `data/bbscache/config.php` 直读源码：nginx 交给 PHP-FPM 执行 → 返回 200/0B，拿不到值（本目标实测）
- 任何 LFI / SQL 注入读 `pw_config` 表 → 直接取 `db_adminfile` 值
- 备份变形 `config.php.bak/.php~/.txt` 一般 404

## 反面教训
无脑枚举根目录文件名上限极低：本目标 1459 个候选名（1433 定向字典 + 26 补充）+ `/admin/` 目录 28 个候选全部落空，
而入口名其实一直躺在 `pw_config.db_adminfile` 里。**先查配置项，再枚举。**