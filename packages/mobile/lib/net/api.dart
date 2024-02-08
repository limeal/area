import 'dart:convert';
import 'dart:core';
import 'package:analyzer_plugin/utilities/pair.dart';
import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';
import 'package:fluttertoast/fluttertoast.dart';
import 'package:http/http.dart';
import 'package:logger/logger.dart';
import 'package:http/http.dart' as http;

import 'package:mobile/classes/server.dart';
import 'package:mobile/classes/server/api_data.dart';
import 'package:mobile/classes/server/applet.dart';
import 'package:mobile/classes/server/area.dart';
import 'package:mobile/store/store.dart';
import 'package:mobile/classes/server/authorization.dart';
import 'package:mobile/classes/server/account.dart';
import 'package:mobile/classes/utils/encrypt.dart';
import 'package:mobile/classes/utils/error_toast.dart';

var logger = Logger();
final api = ApiClient();

class ApiClient {
  /// It returns a string that is the base URL of the server
  ///
  /// Returns:
  ///   A string.
  String getBaseUrl() {
    return 'http://${Store.serverURL}:${Store.port}';
  }

  /// It gets the server URL and port from the store, and then returns a URL
  ///
  /// Args:
  ///   endpoint (String): The endpoint of the API you want to call.
  ///   queryParams (Map<String, dynamic>): This is a map of the query parameters that you want to pass
  /// to the server.
  ///
  /// Returns:
  ///   A string that is the url to the endpoint.
  Future<String> _getUrl(
      String endpoint, Map<String, dynamic> queryParams) async {
    final serverURL = await Store.getServerURL();
    final serverPort = await Store.getPort();

    if (Store.serverURL != serverURL) {
      Store.serverURL = serverURL;
    }

    if (Store.port != serverPort) {
      Store.port = serverPort;
    }

    return Uri.http('$serverURL:$serverPort', endpoint, queryParams).toString();
  }

  /// It gets the server from the server.
  ///
  /// Returns:
  ///   The response body is being returned.
  Future<String> getServer() async {
    final response = await http.get(Uri.parse(await _getUrl('/', {})),
        headers: {
          "Access-Control-Allow-Origin": "*",
          "Content-Type": "application/json"
        });
    if (response.statusCode != 200) {
      showError.createToast('Cannot get server');
      throw Exception('Cannot get server');
    }
    return response.body;
  }

  /// It gets the about.json file from the server and returns a Server object
  ///
  /// Returns:
  ///   A Future<Server?>
  Future<Server?> getAbout() async {
    final response = await http.get(Uri.parse(await _getUrl('/about.json', {})),
        headers: {
          "Access-Control-Allow-Origin": "*",
          "Content-Type": "application/json"
        });
    if (response.statusCode != 200) {
      showError.createToast(
          "Something went wrong.\nCannot get About.json.\nCheck your connection and reload the app.");
      return null;
    }
    final body = jsonDecode(response.body);
    try {
      final server = Server.fromJson(body['server']);
      return server;
    } catch (e) {
      showError.createToast(e);
      logger.e(e);
      return null;
    }
  }

  /// It takes in a mode, email, and password, encrypts the password, and then sends a POST request to
  /// the server with the email and encrypted password
  ///
  /// Args:
  ///   mode: login or register
  ///   email (String): The email of the user
  ///   password (String): The password of the user
  Future<void> postAccount(mode, String email, String password) async {
    final encodedPassword = encrypt(password + email, email.length * 3);

    final response = await http.post(
        Uri.parse(await _getUrl('/auth/$mode', {})),
        headers: {
          "Access-Control-Allow-Origin": "*",
          "Content-Type": "application/json"
        },
        body: jsonEncode({
          "email": email,
          "encoded_password": encodedPassword,
        }));
    final body = jsonDecode(response.body);
    if (!(response.statusCode == 200 || response.statusCode == 201)) {
      showError.createToast('Cannot authenticate the user');
      throw Exception('Cannot authenticate the user');
    }
    if (body['data'] == null) {
      showError.createToast('Wrong body in authAccount');
      throw Exception('Wrong body in authAccount');
    }
    Store.setToken(body['data']['token']);
  }

