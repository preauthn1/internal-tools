const paddingLeftZero = (digits, text) => {
  if (text.length>digits) {
      return text
  }

  const paddingDigits = digits - text.length

  var result = ""
  var i
  for (i = 0; i<paddingDigits; i++) {
      result += "0"
  }
  result = result + text

  return result
};

module.exports = {paddingLeftZero}