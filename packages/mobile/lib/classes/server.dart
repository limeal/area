import 'package:mobile/classes/server/authenticator.dart';
import 'package:mobile/classes/server/service.dart';

class Server {
  final List<Authenticator> authenticators;
  final List<Service> services;

  const Server({
    required this.authenticators,
    required this.services,
  });

  factory Server.fromJson(Map<String, dynamic> json) {
    return Server(
      authenticators: json['authenticators']
          .map<Authenticator>((json) => Authenticator.fromJson(json))
          .toList(),
      services: json['services']
          .map<Service>((json) => Service.fromJson(json))
          .toList(),
    );
  }
}
