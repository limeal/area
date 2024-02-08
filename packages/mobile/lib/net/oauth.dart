import 'package:flutter/material.dart';
import 'package:mobile/classes/server/authenticator.dart';
import 'package:flutter_custom_tabs/flutter_custom_tabs.dart';
import 'package:mobile/store/store.dart';

/// It opens a custom tab with the authorization URL of the authenticator
class OAuthLib {
  static void openAuthentifier(bool mode, Authenticator authenticator) async {
    Store.setRegistrationMode(mode);
    Store.setCurrentAuthenticator(authenticator.name);
    await launch(
      '${authenticator.authorizationUri}&redirect_uri=${Store.redirectURL}',
      customTabsOption: CustomTabsOption(
        toolbarColor: const Color(0xFF1E1E1E),
        enableDefaultShare: true,
        enableUrlBarHiding: true,
        showPageTitle: true,
        animation: CustomTabsSystemAnimation.slideIn(),
        // ignore: prefer_const_literals_to_create_immutables
        extraCustomTabs: <String>[
          'org.mozilla.firefox',
          'com.microsoft.emmx',
        ],
      ),
    );
  }
}
