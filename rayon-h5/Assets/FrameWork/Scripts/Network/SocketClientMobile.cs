using UnityEngine;
using System;
using example.invalid;
using example.invalid.Sockets;

public class SocketClientMobile : ISocketClient
{
    private TcpClient client = null;
    private NetworkStream outStream = null;
    private MemoryStream memStream;
    private BinaryReader reader;
    private ConnectType connectType;
    private ConnectState connectState = ConnectState.Other;

    private const int MAX_READ = 8192;
    private byte[] byteBuffer = new byte[MAX_READ];

    WebSocketEntry webSocketEntry;

    public ConnectState MyConnectState
    {
        get
        {
            return connectState;
        }
        private set
        {
            connectState = value;
        }
    }

    // Use this for initialization
    public SocketClientMobile(ConnectType connectType)
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

        //设置socket状态为正在连接状态
        MyConnectState = ConnectState.Connecting;
        
        if (client != null)
        {
            Close();
        }
        try
        {
            // 依造連線模式(IPV4 or IPV6)轉換成對應格式 By Rayon
            string newServerIp = "";
            AddressFamily newAddressFamily = AddressFamily.InterNetwork;
            Ipv6.getIPType(host, port.ToString(), out newServerIp, out newAddressFamily);
            /**
            Debug.Log("IP協定模式 " + newAddressFamily.ToString());
            Debug.Log("IP原本位址 " + host);
            Debug.Log("IP轉換後位址 " + newServerIp);
            **/
            host = newServerIp;
            if (newAddressFamily == AddressFamily.InterNetworkV6)
            {
                client = new TcpClient(AddressFamily.InterNetworkV6);
            }
            else
            {
                client = new TcpClient(AddressFamily.InterNetwork);
            }
            
            // 陸版舊方法
            /*
            IPAddress[] address = Dns.GetHostAddresses(host);
            if (address.Length == 0)
            {
                Debug.LogError("host invalid");
                MyConnectState = ConnectState.Other;
                return;
            }
            if (address[0].AddressFamily == AddressFamily.InterNetworkV6)
            {
                client = new TcpClient(AddressFamily.InterNetworkV6);
            }
            else
            {
                client = new TcpClient(AddressFamily.InterNetwork);
            }
            */
            client.SendTimeout = 1000;
            client.ReceiveTimeout = 1000;
            client.NoDelay = true;
            client.BeginConnect(host, port, new AsyncCallback(OnConnect), null);
        }
        catch (Exception e)
        {
            MyConnectState = ConnectState.Other;
            Debug.LogError(e.Message);
        }
    }

    /// <summary>
    /// 连接上服务器
    /// </summary>
    void OnConnect(IAsyncResult asr)
    {
        if (client.Connected)
        {
            outStream = client.GetStream();
            client.GetStream().BeginRead(byteBuffer, 0, MAX_READ, new AsyncCallback(OnRead), null);
            MyConnectState = ConnectState.Connect;
        }
        else
        {
            Debug.LogError(string.Format("{0}:{1} 无法连接", Host, Port));
            MyConnectState = ConnectState.Other;
            Close();
        }
    }

    /// <summary>
    /// 写数据
    /// </summary>
    void WriteMessage(int mainID, int secondID, byte[] message)
    {
        if (client != null && client.Connected)
        {
            ByteBuffer buffer = new ByteBuffer();
            buffer.WriteInt((message == null ? 0 : message.Length) + 8);
            buffer.WriteInt(mainID);
            buffer.WriteInt(secondID);

            if (message != null)
                buffer.WriteBytes(message);

            buffer.Flush();
            byte[] payload = buffer.ToBytes();
            // PrintBytes(payload);

            outStream.BeginWrite(payload, 0, payload.Length, new AsyncCallback(OnWrite), null);
            buffer.Close();
        }
        else
        {
            Debug.LogError(string.Format("{0}:{1} 未连接,无法发送消息", Host, Port));
        }
    }

    /// <summary>
    /// 读取消息
    /// </summary>
    void OnRead(IAsyncResult asr)
    {
        if (client == null || !client.Connected)
            return;

        int bytesRead = 0;
        try
        {
            lock (client.GetStream())
            {         //读取字节流到缓冲区
                bytesRead = client.GetStream().EndRead(asr);
            }
            if (bytesRead < 1)
            {
                //包尺寸有问题，断线处理
                Disconnected("bytesRead < 1");
                return;
            }
            OnReceive(byteBuffer, bytesRead);   //分析数据包内容，抛给逻辑层
            lock (client.GetStream())
            {         //分析完，再次监听服务器发过来的新消息
                Array.Clear(byteBuffer, 0, byteBuffer.Length);   //清空数组
                client.GetStream().BeginRead(byteBuffer, 0, MAX_READ, new AsyncCallback(OnRead), null);
            }
        }
        catch (Exception ex)
        {
            //PrintBytes();
            Disconnected(ex.Message);
        }
    }

    /// <summary>
    /// 异常断开连接(会触发断线重连)
    /// </summary>
    public void Disconnected(string msg)
    {
        Debug.LogError("连接异常，" + Host + ":" + Port + "的连接被关闭:>" + msg);
        NetworkManager.AddEvent(this, 0, (int)connectType, null);
        MyConnectState = ConnectState.Other;
        Close();                        //关掉客户端链接
    }

    /// <summary>
    /// 打印字节
    /// </summary>
    /// <param name="bytes"></param>
    void PrintBytes(byte[] bytes)
    {
        string returnStr = string.Empty;
        for (int i = 0; i < bytes.Length; i++)
        {
            returnStr += bytes[i].ToString("X2");
        }
    }

    /// <summary>
    /// 向链接写入数据流
    /// </summary>
    void OnWrite(IAsyncResult r)
    {
        try
        {
            outStream.EndWrite(r);
        }
        catch (Exception ex)
        {
            Debug.LogError("OnWrite--->>>" + ex.Message);
        }
    }

    /// <summary>
    /// 接收到消息
    /// </summary>
    void OnReceive(byte[] bytes, int length)
    {
        memStream.Seek(0, SeekOrigin.End);
        memStream.Write(bytes, 0, length);
        //Reset to beginning
        memStream.Seek(0, SeekOrigin.Begin);

        while (RemainingBytes() > 4)
        {
            ByteBuffer buffer = new ByteBuffer(reader.ReadBytes(4));    // 前面4个字节表示消息长度
            int messageLen = buffer.ReadInt();
            
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
        
        NetworkManager.AddEvent(this, mainId, secondID, buffer.ReadBytes((int)ms.Length - 8));
    }

    public void Close()
    {
        MyConnectState = ConnectState.Other;

        if (client != null)
        {
            if (client != null)
            {
                client.Close();
                memStream.SetLength(0);
                client = null;
            }
        }
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
        if (client == null || !client.Connected)
        {
            return "disconnected";
        }

        return client.Client.RemoteEndPoint.ToString();
    }
}