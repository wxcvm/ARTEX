# 武器库清单

三个发行包自带 2.3 GB 成品工具。**二进制不入库**（体积 + 杀软必报），这里只保留清单：
名字、大小、类别、依赖、是否常被 AV 标记、在原包里的相对路径。

| 项 | 值 |
|---|---|
| 单件工具 | 197 |
| 运行时/大目录（聚合） | 8 |
| 完整数据 | `manifest.json` |

## 单件工具（按类别）

### scanner（117）

| 名称 | 大小 | 依赖 | AV | 来源 | 原包路径 |
|---|---|---|---|---|---|
| `nuclei.exe` | 170.8 MB | none |  | 0.4.2-win | `nuclei.exe` |
| `katana` | 63.9 MB | none |  | lw-linux-amd64 | `katana/katana` |
| `httpx` | 51.9 MB | none |  | lw-linux-amd64 | `httpx/httpx` |
| `httpx` | 51.8 MB | none |  | 0.4.2-win | `httpx` |
| `katana` | 47.8 MB | none |  | 0.4.2-win | `katana` |
| `subfinder.exe` | 38.7 MB | none |  | 0.4.2-win | `subfinder.exe` |
| `naabu.exe` | 33.1 MB | none |  | 0.4.2-win | `naabu.exe` |
| `naabu` | 32.5 MB | none |  | lw-linux-amd64 | `naabu/naabu` |
| `dnsx.exe` | 28.0 MB | none |  | 0.4.2-win | `dnsx.exe` |
| `subfinder` | 27.9 MB | none |  | lw-linux-amd64 | `subfinder/subfinder` |
| `dnsx` | 27.2 MB | none |  | lw-linux-amd64+friend-win | `dnsx/dnsx` |
| `ksubdomain` | 19.1 MB | none |  | lw-linux-amd64 | `ksubdomain/ksubdomain` |
| `ffuf.exe` | 15.1 MB | none |  | 0.4.2-win | `ffuf.exe` |
| `gospider.exe` | 14.4 MB | none |  | 0.4.2-win | `gospider.exe` |
| `gogo_linux_arm64` | 9.8 MB | none | Y | lw-linux-amd64 | `gogo/gogo_linux_arm64` |
| `fscan.exe` | 8.4 MB | none | Y | 0.4.2-win | `fscan.exe` |
| `fscan_windows_x64.exe` | 8.4 MB | none | Y | lw-linux-amd64+friend-win | `fscan/fscan_windows_x64.exe` |
| `fscan` | 8.1 MB | none | Y | lw-linux-amd64+friend-win | `fscan/fscan` |
| `fscan_linux_arm64` | 7.2 MB | none | Y | lw-linux-amd64+friend-win | `fscan/fscan_linux_arm64` |
| `nmap-os-db` | 4.8 MB | none |  | friend-win | `nmap/nmap-os-db` |
| `libgtk-win32-2.0-0.dll` | 4.7 MB | none |  | friend-win | `nmap/py2exe/libgtk-win32-2.0-0.dll` |
| `gogo` | 4.1 MB | none | Y | lw-linux-amd64 | `gogo/gogo` |
| `gogo` | 4.1 MB | none | Y | 0.4.2-win | `gogo` |
| `gogo.exe` | 4.1 MB | none | Y | lw-linux-amd64+0.4.2-win | `gogo.exe` |
| `gogo_windows_amd64.exe` | 4.1 MB | none | Y | lw-linux-amd64 | `gogo/gogo_windows_amd64.exe` |
| `nmap.exe` | 2.6 MB | none |  | friend-win | `nmap/nmap.exe` |
| `python27.dll` | 2.5 MB | python |  | friend-win | `nmap/python27.dll` |
| `nmap-service-probes` | 2.3 MB | none |  | friend-win | `nmap/nmap-service-probes` |
| `gtk._gtk.pyd` | 1.8 MB | none |  | friend-win | `nmap/py2exe/gtk._gtk.pyd` |
| `CRYPT32.dll` | 1.5 MB | none |  | friend-win | `nmap/py2exe/CRYPT32.dll` |
| `_ssl.pyd` | 1.3 MB | none |  | friend-win | `nmap/py2exe/_ssl.pyd` |
| `libcairo-2.dll` | 1.2 MB | none |  | friend-win | `nmap/py2exe/libcairo-2.dll` |
| `libeay32.dll` | 1.2 MB | none |  | friend-win | `nmap/libeay32.dll` |
| `libglib-2.0-0.dll` | 1.2 MB | none |  | friend-win | `nmap/py2exe/libglib-2.0-0.dll` |
| `libxml2-2.dll` | 1.2 MB | none |  | friend-win | `nmap/py2exe/libxml2-2.dll` |
| `libgio-2.0-0.dll` | 1.2 MB | none |  | friend-win | `nmap/py2exe/libgio-2.0-0.dll` |
| `library.zip` | 1.1 MB | none |  | friend-win | `nmap/py2exe/library.zip` |
| `_hashlib.pyd` | 993 KB | none |  | friend-win | `nmap/py2exe/_hashlib.pyd` |
| `nmap-services` | 976 KB | none |  | friend-win | `nmap/nmap-services` |
| `libgdk-win32-2.0-0.dll` | 910 KB | none |  | friend-win | `nmap/py2exe/libgdk-win32-2.0-0.dll` |

