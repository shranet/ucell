# ucell-tpl-match
UCELL kompaniyasi servis SMS shablonlarini tekshirish kutubxonasi

### UCELL qanday tekshiradi
- `%d` - bitta son
- `%d+` - ketma-ket kelgan sonlar
- `%d{1,n}` - ketma-ket n ta sonlar (kamida 1 ta)
- `%w` - son va harfdan iborat so'z
- `%w+` - ketma-ket kelgan son va harflardan iborat so'zlar
- `%w{1,n}` - ketma-ket n ta son va harflardan iborat so'zlar (kamida 1 ta)

### Muammolar
UCELL da shablon yaratguncha `%d`, `%d+`, `%d{1,n}`, `%w`, `%w+`, `%w{1,n}` lar atrofida prefix va suffix kelishi mumkin. Ya'ni:

`prefix%dsuffix` yoki `prefix%w{1,n}suffix`

Agar prefix va suffix `%d`, `%w` bilan birga kelga, tekshirilayotgan matnda ham prefix va suffix yopishib kelishi shart. Masalan: `salom%ddunyo` => `salom123dunyo`

Agar prefix va suffix `%d+`, `%w+` bilan birga kelsa, prefix birinchi son/so'zda, suffix esa oxirgi son/so'zda bo'lishi lozim. Agar son/so'z bitta bo'lsa prefix va suffix yopishib keladi. Masalan: `salom%d+dunyo` => `salom123 456 789dunyo`, `salom123dunyo`

Agar prefix va suffix `%d{1,n}`, `%w{1,n}` bilan birgalikda kelsa, prefix birinchi son/so'zga yopishib, suffix esa oxirgi son/so'zdan bo'shliq bilan ajratilgan holda kelishi kerak. Masalan: `salom%d{1,3}dunyo` => `salom123 456 789 dunyo`, `salom123 dunyo`


### Dasturdan foydalanish

**O'rnatish**

`go get github.com/shranet/ucell-tpl-match`

**Foydalanish**
```go

// Dastlab ucellTemplate obyekt yaratiladi
ut := NewUcellTemplate()

ut.Add("%d")
ut.Add("salom %w")

// Barcha mavjud shablonlar qo'shiladi
// for _, tpl := range templates {
// 	ut.Add(tpl)
// }

// Matnlarni tekshirish
log.Println(ut.IsMatch("123"))
log.Println(ut.IsMatch("salom dunyo"))
```

### Qo'shimcha ma'lumot

Dastur har bir qo'shilgan shablondan prefix va suffix larni hisobga olgan holda BTREE yaratadi. Tekshirish jarayonida esa shu BTREE dan foydalanadi. Agar shablonlarda prefix va suffix ishlatilmasa kod ancha tez ishlaydi. Qancha ko'p prefix/suffix ishlatilsa, kod shuncha sekinlashadi.


$${\color{red} Muhim}$$

UCELL tomonidan berilgan xujjatlarda %d va %w da har xil belgilar aralashib kelishi mumkin deyilgan. Lekin tekshiruvlar natijasida esa a-z (lotin harflari), а-я (krill harflari), raqamlar va bo'shliqdan tashqari barcha belgilar tozalanishi aniqlandi.
`\n` va `\r` lar esa bo'shliq bilan almashtiriladi.

Agar ushbu koida not'g'ri bo'lsa, uni to'g'irlash shart. Chunki ushbu kutubxona a-zа-я0-9 dan tashqari barcha belgilarni tozalab keyin shablonga mosligini tekshiradi.


### Test va Benchmark

```
% go test -bench=. -v

2024/05/16 00:29:35 Fayllarni o'qish
2024/05/16 00:29:35 Fayllarni tahlil qilish va kerakli obyektlarni yaratish
2024/05/16 00:29:36 Barcha shablonlar uchun bitta UcellTemplate yaratish
2024/05/16 00:29:36 Test boshlandi
=== RUN   TestDonalab
--- PASS: TestDonalab (0.06s)
=== RUN   TestHarBirTur
--- PASS: TestHarBirTur (0.33s)
=== RUN   TestAll
--- PASS: TestAll (0.08s)
=== RUN   TestRegex
--- PASS: TestRegex (0.02s)
goos: darwin
goarch: arm64
pkg: github.com/shranet/ucell-tpl-match
BenchmarkAll
BenchmarkAll-10                    71127             16489 ns/op
BenchmarkAllRegex
BenchmarkAllRegex-10                2924            400121 ns/op
BenchmarkUcellTemplate
BenchmarkUcellTemplate-10          85605             13243 ns/op
BenchmarkRegexp
BenchmarkRegexp-10                233331              4879 ns/op
PASS
ok      github.com/shranet/ucell-tpl-match      6.915s
```

**Natijalar:**


| BenchmarkAll     | BenchmarkAllRegex | BenchmarkUcellTemplate | BenchmarkRegexp   |
|------------------|-------------------|------------------------|-------------------|
| 16489 ns/op      | 400121 ns/op      | 13243 ns/op            | 4879 ns/op        |
| 60 646 ta/sekund | 2 499 ta/sekund   | 75 512 ta/sekund       | 204 960 ta/sekund |


Ushbu kutubxona har doim regexlar ro'yxatini shakillantirib, har gal uni for/loop yordamida tekshirgandan ancha tez bo'ladi.

Lekin donalab tekshirishga kelganda regex ancha ustunlikka ega.

**Xulosa**

Agar siz bitta SMS ni N ta shablondan birortasiga mosligini tekshirmoqchi bo'lsangiz, ushbu kutubxonadan foydalaning. Agar siz SMS ni ma'lum 10-100 tagacha shablonga mosligini tekshirmoqchi bo'lsangiz regex ro'yxatdan foydalanganingiz ma'qul.

