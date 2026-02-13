import 'package:flutter/material.dart';
import 'package:flutter_map/flutter_map.dart';
import '../config/constants.dart';
import '../services/api_client.dart';

class MapScreen extends StatefulWidget {
  const MapScreen({super.key});

  @override
  State<MapScreen> createState() => _MapScreenState();
}

class _MapScreenState extends State<MapScreen> {
  final MapController _mapController = MapController();
  final ApiClient _apiClient = ApiClient();

  @override
  void dispose() {
    _mapController.dispose();
    _apiClient.dispose();
    super.dispose();
  }

  Future<void> _testApi() async {
    try {
      final version = await _apiClient.getVersion();
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Backend v${version['version']} (${version['backend']})'),
          backgroundColor: Colors.green,
        ),
      );
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Chyba: $e'),
          backgroundColor: Colors.red,
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Katastr Praha 6'),
      ),
      body: FlutterMap(
        mapController: _mapController,
        options: MapOptions(
          initialCenter: praha6Center,
          initialZoom: defaultZoom,
          minZoom: 10,
          maxZoom: 19,
        ),
        children: [
          TileLayer(
            urlTemplate: 'https://tile.openstreetmap.org/{z}/{x}/{y}.png',
            userAgentPackageName: 'cz.katastrp6.app',
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _testApi,
        tooltip: 'Test API',
        child: const Icon(Icons.cloud_done),
      ),
    );
  }
}
