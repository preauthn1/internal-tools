<x-app-layout>
    <x-slot name="header">
        <div class="flex items-center gap-4">
            <a href="{{ route('game-vendors.index', $vendorId) }}" class="text-zinc-600 hover:text-zinc-900 flex items-center">
                <svg class="size-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
                </svg>
            </a>
            <div>
                <flux:heading>{{ strtoupper($vendorId) }} {{ __('Withdraw') }}</flux:heading>
                <flux:subheading>{{ __('Withdraw funds from vendor account') }}</flux:subheading>
            </div>
        </div>
    </x-slot>

    <div class="py-6">
        <div class="max-w-7xl mx-auto sm:px-6 lg:px-8">
            <div class="bg-white overflow-hidden shadow-sm sm:rounded-lg border border-zinc-200 mb-6">
                <div class="p-6">
                    <div class="flex items-center justify-between gap-4 mb-4">
                        <div class="flex items-center gap-4">
                            @if($vendorIcon)
                                <img
                                    src="{{ $vendorIcon }}"
                                    alt="{{ strtoupper($vendorId) }}"
                                    class="h-16 w-16 object-contain rounded-lg"
                                    onerror="this.src='data:image/svg+xml,%3Csvg xmlns=%22http://example.invalid/2000/svg%22 viewBox=%220 0 24 24%22%3E%3Cpath fill=%22%23999%22 d=%22M12 2L2 7v10c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V7l-10-5z%22/%3E%3C/svg%3E'"
                                >
                            @endif
                            <div>
                                <flux:heading size="lg">{{ strtoupper($vendorId) }}</flux:heading>
                                <flux:text variant="subtle">
                                    Vendor Name: {{ $vendorName }}
                                </flux:text>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div class="bg-white overflow-hidden shadow-sm sm:rounded-lg border border-zinc-200 mb-6">
                <div class="p-6">
                    <flux:heading size="md" class="mb-4">{{ __('Account Balance') }}</flux:heading>
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div class="p-4 bg-zinc-50 rounded-lg">
                            <flux:text variant="subtle" class="text-sm">{{ __('Local Balance') }}</flux:text>
                            <flux:heading size="lg" class="mt-1">{{ number_format($user->balance, 2) }}</flux:heading>
                        </div>
                        <div class="p-4 bg-zinc-50 rounded-lg">
                            <flux:text variant="subtle" class="text-sm">{{ __('Vendor Balance') }}</flux:text>
                            <flux:heading size="lg" class="mt-1">{{ number_format($vendorBalance ?? 0, 2) }}</flux:heading>
                        </div>
                    </div>
                </div>
            </div>

            <div class="bg-white overflow-hidden shadow-sm sm:rounded-lg border border-zinc-200">
                <div class="p-6">
                    <flux:heading size="md" class="mb-4">{{ __('Withdraw Funds') }}</flux:heading>

                    @if(session('success'))
                        <flux:callout variant="success" class="mb-4">
                            {{ session('success') }}
                        </flux:callout>
                    @endif

                    @if($errors->any())
                        <flux:callout variant="danger" class="mb-4">
                            <ul class="list-disc list-inside">
                                @foreach($errors->all() as $error)
                                    <li>{{ $error }}</li>
                                @endforeach
                            </ul>
                        </flux:callout>
                    @endif

                    @if(($vendorBalance ?? 0) > 0)
                        <div class="mb-4">
                            <flux:callout variant="info">
                                {{ __('You are about to withdraw all your vendor balance: ') }}
                                <strong>{{ number_format($vendorBalance ?? 0, 2) }}</strong>
                            </flux:callout>
                        </div>

                        <form method="POST" action="{{ route('game-vendors.withdraw.process', $vendorId) }}" class="space-y-4">
                            @csrf

                            <div class="flex items-center gap-3">
                                <flux:button type="submit" variant="primary" class="text-center">
                                    {{ __('Withdraw') }}
                                </flux:button>
                                <a
                                    href="{{ route('game-vendors.index', $vendorId) }}"
                                    class="inline-flex items-center px-4 py-2 bg-zinc-200 border border-zinc-300 rounded-md font-semibold text-xs text-zinc-800 uppercase tracking-widest hover:bg-zinc-300 focus:bg-zinc-300 active:bg-zinc-400 focus:outline-none focus:ring-2 focus:ring-zinc-500 focus:ring-offset-2 transition ease-in-out duration-150"
                                >
                                    {{ __('Cancel') }}
                                </a>
                            </div>
                        </form>
                    @else
                        <flux:callout variant="warning">
                            {{ __('No balance available to withdraw.') }}
                        </flux:callout>
                        <div class="mt-4">
                            <a
                                href="{{ route('game-vendors.index', $vendorId) }}"
                                class="inline-flex items-center px-4 py-2 bg-zinc-200 border border-zinc-300 rounded-md font-semibold text-xs text-zinc-800 uppercase tracking-widest hover:bg-zinc-300 focus:bg-zinc-300 active:bg-zinc-400 focus:outline-none focus:ring-2 focus:ring-zinc-500 focus:ring-offset-2 transition ease-in-out duration-150"
                            >
                                {{ __('Back to Games') }}
                            </a>
                        </div>
                    @endif
                </div>
            </div>
        </div>
    </div>
</x-app-layout>

