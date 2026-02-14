class CadastralArea {
  final int code;
  final String name;

  const CadastralArea({required this.code, required this.name});

  factory CadastralArea.fromJson(Map<String, dynamic> json) {
    return CadastralArea(
      code: json['kod'] as int? ?? 0,
      name: json['nazev'] as String? ?? '',
    );
  }
}

class Parcel {
  final int id;
  final int baseNumber;
  final int? subdivision;
  final String numberingType;
  final CadastralArea cadastralArea;
  final int area;
  final String? landType;
  final String? usageType;
  final String? ownershipSheet;

  const Parcel({
    required this.id,
    required this.baseNumber,
    this.subdivision,
    required this.numberingType,
    required this.cadastralArea,
    required this.area,
    this.landType,
    this.usageType,
    this.ownershipSheet,
  });

  factory Parcel.fromJson(Map<String, dynamic> json) {
    return Parcel(
      id: json['id'] as int? ?? 0,
      baseNumber: json['kmenoveCislo'] as int? ?? 0,
      subdivision: json['poddeleni'] as int?,
      numberingType: json['druhCislovani'] as String? ?? '',
      cadastralArea: json['katastralniUzemi'] != null
          ? CadastralArea.fromJson(
              json['katastralniUzemi'] as Map<String, dynamic>)
          : const CadastralArea(code: 0, name: ''),
      area: json['vymera'] as int? ?? 0,
      landType: json['druhPozemku'] as String?,
      usageType: json['zpusobVyuziti'] as String?,
      ownershipSheet: json['cisloLV'] as String?,
    );
  }

  /// Human-readable parcel number (e.g., "1234/5" or "1234").
  String get displayNumber {
    if (subdivision != null) {
      return '$baseNumber/$subdivision';
    }
    return '$baseNumber';
  }
}

class ParcelSearchResult {
  final List<Parcel> parcels;
  final int total;

  const ParcelSearchResult({required this.parcels, required this.total});

  factory ParcelSearchResult.fromJson(Map<String, dynamic> json) {
    return ParcelSearchResult(
      parcels: (json['parcely'] as List<dynamic>?)
              ?.map((e) => Parcel.fromJson(e as Map<String, dynamic>))
              .toList() ??
          [],
      total: json['total'] as int? ?? 0,
    );
  }
}