  /// It sends a POST request to the server with the authenticator, code, and redirectUri
  ///
  /// Args:
  ///   authenticator (String): The name of the authenticator you want to use.
  ///   code (String): The code returned by the OAuth provider
  ///   redirectUri (String): The redirect URI that you set in the OAuth provider.
  Future<void> postExternalOAuth(
    String authenticator,
    String code,
    String redirectUri,
  ) async {
    final response = await http.post(
        Uri.parse(await _getUrl('/auth/external', {})),
        headers: {
          "Access-Control-Allow-Origin": "*",
          "Content-Type": "application/json"
        },
        body: jsonEncode({
          "authenticator": authenticator,
          "code": code,
          "redirect_uri": redirectUri
        }));
    final returnBody = jsonDecode(response.body);
    if (!(response.statusCode == 201 || response.statusCode == 200)) {
      throw Exception(returnBody['error'] ?? "Wrong status code");
    }
    Store.setToken(returnBody['data']['token']);
  }

  /// It gets the token from the store, if it's null it throws an exception, if it's not null it makes a
  /// get request to the url with the token as a header, if the status code is not 200 it throws an
  /// exception, if it is 200 it returns the authorizations
  ///
  /// Returns:
  ///   A list of authorizations
  Future<List<Authorization>> getAuthorizations() async {
    final token = await Store.getToken();
    if (token == null) {
      throw Exception("Token must not be null !");
    }

    final response = await http.get(
      Uri.parse(await _getUrl('/authorization', {})),
      headers: {
        "Access-Control-Allow-Origin": "*",
        "Authorization": "Bearer $token"
      },
    );
    final body = jsonDecode(response.body);
    if (response.statusCode != 200) {
      final code = response.statusCode;
      showError.createToast("Wrong status code : $code");
      throw Exception(body['error'] ?? "Wrong status code");
    }
    var authorizations = body['data']['authorizations'];
    return authorizations
        .map<Authorization>((json) => Authorization.fromJson(json))
        .toList();
  }

  /// It creates an authorization for the user
  ///
  /// Args:
  ///   authenticator (String): The name of the authenticator you want to use.
  ///   code (String): The code you received from the authenticator.
  ///   redirectUri (String): The redirect URI that you specified when you created the authorization.
  ///
  /// Returns:
  ///   A string
  Future<String> createAuthorization(
    String authenticator,
    String code,
    String redirectUri,
  ) async {
    final token = await Store.getToken();
    if (token == null) {
      showError.createToast("Token must not be null !");
      throw Exception("Token must not be null !");
    }

    final response = await http.post(
        Uri.parse(await _getUrl('/authorization', {})),
        headers: {
          "Access-Control-Allow-Origin": "*",
          "Content-Type": "application/json",
          "Authorization": "Bearer $token"
        },
        body: jsonEncode({
          "authenticator": authenticator,
          "code": code,
          "redirect_uri": redirectUri
        }));
    final returnBody = jsonDecode(response.body);
    if (response.statusCode != 201) {
      final code = response.statusCode;
      showError.createToast("Wrong status code : $code");
      throw Exception(returnBody['error'] ?? "Wrong status code");
    }
    return returnBody['data']['message'];
  }

  /// It deletes the authorization of the user.
  ///
  /// Returns:
  ///   A string
  Future<String> deleteAuthorization() async {
    final token = await Store.getToken();
    if (token == null) {
      showError.createToast("Token must not be null !");
      throw Exception("Token must not be null !");
    }

    final response = await http
        .delete(Uri.parse(await _getUrl('/auth/external', {})), headers: {
      "Access-Control-Allow-Origin": "*",
      "Content-Type": "application/json",
      "Authorization": "Bearer $token"
    });
    final returnBody = jsonDecode(response.body);
    if (response.statusCode != 200) {
      final code = response.statusCode;
      showError.createToast("Wrong status code : $code");
      throw Exception(returnBody['error'] ?? "Wrong status code");
    }
    return returnBody['data']['message'];
  }

  /// It gets the services authorized by the user
  ///
  /// Returns:
  ///   A map of service_name => bool
  Future<Map<String, bool>> getServicesAuthorized() async {
    final token = await Store.getToken();
    if (token == null) {
      showError.createToast("Token must not be null !");
      throw Exception("Token must not be null !");
    }

    final response = await http.get(
      Uri.parse(await _getUrl('/authorization/services', {})),
      headers: {
        "Access-Control-Allow-Origin": "*",
        "Authorization": "Bearer $token"
      },
    );
    final body = jsonDecode(response.body);
    if (response.statusCode != 200) {
      final code = response.statusCode;
      showError.createToast("Wrong status code : $code");
      throw Exception(body['error'] ?? "Wrong status code");
    }
    var services = body['data']['services'];
    // services is a map of service_name => bool
    return Map<String, bool>.from(
      Map<String, dynamic>.from(services).map(
        (key, value) => MapEntry(key, value as bool),
      ),
    );
  }

