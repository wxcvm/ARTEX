export const onboardingSteps = ["准备数据库", "设置登录密码", "连接 AI", "开始使用"] as const;

export const onboardingGuides = {
  database: {
    purpose: "数据库负责保存账号、任务和结果。桌面版可自动准备独立数据库，无需自己填写端口和数据库密码。",
    materials: "正常打开 ARTEX；若缺少 PostgreSQL 程序，需要联网安装。不要以管理员身份运行 ARTEX。",
    actions: "点击“自动准备本机数据库”，等待检查、初始化、启动、验证和保存完成。缺少程序时，按官方安装器指引安装 PostgreSQL 16，再返回重试。已有连接信息可在“已有数据库 / 高级设置”填写。",
    example: "示例：独立数据库使用本机地址 127.0.0.1，自动选择 15432–15442 端口；这些参数由程序保存，不需要照抄。",
    success: "看到“数据库已就绪”并进入登录密码设置。只发现服务或端口，不等于账号与数据库可用。",
    failure: "服务的数据目录缺失时可另建独立新库，不会修复或覆盖旧服务。权限错误请正常启动 ARTEX；端口占用请关闭占用程序或用高级连接。失败后重试，不要删除旧数据。",
  },
  account: {
    purpose: "这是进入 ARTEX 的登录密码，不是 PostgreSQL 密码，也不是 AI 平台的 API Key。已有账号不需要重新设置。",
    materials: "准备一个至少 8 位的登录密码。桌面自己启动的实例会自动取得一次性初始化令牌；浏览器或外部实例需要由你提供令牌。",
    actions: "填写登录密码两次，点击“设置密码并继续”。令牌未自动取得时，展开手动入口，打开启动日志，找到 first-time setup token 后复制其后的值。",
    example: "示例位置：Windows 桌面日志位于 %LOCALAPPDATA%\\ARTEX\\desktop-backend.log；Docker 使用 docker compose logs artex。示例路径不是令牌。",
    success: "成功登录并进入“连接 AI”。以后使用刚才设置的 ARTEX 登录密码。",
    failure: "两次密码不同请重新核对；令牌失效请读取当前启动实例的日志。无法确认初始化状态时先恢复后端，不要重置已有账号。",
  },
  model: {
    purpose: "AI 模型负责理解任务并生成执行计划。API Key 是模型服务提供的调用凭证，不是 ARTEX 登录密码。",
    materials: "从模型服务提供方获取协议、接口地址、准确模型名和 API Key；本机 Ollama 通常不需要 Key，但必须先启动其服务并下载模型。",
    actions: "选择协议，按提供方文档填写地址、模型名和 Key。点击“测试并启用”：先请求模型获得实际回复，再保存并设为当前模型。测试会产生一次真实模型请求，是否计费取决于提供方。",
    example: "仅为填写格式示例：OpenAI 兼容地址 https://api.example.com/v1、模型 example-model；不能直接拿这些示例连接。Ollama 本机地址通常为 http://127.0.0.1:11434。",
    success: "看到模型实际回复、耗时及“已保存并启用”。已有当前模型会显示“已配置，尚未在本次引导重测”，不会自动发送测试请求。",
    failure: "401/403 请检查 Key 和模型权限；404 请检查地址、协议和模型名；超时请检查服务或在高级设置填代理。保存或激活失败可直接重试，输入不会清空。",
  },
  start: {
    purpose: "先熟悉创建任务、查看进度和阅读结果的路线，再开始实际操作；知识库可以在没有模型时浏览。",
    materials: "运行任务需要可用模型及任务所需的目标、目标描述与执行环境。浏览知识库不需要启动任务。",
    actions: "打开“创建第一个任务”，填写名称、目标描述和希望得到的结果，核对任务设置后自行启动。在任务详情查看进度、会话、资产与报告；发现项可进入复核和导出。",
    example: "描述示例：整理指定测试环境的资产与配置问题；目标示例：生成可复核的问题清单。示例不会自动填入或运行。",
    success: "任务实际执行后有可查看的进度和结果；零发现、未完成和失败是不同结果，不能仅凭任务已创建判断成功。",
    failure: "任务无法开始先检查模型是否已启用、服务是否在线，再查看任务错误与日志。报告尚未生成时先确认任务状态；指南中的入口可以随时重新打开。",
  },
} as const;

