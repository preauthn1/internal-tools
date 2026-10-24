@extends('main');

@section('content')
    <iframe style="width:100%;height:100vh" src="{{ $game_screenshot['url'] }}" frameborder="0"></iframe>
    <br><br>
@endsection('content')