  /// > This function will get the account of the user from the API
  ///
  /// Returns:
  ///   Account.fromJson(body['data']['account']);
  Future<Account> getAccount() async {
    final token = await Store.getToken();
    if (token == null) {
      showError.createToast("Token must not be null !");
      throw Exception("Token must not be null !");
    }

    final response = await http.get(
      Uri.parse(await _getUrl('/me', {})),
      headers: {
        "Access-Control-Allow-Origin": "*",
        "Authorization": "Bearer $token"
      },
    );
    final body = jsonDecode(response.body);
    if (response.statusCode != 200) {
      final code = response.statusCode;
      showError.createToast("Wrong status code : $code");
      throw Exception(body['error'] ?? "Wrong status code");
    }
    return Account.fromJson(body['data']['account']);
  }

  /// It modifies the username of the user
  ///
  /// Args:
  ///   newUsername (String): The new username you want to change to.
  ///
  /// Returns:
  ///   A string
  Future<String> modifyAccount(
    String newUsername,
  ) async {
    final token = await Store.getToken();
    if (token == null) {
      showError.createToast("Token must not be null !");
      throw Exception("Token must not be null !");
    }

    final response = await http.put(
      Uri.parse(await _getUrl('/me', {})),
      headers: {
        "Access-Control-Allow-Origin": "*",
        "Authorization": "Bearer $token",
        "Content-Type": "application/json"
      },
      body: jsonEncode({
        "username": newUsername,
      }),
    );
    final body = jsonDecode(response.body);
    if (response.statusCode != 200) {
      final code = response.statusCode;
      showError.createToast("Wrong status code : $code");
      throw Exception(body['error'] ?? "Wrong status code");
    }
    return body['data']['message'];
  }

  /// It deletes the account of the user.
  ///
  /// Returns:
  ///   A string
  Future<String> deleteAccount() async {
    final token = await Store.getToken();
    if (token == null) {
      showError.createToast("Token must not be null !");
      throw Exception("Token must not be null !");
    }

    final response = await http.delete(
      Uri.parse(await _getUrl('/me', {})),
      headers: {
        "Access-Control-Allow-Origin": "*",
        "Authorization": "Bearer $token"
      },
    );
    final body = jsonDecode(response.body);
    if (response.statusCode != 200) {
      final code = response.statusCode;
      showError.createToast("Wrong status code : $code");
      throw Exception(body['error'] ?? "Wrong status code");
    }
    return body['data']['message'];
  }

  /// It gets the user's avatar from the server
  ///
  /// Returns:
  ///   The avatar of the user.
  Future<String> getAccountAvatar() async {
    final token = await Store.getToken();
    if (token == null) {
      showError.createToast("Token must not be null !");
      throw Exception("Token must not be null !");
    }

    final response = await http.get(
      Uri.parse(await _getUrl('/me/avatar', {
        "base": 'http://${Store.serverURL}:${Store.port}',
      })),
      headers: {
        "Access-Control-Allow-Origin": "*",
        "Authorization": "Bearer $token"
      },
    );
    final body = jsonDecode(response.body);
    if (response.statusCode != 200) {
      final code = response.statusCode;
      showError.createToast("Wrong status code : $code");
      throw Exception(body['error'] ?? "Wrong status code");
    }
    return body['data']['uri'];
  }

  /// It uploads the avatar image to the server.
  ///
  /// Args:
  ///   file (PlatformFile): The file to upload.
  ///
  /// Returns:
  ///   A string.
  Future<String> updateAccountAvatar(
    PlatformFile file,
  ) async {
    final token = await Store.getToken();
    if (token == null) {
      showError.createToast("Token must not be null !");
      throw Exception("Token must not be null !");
    }

    final request = http.MultipartRequest(
      'PUT',
      Uri.parse(await _getUrl('/me/avatar', {})),
    );
    request.headers["Access-Control-Allow-Origin"] = "*";
    request.headers["Authorization"] = "Bearer $token";
    request.files.add(await http.MultipartFile.fromPath(
      'avatar',
      file.path ?? "",
    ));
    StreamedResponse response = await request.send();
    final body = await response.stream.bytesToString();
    final json = jsonDecode(body);

    if (response.statusCode != 200) {
      showError.createToast(
          "Something went wrong.\nCannot get upload the image.\nCheck your connection and reload the app.");
      throw Exception(json['error'] ?? "Cannot upload");
    }
    return json['data']['message'];
  }