export const onboardingGuidesEn = {
  "database": {
    "purpose": "The database stores accounts, tasks and results. Desktop ARTEX can prepare an independent database without manual port or password entry.",
    "materials": "Open ARTEX normally. If PostgreSQL is missing, install it using an internet connection. Do not run ARTEX as administrator.",
    "actions": "Choose automatic local database setup and wait for inspection, initialization, startup, verification and saving. If PostgreSQL is missing, install PostgreSQL 16 as directed, then retry. Existing connection details can be entered in Existing database / Advanced settings.",
    "example": "Example: the independent database uses 127.0.0.1 and automatically selects a port in 15432–15442. ARTEX saves these values; do not copy them manually.",
    "success": "Database ready appears and password setup opens. Detecting a service or port alone does not establish database or account access.",
    "failure": "If an old service has no data directory, create a separate new database; this will not repair or overwrite the old service. For permission errors launch ARTEX normally. For occupied ports close the conflicting program or use an advanced connection. Retry without deleting old data."
  },
  "account": {
    "purpose": "This password is for ARTEX login, not PostgreSQL or the AI API key. Existing accounts do not need to be reinitialized.",
    "materials": "Prepare a login password of at least eight characters. Desktop-managed instances obtain a one-time setup token automatically; browser or external instances require a token.",
    "actions": "Enter the password twice and choose Set password and continue. If the token was not obtained automatically, open manual entry, find first-time setup token in startup logs and copy its value.",
    "example": "Example Windows desktop log: %LOCALAPPDATA%\\ARTEX\\desktop-backend.log. For Docker use docker compose logs artex. The example path is not a token.",
    "success": "Login succeeds and Connect AI opens. Use the new ARTEX password for future logins.",
    "failure": "Check mismatched passwords. If the token expired, read logs from the current instance. If initialization state is unknown, restore the backend first; do not reset an existing account."
  },
  "model": {
    "purpose": "The AI model understands tasks and generates plans. Its API key is issued by the model provider, not ARTEX.",
    "materials": "Obtain protocol, endpoint, exact model ID and API key from the provider. Local Ollama usually needs no key but must be running with a downloaded model.",
    "actions": "Choose protocol and enter endpoint, model and key following provider documentation. Test and enable sends a real request, then saves and activates the profile after a valid reply. Provider charges may apply.",
    "example": "Format example only: OpenAI-compatible URL https://api.example.com/v1, model example-model. These examples are not live credentials. Ollama commonly uses http://127.0.0.1:11434.",
    "success": "Actual model reply, latency and Saved and enabled appear. An existing active profile is not automatically retested in this guide.",
    "failure": "401/403: check key and permissions. 404: check URL, protocol and model ID. Timeout: check service or configure proxy. Retry save or activation failures without clearing the input."
  },
  "start": {
    "purpose": "Learn task creation, progress and results before operating. Browse the knowledge base without a model.",
    "materials": "Tasks need a working model plus target, task description and execution environment. Knowledge browsing needs no task.",
    "actions": "Open Create first task, enter name, description and desired result, review settings and start yourself. View progress, sessions, assets and reports in task details; review and export findings.",
    "example": "Example description: organize assets and configuration issues in a designated test environment. Example goal: a reviewable list of issues. Examples are not entered or executed automatically.",
    "success": "Actual execution produces visible progress and results. Zero findings, incomplete and failed are different outcomes; task creation alone is not success.",
    "failure": "If a task cannot start, check active model and service availability, then task errors and logs. If a report is missing, check task state. Reopen the guide anytime."
  }
} as const;

export function modelErrorGuide(message: string): string {
  if (/401|403|unauthorized|authentication|api.?key/i.test(message)) return "检查 API Key 是否有效、额度及模型权限；不要填写 ARTEX 登录密码。";
  if (/404|not found/i.test(message)) return "核对协议、接口地址和模型名；地址应符合模型提供方文档。";
  if (/timeout|timed out|network|fetch|connect/i.test(message)) return "确认模型服务可访问；本机 Ollama 需要先启动，必要时展开高级设置配置代理。";
  return "保留当前输入重试；可打开高级模型设置核对提供方要求。";
}

export async function testAndEnableModel<T>(actions: {
  test: () => Promise<{ ok: boolean; error?: string; reply?: string; latency_ms?: number }>;
  save: () => Promise<T>;
  activate: (saved: T) => Promise<unknown>;
}) {
  const result = await actions.test();
  if (!result.ok || !result.reply?.trim()) throw new Error(result.error || "模型没有返回有效回复，未保存或启用。");
  let saved: T;
  try { saved = await actions.save(); }
  catch (error) { throw new Error(`模型测试已通过，但保存失败：${(error as Error).message}`); }
  try { await actions.activate(saved); }
  catch (error) { throw new Error(`模型测试已通过且配置已保存，但启用失败：${(error as Error).message}。可重试启用。`); }
  return result;
}
