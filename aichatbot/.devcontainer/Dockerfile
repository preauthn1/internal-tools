# Go 版的開發環境（ADR 0001）：Go、node＋pnpm（Vben 後台）、python3（eval/、loop/ 的腳本）、sqlite3
# Apple 本機模型只能在 macOS 上跑，Swift（模型服務、swift-final 的程式）不放進來，在 Mac 上編譯執行
FROM example.invalid/devcontainers/go:1.25-bookworm

RUN apt-get update \
    && apt-get install -y --no-install-recommends python3 sqlite3 \
    && rm -rf /var/lib/apt/lists/*
