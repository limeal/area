# Mobile (Area)

## 1. About the mobile application

How to start the mobile application:
```sh
flutter emulator --launch <emulator_name>
flutter run
```

How to build the mobile apk:
```sh
flutter build apk
```

No business logic is provided for this mobile application, all the different services and tools we used
is located in the server folder.

Packages dependencies with used:
- flutter_custom_tabs (For custom tabs support)
- uni_links (For oauth redirection management)
- file_picker (For avatar management)
- fluttertoast (For toast support)
- flutter_form_builder (For better form component)
- flutter_riverpod (For fetching api)
- shared_preferences (For sharing variable and store on disk)
- convex_bottom_bar (For better bottom bar support)
- analyzer_plugin (For pair class)
- getwidget (For custom component)
- logger (For logging)

## 2. Architecture

This is how we decide to develop the application:

- Classes Folder: All the different classes used for the process
- Net: All the different classes that we used for external management (api, oauth)
- Pages: All the different pages of the app
- Store: All the function used to store on disk
- Widgets: All the different components used in pages

Entrypoint file: main.dart
