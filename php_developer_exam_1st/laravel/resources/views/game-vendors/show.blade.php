<x-app-layout>
    <x-slot name="header">
        <div class="flex items-center gap-4">
            <a href="{{ route('game-vendors.index', $vendorId) }}" class="text-zinc-600 hover:text-zinc-900 flex items-center">
                <svg class="size-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
                </svg>
            </a>
            <div>
                <flux:heading>{{ strtoupper($vendorId) }} {{ __('Game') }}</flux:heading>
                <flux:subheading>{{ __('Loading game...') }}</flux:subheading>
            </div>
        </div>
    </x-slot>

    <div class="py-6">
        <div class="max-w-7xl mx-auto sm:px-6 lg:px-8">
            <div class="bg-white overflow-hidden shadow-sm sm:rounded-lg border border-zinc-200">
                <div class="p-6">
                    <div class="flex flex-col items-center text-center">
                        <div class="mb-4">
                            <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-zinc-900"></div>
                        </div>
                        <flux:text class="text-center">
                            {{ __('Redirecting to game...') }}
                        </flux:text>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <script>
        @if(isset($gameUrl) && !empty($gameUrl))
            // Redirect to game URL when page loads
            window.location.href = {!! json_encode($gameUrl) !!};
        @else
            // If no game URL, show error or redirect back
            console.error('Game URL is not available');
            window.location.href = {!! json_encode(route('game-vendors.index', $vendorId)) !!};
        @endif
    </script>
</x-app-layout>