_（其余 77 条见 manifest.json）_

### tunnel（18）

| 名称 | 大小 | 依赖 | AV | 来源 | 原包路径 |
|---|---|---|---|---|---|
| `mihomo-windows-amd64.exe` | 31.3 MB | none |  | 0.4.2-win | `mihomo/mihomo-windows-amd64.exe` |
| `ligolo-proxy` | 20.4 MB | none | Y | 0.4.2-win | `ligolo-proxy` |
| `frps` | 19.4 MB | none |  | lw-linux-amd64+friend-win | `frp/frps` |
| `frps.exe` | 18.8 MB | none |  | 0.4.2-win | `frps.exe` |
| `frpc` | 15.8 MB | none |  | lw-linux-amd64+friend-win | `frp/frpc` |
| `frpc.exe` | 14.8 MB | none |  | 0.4.2-win | `frpc.exe` |
| `chisel.exe` | 11.4 MB | none | Y | 0.4.2-win | `chisel.exe` |
| `suo5` | 11.4 MB | none | Y | 0.4.2-win | `suo5` |
| `suo5-windows-amd64.exe` | 11.4 MB | none | Y | lw-linux-amd64 | `suo5-windows-amd64.exe` |
| `chisel` | 11.1 MB | none | Y | lw-linux-amd64+friend-win | `chisel/chisel` |
| `suo5-linux-amd64` | 10.9 MB | none | Y | lw-linux-amd64 | `suo5/suo5-linux-amd64` |
| `ligolo-agent.exe` | 7.0 MB | none | Y | 0.4.2-win | `ligolo-agent.exe` |
| `penelope.py` | 185 KB | python | Y | 0.4.2-win | `penelope.py` |
| `suo5.jsp` | 45 KB | none | Y | lw-linux-amd64+0.4.2-win | `suo5-payloads/suo5.jsp` |
| `suo5.aspx` | 42 KB | none | Y | lw-linux-amd64+0.4.2-win | `suo5-payloads/suo5.aspx` |
| `suo5.php` | 33 KB | none | Y | lw-linux-amd64+0.4.2-win | `suo5-payloads/suo5.php` |
| `suo5_keepalive.cmd` | 612 B | none | Y | lw-linux-amd64+friend-win | `suo5_keepalive.cmd` |
| `suo5_run_176.bat` | 342 B | none | Y | lw-linux-amd64+friend-win | `suo5_run_176.bat` |

### credential（7）

| 名称 | 大小 | 依赖 | AV | 来源 | 原包路径 |
|---|---|---|---|---|---|
| `laZagne.exe` | 9.7 MB | none | Y | 0.4.2-win | `laZagne.exe` |
| `mimikatz.exe` | 1.3 MB | none | Y | 0.4.2-win | `mimikatz.exe` |
| `pypykatz` | 105 KB | none | Y | 0.4.2-win | `pypykatz` |
| `nxc` | 105 KB | none | Y | 0.4.2-win | `nxc` |
| `pypykatz.cmd` | 95 B | none | Y | 0.4.2-win | `pypykatz.cmd` |
| `nxc.cmd` | 90 B | none | Y | 0.4.2-win | `nxc.cmd` |
| `impacket` | 38 B | none | Y | 0.4.2-win | `impacket` |

### exploit（24）

