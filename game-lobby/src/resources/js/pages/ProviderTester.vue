<script setup lang="ts">
import { Head } from '@inertiajs/vue3';
import { computed, reactive, ref } from 'vue';
import { Button } from '@/components/ui/button';

interface ProviderOption {
    key: string;
}

const props = defineProps<{
    providers: ProviderOption[];
    actions: string[];
    defaultPlayerName: string | null;
    csrfToken: [REDACTED];
}>();

// 各動作的中文標籤與所需輸入欄位。
const ACTION_META: Record<
    string,
    { label: string; fields: Field[]; danger?: boolean }
> = {
    games: { label: '遊戲列表　games', fields: [] },
    playerExists: {
        label: '玩家是否存在　playerExists',
        fields: ['playerName'],
    },
    createPlayer: { label: '建立玩家　createPlayer', fields: ['playerName'] },
    walletAmount: {
        label: '遊戲錢包餘額　walletAmount',
        fields: ['playerName'],
    },
    gameLink: {
        label: '遊戲啟動連結　gameLink',
        fields: ['playerName', 'gameId'],
    },
    transferIn: {
        label: '轉入遊戲錢包　transferIn',
        fields: ['playerName', 'amount', 'reference'],
        danger: true,
    },
    transferOut: {
        label: '轉出遊戲錢包　transferOut',
        fields: ['playerName', 'amount', 'reference'],
        danger: true,
    },
    transferRecord: {
        label: '轉帳紀錄　transferRecord',
        fields: ['reference'],
    },
    logoutPlayer: { label: '登出玩家　logoutPlayer', fields: ['playerName'] },
    gameRecords: {
        label: '遊戲紀錄　gameRecords',
        fields: ['startTime', 'endTime', 'playerName', 'pageIndex', 'pageSize'],
    },
};

type Field =
    | 'playerName'
    | 'gameId'
    | 'amount'
    | 'reference'
    | 'startTime'
    | 'endTime'
    | 'pageIndex'
    | 'pageSize';

const FIELD_META: Record<Field, { label: string; placeholder: string }> = {
    playerName: { label: 'playerName', placeholder: 'pgs_xxxx_1' },
    gameId: { label: 'gameId（數字，非 code）', placeholder: '290211' },
    amount: { label: 'amount', placeholder: '10.00（最低 0.10）' },
    reference: {
        label: 'reference（留空自動產生 UUID）',
        placeholder: '留空自動產生',
    },
    startTime: { label: 'startTime', placeholder: '2026-07-21T00:00:00Z' },
    endTime: { label: 'endTime', placeholder: '2026-07-21T01:00:00Z' },
    pageIndex: { label: 'pageIndex', placeholder: '1' },
    pageSize: { label: 'pageSize', placeholder: '200' },
};

const form = reactive({
    provider: props.providers[0]?.key ?? 'target',
    action: 'games',
    playerName: props.defaultPlayerName ?? '',
    gameId: '',
    amount: '10.00',
    reference: '',
    startTime: '',
    endTime: '',
    pageIndex: '1',
    pageSize: '200',
});

const activeFields = computed<Field[]>(
    () => (ACTION_META[form.action]?.fields as Field[]) ?? [],
);
const isDanger = computed(() => ACTION_META[form.action]?.danger ?? false);

interface CallResult {
    ok: boolean;
    action: string;
    apiCode: number | null;
    elapsedMs?: number;
    message?: string;
    providerMessage?: string | null;
    data?: unknown;
}

const loading = ref(false);
const result = ref<CallResult | null>(null);
const log = ref<CallResult[]>([]);

const games = computed(() => {
    if (result.value?.ok && result.value.action === 'games') {
        return result.value.data as Array<{
            gameId: string;
            code: string;
            name: string;
        }>;
    }

    return null;
});

const prettyData = computed(() =>
    result.value?.data !== undefined
        ? JSON.stringify(result.value.data, null, 2)
        : '',
);

const gameLink = computed((): string | null => {
    if (!result.value?.ok || result.value.action !== 'gameLink') {
        return null;
    }

    const data = result.value.data;

    if (typeof data !== 'object' || data === null || !('link' in data)) {
        return null;
    }

    return typeof data.link === 'string' ? data.link : null;
});

