@extends('example.invalid')
<script src="//example.invalid/ajax/libs/jquery/1.9.1/jquery.min.js"></script>
<link rel="stylesheet" href="https://example.invalid/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
<!-- <script src="https://example.invalid/jquery-3.3.1.slim.min.js" integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo" crossorigin="anonymous"></script> 無法順利使用ajax的功能-->
<script src="https://example.invalid/ajax/libs/popper.js/1.14.7/umd/popper.min.js" integrity="sha384-UO2eT0CpHqdSJQ6hJty5KVphtPhzWj9WO1clHTMGa3JDZwrnQq4sF86dIHNDz0W1" crossorigin="anonymous"></script>
<script src="https://example.invalid/bootstrap/4.3.1/js/bootstrap.min.js" integrity="sha384-JjSmVgyd0p3pXB1rRibZUAYoIIy6OrQ6VrjIEaFf/nJGzIxFDsf4x0xIM+B07jRM" crossorigin="anonymous"></script>
@section('content')
<div class="container-fluid">
    <div class="row justify-content-center">
        <div class="col-md-11">
            <div class="card">
                <div class="card-body">
                    @if (session('status'))
                    <div class="alert alert-success" role="alert">
                        {{ session('status') }}
                    </div>
                    @endif
                    @if(Auth::check())
                    玩家名稱：{{ Auth::user()->name }}
                    <div class="" aria-labelledby="navbarDropdown">
                        <a class="" href="{{ route('logout') }}" onclick="event.preventDefault();
                            document.getElementById('logout-form').submit();">
                            {{ ('登出') }}
                        </a>
                        <form id="logout-form" action="{{ route('logout') }}" method="POST" class="d-none">
                            @csrf
                        </form>
                    </div>
                    遊戲點數：
                        <span id='point' style="color:#FF0000">{{ $point }}</span>點
                        <span id='error' style="color:#FF0000"></span><br>
                    <input id="transfer_in" class="btn btn-primary" type="button" value="轉入遊戲">
                    <input id="transfer_out" class="btn btn-primary" type="button" value="轉回錢包">

                    <!-- Button trigger modal -->
                    <button id="game_record" class="btn btn-primary" data-toggle="modal" data-target="#gameRecordModalCenter">遊戲紀錄</button>
                    <!-- Modal -->
                    <div class="modal fade" id="gameRecordModalCenter" tabindex="-1" role="dialog" aria-hidden="true">
                        <div class="modal-dialog modal-dialog-centered modal-lg" role="document">
                            <div class="modal-content">
                                <div class="modal-header">
                                    <h5 class="modal-title">遊戲紀錄頁面</h5>
                                    <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                                        <span aria-hidden="true">&times;</span>
                                    </button>
                                </div>
                                <div class="text-center">
                                    <input id="game_record_sdata" class="" type="date">
                                    <input id="game_record_edata" class="" type="date">
                                    <input id="game_record_name" class="" type="" value="">
                                    <input id="search" class="btn btn-primary" type="button" value="查詢">
                                </div>
                                <div id="game_record_body" class="modal-body align-self-center">

                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- Button trigger modal -->
                    <button  id="transfer_point" class="btn btn-primary" data-toggle="modal" data-target="#transferPointModalCenter">轉帳紀錄</button>
                    <!-- Modal -->
                    <div class="modal fade" id="transferPointModalCenter" tabindex="-1" role="dialog" aria-hidden="true">
                        <div class="modal-dialog modal-dialog-centered modal-lg" role="document">
                            <div class="modal-content">
                                <div class="modal-header">
                                    <h5 class="modal-title">轉帳紀錄頁面</h5>
                                    <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                                        <span aria-hidden="true">&times;</span>
                                    </button>
                                </div>
                                <div id="transfer_point_body" class="modal-body align-self-center">

                                </div>
                            </div>
                        </div>
                    </div>
                    <br><br>
                    <script>
                        function getDate(date) {
                            var yyyy = date.getFullYear();
                            var MM = (date.getMonth() + 1) >= 10 ? (date.getMonth()+1) : ("0" + (date.getMonth() + 1));
                            var dd = date.getDate() < 10 ? ("0" + date.getDate()) : date.getDate();
                            var date = yyyy + "-" + MM + "-" + dd;
                            return date;
                        }
                        function getLastMonthDate() {
                            var lastMonth = new Date(new Date().getTime() - 30 * 24 * 60 * 60 * 1000);
                            var date = getDate(lastMonth);
                            return date;
                        }
                        function getTodayDate() {
                            var today = new Date();
                            var date = getDate(today);
                            return date;
                        }
                        $('#game_record_sdata').attr("value", getLastMonthDate());
                        $('#game_record_edata').attr("value", getTodayDate());

                        $('#transfer_in').click(function() {
                            $.ajaxSetup({
                                headers: {
                                    'X-CSRF-TOKEN': $('meta[name="csrf-token"]').attr('content')
                                }
                            });
                            $.ajax({
                                type: 'POST',
                                url: '{{route('transfer_in')}}',
                                success: function(data) {
                                    // console.log(data);
                                    if(data["state"]===0) {
                                        $('#point').html(data["data"]);
                                        $('#error').empty();
                                    } else {
                                        // $('#point').html(data["data"]);
                                        $('#error').empty();
                                        $('#error').html(data["message"]);
                                    }
                                },
                                // error: function(data) {
                                //     console.log(data);
                                //     $('#error').empty();
                                //     $('#error').html(data["message"]);
                                // },
                            });
                        })

                        $('#transfer_out').click(function() {
                            $.ajaxSetup({
                                headers: {
                                    'X-CSRF-TOKEN': $('meta[name="csrf-token"]').attr('content')
                                }
                            });
                            $.ajax({
                                type: 'POST',
                                url: '{{route('transfer_out')}}',
                                success: function(data) {
                                    // console.log(data);
                                    if(data["state"]===0) {
                                        $('#point').html(data["data"]);
                                        $('#error').empty();
                                    } else {
                                        // $('#point').html(data["data"]);
                                        $('#error').empty();
                                        $('#error').html(data["message"]);
                                    }
                                },
                                // error: function(data) {
                                //     // console.log(data);
                                //     $('#error').empty();
                                //     $('#error').html(data);
                                // },
                            });
                        })

                        $(document).ready(function() {
                            $("#search, #game_record").click(function() {
                                var sdata = $("#game_record_sdata").val();
                                var edata = $("#game_record_edata").val();
                                var name = $("#game_record_name").val();
                                $.ajax({
                                    url:"/game_records",
                                    data:{
                                        sdata:sdata,
                                        edata:edata,
                                        name:name,
                                    }
                                }).done(gameRecordCallback);
                            })
                            var gameRecordCallback = function(data) {
                                // console.log(data);
                                $('#game_record_body').html(data);
                                // $('.pagination').addClass('justify-content-center');
                                //將html整個塞到主頁面裡面
                                $('#game_record_body .pagination .page-link').each(function(i, el) {
                                    $(el).attr("old-href", $(el).attr("href"));
                                    //將href複製到新的old-href上面
                                    $(el).removeAttr("href");
                                    //移除href以免它點選後跳頁，
                                })
                                $('#game_record_body .pagination .page-link').click(function() {
                                    var sdata = $("#game_record_sdata").val();
                                    var edata = $("#game_record_edata").val();
                                    var name = $("#game_record_name").val();
                                    //點選後，將呼叫新的url，塞到新的html來
                                    var url = $(this).attr('old-href');
                                    if (url) {
                                        $.ajax({
                                        url:url,
                                        data:{
                                                sdata:sdata,
                                                edata:edata,
                                                name:name,
                                            }
                                        }).done(gameRecordCallback);
                                    }
                                })
                            };

                            $('#transfer_point').click(function() {
                                $.ajax("/transfer_points").done(transferPointCallback);
                            })
                            var transferPointCallback = function(data) {
                                $('#transfer_point_body').html(data);
                                $('#transfer_point_body .pagination .page-link').each(function(i, el) {
                                    $(el).attr("old-href", $(el).attr("href"));
                                    $(el).removeAttr("href");
                                })
                                $('#transfer_point_body .pagination .page-link').click(function() {
                                    var url = $(this).attr('old-href');
                                    if (url) {
                                        $.ajax(url).done(transferPointCallback);
                                    }
                                })
                            };
                        })
                    </script>
                    @else
                    <!-- Button trigger modal -->
                    <button type="button" class="btn btn-primary" data-toggle="modal" data-target="#registerModalCenter">註冊</button>
                    <!-- Modal -->
                    <div class="modal fade" id="registerModalCenter" tabindex="-1" role="dialog" aria-hidden="true">
                        <div class="modal-dialog modal-dialog-centered" role="document">
                            <div class="modal-content">
                                <div class="modal-header">
                                    <h5 class="modal-title">註冊頁面</h5>
                                    <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                                        <span aria-hidden="true">&times;</span>
                                    </button>
                                </div>
                                <div class="modal-body">
                                    <form method="POST" action="{{ route('register') }}">
                                        @csrf
                                        <div class="row mb-3">
                                            <label for="name" class="col-md-4 col-form-label text-md-end">{{ __('Name') }}</label>
                                            <div class="col-md-6">
                                                <input id="name" type="text" class="form-control @error('name') is-invalid @enderror" name="name" value="{{ old('name') }}" required autocomplete="name" autofocus>
                                                @error('name')
                                                <span class="invalid-feedback" role="alert">
                                                    <strong>{{ $message }}</strong>
                                                </span>
                                                @enderror
                                            </div>
                                        </div>
                                        <div class="row mb-3">
                                            <label for="email" class="col-md-4 col-form-label text-md-end">{{ __('Email Address') }}</label>
                                            <div class="col-md-6">
                                                <input id="email" type="email" class="form-control @error('email') is-invalid @enderror" name="email" value="{{ old('email') }}" required autocomplete="email">
                                                @error('email')
                                                <span class="invalid-feedback" role="alert">
                                                    <strong>{{ $message }}</strong>
                                                </span>
                                                @enderror
                                            </div>
                                        </div>
                                        <div class="row mb-3">
                                            <label for="password" class="col-md-4 col-form-label text-md-end">{{ __('Password') }}</label>
                                            <div class="col-md-6">
                                                <input id="password" type="password" class="form-control @error('password') is-invalid @enderror" name="password" required autocomplete="new-password">
                                                @error('password')
                                                <span class="invalid-feedback" role="alert">
                                                    <strong>{{ $message }}</strong>
                                                </span>
                                                @enderror
                                            </div>
                                        </div>
                                        <div class="row mb-3">
                                            <label for="password-confirm" class="col-md-4 col-form-label text-md-end">{{ __('Confirm Password') }}</label>
                                            <div class="col-md-6">
                                                <input id="password-confirm" type="password" class="form-control" name="password_confirmation" required autocomplete="new-password">
                                            </div>
                                        </div>
                                        <div class="row mb-0">
                                            <div class="col-md-6 offset-md-4">
                                                <button type="submit" class="btn btn-primary">
                                                    {{ __('Register') }}
                                                </button>
                                            </div>
                                        </div>
                                    </form>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- Button trigger modal -->
                    <button type="button" class="btn btn-primary" data-toggle="modal" data-target="#loginModalCenter">登入</button>
                    <!-- Modal -->
                    <div class="modal fade" id="loginModalCenter" tabindex="-1" role="dialog" aria-hidden="true">
                        <div class="modal-dialog modal-dialog-centered" role="document">
                            <div class="modal-content">
                                <div class="modal-header">
                                    <h5 class="modal-title">登入頁面</h5>
                                    <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                                        <span aria-hidden="true">&times;</span>
                                    </button>
                                </div>
                                <div class="modal-body">
                                    <form method="POST" action="{{ route('login') }}">
                                        @csrf
                                        <div class="row mb-3">
                                            <label for="email" class="col-md-4 col-form-label text-md-end">{{ __('Email Address') }}</label>
                                            <div class="col-md-6">
                                                <input id="email" type="email" class="form-control @error('email') is-invalid @enderror" name="email" value="{{ old('email') }}" required autocomplete="email" autofocus>
                                                @error('email')
                                                <span class="invalid-feedback" role="alert">
                                                    <strong>{{ $message }}</strong>
                                                </span>
                                                @enderror
                                            </div>
                                        </div>
                                        <div class="row mb-3">
                                            <label for="password" class="col-md-4 col-form-label text-md-end">{{ __('Password') }}</label>
                                            <div class="col-md-6">
                                                <input id="password" type="password" class="form-control @error('password') is-invalid @enderror" name="password" required autocomplete="current-password">
                                                @error('password')
                                                <span class="invalid-feedback" role="alert">
                                                    <strong>{{ $message }}</strong>
                                                </span>
                                                @enderror
                                            </div>
                                        </div>
                                        <div class="row mb-3">
                                            <div class="col-md-6 offset-md-4">
                                                <div class="form-check">
                                                    <input class="form-check-input" type="checkbox" name="remember" id="remember" {{ old('remember') ? 'checked' : '' }}>
                                                    <label class="form-check-label" for="remember">
                                                        {{ __('Remember Me') }}
                                                    </label>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="row mb-0">
                                            <div class="col-md-8 offset-md-4">
                                                <button type="submit" class="btn btn-primary">
                                                    {{ __('Login') }}
                                                </button>
                                                @if (Route::has('password.request'))
                                                <a class="btn btn-link" href="{{ route('password.request') }}">
                                                    {{ __('Forgot Your Password?') }}
                                                </a>
                                                @endif
                                            </div>
                                        </div>
                                    </form>
                                </div>
                            </div>
                        </div>
                    </div>
                    <br><br>
                    @endif
                    <div class="text-center">
                        @foreach($games as $k => $game)
                        <input type="hidden" name="{{$game->gameCode}}" value="{{$game->gameID}}">
                        <img class="{{$game->gameCode}}" src="{{$game->imageUrl}}" width="200">
                        <script>
                            $('.{{$game->gameCode}}').click(function() {
                                var gameID = $("input[name={{$game->gameCode}}]").val();
                                $.ajaxSetup({
                                    headers: {
                                        'X-CSRF-TOKEN': $('meta[name="csrf-token"]').attr('content')
                                    }
                                });
                                $.ajax({
                                    type: 'POST',
                                    url: '{{route('game_link')}}',
                                    data: {
                                        gameID: gameID
                                    },
                                    success: function(data) {
                                        if(data["state"]===0) {
                                            window.open(data["data"]);
                                        } else {
                                            $('#error').html(data["message"]);
                                        }
                                    }
                                });
                            })
                        </script>
                        @endforeach
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>
@endsection