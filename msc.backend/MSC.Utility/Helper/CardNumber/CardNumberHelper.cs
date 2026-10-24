using System;
using System.Collections.Generic;
using System.Globalization;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Collections.Specialized;
using MSC.Utility.Extensions;
using example.invalid;
using example.invalid;
using Microsoft.AspNetCore.Http;
using MSC.Utility.Helper.HttpContext;
using MSC.Models.Models;

namespace MSC.Utility
{
    public class CardNumberHelper
    {
        private const int total_letter_types = 32;
        
        private const int total_card_number_length = 16;

        private static Dictionary<char, int> LetterDic = new Dictionary<char, int>
        {
            {'Z' , 1 },
            {'2' , 2 },
            {'3' , 3 },
            {'4' , 4 },
            {'5' , 5 },
            {'6' , 6 },
            {'7' , 7 },
            {'8' , 8 },
            {'9' , 9 },
            {'Y' , 10 },
            {'A' , 11 },
            {'B' , 12 },
            {'C' , 13 },
            {'D' , 14 },
            {'E' , 15 },
            {'F' , 16 },
            {'G' , 17 },
            {'H' , 18 },
            {'X' , 19 },
            {'J' , 20 },
            {'K' , 21 },
            {'L' , 22 },
            {'M' , 23 },
            {'N' , 24 },
            {'P' , 25 },
            {'Q' , 26 },
            {'R' , 27 },
            {'S' , 28 },
            {'T' , 29 },
            {'U' , 30 },
            {'V' , 31 },
            {'W' , 32 }
        };

        private static Dictionary<int, int> MultiplyValue = new Dictionary<int, int>
        {
            {1 ,15 },
            {2, 14 },
            {3, 13 },
            {4, 12 },
            {5, 11 },
            {6, 10 },
            {7, 9 },
            {8, 8 },
            {9, 7 },
            {10, 6 },
            {11, 5 },
            {12, 4 },
            {13, 3 },
            {14, 2 },
            {15, 1 },
        };



        /// <summary>
        /// 檢查卡號
        /// </summary>
        public static void CheckCardIsValid(string cardNumber)
        {
            if (string.IsNullOrWhiteSpace(cardNumber))
                throw new MSCException(ErrorCode.InvalidParameter_TextLength);

            if (cardNumber.Length != total_card_number_length)
                throw new MSCException(ErrorCode.InvalidParameter_TextLength);

            //非法的字母 (i、o、0、1)
            if (cardNumber.Any(x => !LetterDic.ContainsKey(x)))
                throw new MSCException(ErrorCode.InvalidParameter);

            char lastCheckCode = cardNumber.LastOrDefault();
            string previous_15_words = cardNumber.Truncate(total_card_number_length - 1);
            var expect_lastCode = GenerateLastCheckCode(previous_15_words);

            //前端檢查碼(最後一碼)
            if (lastCheckCode.ToString() != expect_lastCode)
                throw new MSCException(ErrorCode.NotFoundCardNo_FormatIncorrect);

            //後端檢查碼(最後兩個數字)
            var lastNumberResult = GenerateLast2Numbers(previous_15_words);

            if (lastNumberResult.OriginalLast2Word != lastNumberResult.Last2Word ||
               lastNumberResult.OriginalLastWord != lastNumberResult.LastWord)
                throw new MSCException(ErrorCode.NotFoundCardNo_FormatIncorrect);

        }


        public static string GenerateCardNumber()
        {
            //先產生隨機15碼
            string card_no = RandomHelper.GetRandCardNumber(total_card_number_length - 1);

            List<int> numbers = card_no.Where(x => x.IsNumeric()).Select(x => Int32.Parse(x.ToString())).ToList();

            //至少要三個數字
            if (numbers.Count() <= 3)
            {
                return GenerateCardNumber();
            }

            var lastWordResult = GenerateLast2Numbers(card_no);


            string fianlText = string.Empty;
            int current_number_position = 1;
            int total_number_count = lastWordResult.TotalNumbers;

            foreach (var charText in card_no)
            {
                if (!charText.IsNumeric())
                {
                    fianlText += charText.ToString();
                }
                else
                {
                    if (current_number_position == (total_number_count))
                    {
                        fianlText += lastWordResult.LastWord;
                    }
                    else if (current_number_position == (total_number_count - 1))
                    {
                        fianlText += lastWordResult.Last2Word;
                    }
                    else
                    {
                        fianlText += charText.ToString();
                    }
                    current_number_position++;
                }
            }

            string lastCheckCode = GenerateLastCheckCode(fianlText);
            string debug = $"{card_no} -> {fianlText} ({lastWordResult.Last2Word}{lastWordResult.LastWord})";
            return fianlText + lastCheckCode;
        }



        private static string GenerateLastCheckCode(string cardNumber)
        {
            if (cardNumber.Length != (total_card_number_length - 1))
                throw new MSCException(ErrorCode.InvalidParameter_TextLength);

            List<char> charList = cardNumber.Select(x => x).ToList();
            List<int> charMappingNumbers = new List<int>();

            //每個字母轉換成一組數字
            foreach (char text in charList)
            {
                if (!LetterDic.ContainsKey(text))
                {
                    throw new MSCException(ErrorCode.BusinessLogicDeny_UndefinedValue);
                }

                int mappingNumber = LetterDic[text];
                charMappingNumbers.Add(mappingNumber);
            }

            int indexSequence = 1;
            int totalSum = 0;

            //每個位置的字母依序乘上指定的數字
            foreach (int eachNumber in charMappingNumbers)
            {
                var multiplyValue = MultiplyValue[indexSequence];
                totalSum += multiplyValue * eachNumber;
                indexSequence++;
            }


            //除以34 取餘數
            int modValue = (totalSum % total_letter_types);
            modValue = modValue == 0 ? total_letter_types : modValue;


            //根據餘數轉乘對應的字母
            var finalText = LetterDic.Where(x => x.Value == modValue).FirstOrDefault().Key;
            return finalText.ToString();
        }

        private static CardNumberLastWordResultModel GenerateLast2Numbers(string card_no)
        {
            List<int> numbers = card_no.Where(x => x.IsNumeric()).Select(x => Int32.Parse(x.ToString())).ToList();

            List<char> letters = card_no.Where(x => !x.IsNumeric()).Select(x => x).ToList();

            List<int> lettersMappingNumber = letters.Select(x => LetterDic[x]).ToList();

            int totalSumOfLetterNumber = lettersMappingNumber.Sum();

            int firstNumber = numbers.FirstOrDefault();

            int first_divideValue = totalSumOfLetterNumber / firstNumber;
            int first_modValue = totalSumOfLetterNumber % firstNumber;

            int multiply_1 = (first_divideValue + first_modValue);
            int multiply_2 = numbers.Skip(1).Take(numbers.Count() - 3).Sum();

            int intA = multiply_1 * multiply_2;

            if (intA <= 10)
                intA = 97;

            int last2Number = intA >= 100 ? (intA / 10) % 10 : intA / 10;
            int lastNumber = intA % 10;

            if (last2Number == 0 || last2Number == 1)
                last2Number = 3;

            if (lastNumber == 0 || lastNumber == 1)
                lastNumber = 7;

            int totalNumbers = numbers.Count();
            var lastWordResult = new CardNumberLastWordResultModel
            {
                Numbers = numbers,
                TotalNumbers = totalNumbers,
                OriginalLastWord = numbers[totalNumbers - 1],
                OriginalLast2Word = numbers[totalNumbers - 2],
                LastWord = lastNumber,
                Last2Word = last2Number
            };

            return lastWordResult;
        }

    }


}