| 名称 | 大小 | 依赖 | AV | 来源 | 原包路径 |
|---|---|---|---|---|---|
| `AntSword` | 109.6 MB | none |  | lw-linux-amd64+friend-win | `AntSword/AntSword-Loader-v4.0.3-linux-x64/AntSword` |
| `godzilla.jar` | 30.4 MB | java | Y | lw-linux-amd64 | `Godzilla/godzilla.jar` |
| `icudtl.dat` | 9.7 MB | none |  | lw-linux-amd64+friend-win | `AntSword/AntSword-Loader-v4.0.3-linux-x64/icudtl.dat` |
| `resources.pak` | 8.3 MB | none |  | lw-linux-amd64+friend-win | `AntSword/AntSword-Loader-v4.0.3-linux-x64/resources.pak` |
| `libVkICD_mock_icd.so` | 6.9 MB | none |  | lw-linux-amd64+friend-win | `AntSword/AntSword-Loader-v4.0.3-linux-x64/libVkICD_mock_` |
| `libGLESv2.so` | 5.8 MB | none |  | lw-linux-amd64+friend-win | `AntSword/AntSword-Loader-v4.0.3-linux-x64/libGLESv2.so` |
| `libffmpeg.so` | 3.0 MB | none |  | lw-linux-amd64+friend-win | `AntSword/AntSword-Loader-v4.0.3-linux-x64/libffmpeg.so` |
| `logo.png` | 1.0 MB | none |  | lw-linux-amd64+friend-win | `AntSword/AntSword-Loader-v4.0.3-linux-x64/logo.png` |
| `v8_context_snapshot.bin` | 1015 KB | none |  | lw-linux-amd64+friend-win | `AntSword/AntSword-Loader-v4.0.3-linux-x64/v8_context_sna` |
| `snapshot_blob.bin` | 626 KB | none |  | lw-linux-amd64+friend-win | `AntSword/AntSword-Loader-v4.0.3-linux-x64/snapshot_blob.` |
| `libEGL.so` | 249 KB | none |  | lw-linux-amd64+friend-win | `AntSword/AntSword-Loader-v4.0.3-linux-x64/libEGL.so` |
| `chrome_200_percent.pak` | 244 KB | none |  | lw-linux-amd64+friend-win | `AntSword/AntSword-Loader-v4.0.3-linux-x64/chrome_200_per` |
| `chrome_100_percent.pak` | 163 KB | none |  | lw-linux-amd64+friend-win | `AntSword/AntSword-Loader-v4.0.3-linux-x64/chrome_100_per` |
| `natives_blob.bin` | 122 KB | none |  | lw-linux-amd64+friend-win | `AntSword/AntSword-Loader-v4.0.3-linux-x64/natives_blob.b` |
| `wxgroup.png` | 35 KB | none | Y | lw-linux-amd64 | `Godzilla/wxgroup.png` |
| `base64.png` | 26 KB | none | Y | lw-linux-amd64 | `Godzilla/base64.png` |
| `raw.png` | 21 KB | none | Y | lw-linux-amd64 | `Godzilla/raw.png` |
| `README.md` | 2 KB | none | Y | lw-linux-amd64 | `Godzilla/README.md` |
| `hashsumJar` | 1 KB | none | Y | lw-linux-amd64 | `Godzilla/hashsumJar` |
| `LICENSE` | 1 KB | none |  | lw-linux-amd64+friend-win | `AntSword/AntSword-Loader-v4.0.3-linux-x64/LICENSE` |
| `.gitignore` | 301 B | none | Y | lw-linux-amd64 | `Godzilla/.gitignore` |
| `application.config` | 258 B | none | Y | lw-linux-amd64 | `Godzilla/application.config` |
| `sqlmap.cmd` | 93 B | none |  | 0.4.2-win | `sqlmap.cmd` |
| `version` | 6 B | none |  | lw-linux-amd64+friend-win | `AntSword/AntSword-Loader-v4.0.3-linux-x64/version` |

### other（31）