const submit = async () => {
    loading.value = true;
    const payload: Record<string, unknown> = {
        provider: form.provider,
        action: form.action,
    };

    for (const field of activeFields.value) {
        const value = form[field];

        if (value === '') {
            continue;
        }

        if (field === 'pageIndex' || field === 'pageSize') {
            payload[field] = Number(value);
        } else {
            payload[field] = value;
        }
    }

    try {
        const response = await fetch('/dev/provider-tester/call', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                Accept: 'application/json',
                'X-Requested-With': 'XMLHttpRequest',
                'X-CSRF-TOKEN': props.csrfToken,
            },
            body: JSON.stringify(payload),
        });
        const json = (await response.json()) as CallResult;
        result.value = json;
        log.value.unshift(json);

        if (log.value.length > 20) {
            log.value.pop();
        }
    } catch (error) {
        result.value = {
            ok: false,
            action: form.action,
            apiCode: null,
            message: '前端請求失敗：' + String(error),
        };
    } finally {
        loading.value = false;
    }
};

const useGameId = (gameId: string) => {
    form.gameId = gameId;
    form.action = 'gameLink';
};

const fillRandomReference = () => {
    form.reference = 'test-' + Math.random().toString(36).slice(2, 12);
};

const openLink = (link: string) => window.open(link, '_blank');
</script>

