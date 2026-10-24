<script setup lang="ts">
import { Link, usePage } from '@inertiajs/vue3';
import { Gamepad2 } from '@lucide/vue';
import { computed } from 'vue';
import { Button } from '@/components/ui/button';
import { login, logout, records, wallet } from '@/routes';

const page = usePage();
const user = computed(() => page.props.auth.user);
const balance = computed(() => page.props.auth.balance);

const formattedBalance = computed(() =>
    balance.value != null
        ? Number(balance.value).toLocaleString('en-US', {
              minimumFractionDigits: 2,
              maximumFractionDigits: 2,
          })
        : null,
);
</script>

<template>
    <div class="min-h-screen bg-background text-foreground">
        <header class="border-b">
            <div
                class="mx-auto flex h-16 max-w-6xl items-center justify-between px-4"
            >
                <Link href="/" class="flex items-center gap-2 font-semibold">
                    <span
                        class="flex size-8 items-center justify-center rounded-md bg-primary text-primary-foreground"
                    >
                        <Gamepad2 class="size-5" />
                    </span>
                    <span>遊戲平台</span>
                </Link>

                <div class="flex items-center gap-3">
                    <template v-if="user">
                        <div class="text-right">
                            <div class="text-sm leading-tight font-medium">
                                {{ user.name }}
                            </div>
                            <div class="text-xs text-muted-foreground">
                                餘額 {{ formattedBalance }}
                            </div>
                        </div>
                        <Button as-child variant="ghost" size="sm">
                            <Link :href="wallet()">錢包紀錄</Link>
                        </Button>
                        <Button as-child variant="ghost" size="sm">
                            <Link :href="records()">遊戲紀錄</Link>
                        </Button>
                        <Button as-child variant="ghost" size="sm">
                            <Link href="/dev/provider-tester">API 測試</Link>
                        </Button>
                        <Button as-child variant="ghost" size="sm">
                            <Link href="/settings/profile">設定</Link>
                        </Button>
                        <Button as-child variant="outline" size="sm">
                            <Link :href="logout()" as="button">登出</Link>
                        </Button>
                    </template>
                    <template v-else>
                        <Button as-child size="sm">
                            <Link :href="login()">登入</Link>
                        </Button>
                    </template>
                </div>
            </div>
        </header>

        <main class="mx-auto max-w-6xl px-4 py-8">
            <slot />
        </main>
    </div>
</template>
