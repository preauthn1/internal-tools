using System;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityWebSocket;

class WebSocketEntry
{
    public WebSocket socket { get; private set; }
    
    Action<byte[]> mReceiveAction;

    public WebSocketEntry(string address, Action<byte[]> receiveAction)
    {
        socket = new WebSocket(address);
        socket.onOpen += OnOpen;
        socket.onClose += OnClose;
        socket.onMessage += OnReceive;
        socket.onError += OnError;

        mReceiveAction = receiveAction;
    }

    public void Connect()
    {
        if (socket == null
            || socket.readyState == WebSocketState.Open
            || socket.readyState == WebSocketState.Closing)
            return;

        socket.Connect();
        Debug.Log("Connect");
    }

    public void Close()
    {
        if (socket.readyState == WebSocketState.Connecting
            || socket.readyState == WebSocketState.Open)
        {
            socket.Close();
        }
        Debug.Log("Close");
    }

    public void Send(string text)
    {
        if (socket.readyState == WebSocketState.Open)
        {
            socket.Send(text);
            Debug.Log("Send A : " + text);
        }
        Debug.Log("Send B : " + text);
    }

    private void OnOpen(object sender, EventArgs e)
    {
        Debug.Log("OnOpen");
    }

    private void OnClose(object sender, CloseEventArgs e)
    {
        if (e.statusCode != CloseStatusCode.NoStatus && e.statusCode != CloseStatusCode.Normal)
        {
            Debug.Log("OnClose A : " + e.Reason + " " + e.statusCode + "\n");
        }
        else
        {
            Debug.Log("OnClose B");
        }
    }

    private void OnReceive(object sender, MessageEventArgs e)
    {
        Debug.Log("OnReceive : " + e.Data + "\n");

        mReceiveAction.Invoke(e.RawData);
    }

    private void OnError(object sender, ErrorEventArgs e)
    {
        Debug.Log("OnError : " + e.Message + "\n");
    }
}