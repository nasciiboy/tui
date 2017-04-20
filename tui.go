package tui

import (
  "time"

  "github.com/nasciiboy/tui/term"
)

var Colors [256]ColorPair

func Init() (*Window, error) {
  err := term.Init()
  if err != nil { return nil, err }

  Colors = [256]ColorPair{
    { Bg: termColorDefault, Fg: termColorDefault },
    { Bg: termColorBlack,   Fg: termColorBlack }, { Bg: termColorBlack, Fg: termColorRed }, { Bg: termColorBlack, Fg: termColorGreen }, { Bg: termColorBlack, Fg: termColorYellow }, { Bg: termColorBlack, Fg: termColorBlue }, { Bg: termColorBlack, Fg: termColorMagenta }, { Bg: termColorBlack, Fg: termColorCyan }, { Bg: termColorBlack, Fg: termColorWhite },
    { Bg: termColorRed,     Fg: termColorBlack }, { Bg: termColorRed, Fg: termColorRed }, { Bg: termColorRed, Fg: termColorGreen }, { Bg: termColorRed, Fg: termColorYellow }, { Bg: termColorRed, Fg: termColorBlue }, { Bg: termColorRed, Fg: termColorMagenta }, { Bg: termColorRed, Fg: termColorCyan }, { Bg: termColorRed, Fg: termColorWhite },
    { Bg: termColorGreen,   Fg: termColorBlack }, { Bg: termColorGreen, Fg: termColorRed }, { Bg: termColorGreen, Fg: termColorGreen }, { Bg: termColorGreen, Fg: termColorYellow }, { Bg: termColorGreen, Fg: termColorBlue }, { Bg: termColorGreen, Fg: termColorMagenta }, { Bg: termColorGreen, Fg: termColorCyan }, { Bg: termColorGreen, Fg: termColorWhite },
    { Bg: termColorYellow,  Fg: termColorBlack }, { Bg: termColorYellow, Fg: termColorRed }, { Bg: termColorYellow, Fg: termColorGreen }, { Bg: termColorYellow, Fg: termColorYellow }, { Bg: termColorYellow, Fg: termColorBlue }, { Bg: termColorYellow, Fg: termColorMagenta }, { Bg: termColorYellow, Fg: termColorCyan }, { Bg: termColorYellow, Fg: termColorWhite },
    { Bg: termColorBlue,    Fg: termColorBlack }, { Bg: termColorBlue, Fg: termColorRed }, { Bg: termColorBlue, Fg: termColorGreen }, { Bg: termColorBlue, Fg: termColorYellow }, { Bg: termColorBlue, Fg: termColorBlue }, { Bg: termColorBlue, Fg: termColorMagenta }, { Bg: termColorBlue, Fg: termColorCyan }, { Bg: termColorBlue, Fg: termColorWhite },
    { Bg: termColorMagenta, Fg: termColorBlack }, { Bg: termColorMagenta, Fg: termColorRed }, { Bg: termColorMagenta, Fg: termColorGreen }, { Bg: termColorMagenta, Fg: termColorYellow }, { Bg: termColorMagenta, Fg: termColorBlue }, { Bg: termColorMagenta, Fg: termColorMagenta }, { Bg: termColorMagenta, Fg: termColorCyan }, { Bg: termColorMagenta, Fg: termColorWhite },
    { Bg: termColorCyan,    Fg: termColorBlack }, { Bg: termColorCyan, Fg: termColorRed }, { Bg: termColorCyan, Fg: termColorGreen }, { Bg: termColorCyan, Fg: termColorYellow }, { Bg: termColorCyan, Fg: termColorBlue }, { Bg: termColorCyan, Fg: termColorMagenta }, { Bg: termColorCyan, Fg: termColorCyan }, { Bg: termColorCyan, Fg: termColorWhite },
    { Bg: termColorWhite,   Fg: termColorBlack }, { Bg: termColorWhite, Fg: termColorRed }, { Bg: termColorWhite, Fg: termColorGreen }, { Bg: termColorWhite, Fg: termColorYellow }, { Bg: termColorWhite, Fg: termColorBlue }, { Bg: termColorWhite, Fg: termColorMagenta }, { Bg: termColorWhite, Fg: termColorCyan }, { Bg: termColorWhite, Fg: termColorWhite },
  }

  height, width := Size()

  buff := make([][]Cell, height)
  for i := 0; i < height; i++ {
    buff[i] = make([]Cell, width)
  }

  stdscr := Window {
    Height: height,
    Width : width,
    Buffer: buff,
    BGChar: ' ',
    Looper: true,
    Echo  : true,
    Curs  : true,
    Delay : true,
    Resize: true,
  }

  return &stdscr, nil
}

