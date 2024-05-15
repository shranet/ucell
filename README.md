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

log.Println(ut.IsMatch("123"))
log.Println(ut.IsMatch("salom dunyo"))
```

### Qo'shimcha ma'lumot

Dastur har bir qo'shilgan shablondan prefix va suffix larni hisobga olgan holda BTREE yaratadi. Tekshirish jarayonida esa shu BTREE dan foydalanadi. Agar shablonlarda prefix va suffix ishlatilmasa kod ancha tez ishlaydi. Qancha ko'p prefix/suffix ishlatilsa, kod shuncha sekinlashadi.