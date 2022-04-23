package engine

const (
	MIDI_CHAN_MIDDLE       = 0
	MIDI_CHAN_LEFT         = 1
	MIDI_CHAN_RIGHT        = 2
	MIDI_CHAN_HOT_CUE_LEFT = 6
	MIDI_CTRL_RIM          = 9
	MIDI_CTRL_PAN          = 10

	MIDI_VAL_PAN_CCW = 127
	MIDI_VAL_PAN_CW  = 1

	MIDI_CTRL_BANK_SELECT_MSB    = 0 // Volume
	MIDI_CTRL_BANK_SELECT_LSB    = 32
	MIDI_CTRL_BALANCE_MSB        = 8
	MIDI_CTRL_BALANCE_LSB        = 40
	MIDI_CTRL_BREATH_CONTROL_MSB = 2 // Bass/filter
	MIDI_CTRL_BREATH_CONTROL_LSB = 34
	MIDI_CTRL_MSB                = 3 // Master
	MIDI_CTRL_LSB                = 35

	MIDI_KEY_SYNC     = 5
	MIDI_KEY_PLAY     = 7
	MIDI_KEY_LED_SHOW = 36

	MIDI_KEY_BANK_1 = 0
	MIDI_KEY_BANK_2 = 1
	MIDI_KEY_BANK_3 = 2
	MIDI_KEY_BANK_4 = 3
)
