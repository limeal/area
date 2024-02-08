import 'package:mobile/classes/server/service.dart';
import 'package:mobile/classes/server/service_area.dart';

/// It's a class that holds the code, the authenticator name, and the redirect URL
class CodeArgument {
  final String authenticatorName;
  final String code;
  final String redirectURL;

  const CodeArgument({
    required this.authenticatorName,
    required this.code,
    required this.redirectURL,
  });
}

/// `AreaArgument` is a class that contains a `type` property, a `service` property, and an `item`
/// property
class AreaArgument {
  final String type;
  final Service? service;
  final ServiceArea? item;

  const AreaArgument({
    required this.type,
    required this.service,
    required this.item,
  });
}

/// AppletArgument is a class that represents an argument to an applet.
class AppletArgument {
  final String appletId;

  const AppletArgument({
    required this.appletId,
  });
}
