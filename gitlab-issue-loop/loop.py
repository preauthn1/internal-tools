#!/usr/bin/env python3
"""GitLab issue → MR 自動開發迴圈。

只處理指派給 bot 的 issue。狀態以 label 記在 GitLab 上（人類看得到、也能手動改），分兩組：

triage label（與專案共用的五個，同時只留一個）：
  無／needs-triage／ready-for-agent  → 釐清（triage）
  needs-info       AI 有疑問，等人類回覆；人類回覆後重新釐清
  ready-for-human  AI 判斷不適合自動處理；人類改回 ready-for-agent 才會再接
  wontfix          不處理
執行狀態（迴圈專用）：
  ai::doing    開發中
  ai::review   MR 已開，等人類審閱
  ai::blocked  開發或測試失敗；人類留言後重新釐清

用法：
  ./loop.py once                 處理所有待辦 issue 一輪
  ./loop.py watch [--interval N] 每 N 秒輪詢一次（touch STOP 停止）
  ./loop.py plan                 只列出會做什麼，不動任何東西
  ./loop.py issue group/proj#12  只處理指定 issue
"""
import argparse
import datetime
import fcntl
import json
import os
import re
import shutil
import subprocess
import sys
import time
import urllib.error
import urllib.parse
import urllib.request

ROOT = os.path.dirname(os.path.abspath(__file__))
WORK = os.path.join(ROOT, "work")
RUNS = os.path.join(ROOT, "runs")
MARK = "<!-- ai-loop -->"
TRIAGE_LABELS = ["needs-triage", "needs-info", "ready-for-agent", "ready-for-human", "wontfix"]
STATUS_LABELS = ["ai::doing", "ai::review", "ai::blocked", "ai::asking"]  # ai::asking 是舊版留下的，改用 needs-info
# issue 內容是不可信輸入：不讓 AI 有能力把東西送出本機
DEFAULT_DENY = ["Bash(git push:*)", "Bash(git remote:*)", "Bash(curl:*)", "Bash(wget:*)",
                "Bash(glab:*)", "Bash(ssh:*)", "WebFetch", "WebSearch"]

TRIAGE_SCHEMA = {
    "type": "object",
    "properties": {
        "status": {"enum": ["need_info", "ready", "reject"]},
        "comment": {"type": "string", "description": "要貼到 issue 的 Markdown 回覆"},
        "plan": {"type": "string", "description": "ready 時的實作計畫，交給開發階段"},
    },
    "required": ["status", "comment", "plan"],
}
DEVELOP_SCHEMA = {
    "type": "object",
    "properties": {
        "status": {"enum": ["done", "blocked"]},
        "commit_title": {"type": "string"},
        "summary": {"type": "string", "description": "MR 說明：改了什麼、為什麼"},
        "test_notes": {"type": "string", "description": "如何驗證、跑了哪些測試"},
        "blocked_reason": {"type": "string"},
        "unmet": {"type": "string", "description": "數字門檻沒達到時寫「門檻 X，實際 Y」和試過的做法；都達標或沒有數字門檻時是空字串"},
    },
    "required": ["status", "commit_title", "summary", "test_notes", "blocked_reason", "unmet"],
}


def log(msg):
    print(f"[{datetime.datetime.now():%H:%M:%S}] {msg}", flush=True)


# ---------- GitLab API ----------

class GitLab:
    def __init__(self, url, token):
        self.base = url.rstrip("/") + "/api/v4"
        self.token = [REDACTED]

    def req(self, method, path, data=None, params=None):
        url = self.base + path
        if params:
            url += "?" + urllib.parse.urlencode(params)
        body = json.dumps(data).encode() if data is not None else None
        r = urllib.request.Request(url, data=body, method=method)
        r.add_header("PRIVATE-TOKEN", self.token)
        r.add_header("Content-Type", "application/json")
        try:
            with urllib.request.urlopen(r, timeout=60) as resp:
                return json.loads(resp.read() or b"null"), resp.headers
        except urllib.error.HTTPError as e:
            raise RuntimeError(f"GitLab {method} {path} → {e.code}: {e.read()[:500]!r}")

    def get(self, path, **params):
        return self.req("GET", path, params=params)[0]

    def get_all(self, path, **params):
        out, page = [], 1
        while True:
            items, headers = self.req("GET", path, params={**params, "per_page": 100, "page": page})
            out += items
            if not headers.get("X-Next-Page"):
                return out
            page = int(headers["X-Next-Page"])

    def post(self, path, data):
        return self.req("POST", path, data=data)[0]

    def put(self, path, data):
        return self.req("PUT", path, data=data)[0]


