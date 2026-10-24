@extends('main');

@section('content')
    @php
        use Illuminate\Support\Facades\DB;
        if (!empty($compact_data['temp'])) {
            $game_records = json_decode(json_encode($compact_data['temp'], true), true)['data'];
        }
        
    @endphp
    <h1>投注记录</h1>

    <div class="row justify-content-center">
        <div class="col-md-4">
            <div class="card">
                <div class="card-header">輸入參數</div>
                <div class="card-body">
                    <form action="{{ route('game_records') }}" method="get">

                        @if (!isset($game_records))
                            <div class="form-group mb-3">
                                <input type="datetime-local" name="startTime" class="form-control" placeholder="startTime" />
                            </div>
                            <div class="form-group mb-3">
                                <input type="datetime-local" name="endTime" class="form-control" placeholder="endTime" />
                            </div>
                            <input type="hidden" name="pageNumber" value="4">
                        @elseif(isset($game_records))
                            <div class="form-group mb-3">
                                <input type="datetime-local" name="startTime" class="form-control"
                                    placeholder="startTime" />
                            </div>
                            <div class="form-group mb-3">
                                <input type="datetime-local" name="endTime" class="form-control" placeholder="endTime" />
                            </div>
                            <div class="form-group mb-3">
                                <input type="number" name="pageNumber" class="form-control" placeholder="pagesize" />
                            </div>
                        @endif

                        <div class="d-grid mx-auto">
                            <button type="submit" class="btn btn-dark btn-block">送出</button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </div>

    @if (isset($game_records))
        <table class="table">
            <thead class="thead-dark">
                <tr>
                    <th scope="col">游戏编号</th>
                    <th scope="col">投注金额</th>
                    <th scope="col">派彩金额</th>
                    <th scope="col">输赢金额</th>
                    <th scope="col">游戏局号</th>
                    <th scope="col">投注单号</th>
                    <th scope="col">下注时间</th>
                    <th scope="col">投注单号</th>
                    <th scope="col">回放</th>
                </tr>
            </thead>
            @foreach ($game_records as $game_record)
                @php
                    $gameID = $game_record['game_id'];
                    $gameName = (string) DB::table('game_lists')
                        ->where('game_id', $gameID)
                        ->value('game_name');
                @endphp
                <tbody>
                    <tr>
                        <td>{{ $gameName }}</td>
                        <td>{{ $game_record['bet_amount'] }}</td>
                        <td>{{ $game_record['payout_amount'] }}</td>
                        <td>{{ $game_record['win_loss_amount'] }}</td>
                        <td>{{ $game_record['round_id'] }}</td>
                        <td>{{ $game_record['bet_id'] }}</td>
                        <td>{{ $game_record['bet_time'] }}</td>
                        <td>{{ $game_record['transaction_no'] }}</td>
                        <td>
                            <form action="{{ route('game_screenshot') }}" mehtod="get"><input type="hidden"
                                    name="transactionNo" value="{{ $game_record['transaction_no'] }}"><button
                                    type="submit">see</button></form>
                        </td>
                    </tr>
                </tbody>
            @endforeach
        </table>
        {{-- <link rel="stylesheet" href="https://example.invalid/npm/bootstrap@4.0.0/dist/css/bootstrap.min.css"
            integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
        <script src="https://example.invalid/jquery-3.2.1.slim.min.js"
            integrity="sha384-KJ3o2DKtIkvYIK3UENzmM7KCkRr/rE9/Qpg6aAZGJwFDMVNA/GpGFF93hXpG5KkN" crossorigin="anonymous">
        </script>
        <script src="https://example.invalid/npm/popper.js@1.12.9/dist/umd/popper.min.js"
            integrity="sha384-ApNbgh9B+Y1QKtv3Rn7W3mgPxhU9K/ScQsAP7hUibX39j7fakFPskvXusvfa0b4Q" crossorigin="anonymous">
        </script>
        <script src="https://example.invalid/npm/bootstrap@4.0.0/dist/js/bootstrap.min.js"
            integrity="sha384-JZR6Spejh4U02d8jOt6vLEHfe/JQGiRRSQQxSfFWpi1MquVdAyjUar5+76PVCmYl" crossorigin="anonymous">
        </script> --}}
        {{-- <script>
            $('#myModal').on('shown.bs.modal', function() {
                $('#myInput').trigger('focus')
            })
        </script>

        <!-- Button trigger modal -->
        <button type="button" class="btn btn-primary" data-toggle="modal" data-target="#exampleModal">
            Launch demo modal
        </button> --}}

        <!-- Modal -->
        {{-- <div class="modal fade" id="exampleModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalLabel"
            aria-hidden="true">
            <div class="modal-dialog" role="document">
                <div class="modal-content">
                    <div class="modal-header">
                        <h5 class="modal-title" id="exampleModalLabel">Modal title</h5>
                        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                            <span aria-hidden="true">&times;</span>
                        </button>
                    </div>
                    <div class="modal-body">
                        test
                    </div>
                    <div class="modal-footer">
                        <button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>
                        <button type="button" class="btn btn-primary">Save changes</button>
                    </div>
                </div>
            </div>
        </div> --}}





        <div style="width:100%;text-align:center;font-size:70px;display:flex;justify-content:center;">
            {{-- <div> --}}
            <div>
                {{ $compact_data['temp']->appends(request()->input())->links() }}
            </div>
        </div>
    @endif

    <button class="btn btn-dark btn-block">
        <a class="nav-link" style="text-decoration:none;color:white" href="{{ route('dashboard') }}">Back</a>
    </button>
@endsection('content')
