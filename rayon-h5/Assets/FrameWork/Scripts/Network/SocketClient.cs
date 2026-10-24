public interface ISocketClient
{
    ConnectState MyConnectState { get; }
    string Host { get; }
    int Port { get; }
    void OnRegister();
    void SendConnect(string host, int port);
    void Close();
    void SendMessage(int mainID, int secondID, byte[] buffer);
    void OnRemove();
    void Disconnected(string msg);
}

public enum ConnectState
{
    Connect = 0,        //连接服务器
    Connecting = 1,     //正在连接 
    Other = 2
}

public enum ConnectType
{
    Hall = 0,
    Room = 1
}