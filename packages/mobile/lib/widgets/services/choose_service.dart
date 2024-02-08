import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/widgets/containers/myscaffold.dart';

import 'package:mobile/classes/server/service.dart';
import 'package:mobile/classes/utils/arguments.dart';
import 'package:mobile/net/api.dart';

/// It displays a list of services that the user has authorized, and when the user taps on one of them,
/// it navigates to the next page
class ChooseService extends ConsumerWidget {
  final AreaArgument argument;
  final serviceProvider = FutureProvider<List<Service>>((ref) async {
    final about = await api.getAbout();
    final umap = await api.getServicesAuthorized();

    return about!.services.where((service) {
      final tmp = umap[service.name];
      if (tmp == null) return false;
      return tmp == true;
    }).toList();
  });

  ChooseService(this.argument, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    AsyncValue<List<Service>> services = ref.watch(serviceProvider);
    return MyScaffold(
      title: "Services",
      child: services.when(
        loading: () => const CircularProgressIndicator(),
        error: (error, stack) => Text('Error: $error'),
        data: (services) => ListView.builder(
          itemCount: services.length,
          itemBuilder: (context, index) {
            if (argument.type == 'action' &&
                services[index]
                        .actions
                        ?.where((element) => !element.wip)
                        .isEmpty ==
                    true) {
              return const SizedBox.shrink();
            } else if (argument.type == 'reaction' &&
                services[index]
                        .reactions
                        ?.where((element) => !element.wip)
                        .isEmpty ==
                    true) {
              return const SizedBox.shrink();
            }
            return ListTile(
              leading: Image(
                image: NetworkImage(services[index].getAvatar()),
                width: 50,
                height: 50,
              ),
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
              onTap: () {
                Navigator.of(context).pushReplacementNamed('/create/area',
                    arguments: AreaArgument(
                        type: argument.type,
                        service: services[index],
                        item: null));
              },
            );
          },
        ),
      ),
    );
  }
}
