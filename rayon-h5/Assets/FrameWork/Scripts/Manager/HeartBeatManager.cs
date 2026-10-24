using System;
using System.Collections.Generic;
using UnityEngine;


//心跳来源
public enum HeartBeatSource
{
    //大厅服务器
    HallServer = 0,
    //游戏服务器
    GameServer = 1
}

public class HeartBeatCfg
{
    //心跳时间间隔
    public Int32 BeatInterval;
    //心跳超时时间
    public Int32 ConnLoseTime;
}

/// <summary>
/// 心跳逻辑
/// 客户端定时发送心跳，服务器返回心跳
/// 如果超过指定的时间未收到服务器返回心跳，则认为连接断开
/// </summary>
public class HeartBeat
{
    private ISocketClient socket;
    private HeartBeatCfg beatCfg;
    private HeartBeatSource bt;

    //心跳一级协议号
    private const Int32 MSGCLASSID_B_COMMID = 1;
    //客户端心跳二级协议号
    private const Int32 MSGCLASSID_S_HEARTBEATID = 2;
    //服务器心跳二级协议号
    private const Int32 MSGCLASSID_S_HEARTBEATRETURNID = 3;

    //最后一次收到服务器心跳包时间
    private DateTime lastServerBeatTime = DateTime.Now;
    //最后一次发送客户端心跳包时间
    private DateTime lastClientBeatTime = DateTime.Now.AddMinutes(-10);


    /// <summary>
    /// ctor
    /// </summary>
    /// <param name="socket"></param>
    /// <param name="beatCfg"></param>
    public HeartBeat(ISocketClient socket, HeartBeatCfg beatCfg, HeartBeatSource bt)
    {
        this.socket = socket;
        this.beatCfg = beatCfg;
        this.bt = bt;
    }

    /// <summary>
    /// 是否心跳消息
    /// </summary>
    /// <param name="mainID"></param>
    /// <param name="secondID"></param>
    /// <returns></returns>
    public static bool IsHeartBeatMsg(int mainID, int secondID)
    {
        return mainID == MSGCLASSID_B_COMMID && secondID == MSGCLASSID_S_HEARTBEATRETURNID;
    }

    //检查是否心跳超时
    public void CheckBeatTimeout()
    {
        //不在连接状态
        if (socket.MyConnectState != ConnectState.Connect)
        {
            return;
        }

        var eltime = ((TimeSpan)(DateTime.Now - lastServerBeatTime)).TotalSeconds;

        if (eltime >= this.beatCfg.ConnLoseTime)
        {
            this.lastServerBeatTime = DateTime.Now;         // 避免超时后，断线重连时反复检测心跳超时时间。导致一直超时一直关闭，然后一直重连的情况。
            socket.Disconnected("心跳超时");
        }
    }

    //发送心跳消息
    public void SendBeatMsg()
    {
        //不在连接状态
        if (socket.MyConnectState != ConnectState.Connect)
        {
            this.lastServerBeatTime = DateTime.Now; // 避免重连之后判定为超时
            return;
        }

        //检查是否达到心跳时间
        if (((TimeSpan)(DateTime.Now - lastClientBeatTime)).TotalSeconds >= this.beatCfg.BeatInterval)
        {
            //发送心跳消息
            this.socket.SendMessage(MSGCLASSID_B_COMMID, MSGCLASSID_S_HEARTBEATID, null);

            //设置客户端心跳时间
            this.lastClientBeatTime = DateTime.Now;

            //Debug.LogFormat("发送心跳包,{0}",bt);
        }
    }

    /// <summary>
    /// 收到心跳返回消息
    /// </summary>
    /// <param name="socket"></param>
    public void OnHeartBeartReturnMsg(ISocketClient socket)
    {
        if (this.socket == socket)
        {
            this.lastServerBeatTime = DateTime.Now;
        }
    }

    /// <summary>
    /// 外部重置上次发送时间 
    /// 1. Android 环境下，点击Home之后返回游戏
    /// </summary>
    public void ResetLastBeatTime()
    {
        this.lastServerBeatTime = DateTime.Now; // 避免重连之后判定为超时
    }
}

/// <summary>
/// 心跳管理器
/// </summary>
public class HeartBeatManager : Singleton<HeartBeatManager>
{
    //心跳
    private Dictionary<HeartBeatSource, HeartBeat> dctBeats = new Dictionary<HeartBeatSource, HeartBeat>();

    void Start()
    {
        this.InvokeRepeating("checkHeartBeat", 0.1f, 3.1f);
    }

    private void checkHeartBeat()
    {
        foreach (var bt in dctBeats)
        {
            //发送心跳消息
            bt.Value.SendBeatMsg();

            //检查心跳超时
            bt.Value.CheckBeatTimeout();
        }

    }

    /// <summary>
    /// 启动心跳
    /// </summary>
    /// <param name="beatType"></param>
    /// <param name="beatInterval"></param>
    /// <param name="connLoseTimeout"></param>
    public void StartHeartBeat(int beatType, Int32 beatInterval, Int32 connLoseTimeout)
    {
        var bs = (HeartBeatSource)beatType;

        if (bs == HeartBeatSource.HallServer)
        {
            Debug.Log("开启hallserver 心跳");
        }
        else
        {
            Debug.Log("开启gameserver 心跳");
        }

        if (dctBeats.ContainsKey(bs))
        {
            return;
        }

        var cfg = new HeartBeatCfg() { BeatInterval = beatInterval, ConnLoseTime = connLoseTimeout };

        if (bs == HeartBeatSource.HallServer)
        {
            dctBeats.Add(HeartBeatSource.HallServer, new HeartBeat(NetworkManager.Instance.HallSocket, cfg, HeartBeatSource.HallServer));
        }
        else
        {
            dctBeats.Add(HeartBeatSource.GameServer, new HeartBeat(NetworkManager.Instance.RoomSocket, cfg, HeartBeatSource.GameServer));
        }

    }

    /// <summary>
    /// 接收到心跳返回消息
    /// </summary>
    /// <param name="socket"></param>
    public void OnHeartBeatReturnMsg(ISocketClient socket)
    {
        foreach (var hb in dctBeats)
        {
            hb.Value.OnHeartBeartReturnMsg(socket);
        }
    }

    public void OnApplicationPause(bool pause)
    {
        if (!pause)
        {
            foreach (var bt in dctBeats)
                bt.Value.ResetLastBeatTime();       // 重新回到游戏之后重置超时时间，避免判定为心跳超时
        }
    }
}

