package do

import "testing"

func TestStateMachine(t *testing.T) {
	m := NewStateMachineByMap(map[int]StateFunc[int]{
		1: func(pre int) int {
			return 2
		},
		2: func(pre int) int {
			return 3
		},
		3: func(pre int) int {
			return 1
		},
	})
	Assert(t, m.Run(1), 2)
	Assert(t, m.Run(2), 3)
	Assert(t, m.Run(3), 1)

	// branch
	{
		m := NewStateMachineByMap(map[int]StateFunc[int]{
			1: func(pre int) int {
				return 2
			},
			2: func(pre int) int {
				return 3
			},
			3: func(pre int) int {
				// 根据某些条件决定状态是否修改
				if true {
					return 3
				}
				return 1
			},
		})
		Assert(t, m.Run(1), 2)
		Assert(t, m.Run(2), 3)
		Assert(t, m.Run(3), 3)
	}
}
