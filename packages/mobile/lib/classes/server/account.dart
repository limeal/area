/// `Account` is a class that has three properties: `id`, `username`, and `email`
class Account {
  final String id;
  final String username;
  final String email;

  const Account(
      {required this.id, required this.username, required this.email});

  factory Account.fromJson(Map<String, dynamic> json) {
    return Account(
      id: json['id'],
      username: json['username'],
      email: json['email'],
    );
  }
}