func Close() {
  term.Close()
}

func Size() (height, width int ) {
  width, height = term.Size()
  return
}

func (w *Window) Size() (height, width int) {
  return w.Height, w.Width
}

func (w *Window) AddCh( ch uint64 ) {
  w.Touch = true

  cell := &w.Buffer[w.CurY][w.CurX]
  chAttrs, chColor, _, r := extractData( ch )

  if r != 0 && r != '\n' {
    cell.Ch = r

    if chAttrs == 0 { cell.Attrs = w.Attrs
    }  else         { cell.Attrs = chAttrs | w.Attrs }

    if chColor == 0 { cell.Color = w.Color
    }  else         { cell.Color = chColor }

    cell.Touch = true

    if w.Echo {
      printCell( cell, w.CurY, w.CurX )
      w.Touch = false
    }

  }

  w.mvCurs( r == '\n' )

  if w.Curs {
    term.SetCursor( w.CurX, w.CurY )
    term.Flush()
  } else {
    term.HideCursor()
    term.Flush()
  }
}

func (w *Window) AddCell( cell Cell ) {
  w.Touch = true

  cell.Touch = true
  w.Buffer[w.CurY][w.CurX] = cell

  if w.Echo {
    printCell( &cell, w.CurY, w.CurX )
    w.Touch = false
  }

  w.mvCurs( cell.Ch == '\n' )

  if w.Curs {
    term.SetCursor( w.CurX, w.CurY )
    term.Flush()
  } else {
    term.HideCursor()
    term.Flush()
  }
}

func (w *Window) mvCurs( nl bool ) {
  if w.Looper {
    w.mvLooper( nl )
    return
  }

  if w.Scroll {
    w.mvScroll( nl )
    return
  }

  if w.CurX == w.Width && w.CurY == w.Height { return }
}

func (w *Window) mvLooper( nl bool ) {
  if nl || w.CurX + 1 == w.Width {
    w.CurX = 0
  } else {
    w.CurX++
    return
  }

  if w.CurY + 1 < w.Height { w.CurY++
  } else                   { w.CurY = 0 }
}

func (w *Window) mvScroll( nl bool ) {
  if nl || w.CurX + 1 == w.Width {
    w.CurX = 0
  } else {
    w.CurX++
    return
  }

  if w.CurY + 1 < w.Height {
    w.CurY++
  } else {
    for y := 1; y < w.Height; y++ {
      for x := 0; x < w.Width; x++ {
        w.Buffer[y-1][x] = w.Buffer[y][x]
        w.Buffer[y-1][x].Touch = true
      }
    }

    cleanCell := Cell{ w.BGChar, w.Color, w.Attrs, true }
    for y := w.Height - 1; y < w.Height; y++ {
      for x := 0; x < w.Width; x++ {
        w.Buffer[y][x] = cleanCell
      }
    }

    w.Touch = true
    w.Refresh()
  }
}

func printCell( cell *Cell, y, x int ){
  tbg, tfg := cell.makeColors()
  term.SetCell( x, y, cell.Ch, tfg, tbg )
  term.Flush()
  cell.Touch  = false
}

func (w *Window) AddStr( str string ) {
  for _, c := range( str ) {
    w.AddCh( uint64(c) )
  }
}

