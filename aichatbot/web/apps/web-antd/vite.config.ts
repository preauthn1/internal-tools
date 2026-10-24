import { defineConfig } from '@vben/vite-config';

export default defineConfig(async () => {
  return {
    application: {},
    vite: {
      server: {
        proxy: {
          // 開發時轉到 Go 的後台伺服器（bot serve，BOT_ADMIN_PORT 預設 8081）
          '/api': {
            changeOrigin: true,
            target: 'http://0.0.0.0:8081',
          },
        },
      },
    },
  };
});
