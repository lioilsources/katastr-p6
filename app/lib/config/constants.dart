import 'package:latlong2/latlong.dart';
import 'package:flutter_map/flutter_map.dart';

const String apiBaseUrl = 'http://localhost:8080';

// Praha 6 bounding box
final praha6Bounds = LatLngBounds(
  const LatLng(50.0650, 14.3000), // SW
  const LatLng(50.1150, 14.4100), // NE
);

// Praha 6 center (roughly Dejvice)
const LatLng praha6Center = LatLng(50.1000, 14.3900);
const double defaultZoom = 14.0;

// Katastrální území Prahy 6 (kódy k ověření přes API)
const Map<String, int> katastralniUzemiP6 = {
  'Dejvice': 729272,
  'Bubeneč': 730122,
  'Břevnov': 729582,
  'Střešovice': 730955,
  'Vokovice': 731001,
  'Veleslavín': 730963,
  'Liboc': 730751,
  'Sedlec': 730904,
  'Ruzyně': 730882,
};