func resize( buff [][]Cell, w, h int ) [][]Cell {
  // hs := 0

  hLen := len(buff)
  if hLen > h {
    buff = buff[:h]
  } else if hLen < h {
    hCap := cap(buff)

    if hCap > h {
      buff = buff[:h]
    } else if hCap < h {
      // hs = h - hCap
    }
  }

  //   // There is room to grow.  Extend the slice.
  //   z = x[:zlen]
  // } else {
  //   // There is insufficient space.  Allocate a new array.
  //   // Grow by doubling, for amortized linear complexity.
  //   zcap := zlen
  //   if zcap < 2*len(x) {
  //     zcap = 2 * len(x)
  //   }
  //   z = make([]int, zlen, zcap)
  //   copy(z, x) // a built-in function; see text
  // }
  // z[len(x)] = y
  // return z


  // newBuff :=

  return [][]Cell{}
}

func (w *Window) Refresh() {
  if w.Touch {
    for y, row := range( w.Buffer ) {
      for x, cell := range( row ) {
        if cell.Touch {
          tbg, tfg := cell.makeColors()
          term.SetCell( x, y, cell.Ch, tfg, tbg )
          w.Buffer[y][x].Touch = false
        }
      }
    }

    w.Touch = false
    term.Flush()
  }
}

func (w *Window) Clear(){
  defaultCell := Cell{ Ch: w.BGChar, Touch: true }

  for y := 0; y < len( w.Buffer ); y++ {
    for x := 0; x < len( w.Buffer[y] ); x++ {
      w.Buffer[y][x] = defaultCell
    }
  }

  w.Touch = true
}

// int addch(const chtype ch);
// int waddch(WINDOW *win, const chtype ch);
// int mvaddch(int y, int x, const chtype ch);
// int mvwaddch(WINDOW *win, int y, int x, const chtype ch);
// int addchstr(const chtype *chstr);
// int addchnstr(const chtype *chstr, int n);
// int waddchstr(WINDOW *win, const chtype *chstr);
// int waddchnstr(WINDOW *win, const chtype *chstr, int n);
// int mvaddchstr(int y, int x, const chtype *chstr);
// int mvaddchnstr(int y, int x, const chtype *chstr, int n);
// int mvwaddchstr(WINDOW *win, int y, int x, const chtype *chstr);
// int mvwaddchnstr(WINDOW *win, int y, int x, const chtype *chstr, int n);
//# int addstr(const char *str);
//# int addnstr(const char *str, int n);
//# int waddstr(WINDOW *win, const char *str);
//# int waddnstr(WINDOW *win, const char *str, int n);
// int mvaddstr(int y, int x, const char *str);
// int mvaddnstr(int y, int x, const char *str, int n);
// int mvwaddstr(WINDOW *win, int y, int x, const char *str);
// int mvwaddnstr(WINDOW *win, int y, int x, const char *str, int n);
// int assume_default_colors(int fg, int bg);
// int attr_get(attr_t *attrs, short *pair, void *opts);
// int wattr_get(WINDOW *win, =attr_t= *attrs, short *pair,void *opts);
// int wattroff(WINDOW *win, int attrs);
// int attr_off(attr_t attrs, void *opts);
// int wattr_off(WINDOW *win, attr_t attrs, void *opts);
// int wattron(WINDOW *win, int attrs);
func  WAttron( win *Window, attrs uint8 ){
  win.Attrs |= attrs
}

// int attron(int attrs);
func  Attron( attrs uint32 ){
//  StdScr.Attrs |= attrs
}

// int attr_on(attr_t attrs, void *opts);
// int wattr_on(WINDOW *win, attr_t attrs, void *opts);
// int wattrset(WINDOW *win, int attrs);
func  WAttrSet( win *Window, attrs uint8 ){
  win.Attrs = attrs
}

// int attrset(int attrs);
func  AttrSet( attrs uint32 ){
  //StdScr.Attrs = attrs
}

// int attr_set(attr_t attrs, short pair, void *opts);
// int wattr_set(WINDOW *win, attr_t attrs, short pair, void *opts);
// int baudrate(void);
// int beep(void);
// int bkgd(chtype ch);
func (w *Window) SetColor( color uint64 ){
  w.Attrs, w.Color, _, _ = extractData( color )
}

