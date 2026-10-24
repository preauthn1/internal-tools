using DG.Tweening;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using UnityEngine;

public static class DOTweenExtend
{
    /// <summary>
    ///  贝塞尔曲线移动效果
    /// </summary>
    /// <param name="transform"></param>
    /// <param name="target">目标位置</param>
    /// <param name="height">两点之间成线后，最高点的距离</param>
    /// <param name="duration">动画时间</param>
    /// <param name="direction">方向 0自动,1向左做曲线 2向右曲线</param>
    /// <returns></returns>
    public static Tweener DOMoveBezier(this Transform transform, Vector3 target, int height, float duration, byte direction = 0)
    {
        return transform.DOPath(GetCurvePath(transform, target, height, duration, direction), duration, DG.Tweening.PathType.CatmullRom);
    }

    static Vector3[] GetCurvePath(Transform transform, Vector3 target, float r, float duration, byte direction = 0)
    {
        Vector3 center = GetPos(transform, target, r * 2, direction);
        List<Vector3> points = GetBezierTwoPoints(transform.position, center, target);
        return points.ToArray();
    }

    static Vector3 GetMidPos(Vector3 start, Vector3 target)
    {
        return new Vector3(start.x + (target.x - start.x) / 2, start.y + (target.y - start.y) / 2);
    }
    static Vector3 GetPos(Transform transform, Vector3 _target, float r, byte direction = 0)
    {
        Vector3 _start = transform.position;
        switch (direction)
        {
            case 0:
                if (_start.x > _target.x)
                {
                    r = r * -1;
                }
                break;
            case 1:
                r = r * -1;
                break;
        }

        _start = transform.InverseTransformPoint(_start);
        _target = transform.InverseTransformPoint(_target);
        Vector3 midPos = GetMidPos(_start, _target);
        float x = midPos.x - (midPos.y - _start.y) / Mathf.Sqrt(Mathf.Pow(midPos.y - _start.y, 2) + Mathf.Pow(midPos.x - _start.x, 2)) * r;
        float y = midPos.y + (midPos.x - _start.x) / Mathf.Sqrt(Mathf.Pow(midPos.y - _start.y, 2) + Mathf.Pow(midPos.x - _start.x, 2)) * r;
        return transform.TransformPoint(new Vector3(x, y, _start.z));
    }

    static List<Vector3> GetBezierTwoPoints(Vector3 start, Vector3 center, Vector3 target)
    {
        List<Vector3> list = new List<Vector3>();
        for (int i = 0; i < 100; i++)//for循环计算100个点，  
        {
            Vector3 v1 = Vector3.Lerp(start, center, i / 100f);//注意：这里差值运算时间t必须相同  
            Vector3 v2 = Vector3.Lerp(center, target, i / 100f);

            var find = Vector3.Lerp(v1, v2, i / 100f);

            list.Add(find);
        }
        list.Add(target);
        return list;
    }
}



