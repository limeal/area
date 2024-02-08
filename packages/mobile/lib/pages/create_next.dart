import 'package:flutter/material.dart';
import 'package:mobile/net/api.dart';
import 'package:mobile/widgets/containers/myscaffold.dart';
import 'package:mobile/widgets/extensions/hex_color.dart';

class CreateNextPage extends StatefulWidget {
  const CreateNextPage({super.key});

  @override
  State<CreateNextPage> createState() => CreateNextPageState();
}

/// It's a stateful widget that has a text field for the applet name, a text field for the applet
/// description, a dropdown for the applet visibility, and a button to submit the applet
class CreateNextPageState extends State<CreateNextPage> {
  final TextEditingController _appNameController = TextEditingController();
  final TextEditingController _appDescriptionController =
      TextEditingController();
  String _selectedVisibility = 'private';
  //bool _isSubmitting = false;

  @override
  Widget build(BuildContext context) {
    return MyScaffold(
      title: 'Create (Step 2)',
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            TextField(
              controller: _appNameController,
              maxLines: 1,
              decoration: const InputDecoration(
                labelText: 'Application Name',
                filled: true,
                fillColor: Color(0xff343434),
              ),
              style: TextStyle(
                color: HexColor.fromHex("#757575"),
              ),
            ),
            const SizedBox(height: 16.0),
            TextField(
              controller: _appDescriptionController,
              maxLines: 8,
              decoration: const InputDecoration(
                labelText: 'Application Description',
                alignLabelWithHint: true,
                border: OutlineInputBorder(),
              ),
            ),
            const SizedBox(height: 16.0),
            DropdownButtonFormField<String>(
              value: _selectedVisibility,
              alignment: Alignment.topLeft,
              onChanged: (newValue) {
                setState(() {
                  _selectedVisibility = newValue!;
                });
              },
              items: <String>['private', 'public'].map((String value) {
                return DropdownMenuItem<String>(
                  value: value,
                  child: Text(value),
                );
              }).toList(),
              decoration: const InputDecoration(
                labelText: 'Visibility',
              ),
            ),
            const SizedBox(height: 16.0),
            ElevatedButton(
              onPressed: () {
                //_isSubmitting = true;
                api
                    .submitNewApplet(
                      _appNameController.text,
                      _appDescriptionController.text,
                      _selectedVisibility == 'private' ? false : true,
                    )
                    .then((value) =>
                        Navigator.of(context).pushReplacementNamed('/applets'));
              },
              child: const Text('Submit'),
            ),
          ],
        ),
      ),
    );
  }
}
