import 'package:fluttertoast/fluttertoast.dart';
import 'package:flutter/material.dart';

final showError = ErrorToast();

/// This class is used to create a toast message with a red background and white text
class ErrorToast {
  void createToast(errorMsg) {
    Fluttertoast.showToast(
        msg: errorMsg,
        toastLength: Toast.LENGTH_SHORT,
        gravity: ToastGravity.BOTTOM,
        backgroundColor: Colors.red,
        textColor: Colors.white,
        fontSize: 16.0);
  }
}
