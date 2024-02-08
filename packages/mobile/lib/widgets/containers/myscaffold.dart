import 'package:flutter/material.dart';
import 'package:getwidget/getwidget.dart';
import 'package:mobile/widgets/extensions/hex_color.dart';

/// `MyScaffold` is a `StatefulWidget` that takes a `title`, `child`, `actions`, `bottomNavigationBar`,
/// and `backButton` as parameters
class MyScaffold extends StatefulWidget {
  final String title;
  final Widget child;
  final List<Widget> actions;
  final Widget? bottomNavigationBar;
  final bool backButton;

  const MyScaffold(
      {Key? key,
      required this.title,
      required this.child,
      this.backButton = true,
      this.actions = const [],
      this.bottomNavigationBar})
      : super(key: key);

  @override
  MyScaffoldState createState() {
    return MyScaffoldState();
  }
}

/// `MyScaffoldState` is a class that extends `State` and implements the `build` method
class MyScaffoldState extends State<MyScaffold> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: GFAppBar(
          automaticallyImplyLeading: false,
          leading: widget.backButton == true
              ? GFIconButton(
                  icon: const Icon(
                    Icons.arrow_back,
                    color: Colors.white,
                  ),
                  onPressed: () {
                    Navigator.pop(context);
                  },
                  type: GFButtonType.transparent,
                )
              : null,
          title: Text(widget.title),
          actions: widget.actions),
      body: widget.child,
      backgroundColor: HexColor.fromHex('#222222'),
      bottomNavigationBar: widget.bottomNavigationBar,
    );
  }
}
