import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/classes/utils/arguments.dart';
import 'package:mobile/widgets/services/area_settings.dart';
import 'package:mobile/widgets/services/choose_area.dart';
import 'package:mobile/widgets/services/choose_service.dart';
import 'package:mobile/classes/utils/error_toast.dart';

/// It's a stateful widget that displays a page that allows the user to choose a service, area, and area
/// settings
class AreaPage extends ConsumerStatefulWidget {
  const AreaPage({Key? key}) : super(key: key);

  @override
  AreaPageState createState() => AreaPageState();
}

class AreaPageState extends ConsumerState<AreaPage> {
  @override
  Widget build(BuildContext context) {
    final args = ModalRoute.of(context)!.settings.arguments as AreaArgument?;
    if (args == null) {
      showError.createToast("Invalid area Page");
      return const Text('Invalid area Page');
    }

    if (args.service == null && args.item == null) {
      return ChooseService(args);
    }

    if (args.item == null) {
      return ChooseArea(args);
    }

    return ChooseAreaSettings(args);
  }
}
