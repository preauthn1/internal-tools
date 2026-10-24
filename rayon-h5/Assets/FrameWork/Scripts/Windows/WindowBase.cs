using System;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class WindowBase : LuaCoroutine
{
    public string ProjectName
    {
        get;
        private set;
    }

    public virtual void BindingLuaScript(string projectName)
    {
        ProjectName = projectName;
    }

    /// <summary>
    /// 设置窗体节点
    /// </summary>
    /// <param name="windowNode"></param>
    internal void SetWindowNode(WindowNode windowNode)
    {
        WindowNode = windowNode;
    }

    public WindowNode WindowNode { get; private set; }

    /// <summary>
    /// 刷新界面数据
    /// </summary>
    /// <param name="windowData"></param>
    public virtual void RefreshWindowData(object windowData)
    {
    }

    /// <summary>
    /// 处理各个界面在点击Esc时候的处理，由WindowManager 触发
    /// </summary>
    public virtual void EscapeWindow()
    {

    }
}


/// <summary>
/// 处理UI延时调用
/// </summary>
public class LuaCoroutine : MonoBehaviour
{
    #region 倒计时

    class CountDownInfo
    {
        /// <summary>
        /// 关键字
        /// </summary>
        public string KeyName { get; set; }

        /// <summary>
        /// 倒计时
        /// </summary>
        public float CountDown { get; set; }

        /// <summary>
        /// 下一次更新时间
        /// </summary>
        public float NextCountDown { get; set; }

        /// <summary>
        /// 更新间隔
        /// </summary>
        public float NotifyInterval { get; set; }

        /// <summary>
        /// 更新通知回掉
        /// </summary>
        public System.Action<float> UpdateNotify { get; set; }

        /// <summary>
        /// 是否已经过时的
        /// </summary>
        public bool IsObsolete { get; internal set; }
    }

    List<CountDownInfo> m_CountDownInfoList = new List<CountDownInfo>();

    /// <summary>
    /// 
    /// </summary>
    /// <param name="countDown"></param>
    /// <param name="keyName"></param>
    /// <param name="notifyInterval">更新的最小频率为 0.1s 每次</param>
    public void StartCountDown(float countDown, string keyName, System.Action<float> updateCallBack, float notifyInterval = 1)
    {
        if (notifyInterval < 0.1f)
            notifyInterval = 0.1f;
        CountDownInfo item = new CountDownInfo();
        item.CountDown = item.NextCountDown = countDown;
        item.NotifyInterval = notifyInterval;
        item.KeyName = keyName;
        item.UpdateNotify = updateCallBack;
        m_CountDownInfoList.Add(item);
    }

    public void RemoveCountDown(string keyName)
    {
        CountDownInfo findItem = m_CountDownInfoList.Find((temp) => temp.KeyName == keyName && !temp.IsObsolete);
        if (findItem != null)
        {
            findItem.IsObsolete = true;
            findItem.UpdateNotify = null;
        }
    }

    protected void Update()
    {
        if (m_CountDownInfoList.Count > 0)
        {
            float deltaTime = Time.deltaTime;
            for (int i = m_CountDownInfoList.Count - 1; i >= 0; i--)
            {
                CountDownInfo item = m_CountDownInfoList[i];
                item.CountDown -= deltaTime;
                if (item.CountDown <= item.NextCountDown)
                {
                    if (item.UpdateNotify != null && !item.IsObsolete)
                    {
                        item.UpdateNotify(item.NextCountDown);
                    }
                    if (item.CountDown <= 0 || item.IsObsolete)
                    {
                        item.UpdateNotify = null;
                        m_CountDownInfoList.RemoveAt(i);
                    }
                    else
                    {
                        item.NextCountDown -= item.NotifyInterval;
                        if (item.NextCountDown < 0)
                        {
                            item.NextCountDown = 0;
                        }
                    }
                }
            }
        }
    }

    #endregion

    #region InvokeRepeat

    private Dictionary<Action, IEnumerator> RepeatEnumeratorDic = new Dictionary<Action, IEnumerator>();

