/// More is a class that has two properties, avatar and color, and a factory constructor that takes a
/// Map<String, dynamic> and returns a More instance.
class More {
  final bool avatar;
  final String color;

  const More({required this.avatar, required this.color});

  factory More.fromJson(Map<String, dynamic> json) {
    return More(
      avatar: json['avatar'] ?? false,
      color: json['color'],
    );
  }
}
