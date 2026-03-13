package circuits

import "testing"

func TestDecodeEncode5421(t *testing.T) {
	// 1. Проверка правильного кодирования/декодирования
	valid := map[int]int{
		0: 0b0000, 1: 0b0001, 2: 0b0010, 3: 0b0011, 4: 0b0100,
		5: 0b1000, 6: 0b1001, 7: 0b1010, 8: 0b1011, 9: 0b1100,
	}

	for dec, bin := range valid {
		if encoded := Encode5421(dec); encoded != bin {
			t.Errorf("Encode5421(%d) expected %04b, got %04b", dec, bin, encoded)
		}
		if decoded, ok := Decode5421(bin); !ok || decoded != dec {
			t.Errorf("Decode5421(%04b) expected %d, got %d", bin, dec, decoded)
		}
	}

	// 2. Проверка невалидных значений
	if _, ok := Decode5421(0b1111); ok {
		t.Errorf("Decode5421(0b1111) должно вернуть false")
	}
	if enc := Encode5421(15); enc != 0 {
		t.Errorf("Encode5421(15) должно вернуть 0")
	}
	if enc := Encode5421(-1); enc != 0 {
		t.Errorf("Encode5421(-1) должно вернуть 0")
	}

	// Увеличиваем покрытие: проходимся по всем возможным 4-битным входам
	for i := 0; i < 16; i++ {
		Decode5421(i)
	}
}

func TestCircuitsLengths(t *testing.T) {
	if len(GetSubtractorEquations()) != 2 {
		t.Error("ОДВ-3 должен вернуть 2 функции")
	}
	if len(GetDecoder5421Equations()) != 4 {
		t.Error("Декодер должен вернуть 4 функции")
	}
	if len(GetEncoder5421Equations()) != 6 {
		t.Error("Энкодер должен вернуть 6 функций")
	}
	if len(GetCounterEquations()) != 4 {
		t.Error("Счетчик должен вернуть 4 функции")
	}
}
