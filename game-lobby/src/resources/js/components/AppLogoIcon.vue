<script setup lang="ts">
import { Gamepad2 } from '@lucide/vue';
import type { HTMLAttributes } from 'vue';

defineOptions({
    inheritAttrs: false,
});

type Props = {
    className?: HTMLAttributes['class'];
};

defineProps<Props>();
</script>

<template>
    <Gamepad2 :class="className" v-bind="$attrs" />
</template>
