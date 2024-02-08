import 'package:analyzer_plugin/utilities/pair.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fluttertoast/fluttertoast.dart';
import 'package:mobile/classes/server/area.dart';
import 'package:mobile/classes/utils/arguments.dart';
import 'package:mobile/net/api.dart';
import 'package:mobile/widgets/containers/page.dart';
import 'package:mobile/widgets/services/area_settings.dart';

class CreatePage extends ConsumerStatefulWidget {
  const CreatePage({Key? key}) : super(key: key);

  @override
  CreatePageState createState() => CreatePageState();
}

/// It's a stateful widget that displays the first page of the create applet flow
class CreatePageState extends ConsumerState<CreatePage> {
  final selectProvider =
      StateNotifierProvider<ReactionNotifier, Pair<Area?, int>>(
          (ref) => ReactionNotifier());
  late FutureProvider<Pair<Area?, List<Area>?>> createProvider;

  @override
  void initState() {
    super.initState();

    createProvider = FutureProvider<Pair<Area?, List<Area>?>>((ref) async {
      final about = await api.getAbout();
      var areas = await api.getNewApplet();
      final services = about!.services;

      if (areas.first != null) {
        areas.first?.serviceObj = services
            .firstWhere((element) => element.name == areas.first?.service);
      }
      for (var e in areas.last) {
        e.serviceObj =
            services.firstWhere((element) => element.name == e.service);
      }

      if (areas.last.isNotEmpty) {
        ref.read(selectProvider.notifier).set(Pair(areas.last.first, 0));
      }

      return Pair(areas.first, areas.last.isNotEmpty ? areas.last : null);
    });
  }