// int wbkgd(WINDOW *win, chtype ch);
// int bkgdset(chtype ch);
// void wbkgdset(WINDOW *win, chtype ch);
// int border(chtype ls, chtype rs, chtype ts, chtype bs, chtype tl, chtype tr, chtype bl, chtype br);
// int wborder(WINDOW *win, chtype ls, chtype rs, chtype ts, chtype bs, chtype tl, chtype tr, chtype bl, chtype br);
// int box(WINDOW *win, chtype verch, chtype horch);
// bool can_change_color(void);
// int cbreak(void);
// int nocbreak(void);
// int chgat(int n, attr_t attr, short color, const void *opts);
// int wchgat(WINDOW *win, int n, attr_t attr,short color, const void *opts);
// int mvchgat(int y, int x, int n, attr_t attr,short color, const void *opts);
// int mvwchgat(WINDOW *win, int y, int x, int n,attr_t attr, short color, const void *opts);
// int clear(void);
// int wclear(WINDOW *win);
// int clearok(WINDOW *win, bool bf);
// int clrtobot(void);
// int wclrtobot(WINDOW *win);
// int clrtoeol(void);
// int wclrtoeol(WINDOW *win);
func MakeColorPair( pair uint8 ) uint64 {
  return uint64(pair) << 48
}

func extractData( ch uint64 ) (attrs, color uint8, keyMod uint16, r rune) {
  color  = uint8((ch & hasColor) >> 48)
  attrs  = uint8(ch >> 56)
  keyMod = uint16((ch & modMask) >> 32)

  if (ch & hasKey) == 0 {
    r = rune(ch & runeMask)
  }

  return
}

func extractAttrs( ch uint64 ) uint64 {
  return ch & hasAttr
}

func (c *Cell) makeColors() (bg, fg term.Attribute) {
  color := Colors[ c.Color ]
  fg = term.Attribute(color.Fg) | term.Attribute(c.Attrs) << 8
  bg = term.Attribute(color.Bg)

  return
}

// int color_content(short color, short *r, short *g, short *b);
// int color_set(short color_pair_number, void* opts);
// int copywin(const WINDOW *srcwin, WINDOW *dstwin, int sminrow, int smincol, int dminrow, int dmincol, int dmaxrow, int dmaxcol, int overlay);
// int curs_set(int visibility);
func (w *Window) CursSet( visibility bool ){
  if visibility {
    term.SetCursor( w.CurX, w.CurY )
  } else {
    term.HideCursor()
  }

  w.Curs = visibility
  term.Flush()
}

