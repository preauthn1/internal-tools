<script lang="ts" setup>
import type { DropdownMenuProps } from './interface';

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '../../ui';

interface Props extends DropdownMenuProps {}

defineOptions({ name: 'DropdownRadioMenu' });
withDefaults(defineProps<Props>(), {});

const emit = defineEmits<{ select: [value: string] }>();

const modelValue = defineModel<string>();

/** update:modelValue 只在值變了才發；select 每次點都發（aichatbot：分辨手動選的語言） */
function handleItemClick(value: string) {
  modelValue.value = value;
  emit('select', value);
}
</script>
<template>
  <DropdownMenu>
    <DropdownMenuTrigger as-child class="flex items-center gap-1">
      <slot></slot>
    </DropdownMenuTrigger>
    <DropdownMenuContent align="start">
      <DropdownMenuGroup>
        <template v-for="menu in menus" :key="menu.value">
          <DropdownMenuItem
            :class="
              menu.value === modelValue
                ? 'bg-accent text-accent-foreground'
                : ''
            "
            class="text-foreground/80 data-[state=checked]:bg-accent data-[state=checked]:text-accent-foreground mb-1 cursor-pointer"
            @click="handleItemClick(menu.value)"
          >
            <component :is="menu.icon" v-if="menu.icon" class="mr-2 size-4" />
            <span
              v-if="!menu.icon"
              :class="menu.value === modelValue ? 'bg-foreground' : ''"
              class="mr-2 size-1.5 rounded-full"
            ></span>
            {{ menu.label }}
          </DropdownMenuItem>
        </template>
      </DropdownMenuGroup>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
