import 'package:flutter_test/flutter_test.dart';
import 'package:katastr_p6/app.dart';

void main() {
  testWidgets('App renders map screen', (WidgetTester tester) async {
    await tester.pumpWidget(const KatastrApp());
    expect(find.text('Katastr Praha 6'), findsOneWidget);
  });
}