def pid(path):
    return "/projects/" + urllib.parse.quote(path, safe="")


# ---------- 工具 ----------

def sh(cmd, cwd, timeout=None, env=None, check=True):
    p = subprocess.run(cmd, cwd=cwd, shell=isinstance(cmd, str), capture_output=True,
                       text=True, timeout=timeout, env=env)
    if check and p.returncode != 0:
        raise RuntimeError(f"指令失敗：{cmd}\n{p.stdout[-2000:]}\n{p.stderr[-2000:]}")
    return p


def missing_process(result):
    """develop.md 要求的流程留下的痕跡：summary 要有「## 自我審閱」、test_notes 要有「## TDD」。回傳缺了哪些（空的就是都有）。"""
    missing = []
    if "## 自我審閱" not in result.get("summary", ""):
        missing.append("- summary 沒有「## 自我審閱」：要用 code-review skill 審閱（Standards 和 Spec 兩份結果都等回來、範圍內的修掉），再寫這一節")
    if "## TDD" not in result.get("test_notes", ""):
        missing.append("- test_notes 沒有「## TDD」：要依順序列出用 tdd skill 寫的每個測試（先失敗再通過），或寫明為什麼不適用")
    return missing


# issue 描述裡的等待：「Blocked by: #a」是需求依賴，「排在後面：#b（原因）」是排程順序（避免同時改同一塊程式）。
# 兩種迴圈都一樣處理：列的 issue 還沒關閉（MR 還沒合併）就先不做。只認同專案的 #N；「## Blocked by」小節不讀
WAIT_RE = re.compile(r"^(Blocked by|排在後面)\s*[:：]\s*(.+)$", re.M)
WAIT_KIND = {"Blocked by": "需求", "排在後面": "排程"}


def waits_on(description):
    """回傳 [(種類, 編號)]，種類是「需求」或「排程」。"""
    out = []
    for kind, rest in WAIT_RE.findall(description or ""):
        for n in re.findall(r"(?<![\w/])#(\d+)", rest):
            out.append((WAIT_KIND[kind], int(n)))
    return out


def notify(cfg, text):
    """Telegram 私訊通知（config 的 notify：telegram_token_file、chat_id）；沒設定或送不出去只記 log，不影響迴圈。"""
    n = cfg.get("notify") or {}
    path = os.path.expanduser(n.get("telegram_token_file", ""))
    if not (path and n.get("chat_id") and os.path.isfile(path)):
        return
    try:
        with open(path) as f:
            token = [REDACTED]
        data = urllib.parse.urlencode({"chat_id": n["chat_id"], "text": text, "disable_web_page_preview": "true"}).encode()
        urllib.request.urlopen(f"https://example.invalid/bot{token}/sendMessage", data=data, timeout=15).read()
        log(f"  📨 已通知：{text.splitlines()[0][:80]}")   # 送成功也記一筆，之後查得到有沒有送出
    except Exception as e:  # 通知失敗不能讓迴圈停下來；錯誤訊息不含 token
        log(f"  （通知送不出去：{type(e).__name__}）")


def tail(text, n=150):
    lines = text.splitlines()
    return "\n".join(lines[-n:])


def save_run(name, data):
    os.makedirs(RUNS, exist_ok=True)
    ts = datetime.datetime.now().strftime("%Y%m%d-%H%M%S")
    with open(os.path.join(RUNS, f"{ts}-{name}.json"), "w") as f:
        json.dump(data, f, ensure_ascii=False, indent=2)


def render(template, **kv):
    with open(os.path.join(ROOT, "prompts", template)) as f:
        text = f.read()
    for k, v in kv.items():
        text = text.replace("{{" + k + "}}", v)
    return text


