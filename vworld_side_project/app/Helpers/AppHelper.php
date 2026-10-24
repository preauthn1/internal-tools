<?php
use GuzzleHttp\Client;

if (!function_exists('arrang_data')) {
    function arrange_data($data)
    {
        $tests = explode(",", env('GUZZLE_PARAMETERS'));
        foreach ($tests as $test) {
            $test_finals[] = explode(' => ', $test);
        }
        foreach ($test_finals as $test_final) {
            $result[$test_final[0]] = $test_final[1];
        }
        return array_merge($result, $data);
    }
}


if (!function_exists('guzzle_get')) {
    function guzzle_get($data, $url)
    {
        $client = new Client();
        $params = [
            'query' => arrange_data($data)
        ];
        $response = $client->request('GET', env('GUZZLE_TEST_URL') . $url, $params);
        return check_status($response);
    }
}
if (!function_exists('guzzle_post')) {
    function guzzle_post($data, $url)
    {
        $client = new Client([
            'headers' => ['Content-Type' => 'application/json']
        ]);
        $response = $client->request('POST', env('GUZZLE_TEST_URL') . $url, [
            'body' => json_encode(arrange_data($data))
        ]);
        return check_status($response);
    }
}

if (!function_exists('check_status')) {
    function check_status($response): null|array
    {
        if ($response->getStatusCode() == 200)
            if (!(json_decode($response->getbody(), true)['code']))
                return json_decode($response->getbody(), true)['data'];
            else
                return null;
        else
            return null;
    }
}
?>