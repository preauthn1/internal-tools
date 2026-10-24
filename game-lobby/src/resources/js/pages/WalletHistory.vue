<script setup lang="ts">
import { Head, Link } from '@inertiajs/vue3';
import type { Paginated, WalletTransaction } from '@/types';

defineProps<{
    transactions: Paginated<WalletTransaction>;
    balance: string;
}>();

const formatAmount = (value: string): string => {
    const n = Number(value);
    const formatted = Math.abs(n).toLocaleString('en-US', {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
    });

    return `${n < 0 ? '-' : '+'}${formatted}`;
};

const formatBalance = (value: string): string =>
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
    <Head title="錢包紀錄" />

    <div class="mb-6 flex items-baseline justify-between">
        <h1 class="text-2xl font-semibold">錢包紀錄</h1>
        <div class="text-sm text-muted-foreground">
            目前餘額
            <span class="font-medium text-foreground tabular-nums">
                {{ formatBalance(balance) }}
            </span>
        </div>
    </div>

    <div
        v-if="!transactions.data.length"
        class="rounded-lg border border-dashed p-10 text-center text-muted-foreground"
    >
        尚無錢包紀錄。
    </div>

    <div v-else class="overflow-x-auto rounded-xl border bg-card">
        <table class="w-full text-sm">
            <thead class="border-b bg-muted/50 text-left text-muted-foreground">
                <tr>
                    <th class="px-4 py-3 font-medium">時間</th>
                    <th class="px-4 py-3 font-medium">類型</th>
                    <th class="px-4 py-3 font-medium">遊戲商</th>
                    <th class="px-4 py-3 text-right font-medium">金額</th>
                    <th class="px-4 py-3 text-right font-medium">結餘</th>
                </tr>
            </thead>
            <tbody>
                <tr
                    v-for="tx in transactions.data"
                    :key="tx.id"
                    class="border-b last:border-0"
                >
                    <td class="px-4 py-3 text-muted-foreground">
                        {{ formatTime(tx.created_at) }}
                    </td>
                    <td class="px-4 py-3">{{ tx.type_label }}</td>
                    <td class="px-4 py-3 text-muted-foreground">
                        {{ tx.provider_name ?? '—' }}
                    </td>
                    <td
                        class="px-4 py-3 text-right tabular-nums"
                        :class="
                            Number(tx.amount) < 0
                                ? 'text-red-600 dark:text-red-400'
                                : 'text-green-600 dark:text-green-400'
                        "
                    >
                        {{ formatAmount(tx.amount) }}
                    </td>
                    <td class="px-4 py-3 text-right tabular-nums">
                        {{ formatBalance(tx.balance_after) }}
                    </td>
                </tr>
            </tbody>
        </table>
    </div>

    <div
        v-if="transactions.last_page > 1"
        class="mt-4 flex flex-wrap items-center gap-1"
    >
        <template v-for="(link, index) in transactions.links" :key="index">
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
