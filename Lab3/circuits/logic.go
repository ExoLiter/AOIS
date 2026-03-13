package circuits

import (
	"lab3/qm"
)

// Глобальные константы
const (
	SubtractorInputs = 3
	DecoderInputs    = 4
	EncoderInputs    = 5
	CounterInputs    = 4
	CounterMaxState  = 16
)

type Equation struct {
	Name      string
	SDNF      string
	Minimized string
}

// GetSubtractorEquations - 1. Синтез 1-разрядного вычитателя (ОДВ-3)
func GetSubtractorEquations() []Equation {
	vars := []string{"X1", "X2", "X3"}

	dMinterms := []int{1, 2, 4, 7} // Разность d
	bMinterms := []int{1, 2, 3, 7} // Заем b

	return []Equation{
		{
			Name:      "d (Разность)",
			SDNF:      qm.GenerateSDNF(SubtractorInputs, dMinterms, vars),
			Minimized: qm.Minimize(SubtractorInputs, dMinterms, nil, vars),
		},
		{
			Name:      "b (Заем)",
			SDNF:      qm.GenerateSDNF(SubtractorInputs, bMinterms, vars),
			Minimized: qm.Minimize(SubtractorInputs, bMinterms, nil, vars),
		},
	}
}

// Decode5421 декодирует тетраду 5421
func Decode5421(v int) (int, bool) {
	switch v {
	case 0b0000:
		return 0, true
	case 0b0001:
		return 1, true
	case 0b0010:
		return 2, true
	case 0b0011:
		return 3, true
	case 0b0100:
		return 4, true
	case 0b1000:
		return 5, true
	case 0b1001:
		return 6, true
	case 0b1010:
		return 7, true
	case 0b1011:
		return 8, true
	case 0b1100:
		return 9, true
	}
	return -1, false
}

// Encode5421 кодирует число в код 5421
func Encode5421(v int) int {
	map5421 := []int{0b0000, 0b0001, 0b0010, 0b0011, 0b0100, 0b1000, 0b1001, 0b1010, 0b1011, 0b1100}
	if v >= 0 && v < len(map5421) {
		return map5421[v]
	}
	return 0
}

// GetDecoder5421Equations - 2.1 Декодер 5421 -> Bin (Блок 1)
func GetDecoder5421Equations() []Equation {
	vars := []string{"I3", "I2", "I1", "I0"}
	minterms := make(map[string][]int)
	var dontCares []int

	for i := 0; i < 16; i++ {
		val, ok := Decode5421(i)
		if !ok {
			// Недопустимые коды уходят в Don't Cares
			dontCares = append(dontCares, i)
		} else {
			if (val & 8) != 0 {
				minterms["O3"] = append(minterms["O3"], i)
			}
			if (val & 4) != 0 {
				minterms["O2"] = append(minterms["O2"], i)
			}
			if (val & 2) != 0 {
				minterms["O1"] = append(minterms["O1"], i)
			}
			if (val & 1) != 0 {
				minterms["O0"] = append(minterms["O0"], i)
			}
		}
	}

	outNames := []string{"O3", "O2", "O1", "O0"}
	var result []Equation
	for _, out := range outNames {
		result = append(result, Equation{
			Name:      out,
			Minimized: qm.Minimize(DecoderInputs, minterms[out], dontCares, vars),
		})
	}
	return result
}

// GetEncoder5421Equations - 2.2 Энкодер Bin -> 5421 (Блок 2)
func GetEncoder5421Equations() []Equation {
	vars := []string{"S4", "S3", "S2", "S1", "S0"}
	minterms := make(map[string][]int)
	var dontCares []int

	// Значения > 27 невозможны (макс сумма)
	for i := 28; i < 32; i++ {
		dontCares = append(dontCares, i)
	}

	for i := 0; i < 28; i++ {
		tens := i / 10
		units := i % 10
		t_b := Encode5421(tens)
		u_b := Encode5421(units)

		// T1 (вес 2), T0 (вес 1) для десятков
		if (t_b & 2) != 0 {
			minterms["T1"] = append(minterms["T1"], i)
		}
		if (t_b & 1) != 0 {
			minterms["T0"] = append(minterms["T0"], i)
		}

		// U3..U0 для единиц
		if (u_b & 8) != 0 {
			minterms["U3"] = append(minterms["U3"], i)
		}
		if (u_b & 4) != 0 {
			minterms["U2"] = append(minterms["U2"], i)
		}
		if (u_b & 2) != 0 {
			minterms["U1"] = append(minterms["U1"], i)
		}
		if (u_b & 1) != 0 {
			minterms["U0"] = append(minterms["U0"], i)
		}
	}

	outNames := []string{"T1", "T0", "U3", "U2", "U1", "U0"}
	var result []Equation
	for _, out := range outNames {
		result = append(result, Equation{
			Name:      out,
			Minimized: qm.Minimize(EncoderInputs, minterms[out], dontCares, vars),
		})
	}
	return result
}

// GetCounterEquations - 3. Вычитающий счетчик, 16 состояний, Т-триггер
func GetCounterEquations() []Equation {
	vars := []string{"Q4", "Q3", "Q2", "Q1"}
	minterms := make([][]int, 4)

	for q := 0; q < CounterMaxState; q++ {
		nextQ := (q - 1 + CounterMaxState) % CounterMaxState
		tVals := q ^ nextQ
		for i := 0; i < 4; i++ {
			if ((tVals >> (3 - i)) & 1) == 1 {
				minterms[i] = append(minterms[i], q)
			}
		}
	}

	outNames := []string{"T4", "T3", "T2", "T1"}
	var result []Equation
	for i := 0; i < 4; i++ {
		result = append(result, Equation{
			Name:      outNames[i],
			Minimized: qm.Minimize(CounterInputs, minterms[i], nil, vars),
		})
	}
	return result
}