def agent_env():
    # 不把任何 GitLab token 傳給 AI：所有對外動作（留言、push、開 MR）都由本腳本做
    env = {k: v for k, v in os.environ.items() if "GITLAB" not in k and "TOKEN" not in k}
    # 不讓 AI 把指令或子代理丟到背景：背景工作結束時會再觸發一輪回覆，最後一輪沒有結構化輸出，迴圈會誤判失敗
    # 量測這種長指令改在前景跑，所以把單一指令的逾時上限放寬到 30 分鐘（預設只有 10 分鐘）
    env["CLAUDE_CODE_DISABLE_BACKGROUND_TASKS"] = "1"
    env.setdefault("BASH_MAX_TIMEOUT_MS", "1800000")
    return env


# 這張 issue 每次呼叫 claude 的執行數字（handle 開始時清空），附在 MR 說明和 issue 留言
RUN_STATS = []


def record_stats(name, out, seconds):
    denials = [d.get("tool_input", {}).get("command") or d.get("tool_name", "?") for d in out.get("permission_denials") or []]
    RUN_STATS.append({"name": name.rsplit("-", 1)[-1], "turns": out.get("num_turns") or 0,
                      "cost": out.get("total_cost_usd") or 0.0, "seconds": seconds, "denials": denials})


def stats_md():
    """執行數字的 Markdown；被擋超過 5 次或某一步超過 200 輪時，最上面加警告。"""
    if not RUN_STATS:
        return ""
    denials = [d for s in RUN_STATS for d in s["denials"]]
    rows = [f"| {s['name']} | {s['turns']} | {s['seconds'] / 60:.0f} 分 | ${s['cost']:.2f} | {len(s['denials'])} |" for s in RUN_STATS]
    warn = []
    if len(denials) > 5:
        warn.append(f"權限擋下 {len(denials)} 次")
    if any(s["turns"] > 200 for s in RUN_STATS):
        warn.append("有一步超過 200 輪")
    out = (f"> ⚠️ {'、'.join(warn)}：可能要調整迴圈的權限或拆小 issue\n\n" if warn else "") + \
        "## 執行紀錄\n\n| 步驟 | 輪數 | 時間 | 估算費用 | 權限擋下 |\n|---|---|---|---|---|\n" + "\n".join(rows) + \
        f"\n| **合計** | {sum(s['turns'] for s in RUN_STATS)} | {sum(s['seconds'] for s in RUN_STATS) / 60:.0f} 分 | " \
        f"${sum(s['cost'] for s in RUN_STATS):.2f} | {len(denials)} |\n\n費用是依 API 牌價的估算；用訂閱登入時不另外收費。"
    if denials:
        uniq = list(dict.fromkeys(str(d)[:120].replace("`", "'") for d in denials))
        out += "\n\n<details><summary>被擋下的指令</summary>\n\n" + "\n".join(f"- `{d}`" for d in uniq[:30]) + "\n\n</details>"
    return out


FINAL_PROMPT = "請照 JSON schema 再交一次這次工作的最終結果（內容和你剛才交的相同；沒交過就現在整理），不要再做其他事。"


def claude(cfg, prompt, cwd, schema, tools, name, resume=None, deny_extra=(), retry=True):
    cmd = ["claude", "-p", "--model", cfg.get("model", "opus"),
           "--output-format", "json", "--json-schema", json.dumps(schema),
           "--permission-mode", "acceptEdits",
           "--allowedTools", ",".join(tools),
           "--disallowedTools", ",".join(list(cfg.get("deny_tools", DEFAULT_DENY)) + list(deny_extra)),
           "--max-turns", str(cfg.get("max_turns", 80))]
    if resume:
        cmd += ["--resume", resume]
    t0 = time.time()
    p = subprocess.run(cmd, cwd=cwd, input=prompt, capture_output=True, text=True,
                       timeout=cfg.get("agent_timeout", 3600), env=agent_env())
    try:
        out = json.loads(p.stdout)
    except json.JSONDecodeError:
        save_run(name, {"stdout": p.stdout, "stderr": p.stderr})
        raise RuntimeError(f"claude 輸出無法解析：{p.stderr[-1000:]}")
    save_run(name, out)
    record_stats(name, out, time.time() - t0)
    if not out.get("is_error") and not out.get("structured_output") and retry and out.get("session_id"):
        # 正常結束但最後一輪沒有結構化輸出（例如背景工作結束後又回了一段文字）：接著同一個 session 要一次
        log("  沒有結構化輸出，接續 session 要求重交")
        return claude(cfg, FINAL_PROMPT, cwd, schema, tools, name + "-final", resume=out["session_id"],
                      deny_extra=deny_extra, retry=False)
    if out.get("is_error") or not out.get("structured_output"):
        raise RuntimeError(f"claude 執行失敗：{out.get('subtype')} {str(out.get('result'))[:500]}")
    return out["structured_output"], out.get("session_id")


