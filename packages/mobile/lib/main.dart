import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/classes/utils/arguments.dart';
import 'package:mobile/pages/applet.dart';
import 'package:mobile/pages/applets.dart';
import 'package:mobile/pages/area.dart';
import 'package:mobile/pages/auth.dart';
import 'package:mobile/pages/authorize.dart';
import 'package:mobile/pages/create.dart';
import 'package:mobile/pages/create_next.dart';
import 'package:mobile/pages/home.dart';
import 'package:mobile/pages/profile.dart';
import 'package:mobile/store/store.dart';
import 'package:mobile/widgets/extensions/hex_color.dart';
import 'package:mobile/widgets/server_box.dart';
import 'package:uni_links/uni_links.dart';
import 'net/api.dart';
import 'dart:async';

final GlobalKey<NavigatorState> navigatorKey = GlobalKey<NavigatorState>();
void main() => runApp(const ProviderScope(child: MyApp()));

bool _initialUriIsHandled = false;

class MyApp extends StatefulWidget {
  const MyApp({Key? key}) : super(key: key);

  @override
  MyAppState createState() => MyAppState();
}

/// It listens for incoming links, and if it finds a code in the link, it will navigate to the
/// appropriate page
class MyAppState extends State<MyApp> {
  StreamSubscription? _sub;

  @override
  void initState() {
    super.initState();
    _handleIncomingLinks();
    _handleInitialUri();
  }

  @override
  void dispose() {
    _sub?.cancel();
    super.dispose();
  }

  /// > When the app is opened from a link, the link is parsed and the code is extracted. The code is
  /// then used to complete the authentication process
  ///
  /// Returns:
  ///   A StreamSubscription object.
  void _handleIncomingLinks() {
    if (!kIsWeb) {
      _sub = uriLinkStream.listen((Uri? uri) async {
        if (!mounted) return;
        final code = uri?.queryParameters["code"];
        if (code == null) return;
        final isRegistration = await Store.getRegistrationMode();
        final authenticatorName = await Store.getCurrentAuthenticator();
        if (authenticatorName == null) return;
        if (isRegistration) {
          navigatorKey.currentState?.pushReplacementNamed('/home',
              arguments: CodeArgument(
                  authenticatorName: authenticatorName,
                  code: code,
                  redirectURL: Store.redirectURL));
        } else {
          navigatorKey.currentState?.pushReplacementNamed('/authorize',
              arguments: CodeArgument(
                  authenticatorName: authenticatorName,
                  code: code,
                  redirectURL: Store.redirectURL));
        }
      }, onError: (Object err) {
        if (!mounted) return;
        logger.e('got err: $err');
      });
    }
  }

  /// _handleInitialUri() is called when the app is launched from a deep link. It checks if the app is
  /// launched from a deep link, and if so, it checks if the deep link contains a code parameter. If it
  /// does, it checks if the app is in registration mode or authorization mode, and then navigates to
  /// the appropriate page
  ///
  /// Returns:
  ///   The code is being returned.
  Future<void> _handleInitialUri() async {
    if (!_initialUriIsHandled) {
      _initialUriIsHandled = true;
      try {
        final uri = await getInitialUri();
        if (uri == null) return;
        if (!mounted) return;
        final code = uri.queryParameters["code"];
        if (code == null) return;
        final isRegistration = await Store.getRegistrationMode();
        final authenticatorName = await Store.getCurrentAuthenticator();
        if (authenticatorName == null) return;
        if (isRegistration) {
          navigatorKey.currentState?.pushReplacementNamed('/home',
              arguments: CodeArgument(
                  authenticatorName: authenticatorName,
                  code: code,
                  redirectURL: Store.redirectURL));
        } else {
          navigatorKey.currentState?.pushReplacementNamed('/authorize',
              arguments: CodeArgument(
                  authenticatorName: authenticatorName,
                  code: code,
                  redirectURL: Store.redirectURL));
        }
        // Go to home page
      } on PlatformException catch (err) {
        if (!mounted) return;
      } on FormatException catch (err) {
        if (!mounted) return;
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
        theme: ThemeData(
          primarySwatch: Colors.blueGrey,
          inputDecorationTheme: const InputDecorationTheme(
            filled: true,
            fillColor: Colors.white,
            border: OutlineInputBorder(),
          ),
          textTheme: Theme.of(context).textTheme.apply(
                bodyColor: HexColor.fromHex('#222222'),
                displayColor: HexColor.fromHex('#222222'),
              ),
        ),

        /// Defining the routes for the app.
        navigatorKey: navigatorKey,
        routes: <String, WidgetBuilder>{
          '/': (BuildContext context) => ServerBox(),
          '/auth': (BuildContext context) => const AuthPage(),
          '/authorize': (BuildContext context) => const AuthorizePage(),
          '/home': (BuildContext context) => HomePage(),
          '/profile': (BuildContext context) => const ProfilePage(),
          '/applet': (BuildContext context) => const AppletPage(),
          '/applets': (BuildContext context) => AppletsPage(),
          '/create': (BuildContext context) => const CreatePage(),
          '/create/next': (BuildContext context) => const CreateNextPage(),
          '/create/area': (BuildContext context) => const AreaPage(),
        });
  }
}
