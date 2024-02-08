import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/classes/server/authenticator.dart';
import 'package:mobile/net/api.dart';
import 'package:mobile/net/oauth.dart';
import 'package:mobile/widgets/containers/myscaffold.dart';
import 'package:mobile/widgets/extensions/hex_color.dart';
import 'package:mobile/classes/utils/error_toast.dart';

class AuthPage extends ConsumerStatefulWidget {
  const AuthPage({Key? key}) : super(key: key);

  @override
  AuthPageState createState() => AuthPageState();
}

/// It's a stateful widget that displays a login page
class AuthPageState extends ConsumerState<AuthPage> {
  final TextEditingController emailController = TextEditingController();
  final TextEditingController passwordController = TextEditingController();
  final authenticatorProvider =
      FutureProvider<List<Authenticator>>((ref) async {
    final server = await api.getAbout();
    if (server == null) {
      throw Exception('Failed to load server');
    }
    return server.authenticators.where((element) => element.enabled).toList();
  });

  @override
  Widget build(BuildContext context) {
    AsyncValue<List<Authenticator>> authenticators =
        ref.watch(authenticatorProvider);
    return authenticators.when(
        loading: () => const CircularProgressIndicator(
              color: Colors.red,
              strokeWidth: 5,
            ),
        error: (error, stack) => Text('Error: $error'),
        data: (authenticators) => MyScaffold(
            title: 'Sign In / Sign Up',
            backButton: false,
            child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                crossAxisAlignment: CrossAxisAlignment.center,
                children: [
                  const Image(
                    image: AssetImage('assets/logo.png'),
                    width: 100,
                    height: 100,
                    fit: BoxFit.scaleDown,
                  ),
                  // Input Email + Password with label
                  Container(
                    padding: const EdgeInsets.all(20),
                    child: Column(
                        mainAxisAlignment: MainAxisAlignment.spaceAround,
                        crossAxisAlignment: CrossAxisAlignment.center,
                        children: [
                          TextField(
                            controller: emailController,
                            decoration: const InputDecoration(
                              border: OutlineInputBorder(),
                              labelText: 'Email',
                              filled: true,
                              fillColor: Color(0xff343434),
                            ),
                            style: TextStyle(
                              color: HexColor.fromHex("#757575"),
                            ),
                          ),
                          TextField(
                            controller: passwordController,
                            obscureText: true,
                            autocorrect: false,
                            enableSuggestions: false,
                            decoration: const InputDecoration(
                              border: OutlineInputBorder(),
                              labelText: 'Password',
                              filled: true,
                              fillColor: Color(0xff343434),
                            ),
                            style: TextStyle(
                              color: HexColor.fromHex("#757575"),
                            ),
                          ),
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                            children: [
                              ElevatedButton(
                                  onPressed: () {
                                    if (emailController.text.isEmpty ||
                                        passwordController.text.isEmpty) {
                                      showError.createToast(
                                          "The email and the password input should not be empty.");
                                    } else {
                                      api
                                          .postAccount(
                                              'login',
                                              emailController.value.text,
                                              passwordController.value.text)
                                          .then((value) =>
                                              Navigator.pushReplacementNamed(
                                                  context, '/home'));
                                    }
                                  },
                                  style: ButtonStyle(
                                    backgroundColor:
                                        MaterialStateProperty.all<Color>(
                                            Colors.white),
                                    foregroundColor:
                                        MaterialStateProperty.all<Color>(
                                            Colors.black),
                                  ),
                                  child: const Text('Login')),
                              ElevatedButton(
                                  onPressed: () {
                                    if (emailController.text.isEmpty ||
                                        passwordController.text.isEmpty) {
                                      showError.createToast(
                                          "The email and the password input should not be empty.");
                                    } else {
                                      api
                                          .postAccount(
                                              'register',
                                              emailController.value.text,
                                              passwordController.value.text)
                                          .then((value) =>
                                              Navigator.pushReplacementNamed(
                                                  context, '/home'));
                                    }
                                  },
                                  style: ButtonStyle(
                                    backgroundColor:
                                        MaterialStateProperty.all<Color>(
                                            Colors.white),
                                    foregroundColor:
                                        MaterialStateProperty.all<Color>(
                                            Colors.black),
                                  ),
                                  child: const Text('Register'))
                            ],
                          ),
                        ]),
                  ),
                  const SizedBox(
                    height: 30,
                    child: Divider(
                      color: Colors.white,
                    ),
                  ),
                  Container(
                      padding: const EdgeInsets.all(20),
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                        children: [
                          for (final auth in authenticators)
                            ElevatedButton(
                              style: ElevatedButton.styleFrom(
                                  fixedSize: const Size(300, 20),
                                  textStyle: const TextStyle(
                                      fontSize: 15, color: Colors.white),
                                  backgroundColor:
                                      HexColor.fromHex(auth.more.color),
                                  shape: RoundedRectangleBorder(
                                      borderRadius: BorderRadius.circular(20.0))
                                  //backgroundColor: Colors.blue,
                                  ),
                              child: Text(auth.name),
                              onPressed: () {
                                OAuthLib.openAuthentifier(true, auth);
                              },
                            ),
                        ],
                      )),
                  ElevatedButton(
                      onPressed: () {
                        Navigator.pushNamedAndRemoveUntil(
                            context, '/', (route) => false);
                      },
                      child: const Text('Return to server selection'))
                ])));
  }
}