  /// It gets a new applet from the server
  ///
  /// Returns:
  ///   A pair of an action and a list of areas.
  Future<Pair<Area?, List<Area>>> getNewApplet() async {
    final token = await Store.getToken();
    if (token == null) {
      showError.createToast("Token must not be null !");
      throw Exception("Token must not be null !");
    }

    final response =
        await http.get(Uri.parse(await _getUrl('/applet/new', {})), headers: {
      "Access-Control-Allow-Origin": "*",
      "Content-Type": "application/json",
      "Authorization": "Bearer $token"
    });
    final returnBody = jsonDecode(response.body);

    if (response.statusCode != 200) {
      final code = response.statusCode;
      showError.createToast("Wrong status code : $code");
      throw Exception(returnBody['error'] ?? "Wrong status code");
    }

    Area? action;
    List<Area>? areas = [];
    if (returnBody['data']['action'] != null) {
      action = Area.fromJson(returnBody['data']['action']);
    }
    if (returnBody['data']['reactions'] != null) {
      areas = (returnBody['data']['reactions'] as List)
          .map((e) => Area.fromJson(e))
          .toList();
    }

    return Pair(action, areas);
  }

  /// It adds a state to a new applet
  ///
  /// Args:
  ///   service (String): The service you want to use.
  ///   areaType (String): The type of area you want to add.
  ///   areaItem (String): The name of the item you want to add to the applet.
  ///   areaSettings (Map<String, dynamic>):
  ///
  /// Returns:
  ///   A string
  Future<String> addStateToNewApplet(
    String service,
    String areaType,
    String areaItem,
    Map<String, dynamic> areaSettings,
  ) async {
    final token = await Store.getToken();
    if (token == null) {
      showError.createToast("Token must not be null !");
      throw Exception("Token must not be null !");
    }

    final response = await http.put(Uri.parse(await _getUrl('/applet/new', {})),
        headers: {
          "Access-Control-Allow-Origin": "*",
          "Content-Type": "application/json",
          "Authorization": "Bearer $token"
        },
        body: jsonEncode({
          "service": service,
          "area_type": areaType,
          "area_item": areaItem,
          "area_settings": areaSettings,
        }));
    final returnBody = jsonDecode(response.body);

    if (response.statusCode == 406) {
      Fluttertoast.showToast(
          msg: "Some argument are invalid !",
          toastLength: Toast.LENGTH_SHORT,
          gravity: ToastGravity.BOTTOM,
          backgroundColor: Colors.red,
          textColor: Colors.white,
          fontSize: 32.0);
    }

    if (response.statusCode != 200) {
      final code = response.statusCode;
      showError.createToast("Wrong status code : $code");
      throw Exception(returnBody['error'] ?? "Wrong status code");
    }
    return returnBody['data']['message'];
  }

  /// It sends a POST request to the server with the applet's name, description and public status
  ///
  /// Args:
  ///   name (String): The name of the applet
  ///   description (String): The description of the applet
  ///   isPublic (bool): Whether the applet is public or not
  ///
  /// Returns:
  ///   The id of the applet
  Future<String> submitNewApplet(
    String name,
    String description,
    bool isPublic,
  ) async {
    final token = await Store.getToken();
    if (token == null) {
      showError.createToast("Token must not be null !");
      throw Exception("Token must not be null !");
    }

    final response =
        await http.post(Uri.parse(await _getUrl('/applet/new', {})),
            headers: {
              "Access-Control-Allow-Origin": "*",
              "Content-Type": "application/json",
              "Authorization": "Bearer $token"
            },
            body: jsonEncode({
              "name": name,
              "description": description,
              "public": isPublic,
            }));
    final returnBody = jsonDecode(response.body);
    if (response.statusCode != 200) {
      final code = response.statusCode;
      showError.createToast("Wrong status code : $code");
      throw Exception(returnBody['error'] ?? "Wrong status code");
    }
    return returnBody['data']['message'];
  }

