import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/classes/utils/arguments.dart';
import 'package:mobile/net/api.dart';
import 'package:mobile/store/store.dart';
import 'package:mobile/widgets/containers/page.dart';

/// `HomePage` is a `ConsumerWidget` that uses a `FutureProvider` to create an authorization provider
class HomePage extends ConsumerWidget {
  final createAuthorizationProvider =
      FutureProvider.family<void, CodeArgument?>((ref, argument) async {
    if (argument != null) {
      await api.postExternalOAuth(
          argument.authenticatorName, argument.code, argument.redirectURL);
      Store.setCurrentAuthenticator('');
    }
  });

  HomePage({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final createAuthorization = ref.watch(createAuthorizationProvider(
        ModalRoute.of(context)!.settings.arguments as CodeArgument?));

    return PageContainer(
        title: 'Home',
        index: 0,
        child: createAuthorization.when(
            loading: () => const Center(
                    child: CircularProgressIndicator(
                  color: Colors.red,
                  strokeWidth: 2,
                )),
            error: (error, stack) => const Text('Success'),
            data: (_) => SafeArea(
                  child: Center(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.center,
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        SizedBox(
                            height: 300,
                            width: double.infinity,
                            child: Column(
                                crossAxisAlignment: CrossAxisAlignment.center,
                                mainAxisAlignment:
                                    MainAxisAlignment.spaceEvenly,
                                children: const [
                                  // Add image here
                                  Image(
                                    image: AssetImage('assets/logo.png'),
                                    width: 100,
                                    height: 100,
                                    fit: BoxFit.scaleDown,
                                  ),
                                  Center(
                                    child: Text(
                                      'Everything Works Better Together',
                                      textAlign: TextAlign.center,
                                      style: TextStyle(
                                        color: Colors.white,
                                        fontSize: 30,
                                        fontWeight: FontWeight.bold,
                                      ),
                                    ),
                                  ),
                                  Center(
                                    child: Text(
                                      'Quickly and easily automate your favorite apps and devices.',
                                      textAlign: TextAlign.center,
                                      style: TextStyle(
                                        fontSize: 15,
                                        color: Colors.white,
                                      ),
                                    ),
                                  )
                                ])),
                        const SizedBox(height: 30),
                        ElevatedButton(
                          onPressed: () {
                            Navigator.pushNamed(context, '/create');
                          },
                          child: const Text('Get Started'),
                        ),
                      ],
                    ),
                  ),
                )));
  }
}
