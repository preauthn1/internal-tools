using UnityEngine;
using System;
using System.Collections;
using System.Collections.Generic;

public struct NetMessage
{
    public int MainID;
    public int SecondID;
    public byte[] Content;
}

public class NetworkManager : Singleton<NetworkManager>
{
    private static bool s_ShowLog = false;

    public ISocketClient HallSocket
    {
        get { return hallSocket; }
    }

    public ISocketClient RoomSocket
    {
        get { return roomSocket; }
    }

    private ISocketClient hallSocket;
    private ISocketClient roomSocket;

    // 消息处理挂起 用于处理实例化对象或加载新场景的过程中，分发出去的消息可能没有接收到的情况
    private bool m_IsReceiveMessagePause = false;

    public bool GetPause()
    {
        return m_IsReceiveMessagePause;
    }

    public void SetReceiveMessagePause(bool isPause)
    {
        m_IsReceiveMessagePause = isPause;
    }

    //游戏消息同步变量
    static readonly object m_lockObject = new object();
    //心跳消息同步变量
    static readonly object m_beatLockObject = new object();
    //游戏消息
    static Queue<NetMessage> mEvents = new Queue<NetMessage>();
    //心跳消息
    static Queue<ISocketClient> mEventsHeartBeat = new Queue<ISocketClient>();

    //socket 断网消息
    static Queue<NetMessage> mSocketEvents = new Queue<NetMessage>();

    public Action<int, int, byte[]> OnReceiveNetworkMessage = (first, second, data) => { };

    public const int TIMEOUTSECOND = 8;     // Socket 连接超时时间  Hall & Room

    public override void Awake()
    {
        Init();
        base.Awake();
    }

    void Init()
    {
        switch(Application.platform)
        {
            case RuntimePlatform.WebGLPlayer:
                Debug.Log("--- SocketClient - WebGL ---");
                hallSocket = new SocketClientWebGL(ConnectType.Hall);
                roomSocket = new SocketClientWebGL(ConnectType.Room);
                break;

            default:
                Debug.Log("--- SocketClient - Mobile ---");
                hallSocket = new SocketClientMobile(ConnectType.Hall);
                roomSocket = new SocketClientMobile(ConnectType.Room);
                break;
        }
        hallSocket.OnRegister();
        roomSocket.OnRegister();
    }

    public void ReLaunch()
    {
        m_IsReceiveMessagePause = false;
        mEvents = new Queue<NetMessage>();
        mSocketEvents = new Queue<NetMessage>();
        OnReceiveNetworkMessage = (first, second, data) => { };
    }
    ///------------------------------------------------------------------------------------
    public static void AddEvent(ISocketClient socket, int mainID, int secondID, byte[] data)
    {
        if (s_ShowLog == true)
            Debug.Log(string.Format("接收到消息:{0},{1}", mainID, secondID));

        //心跳消息
        if (HeartBeat.IsHeartBeatMsg(mainID, secondID))
        {
            lock (m_beatLockObject)
            {
                mEventsHeartBeat.Enqueue(socket);
            }
        }
        else
        {
            lock (m_lockObject)
            {
                if (mainID == 0 && (secondID == 0 || secondID == 1))
                    mSocketEvents.Enqueue(new NetMessage() { MainID = mainID, SecondID = secondID, Content = data });
                else
                {
                    if (s_ShowLog == true)
                        Debug.Log(string.Format("将消息写到队列中:{0},{1}", mainID, secondID));
                    mEvents.Enqueue(new NetMessage() { MainID = mainID, SecondID = secondID, Content = data });
                }
            }
        }
    }

    //移除符合条件的消息队列
    public void RemoveEvents(Func<NetMessage, bool> fun)
    {
        Queue<NetMessage> newEvents = new Queue<NetMessage>();
        foreach (var s in mEvents)
        {
            if (!fun(s))
            {
                newEvents.Enqueue(s);
            }
        }
        mEvents = newEvents;
    }

