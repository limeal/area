import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/classes/server/service_area.dart';
import 'package:mobile/classes/utils/arguments.dart';
import 'package:mobile/widgets/containers/myscaffold.dart';
import 'package:mobile/net/api.dart';

/// It's a widget that displays a list of areas for a service
class ChooseArea extends ConsumerWidget {
  final AreaArgument argument;
  final areasProvider =
      Provider.family<List<ServiceArea>, AreaArgument>((ref, argument) {
    if (argument.service == null) throw Exception('Must select a service');
    if (argument.type == 'action') {
      return argument.service?.actions
              ?.where((element) => !element.wip)
              .toList() ??
          [];
    }
    return argument.service?.reactions
            ?.where((element) => !element.wip)
            .toList() ??
        [];
  });

  ChooseArea(this.argument, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final areas = ref.watch(areasProvider(argument));
    return MyScaffold(
        title: argument.type == 'action' ? "Actions" : "Reactions",
        child: ListView.builder(
          itemCount: areas.length,
          itemBuilder: (context, index) {
            return ListTile(
              leading: Image(
                image: NetworkImage(argument.service?.getAvatar() ?? ''),
                width: 50,
                height: 50,
              ),
              title: Text(areas[index].name,
                  style: const TextStyle(
                    fontWeight: FontWeight.bold,
                    color: Colors.white,
                  )),
              subtitle: Text(areas[index].description,
                  style: const TextStyle(
                    fontWeight: FontWeight.bold,
                    color: Colors.white,
                  )),
              onTap: () {
                if (areas[index].store != null) {
                  Navigator.of(context).pushReplacementNamed('/create/area',
                      arguments: AreaArgument(
                          type: argument.type,
                          service: argument.service,
                          item: areas[index]));
                } else {
                  api.addStateToNewApplet(
                      argument.service?.name ?? '',
                      argument.type,
                      areas[index].name, {}).then((value) => Navigator.of(
                          context)
                      .pushReplacementNamed('/create'));
                }
              },
            );
          },
        ));
  }
}
