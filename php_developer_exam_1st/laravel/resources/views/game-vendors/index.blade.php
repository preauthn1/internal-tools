<x-app-layout>
    <x-slot name="header">
        <div class="flex items-center gap-4">
            <a href="{{ route('dashboard') }}" class="text-zinc-600 hover:text-zinc-900 flex items-center">
                <svg class="size-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
                </svg>
            </a>
            <div>
                <flux:heading>{{ strtoupper($vendorId) }} {{ __('Vendor') }}</flux:heading>
                <flux:subheading>{{ __('Game vendor operations and management') }}</flux:subheading>
            </div>
        </div>
    </x-slot>

    <div class="py-6">
        <div class="max-w-7xl mx-auto sm:px-6 lg:px-8">
            <div class="bg-white overflow-hidden shadow-sm sm:rounded-lg border border-zinc-200 mb-6">
                <div class="p-6">
                    <div class="flex items-center justify-between gap-4 mb-4">
                        <div class="flex items-center gap-4">
                            <img
                                src="{{ $vendorIcon }}"
                                alt="{{ strtoupper($vendorId) }}"
                                class="h-16 w-16 object-contain rounded-lg"
                                onerror="this.src='data:image/svg+xml,%3Csvg xmlns=%22http://example.invalid/2000/svg%22 viewBox=%220 0 24 24%22%3E%3Cpath fill=%22%23999%22 d=%22M12 2L2 7v10c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V7l-10-5z%22/%3E%3C/svg%3E'"
                            >
                            <div>
                                <flux:heading size="lg">{{ strtoupper($vendorId) }}</flux:heading>
                                <flux:text variant="subtle">
                                    Vendor Name: {{ $vendorName }}
                                </flux:text>
                            </div>
                        </div>
                        <div class="flex items-center gap-3">
                            <a
                                href="{{ route('game-vendors.withdraw', $vendorId) }}"
                                class="inline-flex items-center px-4 py-2 bg-zinc-200 border border-zinc-300 rounded-md font-semibold text-xs text-zinc-800 uppercase tracking-widest hover:bg-zinc-300 focus:bg-zinc-300 active:bg-zinc-400 focus:outline-none focus:ring-2 focus:ring-zinc-500 focus:ring-offset-2 transition ease-in-out duration-150"
                            >
                                Withdraw
                            </a>
                            <a
                                href="{{ route('game-vendors.game-records', $vendorId) }}"
                                class="inline-flex items-center px-4 py-2 bg-zinc-200 border border-zinc-300 rounded-md font-semibold text-xs text-zinc-800 uppercase tracking-widest hover:bg-zinc-300 focus:bg-zinc-300 active:bg-zinc-400 focus:outline-none focus:ring-2 focus:ring-zinc-500 focus:ring-offset-2 transition ease-in-out duration-150"
                            >
                                Game Records
                            </a>
                        </div>
                    </div>
                </div>
            </div>

            <div class="bg-white overflow-hidden shadow-sm sm:rounded-lg border border-zinc-200">
                <div class="p-6">
                    <flux:heading size="md" class="mb-4">{{ __('Available Games') }}</flux:heading>
                    @if(empty($gameList))
                        <flux:callout variant="warning">
                            {{ __('No games available.') }}
                            <div class="mt-2 text-xs text-zinc-500">
                                Debug: gameList is empty. Type: {{ gettype($gameList) }}, Count: {{ is_array($gameList) ? count($gameList) : 'N/A' }}
                            </div>
                        </flux:callout>
                    @else
                        <div class="mb-4 text-xs text-zinc-500">
                            Debug: Found {{ is_array($gameList) ? count($gameList) : 'N/A' }} items. Type: {{ gettype($gameList) }}
                        </div>
                        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                            @include('game-vendors.partials.game-list', ['games' => $gameList, 'vendorId' => $vendorId])
                        </div>
                    @endif
                </div>
            </div>
        </div>
    </div>
</x-app-layout>

