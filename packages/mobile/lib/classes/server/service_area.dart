/// It's a Dart class that represents a service area
class ServiceArea {
  final String name;
  final String description;
  final bool useGateway;
  final bool wip;
  final List<String>? components;
  final Map<String, AreaStoreItem>? store;

  const ServiceArea({
    required this.name,
    this.description = "",
    this.useGateway = false,
    this.wip = false,
    this.components,
    this.store,
  });

  factory ServiceArea.fromJson(Map<String, dynamic> json) {
    return ServiceArea(
      name: json['name'],
      description: json['description'] ?? "",
      wip: json['wip'] ?? false,
      components: json['components'] != null
          ? List<String>.from(json['components'])
          : null,
      useGateway: json['use_gateway'] ?? false,
      store: json['store'] != null
          ? Map<String, AreaStoreItem>.from(
              Map<String, dynamic>.from(json['store']).map(
                (key, value) => MapEntry(key, AreaStoreItem.fromJson(value)),
              ),
            )
          : null,
    );
  }
}

/// It's a Dart class that represents the JSON data that we get from the API
class AreaStoreItem {
  final String description;
  final bool isRequired;
  final String type;
  final List<String>? needFields;
  final List<String>? allowedComponents;
  final List<String>? values;
  final int priority;

  const AreaStoreItem({
    this.description = "",
    this.isRequired = false,
    this.type = "string",
    this.needFields,
    this.allowedComponents,
    this.values,
    this.priority = 0,
  });

  factory AreaStoreItem.fromJson(Map<String, dynamic> json) {
    return AreaStoreItem(
        description: json['description'] ?? "",
        isRequired: json['required'] ?? false,
        type: json['type'] ?? "string",
        needFields: json['need_fields'] != null
            ? List<String>.from(json['need_fields'])
            : null,
        allowedComponents: json['allowed_components'] != null
            ? List<String>.from(json['allowed_components'])
            : null,
        values:
            json['values'] != null ? List<String>.from(json['values']) : null,
        priority: json['priority'] ?? 0);
  }
}
