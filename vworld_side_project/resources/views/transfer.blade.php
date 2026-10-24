@extends('main')

@section('content')
    @php
        if (!empty($compact_data['temp'])) {
            $transfer_records = json_decode(json_encode($compact_data['temp'], true), true)['data'];
        }
        
    @endphp
    <h1>交易紀錄</h1>
    <div class="row justify-content-center">
        <div class="col-md-4">
            <div class="card">
                <div class="card-header">輸入參數</div>
                <div class="card-body">
                    <form action="{{ route('transfer') }}" method="get">
                        @if (!isset($transfer_records))
                            <div class="form-group mb-3">
                                <input type="datetime-local" name="startTime" class="form-control" placeholder="startTime" />
                            </div>
                            <div class="form-group mb-3">
                                <input type="datetime-local" name="endTime" class="form-control" placeholder="endTime" />
                            </div>
                            <input type="hidden" name="pageNumber" value="5">
                        @elseif(isset($transfer_records))
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

    @if (isset($transfer_records))
        <table class="table">
            <thead class="thead-dark">
                <tr>
                    <th scope="col">交易編號</th>
                    <th scope="col">轉帳金額</th>
                    <th scope="col">轉帳型態</th>
                    <th scope="col">轉帳狀態</th>
                    <th scope="col">交易前金額</th>
                    <th scope="col">交易後金額</th>
                    <th scope="col">交易日期</th>
                </tr>
            </thead>
            @foreach ($transfer_records as $transfer_record)
                <tbody>
                    <tr>
                        <td>{{ $transfer_record['transaction_no'] }}</td>
                        <td>{{ $transfer_record['amount'] }}</td>
                        <td>{{ $transfer_record['transfer_type'] }}</td>
                        <td>{{ $transfer_record['status'] }}</td>
                        <td>{{ $transfer_record['before_amount'] }}</td>
                        <td>{{ $transfer_record['after_amount'] }}</td>
                        <td>{{ $transfer_record['created_at'] }}</td>
                    </tr>
                </tbody>
            @endforeach
        </table>
        <div style="height:50px;width:100%;text-align:center;font-size:70px;display:flex;justify-content:center;">
            <div>
                {{ $compact_data['temp']->appends(request()->input())->links() }}
            </div>
        </div>
    @endif
    <button class="btn btn-dark btn-block">
        <a class="nav-link" style="text-decoration:none;color:white" href="{{ route('dashboard') }}">Back</a>
    </button>
@endsection('content')
