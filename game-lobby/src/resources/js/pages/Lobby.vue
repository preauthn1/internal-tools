<script setup lang="ts">
import { Head, Link, router, usePage } from '@inertiajs/vue3';
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { login, play, settle } from '@/routes';
import type { Game } from '@/types';

defineProps<{
    games: Game[];
    gamesUnavailable: boolean;
    csrfToken: [REDACTED];
}>();

const page = usePage();
const user = computed(() => page.props.auth.user);
const playError = computed(() => page.props.errors.play);
const settlementError = computed(() => page.props.errors.settlement);
const settling = ref(false);
const LAUNCHED_KEY = 'game:launched';

function keyOf(game: Game): string {
    return `${game.provider}:${game.gameId}`;
}

function launchHref(game: Game): string {
    return play({ provider: game.provider, gameId: game.gameId }).url;
}

function markLaunched(): void {
    sessionStorage.setItem(LAUNCHED_KEY, '1');
}

function settleIfReturned(): void {
    if (sessionStorage.getItem(LAUNCHED_KEY) == null || settling.value) {
        return;
    }

    settling.value = true;
    router.post(
        settle().url,
        {},
        {
            preserveScroll: true,
            onSuccess: () => sessionStorage.removeItem(LAUNCHED_KEY),
            onFinish: () => {
                settling.value = false;
            },
        },
    );
}

function onVisible(): void {
    if (document.visibilityState === 'visible') {
        settleIfReturned();
    }
}

onMounted(() => {
    settleIfReturned();
    document.addEventListener('visibilitychange', onVisible);
});

onUnmounted(() => document.removeEventListener('visibilitychange', onVisible));
</script>

<template>
    <Head title="遊戲大廳" />
    <h1 class="mb-6 text-2xl font-semibold">遊戲大廳</h1>

    <div
        v-if="playError"
        class="mb-4 rounded-md border border-destructive/40 bg-destructive/10 px-4 py-3 text-sm text-destructive"
    >
        {{ playError }}
    </div>
    <div
        v-if="settlementError || settling"
        class="mb-4 rounded-md border border-amber-400/40 bg-amber-50 px-4 py-3 text-sm text-amber-800"
    >
        {{ settlementError ?? '結算處理中，請稍候。' }}
    </div>

    <div
        v-if="!games.length"
        class="rounded-lg border border-dashed p-10 text-center text-muted-foreground"
    >
        <template v-if="gamesUnavailable"
            >目前無法取得遊戲清單，請稍後再試。</template
        >
        <template v-else>目前沒有可用遊戲。</template>
    </div>

    <div v-else class="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
        <div
            v-for="game in games"
            :key="keyOf(game)"
            class="overflow-hidden rounded-xl border bg-card"
        >
            <div class="flex aspect-video items-center justify-center bg-muted">
                <img
                    v-if="game.image"
                    :src="game.image"
                    :alt="game.name"
                    loading="lazy"
                    decoding="async"
                    class="h-full w-full object-cover"
                />
                <span
                    v-else
                    class="px-2 text-center text-sm text-muted-foreground"
                    >{{ game.name }}</span
                >
            </div>
            <div class="p-3">
                <div class="truncate text-sm font-medium">{{ game.name }}</div>
                <div class="mt-0.5 text-xs text-muted-foreground uppercase">
                    {{ game.provider }}
                </div>

                <Link
                    v-if="!user"
                    :href="login()"
                    class="mt-2 block w-full rounded-md bg-primary px-3 py-1.5 text-center text-xs text-primary-foreground hover:bg-primary/90"
                >
                    登入後遊玩
                </Link>
                <button
                    v-else-if="!game.isLaunchAvailable"
                    type="button"
                    disabled
                    class="mt-2 w-full cursor-not-allowed rounded-md bg-muted px-3 py-1.5 text-xs text-muted-foreground"
                >
                    暫不開放
                </button>
                <form
                    v-else
                    method="post"
                    :action="launchHref(game)"
                    target="_blank"
                    class="mt-2"
                    @submit="markLaunched"
                >
                    <input type="hidden" name="_token" :value="csrfToken" />
                    <button
                        type="submit"
                        class="block w-full rounded-md bg-primary px-3 py-1.5 text-center text-xs text-primary-foreground hover:bg-primary/90"
                    >
                        進入遊戲
                    </button>
                </form>
            </div>
        </div>
    </div>
</template>
