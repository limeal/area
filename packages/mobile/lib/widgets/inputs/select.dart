import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:mobile/classes/server/service_area.dart';

/// It's a dropdown field that takes a map entry of a string and an area store item
class AreaSelectField extends StatelessWidget {
  final MapEntry<String, AreaStoreItem> item;

  const AreaSelectField(this.item, {super.key});

  @override
  Widget build(BuildContext context) {
    return FormBuilderDropdown(
        name: item.key,
        items: item.value.values
                ?.map((e) => DropdownMenuItem(value: e, child: Text(e)))
                .toList() ??
            [],
        decoration: InputDecoration(
            labelText: item.key, border: const OutlineInputBorder()),
        validator: (value) {
          if (item.value.isRequired && value == null) {
            return 'Please select an option';
          }
          return null;
        });
  }
}
