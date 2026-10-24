using System;
using UnityEngine;

/// <summary>
/// 游戏配置信息
/// </summary>
[Serializable]
public class GameConfig : ScriptableObject
{  
    public string Version = string.Empty;           // 游戏资源版本号 
    public string Skin = string.Empty;           // 游戏资源样版
}