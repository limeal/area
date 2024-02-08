import 'package:flutter/material.dart';
import 'package:convex_bottom_bar/convex_bottom_bar.dart';
import 'package:mobile/widgets/containers/myscaffold.dart';

/// `PageContainer` is a stateful widget that contains a title, an index, and a child widget
class PageContainer extends StatefulWidget {
  final String title;
  final int index;
  final Widget child;

  const PageContainer(
      {super.key,
      required this.title,
      required this.index,
      required this.child});

  @override
  State<PageContainer> createState() {
    return _PageContainerState();
  }
}

/// It's a stateful widget that contains a scaffold with a bottom navigation bar
class _PageContainerState extends State<PageContainer> {
  void _onItemTapped(int index) {
    switch (index) {
      case 0:
        Navigator.pushNamed(context, '/home');
        break;
      case 1:
        Navigator.pushNamed(context, '/applets');
        break;
      case 2:
        Navigator.pushNamed(context, '/create');
        break;
      case 3:
        Navigator.pushNamed(context, '/authorize');
        break;
      case 4:
        Navigator.pushNamed(context, '/profile');
        break;
    }
  }

  @override
  Widget build(BuildContext context) {
    return MyScaffold(
      title: widget.title,
      backButton: false,
      bottomNavigationBar: ConvexAppBar(
        initialActiveIndex: widget.index,
        height: 70,
        items: const [
          TabItem(icon: Icons.home, title: 'Home'),
          TabItem(icon: Icons.grid_view, title: 'Applets'),
          TabItem(icon: Icons.add, title: 'Create'),
          TabItem(icon: Icons.check, title: 'Authorize'),
          TabItem(icon: Icons.person, title: 'Profile'),
        ],
        onTap: _onItemTapped,
      ),
      child: widget.child,
    );
  }
}

/*

  BottomNavigationBar(
          items: const <BottomNavigationBarItem>[
            BottomNavigationBarItem(
              icon: Icon(Icons.home),
              label: 'Home',
            ),
            BottomNavigationBarItem(
              icon: Icon(Icons.grid_view),
              label: 'Applications',
            ),
            BottomNavigationBarItem(
              icon: Icon(Icons.add),
              label: 'New',
            ),
            BottomNavigationBarItem(
                icon: Icon(Icons.shopping_cart), label: 'Explore'),
            BottomNavigationBarItem(
              icon: Icon(Icons.person),
              label: 'Profile',
            ),
          ],
          currentIndex: widget.index,
          selectedItemColor: Colors.amber[800],
          unselectedItemColor: Colors.grey,
          onTap: _onItemTapped)

*/