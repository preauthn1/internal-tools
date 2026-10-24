using UnityEngine;

public abstract class Singleton<T> : MonoBehaviour where T : MonoBehaviour
{
    public static T Instance { private set; get; }

    public virtual void Awake()
    {
        if (Instance != null)
        {
            Debug.LogError(string.Format("管理器{0}重复添加", typeof(T).Name));
            Object.Destroy(this);
            return;
        }

        Instance = gameObject.GetComponent<T>();
    }

    protected virtual void OnDestroy()
    {
        Instance = null;
    }
}