# ---------- 單一專案 ----------

class Project:
    def __init__(self, gl, cfg, pcfg):
        self.gl, self.cfg, self.p = gl, cfg, pcfg
        self.path = pcfg["path"]
        self.api = pid(self.path)
        self.bot = cfg["bot_username"]
        self.slug = self.path.replace("/", "__")
        self.repo = os.path.join(WORK, self.slug, "repo")
        self._info = None

    @property
    def info(self):
        if self._info is None:
            self._info = self.gl.get(self.api)
        return self._info

    @property
    def default_branch(self):
        return self.p.get("target_branch") or example.invalid["default_branch"]

    # --- git ---
    def git(self, *args, cwd):
        # 走 SSH；有設定 ssh_key 就用 bot 自己的金鑰，不用執行者的
        env = dict(os.environ)
        key = self.cfg.get("ssh_key")
        if key:
            env["GIT_SSH_COMMAND"] = f"ssh -i {os.path.expanduser(key)} -o IdentitiesOnly=yes -o BatchMode=yes"
        return sh(["git", *args], cwd=cwd, env=env)

    def sync_repo(self):
        url = self.p.get("ssh_url") or example.invalid["ssh_url_to_repo"]
        if not os.path.isdir(os.path.join(self.repo, ".git")):
            os.makedirs(os.path.dirname(self.repo), exist_ok=True)
            self.git("clone", url, self.repo, cwd=ROOT)
            sh(["git", "config", "user.name", self.cfg.get("git_name", self.bot)], cwd=self.repo)
            sh(["git", "config", "user.email", self.cfg.get("git_email", f"{self.bot}@users.noreply")], cwd=self.repo)
        self.git("fetch", "--prune", "origin", cwd=self.repo)
        sh(["git", "checkout", "-q", "--detach", f"origin/{self.default_branch}"], cwd=self.repo)

    def fresh_worktree(self, iid):
        """所有 issue 共用同一個 worktree：換工單時切分支、清掉未追蹤檔案，
        但保留建置快取（keep_dirs，預設 .build、node_modules），之後只要增量重編。"""
        root = os.path.join(WORK, self.slug)
        wt = os.path.join(root, "worktree")
        branch = f"ai/issue-{iid}"
        # 舊版每個 issue 一個 worktree（issue-N），各帶一份建置快取，留著只佔空間
        for name in sorted(os.listdir(root)):
            old = os.path.join(root, name)
            if name.startswith("issue-") and os.path.isdir(old):
                sh(["git", "worktree", "remove", "--force", old], cwd=self.repo, check=False)
                shutil.rmtree(old, ignore_errors=True)
        sh(["git", "worktree", "prune"], cwd=self.repo)
        # 分支已推過（MR 修改要求）就接著原本的 commit 改，不從預設分支重來
        existing = sh(["git", "rev-parse", "--verify", "--quiet", f"refs/remotes/origin/{branch}"],
                      cwd=self.repo, check=False).returncode == 0
        base = f"origin/{branch}" if existing else f"origin/{self.default_branch}"
        if os.path.isfile(os.path.join(wt, ".git")):
            keep = self.p.get("keep_dirs", [".build", "node_modules"])
            sh(["git", "reset", "-q", "--hard"], cwd=wt)   # 也清掉上次中斷留下的合併狀態
            sh(["git", "checkout", "-q", "-f", "-B", branch, base], cwd=wt)
            sh(["git", "clean", "-q", "-fdx", *[f"--exclude={k}" for k in keep]], cwd=wt)
        else:
            shutil.rmtree(wt, ignore_errors=True)
            sh(["git", "worktree", "add", "-B", branch, wt, base], cwd=self.repo)
        return wt, branch, existing

    # --- issue ---
    def issue_url(self, iid):
        return f"{self.cfg['gitlab_url']}/{self.path}/-/issues/{iid}"

    def issues(self):
        """指派給 bot 的 issue，編號小的先做。"""
        items = self.gl.get_all(self.api + "/issues", state="opened", assignee_username=self.bot)
        return sorted(items, key=lambda i: i["iid"])

    def waiting(self, issue, states):
        """還沒關閉的等待對象，例如 ["#27（需求）"]；states 是這一輪查過的 issue 狀態（快取）。"""
        out = []
        for kind, n in waits_on(issue.get("description")):
            if n == issue["iid"]:
                continue
            if n not in states:
                try:
                    states[n] = self.gl.get(f"{self.api}/issues/{n}")["state"]
                except RuntimeError:
                    states[n] = "missing"  # 找不到的編號不擋
            if states[n] == "opened":
                out.append(f"#{n}（{kind}）")
        return out

    def notes(self, iid):
        notes = self.gl.get_all(f"{self.api}/issues/{iid}/notes", sort="asc", order_by="created_at")
        return [n for n in notes if not n.get("system")]

    def comment(self, iid, body):
        self.gl.post(f"{self.api}/issues/{iid}/notes", {"body": body + "\n\n" + MARK})

    def set_labels(self, issue, triage=None, status=None):
        """triage：換成這個 triage label（None 表示不動）；status：執行狀態，None 表示清掉。"""
        drop = STATUS_LABELS + (TRIAGE_LABELS if triage else [])
        labels = [l for l in issue["labels"] if l not in drop] + [l for l in (triage, status) if l]
        self.gl.put(f"{self.api}/issues/{issue['iid']}", {"labels": ",".join(labels)})
        issue["labels"] = labels

    def state(self, issue):
        found = [l for l in STATUS_LABELS + TRIAGE_LABELS if l in issue["labels"]]
        return found[0] if found else None

    def human_replied(self, issue):
        notes = self.notes(issue["iid"])
        return bool(notes) and notes[-1]["author"]["username"] != self.bot

    def next_action(self, issue):
        """回傳要做的事：triage / None。"""
        st = self.state(issue)
        if st in ("ai::doing", "ai::review", "ready-for-human", "wontfix"):
            return None
        if st in ("ai::blocked", "ai::asking", "needs-info"):
            return "triage" if self.human_replied(issue) else None
        return "triage"  # 無 label、needs-triage、ready-for-agent

    def transcript(self, issue):
        parts = [f"# #{issue['iid']} {issue['title']}",
                 f"作者：@{issue['author']['username']}　label：{', '.join(issue['labels']) or '無'}",
                 "", issue.get("description") or "（無描述）"]
        for n in self.notes(issue["iid"]):
            who = "你（AI bot）" if n["author"]["username"] == self.bot else "@" + n["author"]["username"]
            parts += ["", f"---\n## 留言 by {who}（{n['created_at'][:16]}）", n["body"].replace(MARK, "").strip()]
        return "\n".join(parts)

    # --- 流程 ---
    def handle(self, issue):
        iid = issue["iid"]
        tag = f"{self.path}#{iid}"
        log(f"▶ {tag} {issue['title']}")
        RUN_STATS.clear()
        self.sync_repo()
        result, _ = claude(
            self.cfg,
            render("triage.md", ISSUE=self.transcript(issue), BOT=self.bot,
                   TEST_CMD=self.p.get("test_cmd", "（未設定）")),
            cwd=self.repo, schema=TRIAGE_SCHEMA, deny_extra=self.p.get("deny_tools_extra", []), name=f"{self.slug}-{iid}-triage",
            tools=["Read", "Glob", "Grep", "Bash(git log:*)", "Bash(git show:*)", "Bash(git diff:*)"])
        log(f"  釐清結果：{result['status']}")

        if result["status"] == "need_info":
            self.comment(iid, result["comment"])
            self.set_labels(issue, triage="needs-info")
            notify(self.cfg, f"❓ {tag} 有問題要你回答：{self.issue_url(iid)}")
            return
        if result["status"] == "reject":
            self.comment(iid, result["comment"])
            self.set_labels(issue, triage="ready-for-human")
            notify(self.cfg, f"🙋 {tag} 不適合自動處理，交給人：{self.issue_url(iid)}")
            return

        self.comment(iid, result["comment"])
        self.set_labels(issue, triage="ready-for-agent", status="ai::doing")
        try:
            self.develop(issue, result["plan"])
        except Exception as e:
            log(f"  ✖ {e}")
            self.comment(iid, f"⚠️ 自動開發中斷，需要人類協助：\n\n```\n{tail(str(e), 40)}\n```\n\n"
                              "補充說明後留言，我會重新評估。\n\n" + stats_md())
            self.set_labels(issue, status="ai::blocked")
            notify(self.cfg, f"⚠️ {tag} 卡住了：{tail(str(e), 3)[:300]}\n{self.issue_url(iid)}")

    def develop(self, issue, plan):
        iid = issue["iid"]
        wt, branch, existing = self.fresh_worktree(iid)
        if self.p.get("setup_cmd"):
            sh(self.p["setup_cmd"], cwd=wt, timeout=self.p.get("test_timeout", 1800))
        # Skill、Agent：develop.md 要求用 tdd、code-review 兩個 skill，code-review 會開子代理平行審閱
        tools = ["Read", "Edit", "Write", "Glob", "Grep", "Skill", "Agent",
                 "Bash(git status:*)", "Bash(git diff:*)", "Bash(git log:*)"] + self.p.get("agent_bash", [])
        result, session = claude(
            self.cfg,
            render("develop.md", ISSUE=self.transcript(issue), PLAN=plan,
                   BASE=(f"這個分支已經有先前推上去、MR 審閱中的 commit（`git log origin/{self.default_branch}..HEAD` 可以看），"
                         "請接在上面修改，不要重做已完成的部分。summary 會整段取代 MR 說明，所以要描述整個分支相對於"
                         f" {self.default_branch} 的全部變更，不只這次的修改。") if existing else
                        f"分支從 {self.default_branch} 的最新 commit 開始。",
                   TEST_CMD=self.p.get("test_cmd", "（未設定，請自行找出專案的測試方式）"),
                   DEFAULT_BRANCH=self.default_branch),
            cwd=wt, schema=DEVELOP_SCHEMA, tools=tools, deny_extra=self.p.get("deny_tools_extra", []), name=f"{self.slug}-{iid}-develop")

        # 流程檢查：沒做自我審閱或沒交代 TDD，退回補做一次；還是沒有就不開 MR
        if result["status"] != "blocked":
            missing = missing_process(result)
            if missing:
                log("  流程不完整，退回補做：" + "；".join(m.split("：")[0][2:] for m in missing))
                result, session = claude(
                    self.cfg, render("process.md", MISSING="\n".join(missing)), cwd=wt, schema=DEVELOP_SCHEMA,
                    tools=tools, deny_extra=self.p.get("deny_tools_extra", []), name=f"{self.slug}-{iid}-process", resume=session)
                missing = missing_process(result) if result["status"] != "blocked" else []
                if missing:
                    raise RuntimeError("補做之後開發流程仍不完整：\n" + "\n".join(missing))

        rounds = 0
        while True:
            if result["status"] == "blocked":
                raise RuntimeError("AI 回報無法完成：" + result["blocked_reason"])
            if not sh(["git", "status", "--porcelain"], cwd=wt).stdout.strip():
                raise RuntimeError("AI 回報完成，但沒有任何檔案變更")
            ok, output = self.verify(wt)
            if ok:
                break
            rounds += 1
            if rounds > self.cfg.get("max_fix_rounds", 2):
                raise RuntimeError(f"修了 {rounds - 1} 輪測試仍失敗：\n{tail(output, 40)}")
            log(f"  測試失敗，第 {rounds} 輪修正")
            result, session = claude(
                self.cfg, render("fix.md", OUTPUT=tail(output)), cwd=wt, schema=DEVELOP_SCHEMA,
                tools=tools, deny_extra=self.p.get("deny_tools_extra", []), name=f"{self.slug}-{iid}-fix{rounds}", resume=session)

        sh(["git", "add", "-A"], cwd=wt)
        sh(["git", "commit", "-q", "-m", f"{result['commit_title']}\n\nRefs #{iid}"], cwd=wt)
        output, warning = self.merge_default(wt, output)
        if warning:
            result["summary"] = f"> ⚠️ {warning}\n\n" + result["summary"]
        self.git("push", "-f", "origin", f"HEAD:refs/heads/{branch}", cwd=wt)
        mr = self.upsert_mr(issue, branch, result, output)
        log(f"  ✔ MR !{mr['iid']} {mr['web_url']}" + ("（未達標）" if result.get("unmet") else ""))
        if result.get("unmet"):
            notify(self.cfg, f"📉 {self.path}#{iid} 開了 MR !{mr['iid']} 但沒達到門檻，要你決定：{result['unmet'][:300]}\n{mr['web_url']}")
        self.comment(iid, render("review_comment.md", AUTHOR=issue["author"]["username"],
                                 MR=f"!{mr['iid']}", URL=mr["web_url"], SUMMARY=result["summary"]))
        self.set_labels(issue, status="ai::review")

    def merge_default(self, wt, output):
        """開發期間 main 可能已經合了別的 MR：先合併最新的 main 再驗證一次，合併時的衝突在開 MR 前就看得到。
        有衝突或合併後測試失敗時退回合併前的 commit，照樣開 MR，回傳警告寫進說明（不丟掉這次的成果）。"""
        base = f"origin/{self.default_branch}"
        self.git("fetch", "-q", "origin", self.default_branch, cwd=wt)
        if sh(["git", "merge-base", "--is-ancestor", base, "HEAD"], cwd=wt, check=False).returncode == 0:
            return output, None
        before = sh(["git", "rev-parse", "HEAD"], cwd=wt).stdout.strip()
        m = sh(["git", "merge", "--no-edit", base], cwd=wt, check=False)
        if m.returncode != 0:
            sh(["git", "merge", "--abort"], cwd=wt, check=False)
            log("  ⚠️ 和最新的 main 有衝突")
            return output, f"和最新的 {self.default_branch} 有衝突，合併前需要處理：\n```\n{tail(m.stdout + m.stderr, 15)}\n```"
        ok, merged_output = self.verify(wt)
        if not ok:
            sh(["git", "reset", "-q", "--hard", before], cwd=wt)
            log("  ⚠️ 合併最新的 main 後測試失敗")
            return output, f"合併最新的 {self.default_branch} 後測試失敗（MR 是合併前的版本）：\n```\n{tail(merged_output, 30)}\n```"
        log("  已合併最新的 main 並重新驗證")
        return merged_output, None

    def verify(self, wt):
        cmd = self.p.get("test_cmd")
        if not cmd:
            return True, "（未設定 test_cmd，僅依 AI 自行驗證）"
        try:
            p = sh(cmd, cwd=wt, timeout=self.p.get("test_timeout", 1800), check=False)
        except subprocess.TimeoutExpired:
            return False, f"測試逾時：{cmd}"
        return p.returncode == 0, f"$ {cmd}\n{p.stdout}\n{p.stderr}"

    def upsert_mr(self, issue, branch, result, test_output):
        unmet = result.get("unmet", "").strip()
        body = render("mr_description.md", SUMMARY=result["summary"], TEST_NOTES=result["test_notes"],
                      TEST_OUTPUT=tail(test_output, 30), IID=str(issue["iid"]),
                      UNMET=f"> ⚠️ **未達標**：{unmet}\n\n" if unmet else "", RUN_STATS=stats_md())
        existing = self.gl.get(self.api + "/merge_requests", source_branch=branch, state="opened")
        data = {"title": ("未達標：" if unmet else "") + result["commit_title"], "description": body}
        if self.cfg.get("request_review_from_author", True):
            data["reviewer_ids"] = [issue["author"]["id"]]
        if existing:
            # 已有 MR 就保留原標題，不換成這次修改的 commit 標題
            data.pop("title")
            return self.gl.put(f"{self.api}/merge_requests/{existing[0]['iid']}", data)
        # 剛 push 完 GitLab 可能還沒登記這個分支（回 400 source_branch does not exist）：等一下再試，最多約 1 分鐘
        for attempt in range(12):
            try:
                return self.gl.post(self.api + "/merge_requests", {
                    **data, "source_branch": branch, "target_branch": self.default_branch,
                    "remove_source_branch": True})
            except RuntimeError as e:
                if "does not exist" not in str(e) or attempt == 11:
                    raise
                time.sleep(5)