    public void InvokeRepeating(Action repeatAction, float time, float repeatRate)
    {
        if (RepeatEnumeratorDic.ContainsKey(repeatAction))
        {
            Debug.Log("不能Repeating同一个方法");
            return;
        }

        IEnumerator enumerator = RepeatInvoke(repeatAction, time, repeatRate);
        RepeatEnumeratorDic.Add(repeatAction, enumerator);

        StartCoroutine(enumerator);
    }

    public void CancelInvokeRepeat(Action repeatAction)
    {
        if (!RepeatEnumeratorDic.ContainsKey(repeatAction))
            return;

        StopCoroutine(RepeatEnumeratorDic[repeatAction]);
        RepeatEnumeratorDic.Remove(repeatAction);
    }

    public void StopAllRepeatInvoke()
    {
        foreach (var item in RepeatEnumeratorDic)
            StopCoroutine(item.Value);

        RepeatEnumeratorDic.Clear();
    }

    private IEnumerator RepeatInvoke(Action repeatAction, float time, float repeatRate)
    {
        yield return new WaitForSeconds(time);

        while (repeatAction != null)
        {
            repeatAction();
            yield return new WaitForSeconds(repeatRate);
        }
    }

    #endregion

    #region DelayInvoke

    private int mCoroutineID = 0;
    private int GetCoroutineID()
    {
        mCoroutineID++;
        if (mCoroutineID == 100000)
            mCoroutineID = 1;
        return mCoroutineID;
    }

    Dictionary<int, Coroutine> coroutineDict = new Dictionary<int, Coroutine>();

    public int DelayInvoke(float delayTime, Action callBack)
    {
        if (delayTime > 0)
        {
            int delayID = GetCoroutineID();
            Coroutine cor = StartCoroutine(DelayCoroutine(delayID, delayTime, callBack));
            coroutineDict.Add(delayID, cor);
            return delayID;
        }
        else
        {
            if (callBack != null)
                callBack();

            return 0;
        }
    }

    public void StopDelayInvoke(int delayID)
    {
        StopCoroutine(coroutineDict[delayID]);
    }

    public void StopAllDelayInvoke()
    {
        foreach (var item in coroutineDict)
        {
            if (item.Value != null)
            {
                StopCoroutine(item.Value);
            }
        }
        coroutineDict.Clear();
    }

    protected IEnumerator DelayCoroutine(int delayID, float delayTime, Action callBack)
    {
        if (delayTime > 0)
            yield return new WaitForSeconds(delayTime);

        if (coroutineDict.ContainsKey(delayID))
        {
            if (callBack != null)
                callBack();

            coroutineDict[delayID] = null;
            coroutineDict.Remove(delayID);
        }
    }

    #endregion

    #region 协程

    public void YieldAndCallback(object to_yield, Action callback)
    {
        StartCoroutine(CoBody(to_yield, callback));
    }

    private IEnumerator CoBody(object to_yield, Action callback)
    {
        if (to_yield is IEnumerator)
            yield return StartCoroutine((IEnumerator)to_yield);
        else
            yield return to_yield;
        callback();
    }

    Dictionary<string, List<Coroutine>> signCoroutineDict = new Dictionary<string, List<Coroutine>>();

    public void YieldAndCallback(string sign, object to_yield, Action callback)
    {
        Coroutine cor = StartCoroutine(CoBody(to_yield, callback));

        if (!signCoroutineDict.ContainsKey(sign))
            signCoroutineDict.Add(sign, new List<Coroutine>());

        signCoroutineDict[sign].Add(cor);
    }

    public void StopAllCoroutines(string sign)
    {
        foreach (var item in signCoroutineDict)
        {
            if (item.Key == sign && item.Value != null)
            {
                foreach (var cor in signCoroutineDict[sign])
                    StopCoroutine(cor);

                Debug.Log("成功停掉" + sign + "的协程");
                signCoroutineDict[sign].Clear();
            }
        }
    }

    public void BreakCoroutine(object to_yield, Action callback)
    {
        Coroutine cor = StartCoroutine(CoBody(to_yield, callback));
        StopCoroutine(cor);
    }


    #endregion
}