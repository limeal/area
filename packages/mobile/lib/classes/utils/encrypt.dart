/// We're creating a map from the list of alphabets, where the key is the current character and the
/// value is the character at the index of the current character plus the shift
///
/// Args:
///   alphabets (List<String>): A list of alphabets.
///   shift (int): The number of characters to shift the alphabet by.
///
/// Returns:
///   A map with the key being the current character and the value being the shifted character.
Map<String, String> createMap(List<String> alphabets, int shift) {
  return alphabets.asMap().map((charIndex, currentChar) {
    int ind = (charIndex + shift) % alphabets.length;
    if (ind < 0) {
      ind += alphabets.length;
    }
    return MapEntry(currentChar, alphabets[ind]);
  });
}

/// It encrypts a string using the Caesar cipher.
///
/// Args:
///   org (String): The string to be encrypted.
///   shift (int): The amount of shift to be applied to the alphabets. Defaults to 0
///
/// Returns:
///   A map of the alphabets with the shift applied.
String encrypt(String org, [int shift = 0]) {
  final alphabets =
      'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890'
          .split('');
  final map = createMap(alphabets, shift);
  return org.toLowerCase().split('').map((char) => map[char] ?? char).join('');
}