  /// It deletes a new applet
  ///
  /// Args:
  ///   type (String): action or reaction
  ///   number (String): The number of the applet you want to delete.
  ///
  /// Returns:
  ///   A string
  Future<String> deleteStateNewApplet(
    String type, // action or reaction
    String number,
  ) async {
    final token = await Store.getToken();
    if (token == null) {
      showError.createToast("Token must not be null !");
      throw Exception("Token must not be null !");
    }

    final response = await http.delete(
        Uri.parse(await _getUrl('/applet/new', {
          "type": type,
          "number": number,
        })),
        headers: {
          "Access-Control-Allow-Origin": "*",
          "Content-Type": "application/json",
          "Authorization": "Bearer $token"
        });
    final returnBody = jsonDecode(response.body);
    if (response.statusCode != 200) {
      final code = response.statusCode;
      showError.createToast("Wrong status code : $code");
      throw Exception(returnBody['error'] ?? "Wrong status code");
    }
    return returnBody['data']['message'];
  }

  /// It gets all the applets from the server.
  ///
  /// Returns:
  ///   A list of applets
  Future<List<Applet>> getApplets() async {
    final token = await Store.getToken();
    if (token == null) {
      showError.createToast("Token must not be null !");
      throw Exception("Token must not be null !");
    }

    final response = await http.get(
      Uri.parse(await _getUrl('/applet', {})),
      headers: {
        "Access-Control-Allow-Origin": "*",
        "Authorization": "Bearer $token"
      },
    );
    final body = jsonDecode(response.body);
    if (!(response.statusCode == 200 || response.statusCode == 204)) {
      throw Exception(body['error'] ?? "Wrong status code");
    }
    var applets = body['data']['applets'];
    try {
      applets =
          applets.map<Applet>((applet) => Applet.fromJson(applet)).toList();
      return applets;
    } catch (e) {
      logger.e(e);
      showError.createToast(e);
      return [];
    }
  }

  /// It gets an applet from the server
  ///
  /// Args:
  ///   id (String): The id of the applet you want to get.
  ///
  /// Returns:
  ///   A Future<Applet>
  Future<Applet> getApplet(
    String id,
  ) async {
    final token = await Store.getToken();
    if (token == null) {
      showError.createToast("Token must not be null !");
      throw Exception("Token must not be null !");
    }

    final response = await http.get(
      Uri.parse(await _getUrl('/applet/$id', {})),
      headers: {
        "Access-Control-Allow-Origin": "*",
        "Authorization": "Bearer $token"
      },
    );
    final body = jsonDecode(response.body);
    if (response.statusCode != 200) {
      final code = response.statusCode;
      showError.createToast("Wrong status code : $code");
      throw Exception(body['error'] ?? "Wrong status code");
    }
    return Applet.fromJson(body['data']['applet']);
  }

  /// It gets the reactions of an applet
  ///
  /// Args:
  ///   id (String): The id of the applet
  ///
  /// Returns:
  ///   A list of Area objects.
  Future<List<Area>?> getAppletReactions(
    String id,
  ) async {
    final token = await Store.getToken();
    if (token == null) {
      showError.createToast("Token must not be null !");
      throw Exception("Token must not be null !");
    }

    final response = await http.get(
      Uri.parse(await _getUrl('/applet/$id/reactions', {})),
      headers: {
        "Access-Control-Allow-Origin": "*",
        "Authorization": "Bearer $token",
      },
    );
    final body = jsonDecode(response.body);
    if (response.statusCode != 200) {
      final code = response.statusCode;
      showError.createToast("Wrong status code : $code");
      throw Exception(body['error'] ?? "Wrong status code");
    }
    return (body['data']['reactions'] as List)
        .map((e) => Area.fromJson(e))
        .toList();
  }

  /// It modifies an applet
  ///
  /// Args:
  ///   id (String): The id of the applet you want to modify
  ///
  /// Returns:
  ///   A string
  Future<String> modifyApplet(
    String id,
  ) async {
    final token = await Store.getToken();
    if (token == null) {
      showError.createToast("Token must not be null !");
      throw Exception("Token must not be null !");
    }

    final response = await http.put(
      Uri.parse(await _getUrl('/applet/$id', {})),
      headers: {
        "Access-Control-Allow-Origin": "*",
        "Authorization": "Bearer $token"
      },
    );
    final body = jsonDecode(response.body);
    if (response.statusCode != 200) {
      final code = response.statusCode;
      showError.createToast("Wrong status code : $code");
      throw Exception(body['error'] ?? "Wrong status code");
    }
    return body['data']['message'];
  }

