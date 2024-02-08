/// It's a class that represents an Authenticator
import 'more.dart';

class Authenticator {
  final String name;
  final String authorizationUri;
  final bool enabled;
  final More more;

  const Authenticator(
      {required this.name,
      required this.authorizationUri,
      required this.enabled,
      required this.more});

  factory Authenticator.fromJson(Map<String, dynamic> json) {
    return Authenticator(
      name: json['name'],
      authorizationUri: json['authorization_uri'],
      enabled: json['enabled'] ?? false,
      more: More.fromJson(json['more']),
    );
  }
}
