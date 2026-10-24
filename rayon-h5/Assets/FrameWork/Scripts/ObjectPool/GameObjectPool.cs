using UnityEngine;
using System;
using System.Collections;
using System.Collections.Generic;

namespace FrameWork
{

    [Serializable]
    public class PoolInfo
    {
        public string poolName;
        public GameObject prefab;
        public int poolSize;
        public bool fixedSize;
    }

    public class GameObjectPool
    {
        private int HistoryMaxSize = 0;        // 最高峰时候使用对象的数量

        private int maxSize;
        private int poolSize;
        private string poolName;
        private Transform poolRoot;
        private GameObject poolObjectPrefab;
        private Stack<GameObject> availableObjStack = new Stack<GameObject>();
        private List<GameObject> allocatedObjList = new List<GameObject>();

        public GameObjectPool(string poolName, GameObject poolObjectPrefab, int initCount, int maxSize, Transform pool)
        {
            this.poolName = poolName;
            this.poolSize = initCount;
            this.maxSize = maxSize;
            this.poolRoot = pool;
            this.poolObjectPrefab = poolObjectPrefab;

            //populate the pool
            for (int index = 0; index < initCount; index++)
            {
                AddObjectToPool(NewObjectInstance());
            }
        }

        private GameObject NewObjectInstance()
        {
            return GameObject.Instantiate(poolObjectPrefab) as GameObject;
        }

        public GameObject NextAvailableObject()
        {
            GameObject go = null;
            if (availableObjStack.Count == 0)
            {
#if UNITY_EDITOR
                Debug.LogWarning("缓存池缓存对象数量不足: " + poolName);
#endif
                AddObjectToPool(NewObjectInstance());
            }

            go = availableObjStack.Pop();


            if (go)
            {
                allocatedObjList.Add(go);
                go.SetActive(true);

                if (HistoryMaxSize < allocatedObjList.Count)
                    HistoryMaxSize = allocatedObjList.Count;

                return go;
            }
            else
                return NextAvailableObject();
        }

        //o(1)
        public void AddObjectToPool(GameObject go)
        {
            if (allocatedObjList.Contains(go))
            {
                allocatedObjList.Remove(go);
            }

            if (go)
            {
                go.RemoveAllScriptComponents();
                //add to pool
                go.SetActive(false);
                availableObjStack.Push(go);
                go.transform.SetParent(poolRoot, false);
            }
        }

        public void Delate()
        {
            while (availableObjStack.Count != 0)
            {
                UnityEngine.Object.Destroy(availableObjStack.Pop());
            }

            foreach (var allocatedObj in allocatedObjList)
            {
                UnityEngine.Object.Destroy(allocatedObj);
            }

            allocatedObjList.Clear();
            availableObjStack.Clear();
        }

        // 将所有已经分配出去的对象还原回对象池
        public void RetrieveAllocatedObj()
        {
            for (int i = 0; i < allocatedObjList.Count; i++)
            {
                AddObjectToPool(allocatedObjList[i]);
                i--;
            }
        }

        public int GetHistoryMaxSize()
        {
            return HistoryMaxSize;
        }
    }
}
