import 'package:mobile/classes/server/authenticator.dart';
import 'package:shared_preferences/shared_preferences.dart';

/// It's a class that stores data in the shared preferences
class Store {
  static final Future<SharedPreferences> _prefs =
      SharedPreferences.getInstance();

  /// > Get the value of the key 'isRegistration' from the SharedPreferences object. If the key doesn't
  /// exist, return false
  ///
  /// Returns:
  ///   A Future<bool>
  static Future<bool> getRegistrationMode() async {
    return _prefs.then((SharedPreferences prefs) {
      return prefs.getBool('isRegistration') ?? false;
    });
  }

  /// "Get the current authenticator from the shared preferences."
  ///
  /// The function is asynchronous, so it returns a Future. The function is also marked as async, which
  /// means that it can use the await keyword
  ///
  /// Returns:
  ///   A Future<String?>
  static Future<String?> getCurrentAuthenticator() async {
    return _prefs.then((SharedPreferences prefs) {
      return prefs.getString('currentAuthenticator');
    });
  }

  /// It sets the value of the key 'isRegistration' to the value passed in the function.
  ///
  /// Args:
  ///   value (bool): The value to store.
  static setRegistrationMode(bool value) async {
    final SharedPreferences prefs = await _prefs;
    prefs.setBool('isRegistration', value);
  }

  /// It sets the current authenticator to the value passed in.
  ///
  /// Args:
  ///   value (String): The value to be stored.
  static setCurrentAuthenticator(String value) async {
    final SharedPreferences prefs = await _prefs;
    prefs.setString('currentAuthenticator', value);
  }

  /// "Get the server URL from the shared preferences, or return the default server URL if it's not
  /// set."
  ///
  /// The `_prefs` variable is a `Future<SharedPreferences>` that is initialized in the `initState`
  /// function
  ///
  /// Returns:
  ///   A Future<String>
  static Future<String> getServerURL() async {
    return _prefs.then((SharedPreferences prefs) {
      return prefs.getString('serverURL') ?? defaultServerURL;
    });
  }

  /// "Get the port from the shared preferences, or return the default port if it's not set."
  ///
  /// The `getPort()` function returns a `Future<String>`. This means that it returns a `Future` that
  /// will eventually contain a `String`
  ///
  /// Returns:
  ///   A Future<String>
  static Future<String> getPort() async {
    return _prefs.then((SharedPreferences prefs) {
      return prefs.getString('port') ?? defaultPort;
    });
  }

  /// It sets the serverURL in the shared preferences.
  ///
  /// Args:
  ///   value (String): The value to be stored.
  static setServerURL(String value) async {
    final SharedPreferences prefs = await _prefs;
    prefs.setString('serverURL', value);
  }

  /// It sets the port value in the shared preferences.
  ///
  /// Args:
  ///   value (String): The value to store.
  static setPort(String value) async {
    final SharedPreferences prefs = await _prefs;
    prefs.setString('port', value);
  }

  /// "Get the token from the shared preferences, and return it as a future."
  ///
  /// The `_prefs` variable is a `Future<SharedPreferences>` that is initialized in the `initState`
  /// function
  ///
  /// Returns:
  ///   A Future<String?>
  static Future<String?> getToken() async {
    return _prefs.then((SharedPreferences prefs) {
      return prefs.getString('token');
    });
  }

  /// It sets the token value in the shared preferences.
  ///
  /// Args:
  ///   value (String): The value to store
  static setToken(String value) async {
    final SharedPreferences prefs = await _prefs;
    prefs.setString('token', value);
  }

  /// > Get the value of the 'open' key from the SharedPreferences object. If the key doesn't exist,
  /// return false
  ///
  /// Returns:
  ///   A Future<bool>
  static Future<bool> getOpen() async {
    return _prefs.then((SharedPreferences prefs) {
      return prefs.getBool('open') ?? false;
    });
  }

  /// It sets the value of the key 'open' to the value of the variable 'value'
  ///
  /// Args:
  ///   value (bool): The value to store.
  static setOpen(bool value) async {
    final SharedPreferences prefs = await _prefs;
    prefs.setBool('open', value);
  }

  /// It removes all the data from the shared preferences.
  static void clear() async {
    final SharedPreferences prefs = await _prefs;
    prefs.remove('token');
    prefs.remove('serverURL');
    prefs.remove('port');
    prefs.remove('open');
  }

  static String defaultServerURL = '10.0.2.2';
  static String defaultPort = '8080';
  static String serverURL = defaultServerURL;
  static String port = defaultPort;
  static String redirectURL = 'http://localhost:8081';
}