  /// It updates the activity of an applet
  ///
  /// Args:
  ///   id (String): The id of the applet you want to update
  ///   active (bool): true/false
  ///
  /// Returns:
  ///   A string
  Future<String> updateAppletActivity(
    String id,
    bool active,
  ) async {
    final token = await Store.getToken();
    if (token == null) {
      showError.createToast("Token must not be null !");
      throw Exception("Token must not be null !");
    }

    final response = await http.patch(
      Uri.parse(await _getUrl('/applet/$id', {
        "active": active,
      })),
      headers: {
        "Access-Control-Allow-Origin": "*",
        "Authorization": "Bearer $token"
      },
    );
    final body = jsonDecode(response.body);
    if (response.statusCode != 200) {
      final code = response.statusCode;
      showError.createToast("Wrong status code : $code");
      throw Exception(body['error'] ?? "Wrong status code");
    }
    return body['data']['message'];
  }

  /// This function deletes an applet
  ///
  /// Args:
  ///   id (String): The id of the applet you want to delete.
  ///
  /// Returns:
  ///   A string
  Future<String> deleteApplet(
    String id,
  ) async {
    final token = await Store.getToken();
    if (token == null) {
      showError.createToast("Token must not be null !");
      throw Exception("Token must not be null !");
    }

    final response = await http.delete(
      Uri.parse(await _getUrl('/applet/$id', {})),
      headers: {
        "Access-Control-Allow-Origin": "*",
        "Authorization": "Bearer $token"
      },
    );
    final body = jsonDecode(response.body);
    if (response.statusCode != 200) {
      final code = response.statusCode;
      showError.createToast("Wrong status code : $code");
      throw Exception(body['error'] ?? "Wrong status code");
    }
    return body['data']['message'];
  }

  /// It starts an applet
  ///
  /// Args:
  ///   id (String): The id of the applet you want to start.
  ///
  /// Returns:
  ///   A string
  Future<String> startApplet(
    String id,
  ) async {
    final token = await Store.getToken();
    if (token == null) {
      showError.createToast("Token must not be null !");
      throw Exception("Token must not be null !");
    }

    final response = await http.put(
      Uri.parse(await _getUrl('/applet/$id/start', {})),
      headers: {
        "Access-Control-Allow-Origin": "*",
        "Authorization": "Bearer $token"
      },
    );

    final body = jsonDecode(response.body);
    if (response.statusCode != 200) {
      final code = response.statusCode;
      showError.createToast("Wrong status code : $code");
      throw Exception(body['error'] ?? "Wrong status code");
    }
    return body['data']['message'];
  }

  /// It stops an applet.
  ///
  /// Args:
  ///   id (String): The id of the applet you want to stop.
  Future<String> stopApplet(
    String id,
  ) async {
    final token = await Store.getToken();
    if (token == null) {
      showError.createToast("Token must not be null !");
      throw Exception("Token must not be null !");
    }

    final response = await http.put(
      Uri.parse(await _getUrl('/applet/$id/stop', {})),
      headers: {
        "Access-Control-Allow-Origin": "*",
        "Authorization": "Bearer $token"
      },
    );

    final body = jsonDecode(response.body);
    if (response.statusCode != 200) {
      final code = response.statusCode;
      showError.createToast("Wrong status code : $code");
      throw Exception(body['error'] ?? "Wrong status code");
    }
    return body['data']['message'];
  }

  /// It fetches the API data from the service
  ///
  /// Args:
  ///   service (String): The name of the service you want to call.
  ///   endpoint (String): The endpoint of the API you want to call.
  ///   method (String): The HTTP method to use (GET, POST, PUT, DELETE, etc.)
  ///   queryParameters (Map<String, String>): A map of query parameters to add to the request.
  ///
  /// Returns:
  ///   ApiData
  Future<ApiData?> fetchServiceAPI(
    String service,
    String endpoint,
    String method,
    Map<String, String> queryParameters,
  ) async {
    final token = await Store.getToken();
    if (token == null) {
      showError.createToast("Token must not be null !");
      throw Exception("Token must not be null !");
    }

    http.Request request;
    request = http.Request(
        method,
        Uri.parse(
            await _getUrl('/services/$service/api$endpoint', queryParameters)));
    request.headers["Access-Control-Allow-Origin"] = "*";
    request.headers["Authorization"] = "Bearer $token";

    final response = await request.send();
    final body = await response.stream.bytesToString();
    final json = jsonDecode(body);

    if (response.statusCode != 200) {
      final code = response.statusCode;
      showError.createToast("Wrong status code : $code");
      throw Exception(json['error'] ?? "Wrong status code");
    }

    try {
      final apiData = ApiData.fromJson(json['data']);
      return apiData;
    } catch (e) {
      logger.e(e);
      return null;
    }
  }
}
