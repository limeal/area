/// It's a Dart class that represents the JSON data returned by the API
class ApiData {
  final List<dynamic> elements;
  final List<String> fields;

  const ApiData({
    required this.elements,
    required this.fields,
  });

  factory ApiData.fromJson(Map<String, dynamic> json) {
    return ApiData(
        elements: json['data'],
        fields:
            json['fields']?.map<String>((elem) => elem.toString()).toList());
  }
}