| 名称 | 大小 | 依赖 | AV | 来源 | 原包路径 |
|---|---|---|---|---|---|
| `winPEASany.exe` | 11.2 MB | none | Y | 0.4.2-win | `peass/winPEASany.exe` |
| `winPEASx64.exe` | 11.2 MB | none | Y | 0.4.2-win | `peass/winPEASx64.exe` |
| `linpeas.sh` | 1.7 MB | none |  | 0.4.2-win | `peass/linpeas.sh` |
| `ncat.exe` | 425 KB | none |  | lw-linux-amd64 | `vps/ncat.exe` |
| `vps.sh` | 10 KB | none |  | lw-linux-amd64 | `vps/vps.sh` |
| `dga_jumper.py` | 10 KB | python |  | lw-linux-amd64+friend-win | `dga_jumper.py` |
| `fofa_mcp.py` | 9 KB | python |  | 0.4.2-win | `mcp/fofa_mcp.py` |
| `jumper.js` | 5 KB | none |  | lw-linux-amd64+friend-win | `jumper.js` |
| `p79.py` | 4 KB | python |  | lw-linux-amd64+friend-win | `p79.py` |
| `vps_proxy.py` | 4 KB | python |  | lw-linux-amd64 | `vps/vps_proxy.py` |
| `p78_client.py` | 3 KB | python |  | lw-linux-amd64+friend-win | `p78_client.py` |
| `VPS-│÷┐┌╦╡├≈.md` | 3 KB | none |  | lw-linux-amd64 | `vps/VPS-│÷┐┌╦╡├≈.md` |
| `VPS-出口说明.md` | 3 KB | none |  | lw-linux-amd64 | `vps/VPS-出口说明.md` |
| `VPS-│÷┐┌╦╡├≈.md.bak-fofa-20261010115636` | 3 KB | none |  | lw-linux-amd64 | `vps/VPS-│÷┐┌╦╡├≈.md.bak-fofa-20261010115636` |
| `dga_jumper.xml` | 1 KB | none |  | lw-linux-amd64+friend-win | `dga_jumper.xml` |
| `dga_jumper_nocmd.xml` | 1 KB | none |  | lw-linux-amd64+friend-win | `dga_jumper_nocmd.xml` |
| `id_rsa` | 1 KB | none |  | lw-linux-amd64 | `vps/id_rsa` |
| `╟σ╡Ñ.md` | 1 KB | none |  | lw-linux-amd64 | `╟σ╡Ñ.md` |
| `清单.md` | 1 KB | none |  | lw-linux-amd64+friend-win+friend-win | `清单.md` |
| `rt-proxy-on.cmd` | 1 KB | none |  | lw-linux-amd64 | `vps/rt-proxy-on.cmd` |
| `rt-proxy-on.cmd.bak-fofa-20261010115636` | 1 KB | none |  | lw-linux-amd64 | `vps/rt-proxy-on.cmd.bak-fofa-20261010115636` |
| `rt-proxy-off.cmd` | 962 B | none |  | lw-linux-amd64 | `vps/rt-proxy-off.cmd` |
| `ssh-proxy.cmd` | 891 B | none |  | lw-linux-amd64 | `vps/ssh-proxy.cmd` |
| `s5loop.cmd` | 542 B | none |  | lw-linux-amd64+friend-win | `s5loop.cmd` |
| `p78.out` | 497 B | none |  | lw-linux-amd64+friend-win | `p78.out` |
| `rt-tunnel-key` | 464 B | none |  | lw-linux-amd64 | `vps/rt-tunnel-key` |
| `id_rsa.153-newserver` | 411 B | none |  | lw-linux-amd64 | `vps/id_rsa.153-newserver` |
| `hidjumper.vbs` | 274 B | none |  | lw-linux-amd64+friend-win | `hidjumper.vbs` |
| `hid.vbs` | 226 B | none |  | lw-linux-amd64+friend-win | `hid.vbs` |
| `rt-tunnel-key.pub` | 100 B | none |  | lw-linux-amd64 | `vps/rt-tunnel-key.pub` |
| `run_jumper.cmd` | 82 B | none |  | lw-linux-amd64+friend-win | `run_jumper.cmd` |

## 运行时与大目录（聚合，不逐个列）

| 目录 | 文件数 | 大小 | 来源 |
|---|---|---|---|
| `Behinder/` | 84 | 270.8 MB | friend-win+lw-linux-amd64 |
| `pgsql/` | 156 | 146.9 MB | 0.4.2-win+friend-win |
| `oneforall/` | 284 | 107.2 MB | 0.4.2-win+lw-linux-amd64 |
| `jdk11/` | 201 | 76.4 MB | friend-win+lw-linux-amd64 |
| `python313/` | 33 | 20.3 MB | 0.4.2-win |
| `dirsearch/` | 459 | 2.3 MB | 0.4.2-win+friend-win+lw-linux-amd64 |
| `wordlists/` | 6 | 575 KB | 0.4.2-win+lw-linux-amd64 |
| `nmap/ncat/` | 1 | 425 KB | friend-win |

## 用的时候注意

1. **加排除**：`mimikatz` / `laZagne` / `winPEAS` / `suo5` / `behinder` 会被 Defender 直接删；先加目录排除再解压。
2. **别进 git**：仓库里没有这些文件；`.gitignore` 的 `*.exe`/`*.db` 也不会自动挡住它们，靠"不在仓库里"保证。
3. **Java**：老魏包 `redteam/toolkit/jdk11/`、0.4.2 包 `java/`、Friend 包 `java/`（JDK21）。
4. **Python**：impacket / nxc / pypykatz / sqlmap 在各自 venv 里，不是单文件。
5. **Linux 工具不能直接在 Windows 跑**：老魏包 toolkit 是 linux/amd64。
