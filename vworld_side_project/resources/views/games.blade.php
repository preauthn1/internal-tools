@php
    use GuzzleHttp\Client;
@endphp
@extends('main');

@section('content')
    <html>

    <head>
        <style>
            .container {
                display: flex;
                width: 100%;
                flex-wrap: wrap;
                justify-content: center;
                align-items: center;
                margin: 0 auto;
            }
        </style>
    </head>

    <body>
        <div class="container">
            @php
                // dd($game_lists);
            @endphp
            @foreach ($game_lists as $game_list)
                <form action="{{ route('game_link') }}" method="get">
                    <div class="card" style="width: 200px;">
                        <img src="{{ $game_list->image_url }}" class="card-img-top" alt="...">
                        <div class="card-body">
                            <h5 class="card-title">{{ $game_list->game_name }}</h5>
                            <input type="hidden" name="gameID" value="{{ $game_list->game_id }}">
                            <div class="d-grid mx-auto">
                                <button type="submit" class="btn btn-dark btn-block">GO</button>
                            </div>
                        </div>
                    </div>
                </form>
            @endforeach

            <button class="btn btn-dark btn-block">
                <a class="nav-link" style="text-decoration:none;color:white" href="{{ route('dashboard') }}">Back</a>
            </button><br>
        </div>
    </body>

    </html>
@endsection('content')