  @override
  Widget build(BuildContext context) {
    AsyncValue<Pair<Area?, List<Area>?>> areas = ref.watch(createProvider);
    final selected = ref.watch(selectProvider);

    return PageContainer(
      title: 'Create',
      index: 2,
      child: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Text('Help: One tap on an area to start editing it',
                style: TextStyle(fontSize: 20, color: Colors.white)),
            const Text('Help: Double tap on an area to remove it',
                style: TextStyle(fontSize: 20, color: Colors.white)),
            const SizedBox(height: 20),
            SizedBox(
              width: 300,
              height: 100,
              child: Card(
                shape: const RoundedRectangleBorder(
                    borderRadius: BorderRadius.all(Radius.circular(16.0))),
                elevation: 4.0,
                color: areas.asData?.value.first?.serviceObj?.getColor() ??
                    Colors.white,
                child: InkWell(
                  onTap: () {
                    if (areas.asData?.value.first == null) {
                      Navigator.of(context).pushNamed('/create/area',
                          arguments: const AreaArgument(
                              type: 'action', service: null, item: null));
                    } else {
                      // Remove the area and refresh the page
                      Fluttertoast.showToast(
                          msg: "You can have only one action !",
                          toastLength: Toast.LENGTH_SHORT,
                          gravity: ToastGravity.BOTTOM,
                          backgroundColor: Colors.red,
                          textColor: Colors.white,
                          fontSize: 32.0);
                    }
                  },
                  onDoubleTap: () {
                    if (areas.asData?.value.first == null) return;

                    // Add a confirmation dialog

                    showModalBottomSheet(
                        context: context,
                        builder: (context) {
                          return SizedBox(
                            height: 150,
                            child: Column(
                              children: [
                                const Text(
                                  'Are you sure you want to delete this area ?. This will also remove all reactions.',
                                  style: TextStyle(fontSize: 20),
                                ),
                                const SizedBox(height: 20),
                                Row(
                                  mainAxisAlignment:
                                      MainAxisAlignment.spaceEvenly,
                                  children: [
                                    ElevatedButton(
                                      onPressed: () {
                                        Navigator.of(context).pop();
                                      },
                                      child: const Text('Cancel'),
                                    ),
                                    ElevatedButton(
                                      style: ElevatedButton.styleFrom(
                                          backgroundColor: Colors.red),
                                      onPressed: () {
                                        Navigator.of(context).pop();
                                        // Remove the area and refresh the page
                                        api
                                            .deleteStateNewApplet('action', "")
                                            .then((value) => {
                                                  areas = ref
                                                      .refresh(createProvider),
                                                  Fluttertoast.showToast(
                                                      msg:
                                                          "Deletion successful !",
                                                      toastLength:
                                                          Toast.LENGTH_SHORT,
                                                      gravity:
                                                          ToastGravity.BOTTOM,
                                                      backgroundColor:
                                                          Colors.green,
                                                      textColor: Colors.white,
                                                      fontSize: 32.0)
                                                });
                                      },
                                      child: const Text('Delete'),
                                    ),
                                  ],
                                ),
                              ],
                            ),
                          );
                        });
                  },
                  child: Padding(
                    padding: const EdgeInsets.all(16.0),
                    child: Row(
                      mainAxisSize: MainAxisSize.min,
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        if (areas.asData?.value.first != null)
                          Image(
                              width: 50,
                              height: 50,
                              image: NetworkImage(areas
                                      .asData?.value.first?.serviceObj
                                      ?.getAvatar() ??
                                  '')),
                        if (areas.asData?.value.first == null)
                          const Text(
                            'If This',
                            style: TextStyle(fontSize: 20.0),
                          ),
                        if (areas.asData?.value.first != null)
                          Text(
                            areas.asData?.value.first?.name ?? '',
                            style: const TextStyle(
                                fontSize: 15.0, color: Colors.white),
                          ),
                      ],
                    ),
                  ),
                ),
              ),
            ),
            const SizedBox(height: 16.0),
            SizedBox(
              width: 300,
              height: 100,
              child: Card(
                shape: const RoundedRectangleBorder(
                    borderRadius: BorderRadius.all(Radius.circular(16.0))),
                color: selected.first?.serviceObj?.getColor() ?? Colors.white,
                child: InkWell(
                  onTap: () {
                    if (areas.asData?.value.first == null) return;
                    Navigator.of(context).pushNamed('/create/area',
                        arguments: const AreaArgument(
                            type: 'reaction', service: null, item: null));
                  },
                  onDoubleTap: () {
                    if (areas.asData?.value.first == null) return;
                    if (areas.asData?.value.last == null) return;
                    // Remove the area and refresh the page
                    api
                        .deleteStateNewApplet('reaction', "")
                        .then((value) => areas = ref.refresh(createProvider));
                    ref.read(selectProvider.notifier).reset();
                  },
                  child: Padding(
                    padding: const EdgeInsets.all(16.0),
                    child: Row(
                      mainAxisSize: MainAxisSize.min,
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        if (selected.first != null)
                          Image(
                              width: 50,
                              height: 50,
                              image: NetworkImage(
                                  selected.first?.serviceObj?.getAvatar() ??
                                      '')),
                        if (areas.asData?.value.first != null &&
                            selected.first == null)
                          const Text(
                            'Then That',
                            style: TextStyle(fontSize: 20.0),
                          ),
                        if (areas.asData?.value.first == null)
                          const Icon(Icons.lock),
                        if (selected.first != null)
                          Text(
                            '${selected.last}.',
                            style: const TextStyle(
                                fontSize: 15.0, color: Colors.white),
                          ),
                        if (selected.first != null)
                          DropdownButton<Area>(
                            value: selected.first,
                            icon: const Icon(Icons.arrow_downward),
                            iconSize: 24,
                            elevation: 16,
                            style: const TextStyle(color: Colors.deepPurple),
                            underline: Container(
                              height: 2,
                              color: Colors.deepPurpleAccent,
                            ),
                            onChanged: (Area? newValue) {
                              if (newValue == null) return;
                              ref.read(selectProvider.notifier).set(Pair(
                                  newValue,
                                  areas.asData?.value.last?.indexOf(newValue) ??
                                      0));
                            },
                            items: areas.asData?.value.last
                                    ?.map<DropdownMenuItem<Area>>((Area value) {
                                  return DropdownMenuItem<Area>(
                                    value: value,
                                    child: Text(value.name),
                                  );
                                }).toList() ??
                                [],
                          ),
                      ],
                    ),
                  ),
                ),
              ),
            ),
            const SizedBox(height: 16.0),
            if (areas.asData?.value.first != null && selected.first != null)
              ElevatedButton(
                onPressed: () {
                  Navigator.pushNamed(context, '/create/next');
                },
                child: const Text('Next Step'),
              ),
          ],
        ),
      ),
    );
  }
}
