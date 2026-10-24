<script setup lang="ts">
import { Form, Head } from '@inertiajs/vue3';
import InputError from '@/components/InputError.vue';
import PasswordInput from '@/components/PasswordInput.vue';
import TextLink from '@/components/TextLink.vue';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Spinner } from '@/components/ui/spinner';
import { login } from '@/routes';
import { store } from '@/routes/register';

defineProps<{
    passwordRules: string;
}>();

defineOptions({
    layout: {
        title: '建立帳號',
        description: '請在下方輸入資料以建立你的帳號',
    },
});
</script>

<template>
    <Head title="註冊" />

    <Form
        v-bind="store.form()"
        :reset-on-success="['password', 'password_confirmation']"
        v-slot="{ errors, processing }"
        class="flex flex-col gap-6"
    >
        <div class="grid gap-6">
            <div class="grid gap-2">
                <Label for="email">電子郵件</Label>
                <Input
                    id="email"
                    type="email"
                    required
                    autofocus
                    :tabindex="1"
                    autocomplete="email"
                    name="email"
                    placeholder="user@example.invalid"
                />
                <InputError :message="errors.email" />
            </div>

            <div class="grid gap-2">
                <Label for="account">帳號</Label>
                <Input
                    id="account"
                    type="text"
                    required
                    :tabindex="2"
                    autocomplete="username"
                    name="account"
                    placeholder="你的帳號"
                />
                <InputError :message="errors.account" />
            </div>

            <div class="grid gap-2">
                <Label for="password">密碼</Label>
                <PasswordInput
                    id="password"
                    required
                    :tabindex="3"
                    autocomplete="new-password"
                    name="password"
                    placeholder="密碼"
                    :passwordrules="passwordRules"
                />
                <InputError :message="errors.password" />
            </div>

            <div class="grid gap-2">
                <Label for="password_confirmation">確認密碼</Label>
                <PasswordInput
                    id="password_confirmation"
                    required
                    :tabindex="4"
                    autocomplete="new-password"
                    name="password_confirmation"
                    placeholder="確認密碼"
                    :passwordrules="passwordRules"
                />
                <InputError :message="errors.password_confirmation" />
            </div>

            <Button
                type="submit"
                class="mt-2 w-full"
                tabindex="5"
                :disabled="processing"
                data-test="register-user-button"
            >
                <Spinner v-if="processing" />
                建立帳號
            </Button>
        </div>

        <div class="text-center text-sm text-muted-foreground">
            已經有帳號了？
            <TextLink
                :href="login()"
                class="underline underline-offset-4"
                :tabindex="6"
                >登入</TextLink
            >
        </div>
    </Form>
</template>
