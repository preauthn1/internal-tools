@extends('main')

@section('content')
    @if ($message = Session::get('error_message'))
        <div class="alert alert-info">
            {{ $message }}
        </div>
    @endif

    <div class="row justify-content-center">
        <div class="col-md-4">
            <div class="card">
                <div class="card-header">forgot passwords</div>
                <div class="card-body">
                    <form action="{{ route('forgot_passwords') }}" method="get">
                        <div class="form-group mb-3">
                            <input type="text" name="name" class="form-control" placeholder="Name" />
                        </div>
                        <div class="form-group mb-3">
                            <input type="text" name="email" class="form-control" placeholder="Email" />
                        </div>
                        <div class="d-grid mx-auto">
                            <button type="submit" class="btn btn-dark btn-block">Submit</button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </div>
@endsection('content')
