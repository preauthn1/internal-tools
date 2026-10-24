using UnityEngine;
using System.Collections;
using example.invalid.Sockets;
using System;
using System.Runtime.InteropServices;

public class Ipv6 : MonoBehaviour {

#if UNITY_IOS && !UNITY_EDITOR
	[DllImport("__Internal")]
	private static extern string getIPv6(string mHost, string mPort);  
#endif
	
	//"0.0.0.0&&ipv4"
  	public static string GetIPv6(string mHost, string mPort)
  	{
#if UNITY_IOS && !UNITY_EDITOR
		string mIPv6 = getIPv6(mHost, mPort);
		return mIPv6;
#else
   		return mHost + "&&ipv4";
#endif
  	}

	public static void getIPType(string serverIp, string serverPorts, out string newServerIp, out AddressFamily  mIPType)
    {
		newServerIp = serverIp;
    	mIPType = AddressFamily.InterNetwork;
		
      	try
      	{
        	string mIPv6 = GetIPv6(serverIp, serverPorts);
        	if (!string.IsNullOrEmpty(mIPv6))
        	{
          		string[] m_StrTemp = System.Text.RegularExpressions.Regex.Split(mIPv6, "&&");
          		if (m_StrTemp != null && m_StrTemp.Length >= 2)
          		{
            		string IPType = m_StrTemp[1];
            		if (IPType == "ipv6")
            		{
              			newServerIp = m_StrTemp[0];
              			mIPType = AddressFamily.InterNetworkV6;
            		}
          		}
        	}
      	}
      	catch (Exception e)
     	{
      		Debug.Log("GetIPv6 error:" + e);
      	}
    }

    public void SocketClient(string serverIp, string serverPorts)
    {
	      string newServerIp = "";
	      AddressFamily newAddressFamily = AddressFamily.InterNetwork;
	      getIPType(serverIp, serverPorts, out newServerIp, out newAddressFamily);
	      if (!string.IsNullOrEmpty(newServerIp)){ serverIp = newServerIp;}      
	      //Socket socketClient = new Socket(newAddressFamily, SocketType.Stream, ProtocolType.Tcp);
	      Debug.Log("Socket AddressFamily :" + newAddressFamily.ToString() + "ServerIp:" + serverIp);
     }
}
