<div class="container">
    <table class="table table-sm table-bordered table-striped">
        <thead class="text-center">
            <th>id</th>
            <th>遊戲名稱</th>
            <th>投注金额</th>
            <th>贏與輸金额</th>
            <th>建立日期</th>
        </thead>
        <tbody id="">
            @if (isset($game_records))
                @foreach($game_records as $k1 => $game_record)
                    <tr class="text-center">
                        <td class="border border-gray-900 p-2">{{ $game_record->id }}</td>
                        <td class="border border-gray-900 p-2">{{ $game_record->gameName }}</td>
                        <td class="border border-gray-900 p-2">{{ $game_record->betAmount }}</td>
                        <td class="border border-gray-900 p-2">{{ $game_record->winLossAmount }}</td>
                        <td class="border border-gray-900 p-2">{{ $game_record->betTime }}</td>
                    </tr>
                @endforeach
            @endif
        </tbody>
    </table>
</div>
<div class="text-center">共{{ $count }}筆資料</div>
{{ $game_records->links() }}
