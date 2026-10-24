using System;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;

[XLua.LuaCallCSharp]
public class TimeManager : Singleton<TimeManager>
{
    private LuaCoroutine m_LuaCoroutine;

    void Start()
    {
        m_LuaCoroutine = gameObject.AddComponent<LuaCoroutine>();
    }

    /// <summary>
    /// 
    /// </summary>
    /// <param name="countDown"></param>
    /// <param name="keyName"></param>
    /// <param name="notifyInterval">更新的最小频率为 0.1s 每次</param>
    public void StartCountDown(float countDown, string keyName, System.Action<float> updateCallBack, float notifyInterval = 1)
    {
        //Debug.Log("请尽量使用This的实例替换TimeManager调用该方法。");
        m_LuaCoroutine.StartCountDown(countDown, keyName, updateCallBack, notifyInterval);
    }

    public void RemoveCountDown(string keyName)
    {
        //Debug.Log("请尽量使用This的实例替换TimeManager调用该方法。");
        m_LuaCoroutine.RemoveCountDown(keyName);
    }

    public void InvokeRepeating(Action repeatAction, float time, float repeatRate)
    {
        //Debug.Log("请尽量使用This的实例替换TimeManager调用该方法。");
        m_LuaCoroutine.InvokeRepeating(repeatAction, time, repeatRate);
    }

    public void CancelInvokeRepeat(Action repeatAction)
    {
        //Debug.Log("请尽量使用This的实例替换TimeManager调用该方法。");
        m_LuaCoroutine.CancelInvokeRepeat(repeatAction);
    }

    public void StopAllRepeatInvoke()
    {
        //Debug.Log("请尽量使用This的实例替换TimeManager调用该方法。");
        m_LuaCoroutine.StopAllRepeatInvoke();
    }

    /// <summary>
    /// 延迟触发方法
    /// </summary>
    /// <param name="delayTime">延迟时间</param>
    /// <param name="action">触发的方法</param>
    public int DelayInvoke(float delayTime, System.Action action)
    {
        //Debug.Log("请尽量使用This的实例替换TimeManager调用该方法。");
        return m_LuaCoroutine.DelayInvoke(delayTime, action);
    }

    public void StopDelayInvoke(int delayID)
    {
        //Debug.Log("请尽量使用This的实例替换TimeManager调用该方法。");
        m_LuaCoroutine.StopDelayInvoke(delayID);
    }

    public void StopAllDelayInvoke()
    {
        //Debug.Log("请尽量使用This的实例替换TimeManager调用该方法。");
        m_LuaCoroutine.StopAllDelayInvoke();
    }

    public void YieldAndCallback(object to_yield, Action callback)
    {
        //Debug.Log("请尽量使用This的实例替换TimeManager调用该方法。");
        m_LuaCoroutine.YieldAndCallback(to_yield, callback);
    }

    public void YieldAndCallback(string sign, object to_yield, Action callback)
    {
        m_LuaCoroutine.YieldAndCallback(sign, to_yield, callback);
    }

    public new void StopAllCoroutines()
    {
        m_LuaCoroutine.StopAllCoroutines();
    }

    public void StopAllCoroutines(string sign)
    {
        m_LuaCoroutine.StopAllCoroutines(sign);
    }
    
    /// <summary>
    /// 用于www超时判断
    /// </summary>
    /// <param name="complateRuleFun">正常情况下执行完成的判定条件</param>
    /// <param name="timeOut">超时时间(秒)</param>
    /// <returns></returns>
    public IEnumerator TimeOutReturn(Func<bool> complateRuleFun, float timeOut)
    {
        float beginTime = Time.realtimeSinceStartup;

        if (complateRuleFun != null)
        {
            while (!complateRuleFun() && Time.realtimeSinceStartup - beginTime < timeOut)
            {
                yield return 1;
            }
        }
    }
}