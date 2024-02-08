import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/classes/server/account.dart';
import 'package:mobile/net/api.dart';
import 'package:mobile/store/store.dart';
import 'package:mobile/widgets/containers/page.dart';
import 'package:mobile/widgets/extensions/hex_color.dart';

class ProfilePage extends ConsumerStatefulWidget {
  const ProfilePage({Key? key}) : super(key: key);

  @override
  ProfilePageState createState() => ProfilePageState();
}

/// `ProfilePageState` is a `ConsumerState` that uses a `FutureProvider` to get the user's account
/// information and display it
class ProfilePageState extends ConsumerState<ProfilePage> {
  final TextEditingController _usernameController = TextEditingController();
  bool _isEditingUsername = false;

  final profileProvider = FutureProvider<Account>((ref) async {
    return await api.getAccount();
  });

  final profileAvatarProvider = FutureProvider<String>((ref) async {
    return await api.getAccountAvatar();
  });

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    AsyncValue<Account> account = ref.watch(profileProvider);
    AsyncValue<String> avatar = ref.watch(profileAvatarProvider);

    return PageContainer(
        title: 'My Profile',
        index: 4,
        child: account.when(
          loading: () => const CircularProgressIndicator(),
          error: (error, stack) => Text('Error: $error'),
          data: (account) => Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.center,
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                avatar.when(
                    loading: () => const CircularProgressIndicator(),
                    error: (error, stack) => Text('Error: $error'),
                    data: (avatar) => TextButton(
                        style: ButtonStyle(
                          backgroundColor: MaterialStateProperty.all<Color>(
                              HexColor.fromHex('#222222')),
                        ),
                        onPressed: () async {
                          FilePickerResult? result =
                              await FilePicker.platform.pickFiles(
                            type: FileType.custom,
                            allowedExtensions: ['jpg', 'png'],
                          );

                          if (result != null) {
                            api.updateAccountAvatar(result.files.single);
                          }
                        },
                        child: CircleAvatar(
                          radius: 64,
                          backgroundImage: NetworkImage(avatar),
                        ))),
                const SizedBox(height: 16),
                const Text(
                  "Email",
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                    color: Colors.white,
                  ),
                ),
                const SizedBox(height: 8),
                Text(
                  account.email,
                  style: const TextStyle(fontSize: 16, color: Colors.white),
                ),
                const SizedBox(height: 16),
                const Text(
                  "Username",
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                    color: Colors.white,
                  ),
                ),
                const SizedBox(height: 8),
                _isEditingUsername
                    ? TextFormField(
                        style: TextStyle(
                          fontSize: 16,
                          color: HexColor.fromHex('#222222'),
                        ),
                        controller: _usernameController
                          ..text = account.username,
                        decoration: InputDecoration(
                          filled: true,
                          fillColor: Colors.white,
                          hintText: "Enter username",
                          suffixIcon: IconButton(
                            icon: const Icon(Icons.check),
                            onPressed: () {
                              api
                                  .modifyAccount(_usernameController.text)
                                  .then((value) {
                                setState(() {
                                  _isEditingUsername = false;
                                });
                              });
                            },
                          ),
                        ),
                      )
                    : Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Text(
                            account.username,
                            style: const TextStyle(
                                fontSize: 20, color: Colors.white),
                          ),
                          const SizedBox(width: 8),
                          IconButton(
                            icon: const Icon(Icons.edit,
                                size: 20, color: Colors.white),
                            onPressed: () {
                              setState(() {
                                _isEditingUsername = true;
                              });
                            },
                          ),
                        ],
                      ),
                // Logout button and delete account button
                ElevatedButton(
                  onPressed: () {
                    Store.setToken('');
                    Navigator.pushNamedAndRemoveUntil(
                        context, '/auth', (route) => false);
                  },
                  child: const Text('Logout'),
                ),
                ElevatedButton(
                  onPressed: () {
                    showModalBottomSheet(
                        context: context,
                        builder: (context) {
                          return SizedBox(
                            height: 150,
                            child: Column(
                              children: [
                                const Text(
                                  'Are you sure you want to delete this account ?',
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
                                        // Remove the area and refresh the page
                                        api.deleteAccount().then((value) {
                                          Store.setToken('');
                                          Navigator.pushNamedAndRemoveUntil(
                                              context,
                                              '/auth',
                                              (route) => false);
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
                  child: const Text('Delete Account'),
                ),
              ],
            ),
          ),
        ));
  }
}
