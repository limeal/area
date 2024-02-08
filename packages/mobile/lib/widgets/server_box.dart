import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/net/api.dart';
import 'package:mobile/store/store.dart';
import 'package:fluttertoast/fluttertoast.dart';
import 'package:mobile/classes/utils/error_toast.dart';

/// It's a widget that displays a text field and a button
class ServerBox extends ConsumerWidget {
  final TextEditingController serverController = TextEditingController();
  final TextEditingController portController = TextEditingController();

  ServerBox({super.key});

  get showT => showError.createToast("Cannot connect to server");

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Container(
      alignment: Alignment.center,
      decoration: const BoxDecoration(
        image: DecorationImage(
          image: AssetImage('assets/logo.png'),
          fit: BoxFit.cover,
        ),
      ),
      child: Material(
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
            Center(
              child: SizedBox(
                width: 300,
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: <Widget>[
                    const Text(
                      'Area - Server connection',
                      style: TextStyle(fontSize: 30),
                    ),
                    const SizedBox(height: 20),
                    TextField(
                      controller: serverController
                        ..text = Store.defaultServerURL,
                      decoration: const InputDecoration(
                        border: OutlineInputBorder(),
                        labelText: 'Server address',
                      ),
                      onChanged: (String value) {
                        Store.setServerURL(serverController.value.text);
                      },
                    ),
                    const SizedBox(height: 20),
                    TextField(
                      controller: portController..text = Store.defaultPort,
                      decoration: const InputDecoration(
                        border: OutlineInputBorder(),
                        labelText: 'Port',
                      ),
                      onChanged: (String value) {
                        Store.setPort(portController.value.text);
                      },
                    ),
                    ElevatedButton(
                      onPressed: () {
                        api
                            .getServer()
                            .then((value) => {
                                  if (value != "Area Server - v0.1.0")
                                    {
                                      showError.createToast(
                                          "That's not a valid area server !")
                                    }
                                  else
                                    {
                                      Fluttertoast.showToast(
                                          msg: "Connected to server !",
                                          toastLength: Toast.LENGTH_SHORT,
                                          gravity: ToastGravity.BOTTOM,
                                          backgroundColor: Colors.green,
                                          textColor: Colors.white,
                                          fontSize: 32.0),
                                      Navigator.pushReplacementNamed(
                                          context, '/auth')
                                    }
                                })
                            .catchError((error) => showT);
                      },
                      child: const Text('Connect'),
                    ),
                  ],
                ),
              ),
            ),
          ])),
    );
  }
}
