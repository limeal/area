import 'package:analyzer_plugin/utilities/pair.dart';
import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:getwidget/getwidget.dart';
import 'package:mobile/classes/server/area.dart';
import 'package:mobile/classes/server/service_area.dart';
import 'package:mobile/classes/utils/arguments.dart';
import 'package:mobile/net/api.dart';
import 'package:mobile/widgets/containers/myscaffold.dart';
import 'package:mobile/widgets/inputs/select.dart';
import 'package:mobile/widgets/inputs/select_uri.dart';
import 'package:mobile/widgets/inputs/text.dart';

class ChooseAreaSettings extends ConsumerStatefulWidget {
  final AreaArgument argument;
  final _formKey = GlobalKey<FormBuilderState>();

  ChooseAreaSettings(this.argument, {Key? key}) : super(key: key);

  @override
  ChooseAreaSettingsState createState() => ChooseAreaSettingsState();
}

/// It's a stateful widget that displays a form with fields that are dynamically generated based on the
/// service and area selected
class ChooseAreaSettingsState extends ConsumerState<ChooseAreaSettings> {
  final optionalFields = <String>[];
  final fieldsProvider =
      Provider.family<List<MapEntry<String, AreaStoreItem>>, AreaArgument>(
          (ref, argument) {
    if (argument.service == null) throw Exception('Must select a service');
    if (argument.item == null) throw Exception('Must select an area');
    return argument.item?.store?.entries.toList() ?? [];
  });
  final fieldProvider =
      StateNotifierProvider<FieldsNotifier, List<MapEntry<String, dynamic>>>(
          (ref) => FieldsNotifier());

  List<Widget> getFields(List<MapEntry<String, AreaStoreItem>> fields) {
    List<Widget> widgets = [];

    bool canBeDisplayed(MapEntry<String, AreaStoreItem> field) {
      if (field.value.needFields == null) return true;
      for (var needField in field.value.needFields!) {
        if (widget._formKey.currentState?.fields[needField]?.value == null ||
            widget._formKey.currentState?.fields[needField]?.value == '') {
          return false;
        }
      }
      return true;
    }

    final filledFields = ref.read(fieldProvider.notifier).get();
    fields.sort((a, b) => a.value.priority.compareTo(b.value.priority));
    for (var field in fields) {
      if (!canBeDisplayed(field)) continue;
      if (!field.value.isRequired && !optionalFields.contains(field.key)) {
        continue;
      }
      if (field.value.type == 'select') {
        widgets.add(AreaSelectField(field));
      } else if (field.value.type == 'select_uri') {
        widgets.add(AreaSelectUriField(widget.argument.service?.name ?? '',
            widget._formKey.currentState, filledFields, field));
      } else {
        widgets.add(AreaTextField(field));
      }
    }
    return widgets;
  }

  @override
  Widget build(BuildContext context) {
    var fields = ref.watch(fieldsProvider(widget.argument));
    final nonRequiredFields =
        fields.where((element) => element.value.isRequired == false).toList();
    return MyScaffold(
        title: widget.argument.item?.name ?? 'Title',
        child: FormBuilder(
            key: widget._formKey,
            onChanged: () async {
              final values = widget._formKey.currentState?.fields.entries
                      .map((e) => MapEntry(e.key, e.value.value))
                      .toList() ??
                  [];
              ref.read(fieldProvider.notifier).set(values);
              fields = ref.refresh(fieldsProvider(widget.argument));
            },
            child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                crossAxisAlignment: CrossAxisAlignment.center,
                children: [
                  ...getFields(fields),
                  if (optionalFields.length != nonRequiredFields.length)
                    DropdownButton(
                        hint: const Text('Add parameter'),
                        items: nonRequiredFields
                            .where((element) =>
                                !optionalFields.contains(element.key))
                            .map((e) => DropdownMenuItem(
                                value: e.key, child: Text(e.key)))
                            .toList(),
                        onChanged: (value) {
                          if (value == null) return;
                          optionalFields.add(value);
                          setState(() {});
                        }),
                  ElevatedButton(
                    onPressed: () {
                      widget._formKey.currentState!.save();
                      if (widget._formKey.currentState!.validate()) {
                        final formData = widget._formKey.currentState?.value;
                        api
                            .addStateToNewApplet(
                                widget.argument.service?.name ?? '',
                                widget.argument.type,
                                widget.argument.item?.name ?? '',
                                formData!)
                            .then((value) => Navigator.of(context)
                                .pushReplacementNamed('/create'));
                      }
                    },
                    child: const Text('Submit'),
                  ),
                ])));
  }
}

/// `FieldsNotifier` is a `StateNotifier` that holds a list of `MapEntry`s
class FieldsNotifier extends StateNotifier<List<MapEntry<String, dynamic>>> {
  FieldsNotifier() : super([]);

  void set(List<MapEntry<String, dynamic>> value) {
    state = value;
  }

  List<MapEntry<String, dynamic>> get() {
    return state;
  }
}

/// It's a state notifier that holds a pair of an Area and an int
class ReactionNotifier extends StateNotifier<Pair<Area?, int>> {
  ReactionNotifier() : super(Pair(null, 0));

  void set(Pair<Area?, int> value) {
    state = value;
  }

  Pair<Area?, int> get() {
    return state;
  }

  void reset() {
    state = Pair(null, 0);
  }
}
