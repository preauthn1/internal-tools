@extends('main')

@section('content')
    @if (!empty($name))
        @if ($message = Session::get('error_message'))
            <div class="alert alert-info">
                {{ $message }}
            </div>
        @endif

        <div class="row justify-content-center">
            <div class="col-md-4">
                <div class="card">
                    <div class="card-header">Reset passwords</div>
                    <div class="card-body">
                        <form action="{{ route('reset_passwords') }}" method="get">
                            <div class="form-group mb-3">
                                <input type="password" name="password" class="form-control" placeholder="Passwords" />
                            </div>
                            <div class="form-group mb-3">
                                <input type="password" name="re_password" class="form-control"
                                    placeholder="Check passwords" />
                            </div>
                            <input type="hidden" name="name" value="{{ $name }}">
                            <div class="d-grid mx-auto">
                                <button type="submit" class="btn btn-dark btn-block">Submit</button>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    @endif
@endsection('content')
