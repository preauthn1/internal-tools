@php
    // Ensure games is an array
    if (!is_array($games)) {
        $games = [];
    }

    // Recursively collect all game objects
    function collectGames($data, &$result = []) {
        if (!is_array($data)) {
            return;
        }

        // Check if this array itself looks like a game (associative array with game properties)
        $hasName = isset($data['name']) || isset($data['game_name']) || isset($data['title']) || isset($data['gameName']) || isset($data['game_title']) || isset($data['gameNameMappings']);
        $hasId = isset($data['id']) || isset($data['game_id']) || isset($data['gameId']) || isset($data['gameID']) || isset($data['game_code']) || isset($data['gameCode']) || isset($data['code']);
        $isAssociative = !empty($data) && array_keys($data) !== range(0, count($data) - 1);

        // If this is an associative array with game properties, add it
        if ($isAssociative && ($hasName || $hasId)) {
            $result[] = $data;
            return;
        }

        // Otherwise, iterate through items
        foreach ($data as $item) {
            if (is_array($item)) {
                // Check if this item looks like a game object
                $itemHasName = isset($item['name']) || isset($item['game_name']) || isset($item['title']) || isset($item['gameName']) || isset($item['game_title']) || isset($item['gameNameMappings']);
                $itemHasId = isset($item['id']) || isset($item['game_id']) || isset($item['gameId']) || isset($item['gameID']) || isset($item['game_code']) || isset($item['gameCode']) || isset($item['code']);
                $itemIsAssociative = array_keys($item) !== range(0, count($item) - 1);

                // If it's an associative array with game properties, add it
                if ($itemIsAssociative && ($itemHasName || $itemHasId)) {
                    $result[] = $item;
                } else {
                    // Otherwise, recurse into it
                    collectGames($item, $result);
                }
            }
        }
    }

    $allGames = [];
    collectGames($games, $allGames);
@endphp

@if(empty($allGames))
    <div class="col-span-full text-center text-zinc-500 py-8">
        No games found in the response.
    </div>
@else
    @foreach($allGames as $game)
    @php
        // Extract game data
        $gameId = $game['gameID'] ?? $game['gameCode'] ?? $game['id'] ?? $game['game_id'] ?? $game['gameId'] ?? $game['game_code'] ?? $game['code'] ?? null;
        
        // Extract game name - prioritize zh-cn from gameNameMappings
        $gameName = 'Unknown Game';
        if (isset($game['gameNameMappings']) && is_array($game['gameNameMappings'])) {
            foreach ($game['gameNameMappings'] as $mapping) {
                if (isset($mapping['languageCode']) && $mapping['languageCode'] === 'zh-cn' && isset($mapping['text'])) {
                    $gameName = $mapping['text'];
                    break;
                }
            }
        }
        // Fallback to other name fields if zh-cn not found
        if ($gameName === 'Unknown Game') {
            $gameName = $game['name'] ?? $game['game_name'] ?? $game['title'] ?? $game['gameName'] ?? $game['game_title'] ?? 'Unknown Game';
        }
        
        // Extract game icon - prioritize imageUrl
        $gameIcon = $game['imageUrl'] ?? $game['image'] ?? $game['thumbnail'] ?? $game['icon_url'] ?? $game['image_url'] ?? $game['logo'] ?? $game['logo_url'] ?? null;

        // Skip if we don't have at least a name
        if (empty($gameName) || $gameName === 'Unknown Game') {
            continue;
        }

        // Use name as ID if no ID found
        if (empty($gameId)) {
            $gameId = $gameName;
        }
        
        // Ensure gameId is a string for the URL
        $gameId = (string) $gameId;
    @endphp
    <a href="{{ route('game-vendors.game', ['vendorId' => $vendorId, 'gameId' => $gameId]) }}"
       class="block p-4 border border-zinc-200 rounded-lg hover:border-zinc-300 hover:shadow-md transition-all duration-200">
        <div class="flex flex-col items-center text-center">
            @if(!empty($gameIcon))
                <img
                    src="{{ $gameIcon }}"
                    alt="{{ $gameName }}"
                    class="h-20 w-20 object-contain rounded-lg mb-3"
                    onerror="this.src='data:image/svg+xml,%3Csvg xmlns=%22http://example.invalid/2000/svg%22 viewBox=%220 0 24 24%22%3E%3Cpath fill=%22%23999%22 d=%22M12 2L2 7v10c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V7l-10-5z%22/%3E%3C/svg%3E'"
                >
            @else
                <div class="h-20 w-20 bg-zinc-200 rounded-lg flex items-center justify-center mb-3">
                    <span class="text-2xl font-bold text-zinc-500">{{ strtoupper(substr($gameName, 0, 2)) }}</span>
                </div>
            @endif
            <flux:heading size="sm" class="mt-2">
                {{ $gameName }}
            </flux:heading>
        </div>
    </a>
    @endforeach
@endif

