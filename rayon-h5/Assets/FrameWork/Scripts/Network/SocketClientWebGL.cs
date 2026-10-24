using UnityEngine;
using System;
using example.invalid;
using example.invalid.Sockets;

public class SocketClientWebGL : ISocketClient
{
    private NetworkStream outStream = null;
    private MemoryStream memStream;
    private BinaryReader reader;
    private ConnectType connectType;
    private const int MAX_READ = 8192;
    private byte[] byteBuffer = new byte[MAX_READ];

    WebSocketEntry webSocketEntry;
    
    public ConnectState MyConnectState
    {
        get
        {
            switch (webSocketEntry.socket.readyState)
            {
                case UnityWebSocket.WebSocketState.Open:
                    return ConnectState.Connect;

                case UnityWebSocket.WebSocketState.Connecting:
                    return ConnectState.Connecting;
            }
            return ConnectState.Other;
        }
    }

    // Use this for initialization
    public SocketClientWebGL(ConnectType connectType)
    {
        this.connectType = connectType;
    }

    /// <summary>
    /// 注册代理
    /// </summary>
    public void OnRegister()
    {
        memStream = new MemoryStream();
        reader = new BinaryReader(memStream);
    }

    /// <summary>
    /// 移除代理
    /// </summary>
    public void OnRemove()
    {
        if (!string.IsNullOrEmpty(Host))
            Debug.Log(Host + ":" + Port + "的连接被移除");

        this.Close();
        reader.Close();
        memStream.Close();
    }

    public string Host
    {
        private set;
        get;
    }

    public int Port
    {
        private set;
        get;
    }

    /// <summary>
    /// 连接服务器
    /// </summary>
    public void SendConnect(string host, int port)
    {
        Host = host;
        Port = port;
        
        try
        {
            webSocketEntry = new WebSocketEntry("ws://" + host + ":" + port, OnReceive);
            webSocketEntry.Connect();
            
            switch (webSocketEntry.socket.readyState)
            {
                case UnityWebSocket.WebSocketState.Closed:
                case UnityWebSocket.WebSocketState.Closing:
                    {
                        Debug.LogError(string.Format("{0}:{1} 无法连接", Host, Port));
                    }
                    break;
            }
        }
        catch (Exception e)
        {
            Debug.LogError(e.Message);
        }
    }

    /// <summary>
    /// 写数据
    /// </summary>
    void WriteMessage(int mainID, int secondID, byte[] message)
    {
        Debug.Log("WriteMessage A : " + mainID + "/" + secondID);
        Debug.Log("WriteMessage B : " + webSocketEntry.socket.readyState);

        switch (webSocketEntry.socket.readyState)
        {
            case UnityWebSocket.WebSocketState.Closed:
            case UnityWebSocket.WebSocketState.Closing:
                {
                    webSocketEntry.Connect();
                    WriteMessage(mainID, secondID, message);
                }
                break;

            case UnityWebSocket.WebSocketState.Connecting:
                {
                    WriteMessage(mainID, secondID, message);
                }
                break;

            case UnityWebSocket.WebSocketState.Open:
                {
                    ByteBuffer buffer = new ByteBuffer();
                    buffer.WriteInt((message == null ? 0 : message.Length) + 8);
                    buffer.WriteInt(mainID);
                    buffer.WriteInt(secondID);

                    if (message != null)
                        buffer.WriteBytes(message);

                    buffer.Flush();
                    byte[] payload = buffer.ToBytes();

                    // 目前測試階段 所以只送協定
                    string data = mainID + "_" + secondID;
                    webSocketEntry.Send(data);
                }
                break;
        }
    }

    /// <summary>
    /// 异常断开连接(会触发断线重连)
    /// </summary>
    public void Disconnected(string msg)
    {
        Debug.LogError("连接异常，" + Host + ":" + Port + "的连接被关闭:>" + msg);
        NetworkManager.AddEvent(this, 0, (int)connectType, null);
        //关掉客户端链接
        Close();
    }
    
    /// <summary>
    /// 剩余的字节
    /// </summary>
    private long RemainingBytes()
    {
        return memStream.Length - memStream.Position;
    }

    /// <summary>
    /// 接收到消息
    /// </summary>
    /// <param name="ms"></param>
    void OnReceivedMessage(MemoryStream ms)
    {
        BinaryReader r = new BinaryReader(ms);
        byte[] message = r.ReadBytes((int)(ms.Length - ms.Position));
        //int msglen = message.Length;

        ByteBuffer buffer = new ByteBuffer(message);
        int mainId = buffer.ReadInt();
        int secondID = buffer.ReadInt();

        Debug.Log("AAA OnReceivedMessage : " + mainId + "/" + secondID);

        NetworkManager.AddEvent(this, mainId, secondID, buffer.ReadBytes((int)ms.Length - 8));
    }

    public void Close()
    {
        webSocketEntry.Close();
    }

    /// <summary>
    /// 发送消息
    /// </summary>
    public void SendMessage(int mainID, int secondID, byte[] buffer)
    {
        WriteMessage(mainID, secondID, buffer);
    }

    public override string ToString()
    {
        if (webSocketEntry.socket.readyState == UnityWebSocket.WebSocketState.Closed ||
            webSocketEntry.socket.readyState == UnityWebSocket.WebSocketState.Closing)
        {
            return "disconnected";
        }
        return webSocketEntry.socket.address;
    }

    /// <summary>
    /// 接收
    /// </summary>
    /// <param name="data"></param>
    void OnReceive(byte[] bytes)
    {
        Debug.Log("OnReceive : " + bytes.Length);

        memStream.Seek(0, SeekOrigin.End);
        memStream.Write(bytes, 0, bytes.Length);
        memStream.Seek(0, SeekOrigin.Begin);

        while (RemainingBytes() > 4)
        {
            ByteBuffer buffer = new ByteBuffer(reader.ReadBytes(4));    // 前面4个字节表示消息长度
            int messageLen = buffer.ReadInt();

            Debug.Log("OnReceive BBB : " + RemainingBytes() + "/" + messageLen);

            if (RemainingBytes() >= messageLen)
            {
                MemoryStream ms = new MemoryStream();
                BinaryWriter writer = new BinaryWriter(ms);
                writer.Write(reader.ReadBytes(messageLen));
                ms.Seek(0, SeekOrigin.Begin);
                OnReceivedMessage(ms);
            }
            else
            {
                //Back up the position 4 bytes
                memStream.Position = memStream.Position - 4;
                break;
            }
        }
        //Create a new stream with any leftover bytes
        byte[] leftover = reader.ReadBytes((int)RemainingBytes());
        memStream.SetLength(0);     //Clear
        memStream.Write(leftover, 0, leftover.Length);

    }
}