# ---------- 主程式 ----------

def load(args):
    path = os.path.join(ROOT, args.config)
    if not os.path.isfile(path):
        sys.exit(f"找不到 {args.config}，先執行：cp config.example.json {args.config}")
    with open(path) as f:
        cfg = json.load(f)
    token = [REDACTED]"token_env", "GITLAB_BOT_TOKEN"))
    token_file = os.path.expanduser(cfg.get("token_file", ""))
    if not token and token_file and os.path.isfile(token_file):
        with open(token_file) as f:
            token = [REDACTED]
    if not token:
        [REDACTED]"缺少 bot 的 access token：設定環境變數 {cfg.get('token_env', 'GITLAB_BOT_TOKEN')} 或檔案 {token_file or '（token_file 未設定）'}")
    gl = GitLab(cfg["gitlab_url"], token)
    return cfg, [Project(gl, cfg, p) for p in cfg["projects"]]


def one_round(projects, dry=False, only=None):
    for proj in projects:
        if only and only[0] != proj.path:
            continue
        states = {}
        for issue in proj.issues():
            if only and issue["iid"] != only[1]:
                continue
            # 同一時間只有一個迴圈（.lock），輪到這裡還標著開發中，表示上次處理到一半就中斷了
            # （例如斷網：開發失敗後連留言、標「卡住」都送不出去，就一直停在開發中、再也不會被處理）：清掉狀態重做
            if not dry and proj.state(issue) == "ai::doing":
                log(f"↻ {proj.path}#{issue['iid']} 上次開發中斷（還標著開發中），重新處理")
                proj.set_labels(issue, status=None)
            waits = proj.waiting(issue, states)
            if waits and not only:
                if proj.next_action(issue) is not None:
                    log(f"⏸ {proj.path}#{issue['iid']} 等 {'、'.join(waits)} 關閉")
                continue
            if os.path.exists(os.path.join(ROOT, "STOP")):
                log("🛑 偵測到 STOP")
                return
            action = proj.next_action(issue)
            if only and action is None:
                action = "triage"  # 指定 issue 時強制處理
            if action is None:
                continue
            if dry:
                log(f"會處理 {proj.path}#{issue['iid']}（目前 {proj.state(issue) or '新'}）{issue['title']}")
                continue
            try:
                proj.handle(issue)
            except Exception as e:
                log(f"✖ {proj.path}#{issue['iid']} 失敗：{e}")


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("cmd", choices=["once", "watch", "plan", "issue"])
    ap.add_argument("target", nargs="?", help="issue 指令用：group/proj#12")
    ap.add_argument("--interval", type=int, default=300)
    ap.add_argument("--config", default="config.json")
    args = ap.parse_args()
    cfg, projects = load(args)

    lock = open(os.path.join(ROOT, ".lock"), "w")
    try:
        fcntl.flock(lock, fcntl.LOCK_EX | fcntl.LOCK_NB)
    except BlockingIOError:
        sys.exit("已有另一個迴圈在跑")

    if args.cmd == "plan":
        one_round(projects, dry=True)
    elif args.cmd == "once":
        one_round(projects)
    elif args.cmd == "issue":
        m = re.fullmatch(r"(.+)#(\d+)", args.target or "")
        if not m:
            sys.exit("格式：group/proj#12")
        one_round(projects, only=(m.group(1), int(m.group(2))))
    else:
        log(f"開始輪詢，每 {args.interval} 秒一次；touch STOP 停止")
        while not os.path.exists(os.path.join(ROOT, "STOP")):
            try:
                one_round(projects)
            except Exception as e:
                log(f"✖ 本輪失敗：{e}")
            time.sleep(args.interval)
        os.remove(os.path.join(ROOT, "STOP"))
        log("🛑 停止")


if __name__ == "__main__":
    main()
