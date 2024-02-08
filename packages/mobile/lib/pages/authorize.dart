import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/classes/server/service.dart';
import 'package:mobile/classes/utils/arguments.dart';
import 'package:mobile/net/api.dart';
import 'package:mobile/net/oauth.dart';
import 'package:mobile/store/store.dart';
import 'package:mobile/widgets/containers/page.dart';

class AuthorizePage extends ConsumerStatefulWidget {
  const AuthorizePage({Key? key}) : super(key: key);

  @override
  AuthorizePageState createState() => AuthorizePageState();
}

/// It's a stateful widget that displays a list of services that the user can authorize
class AuthorizePageState extends ConsumerState<AuthorizePage> {
  final authorizeProvider = FutureProvider.family<List<Service>, CodeArgument?>(
      (ref, argument) async {
    if (argument != null) {
      await api.createAuthorization(
          argument.authenticatorName, argument.code, argument.redirectURL);
      Store.setCurrentAuthenticator('');
    }
    final about = await api.getAbout();
    final umap = await api.getServicesAuthorized();

    return about!.services.where((service) {
      final tmp = umap[service.name];
      if (tmp == null) return false;
      return tmp == false;
    }).toList();
  });

  @override
  Widget build(BuildContext context) {
    AsyncValue<List<Service>> services = ref.watch(authorizeProvider(
        ModalRoute.of(context)!.settings.arguments as CodeArgument?));

    return PageContainer(
      title: 'Authorize services',
      index: 3,
      child: services.when(
        loading: () => const CircularProgressIndicator(),
        error: (error, stack) => Text('Error: $error'),
        data: (services) => ListView.builder(
          itemCount: services.length,
          itemBuilder: (context, index) {
            return ListTile(
              title: Text(services[index].name,
                  style: const TextStyle(
                    fontWeight: FontWeight.bold,
                    color: Colors.white,
                  )),
              subtitle: Text(services[index].description ?? 'Basic description',
                  style: const TextStyle(
                    fontWeight: FontWeight.bold,
                    color: Colors.white,
                  )),
              leading: Image(image: NetworkImage(services[index].getAvatar())),
              trailing: const Icon(Icons.arrow_forward_ios),
              onTap: () async {
                OAuthLib.openAuthentifier(
                    false, services[index].authenticator!);
              },
            );
          },
        ),
      ),
    );
  }
}
