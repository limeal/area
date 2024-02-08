import 'package:flutter/material.dart';
import 'package:mobile/net/api.dart';
import 'package:mobile/store/store.dart';
import 'package:mobile/widgets/extensions/hex_color.dart';

import 'more.dart';
import 'service_area.dart';
import 'authenticator.dart';

/// It's a Dart class that represents a service
class Service {
  final String name;
  final String? description;
  final Authenticator? authenticator;
  final More? more;
  final List<ServiceArea>? actions;
  final List<ServiceArea>? reactions;

  const Service({
    required this.name,
    this.description,
    this.authenticator,
    this.more,
    this.actions,
    this.reactions,
  });

  factory Service.fromJson(Map<String, dynamic> json) {
    return Service(
      name: json['name'],
      description: json['description'],
      authenticator: json['authenticator'] != null
          ? Authenticator.fromJson(json['authenticator'])
          : null,
      more: json['more'] != null ? More.fromJson(json['more']) : null,
      actions: json['actions']?.map<ServiceArea>((json) {
        return ServiceArea.fromJson(json);
      }).toList(),
      reactions: json['reactions']
          ?.map<ServiceArea>((json) => ServiceArea.fromJson(json))
          .toList(),
    );
  }

  Color getColor() {
    if (more != null) {
      return HexColor.fromHex(more!.color);
    }

    if (authenticator == null) {
      return HexColor.fromHex('#222222');
    }

    return HexColor.fromHex(authenticator!.more.color);
  }

  String getAvatar() {
    String baseUrl = api.getBaseUrl();
    if (more != null) {
      return '$baseUrl/assets/$name.png';
    }

    if (authenticator == null) {
      return 'https://via.placeholder.com/50';
    }

    return '$baseUrl/assets/${authenticator!.name}.png';
  }
}
