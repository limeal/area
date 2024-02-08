/// `Authorization` is a class that has a `type` property that is a `String` and is optional, a `name`
/// property that is a `String` and is required, and a `permanent` property that is a `bool` and is
/// required
class Authorization {
  final String type;
  final String name;
  final bool permanent;

  const Authorization(
      {this.type = "", required this.name, required this.permanent});

  factory Authorization.fromJson(Map<String, dynamic> json) {
    return Authorization(
        type: json['type'] ? json['type'] : "",
        name: json['name'],
        permanent: json['permanent']);
  }
}
