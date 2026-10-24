using System;

namespace FrameWork
{ 
    public class Random
    {
        private const int MBIG = 2147483647;

        private const int MSEED = 161803398;

        private const int MZ = 0;

        private int inext;

        private int inextp;

        private int[] SeedArray = new int[56];

        public Random() : this(Environment.TickCount)
        { }

        public Random(int Seed)
        {
            int num = 161803398 - Math.Abs(Seed);
            this.SeedArray[55] = num;
            int num2 = 1;
            for (int i = 1; i < 55; i++)
            {
                int num3 = 21 * i % 55;
                this.SeedArray[num3] = num2;
                num2 = num - num2;
                if (num2 < 0)
                {
                    num2 += 2147483647;
                }
                num = this.SeedArray[num3];
            }
            for (int j = 1; j < 5; j++)
            {
                for (int k = 1; k < 56; k++)
                {
                    this.SeedArray[k] -= this.SeedArray[1 + (k + 30) % 55];
                    if (this.SeedArray[k] < 0)
                    {
                        this.SeedArray[k] += 2147483647;
                    }
                }
            }
            this.inext = 0;
            this.inextp = 21;
            Seed = 1;
        }

        protected virtual double Sample()
        {
            return (double)this.InternalSample() * 4.6566128752457969E-10;
        }

        private int InternalSample()
        {
            int num = this.inext;
            int num2 = this.inextp;
            if (++num >= 56)
            {
                num = 1;
            }
            if (++num2 >= 56)
            {
                num2 = 1;
            }
            int num3 = this.SeedArray[num] - this.SeedArray[num2];
            if (num3 < 0)
            {
                num3 += 2147483647;
            }
            this.SeedArray[num] = num3;
            this.inext = num;
            this.inextp = num2;
            return num3;
        }

        //public virtual int Next()
        //{
        //    return this.InternalSample();
        //}

        private double GetSampleForLargeRange()
        {
            int num = this.InternalSample();
            bool flag = this.InternalSample() % 2 == 0;
            if (flag)
            {
                num = -num;
            }
            double num2 = (double)num;
            num2 += 2147483646.0;
            return num2 / 4294967293.0;
        }

        //public virtual int Next(int maxValue)
        //{
        //    if (maxValue < 0)
        //    {
        //        throw new ArgumentOutOfRangeException("maxValue", "随机数最大值不能小于0");
        //    }
        //    return (int)(this.Sample() * (double)maxValue);
        //}

        //public virtual double NextDouble()
        //{
        //    return this.Sample();
        //}

        //public virtual void NextBytes(byte[] buffer)
        //{
        //    if (buffer == null)
        //    {
        //        throw new ArgumentNullException("buffer");
        //    }
        //    for (int i = 0; i < buffer.Length; i++)
        //    {
        //        buffer[i] = (byte)(this.InternalSample() % 256);
        //    }
        //}

        public int RangeInt(int minValue, int maxValue)
        {
            if (minValue > maxValue)
            {
                throw new ArgumentOutOfRangeException("minValue", "范围随机数最小值不能大于最大值");
            }
            long num = (long)maxValue - (long)minValue;
            if (num <= 2147483647L)
            {
                return (int)(this.Sample() * (double)num) + minValue;
            }
            return (int)((long)(this.GetSampleForLargeRange() * (double)num) + (long)minValue);
        }

        public float RangeFloat(float minValue, float maxValue)
        {
            if (minValue > maxValue)
            {
                throw new ArgumentOutOfRangeException("minValue", "范围随机数最小值不能大于最大值");
            }
            float num = maxValue - minValue;
            if (num <= 2147483647L)
            {
                return (float)this.Sample() * num + minValue;
            }

            return (float)this.GetSampleForLargeRange() * num + minValue;
        }
    }
}