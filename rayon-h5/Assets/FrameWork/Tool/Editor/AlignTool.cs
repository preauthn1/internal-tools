using System.Linq;
using System.Collections;
using UnityEngine;
using UnityEditor;

public class AlignTool : Editor
{
    [MenuItem("对齐工具/水平左对齐 ←")]
    public static void alignInHorziontalLeft()
    {
        float x = Mathf.Min(Selection.gameObjects.Select(obj => obj.transform.localPosition.x -
        ((RectTransform)obj.transform).sizeDelta.x / 2
        ).ToArray());

        foreach (GameObject gameObject in Selection.gameObjects)
        {
            gameObject.transform.localPosition = new Vector3(x + ((RectTransform)gameObject.transform).sizeDelta.x / 2,
                gameObject.transform.localPosition.y);
            // MathUtil.setPositionX(gameObject.transform.localPosition, x);
        }
    }

    [MenuItem("对齐工具/水平右对齐 →")]
    public static void alignInHorziontalRight()
    {
        float x = Mathf.Max(Selection.gameObjects.Select(obj => obj.transform.localPosition.x +
        ((RectTransform)obj.transform).sizeDelta.x / 2).ToArray());
        foreach (GameObject gameObject in Selection.gameObjects)
        {
            gameObject.transform.localPosition = new Vector3(x - ((RectTransform)gameObject.transform).sizeDelta.x / 2, gameObject.transform.localPosition.y);
            // MathUtil.setPositionX(gameObject.transform.localPosition, x);
        }
    }

    [MenuItem("对齐工具/垂直上对齐 ↑")]
    public static void alignInVerticalUp()
    {
        float y = Mathf.Max(Selection.gameObjects.Select(obj => obj.transform.localPosition.y +
        ((RectTransform)obj.transform).sizeDelta.y / 2).ToArray());

        foreach (GameObject gameObject in Selection.gameObjects)
        {
            gameObject.transform.localPosition = new Vector3(gameObject.transform.localPosition.x, y - ((RectTransform)gameObject.transform).sizeDelta.y / 2);
            //MathUtil.setPositionY(gameObject.transform.localPosition, y);
        }
    }

    [MenuItem("对齐工具/垂直下对齐 ↓")]
    public static void alignInVerticalDown()
    {
        float y = Mathf.Min(Selection.gameObjects.Select(obj => obj.transform.localPosition.y -
        ((RectTransform)obj.transform).sizeDelta.y / 2).ToArray());

        foreach (GameObject gameObject in Selection.gameObjects)
        {
            gameObject.transform.localPosition = new Vector3(gameObject.transform.localPosition.x, y + ((RectTransform)gameObject.transform).sizeDelta.y / 2);
            //MathUtil.setPositionY(gameObject.transform.localPosition, y);
        }
    }


    [MenuItem("对齐工具/水平均匀 |||")]
    public static void uniformDistributionInHorziontal()
    {
        int count = Selection.gameObjects.Length;
        float firstX = Mathf.Min(Selection.gameObjects.Select(obj => obj.transform.localPosition.x).ToArray());
        float lastX = Mathf.Max(Selection.gameObjects.Select(obj => obj.transform.localPosition.x).ToArray());
        float distance = (lastX - firstX) / (count - 1);
        var objects = Selection.gameObjects.ToList();
        objects.Sort((x, y) => (int)(x.transform.localPosition.x - y.transform.localPosition.x));
        for (int i = 0; i < count; i++)
        {
            objects[i].transform.localPosition = new Vector3(firstX + i * distance, objects[i].transform.localPosition.y);
            //MathUtil.setPositionX(position, firstX + i * distance);
        }
    }

    [MenuItem("对齐工具/垂直均匀 ☰")]
    public static void uniformDistributionInVertical()
    {
        int count = Selection.gameObjects.Length;
        float firstY = Mathf.Min(Selection.gameObjects.Select(obj => obj.transform.localPosition.y).ToArray());
        float lastY = Mathf.Max(Selection.gameObjects.Select(obj => obj.transform.localPosition.y).ToArray());
        float distance = (lastY - firstY) / (count - 1);
        var objects = Selection.gameObjects.ToList();
        objects.Sort((x, y) => (int)(x.transform.localPosition.y - y.transform.localPosition.y));
        for (int i = 0; i < count; i++)
        {
            objects[i].transform.localPosition = new Vector3(objects[i].transform.localPosition.x, firstY + i * distance);
            //MathUtil.setPositionX(position, firstX + i * distance);
        }
    }

    [MenuItem("对齐工具/重置大小(大) ■")]
    public static void ResizeMax()
    {
        var height = Mathf.Max(Selection.gameObjects.Select(obj => ((RectTransform)obj.transform).sizeDelta.y).ToArray());
        var width = Mathf.Max(Selection.gameObjects.Select(obj => ((RectTransform)obj.transform).sizeDelta.x).ToArray());
        foreach (GameObject gameObject in Selection.gameObjects)
        {
            ((RectTransform)gameObject.transform).sizeDelta = new Vector2(width, height);
        }
    }

    [MenuItem("对齐工具/重置大小(小) ●")]
    public static void ResizeMin()
    {
        var height = Mathf.Min(Selection.gameObjects.Select(obj => ((RectTransform)obj.transform).sizeDelta.y).ToArray());
        var width = Mathf.Min(Selection.gameObjects.Select(obj => ((RectTransform)obj.transform).sizeDelta.x).ToArray());
        foreach (GameObject gameObject in Selection.gameObjects)
        {
            ((RectTransform)gameObject.transform).sizeDelta = new Vector2(width, height);
        }
    }

}