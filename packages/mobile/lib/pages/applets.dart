import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/classes/server/applet.dart';
import 'package:mobile/classes/utils/arguments.dart';
import 'package:mobile/net/api.dart';
import 'package:mobile/widgets/containers/page.dart';

/// `AppletsPage` is a `ConsumerWidget` that displays a list of `Applet`s
class AppletsPage extends ConsumerWidget {
  final appletsProvider = FutureProvider<List<Applet>>((ref) async {
    final applets = await api.getApplets();
    return applets;
  });

  AppletsPage({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    AsyncValue<List<Applet>> applets = ref.watch(appletsProvider);
    return PageContainer(
      title: 'Applets',
      index: 1,
      child: applets.when(
        loading: () => const CircularProgressIndicator(),
        error: (error, stack) => Text('Error: $error'),
        data: (applets) => ListView.builder(
          itemCount: applets.length,
          itemBuilder: (context, index) {
            return ListTile(
              title: Text(applets[index].name,
                  style: const TextStyle(
                    fontWeight: FontWeight.bold,
                    color: Colors.white,
                  )),
              subtitle: Text(applets[index].description,
                  style: const TextStyle(
                    fontWeight: FontWeight.bold,
                    color: Colors.white,
                  )),
              trailing: const Icon(Icons.arrow_forward_ios),
              onTap: () async {
                Navigator.of(context).pushNamed('/applet',
                    arguments: AppletArgument(appletId: applets[index].id));
              },
            );
          },
        ),
      ),
    );
  }
}