<template>
    <Head title="遊戲商 API 測試" />

    <div class="mb-6">
        <h1 class="text-2xl font-semibold">遊戲商 API 測試</h1>
        <p class="mt-1 text-sm text-muted-foreground">
            直接呼叫遊戲商 driver
            的各端點並顯示原始結果／錯誤碼。此頁不經過平台錢包帳本，
            <span class="font-medium text-amber-600 dark:text-amber-500"
                >transferIn／transferOut
                會實際變動遊戲商錢包餘額，請謹慎使用。</span
            >
        </p>
    </div>

    <div class="grid gap-6 lg:grid-cols-[360px_1fr]">
        <!-- 表單 -->
        <div class="space-y-4 rounded-xl border bg-card p-4">
            <div class="grid gap-1.5">
                <label class="text-sm font-medium">遊戲商 provider</label>
                <select
                    v-model="form.provider"
                    class="rounded-md border bg-background px-3 py-2 text-sm"
                >
                    <option v-for="p in providers" :key="p.key" :value="p.key">
                        {{ p.key }}
                    </option>
                </select>
            </div>

            <div class="grid gap-1.5">
                <label class="text-sm font-medium">動作 action</label>
                <select
                    v-model="form.action"
                    class="rounded-md border bg-background px-3 py-2 text-sm"
                >
                    <option v-for="a in actions" :key="a" :value="a">
                        {{ ACTION_META[a]?.label ?? a }}
                    </option>
                </select>
            </div>

            <div
                v-for="field in activeFields"
                :key="field"
                class="grid gap-1.5"
            >
                <label class="text-sm font-medium">{{
                    FIELD_META[field].label
                }}</label>
                <div v-if="field === 'reference'" class="flex gap-2">
                    <input
                        v-model="form[field]"
                        :placeholder="FIELD_META[field].placeholder"
                        class="min-w-0 flex-1 rounded-md border bg-background px-3 py-2 text-sm"
                    />
                    <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        @click="fillRandomReference"
                        >亂數</Button
                    >
                </div>
                <input
                    v-else
                    v-model="form[field]"
                    :placeholder="FIELD_META[field].placeholder"
                    class="rounded-md border bg-background px-3 py-2 text-sm"
                />
            </div>

            <Button
                class="w-full"
                :disabled="loading"
                :variant="isDanger ? 'destructive' : 'default'"
                @click="submit"
            >
                {{ loading ? '呼叫中…' : '送出呼叫' }}
            </Button>
        </div>

        <!-- 結果 -->
        <div class="space-y-4">
            <div
                v-if="result"
                class="rounded-xl border bg-card p-4"
                :class="
                    result.ok
                        ? 'border-emerald-500/40'
                        : 'border-destructive/50'
                "
            >
                <div class="flex flex-wrap items-center gap-3 text-sm">
                    <span
                        class="rounded-md px-2 py-0.5 font-medium"
                        :class="
                            result.ok
                                ? 'bg-emerald-500/15 text-emerald-600 dark:text-emerald-400'
                                : 'bg-destructive/15 text-destructive'
                        "
                    >
                        {{ result.ok ? 'OK' : '失敗' }}
                    </span>
                    <span class="text-muted-foreground"
                        >action：<code>{{ result.action }}</code></span
                    >
                    <span class="text-muted-foreground"
                        >apiCode：<code>{{ result.apiCode ?? '—' }}</code></span
                    >
                    <span
                        v-if="result.elapsedMs !== undefined"
                        class="text-muted-foreground"
                        >耗時：{{ result.elapsedMs }} ms</span
                    >
                </div>

                <p v-if="result.message" class="mt-3 text-sm text-destructive">
                    {{ result.message }}
                </p>

                <div v-if="result.providerMessage" class="mt-3">
                    <p class="mb-1 text-xs text-muted-foreground">
                        遊戲商原始回應訊息
                    </p>
                    <pre
                        class="overflow-x-auto rounded-md bg-muted p-3 text-xs whitespace-pre-wrap"
                        >{{ result.providerMessage }}</pre>
                </div>

                <!-- 遊戲列表以表格呈現，方便直接帶入 gameId -->
                <div v-if="games" class="mt-3 overflow-x-auto">
                    <p class="mb-2 text-xs text-muted-foreground">
                        共 {{ games.length }} 款；點「用此 gameId」可帶入並切到
                        gameLink。
                    </p>
                    <table class="w-full text-sm">
                        <thead class="border-b text-left text-muted-foreground">
                            <tr>
                                <th class="py-2 pr-3 font-medium">gameId</th>
                                <th class="py-2 pr-3 font-medium">code</th>
                                <th class="py-2 pr-3 font-medium">name</th>
                                <th class="py-2"></th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr
                                v-for="g in games"
                                :key="g.gameId + g.code"
                                class="border-b last:border-0"
                            >
                                <td class="py-1.5 pr-3 tabular-nums">
                                    {{ g.gameId }}
                                </td>
                                <td class="py-1.5 pr-3 tabular-nums">
                                    {{ g.code }}
                                </td>
                                <td class="py-1.5 pr-3">{{ g.name }}</td>
                                <td class="py-1.5 text-right">
                                    <Button
                                        type="button"
                                        variant="outline"
                                        size="sm"
                                        @click="useGameId(g.gameId)"
                                        >用此 gameId</Button
                                    >
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <!-- gameLink 結果附開啟按鈕 -->
                <div v-else-if="gameLink" class="mt-3 space-y-2">
                    <Button type="button" size="sm" @click="openLink(gameLink)"
                        >於新分頁開啟遊戲</Button
                    >
                    <pre
                        class="overflow-x-auto rounded-md bg-muted p-3 text-xs"
                        >{{ prettyData }}</pre>
                </div>

                <pre
                    v-else-if="result.ok"
                    class="mt-3 overflow-x-auto rounded-md bg-muted p-3 text-xs"
                    >{{ prettyData }}</pre>
            </div>

            <div
                v-else
                class="rounded-xl border border-dashed p-10 text-center text-muted-foreground"
            >
                選擇動作並送出後，結果會顯示在這裡。
            </div>

            <!-- 呼叫紀錄 -->
            <div v-if="log.length" class="rounded-xl border bg-card p-4">
                <h2 class="mb-2 text-sm font-medium">最近呼叫</h2>
                <ul class="space-y-1 text-xs">
                    <li
                        v-for="(item, index) in log"
                        :key="index"
                        class="flex items-center gap-2"
                    >
                        <span
                            class="inline-block h-2 w-2 rounded-full"
                            :class="
                                item.ok ? 'bg-emerald-500' : 'bg-destructive'
                            "
                        />
                        <code>{{ item.action }}</code>
                        <span class="text-muted-foreground"
                            >apiCode={{ item.apiCode ?? '—' }}</span
                        >
                        <span
                            v-if="item.elapsedMs !== undefined"
                            class="text-muted-foreground"
                            >· {{ item.elapsedMs }}ms</span
                        >
                        <span
                            v-if="item.message"
                            class="truncate text-muted-foreground"
                            >· {{ item.message }}</span
                        >
                    </li>
                </ul>
            </div>
        </div>
    </div>
</template>
