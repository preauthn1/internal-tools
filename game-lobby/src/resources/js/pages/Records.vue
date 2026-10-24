<script setup lang="ts">
import { Head, Link } from '@inertiajs/vue3';
import type { GameRecord, Paginated } from '@/types';

const props = defineProps<{
    records: Paginated<GameRecord>;
    gameNames?: Record<string, Record<string, string>>;
}>();

// 遊戲商紀錄常不帶名稱，優先用紀錄本身的 game_name，其次以大廳 games 對照表把編號轉成名稱；
// 一律不顯示遊戲編號，對照不到時以「—」呈現。
const resolveGameName = (record: GameRecord): string =>
    record.game_name ||
    props.gameNames?.[record.provider]?.[record.game_id ?? ''] ||
    props.gameNames?.[record.provider]?.[record.game_code ?? ''] ||
    '—';

const formatAmount = (value: string): string =>
    Number(value).toLocaleString('en-US', {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
    });

const formatTime = (value: string | null): string =>
    value ? new Date(value).toLocaleString() : '-';

const paginationLabel = (value: string): string =>
    value.replace('&laquo;', '«').replace('&raquo;', '»');
</script>

<template>
    <Head title="遊戲紀錄" />

    <h1 class="mb-6 text-2xl font-semibold">遊戲紀錄</h1>

    <div
        v-if="!records.data.length"
        class="rounded-lg border border-dashed p-10 text-center text-muted-foreground"
    >
        尚無紀錄。進入遊戲下注後，系統會定時同步你的遊戲紀錄。
    </div>

    <div v-else class="overflow-x-auto rounded-xl border bg-card">
        <table class="w-full text-sm">
            <thead class="border-b bg-muted/50 text-left text-muted-foreground">
                <tr>
                    <th class="px-4 py-3 font-medium">遊戲</th>
                    <th class="px-4 py-3 text-right font-medium">下注</th>
                    <th class="px-4 py-3 text-right font-medium">派彩</th>
                    <th class="px-4 py-3 text-right font-medium">有效投注</th>
                    <th class="px-4 py-3 font-medium">時間</th>
                </tr>
            </thead>
            <tbody>
                <tr
                    v-for="record in records.data"
                    :key="record.id"
                    class="border-b last:border-0"
                >
                    <td class="px-4 py-3">
                        <div class="font-medium">
                            {{ resolveGameName(record) }}
                        </div>
                        <div class="text-xs text-muted-foreground uppercase">
                            {{ record.provider }}
                        </div>
                    </td>
                    <td class="px-4 py-3 text-right tabular-nums">
                        {{ formatAmount(record.bet_amount) }}
                    </td>
                    <td class="px-4 py-3 text-right tabular-nums">
                        {{ formatAmount(record.payout_amount) }}
                    </td>
                    <td class="px-4 py-3 text-right tabular-nums">
                        {{ formatAmount(record.valid_amount) }}
                    </td>
                    <td class="px-4 py-3 text-muted-foreground">
                        {{ formatTime(record.bet_at) }}
                    </td>
                </tr>
            </tbody>
        </table>
    </div>

    <div
        v-if="records.last_page > 1"
        class="mt-4 flex flex-wrap items-center gap-1"
    >
        <template v-for="(link, index) in records.links" :key="index">
            <Link
                v-if="link.url"
                :href="link.url"
                class="rounded-md border px-3 py-1.5 text-sm"
                :class="
                    link.active
                        ? 'bg-primary text-primary-foreground'
                        : 'hover:bg-muted'
                "
            >
                {{ paginationLabel(link.label) }}
            </Link>
            <span
                v-else
                class="rounded-md border px-3 py-1.5 text-sm text-muted-foreground"
            >
                {{ paginationLabel(link.label) }}
            </span>
        </template>
    </div>
</template>
