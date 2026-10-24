<x-app-layout>
    <div class="py-6">
        <div class="max-w-7xl mx-auto sm:px-6 lg:px-8">
            <div class="bg-white overflow-hidden shadow-sm sm:rounded-lg border border-zinc-200">
                <div class="p-6">
                    <div class="flex flex-col items-center text-center">
                        <flux:callout variant="danger" class="mb-6 max-w-2xl">
                            <flux:text class="text-center">
                                {{ $errorMessage }}
                            </flux:text>
                        </flux:callout>

                        <a
                            href="{{ route('game-vendors.index', $vendorId) }}"
                            class="inline-flex items-center px-4 py-2 bg-zinc-200 border border-transparent rounded-md font-semibold text-xs text-zinc-600 uppercase tracking-widest hover:bg-zinc-300 focus:bg-zinc-300 active:bg-zinc-400 focus:outline-none focus:ring-2 focus:ring-zinc-500 focus:ring-offset-2 transition ease-in-out duration-150"
                        >
                            Back to Game Vendors
                        </a>
                    </div>
                </div>
            </div>
        </div>
    </div>
</x-app-layout>

