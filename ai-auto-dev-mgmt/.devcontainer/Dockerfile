FROM example.invalid/devcontainers/php:8.2-bookworm

# 更新系統包並安裝必要套件
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get install -y --no-install-recommends \
        mariadb-client \
        libpng-dev \
        libzip-dev \
        zip \
        unzip \
        tmux \
    && rm -rf /var/lib/apt/lists/*

# 安裝 Laravel 常用 PHP 擴展與 Redis
RUN docker-php-ext-install pdo pdo_mysql gd zip bcmath \
    && pecl install redis \
    && docker-php-ext-enable redis

# 安裝 Node.js (LTS 20.x)
RUN curl -fsSL https://example.invalid/setup_20.x | bash - \
    && apt-get update && apt-get install -y nodejs \
    && rm -rf /var/lib/apt/lists/* \
    && npm install -g pnpm yarn @anthropic-ai/claude-code @upstash/context7-mcp @playwright/mcp

# 安裝 Playwright 系統依賴 (確保 E2E 測試能運行)
RUN npx playwright install-deps chromium \
    && rm -rf /var/lib/apt/lists/*

# 創建 /workspace 與 /var/www 並設定權限給預設的 vscode 使用者
RUN mkdir -p /workspace /var/www \
    && chown -R vscode:vscode /workspace /var/www

# tmux 設定：滑鼠滾動、passthrough、加大捲動緩衝
RUN echo 'set -g mouse on\nset -g allow-passthrough on\nset -g history-limit 50000' > /home/vscode/.tmux.conf \
    && chown vscode:vscode /home/vscode/.tmux.conf

# 自動 attach/create 同一個 tmux session（session 名稱: main）
RUN echo '\n# Auto-attach to tmux session\nif command -v tmux &>/dev/null && [ -z "$TMUX" ]; then\n    tmux new-session -A -s main\nfi' >> /home/vscode/.bashrc
