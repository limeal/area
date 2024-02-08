import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/classes/server/api_data.dart';

import 'package:mobile/classes/server/service_area.dart';
import 'package:mobile/net/api.dart';

class Wrapper {
  final String service;
  final String endpoint;

  Wrapper(this.service, this.endpoint);
}

class AreaSelectUriField extends ConsumerStatefulWidget {
  final FormBuilderState? formState;
  final List<MapEntry<String, dynamic>> fields;
  final MapEntry<String, AreaStoreItem> item;
  final String service;

  const AreaSelectUriField(this.service, this.formState, this.fields, this.item,
      {Key? key})
      : super(key: key);

  @override
  AreaSelectUriFieldState createState() => AreaSelectUriFieldState();
}

class AreaSelectUriFieldState extends ConsumerState<AreaSelectUriField> {
  late StateProvider<Wrapper> wrapperProvider;
  ApiData? oldData;
  Map<String, dynamic> saved = {};
  late FutureProviderFamily<ApiData, List<MapEntry<String, dynamic>>?>
      apiProvider;

  @override

  /// It takes the endpoint of the service and replaces all variables with the values from the item
  ///
  /// Returns:
  ///   The result of the API call.
  void initState() {
    super.initState();
    wrapperProvider = StateProvider<Wrapper>(
        (ref) => Wrapper(widget.service, widget.item.value.values![0]));
    apiProvider =
        FutureProviderFamily<ApiData, List<MapEntry<String, dynamic>>?>(
            (ref, fields) async {
      List<String> getChangedFields(List<MapEntry<String, dynamic>>? fields) {
        List<String> changedFields = [];
        if (saved.isEmpty) {
          return [];
        }
        for (final field in fields!) {
          if (saved.containsKey(field.key) && saved[field.key] != field.value) {
            changedFields.add(field.key);
          }
        }
        return changedFields;
      }

      final changedFields = getChangedFields(fields);
      if (oldData != null && changedFields.isEmpty) {
        return oldData!;
      }

      final wrapper = ref.read(wrapperProvider);
      // Replace variable in wrapper.endpoint with value from items
      final endpoint = wrapper.endpoint
          .replaceAllMapped(RegExp(r'\${([a-zA-Z:]+)\}'), (match) {
        final key = match.group(1);
        if (key == null) {
          throw Exception('Failed to parse key');
        }
        saved[key] = null;
        for (final field in fields!) {
          if (field.key == key) {
            saved[key] = field.value;
            return field.value.toString();
          }
        }
        return 'default';
      });

      // Remove query parameters and store them in a map
      final queryParameters = <String, String>{};
      final endpointWithoutQuery =
          endpoint.replaceAllMapped(RegExp(r'\?([a-zA-Z0-9=&]+)'), (match) {
        final query = match.group(1);
        if (query == null) {
          throw Exception('Failed to parse query');
        }
        for (final pair in query.split('&')) {
          final split = pair.split('=');
          if (split.length != 2) {
            throw Exception('Failed to parse query');
          }
          queryParameters[split[0]] = split[1];
        }
        return '';
      });

      final result = await api.fetchServiceAPI(
          wrapper.service, endpointWithoutQuery, 'GET', queryParameters);
      if (result == null) {
        throw Exception('Failed to fetch API');
      }
      oldData = result;
      return result;
    });
  }

  @override

  /// > It takes a list of fields, and returns a widget that displays a dropdown with the values of the
  /// first field as the dropdown options, and the values of the second field as the values of the
  /// dropdown options
  ///
  /// Args:
  ///   context (BuildContext): The context of the widget.
  ///
  /// Returns:
  ///   A dropdown menu with the values from the API.
  Widget build(BuildContext context) {
    final apiData = ref.watch(apiProvider(widget.fields));

    return apiData.when(
        data: (data) {
          final values = Map.fromIterables(
              data.elements.map((e) {
                final namePath = data.fields[0].split(':');
                return namePath.fold(e, (value, element) {
                  return value[element];
                });
              }).toSet(),
              data.elements.map((e) {
                final valuePath = data.fields[1].split(':');
                return valuePath.fold(e, (value, element) {
                  return value[element];
                });
              }).toSet());
          final element =
              widget.fields.where((element) => element.key == widget.item.key);

          bool hasValue = true;
          if (element.isNotEmpty) {
            hasValue = values.values.contains(element.first.value);
          }

          if (!hasValue) {
            widget.formState?.fields[widget.item.key]
                ?.didChange(values.values.first);
          }
          return FormBuilderDropdown(
            name: widget.item.key,
            decoration: InputDecoration(labelText: widget.item.key),
            items: values.entries
                .map((e) => DropdownMenuItem(
                      key: ValueKey(e.value),
                      value: e.value,
                      child: Text(e.key),
                    ))
                .toSet()
                .toList(),
            validator: (value) {
              if (value == null) {
                return 'Please select an option';
              }
              return null;
            },
          );
        },
        loading: () => const CircularProgressIndicator(),
        error: (e, s) => Text('Error: $e'));
  }
}
