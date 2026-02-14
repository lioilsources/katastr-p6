import 'package:flutter/material.dart';
import 'package:flutter_map/flutter_map.dart';
import 'package:latlong2/latlong.dart';
import '../config/constants.dart';
import '../models/parcel.dart';
import '../services/api_client.dart';
import '../services/location_service.dart';

class MapScreen extends StatefulWidget {
  const MapScreen({super.key});

  @override
  State<MapScreen> createState() => _MapScreenState();
}

class _MapScreenState extends State<MapScreen> {
  final MapController _mapController = MapController();
  final ApiClient _apiClient = ApiClient();
  final LocationService _locationService = LocationService();

  LatLng? _userPosition;
  bool _showCadastralOverlay = true;
  bool _isLocating = false;

  @override
  void dispose() {
    _mapController.dispose();
    _apiClient.dispose();
    super.dispose();
  }

  Future<void> _locateUser() async {
    if (_isLocating) return;
    setState(() => _isLocating = true);
    try {
      final position = await _locationService.getCurrentPosition();
      if (!mounted) return;
      if (position != null) {
        setState(() => _userPosition = position);
        _mapController.move(position, _mapController.camera.zoom);
      } else {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Location unavailable. Check GPS settings.'),
            backgroundColor: Colors.orange,
          ),
        );
      }
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Location error: $e'),
          backgroundColor: Colors.red,
        ),
      );
    } finally {
      if (mounted) setState(() => _isLocating = false);
    }
  }

  Future<void> _onMapTap(TapPosition tapPosition, LatLng point) async {
    try {
      final result = await _apiClient.searchParcelsByPoint(
        point.latitude,
        point.longitude,
      );
      if (!mounted) return;
      if (result.parcels.isEmpty) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('No parcel found at this location.')),
        );
        return;
      }
      _showParcelBottomSheet(result.parcels.first, point);
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Parcel lookup failed: $e'),
          backgroundColor: Colors.red,
        ),
      );
    }
  }

  void _showParcelBottomSheet(Parcel parcel, LatLng tappedPoint) {
    showModalBottomSheet(
      context: context,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
      ),
      builder: (context) {
        return Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Center(
                child: Container(
                  width: 40,
                  height: 4,
                  margin: const EdgeInsets.only(bottom: 16),
                  decoration: BoxDecoration(
                    color: Colors.grey[300],
                    borderRadius: BorderRadius.circular(2),
                  ),
                ),
              ),
              Text(
                'Parcel ${parcel.displayNumber}',
                style: Theme.of(context).textTheme.titleLarge,
              ),
              const SizedBox(height: 8),
              _infoRow('Cadastral area', parcel.cadastralArea.name),
              _infoRow('Area', '${parcel.area} m\u00B2'),
              if (parcel.landType != null)
                _infoRow('Land type', parcel.landType!),
              if (parcel.usageType != null)
                _infoRow('Usage', parcel.usageType!),
              if (parcel.ownershipSheet != null)
                _infoRow('Ownership sheet', parcel.ownershipSheet!),
              _infoRow(
                'Coordinates',
                '${tappedPoint.latitude.toStringAsFixed(6)}, '
                    '${tappedPoint.longitude.toStringAsFixed(6)}',
              ),
              const SizedBox(height: 16),
            ],
          ),
        );
      },
    );
  }

  Widget _infoRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 2),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 140,
            child: Text(
              label,
              style: const TextStyle(
                fontWeight: FontWeight.w500,
                color: Colors.grey,
              ),
            ),
          ),
          Expanded(child: Text(value)),
        ],
      ),
    );
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
          onTap: _onMapTap,
        ),
        children: [
          TileLayer(
            urlTemplate: 'https://tile.openstreetmap.org/{z}/{x}/{y}.png',
            userAgentPackageName: 'cz.katastrp6.app',
          ),
          if (_showCadastralOverlay)
            TileLayer(
              urlTemplate: cuzkWmtsUrl,
              maxZoom: 20,
              minZoom: 10,
              userAgentPackageName: 'cz.katastrp6.app',
              tileBuilder: (context, tileWidget, tile) {
                return Opacity(opacity: 0.7, child: tileWidget);
              },
              errorTileCallback: (tile, error, stackTrace) {},
            ),
          if (_userPosition != null)
            MarkerLayer(
              markers: [
                Marker(
                  point: _userPosition!,
                  width: 20,
                  height: 20,
                  child: Container(
                    decoration: BoxDecoration(
                      color: Colors.blue,
                      shape: BoxShape.circle,
                      border: Border.all(color: Colors.white, width: 2),
                      boxShadow: [
                        BoxShadow(
                          color: Colors.blue.withValues(alpha: 0.3),
                          blurRadius: 8,
                          spreadRadius: 4,
                        ),
                      ],
                    ),
                  ),
                ),
              ],
            ),
        ],
      ),
      floatingActionButton: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          FloatingActionButton(
            heroTag: 'gps',
            onPressed: _locateUser,
            tooltip: 'My location',
            child: _isLocating
                ? const SizedBox(
                    width: 24,
                    height: 24,
                    child: CircularProgressIndicator(
                      strokeWidth: 2,
                      color: Colors.white,
                    ),
                  )
                : const Icon(Icons.my_location),
          ),
          const SizedBox(height: 12),
          FloatingActionButton.small(
            heroTag: 'layers',
            onPressed: () {
              setState(() => _showCadastralOverlay = !_showCadastralOverlay);
            },
            tooltip: 'Toggle cadastral overlay',
            child: Icon(
              _showCadastralOverlay ? Icons.layers : Icons.layers_clear,
            ),
          ),
        ],
      ),
    );
  }
}
