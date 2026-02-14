import 'dart:convert';
import 'package:http/http.dart' as http;
import '../config/constants.dart';
import '../models/parcel.dart';

class ApiClient {
  final String baseUrl;
  final http.Client _client;

  ApiClient({String? baseUrl, http.Client? client})
      : baseUrl = baseUrl ?? apiBaseUrl,
        _client = client ?? http.Client();

  Future<Map<String, dynamic>> getVersion() async {
    final response = await _client.get(Uri.parse('$baseUrl/api/version'));
    if (response.statusCode == 200) {
      return jsonDecode(response.body) as Map<String, dynamic>;
    }
    throw Exception('Failed to load version: ${response.statusCode}');
  }

  Future<Map<String, dynamic>> checkHealth() async {
    final response = await _client.get(Uri.parse('$baseUrl/health'));
    if (response.statusCode == 200) {
      return jsonDecode(response.body) as Map<String, dynamic>;
    }
    throw Exception('Health check failed: ${response.statusCode}');
  }

  /// Search parcels near a WGS-84 point (backend converts to S-JTSK).
  Future<ParcelSearchResult> searchParcelsByPoint(
    double lat,
    double lon, {
    int radius = 5,
  }) async {
    final uri = Uri.parse(
      '$baseUrl/api/parcels/polygon?lat=$lat&lon=$lon&radius=$radius',
    );
    final response = await _client.get(uri);
    if (response.statusCode == 200) {
      return ParcelSearchResult.fromJson(
        jsonDecode(response.body) as Map<String, dynamic>,
      );
    }
    throw Exception(
      'Parcel search failed (${response.statusCode})',
    );
  }

  void dispose() {
    _client.close();
  }
}