// const char * curses_version(void);
// int delch(void);
// int wdelch(WINDOW *win);
// int mvdelch(int y, int x);
// int mvwdelch(WINDOW *win, int y, int x);
// int deleteln(void);
// int wdeleteln(WINDOW *win);
// void delscreen(SCREEN* sp);
// int delwin(WINDOW *win);
// WINDOW *derwin(WINDOW *orig, int nlines, int ncols, int begin_y, int begin_x);
// int doupdate(void);
// WINDOW *dupwin(WINDOW *win);
// int echo(void);
// int noecho(void);
// int echochar(const chtype ch);
// int wechochar(WINDOW *win, const chtype ch);
// int endwin(void);
// int erase(void);
// int werase(WINDOW *win);
// char erasechar(void);
// void filter(void);
// int flash(void);
// int flushinp(void);
// void getbegyx(WINDOW *win, int y, int x);
// chtype getbkgd(WINDOW *win);
// int getch(void);
func (w *Window) Getch() uint64 {
  for {
    event := term.PollEvent()
    if event.Type == term.EventKey {
      if event.Ch != 0 {
        if w.Echo { w.AddCh( uint64(event.Ch) ) }
        return uint64(event.Ch)
      }

      switch event.Key {
      case term.KeyF1:            return KeyF1
      case term.KeyF2:            return KeyF2
      case term.KeyF3:            return KeyF3
      case term.KeyF4:            return KeyF4
      case term.KeyF5:            return KeyF5
      case term.KeyF6:            return KeyF6
      case term.KeyF7:            return KeyF7
      case term.KeyF8:            return KeyF8
      case term.KeyF9:            return KeyF9
      case term.KeyF10:           return KeyF10
      case term.KeyF11:           return KeyF11
      case term.KeyF12:           return KeyF12
      case term.KeyInsert:        return KeyInsert
      case term.KeyDelete:        return KeyDelete
      case term.KeyHome:          return KeyHome
      case term.KeyEnd:           return KeyEnd
      case term.KeyPgup:          return KeyPgup
      case term.KeyPgdn:          return KeyPgdn
      case term.KeyArrowUp:       return KeyArrowUp
      case term.KeyArrowDown:     return KeyArrowDown
      case term.KeyArrowLeft:     return KeyArrowLeft
      case term.KeyArrowRight:    return KeyArrowRight
      case term.MouseLeft:        return Mouse
      case term.MouseMiddle:      return Mouse
      case term.MouseRight:       return Mouse
      case term.MouseRelease:     return Mouse
      case term.MouseWheelUp:     return Mouse
      case term.MouseWheelDown:   return Mouse
      case term.KeyCtrlSpace:     return ' ' | Ctrl
      case term.KeyCtrlA:         return 'a' | Ctrl
      case term.KeyCtrlB:         return 'b' | Ctrl
      case term.KeyCtrlC:         return 'c' | Ctrl
      case term.KeyCtrlD:         return 'd' | Ctrl
      case term.KeyCtrlE:         return 'e' | Ctrl
      case term.KeyCtrlF:         return 'f' | Ctrl
      case term.KeyCtrlG:         return 'g' | Ctrl
      case term.KeyCtrlH:         return 'h' | Ctrl
      case term.KeyCtrlI:         return 'i' | Ctrl
      case term.KeyCtrlJ:         return 'j' | Ctrl
      case term.KeyCtrlK:         return 'k' | Ctrl
      case term.KeyCtrlL:         return 'l' | Ctrl
      case term.KeyCtrlM:         return 'm' | Ctrl
      case term.KeyCtrlN:         return 'n' | Ctrl
      case term.KeyCtrlO:         return 'o' | Ctrl
      case term.KeyCtrlP:         return 'p' | Ctrl
      case term.KeyCtrlQ:         return 'q' | Ctrl
      case term.KeyCtrlR:         return 'r' | Ctrl
      case term.KeyCtrlS:         return 's' | Ctrl
      case term.KeyCtrlT:         return 't' | Ctrl
      case term.KeyCtrlU:         return 'u' | Ctrl
      case term.KeyCtrlV:         return 'v' | Ctrl
      case term.KeyCtrlW:         return 'w' | Ctrl
      case term.KeyCtrlX:         return 'x' | Ctrl
      case term.KeyCtrlY:         return 'y' | Ctrl
      case term.KeyCtrlZ:         return 'z' | Ctrl
      case term.KeyEsc:           return KeyEsc
      case term.KeyCtrlBackslash: return '\\' | Ctrl
      case term.KeyCtrl5:         return '5' | Ctrl
      case term.KeyCtrl6:         return '6' | Ctrl
      case term.KeyCtrl7:         return '7' | Ctrl
      case term.KeySpace:         return ' '
      case term.KeyBackspace2:
      }
    }
  }
}

