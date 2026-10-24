<x-app-layout>
    <x-slot name="header">
        <div class="flex items-center gap-4">
            <a href="{{ route('game-vendors.index', $vendorId) }}" class="text-zinc-600 hover:text-zinc-900 flex items-center">
                <svg class="size-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
                </svg>
            </a>
            <div>
                <flux:heading>{{ strtoupper($vendorId) }} {{ __('Game Records') }}</flux:heading>
                <flux:subheading>{{ __('Your game records history') }}</flux:subheading>
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

            <div class="bg-white overflow-hidden shadow-sm sm:rounded-lg border border-zinc-200">
                <div class="p-6">
                    <flux:heading size="md" class="mb-4">{{ __('Game Records') }}</flux:heading>
                    
                    @if($gameRecords->isEmpty())
                        <flux:callout variant="info">
                            {{ __('No game records found.') }}
                        </flux:callout>
                    @else
                        <div class="overflow-x-auto">
                            <table class="min-w-full divide-y divide-zinc-200">
                                <thead class="bg-zinc-50">
                                    <tr>
                                        <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase tracking-wider">
                                            {{ __('Game Name') }}
                                        </th>
                                        <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase tracking-wider">
                                            {{ __('Bet ID') }}
                                        </th>
                                        <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase tracking-wider">
                                            {{ __('Bet Time') }}
                                        </th>
                                        <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase tracking-wider">
                                            {{ __('Bet Amount') }}
                                        </th>
                                        <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase tracking-wider">
                                            {{ __('Payout Amount') }}
                                        </th>
                                        <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase tracking-wider">
                                            {{ __('Result Amount') }}
                                        </th>
                                        <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase tracking-wider">
                                            {{ __('Valid Amount') }}
                                        </th>
                                        <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase tracking-wider">
                                            {{ __('Tax Amount') }}
                                        </th>
                                    </tr>
                                </thead>
                                <tbody class="bg-white divide-y divide-zinc-200">
                                    @foreach($gameRecords as $record)
                                        @php
                                            $gameName = $gameNameMap[$record->game_id] ?? $record->game_id ?? __('Unknown Game');
                                        @endphp
                                        <tr class="hover:bg-zinc-50">
                                            <td class="px-6 py-4 whitespace-nowrap text-sm text-zinc-900">
                                                {{ $gameName }}
                                            </td>
                                            <td class="px-6 py-4 whitespace-nowrap text-sm text-zinc-900">
                                                {{ $record->bet_id }}
                                            </td>
                                            <td class="px-6 py-4 whitespace-nowrap text-sm text-zinc-500">
                                                {{ $record->bet_time }}
                                            </td>
                                            <td class="px-6 py-4 whitespace-nowrap text-sm text-zinc-900">
                                                {{ number_format($record->bet_amount, 2) }}
                                            </td>
                                            <td class="px-6 py-4 whitespace-nowrap text-sm text-zinc-900">
                                                {{ number_format($record->payout_amount, 2) }}
                                            </td>
                                            <td class="px-6 py-4 whitespace-nowrap text-sm {{ $record->bet_result_amount >= 0 ? 'text-green-600' : 'text-red-600' }}">
                                                {{ number_format($record->bet_result_amount, 2) }}
                                            </td>
                                            <td class="px-6 py-4 whitespace-nowrap text-sm text-zinc-900">
                                                {{ number_format($record->valid_amount, 2) }}
                                            </td>
                                            <td class="px-6 py-4 whitespace-nowrap text-sm text-zinc-500">
                                                {{ number_format($record->tax_amount, 2) }}
                                            </td>
                                        </tr>
                                    @endforeach
                                </tbody>
                            </table>
                        </div>

                        <!-- Laravel's built-in pagination -->
                        <div class="mt-4">
                            {{ $gameRecords->links() }}
                        </div>
                    @endif
                </div>
            </div>
        </div>
    </div>
</x-app-layout>

