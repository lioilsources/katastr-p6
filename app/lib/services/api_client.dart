import 'dart:convert';
import 'package:http/http.dart' as http;
import '../config/constants.dart';

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

  void dispose() {
    _client.close();
  }
}