// int wgetch(WINDOW *win);
// int mvgetch(int y, int x);
// int mvwgetch(WINDOW *win, int y, int x);
// void getmaxyx(WINDOW *win, int y, int x);
// int getmouse(MEVENT *event);
// void getparyx(WINDOW *win, int y, int x);
// int getstr(char *str);
// int getnstr(char *str, int n);
// int wgetstr(WINDOW *win, char *str);
// int wgetnstr(WINDOW *win, char *str, int n);
// int mvgetstr(int y, int x, char *str);
// int mvwgetstr(WINDOW *win, int y, int x, char *str);
// int mvgetnstr(int y, int x, char *str, int n);
// int mvwgetnstr(WINDOW *, int y, int x, char *str, int n);
// void getsyx(int y, int x);
// WINDOW *getwin(FILE *filep);
// void getyx(WINDOW *win, int y, int x);
// int halfdelay(int tenths);
// bool has_colors(void);
// bool has_ic(void);
// bool has_il(void);
// int hline(chtype ch, int n);
// void idcok(WINDOW *win, bool bf);
// void immedok(WINDOW *win, bool bf);
// chtype inch(void);
// chtype winch(WINDOW *win);
// chtype mvinch(int y, int x);
// chtype mvwinch(WINDOW *win, int y, int x);
// int inchstr(chtype *chstr);
// int inchnstr(chtype *chstr, int n);
// int winchstr(WINDOW *win, chtype *chstr);
// int winchnstr(WINDOW *win, chtype *chstr, int n);
// int mvinchstr(int y, int x, chtype *chstr);
// int mvinchnstr(int y, int x, chtype *chstr, int n);
// int mvwinchstr(WINDOW *win, int y, int x, chtype *chstr);
// int mvwinchnstr(WINDOW *win, int y, int x, chtype *chstr, int n);
// int init_color(short color, short r, short g, short b);
// int init_pair(short pair, short f, short b);
// WINDOW *initscr(void);
// int insch(chtype ch);
// int winsch(WINDOW *win, chtype ch);
// int mvinsch(int y, int x, chtype ch);
// int mvwinsch(WINDOW *win, int y, int x, chtype ch);
// int insdelln(int n);
// int winsdelln(WINDOW *win, int n);
// int insertln(void);
// int winsdelln(WINDOW *win, int n);
// int insstr(const char *str);
// int insnstr(const char *str, int n);
// int winsstr(WINDOW *win, const char *str);
// int winsnstr(WINDOW *win, const char *str, int n);
// int mvinsstr(int y, int x, const char *str);
// int mvinsnstr(int y, int x, const char *str, int n);
// int mvwinsstr(WINDOW *win, int y, int x, const char *str);
// int mvwinsnstr(WINDOW *win, int y, int x, const char *str, int n);
// int instr(char *str);
// int innstr(char *str, int n);
// int winstr(WINDOW *win, char *str);
// int winnstr(WINDOW *win, char *str, int n);
// int mvinstr(int y, int x, char *str);
// int mvinnstr(int y, int x, char *str, int n);
// int mvwinstr(WINDOW *win, int y, int x, char *str);
// int mvwinnstr(WINDOW *win, int y, int x, char *str, int n);
// int intrflush(WINDOW *win, bool bf);
// bool isendwin(void);
// bool is_linetouched(WINDOW *win, int line);
// bool is_wintouched(WINDOW *win);
// char *keyname(int c);
// char *key_name(wchar_t w);
// int keypad(WINDOW *win, bool bf);
// char killchar(void);
// int leaveok(WINDOW *win, bool bf);
// Not applicable.
// char *termname(void);
// int meta(WINDOW *win, bool bf);
// Not applicable.
// bool mouse_trafo(int* pY, int* pX, bool to_screen);
// bool wmouse_trafo(const WINDOW* win, int* pY, int* pX, bool to_screen);
// mmask _t mousemask(mmask_t newmask, mmask_t *oldmask);
// int move(int y, int x);
func Move( y, x int ){
  term.SetCursor( x, y )
}

// int wmove(WINDOW *win, int y, int x);
// int mvderwin(WINDOW *win, int par_y, int par_x);
// int mvwin(WINDOW *win, int y, int x);
// int napms(int ms);
func Napms( ms uint ){
  var t time.Duration
  t = time.Duration( ms ) * time.Millisecond
  time.Sleep( t )
}


