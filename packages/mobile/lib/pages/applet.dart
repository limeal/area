import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/classes/server/applet.dart';
import 'package:mobile/classes/utils/arguments.dart';
import 'package:mobile/net/api.dart';
import 'package:mobile/widgets/containers/myscaffold.dart';

class AppletPage extends ConsumerStatefulWidget {
  const AppletPage({Key? key}) : super(key: key);

  @override
  AppletPageState createState() => AppletPageState();
}

/// `AppletPageState` is a `ConsumerState` that uses a `FutureProvider` to fetch an `Applet` from the
/// API and display it
class AppletPageState extends ConsumerState<AppletPage> {
  /// Creating a `FutureProvider` that will fetch an `Applet` from the API.
  final appletProvider = FutureProvider.family<Applet, String>((ref, id) async {
    return await api.getApplet(id);
  });

  @override
  Widget build(BuildContext context) {
    final args = ModalRoute.of(context)!.settings.arguments as AppletArgument;
    var applet = ref.watch(appletProvider(args.appletId));
    return applet.when(
      data: (applet) {
        return MyScaffold(
          title: applet.name,
          child: Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Text(applet.description),
                const SizedBox(height: 20),
                ElevatedButton(
                  onPressed: () {
                    if (applet.status == 'stopped') {
                      api.startApplet(applet.id).then((value) => ref
                          .refresh(appletProvider(args.appletId).future)
                          .then((value) => applet = value));
                    } else {
                      api.stopApplet(applet.id).then((value) => ref
                          .refresh(appletProvider(args.appletId).future)
                          .then((value) => applet = value));
                    }
                  },
                  child: Text(applet.status == 'stopped' ? 'Start' : 'Stop'),
                ),
                const SizedBox(height: 10),
                ElevatedButton(
                  onPressed: () {
                    api.deleteApplet(applet.id).then((value) =>
                        Navigator.of(context).pushNamedAndRemoveUntil(
                            '/applets', (Route<dynamic> route) => false));
                  },
                  child: const Text('Delete'),
                ),
              ],
            ),
          ),
        );
      },
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (error, stack) => const Center(child: Text('Error')),
    );
  }
}
