<x-app-layout>
    <x-slot name="header">
        <flux:heading>{{ __('Game Vendors') }}</flux:heading>
        <flux:subheading>{{ __('Available game vendor integrations') }}</flux:subheading>
    </x-slot>

    <div class="py-6">
        <div class="max-w-7xl mx-auto sm:px-6 lg:px-8">
            @if(empty($vendors))
                <flux:callout variant="warning">
                    {{ __('No game vendors configured.') }}
                </flux:callout>
            @else
                <div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
                    @foreach($vendors as $vendorKey => $vendorConfig)
                        <a 
                            href="{{ route('game-vendors.index', $vendorKey) }}" 
                            class="block bg-white overflow-hidden shadow-sm sm:rounded-lg border border-zinc-200 hover:border-zinc-300 hover:shadow-md transition-all duration-200"
                        >
                            <div class="p-5">
                                <div class="flex flex-col items-center justify-center text-center">
                                    @if(!empty($vendorConfig['icon']))
                                        <img 
                                            src="{{ $vendorConfig['icon'] }}" 
                                            alt="{{ strtoupper($vendorKey) }}"
                                            class="h-20 w-20 object-contain rounded-lg mb-3"
                                            onerror="this.src='data:image/svg+xml,%3Csvg xmlns=%22http://example.invalid/2000/svg%22 viewBox=%220 0 24 24%22%3E%3Cpath fill=%22%23999%22 d=%22M12 2L2 7v10c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V7l-10-5z%22/%3E%3C/svg%3E'"
                                        >
                                    @else
                                        <div class="h-20 w-20 bg-zinc-200 rounded-lg flex items-center justify-center mb-3">
                                            <span class="text-2xl font-bold text-zinc-500">{{ strtoupper(substr($vendorKey, 0, 2)) }}</span>
                                        </div>
                                    @endif
                                    <flux:heading size="sm">
                                        {{ strtoupper($vendorKey) }}
                                    </flux:heading>
                                </div>
                            </div>
                        </a>
                    @endforeach
                </div>
            @endif
        </div>
    </div>
</x-app-layout>
