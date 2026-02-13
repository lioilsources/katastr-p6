import 'package:flutter/material.dart';
import 'screens/map_screen.dart';

class KatastrApp extends StatelessWidget {
  const KatastrApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Katastr Praha 6',
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        colorSchemeSeed: const Color(0xFF1565C0),
        useMaterial3: true,
      ),
      home: const MapScreen(),
    );
  }
}