    /// <summary>
    /// 交给Command，这里不想关心发给谁。
    /// </summary>
    void Update()
    {
        //处理心跳消息
        if (mEventsHeartBeat.Count > 0)
        {
            while (mEventsHeartBeat.Count > 0)
            {
                ISocketClient socket;

                lock (m_beatLockObject)
                    socket = mEventsHeartBeat.Dequeue();

                HeartBeatManager.Instance.OnHeartBeatReturnMsg(socket);

            }
        }

        //处理LUA消息
        if (mEvents.Count > 0)
        {
            while (mEvents.Count > 0 && !m_IsReceiveMessagePause)
            {
                NetMessage message;

                lock (m_lockObject)
                    message = mEvents.Dequeue();

                if (OnReceiveNetworkMessage != null)
                {
                    if (s_ShowLog == true)
                        Debug.LogFormat("触发网络消息:{0},{1}", message.MainID, message.SecondID);
                    OnReceiveNetworkMessage(message.MainID, message.SecondID, message.Content);
                }
            }
        }

        // 断网消息不受消息暂停设置影响
        if (mSocketEvents.Count > 0)
        {
            while (mSocketEvents.Count > 0)
            {
                NetMessage message;

                lock (m_lockObject)
                    message = mSocketEvents.Dequeue();

                if (OnReceiveNetworkMessage != null)
                    OnReceiveNetworkMessage(message.MainID, message.SecondID, message.Content);
            }
        }
    }

    /// <summary>
    /// 发送链接请求
    /// </summary>
    public IEnumerator SendHallConnect(string host, int port)
    {
        // 是否没有必要去连接
        bool noNeedConnect = hallSocket.Host == host
                            && hallSocket.Port == port
                            && (hallSocket.MyConnectState == ConnectState.Connecting || hallSocket.MyConnectState == ConnectState.Connect);

        if (noNeedConnect)
        {
            Debug.Log("大厅正在连接或者已经连接，不需要重新连。");
        }
        else
        {
            DateTime dt = DateTime.Now;
            hallSocket.SendConnect(host, port);

            while (hallSocket.MyConnectState == ConnectState.Connecting)
            {
                if (dt.AddSeconds(TIMEOUTSECOND) < DateTime.Now)
                {
                    try
                    {
                        print("客户端判断连接超时，断开大厅连接。");
                        hallSocket.Close();
                    }
                    catch (Exception ex)
                    {
                        Debug.LogError(ex);
                    }
                }
                else
                    yield return 1;
            }
        }
    }

    // 请求连接到房间
    public IEnumerator SendRoomConnect(string host, int port)
    {
        // 是否没有必要去连接
        bool noNeedConnect = roomSocket.Host == host
                            && roomSocket.Port == port
                            && (roomSocket.MyConnectState == ConnectState.Connecting || roomSocket.MyConnectState == ConnectState.Connect);

        if (noNeedConnect)
        {
            Debug.Log("游戏正在连接或者已经连接，不需要重新连。");
        }
        else
        {
            DateTime dt = DateTime.Now;
            roomSocket.SendConnect(host, port);

            while (roomSocket.MyConnectState == ConnectState.Connecting)
            {
                if (dt.AddSeconds(TIMEOUTSECOND) < DateTime.Now)
                {
                    try
                    {
                        print("客户端判断连接超时，断开游戏连接。");
                        roomSocket.Close();
                    }
                    catch (Exception ex)
                    {
                        Debug.LogError(ex);
                    }
                }
                else
                    yield return 1;
            }
        }
    }

    /// <summary>
    /// 发送SOCKET消息
    /// </summary>
    public void SendHallMessage(int mainID, int secondID, byte[] buffer)
    {
        hallSocket.SendMessage(mainID, secondID, buffer);
    }

    /// <summary>
    /// 发送SOCKET消息
    /// </summary>
    public void SendRoomMessage(int mainID, int secondID, byte[] buffer)
    {
        roomSocket.SendMessage(mainID, secondID, buffer);
    }

    /// <summary>
    /// 析构函数
    /// </summary>
    protected override void OnDestroy()
    {
        base.OnDestroy();
        hallSocket.OnRemove();
        roomSocket.OnRemove();
    }

    public void AddReceiveNetworkMessage(Action<int, int, byte[]> action)
    {
        OnReceiveNetworkMessage += action;
    }
}