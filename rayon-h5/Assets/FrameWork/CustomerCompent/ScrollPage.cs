using UnityEngine;
using System.Collections.Generic;
using UnityEngine.EventSystems;
using UnityEngine.UI;
using System;
using DG.Tweening;

public class ScrollPage : MonoBehaviour, IBeginDragHandler, IEndDragHandler
{
    ScrollRect rect;
    //页面：0，1，2，3  索引从0开始
    //每页占的比列：0/3=0  1/3=0.333  2/3=0.6666 3/3=1
    //float[] pages = { 0f, 0.333f, 0.6666f, 1f };
    List<float> pages = new List<float>();
    int currentPageIndex = 0;

    //滑动速度
    public float smooting = 4;

    //是否拖拽结束
    bool isDrag = false;

    /// <summary>
    /// 用于返回一个页码，-1说明page的数据为0
    /// </summary>
    public Action<int, int> OnPageChanged = (pageCount, currentIndex) => { };

    void Start()
    {
        rect = transform.GetComponent<ScrollRect>();
    }

    void Update()
    {
        if (!isDrag && pages.Count > 0)
        {
            rect.horizontalNormalizedPosition = Mathf.Lerp(rect.horizontalNormalizedPosition, pages[currentPageIndex], Time.deltaTime * smooting);
        }
    }

    public void OnBeginDrag(PointerEventData eventData)
    {
        isDrag = true;
    }

    public void OnEndDrag(PointerEventData eventData)
    {
        isDrag = false;

        if (pages.Count <= 1)
            return;

        if (currentPageIndex > pages.Count - 1)
            currentPageIndex = pages.Count - 1;

        if (Mathf.Abs(pages[currentPageIndex] - rect.horizontalNormalizedPosition) > 0.2f / pages.Count)
        {
            if (rect.horizontalNormalizedPosition > pages[currentPageIndex])
            {
                ToNextPage();
            }
            else if (rect.horizontalNormalizedPosition < pages[currentPageIndex])
            {
                ToPrePage();
            }
        }
        else
        {
            if (OnPageChanged != null)
                OnPageChanged(pages.Count, currentPageIndex);
        }
    }

    public void ToNextPage()
    {
        if (currentPageIndex == pages.Count - 1)
            return;

        currentPageIndex++;

        if (OnPageChanged != null)
            OnPageChanged(pages.Count, currentPageIndex);
    }

    public void ToPrePage()
    {
        if (currentPageIndex == 0)
            return;

        currentPageIndex--;

        if (OnPageChanged != null)
            OnPageChanged(pages.Count, currentPageIndex);
    }
    public void ToFirstPage()
    {
        currentPageIndex=0;

        if (OnPageChanged != null)
            OnPageChanged(pages.Count, currentPageIndex);
    }
    public void ToLastPage()
    {
        currentPageIndex = pages.Count - 1;

        if (OnPageChanged != null)
            OnPageChanged(pages.Count, currentPageIndex);
    }

    public void SetPageCount(int pageCount)
    {
        if (pageCount == 0)
            return;

        if (pages.Count != pageCount)
        {
            pages.Clear();

            for (int i = 0; i < pageCount; i++)
            {
                float page = 0;
                if (pageCount != 1)
                    page = i / ((float)(pageCount - 1));
                pages.Add(page);
            }

            OnEndDrag(null);
        }
    }
}
