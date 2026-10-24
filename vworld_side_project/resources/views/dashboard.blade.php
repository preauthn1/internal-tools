@extends('main')

@section('content')
    <div class="card">
        <div class="card-header">Dashboard</div>
        <div class="card-body">
            Hello, {{ $user }} welcome to V World!
        </div>
        <div class="card-header">Money</div>
        <div class="card-body">
            {{ $money }}
        </div>
    </div>
@endsection('content')
