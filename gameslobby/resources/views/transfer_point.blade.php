<div class="container">
    <table class="table table-bordered table-striped">
        <thead class="text-center">
            <th>id</th>
            <th>遊戲玩家</th>
            <th>交易金额</th>
            <th>交易前金额</th>
            <th>交易后金额</th>
            <th>建立日期</th>
            <th>是否成功</th>
        </thead>
        <tbody class="">
            @if (isset($transfer_points))
                @foreach($transfer_points as $k1 => $transfer_point)
                    <tr class="text-center">
                        <td class="border border-gray-900 p-2">{{ $transfer_point->id }}</td>
                        <td class="border border-gray-900 p-2">{{ $transfer_point->name }}</td>
                        <td class="border border-gray-900 p-2">{{ $transfer_point->amount }}</td>
                        <td class="border border-gray-900 p-2">{{ $transfer_point->beforeAmount }}</td>
                        <td class="border border-gray-900 p-2">{{ $transfer_point->afterAmount }}</td>
                        <td class="border border-gray-900 p-2">{{ $transfer_point->created_at }}</td>
                        <td class="border border-gray-900 p-2">{{ $transfer_point->state }}</td>
                    </tr>
                @endforeach
            @endif
        </tbody>
    </table>
</div>
<div class="text-center">共{{ $count }}筆資料</div>
{{ $transfer_points->links() }}
