package tui

const (
)

const (
  runeMask  uint64 = 0x000000007FFFFFFF
  hasKey    uint64 = 1 << 32
  modMask   uint64 = 0x0000FFFF00000000
  hasColor  uint64 = 0x00FF000000000000
  hasAttr   uint64 = 0xFF00000000000000
)

const (
	termColorDefault uint8 = iota
	termColorBlack
	termColorRed
	termColorGreen
	termColorYellow
	termColorBlue
	termColorMagenta
	termColorCyan
	termColorWhite
)

const (
  ColorDefault uint64 = iota << 48
  ColorBlack
  ColorRed
  ColorGreen
  ColorYellow
  ColorBlue
  ColorMagenta
  ColorCyan
  ColorWhite
)

const (
  ColorDD uint64 = iota << 48
  ColorBB; ColorBR; ColorBG; ColorBY; ColorBX; ColorBM; ColorBC; ColorBW; // BLACK    B //
  ColorRB; ColorRR; ColorRG; ColorRY; ColorRX; ColorRM; ColorRC; ColorRW; // RED      R //
  ColorGB; ColorGR; ColorGG; ColorGY; ColorGX; ColorGM; ColorGC; ColorGW; // GREEN    G //
  ColorYB; ColorYR; ColorYG; ColorYY; ColorYX; ColorYM; ColorYC; ColorYW; // YELLOW   Y //
  ColorXB; ColorXR; ColorXG; ColorXY; ColorXX; ColorXM; ColorXC; ColorXW; // BLUE     X //
  ColorMB; ColorMR; ColorMG; ColorMY; ColorMX; ColorMM; ColorMC; ColorMW; // MAGENTA  M //
  ColorCB; ColorCR; ColorCG; ColorCY; ColorCX; ColorCM; ColorCC; ColorCW; // CYAN     C //
  ColorWB; ColorWR; ColorWG; ColorWY; ColorWX; ColorWM; ColorWC; ColorWW; // WHITE    W //
)

const (
  cDD uint8 = iota
  cBB; cBR; cBG; cBY; cBX; cBM; cBC; cBW; // BLACK    B //
  cRB; cRR; cRG; cRY; cRX; cRM; cRC; cRW; // RED      R //
  cGB; cGR; cGG; cGY; cGX; cGM; cGC; cGW; // GREEN    G //
  cYB; cYR; cYG; cYY; cYX; cYM; cYC; cYW; // YELLOW   Y //
  cXB; cXR; cXG; cXY; cXX; cXM; cXC; cXW; // BLUE     X //
  cMB; cMR; cMG; cMY; cMX; cMM; cMC; cMW; // MAGENTA  M //
  cCB; cCR; cCG; cCY; cCX; cCM; cCC; cCW; // CYAN     C //
  cWB; cWR; cWG; cWY; cWX; cWM; cWC; cWW; // WHITE    W //
)

const (
	bold   uint8 = 1 << (1 + iota)
	underline
	reverse
  // Blink
  normal uint8 = 0
)

const (
	Bold   uint64 = 1 << (57 + iota)
	Underline
	Reverse
  // Blink
  Normal uint64 = 0
)

const (
	aBold   uint8 = 1 << (1 + iota)
	aUnderline
	aReverse
  // Blink
  aNormal uint8 = 0
)

type ColorPair struct {
  Bg, Fg, Attrs uint8
}

const (
  Ctrl uint64 = 1 << (32 + iota)
  Alt
  Shift

)

const (
  Iotatest uint64 = iota
)

const (
	KeyF1 uint64 = 0xFFFFFFFF - iota
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyInsert
	KeyDelete
	KeyHome
	KeyEnd
	KeyPgup
	KeyPgdn
	KeyArrowUp
	KeyArrowDown
	KeyArrowLeft
	KeyArrowRight
  Mouse
	MouseLeft
	MouseMiddle
	MouseRight
	MouseRelease
	MouseWheelUp
	MouseWheelDown
	KeyCtrlTilde
	KeyCtrl2
	KeyBackspace
	KeyBackspace2
	KeyTab   uint64 = '\t'
	KeyEnter uint64 = '\n'
	KeyEsc   uint64 = 27
)

const (
  Mouse1_PRESSED = 0
  Mouse1_RELEASED
  Mouse1_CLICKED
  Mouse1_DOUBLE_CLICKED
  Mouse1_TRIPLE_CLICKED
  Mouse2_PRESSED
  Mouse2_RELEASED
  Mouse2_CLICKED
  Mouse2_DOUBLE_CLICKED
  Mouse2_TRIPLE_CLICKED
  Mouse3_PRESSED
  Mouse3_RELEASED
  Mouse3_CLICKED
  Mouse3_DOUBLE_CLICKED
  Mouse3_TRIPLE_CLICKED
  Mouse4_PRESSED
  Mouse4_RELEASED
  Mouse4_CLICKED
  Mouse4_DOUBLE_CLICKED
  Mouse4_TRIPLE_CLICKED
)

type Cell struct {
  Ch    rune
  Color uint8
  Attrs uint8
  Touch bool
}

type Gps struct {
  Y, X int
}

type Window struct {
	CurY,    CurX  int    // current cursor position
	Height,  Width int   // maximums of x and y, NOT window size
	InitY,   InitX int   // screen coords of upper-left-hand corner

  Color   uint8
  Attrs   uint8
	BGChar  rune         // current background char

  Buffer [][]Cell

	// bool	_notimeout;	// no time out on function-key entry?

  Looper bool
  Scroll bool
  Touch  bool
  Echo   bool
  Curs   bool
  Delay  bool
  Resize bool
  // Beep   bool
  // Flash  bool

	// bool	_scroll;	// OK to scroll this window?
}
