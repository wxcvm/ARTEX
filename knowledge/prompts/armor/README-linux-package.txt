ARTEX (Linux amd64)

目录
  artex               主程序（linux/amd64，已内嵌前端）
  start.sh            守护启动脚本（底层入口）
  start-artex.sh      推荐入口：导出 KB/工具路径与 venv 后再启动
  install-deps.sh     一键装系统包、初始化 PostgreSQL、装 Python 工具、生成 config.json
  config.example.json 配置模板
  requirements.txt    Python 工具清单
  skills/             技能库
  redteam/            知识库（knowledge.db + pocs/ + skills/ + toolkit/）

首次使用
  1. chmod +x *.sh artex
  2. ./install-deps.sh     # 需要 root/sudo 与联网
  3. ./start-artex.sh
  4. 打开 http://127.0.0.1:8787 ，首次进 /setup 设密码

默认配置
  install-deps.sh 只建库、写密码；破甲与拦截的默认值由 start-artex.sh
  在面板就绪后应用一次（控制台会打印「已应用部署者默认配置（破甲开关全开 /
  预设 4 套 / 拦截规则 8 条）」）。之后你在面板里改的设置持久保存，重启不覆盖。
  想恢复出厂默认：delete from settings where key='defaults.deployment_state' 后重启。

说明
  - config.json 由安装脚本生成，数据库密码写在 .pgpass.generated（0600）
  - knowledge_dir 指向本包 redteam，绝对路径
  - toolkit 内是 Linux 版侦察/扫描/隧道工具；冰蝎、哥斯拉为 jar，需 JDK（脚本会装）
  - 一定要用 start-artex.sh 或 start.sh 启动，别直接跑 ./artex
  - 若从 Windows 复制过来，注意执行 chmod +x；解压工具可能丢失可执行位
