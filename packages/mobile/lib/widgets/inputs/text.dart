import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/classes/server/service_area.dart';
import 'package:mobile/net/api.dart';
import 'package:mobile/widgets/extensions/hex_color.dart';

class AreaTextField extends ConsumerStatefulWidget {
  final MapEntry<String, AreaStoreItem> item;

  const AreaTextField(this.item, {super.key});

  @override
  AreaTextFieldState createState() => AreaTextFieldState();
}

/// It's a stateful widget that renders a text field and a dropdown button. The dropdown button is
/// populated with a list of components that are available for the current action
class AreaTextFieldState extends ConsumerState<AreaTextField> {
  final TextEditingController controller = TextEditingController();
  final getComponentsProvider = FutureProvider<List<String>>((ref) async {
    final about = await api.getAbout();
    final response = await api.getNewApplet();

    if (response.first == null) {
      return [];
    }

    final area = about?.services
        .where((element) => element.name == response.first!.service)
        .first
        .actions
        ?.where((element) => element.name == response.first!.name);
    return area?.first.components ?? [];
  });

  @override
  Widget build(BuildContext context) {
    final components = ref.watch(getComponentsProvider);
    return Column(
      children: [
        FormBuilderTextField(
            name: widget.item.key,
            controller: controller,
            style: TextStyle(
              color: HexColor.fromHex("#757575"),
            ),
            decoration: InputDecoration(
              labelText: widget.item.key,
              border: const OutlineInputBorder(),
              filled: true,
              fillColor: const Color(0xff343434),
            ),
            validator: (value) {
              if (widget.item.value.isRequired &&
                  (value == null || value.isEmpty)) {
                return 'Please enter some text';
              }
              return null;
            }),
        components.when(
          data: (data) {
            return DropdownButton(
              hint: const Text('Add component'),
              items: data
                  .map((e) => DropdownMenuItem(value: e, child: Text(e)))
                  .toList(),
              onChanged: (value) => {
                controller.text += "{{$value}}",
              },
            );
          },
          loading: () => const SizedBox(),
          error: (error, stack) => const SizedBox(),
        ),
      ],
    );
  }
}
