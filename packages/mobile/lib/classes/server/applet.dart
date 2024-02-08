/// It's a Dart class that represents an Applet
class Applet {
  final String id;
  final String name;
  final String description;
  final bool isPublic;
  final String action;
  final bool active;
  final String status;

  const Applet({
    required this.id,
    required this.name,
    required this.description,
    required this.isPublic,
    required this.action,
    required this.active,
    required this.status,
  });

  factory Applet.fromJson(Map<String, dynamic> json) {
    return Applet(
      id: json['id'],
      name: json['name'],
      description: json['description'],
      isPublic: json['public'],
      action: json['action'],
      active: json['active'],
      status: json['status'],
    );
  }
}