// Not applicable.
// WINDOW *newpad(int nlines, int ncols);
// SCREEN *newterm(char *type, FILE *outfd, FILE *infd);
// WINDOW *newwin(int nlines, int ncols, int begin_y, int begin_x);
// int nl(void);
// int nonl(void);
// int nodelay(WINDOW *win, bool bf);
// int notimeout(WINDOW *win, bool bf);
// int overlay(const WINDOW *srcwin, WINDOW *dstwin);
// int overwrite(const WINDOW *srcwin, WINDOW *dstwin);
// int pair_content(short pair, short *f, short *b);
// int pechochar(WINDOW *pad, chtype ch);
// int pnoutrefresh(WINDOW *pad, int pminrow, int pmincol, int sminrow, int smincol, int smaxrow, int smaxcol);
// int prefresh(WINDOW *pad, int pminrow, int pmincol, int sminrow, int smincol, int smaxrow, int smaxcol);
// int printw(const char *fmt, ...);
// int wprintw(WINDOW *win, const char *fmt, ...);
// int mvprintw(int y, int x, const char *fmt, ...);
// int mvwprintw(WINDOW *win, int y, int x, const char *fmt, ...);
// int vw_printw(WINDOW *win, const char *fmt, va_list varglist);
// int vwprintw(WINDOW *win, const char *fmt, va_list varglist);
// int putwin(WINDOW *win, FILE *filep);
// void qiflush(void);
// void noqiflush(void);
// int raw(void);
// int noraw(void);
// int redrawwin(WINDOW *win);
// int refresh(void);
// int wrefresh(WINDOW *win);
// int ripoffline(int line, int (*init)(WINDOW *, int));
// int scanw(char *fmt, ...);
// int wscanw(WINDOW *win, char *fmt, ...);
// int mvscanw(int y, int x, char *fmt, ...);
// int mvwscanw(WINDOW *win, int y, int x, char *fmt, ...);
// int vw_scanw(WINDOW *win, char *fmt, va_list varglist);
// int vwscanw(WINDOW *win, char *fmt, va_list varglist);
// int scr_dump(const char *filename);
// int scr_init(const char *filename);
// int scr_restore(const char *filename);
// int scr_set(const char *filename);
// int scrl(int n);
// int wscrl(WINDOW *win, int n);
// int scroll(WINDOW *win);
// int scrollok(WINDOW *win, bool bf);
// int setscrreg(int top, int bot);
// int wsetscrreg(WINDOW *win, int top, int bot);
// void setsyx(int y, int x);
// SCREEN *set_term(SCREEN *new);
// attr _t slk_attr(void);
// int slk_attroff(const chtype attrs);
// int slk_attr_off(const attr_t attrs, void * opts);
// int slk_attron(const chtype attrs);
// int slk_attr_on(attr_t attrs, void* opts);
// int slk_attrset(const chtype attrs);
// int slk_attr_set(const attr_t attrs, short color_pair_number, ; void* opts);
// int slk_clear(void);
// int slk_color(short color_pair_number);
// int slk_init(int fmt);
// char *slk_label(int labnum);
// int slk_noutrefresh(void);
// int slk_refresh(void);
// int slk_restore(void);
// int slk_set(int labnum, const char *label, int fmt);
// int slk_touch(void);
// int standend(void);
// int wstandend(WINDOW *win);
// int standout(void);
// int wstandout(WINDOW *win);
// int start_color(void);
// WINDOW *subpad(WINDOW *orig, int nlines, int ncols,int begin_y, int begin_x);
// WINDOW *subwin(WINDOW *orig, int nlines, int ncols, int begin_y, int begin_x);
// int syncok(WINDOW *win, bool bf);
// chtype termattrs(void);
// attr _t term_attrs(void);
// char *termname(void);
// void timeout(int delay);
// void wtimeout(WINDOW *win, int delay);
// int touchline(WINDOW *win, int start, int count);
// int touchwin(WINDOW *win);
// int typeahead(int fd);
// char *unctrl(chtype c);
// int ungetch(int ch);
// int untouchwin(WINDOW *win);
// int use_default_colors(void);
// void use_env(bool f);
// int vline(chtype ch, int n);
// int wvline(WINDOW *win, chtype ch, int n);
// int mvvline(int y, int x, chtype ch, int n);
// int mvwvline(WINDOW *, int y, int x, chtype ch, int n);
// void wcursyncup(WINDOW *win);
// bool wenclose(const WINDOW *win, int y, int x);
// int wnoutrefresh(WINDOW *win);
// int wredrawln(WINDOW *win, int beg_line, int num_lines);
// void wsyncdown(WINDOW *win);
// void wsyncup(WINDOW *win);
// int wtouchln(WINDOW *win, int y, int n, int changed);
