/// It's a class that represents an area of the app
import 'service.dart';

class Area {
  final String id;
  final String type; // action or reaction
  final String service;
  Service? serviceObj;
  final String name;
  final Map<String, dynamic>? store;

  Area({
    required this.id,
    required this.type,
    required this.service,
    this.serviceObj,
    required this.name,
    this.store,
  });

  factory Area.fromJson(Map<String, dynamic> json) {
    return Area(
      id: json['id'],
      type: json['type'],
      service: json['service'],
      name: json['name'],
      store: json['store'] != null
          ? Map<String, dynamic>.from(json['store'])
          : null,
    );
  }
}
