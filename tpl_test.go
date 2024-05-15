package ucell_tpl_match

import (
	"encoding/json"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
)

type testData struct {
	Tpl string `json:"tpl"`
	Ok  string `json:"ok"`
	No  string `json:"no"`
}

var testsData [][]*testData

var ut *ucellTemplate

func init() {

	var testFiles []string

	err := filepath.Walk("./tests", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			testFiles = append(testFiles, path)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	for _, file := range testFiles {
		data, err := os.ReadFile(file)
		if err != nil {
			panic(err)
		}

		tests := []*testData{}
		if err := json.Unmarshal(data, &tests); err != nil {
			panic(err)
		}

		testsData = append(testsData, tests)
	}

	// TestDonalabda biz bildik, UcellTemplate to'g'ri ishlayapti
	// Endi shundan foydalanib, No dagi matnlarga mos shablonlarni qidiramiz.
	//
	//for gidx, tds := range testsData {
	//	for idx, td := range tds {
	//		if gidx*50+idx < 1642 {
	//			continue
	//		}
	//
	//		t.Log(gidx*50 + idx)
	//
	//		for _, _tds := range testsData {
	//			for _, _td := range _tds {
	//				ut := NewUcellTemplate(_td.Tpl)
	//				if ut.IsMatch(td.No) {
	//					t.Log(_td.Tpl)
	//					t.Log(td.No)
	//					t.Log("\n")
	//				}
	//			}
	//		}
	//	}
	//}
	// Qidirishlar natijasida No matnlarga mos quyidagi uchta shablon aniqlandi:
	// %w %w %w %w %w
	// %w %w %w %w
	// %d %d %d
	// va ularni qolgan shablonlar ro'yxatiga qo'shmaymiz
	// shunda ut da IsMatch Ok va No holatlar uchun to'g'ri qiymat qaytarishi lozim bo'ladi

	ut = NewUcellTemplate()
	for _, tds := range testsData {
		for _, td := range tds {
			if isNoTpl(td.Tpl) {
				continue
			}

			ut.Add(td.Tpl)
		}
	}
}

func isNoTpl(tpl string) bool {
	return tpl == "%w %w %w %w %w" || tpl == "%w %w %w %w" || tpl == "%d %d %d"
}

func printData(t *testing.T, idx int, td *testData, ok bool) {
	t.Log(idx)
	t.Log(td.Tpl)
	if ok {
		t.Log(td.Ok)
		t.Error("Shablonlarga mos bo'lishi lozim")
	} else {
		t.Log(td.No)
		t.Error("Shablonlarga mos bo'lishi kerak emas")
	}
}

func TestDonalab(t *testing.T) {
	for _, tds := range testsData {
		for idx, td := range tds {
			ut := NewUcellTemplate(td.Tpl)

			if !ut.IsMatch(td.Ok) {
				printData(t, idx, td, true)
				return
			}

			if ut.IsMatch(td.No) {
				printData(t, idx, td, false)
				return
			}
		}
	}
}

func TestHarBirTur(t *testing.T) {
	for gidx, tds := range testsData {
		ut := NewUcellTemplate()
		for _, td := range tds {
			ut.Add(td.Tpl)
		}

		for idx, td := range tds {
			if !ut.IsMatch(td.Ok) {
				printData(t, gidx*100+idx, td, true)
				return
			}

			if ut.IsMatch(td.No) {
				printData(t, gidx*100+idx, td, false)
				return
			}
		}
	}
}

func TestAll(t *testing.T) {
	for gidx, tds := range testsData {
		for idx, td := range tds {
			if isNoTpl(td.Tpl) {
				//Ushbu shablonlarni tekshirmaymiz
				continue
			}

			if !ut.IsMatch(td.Ok) {
				printData(t, gidx*100+idx, td, true)
				break
			}

			if ut.IsMatch(td.No) {
				printData(t, gidx*100+idx, td, false)
				break
			}
		}
	}
}

func BenchmarkAll(b *testing.B) {
	b.SetParallelism(1)

	correctCount := 0

	for i := 0; i < b.N; i++ {
		tds := testsData[rand.Intn(len(testsData))]
		td := tds[rand.Intn(len(tds))]

		if isNoTpl(td.Tpl) {
			correctCount += 1
		} else {
			if rand.Intn(2) == 1 {
				if ut.IsMatch(td.Ok) {
					correctCount += 1
				}
			} else {
				if !ut.IsMatch(td.No) {
					correctCount += 1
				}
			}
		}
	}

	b.Log(correctCount)
}
